package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/degica/barcelona-cli/api"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

var ReviewGroupCommand = cli.Command{
	Name: "group",
	Subcommands: []cli.Command{
		{
			Name:      "create",
			ArgsUsage: "GROUP_NAME",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "base-domain",
					Usage: "Base Domain Name",
				},
				cli.StringFlag{
					Name:  "endpoint",
					Usage: "Endpoint name",
				},
			},
			Action: func(c *cli.Context) error {
				name := c.Args().Get(0)
				if len(name) == 0 {
					return cli.NewExitError("Group name is required", 1)
				}

				req := api.ReviewGroupRequest{
					Name:         name,
					BaseDomain:   c.String("base-domain"),
					EndpointName: c.String("endpoint"),
				}

				j, err := json.Marshal(req)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}

				resp, err := api.DefaultClient.Post("/review_groups", bytes.NewBuffer(j))
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}

				var rResp api.ReviewGroupResponse
				err = json.Unmarshal(resp, &rResp)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}

				return nil
			},
		},
		{
			Name:      "show",
			ArgsUsage: "GROUP_NAME",
			Action: func(c *cli.Context) error {
				groupName := c.Args().Get(0)
				resp, err := api.DefaultClient.Get("/review_groups/"+groupName, nil)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}

				var rResp api.ReviewGroupResponse
				err = json.Unmarshal(resp, &rResp)
				group := rResp.ReviewGroup

				fmt.Println("Name: ", group.Name)
				fmt.Println("Base Domain: ", group.BaseDomain)
				fmt.Println("Endpoint: ", group.Endpoint.Name)
				fmt.Println("Token: ", *group.Token)
				fmt.Println("Apps")

				return nil
			},
		},
		{
			Name: "list",
			Action: func(c *cli.Context) error {
				resp, err := api.DefaultClient.Get("/review_groups", nil)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}

				var rResp api.ReviewGroupResponse
				err = json.Unmarshal(resp, &rResp)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				groups := rResp.ReviewGroups

				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Name", "Base Domain", "Endpoint"})
				table.SetBorder(false)
				for _, g := range groups {
					table.Append([]string{g.Name, g.BaseDomain, g.Endpoint.Name})
				}
				table.Render()

				return nil
			},
		},
		{
			Name:      "delete",
			ArgsUsage: "GROUP_NAME",
			Action: func(c *cli.Context) error {
				groupName := c.Args().Get(0)
				_, err := api.DefaultClient.Delete("/review_groups/"+groupName, nil)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}

				return nil
			},
		},
	},
}
