package operations

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"

	"encoding/json"

	"github.com/degica/barcelona-cli/api"
	"github.com/degica/barcelona-cli/utils"
)

type LoginOperationApiClient interface {
	LoginWithGithub(endpoint string, token string) (*api.User, error)
	LoginWithVault(endpoint string, vault_url string, token string) (*api.User, error)
	ReloadDefaultClient() (*api.Client, error)
	Patch(path string, body io.Reader) ([]byte, error)
}

type LoginConfig interface {
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
	client      LoginOperationApiClient
	cfg         LoginConfig
}

func NewLoginOperation(endpoint string, backend string, gh_token string, vault_token string, vault_url string, client LoginOperationApiClient, cfg LoginConfig) *LoginOperation {
	return &LoginOperation{
		endpoint:    endpoint,
		backend:     backend,
		gh_token:    gh_token,
		vault_token: vault_token,
		vault_url:   vault_url,
		client:      client,
		cfg:         cfg,
	}
}

func githubLogin(oper LoginOperation, user *api.User) *runResult {
	fmt.Println("Logging in with Github")
	token := oper.gh_token
	if len(token) == 0 {
		fmt.Println("Create new GitHub access token with read:org permission here https://github.com/settings/tokens/new")
		token = utils.Ask("GitHub Token", true, true, utils.NewStdinInputReader())
	}

	user, err := oper.client.LoginWithGithub(oper.endpoint, token)
	if err != nil {
		return error_result(err.Error())
	}

	err = oper.cfg.WriteLogin(oper.backend, user.Token, oper.endpoint)
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
		token = utils.Ask("GitHub Token", true, true, utils.NewStdinInputReader())
	}
	if len(url) == 0 {
		fmt.Println("URL of vault server (e.g. https://vault.degica.com)")
		url = utils.Ask("Vault server URL", true, false, utils.NewStdinInputReader())
	}
	user, err := oper.client.LoginWithVault(oper.endpoint, url, token)
	if err != nil {
		return error_result(err.Error())
	}

	err = oper.cfg.WriteLogin(oper.backend, user.Token, oper.endpoint)
	if err != nil {
		return error_result(err.Error())
	}

	return ok_result()
}

func setUpKeys(oper LoginOperation, user *api.User) *runResult {
	keyExists := utils.FileExists(oper.cfg.GetPublicKeyPath())
	if !keyExists {
		fmt.Println("Generating your SSH key pair...")
		cmd := exec.Command("ssh-keygen",
			"-t", "ecdsa",
			"-b", "521",
			"-f", oper.cfg.GetPrivateKeyPath(),
			"-C", "")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return error_result(err.Error())
		}
	}

	if !keyExists || len(user.PublicKey) == 0 {
		fmt.Println("Registering your public key...")

		pubKeyB, err := ioutil.ReadFile(oper.cfg.GetPublicKeyPath())
		if err != nil {
			return error_result(err.Error())
		}

		re := regexp.MustCompile(" *\n$")
		pubKey := re.ReplaceAllString(string(pubKeyB), "")
		reqBody := make(map[string]string)
		reqBody["public_key"] = pubKey
		bodyB, err := json.Marshal(reqBody)
		oper.client, err = oper.client.ReloadDefaultClient()
		if err != nil {
			return error_result(err.Error())
		}

		_, err = oper.client.Patch("/user", bytes.NewBuffer(bodyB))
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
