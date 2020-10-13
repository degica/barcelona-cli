package operations

import (
	"bytes"
	"io"
	"testing"
)

type MockNormalApiOperationApiClient struct{}

func (client MockNormalApiOperationApiClient) Request(method string, path string, body io.Reader) ([]byte, error) {
	return bytes.NewBufferString("{\"say\":\"hello\"}").Bytes(), nil
}

func TestExecutingApiOperation(t *testing.T) {
	client := &MockNormalApiOperationApiClient{}
	oper := NewApiOperation("GET", "https://somewhere", bytes.NewBufferString(""), client)
	result := oper.run()

	if result.is_error != false {
		t.Errorf("Expected no error to be returned.")
	}
}

func ExampleApiOperation_run_output() {
	client := &MockNormalApiOperationApiClient{}
	oper := NewApiOperation("GET", "https://somewhere", bytes.NewBufferString(""), client)

	oper.run()
	// Output:
	// {
	//   "say": "hello"
	// }
}

type MockApiError struct{}

func (e *MockApiError) Error() string { return "some error msg" }

type MockErrorApiOperationApiClient struct{}

func (client MockErrorApiOperationApiClient) Request(method string, path string, body io.Reader) ([]byte, error) {
	return nil, &MockApiError{}
}

func TestExecutingApiOperationReturningError(t *testing.T) {
	client := &MockErrorApiOperationApiClient{}
	oper := NewApiOperation("GET", "https://somewhere", bytes.NewBufferString(""), client)
	result := oper.run()

	if result.is_error != true {
		t.Errorf("Expected an error to be returned.")
	}

	if result.message != "some error msg" {
		t.Errorf("Expected 'some error msg' to be returned. But got %s", result.message)
	}
}

func ExampleApiOperation_run_with_error_output() {
	client := &MockErrorApiOperationApiClient{}
	oper := NewApiOperation("GET", "https://somewhere", bytes.NewBufferString(""), client)

	oper.run()
	// Output:
}
