package operations

import (
	"fmt"

	"github.com/degica/barcelona-cli/config"
)

type LoginInfoOperation struct {
	cfg config.Configuration
}

func NewLoginInfoOperation(cfg config.Configuration) *LoginInfoOperation {
	// set only specific field value with field key
	return &LoginInfoOperation{
		cfg: cfg,
	}
}

func (oper LoginInfoOperation) Run() error {
	login := oper.cfg.LoadLogin()

	fmt.Printf("Endpoint: %s\n", login.Endpoint)
	fmt.Printf("Auth:     %s\n", login.Auth)

	return nil
}
