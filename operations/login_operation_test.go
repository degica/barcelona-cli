package operations

import ()

type mockLoginOperationExternals struct {
}

func ExampleLoginOperationApiClient_run_output() {

	op := NewLoginOperation("https://endpoint", "mybckend", "gh_token", "vault_token", "https://vault_url", &mockLoginOperationExternals{})
	op.run()

	// Output:
	// Endpoint: https://test.example.com
	// Auth:     vault
}
