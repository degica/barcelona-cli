package operations

import (
	"github.com/degica/barcelona-cli/api"
	"io"
	"testing"
)

type mockLoginOperationError struct {
	msg string
}

func (m mockLoginOperationError) Error() string {
	return m.msg
}

type mockLoginOperationExternals struct {
	readString string
	readError  error

	loginWithGithubUser  *api.User
	loginWithGithubError error

	loginWithVaultUser  *api.User
	loginWithVaultError error

	readFileBytes []byte
	readFileError error

	patchBytes []byte
	patchError error

	fileExistsBool bool
}

func (m mockLoginOperationExternals) Read(secret bool) (string, error) {
	return m.readString, m.readError
}

func (m mockLoginOperationExternals) RunCommand(name string, arg ...string) error {
	return nil
}

func (m mockLoginOperationExternals) FileExists(path string) bool {
	return m.fileExistsBool
}

func (m mockLoginOperationExternals) ReadFile(path string) ([]byte, error) {
	return m.readFileBytes, m.readFileError
}

func (m mockLoginOperationExternals) LoginWithGithub(endpoint string, token string) (*api.User, error) {
	return m.loginWithGithubUser, m.loginWithGithubError
}

func (m mockLoginOperationExternals) LoginWithVault(endpoint string, vault_url string, token string) (*api.User, error) {
	return m.loginWithVaultUser, m.loginWithVaultError
}

func (m mockLoginOperationExternals) ReloadDefaultClient() (LoginOperationClient, error) {
	return m, nil
}

func (m mockLoginOperationExternals) Patch(path string, body io.Reader) ([]byte, error) {
	return m.patchBytes, m.patchError
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
		t.Errorf("Expected 'Unrecognized auth backend' to be returned. But got '%s'", result.message)
	}
}

func TestGithubBackend(t *testing.T) {
	ext := &mockLoginOperationExternals{
		readString:           "aw\n",
		readError:            nil,
		loginWithGithubUser:  &api.User{},
		loginWithGithubError: nil,
		readFileBytes:        []byte("stuff"),
	}

	op := NewLoginOperation("https://endpoint", "github", "", "", "", ext)
	result := op.run()

	if result.is_error != false {
		t.Errorf("Expected no error to be returned.")
	}
}

func ExampleLoginOperation_run_with_github_already_has_ssh() {
	ext := &mockLoginOperationExternals{
		readString:           "aw\n",
		readError:            nil,
		loginWithGithubUser:  &api.User{},
		loginWithGithubError: nil,
		readFileBytes:        []byte("stuff"),
		fileExistsBool:       true,
	}

	op := NewLoginOperation("https://endpoint", "github", "", "", "", ext)
	op.run()

	// Output:
	// Logging in with Github
	// Create new GitHub access token with read:org permission here https://github.com/settings/tokens/new
	// GitHub Token: Registering your public key...
}

func ExampleLoginOperation_run_with_github_token_already_has_ssh() {
	ext := &mockLoginOperationExternals{
		readString:           "aw\n",
		readError:            nil,
		loginWithGithubUser:  &api.User{},
		loginWithGithubError: nil,
		readFileBytes:        []byte("stuff"),
		fileExistsBool:       true,
	}

	op := NewLoginOperation("https://endpoint", "github", "gh_token", "", "", ext)
	op.run()

	// Output:
	// Logging in with Github
	// Registering your public key...
}

func ExampleLoginOperation_run_with_github() {
	ext := &mockLoginOperationExternals{
		readString:           "aw\n",
		readError:            nil,
		loginWithGithubUser:  &api.User{},
		loginWithGithubError: nil,
		readFileBytes:        []byte("stuff"),
	}

	op := NewLoginOperation("https://endpoint", "github", "", "", "", ext)
	op.run()

	// Output:
	// Logging in with Github
	// Create new GitHub access token with read:org permission here https://github.com/settings/tokens/new
	// GitHub Token: Generating your SSH key pair...
	// Registering your public key...
}

func ExampleLoginOperation_run_with_github_token() {
	ext := &mockLoginOperationExternals{
		readString:           "aw\n",
		readError:            nil,
		loginWithGithubUser:  &api.User{},
		loginWithGithubError: nil,
		readFileBytes:        []byte("stuff"),
	}

	op := NewLoginOperation("https://endpoint", "github", "gh_token", "", "", ext)
	op.run()

	// Output:
	// Logging in with Github
	// Generating your SSH key pair...
	// Registering your public key...
}

func TestVaultBackend(t *testing.T) {
	ext := &mockLoginOperationExternals{
		readString:          "aw\n",
		readError:           nil,
		loginWithVaultUser:  &api.User{},
		loginWithVaultError: nil,
		readFileBytes:       []byte("stuff"),
	}

	op := NewLoginOperation("https://endpoint", "vault", "", "", "", ext)
	result := op.run()

	if result.is_error != false {
		t.Errorf("Expected no error to be returned.")
	}
}

