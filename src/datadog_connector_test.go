package main

import (
	"os"
	"testing"

	gock "gopkg.in/h2non/gock.v1"
)

func TestValidate(t *testing.T) {
	t.Run("Invalid Validate Test", func(t *testing.T) {
		defer gock.Off()

		respData := make(map[string]interface{})
		respData["errors"] = []string{"test"}
		respData["valid"] = false

		if value, ok := os.LookupEnv("DATADOG_HOST"); ok {
			gock.New(value).
				Get("/api/v1/validate").
				Reply(200).
				JSON(respData)
		} else {
			gock.New("https://app.datadoghq.com").
				Get("/api/v1/validate").
				Reply(200).
				JSON(respData)
		}

		connector := NewDatadogConnector("test", "test", 3)
		valid, err := connector.Validate()

		if err != nil {
			t.Fatalf("Invalid Validate Test should be able to read HTTP: Error: %v", err)
		}
		if valid == true {
			t.Fatalf("Invalid Validate Test returned that it was Valid, WHUT?!?")
		}
	})

	t.Run("Valid Validate Test", func(t *testing.T) {
		defer gock.Off()

		respData := make(map[string]interface{})
		respData["errors"] = nil
		respData["valid"] = true

		if value, ok := os.LookupEnv("DATADOG_HOST"); ok {
			gock.New(value).
				Get("/api/v1/validate").
				Reply(200).
				JSON(respData)
		} else {
			gock.New("https://app.datadoghq.com").
				Get("/api/v1/validate").
				Reply(200).
				JSON(respData)
		}

		connector := NewDatadogConnector("test", "test", 3)
		valid, err := connector.Validate()

		if err != nil {
			t.Fatalf("Invalid Validate Test should be able to read HTTP: Error: %v", err)
		}
		if valid == false {
			t.Fatalf("Valid Validate Test returned that it wasn't Valid, WHUT?!?")
		}
	})
}
