package cmd

import (
	"testing"
)

func TestCheckEnvVars(t *testing.T) {
	result, err := checkEnvVars([]string {"ABC=def", "GHI=jkl"})

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
	result, err := checkEnvVars([]string {"ABCd=def", "GHI=jkl"})

	if err == nil {
		t.Errorf("Expected to be an error but was nil")
	}

	if result != nil {
		t.Errorf("Expected result to be nil, but was %s", result)
	}

}
