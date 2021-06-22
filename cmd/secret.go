package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/degica/barcelona-cli/api"
	"github.com/urfave/cli"
)

var SecretCommand = cli.Command{
	Name:  "secret",
	Usage: "Secret operations",

	Subcommands: []cli.Command{
		{
			Name: "add",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name, n",
					Usage: "Ssm parameter name",
				},
				cli.StringFlag{
					Name:  "value, v",
					Usage: "Ssm parameter value",
				},
				cli.StringFlag{
					Name:  "district, d",
					Usage: "District name",
				},
			},
			Action: func(c *cli.Context) error {
				parameterName := c.String("name")
				secretValue := c.String("value")
				district := c.String("district")

				params := make(map[string]interface{})
				params["name"] = parameterName
				params["value"] = secretValue

				j, err := json.Marshal(params)
				if err != nil {
					return err
				}

				_, err = api.DefaultClient.Post("/districts/"+district+"/ssm_parameters", bytes.NewBuffer(j))

				if err != nil {
					return err
				}

				fmt.Println("success to set " + parameterName)
				return nil
			},
		},
		{
			Name: "delete",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name, n",
					Usage: "Ssm parameter name",
				},
				cli.StringFlag{
					Name:  "district, d",
					Usage: "District name",
				},
			},
			Action: func(c *cli.Context) error {
				district := c.String("district")
				parameterName := url.QueryEscape(c.String("name"))

				resp, err := api.DefaultClient.Delete("/districts/"+district+"/ssm_parameters/"+parameterName, nil)

				if err != nil {
					return err
				}

				fmt.Println(string(resp))
				return nil
			},
		},
	},
}
