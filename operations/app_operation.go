package operations

import (
	"fmt"
	"io"

	"encoding/json"

	"github.com/degica/barcelona-cli/api"
	"github.com/degica/barcelona-cli/utils"
)

type AppOperationApiClient interface {
	Get(path string, body io.Reader) ([]byte, error)
	Delete(path string, body io.Reader) ([]byte, error)
}

type AppOperation struct {
	name         string
	op_type      OperationType
	no_confirm   bool
	client       AppOperationApiClient
	input_reader utils.UserInputReader
}

func NewAppOperation(name string, op_type OperationType, no_confirm bool, client AppOperationApiClient, input_reader utils.UserInputReader) *AppOperation {
	return &AppOperation{
		name:         name,
		op_type:      op_type,
		no_confirm:   no_confirm,
		client:       client,
		input_reader: input_reader,
	}
}

func app_delete(operation AppOperation) *runResult {
	fmt.Printf("You are attempting to delete %s\n", operation.name)
	if !operation.no_confirm && !utils.AreYouSure("This operation cannot be undone. Are you sure?", operation.input_reader) {
		return nil
	}

	_, err := operation.client.Delete("/heritages/"+operation.name, nil)
	if err != nil {
		return error_result(err.Error())
	}
	fmt.Printf("Deleted %s\n", operation.name)

	return ok_result()
}

func app_show(operation AppOperation) *runResult {
	resp, err := operation.client.Get("/heritages/"+operation.name, nil)
	if err != nil {
		return error_result(err.Error())
	}
	var hResp api.HeritageResponse
	err = json.Unmarshal(resp, &hResp)
	if err != nil {
		return error_result(err.Error())
	}
	if hResp.Heritage == nil {
		return error_result("No such heritage")
	}
	hResp.Heritage.Print()

	return ok_result()
}

func (operation AppOperation) run() *runResult {
	if len(operation.name) == 0 {
		return error_result("district name is required")
	}

	if operation.op_type == Delete {
		return app_delete(operation)
	}

	if operation.op_type == Show {
		return app_show(operation)
	}

	return error_result("unknown operation")
}
