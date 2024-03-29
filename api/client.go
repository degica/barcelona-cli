package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"regexp"
	"time"

	"github.com/degica/barcelona-cli/config"
)

type Config struct {
	Version    int
	HttpClient *http.Client
}

var hasAutoRefresh = false

type Client struct {
	login  *config.Login
	config *Config
}

var pathPrefix = "/v1"
var DefaultClient *Client

func init() {
	err := reloadDefaultClient()
	if err != nil {
		panic("Couldn't initialize default client: " + err.Error())
	}
}

func DefaultConfig() *Config {
	return &Config{
		Version:    1,
		HttpClient: &http.Client{Timeout: time.Duration(60) * time.Second},
	}
}

func NewClient(c *Config, l *config.Login) *Client {
	return &Client{login: l, config: c}
}

func newDefaultClient() *Client {
	c := DefaultConfig()
	l := config.Get().LoadLogin()
	return NewClient(c, l)
}

func reloadDefaultClient() error {
	DefaultClient = newDefaultClient()
	return nil
}

func (cli *Client) rawRequest(req *http.Request) ([]byte, error) {
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	if config.Get().IsDebug() {
		q := req.URL.Query()
		q.Add("debug", "true")
		req.URL.RawQuery = q.Encode()
		dump(httputil.DumpRequestOut(req, true))
	}

	resp, err := cli.config.HttpClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if config.Get().IsDebug() {
		//Is dumping response too much? Probably not.
		dump(httputil.DumpResponse(resp, true))
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		var apiErr APIError
		err := json.Unmarshal(b, &apiErr)

		if err != nil {
			// Not a JSON error
			resstr := "HTTP Error:" + resp.Status + "\nType bcn -d <command> for details"

			apiErr = APIError{resstr, string(b[:]), []string{}}
		}

		return b, &apiErr
	}
	return b, nil
}

func (cli *Client) ReloadDefaultClient() (*Client, error) {
	err := reloadDefaultClient()
	return DefaultClient, err
}

func (cli *Client) Request(method string, path string, body io.Reader) ([]byte, error) {
	if !hasAutoRefresh {
		hasAutoRefresh = true
		err := cli.autoRefreshVaultToken()
		if err != nil {
			return nil, err
		}
	}

	url := cli.login.GetEndpoint() + pathPrefix + path

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	switch cli.login.GetAuth() {
	case "github":
		if len(cli.login.GetToken()) > 0 {
			req.Header.Add("X-Barcelona-Token", cli.login.GetToken())
		}
	case "vault":
		req.Header.Add("X-Vault-Token", cli.login.GetToken())
	}

	return cli.rawRequest(req)
}

func (cli *Client) Get(path string, body io.Reader) ([]byte, error) {
	return cli.Request("GET", path, body)
}

func (cli *Client) Post(path string, body io.Reader) ([]byte, error) {
	return cli.Request("POST", path, body)
}

func (cli *Client) Patch(path string, body io.Reader) ([]byte, error) {
	return cli.Request("PATCH", path, body)
}

func (cli *Client) Put(path string, body io.Reader) ([]byte, error) {
	return cli.Request("PUT", path, body)
}

func (cli *Client) Delete(path string, body io.Reader) ([]byte, error) {
	return cli.Request("DELETE", path, body)
}

func dump(dump []byte, err error) {
	s := string(dump)
	regex, err := regexp.Compile("(Token): ([0-9A-Za-z]+)")
	if err != nil {
		panic(err.Error())
	}

	ss := regex.ReplaceAllString(s, "$1: [filtered]")
	fmt.Printf("%s\n", ss)
}

func (cli *Client) autoRefreshVaultToken() error {
	login := config.Get().LoadLogin()
	if login.Auth != "vault" || login.VaultUrl == "" || login.VaultToken == "" {
		return nil
	}

	user, err := cli.LoginWithVault(login.VaultUrl, login.VaultToken)

	if err != nil {
		return err
	}

	backend := login.Auth
	endpoint := login.Endpoint
	vaultUrl := login.VaultUrl
	vaultToken := login.VaultToken

	err = config.Get().WriteLogin(backend, user.Token, endpoint, vaultUrl, vaultToken)
	if err != nil {
		return err
	}

	return nil
}
