package api

import (
	"fmt"
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

	// go maps are unordered, deal with this without having to create a sort function
	var1 := 0
	var2 := 1
	if res := heritage.Environment.Entries[0].Name; res == "TESTVAR2" {
		var1 = 1
		var2 = 0
	}

	if res := heritage.Environment.Entries[var1].Name; res != "TESTVAR" {
		t.Errorf("Expected 'TESTVAR' but got: %s", res)
	}

	if res := *heritage.Environment.Entries[var1].SsmPath; res != "hello/testvar" {
		t.Errorf("Expected 'hello/testvar' but got: %s", res)
	}

	if res := heritage.Environment.Entries[var2].Name; res != "TESTVAR2" {
		t.Errorf("Expected 'TESTVAR2' but got: %s", res)
	}

	if res := *heritage.Environment.Entries[var2].SsmPath; res != "hello/testvar2" {
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

	// go maps are unordered, deal with this without having to create a sort function
	var1 := 0
	var2 := 1
	if res := heritage.Environment.Entries[0].Name; res == "TESTVAR1" {
		var1 = 1
		var2 = 0
	}

	if res := heritage.Environment.Entries[var1].Name; res != "TESTVAR" {
		t.Errorf("Expected 'TESTVAR' but got: %s", res)
	}

	if res := *heritage.Environment.Entries[var1].SsmPath; res != "hello/testvar1" {
		t.Errorf("Expected 'hello/testvar1' but got: %s", res)
	}

	if res := heritage.Environment.Entries[var2].Name; res != "TESTVAR2" {
		t.Errorf("Expected 'TESTVAR2' but got: %s", res)
	}

	if res := *heritage.Environment.Entries[var2].SsmPath; res != "hello/testvar2" {
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

	s1 := `{"name":"TESTVAR","ssm_path":"hello/testvar1"}`
	s2 := `{"name":"TESTVAR2","ssm_path":"hello/testvar2"}`
	expectation1 := fmt.Sprintf("[%s,%s]", s1, s2)
	expectation2 := fmt.Sprintf("[%s,%s]", s2, s1)

	bytes, _ := json.Marshal(heritage.Environment)
	if str := string(bytes[:]); str != expectation1 && str != expectation2 {
		t.Errorf("Expected \n'%s' \nor \n'%s' \nbut got: \n'%s'", expectation1, expectation2, str)
	}
}
