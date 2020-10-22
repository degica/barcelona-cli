package cmd

import (
	"github.com/degica/barcelona-cli/api"
	"github.com/urfave/cli"
)

var CreateCommand = cli.Command{
	Name:  "create",
	Usage: "Create a new Barcelona heritage",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "environment, e",
			Usage: "Environment of heritage",
		},
		cli.StringFlag{
			Name:  "district, d",
			Value: "default",
			Usage: "District name",
		},
		cli.StringFlag{
			Name:  "tag, t",
			Value: "latest",
			Usage: "District name",
		},
	},
	Action: func(c *cli.Context) error {
		h, err := LoadEnvironment(c.String("environment"))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		h.FillinDefaults()
		h.ImageTag = c.String("tag")

		resp, err := api.DefaultClient.CreateHeritage(c.String("district"), h)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		resp.Print()

		return nil
	},
}
