package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/degica/barcelona-cli/api"
	"github.com/urfave/cli"
)

var AppCommand = cli.Command{
	Name:  "app",
	Usage: "Manage heritages",
	Subcommands: []cli.Command{
		{
			Name:      "delete",
			Usage:     "Delete a heritage",
			ArgsUsage: "HERITAGE_NAME",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "no-confirmation",
				},
			},
			Action: func(c *cli.Context) error {
				name := c.Args().Get(0)
				if len(name) == 0 {
					return cli.NewExitError("district name is required", 1)
				}

				fmt.Printf("You are attempting to delete %s\n", name)
				if !c.Bool("no-confirmation") && !areYouSure("This operation cannot be undone. Are you sure?") {
					return nil
				}

				_, err := api.DefaultClient.Delete("/heritages/"+name, nil)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				fmt.Printf("Deleted %s\n", name)

				return nil
			},
		},
		{
			Name:      "show",
			Usage:     "Show a heritage",
			ArgsUsage: "HERITAGE_NAME",
			Action: func(c *cli.Context) error {
				name := c.Args().Get(0)
				if len(name) == 0 {
					return cli.NewExitError("district name is required", 1)
				}

				resp, err := api.DefaultClient.Get("/heritages/"+name, nil)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				var hResp api.HeritageResponse
				err = json.Unmarshal(resp, &hResp)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				PrintHeritage(hResp.Heritage)

				return nil
			},
		},
	},
}
