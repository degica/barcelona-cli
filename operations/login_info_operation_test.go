package operations

import (
	"github.com/degica/barcelona-cli/config"
)

type mockConfig struct{}

func (m mockConfig) LoadLogin() *config.Login {
	return &config.Login{
		Auth:     "vault",
		Endpoint: "https://test.example.com",
	}
}

func ExampleLoginInfoOperation_run_output() {

	op := NewLoginInfoOperation(&mockConfig{})
	op.run()

	// Output:
	// Endpoint: https://test.example.com
	// Auth:     vault
}
