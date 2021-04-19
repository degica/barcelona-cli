package config

type Login struct {
	Auth       string `json:"auth"`
	Token      string `json:"token"`
	Endpoint   string `json:"endpoint"`
	VaultUrl   string `json:"vault_url"`
	VaultToken string `json:"vault_token"`
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
