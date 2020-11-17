package operations

import (
	"fmt"
)

type LoginInfo interface {
	GetEndpoint() string
	GetAuth() string
}

type LoginInfoOperation struct {
	info LoginInfo
}

func NewLoginInfoOperation(info LoginInfo) *LoginInfoOperation {
	// set only specific field value with field key
	return &LoginInfoOperation{
		info: info,
	}
}

func (oper LoginInfoOperation) run() *runResult {
	fmt.Printf("Endpoint: %s\n", oper.info.GetEndpoint())
	fmt.Printf("Auth:     %s\n", oper.info.GetAuth())

	return ok_result()
}
