package cmd

import (
	"bytes"
	"strings"

	"github.com/degica/barcelona-cli/api"
	"github.com/degica/barcelona-cli/operations"
	"github.com/urfave/cli"
)

var APICommand = cli.Command{
	Name:      "api",
	Usage:     "Call Barcelona API",
	ArgsUsage: "METHOD PATH [BODY]",
	Action: func(c *cli.Context) error {
		method := strings.ToUpper(c.Args().Get(0))
		path := c.Args().Get(1)
		body := bytes.NewBufferString(c.Args().Get(2))

		oper := operations.NewApiOperation(method, path, body, api.DefaultClient)
		return operations.Execute(oper)
	},
}
