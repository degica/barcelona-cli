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
var PrivateKeyPath string
var PublicKeyPath string
var CertPath string
var Debug bool

// This is the interface clients should use to extract config info
type Configuration interface {
	LoadLogin() *Login
}

// Clients should get configs using this function
func Get() Configuration {
	return &localConfig{}
}

// Implementation of our Configuration object
type localConfig struct{}

func (m localConfig) LoadLogin() *Login {
	return LoadLogin()
}

func init() {
	path, err := getConfigPath()
	if err != nil {
		panic("Couldn't get login path")
	}
	ConfigDir = path
	LoginFilePath = filepath.Join(ConfigDir, "login")
	PrivateKeyPath = filepath.Join(ConfigDir, "id_ecdsa")
	PublicKeyPath = filepath.Join(ConfigDir, "id_ecdsa.pub")
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
	Auth       string `json:"auth"`
	Token      string `json:"token"`
	VaultToken string `json:"vault_token"`
	Endpoint   string `json:"endpoint"`
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

func WriteLogin(login *Login) error {
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
