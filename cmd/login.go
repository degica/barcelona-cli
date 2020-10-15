package cmd

import (
	"github.com/degica/barcelona-cli/api"
	"github.com/degica/barcelona-cli/config"
	"github.com/degica/barcelona-cli/operations"
	"github.com/degica/barcelona-cli/utils"
	"github.com/urfave/cli"
)

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
			*api.Client
			*config.LocalConfig
			*utils.CommandRunner
		}{
			utils.NewStdinInputReader(),
			api.DefaultClient,
			config.Get(),
			&utils.CommandRunner{},
		}

		oper := operations.NewLoginOperation(endpoint, backend, gh_token, vault_token, vault_url, ext)

		return operations.Execute(oper)
	},
	Subcommands: []cli.Command{
		{
			Name:  "info",
			Usage: "Show login information",
			Action: func(c *cli.Context) error {
				oper := operations.NewLoginInfoOperation(config.Get())
				return operations.Execute(oper)
			},
		},
	},
}
