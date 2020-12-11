package operations

import (
	"github.com/degica/barcelona-cli/config"
)

type ProfileFile struct {
	name       string       `json:"name"`
	login      config.Login `json:"login"`
	privateKey string       `json:"privateKey"`
	publicKey  string       `json:"publicKey"`
	cert       string       `json:"cert"`
}
