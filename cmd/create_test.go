package cmd

import (
	"os"

	"github.com/jarcoal/httpmock"
)

func Example_create_heritage() {
	app := newTestApp(CreateCommand)
	app.Writer = os.Stdout
	testArgs := []string{"bcn", "create", "-e", "test"}
	pwd, _ := os.Getwd()
	endpoint := os.Getenv("BARCELONA_ENDPOINT")

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", endpoint+"/v1/auth/github/login",
		httpmock.NewStringResponder(200, "{}"))

	resJson, _ := readJsonResponse(pwd + "/test/create_heritage.json")
	httpmock.RegisterResponder("POST", endpoint+"/v1/districts/default/heritages",
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
