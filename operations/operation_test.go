package operations

import (
	"testing"
)

type errorOperation struct{}

func (m errorOperation) run() *runResult {
	return error_result("some error")
}

func TestExecutingOperationThatReturnsError(t *testing.T) {
	result := Execute(&errorOperation{})

	if result == nil {
		t.Errorf("Expected an error to be returned.")
	}

	if result.Error() != "some error" {
		t.Errorf("Expected error message to be 'some_error' but got %s", result.Error())
	}
}

type nilOperation struct{}

func (m nilOperation) run() *runResult {
	return nil
}

func TestExecutingOperationThatReturnsNil(t *testing.T) {
	result := Execute(&nilOperation{})

	if result != nil {
		t.Errorf("Expected nil to be returned.")
	}
}

type okOperation struct{}

func (m okOperation) run() *runResult {
	return ok_result()
}

func TestExecutingOperationThatReturnsOk(t *testing.T) {
	result := Execute(&okOperation{})

	if result != nil {
		t.Errorf("Expected nil to be returned.")
	}
}

func TestOkResult(t *testing.T) {
	result := ok_result()

	if result.is_error != false {
		t.Errorf("Expected is_error to be false but is true")
	}

	if result.message != "" {
		t.Errorf("Expected message to be empty string but got %s", result.message)
	}
}

func TestErrorResult(t *testing.T) {
	result := error_result("error message")

	if result.is_error != true {
		t.Errorf("Expected is_error to be true but is false")
	}

	if result.message != "error message" {
		t.Errorf("Expected message to be 'error message' but got %s", result.message)
	}
}
