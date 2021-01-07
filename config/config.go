package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
)

var Debug bool

// Clients should get configs using this function
func Get() *LocalConfig {
	path, err := getConfigPath()
	if err != nil {
		panic("Couldn't get login path")
	}
	return &LocalConfig{
		configDir:      path,
		loginFilePath:  filepath.Join(path, "login"),
		privateKeyPath: filepath.Join(path, "id_ecdsa"),
		publicKeyPath:  filepath.Join(path, "id_ecdsa.pub"),
		certPath:       filepath.Join(path, "id_ecdsa-cert.pub"),
	}
}

// Implementation of our Configuration object
type LocalConfig struct {
	configDir      string
	loginFilePath  string
	privateKeyPath string
	publicKeyPath  string
	certPath       string
}

func (m LocalConfig) GetPrivateKeyPath() string {
	return m.privateKeyPath
}

func (m LocalConfig) GetPublicKeyPath() string {
	return m.publicKeyPath
}

func (m LocalConfig) GetCertPath() string {
	return m.certPath
}

func (m LocalConfig) GetConfigDir() string {
	return m.configDir
}

func (m LocalConfig) IsDebug() bool {
	return Debug
}

func (m LocalConfig) WriteLogin(auth string, token string, endpoint string) error {
	login := &Login{
		Auth:     auth,
		Token:    token,
		Endpoint: endpoint,
	}

	b, err := json.Marshal(login)
	if err != nil {
		return err
	}

	err = os.MkdirAll(m.configDir, 0775)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(m.loginFilePath, b, 0600)
	if err != nil {
		return err
	}
	return nil
}

func (m LocalConfig) LoadLogin() *Login {
	var login Login
	loginJSON, err := ioutil.ReadFile(m.loginFilePath)
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

func getConfigPath() (string, error) {
	homedir, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homedir, ".bcn"), nil
}
