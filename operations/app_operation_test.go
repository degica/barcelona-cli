package operations

import (
	"bytes"
	"io"
)

type MockAppOperationApiClient struct {
}

func (client MockAppOperationApiClient) Get(path string, body io.Reader) ([]byte, error) {
	return bytes.NewBufferString("{\"say\":\"hello\"}").Bytes(), nil
}

func (client MockAppOperationApiClient) Delete(path string, body io.Reader) ([]byte, error) {
	return bytes.NewBufferString("").Bytes(), nil
}

type MockAppOperationInputReaderNo struct {
}

func (_ MockAppOperationInputReaderNo) Read(_ bool) (string, error) {
	return "n\n", nil
}

type MockAppOperationInputReaderYes struct {
}

func (_ MockAppOperationInputReaderYes) Read(_ bool) (string, error) {
	return "y\n", nil
}

func ExampleAppOperation_run_delete_no_confirm_output() {
	client := &MockAppOperationApiClient{}
	oper := NewAppOperation("asd", Delete, true, client, nil /* because it doesnt matter */)

	oper.run()
	// Output:
	// You are attempting to delete asd
	// Deleted asd
}

func ExampleAppOperation_run_delete_confirm_y_output() {
	client := &MockAppOperationApiClient{}
	oper := NewAppOperation("asd", Delete, false, client, &MockAppOperationInputReaderYes{})

	oper.run()
	// Output:
	// You are attempting to delete asd
	// This operation cannot be undone. Are you sure? [y/n]: Deleted asd
}

func ExampleAppOperation_run_delete_confirm_n_output() {
	client := &MockAppOperationApiClient{}
	oper := NewAppOperation("asd", Delete, false, client, &MockAppOperationInputReaderNo{})

	oper.run()
	// Output:
	// You are attempting to delete asd
	// This operation cannot be undone. Are you sure? [y/n]:
}

func ExampleAppOperation_run_show_output_when_no_heritage() {
	client := &MockAppOperationApiClient{}
	oper := NewAppOperation("asd", Show, true, client, &MockAppOperationInputReaderYes{})

	oper.run()
	// Output:
}
