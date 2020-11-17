package config

import ()

type Login struct {
	auth     string `json:"auth"`
	token    string `json:"token"`
	endpoint string `json:"endpoint"`
}

func (login Login) GetAuth() string {
	return login.auth
}

func (login Login) GetToken() string {
	return login.token
}

func (login Login) GetEndpoint() string {
	return login.endpoint
}
