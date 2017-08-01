package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/spf13/afero"
)

var dryRun = flag.Bool("dry-run", false, "Whether or not to run a Dry Run.")

func main() {
	fmt.Println("Starting Greyhound...")

	fmt.Println("Creating Datadog Client...")
	ddConnector := NewDatadogConnector(os.Getenv("DATADOG_API_KEY"), os.Getenv("DATADOG_APP_KEY"), 10)
	isValid, err := ddConnector.Validate()
	if err != nil {
		fmt.Printf("Failed to query datadog: %v\n", err)
		os.Exit(1)
	}
	if !isValid {
		fmt.Printf("Datadog Credentials aren't valid\n")
		os.Exit(1)
	}

	fmt.Println("Creating FileSystem client for Dashboards...")
	fs, err := CreateFileSystem(os.Getenv("GREYDOG_DASH_PATH"), os.Getenv("GREYDOG_CACHE_DASH_PATH"), afero.NewOsFs())
	if err != nil {
		fmt.Printf("Failed to Create FileSystem for Dashs: %v\n", err)
		os.Exit(1)
	}

	fsScreen, err := CreateFileSystem(os.Getenv("GREYDOG_SCREEN_PATH"), os.Getenv("GREYDOG_CACHE_SCREEN_PATH"), afero.NewOsFs())
	if err != nil {
		fmt.Printf("Failed to Create FileSystem for Screens: %v\n", err)
		os.Exit(1)
	}

	if *dryRun {
		fmt.Println("Running a Dry run of Dashboards.")
		err = ddConnector.DryRunDash(fs)
		if err != nil {
			fmt.Println("Ran into an error on dry run dash!")
			fmt.Print(err)
			os.Exit(1)
		} else {
			fmt.Println("Successful!")
		}
		fmt.Println("Running a Dry run of Screens")
		err = ddConnector.DryRunScreen(fsScreen)
		if err != nil {
			fmt.Println("Ran into an error on dry run screen!")
			fmt.Print(err)
			os.Exit(1)
		} else {
			fmt.Println("Successful!")
		}
	} else {
		fmt.Println("Creating Dashboards...")
		err = ddConnector.CreateDashboards(fs)
		if err != nil {
			fmt.Println("Ran into an error Creating Dashboards!")
			fmt.Print(err)
			os.Exit(1)
		} else {
			fmt.Println("Successful!")
		}
		fmt.Println("Creating Screenboareds...")
		err = ddConnector.CreateScreens(fsScreen)
		if err != nil {
			fmt.Println("Ran into an error Creating Screens!")
			fmt.Print(err)
			os.Exit(1)
		} else {
			fmt.Println("Successful.")
		}
	}
}
