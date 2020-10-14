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
	optype       OperationType
	no_confirm   bool
	client       AppOperationApiClient
	input_reader utils.UserInputReader
}

func NewAppOperation(name string, optype OperationType, no_confirm bool, client AppOperationApiClient, input_reader utils.UserInputReader) *AppOperation {
	return &AppOperation{
		name:         name,
		optype:       optype,
		no_confirm:   no_confirm,
		client:       client,
		input_reader: input_reader,
	}
}

func app_delete(oper AppOperation) *runResult {
	fmt.Printf("You are attempting to delete %s\n", oper.name)
	if !oper.no_confirm && !utils.AreYouSure("This operation cannot be undone. Are you sure?", oper.input_reader) {
		return nil
	}

	_, err := oper.client.Delete("/heritages/"+oper.name, nil)
	if err != nil {
		return error_result(err.Error())
	}
	fmt.Printf("Deleted %s\n", oper.name)

	return ok_result()
}

func app_show(oper AppOperation) *runResult {
	resp, err := oper.client.Get("/heritages/"+oper.name, nil)
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

func (oper AppOperation) run() *runResult {
	if len(oper.name) == 0 {
		return error_result("district name is required")
	}

	if oper.optype == Delete {
		return app_delete(oper)
	}

	if oper.optype == Show {
		return app_show(oper)
	}

	return error_result("unknown operation")
}
