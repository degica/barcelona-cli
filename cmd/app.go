package cmd

import (
	"github.com/degica/barcelona-cli/api"
	"github.com/degica/barcelona-cli/operations"
	"github.com/degica/barcelona-cli/utils"
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

				oper := operations.NewAppOperation(name, operations.Delete, c.Bool("no-confirmation"), api.DefaultClient, utils.NewStdinInputReader())
				return operations.Execute(oper)
			},
		},
		{
			Name:      "show",
			Usage:     "Show a heritage",
			ArgsUsage: "HERITAGE_NAME",
			Action: func(c *cli.Context) error {
				name := c.Args().Get(0)

				oper := operations.NewAppOperation(name, operations.Show, false, api.DefaultClient, utils.NewStdinInputReader())
				return operations.Execute(oper)
			},
		},
	},
}
