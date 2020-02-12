package api

import "testing"

func TestError(t *testing.T) {
	var err = APIError{ "errorname", "foobar", []string {} }
	var result = err.Error()
	if result != "errorname" {
		t.Errorf("Expected 'errorname' but got: %s", result)
	}
}
