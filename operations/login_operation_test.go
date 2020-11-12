package operations

import (
	"io"
	"testing"
	"github.com/degica/barcelona-cli/api"
)

type mockLoginOperationError struct {
	msg string
}

func (m mockLoginOperationError) Error() string {
	return m.msg
}

type mockLoginOperationExternals struct {
	readString string
	readError error

	loginWithGithubUser *api.User
	loginWithGithubError error
}

func (m mockLoginOperationExternals) Read(secret bool) (string, error) {
	return m.readString, m.readError
}

func (m mockLoginOperationExternals) RunCommand(name string, arg ...string) error {
	return nil
}

func (m mockLoginOperationExternals) FileExists(path string) bool {
	return false
}

func (m mockLoginOperationExternals) LoginWithGithub(endpoint string, token string) (*api.User, error) {
	return m.loginWithGithubUser, m.loginWithGithubError
}

func (m mockLoginOperationExternals) LoginWithVault(endpoint string, vault_url string, token string) (*api.User, error) {
	return nil, nil
}

func (m mockLoginOperationExternals) ReloadDefaultClient() (*api.Client, error) {
	return nil, nil
}

func (m mockLoginOperationExternals) Patch(path string, body io.Reader) ([]byte, error) {
	return nil, nil
}

func (m mockLoginOperationExternals) WriteLogin(auth string, token string, endpoint string) error {
	return nil
}

func (m mockLoginOperationExternals) GetPublicKeyPath() string {
	return ""
}

func (m mockLoginOperationExternals) GetPrivateKeyPath() string {
	return ""
}

func TestUnknownBackend(t *testing.T) {
	op := NewLoginOperation("https://endpoint", "mybckend", "gh_token", "vault_token", "https://vault_url", &mockLoginOperationExternals{})
	result := op.run()

	if result.is_error != true {
		t.Errorf("Expected no error to be returned.")
	}

	if result.message != "Unrecognized auth backend" {
		t.Errorf("Expected 'Unrecognized auth backend' to be returned.")
	}
}

func TestGithubBackend(t *testing.T) {
	ext := &mockLoginOperationExternals{
		readString: "aw\n",
		readError: nil,
		loginWithGithubUser: &api.User{},
		loginWithGithubError: nil,
	}

	op := NewLoginOperation("https://endpoint", "github", "", "", "", ext)
	result := op.run()

	if result.is_error != true {
		t.Errorf("Expected no error to be returned.")
	}

	if result.message != "Unrecognized auth backend" {
		t.Errorf("Expected 'Unrecognized auth backend' to be returned." + result.message)
	}
}

func ExampleLoginOperationApiClient_run_output() {

	op := NewLoginOperation("https://endpoint", "mybckend", "gh_token", "vault_token", "https://vault_url", &mockLoginOperationExternals{})
	op.run()

	// Output:
	// 
}
