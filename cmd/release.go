package cmd

import "github.com/urfave/cli"

var ReleaseCommand = cli.Command{
	Name:  "release",
	Usage: "Manipulate releases",
	Subcommands: []cli.Command{
		{
			Name:  "list",
			Usage: "List releases",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
		{
			Name:  "rollback",
			Usage: "Roll back to the previous release",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
	},
}
