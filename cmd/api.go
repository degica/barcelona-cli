package cmd

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/degica/barcelona-cli/api"
	"github.com/urfave/cli"
)

var APICommand = cli.Command{
	Name:      "api",
	Usage:     "Call Barcelona API",
	ArgsUsage: "METHOD PATH [BODY]",
	Action: func(c *cli.Context) error {
		method := strings.ToUpper(c.Args().Get(0))
		if len(method) == 0 {
			return cli.NewExitError("method is required", 1)
		}
		path := c.Args().Get(1)
		if len(path) == 0 {
			return cli.NewExitError("path is required", 1)
		}
		body := bytes.NewBufferString(c.Args().Get(2))
		b, err := api.DefaultClient.Request(method, path, body)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		fmt.Println(PrettyJSON(b))

		return nil
	},
}
