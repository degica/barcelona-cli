package operations

import (
	"fmt"

	"github.com/degica/barcelona-cli/config"
)

type LoginInfo interface {
	GetEndpoint() string
	GetAuth() string
}

type LoginConfiguration interface {
	LoadLogin() *config.Login
}

type LoginInfoOperation struct {
	cfg LoginConfiguration
}

func NewLoginInfoOperation(cfg LoginConfiguration) *LoginInfoOperation {
	// set only specific field value with field key
	return &LoginInfoOperation{
		cfg: cfg,
	}
}

func (oper LoginInfoOperation) run() *runResult {
	login := oper.cfg.LoadLogin()

	fmt.Printf("Endpoint: %s\n", login.GetEndpoint())
	fmt.Printf("Auth:     %s\n", login.GetAuth())

	return ok_result()
}
