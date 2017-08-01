package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cenkalti/backoff"
	"gopkg.in/yaml.v2"
)

// DatadogConnector performs a connection to Datadog.
type DatadogConnector struct {
	// The API Key for Datadog.
	apiKey string
	// The APP Key for Datadog.
	appKey string
	// HTTPClient in order to talk to DD.
	HTTPClient *http.Client
	// RetryTimeout specifies the retry timeout
	RetryTimeout time.Duration
}

type validationResponse struct {
	Errors  []string `json:"errors"`
	IsValid bool     `json:"valid"`
}

// Dashboard NOTE this doesn't contian all fields for a dashboard, just the fields
// We care about for HTTP.
type Dashboard struct {
	ID    *string `json:"id,omitempty"`
	Title *string `json:"title,omitempty"`
}

// Screenboard NOTE this doesn't contian all fields for a dashboard, just the fields
// We care about for HTTP.
type Screenboard struct {
	ID    *int    `json:"id,omitempty"`
	Title *string `json:"title,omitempty"`
}

// CreateDashboardResp is a response from CreateDashboard
type CreateDashboardResp struct {
	Resource  *string    `json:"resource,omitempty"`
	URL       *string    `json:"url,omitempty"`
	Dashboard *Dashboard `json:"dash,omitempty"`
}

// DashboardListResp is a list of Dashboards.
type DashboardListResp struct {
	Dashboards []Dashboard `json:"dashes,omitempty"`
}

// ScreensListResp is a list of Screenboards.
type ScreensListResp struct {
	Dashboards []Screenboard `json:"screenboards,omitempty"`
}

// NewDatadogConnector creates a new datadog connector.
func NewDatadogConnector(apiKey string, appKey string, timeoutSeconds int) *DatadogConnector {
	return &DatadogConnector{
		apiKey,
		appKey,
		&http.Client{
			Timeout: time.Duration(timeoutSeconds) * time.Second,
		},
		time.Duration(timeoutSeconds*5) * time.Second,
	}
}

// urlForApi grabs a url for a specific API Path.
func (client *DatadogConnector) uriForAPI(api string) string {
	url := os.Getenv("DATADOG_HOST")
	if url == "" {
		url = "https://app.datadoghq.com"
	}
	if strings.Index(api, "?") > -1 {
		return url + "/api" + api + "&api_key=" +
			client.apiKey + "&application_key=" + client.appKey
	}
	return url + "/api" + api + "?api_key=" +
		client.apiKey + "&application_key=" + client.appKey
}

// Validate checks if the API and application keys are valid.
func (client *DatadogConnector) Validate() (bool, error) {
	var bodyreader io.Reader
	var out validationResponse
	req, err := http.NewRequest("GET", client.uriForAPI("/v1/validate"), bodyreader)

	if err != nil {
		return false, err
	}
	if bodyreader != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	var resp *http.Response
	resp, err = client.HTTPClient.Do(req)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	// Only care about 200 OK or 403 which we'll unmarshal into struct valid.
	// Everything else is of no interest to us.
	if resp.StatusCode != 200 && resp.StatusCode != 403 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return false, err
		}
		return false, fmt.Errorf("API error %s: %s", resp.Status, body)
	}

	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &out)
	if err != nil {
		return false, err
	}

	return out.IsValid, nil
}

// doRequestWithRetries performs an HTTP request repeatedly for maxTime or until
// no error and no acceptable HTTP response code was returned.
func (client *DatadogConnector) doRequestWithRetries(req *http.Request, maxTime time.Duration) (*http.Response, error) {
	var (
		err  error
		resp *http.Response
		bo   = backoff.NewExponentialBackOff()
		body []byte
	)

	bo.MaxElapsedTime = maxTime

	// Save the body for retries
	if req.Body != nil {
		body, err = ioutil.ReadAll(req.Body)
		if err != nil {
			return resp, err
		}
	}

	operation := func() error {
		if body != nil {
			r := bytes.NewReader(body)
			req.Body = ioutil.NopCloser(r)
		}

		resp, err = client.HTTPClient.Do(req)
		if err != nil {
			return err
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			// 2xx all done
			return nil
		} else if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			// 4xx are not retryable
			return nil
		}

		return fmt.Errorf("Received HTTP status code %d", resp.StatusCode)
	}

	err = backoff.Retry(operation, bo)

	return resp, err
}

// DoJSONRequest is the simplest type of request: a method on a URI that returns
// some JSON result which we unmarshal into the passed interface.
func (client *DatadogConnector) DoJSONRequest(method string, api string, reqbody, out interface{}) error {
	var bodyreader io.Reader
	if method != "GET" && reqbody != nil {
		bjson, err := json.Marshal(reqbody)
		if err != nil {
			return err
		}
		bodyreader = bytes.NewReader(bjson)
	}

	req, err := http.NewRequest(method, client.uriForAPI(api), bodyreader)
	if err != nil {
		return err
	}
	if bodyreader != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	// Perform the request and retry it if it's not a POST or PUT or DELETE request
	var resp *http.Response
	if method == "POST" || method == "PUT" || method == "DELETE" {
		resp, err = client.HTTPClient.Do(req)
	} else {
		resp, err = client.doRequestWithRetries(req, client.RetryTimeout)
	}
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("API error %s: %s", resp.Status, body)
	}

	// If they don't care about the body, then we don't care to give them one,
	// so bail out because we're done.
	if out == nil {
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// If we got no body, by default let's just make an empty JSON dict.
	if len(body) == 0 {
		body = []byte{'{', '}'}
	}

	if err := json.Unmarshal(body, &out); err != nil {
		return err
	}
	return nil
}

