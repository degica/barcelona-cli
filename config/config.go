package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
)

var ConfigDir string
var LoginFilePath string
var privateKeyPath string
var publicKeyPath string
var CertPath string
var Debug bool

// Clients should get configs using this function
func Get() *LocalConfig {
	return &LocalConfig{}
}

// Implementation of our Configuration object
type LocalConfig struct{}

func (m LocalConfig) LoadLogin() *Login {
	return LoadLogin()
}

func (m LocalConfig) GetPrivateKeyPath() string {
	return privateKeyPath
}

func (m LocalConfig) GetPublicKeyPath() string {
	return publicKeyPath
}

func (m LocalConfig) WriteLogin(auth string, token string, endpoint string) error {
	login := &Login{
		Auth:     auth,
		Token:    token,
		Endpoint: endpoint,
	}

	return writeLogin(login)
}

func init() {
	path, err := getConfigPath()
	if err != nil {
		panic("Couldn't get login path")
	}
	ConfigDir = path
	LoginFilePath = filepath.Join(ConfigDir, "login")
	privateKeyPath = filepath.Join(ConfigDir, "id_ecdsa")
	publicKeyPath = filepath.Join(ConfigDir, "id_ecdsa.pub")
	CertPath = filepath.Join(ConfigDir, "id_ecdsa-cert.pub")
}

func getConfigPath() (string, error) {
	homedir, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homedir, ".bcn"), nil
}

type Login struct {
	Auth     string `json:"auth"`
	Token    string `json:"token"`
	Endpoint string `json:"endpoint"`
}

func (login Login) GetAuth() string {
	return login.Auth
}

func (login Login) GetToken() string {
	return login.Token
}

func (login Login) GetEndpoint() string {
	return login.Endpoint
}

func LoadLogin() *Login {
	var login Login
	loginJSON, err := ioutil.ReadFile(LoginFilePath)
	if err != nil {
		login = Login{}
	} else {
		err = json.Unmarshal(loginJSON, &login)
		if err != nil {
			login = Login{}
		}
	}

	// Overwrite endpoint with env var if exists
	ep := os.Getenv("BARCELONA_ENDPOINT")
	if len(ep) > 0 {
		login.Endpoint = ep
	}

	return &login
}

func writeLogin(login *Login) error {
	b, err := json.Marshal(login)
	if err != nil {
		return err
	}

	err = os.MkdirAll(ConfigDir, 0775)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(LoginFilePath, b, 0600)
	if err != nil {
		return err
	}
	return nil
}
