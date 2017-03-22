package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/degica/barcelona-cli/api"
	"github.com/urfave/cli"
)

var EnvCommand = cli.Command{
	Name:  "env",
	Usage: "Environment variable operations",
	Subcommands: []cli.Command{
		{
			Name:  "get",
			Usage: "Get environment variables",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "environment, e",
					Usage: "Environment of heritage",
				},
				cli.StringFlag{
					Name:  "heritage-name, H",
					Usage: "Heritage name",
				},
			},
			Action: func(c *cli.Context) error {
				envName := c.String("environment")
				heritageName := c.String("heritage-name")
				if len(envName) > 0 && len(heritageName) > 0 {
					return cli.NewExitError("environment and heritage-name are exclusive", 1)
				}
				if len(envName) > 0 {
					env, err := LoadEnvironment(c.String("environment"))
					if err != nil {
						return cli.NewExitError(err.Error(), 1)
					}
					heritageName = env.Name
				}

				resp, err := api.DefaultClient.Get("/heritages/"+heritageName, nil)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				var respHeritage api.HeritageResponse
				err = json.Unmarshal(resp, &respHeritage)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				printEnv(respHeritage.Heritage.EnvVars)

				return nil
			},
		},
		{
			Name:  "set",
			Usage: "Set environment variables",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "environment, e",
					Usage: "Environment of heritage",
				},
				cli.StringFlag{
					Name:  "heritage-name, H",
					Usage: "Heritage name",
				},
				cli.BoolFlag{
					Name:  "secret, s",
					Usage: "Save values as secret",
				},
			},
			ArgsUsage: "KEY1=VALUE1 [KEY2=VALUE2 ...]",
			Action: func(c *cli.Context) error {
				envName := c.String("environment")
				heritageName := c.String("heritage-name")
				if len(envName) > 0 && len(heritageName) > 0 {
					return cli.NewExitError("environment and heritage-name are exclusive", 1)
				}
				if len(envName) > 0 {
					env, err := LoadEnvironment(c.String("environment"))
					if err != nil {
						return cli.NewExitError(err.Error(), 1)
					}
					heritageName = env.Name
				}
				n := c.NArg()
				if n == 0 {
					return cli.NewExitError("Specify NAME=VALUE pairs", 1)
				}

				pairs := make(map[string]string)
				for i := 0; i < n; i++ {
					line := c.Args().Get(i)
					pair := strings.SplitN(line, "=", 2)
					pairs[pair[0]] = pair[1]
				}
				params := map[string]interface{}{
					"env_vars": pairs,
					"secret":   c.Bool("secret"),
				}
				j, err := json.Marshal(params)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}

				resp, err := api.DefaultClient.Post("/heritages/"+heritageName+"/env_vars", bytes.NewBuffer(j))
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				var respHeritage api.HeritageResponse
				err = json.Unmarshal(resp, &respHeritage)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				printEnv(respHeritage.Heritage.EnvVars)

				return nil
			},
		},
		{
			Name:      "unset",
			Usage:     "Unset environment variables",
			ArgsUsage: "KEY1 [KEY2 ...]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "environment, e",
					Usage: "Environment of heritage",
				},
				cli.StringFlag{
					Name:  "heritage-name, H",
					Usage: "Heritage name",
				},
			},
			Action: func(c *cli.Context) error {
				envName := c.String("environment")
				heritageName := c.String("heritage-name")
				if len(envName) > 0 && len(heritageName) > 0 {
					return cli.NewExitError("environment and heritage-name are exclusive", 1)
				}
				if len(envName) > 0 {
					env, err := LoadEnvironment(c.String("environment"))
					if err != nil {
						return cli.NewExitError(err.Error(), 1)
					}
					heritageName = env.Name
				}

				params := map[string]interface{}{
					"env_keys": c.Args(),
				}
				j, err := json.Marshal(params)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}

				resp, err := api.DefaultClient.Delete("/heritages/"+heritageName+"/env_vars", bytes.NewBuffer(j))
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				var respHeritage api.HeritageResponse
				err = json.Unmarshal(resp, &respHeritage)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				printEnv(respHeritage.Heritage.EnvVars)

				return nil
			},
		},
	},
}

func printEnv(es map[string]string) {
	for name, value := range es {
		fmt.Printf("%s: %s\n", name, value)
	}
}
