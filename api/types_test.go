package api

import (
	"testing"

	"encoding/json"
	yaml "gopkg.in/yaml.v2"
)

func TestError(t *testing.T) {
	var err = APIError{"errorname", "foobar", []string{}}
	var result = err.Error()
	if result != "errorname" {
		t.Errorf("Expected 'errorname' but got: %s", result)
	}
}

func TestUnmarshalEnvironmentArray(t *testing.T) {
	str := `
    name: hello
    image_name: quay.io/dsiawdegica/hello
    scheduled_tasks: []
    environment:
      - name: TESTVAR
        ssm_path: hello/testvar
      - name: TESTVAR2
        ssm_path: hello/testvar2
  `

	var heritage Heritage
	err := yaml.Unmarshal([]byte(str), &heritage)
	if err != nil {
		t.Errorf("Could not Unmarshal heritage %s", err.Error())
	}

	if res := heritage.Environment.Entries[0].Name; res != "TESTVAR" {
		t.Errorf("Expected 'TESTVAR' but got: %s", res)
	}

	if res := *heritage.Environment.Entries[0].SsmPath; res != "hello/testvar" {
		t.Errorf("Expected 'hello/testvar' but got: %s", res)
	}

	if res := heritage.Environment.Entries[1].Name; res != "TESTVAR2" {
		t.Errorf("Expected 'TESTVAR2' but got: %s", res)
	}

	if res := *heritage.Environment.Entries[1].SsmPath; res != "hello/testvar2" {
		t.Errorf("Expected 'hello/testvar2' but got: %s", res)
	}
}

func TestUnmarshalEnvironmentHash(t *testing.T) {
	str := `
    name: hello
    image_name: quay.io/dsiawdegica/hello
    scheduled_tasks: []
    environment:
      TESTVAR:
        ssm_path: hello/testvar1
      TESTVAR2:
        ssm_path: hello/testvar2
  `

	var heritage Heritage
	err := yaml.Unmarshal([]byte(str), &heritage)
	if err != nil {
		t.Errorf("Could not Unmarshal heritage %s", err.Error())
	}

	if res := heritage.Environment.Entries[0].Name; res != "TESTVAR" {
		t.Errorf("Expected 'TESTVAR' but got: %s", res)
	}

	if res := *heritage.Environment.Entries[0].SsmPath; res != "hello/testvar1" {
		t.Errorf("Expected 'hello/testvar1' but got: %s", res)
	}

	if res := heritage.Environment.Entries[1].Name; res != "TESTVAR2" {
		t.Errorf("Expected 'TESTVAR2' but got: %s", res)
	}

	if res := *heritage.Environment.Entries[1].SsmPath; res != "hello/testvar2" {
		t.Errorf("Expected 'hello/testvar2' but got: %s", res)
	}
}

func TestMarshalEnvironmentHash(t *testing.T) {
	str := `
    name: hello
    image_name: quay.io/dsiawdegica/hello
    scheduled_tasks: []
    environment:
      TESTVAR:
        ssm_path: hello/testvar1
      TESTVAR2:
        ssm_path: hello/testvar2
  `

	var heritage Heritage
	err := yaml.Unmarshal([]byte(str), &heritage)
	if err != nil {
		t.Errorf("Could not Unmarshal heritage %s", err.Error())
	}

	expectation := `[{"name":"TESTVAR","ssm_path":"hello/testvar1"},{"name":"TESTVAR2","ssm_path":"hello/testvar2"}]`

	bytes, _ := json.Marshal(heritage.Environment)
	if str := string(bytes[:]); str != expectation {
		t.Errorf("Expected '%s' but got: %s", expectation, str)
	}
}
