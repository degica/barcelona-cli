package cmd

import (
	"fmt"
	"os"

	"github.com/jarcoal/httpmock"
)

func Example_secret_add() {
	pwd, _ := os.Getwd()

	HeritageConfigFilePath = pwd + "/test/test-barcerola.yml"

	app := newTestApp(SecretCommand)
	endpoint := os.Getenv("BARCELONA_ENDPOINT")

	testArgs := []string{
		"bcn", "secret", "add", "-n", "bcn-test",
		"-v", "test", "-d", "staging"}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", endpoint+"/v1/districts/staging/ssm_parameters",
		httpmock.NewStringResponder(200, ""))

	app.Run(testArgs)

	// Output:
	// success to set bcn-test
}

func Example_secret_delete() {
	pwd, _ := os.Getwd()

	HeritageConfigFilePath = pwd + "/test/test-barcerola.yml"

	app := newTestApp(SecretCommand)
	endpoint := os.Getenv("BARCELONA_ENDPOINT")

	testArgs := []string{
		"bcn", "secret", "delete", "-n", "bcn-test", "-d", "staging"}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	resJson, _ := readJsonResponse(pwd + "/test/secret_remove.json")
	httpmock.RegisterResponder("DELETE", endpoint+"/v1/districts/staging/ssm_parameters/bcn-test",
		httpmock.NewStringResponder(200, resJson))

	app.Run(testArgs)

	// Output:
	// {"deleted_parameters":["bcn-tests"],"invalid_parameters":[]}
}

func Example_secret_delete_with_slash_name() {

	pwd, _ := os.Getwd()

	HeritageConfigFilePath = pwd + "/test/test-barcerola.yml"
	app := newTestApp(SecretCommand)
	endpoint := os.Getenv("BARCELONA_ENDPOINT")

	testArgs := []string{
		"bcn", "secret", "delete", "-n", "hello/test", "-d", "staging"}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	resJson, _ := readJsonResponse(pwd + "/test/secret_remove.json")

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("%s/v1/districts/staging/ssm_parameters/hello%%2Ftest", endpoint),
		httpmock.NewStringResponder(200, resJson))

	app.Run(testArgs)
}