// DryRunDash creates, and destroys a whole bunch of example dashboards.
func (client *DatadogConnector) DryRunDash(fs *FileSystem) error {
	docs, err := fs.GetTemplates()
	if err != nil {
		return err
	}
	for _, doc := range docs {
		var out CreateDashboardResp
		if err := client.DoJSONRequest("POST", "/v1/dash", doc, &out); err != nil {
			return err
		}
		if out.Dashboard == nil || out.Dashboard.ID == nil {
			return fmt.Errorf("Response from datadog had no valid dashboard: %+v", out)
		}
		if err := client.DoJSONRequest("DELETE", fmt.Sprintf("/v1/dash/%s", *out.Dashboard.ID), nil, nil); err != nil {
			return err
		}
	}
	return nil
}

// DryRunScreen creates, and destroys a whole bunch of example screens.
func (client *DatadogConnector) DryRunScreen(fs *FileSystem) error {
	docs, err := fs.GetTemplates()
	if err != nil {
		return err
	}
	for _, doc := range docs {
		var out CreateDashboardResp
		if err := client.DoJSONRequest("POST", "/v1/screen", doc, &out); err != nil {
			return err
		}
		if out.Dashboard == nil || out.Dashboard.ID == nil {
			return fmt.Errorf("Response from datadog had no valid screen: %+v", out)
		}
		if err := client.DoJSONRequest("DELETE", fmt.Sprintf("/v1/screen/%s", *out.Dashboard.ID), nil, nil); err != nil {
			return err
		}
	}
	return nil
}

// Fings a matching dashboard based on titles.
func findDashboard(title string, dashboards []Dashboard) string {
	for _, dash := range dashboards {
		if dash.Title == nil {
			continue
		}
		if dash.Title == &title {
			return *dash.ID
		}
	}
	return "-1"
}

// Fings a matching dashboard based on titles.
func findScreenboard(title string, dashboards []Screenboard) int {
	for _, dash := range dashboards {
		if dash.Title == nil {
			continue
		}
		if dash.Title == &title {
			return *dash.ID
		}
	}
	return -1
}

// getDashAsMap gets a dashboard as a map[string]interface{} instead of map[interface{}]interface{}
func getDashAsMap(dash map[string]interface{}) map[string]interface{} {
	dashFrd := dash["dash"]
	if dashFrd == nil {
		return nil
	}
	newMap := make(map[string]interface{})
	if actualDash, okay := dashFrd.(map[interface{}]interface{}); okay {
		for k, v := range actualDash {
			if actualKey, ok := k.(string); ok {
				newMap[actualKey] = v
			}
		}
	}
	return newMap
}

// CreateDashboards actually runs, and creates all the dashboards.
func (client *DatadogConnector) CreateDashboards(fs *FileSystem) error {
	docs, err := fs.GetTemplates()
	if err != nil {
		return err
	}
	for _, doc := range docs {
		var out DashboardListResp
		if err := client.DoJSONRequest("GET", "/v1/dash", nil, &out); err != nil {
			return err
		}
		dashFrd := getDashAsMap(doc)
		if dashFrd == nil {
			continue
		}
		if val := findDashboard(dashFrd["title"].(string), out.Dashboards); val != "-1" {
			if err = client.DoJSONRequest("DELETE", fmt.Sprintf("/v1/dash/%s", val), nil, nil); err != nil {
				return err
			}
		}
		marshaled, err := yaml.Marshal(&dashFrd)
		if err != nil {
			return err
		}
		asBytes, err := yamlToJSON(marshaled, nil)
		if err != nil {
			return err
		}
		var fromJSON map[string]interface{}
		err = json.Unmarshal(asBytes, &fromJSON)
		if err != nil {
			return err
		}
		if err = client.DoJSONRequest("POST", "/v1/dash", fromJSON, nil); err != nil {
			return err
		}
	}
	return nil
}

// CreateScreens actually runs, and creates all the screens.
func (client *DatadogConnector) CreateScreens(fs *FileSystem) error {
	docs, err := fs.GetTemplates()
	if err != nil {
		return err
	}
	for _, doc := range docs {
		var out ScreensListResp
		if err := client.DoJSONRequest("GET", "/v1/screen", nil, &out); err != nil {
			return err
		}
		if val := findScreenboard(doc["board_title"].(string), out.Dashboards); val != -1 {
			if err = client.DoJSONRequest("DELETE", fmt.Sprintf("/v1/screen/%d", val), nil, nil); err != nil {
				return err
			}
		}
		marshaled, err := yaml.Marshal(&doc)
		if err != nil {
			return err
		}
		asBytes, err := yamlToJSON(marshaled, nil)
		if err != nil {
			return err
		}
		var fromJSON map[string]interface{}
		err = json.Unmarshal(asBytes, &fromJSON)
		if err != nil {
			return err
		}
		if err = client.DoJSONRequest("POST", "/v1/screen", fromJSON, nil); err != nil {
			return err
		}
	}
	return nil
}
