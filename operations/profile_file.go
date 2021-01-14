package operations

import (
	"github.com/degica/barcelona-cli/config"
)

type profileFile struct {
	Name       string       `json:"name"`
	Login      config.Login `json:"login"`
	PrivateKey string       `json:"privateKey"`
	PublicKey  string       `json:"publicKey"`
	Cert       string       `json:"cert"`
}
