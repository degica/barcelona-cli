package operations

import (
	"bytes"
	"fmt"
	"io"

	"github.com/degica/barcelona-cli/utils"
)

type ApiOperationApiClient interface {
	Request(method string, path string, body io.Reader) ([]byte, error)
}

type ApiOperation struct {
	method string
	path   string
	body   *bytes.Buffer
	client ApiOperationApiClient
}

func NewApiOperation(method string, path string, body *bytes.Buffer, client ApiOperationApiClient) *ApiOperation {
	return &ApiOperation{
		method: method,
		path:   path,
		body:   body,
		client: client,
	}
}

func (apiOperation ApiOperation) run() *runResult {
	if len(apiOperation.method) == 0 {
		return error_result("method is required")
	}
	if len(apiOperation.path) == 0 {
		return error_result("path is required")
	}

	response, err := apiOperation.client.Request(apiOperation.method, apiOperation.path, apiOperation.body)
	if err != nil {
		return error_result(err.Error())
	}

	fmt.Println(utils.PrettyJSON(response))

	return ok_result()
}
