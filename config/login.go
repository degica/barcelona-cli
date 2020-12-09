package config

import ()

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
