package operations

import (
	"fmt"

	"github.com/degica/barcelona-cli/config"
)

type LoginInfoOperation struct {
}

func NewLoginInfoOperation() *LoginInfoOperation {
    // set only specific field value with field key
    return &LoginInfoOperation{}
}

func (oper LoginInfoOperation) Run() error {
	login := config.LoadLogin()

	fmt.Printf("Endpoint: %s\n", login.Endpoint)
  fmt.Printf("Auth:     %s\n", login.Auth)

	return nil
}
