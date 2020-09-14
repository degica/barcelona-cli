package operations

import (
	"github.com/degica/barcelona-cli/config"
)

type mockConfig struct{}

func (m mockConfig) LoadLogin() *config.Login {
	return &config.Login{
		Token:    "",
		Endpoint: "https://test.example.com",
	}
}

func ExampleLoginInfoOperationOutput() {

	op := NewLoginInfoOperation(&mockConfig{})
	op.Run()

	// Output:
	// Endpoint: https://test.example.com
}
