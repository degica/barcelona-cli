package operations

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"

	"encoding/json"

	"github.com/degica/barcelona-cli/api"
	"github.com/degica/barcelona-cli/utils"
)

type LoginOperationExternals interface {
	// User Input Reader
	Read(secret bool) (string, error)

	// CommandRunner
	RunCommand(name string, arg ...string) error

	LoginWithGithub(endpoint string, token string) (*api.User, error)
	LoginWithVault(endpoint string, vault_url string, token string) (*api.User, error)
	ReloadDefaultClient() (*api.Client, error)
	Patch(path string, body io.Reader) ([]byte, error)

	// Config stuff
	WriteLogin(auth string, token string, endpoint string) error
	GetPublicKeyPath() string
	GetPrivateKeyPath() string
}

type LoginOperation struct {
	endpoint    string
	backend     string
	gh_token    string
	vault_token string
	vault_url   string
	ext         LoginOperationExternals
}

func NewLoginOperation(endpoint string, backend string, gh_token string, vault_token string, vault_url string, ext LoginOperationExternals) *LoginOperation {
	return &LoginOperation{
		endpoint:    endpoint,
		backend:     backend,
		gh_token:    gh_token,
		vault_token: vault_token,
		vault_url:   vault_url,
		ext:         ext,
	}
}

func githubLogin(oper LoginOperation, user *api.User) *runResult {
	fmt.Println("Logging in with Github")
	token := oper.gh_token
	if len(token) == 0 {
		fmt.Println("Create new GitHub access token with read:org permission here https://github.com/settings/tokens/new")
		token = utils.Ask("GitHub Token", true, true, oper.ext)
	}

	user, err := oper.ext.LoginWithGithub(oper.endpoint, token)
	if err != nil {
		return error_result(err.Error())
	}

	err = oper.ext.WriteLogin(oper.backend, user.Token, oper.endpoint)
	if err != nil {
		return error_result(err.Error())
	}

	return ok_result()
}

func vaultLogin(oper LoginOperation, user *api.User) *runResult {
	fmt.Println("Logging in with Vault")
	token := oper.vault_token
	url := oper.vault_url
	if len(token) == 0 {
		fmt.Println("Create new GitHub access token with read:org permission here https://github.com/settings/tokens/new")
		token = utils.Ask("GitHub Token", true, true, oper.ext)
	}
	if len(url) == 0 {
		fmt.Println("URL of vault server (e.g. https://vault.degica.com)")
		url = utils.Ask("Vault server URL", true, false, oper.ext)
	}
	user, err := oper.ext.LoginWithVault(oper.endpoint, url, token)
	if err != nil {
		return error_result(err.Error())
	}

	err = oper.ext.WriteLogin(oper.backend, user.Token, oper.endpoint)
	if err != nil {
		return error_result(err.Error())
	}

	return ok_result()
}

func setUpKeys(oper LoginOperation, user *api.User) *runResult {
	keyExists := utils.FileExists(oper.ext.GetPublicKeyPath())
	if !keyExists {
		fmt.Println("Generating your SSH key pair...")
		err := oper.ext.RunCommand("ssh-keygen",
			"-t", "ecdsa",
			"-b", "521",
			"-f", oper.ext.GetPrivateKeyPath(),
			"-C", "")
		if err != nil {
			return error_result(err.Error())
		}
	}

	if !keyExists || len(user.PublicKey) == 0 {
		fmt.Println("Registering your public key...")

		pubKeyB, err := ioutil.ReadFile(oper.ext.GetPublicKeyPath())
		if err != nil {
			return error_result(err.Error())
		}

		re := regexp.MustCompile(" *\n$")
		pubKey := re.ReplaceAllString(string(pubKeyB), "")
		reqBody := make(map[string]string)
		reqBody["public_key"] = pubKey
		bodyB, err := json.Marshal(reqBody)
		reloaded_client, err := oper.ext.ReloadDefaultClient()
		if err != nil {
			return error_result(err.Error())
		}

		_, err = reloaded_client.Patch("/user", bytes.NewBuffer(bodyB))
		if err != nil {
			return error_result(err.Error())
		}
	}

	return ok_result()
}

func (oper LoginOperation) run() *runResult {
	if len(oper.endpoint) == 0 {
		return error_result("endpoint is required")
	}

	var user api.User

	switch oper.backend {
	case "github":
		result := githubLogin(oper, &user)
		if result.is_error {
			return result
		}
	case "vault":
		result := vaultLogin(oper, &user)
		if result.is_error {
			return result
		}
	default:
		return error_result("Unrecognized auth backend")
	}

	return setUpKeys(oper, &user)
}
