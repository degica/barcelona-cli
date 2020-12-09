package operations

import ()

type mockLogin struct{}

func (m mockLogin) GetAuth() string {
	return "fooauth"
}

func (m mockLogin) GetEndpoint() string {
	return "https://example.com"
}

func ExampleLoginInfoOperation_run_output() {

	op := NewLoginInfoOperation(&mockLogin{})
	op.run()

	// Output:
	// Endpoint: https://example.com
	// Auth:     fooauth
}
