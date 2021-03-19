package cmd

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/urfave/cli"
)

func TestCheckEnvVars(t *testing.T) {
	result, err := checkEnvVars([]string{"ABC=def", "GHI=jkl"})

	if err != nil {
		t.Errorf("Expected there to be no error but got: %s", err)
	}

	if result["ABC"] != "def" {
		t.Errorf("Expected ABC to be set to def but got: %s", result)
	}

	if result["GHI"] != "jkl" {
		t.Errorf("Expected ABC to be set to def but got: %s", result)
	}
}

func TestCheckEnvVarsError(t *testing.T) {
	result, err := checkEnvVars([]string{"ABCd=def", "GHI=jkl"})

	if err == nil {
		t.Errorf("Expected to be an error but was nil")
	}

	if result != nil {
		t.Errorf("Expected result to be nil, but was %s", result)
	}
}

func TestRunByBranchName(t *testing.T) {
	pwd, _ := os.Getwd()

	HeritageConfigFilePath = pwd + "/test/test-barcerola.yml"

	app := newTestApp(RunCommand)
	endpoint := os.Getenv("BARCELONA_ENDPOINT")

	testArgs := []string{"bcn", "run", "-b", "test-branch", "--D", "bash"}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	resJson, err := readJsonResponse(pwd + "/test/review_group.json")
	if err != nil {
		t.Fatalf("error running command `loading json %v", err)
	}

	httpmock.RegisterResponder("GET", endpoint+"/v1/review_groups/test1/apps",
		httpmock.NewStringResponder(200, resJson))

	resJson, err = readJsonResponse(pwd + "/test/oneoffs.json")
	if err != nil {
		t.Fatalf("error running command `loading json %v", err)
	}

	httpmock.RegisterResponder("POST", endpoint+"/v1/heritages/review-heritage/oneoffs",
		httpmock.NewStringResponder(200, resJson))

	app.Run(testArgs)
}

func TestWrongFlag(t *testing.T) {
	app := newTestApp(RunCommand)

	testArgs := []string{"bcn", "run", "-B", "not-found", "--D", "bash"}

	err := app.Run(testArgs)

	if err.Error() != "flag provided but not defined: -B" {
		t.Fatalf("error running command `bcn run %v", err)
	}
}

func newTestApp(command cli.Command) *cli.App {
	a := cli.NewApp()
	a.Name = "bcn"
	a.Writer = ioutil.Discard
	a.Commands = []cli.Command{command}

	return a
}

func readJsonResponse(path string) (string, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return "", err
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	defer jsonFile.Close()

	return string(byteValue), err
}