func ExampleLoginOperation_run_with_vault_already_has_ssh() {
	ext := &mockLoginOperationExternals{
		readString:          "aw\n",
		readError:           nil,
		loginWithVaultUser:  &api.User{},
		loginWithVaultError: nil,
		readFileBytes:       []byte("stuff"),
		fileExistsBool:      true,
	}

	op := NewLoginOperation("https://endpoint", "vault", "", "", "", ext)
	op.run()

	// Output:
	// Logging in with Vault
	// Create new GitHub access token with read:org permission here https://github.com/settings/tokens/new
	// GitHub Token: URL of vault server (e.g. https://vault.degica.com)
	// Vault server URL: Registering your public key...
}

func ExampleLoginOperation_run_with_vault_token_already_has_ssh() {
	ext := &mockLoginOperationExternals{
		readString:          "aw\n",
		readError:           nil,
		loginWithVaultUser:  &api.User{},
		loginWithVaultError: nil,
		readFileBytes:       []byte("stuff"),
		fileExistsBool:      true,
	}

	op := NewLoginOperation("https://endpoint", "vault", "", "gh_token", "", ext)
	op.run()

	// Output:
	// Logging in with Vault
	// URL of vault server (e.g. https://vault.degica.com)
	// Vault server URL: Registering your public key...
}

func ExampleLoginOperation_run_with_vault() {
	ext := &mockLoginOperationExternals{
		readString:          "aw\n",
		readError:           nil,
		loginWithVaultUser:  &api.User{},
		loginWithVaultError: nil,
		readFileBytes:       []byte("stuff"),
	}

	op := NewLoginOperation("https://endpoint", "vault", "", "", "", ext)
	op.run()

	// Output:
	// Logging in with Vault
	// Create new GitHub access token with read:org permission here https://github.com/settings/tokens/new
	// GitHub Token: URL of vault server (e.g. https://vault.degica.com)
	// Vault server URL: Generating your SSH key pair...
	// Registering your public key...
}

func ExampleLoginOperation_run_with_vault_token() {
	ext := &mockLoginOperationExternals{
		readString:          "aw\n",
		readError:           nil,
		loginWithVaultUser:  &api.User{},
		loginWithVaultError: nil,
		readFileBytes:       []byte("stuff"),
	}

	op := NewLoginOperation("https://endpoint", "vault", "", "gh_token", "", ext)
	op.run()

	// Output:
	// Logging in with Vault
	// URL of vault server (e.g. https://vault.degica.com)
	// Vault server URL: Generating your SSH key pair...
	// Registering your public key...
}

func ExampleLoginOperation_run_with_vault_already_has_ssh_given_url() {
	ext := &mockLoginOperationExternals{
		readString:          "aw\n",
		readError:           nil,
		loginWithVaultUser:  &api.User{},
		loginWithVaultError: nil,
		readFileBytes:       []byte("stuff"),
		fileExistsBool:      true,
	}

	op := NewLoginOperation("https://endpoint", "vault", "", "", "https://vaultserv", ext)
	op.run()

	// Output:
	// Logging in with Vault
	// Create new GitHub access token with read:org permission here https://github.com/settings/tokens/new
	// GitHub Token: Registering your public key...
}

func ExampleLoginOperation_run_with_vault_token_already_has_ssh_given_url() {
	ext := &mockLoginOperationExternals{
		readString:          "aw\n",
		readError:           nil,
		loginWithVaultUser:  &api.User{},
		loginWithVaultError: nil,
		readFileBytes:       []byte("stuff"),
		fileExistsBool:      true,
	}

	op := NewLoginOperation("https://endpoint", "vault", "", "gh_token", "https://vaultserv", ext)
	op.run()

	// Output:
	// Logging in with Vault
	// Registering your public key...
}

func ExampleLoginOperation_run_with_vault_given_url() {
	ext := &mockLoginOperationExternals{
		readString:          "aw\n",
		readError:           nil,
		loginWithVaultUser:  &api.User{},
		loginWithVaultError: nil,
		readFileBytes:       []byte("stuff"),
	}

	op := NewLoginOperation("https://endpoint", "vault", "", "", "https://vaultserv", ext)
	op.run()

	// Output:
	// Logging in with Vault
	// Create new GitHub access token with read:org permission here https://github.com/settings/tokens/new
	// GitHub Token: Generating your SSH key pair...
	// Registering your public key...
}

func ExampleLoginOperation_run_with_vault_token_given_url() {
	ext := &mockLoginOperationExternals{
		readString:          "aw\n",
		readError:           nil,
		loginWithVaultUser:  &api.User{},
		loginWithVaultError: nil,
		readFileBytes:       []byte("stuff"),
	}

	op := NewLoginOperation("https://endpoint", "vault", "", "gh_token", "https://vaultserv", ext)
	op.run()

	// Output:
	// Logging in with Vault
	// Generating your SSH key pair...
	// Registering your public key...
}

func ExampleLoginOperation_run_output() {

	op := NewLoginOperation("https://endpoint", "somethingrando", "gh_token", "vault_token", "https://vault_url", &mockLoginOperationExternals{})
	op.run()

	// Output:
	//
}
