package cmd

import (
	"github.com/degica/barcelona-cli/api"
	"github.com/degica/barcelona-cli/config"
	"github.com/degica/barcelona-cli/operations"
	"github.com/degica/barcelona-cli/utils"
	"github.com/urfave/cli"
)

func AutoRefreshVaultToken(app *cli.App) {
	login := config.Get().LoadLogin()
	if login.Auth != "vault" || login.VaultUrl == "" || login.VaultToken == "" {
		return
	}

	args := []string{
		app.Name, "login", "--auth", login.Auth, "--vault-token",
		login.VaultToken, "--vault-url", login.VaultUrl, login.Endpoint}

	app.Run(args)
}

var LoginCommand = cli.Command{
	Name:      "login",
	Usage:     "Login Barcelona",
	ArgsUsage: "https://endpoint GITHUB_TOKEN",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "auth, a",
			Usage: "Auth backend",
			Value: "github",
		},
		cli.StringFlag{
			Name:  "github-token",
			Usage: "GitHub Token",
		},
		cli.StringFlag{
			Name:  "vault-token",
			Usage: "Vault Token",
		},
		cli.StringFlag{
			Name:  "vault-url",
			Usage: "Vault URL",
		},
	},
	Action: func(c *cli.Context) error {
		endpoint := c.Args().Get(0)
		backend := c.String("auth")
		gh_token := c.String("github-token")
		vault_token := c.String("vault-token")
		vault_url := c.String("vault-url")

		ext := struct {
			utils.UserInputReader
			*operations.ProxyLoginOperationClient
			*config.LocalConfig
			*utils.CommandRunner
			*utils.FileOps
		}{
			utils.NewStdinInputReader(),
			&operations.ProxyLoginOperationClient{Client: api.DefaultClient},
			config.Get(),
			&utils.CommandRunner{},
			&utils.FileOps{},
		}

		operation := operations.NewLoginOperation(endpoint, backend, gh_token, vault_token, vault_url, ext)
		return operations.Execute(operation)
	},
	Subcommands: []cli.Command{
		{
			Name:  "info",
			Usage: "Show login information",
			Action: func(c *cli.Context) error {
				operation := operations.NewLoginInfoOperation(config.Get().LoadLogin())
				return operations.Execute(operation)
			},
		},
	},
}
