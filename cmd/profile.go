package cmd

import (
	"strings"
	"github.com/degica/barcelona-cli/config"
	"github.com/degica/barcelona-cli/operations"
	"github.com/degica/barcelona-cli/utils"
	"github.com/urfave/cli"
)

func profileSubcommands(opnames []string) []cli.Command {
	var array []cli.Command

  for _, opname := range opnames {
  	subcommand := cli.Command{
			Name:      opname,
			Usage:     strings.Title(opname) + " a profile",
			ArgsUsage: "PROFILE_NAME",
			Flags: []cli.Flag{},
			Action: func(c *cli.Context) error {
				name := c.Args().Get(0)

				ext := struct {
					utils.UserInputReader
					*config.LocalConfig
					*utils.FileOps
				}{
					utils.NewStdinInputReader(),
					config.Get(),
					&utils.FileOps{},
				}

				oper := operations.NewProfileOperation(c.Command.Name, name, ext)
				return operations.Execute(oper)
			},
  	}
    array = append(array, subcommand)
  }
	return array
}

var ProfileCommand = cli.Command{
	Name:  "profile",
	Usage: "Manage profiles",
	Subcommands: profileSubcommands([]string{
		"create",
		"delete",
		"use",
		"show",
	}),
}
