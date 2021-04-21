package cmd

import (
	"fmt"
	"os"

	"github.com/jarcoal/httpmock"
)

func Example_create_heritage_quiet_false() {
	app := newTestApp(CreateCommand)
	app.Writer = os.Stdout
	testArgs := []string{"bcn", "create", "-e", "test", "--quiet=false"}
	pwd, _ := os.Getwd()

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://vault.degica.com/v1/auth/github/login",
		httpmock.NewStringResponder(200, "{}"))

	resJson, err := readJsonResponse(pwd + "/test/create_heritage.json")
	if err != nil {
		fmt.Errorf("error running command `loading json %s", err)
	}
	httpmock.RegisterResponder("POST", "https://barcelona.degica.com/v1/districts/default/heritages",
		httpmock.NewStringResponder(200, resJson))

	app.Run(testArgs)

	// Output:
	// Name:          nginx
	// Image Name:    nginx
	// Image Tag :    latest
	// Version:       1
	// Before Deploy: echo hello
	// Token:         560d9e10-70ce-4f47-82b8-37d47761116d
	// Scheduled Tasks:
	// rate(1 minute)       echo hello
	// Environment Variables
}

func Example_create_heritage_quiet_true() {
	app := newTestApp(CreateCommand)
	testArgs := []string{"bcn", "create", "-e", "test"}
	pwd, _ := os.Getwd()
	app.Writer = os.Stdout

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://vault.degica.com/v1/auth/github/login",
		httpmock.NewStringResponder(200, "{}"))

	resJson, err := readJsonResponse(pwd + "/test/create_heritage.json")
	if err != nil {
		fmt.Errorf("error running command `loading json %s", err)
	}
	httpmock.RegisterResponder("POST", "https://barcelona.degica.com/v1/districts/default/heritages",
		httpmock.NewStringResponder(200, resJson))

	app.Run(testArgs)

	// Output:
}
