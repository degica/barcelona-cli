package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"

	"github.com/degica/barcelona-cli/api"
	"github.com/degica/barcelona-cli/config"
	"github.com/urfave/cli"
)

var LoginCommand = cli.Command{
	Name:      "login",
	Usage:     "Login Barcelona",
	ArgsUsage: "https://endpoint GITHUB_TOKEN",
	Action: func(c *cli.Context) error {
		endpoint := c.Args().Get(0)
		if len(endpoint) == 0 {
			return cli.NewExitError("endpoint is required", 1)
		}
		token := c.Args().Get(1)
		if len(token) == 0 {
			return cli.NewExitError("token is required", 1)
		}
		user, err := api.DefaultClient.Login(endpoint, token)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		login := config.Login{
			Token:    user.Token,
			Endpoint: endpoint,
		}
		err = config.WriteLogin(&login)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		keyExists := fileExists(config.PublicKeyPath)
		if !keyExists {
			fmt.Println("Generating your SSH key pair...")
			cmd := exec.Command("ssh-keygen",
				"-t", "ecdsa",
				"-b", "521",
				"-f", config.PrivateKeyPath,
				"-C", "")
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
		}

		if !keyExists || len(user.PublicKey) == 0 {
			fmt.Println("Registering your public key...")

			pubKeyB, err := ioutil.ReadFile(config.PublicKeyPath)
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}

			re := regexp.MustCompile(" *\n$")
			pubKey := re.ReplaceAllString(string(pubKeyB), "")
			reqBody := make(map[string]string)
			reqBody["public_key"] = pubKey
			bodyB, err := json.Marshal(reqBody)
			err = api.ReloadDefaultClient()
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}

			_, err = api.DefaultClient.Patch("/user", bytes.NewBuffer(bodyB))
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
		}

		return nil
	},
	Subcommands: []cli.Command{
		{
			Name:  "info",
			Usage: "Show login information",
			Action: func(c *cli.Context) error {
				login := config.LoadLogin()

				fmt.Printf("Endpoint: %s\n", login.Endpoint)
				return nil
			},
		},
	},
}
