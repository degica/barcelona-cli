package cmd

import (
	"github.com/degica/barcelona-cli/api"
	"github.com/degica/barcelona-cli/config"
	"github.com/degica/barcelona-cli/operations"
	"github.com/degica/barcelona-cli/utils"
	"github.com/urfave/cli"
)

var SSHCommand = cli.Command{
	Name:      "ssh",
	Usage:     "SSH into Barcelona container instance",
	ArgsUsage: "DISTRICT_NAME CONTAINER_INSTANCE_PRIVATE_IP",
	Action: func(c *cli.Context) error {
		districtName := c.Args().Get(0)
		ip := c.Args().Get(1)
		oper := operations.NewSshcmdOperation(
			api.DefaultClient,
			districtName,
			ip,
			config.Get(),
			&utils.CommandRunner{},
		)
		return operations.Execute(oper)
	},
}
