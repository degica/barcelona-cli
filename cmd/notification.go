package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/degica/barcelona-cli/api"
	"github.com/urfave/cli"
)

var NotificationCommand = cli.Command{
	Name:  "notification",
	Usage: "Operate notification",
	Subcommands: []cli.Command{
		{
			Name:  "create",
			Usage: "Create a new notification",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "district, d",
					Value: "default",
					Usage: "District name",
				},
				cli.StringFlag{
					Name:  "target",
					Usage: "Notification Target",
				},
				cli.StringFlag{
					Name:  "endpoint",
					Usage: "Notification Endpoint",
				},
			},
			Action: func(c *cli.Context) error {
				target := c.String("target")
				if len(target) == 0 {
					return cli.NewExitError("target is required", 1)
				}
				endpoint := c.String("endpoint")
				if len(endpoint) == 0 {
					return cli.NewExitError("endpoint is required", 1)
				}

				request := api.Notification{
					Target:   target,
					Endpoint: endpoint,
				}

				b, err := json.Marshal(&request)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}

				resp, err := api.DefaultClient.Request("POST", fmt.Sprintf("/districts/%s/notifications", c.String("district")), bytes.NewBuffer(b))
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				var eResp api.NotificationResponse
				err = json.Unmarshal(resp, &eResp)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				printNotification(eResp.Notification)

				return nil
			},
		},
		{
			Name:      "show",
			Usage:     "Show notification",
			ArgsUsage: "ID",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "district, d",
					Value: "default",
					Usage: "District name",
				},
			},
			Action: func(c *cli.Context) error {
				id, err := parseIDArg(c)
				if err != nil {
					return err
				}

				resp, err := api.DefaultClient.Request("GET", fmt.Sprintf("/districts/%s/notifications/%d", c.String("district"), id), nil)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				var eResp api.NotificationResponse
				err = json.Unmarshal(resp, &eResp)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				printNotification(eResp.Notification)

				return nil
			},
		},
		{
			Name:      "update",
			Usage:     "Update notification",
			ArgsUsage: "ID",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "district, d",
					Value: "default",
					Usage: "District name",
				},
				cli.StringFlag{
					Name:  "target",
					Usage: "Notification Target",
				},
				cli.StringFlag{
					Name:  "endpoint",
					Usage: "Notification Endpoint",
				},
			},
			Action: func(c *cli.Context) error {
				id, err := parseIDArg(c)
				if err != nil {
					return err
				}

				request := api.Notification{
					ID:       id,
					Target:   c.String("target"),
					Endpoint: c.String("endpoint"),
				}

				b, err := json.Marshal(&request)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				fmt.Println(string(b))

				resp, err := api.DefaultClient.Request("PATCH", fmt.Sprintf("/districts/%s/notifications/%d", c.String("district"), id), bytes.NewBuffer(b))
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				var eResp api.NotificationResponse
				err = json.Unmarshal(resp, &eResp)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				printNotification(eResp.Notification)

				return nil
			},
		},
		{
			Name:      "delete",
			Usage:     "Delete notification",
			ArgsUsage: "ID",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "district, d",
					Value: "default",
					Usage: "District name",
				},
			},
			Action: func(c *cli.Context) error {
				id, err := parseIDArg(c)
				if err != nil {
					return err
				}

				_, err = api.DefaultClient.Request("DELETE", fmt.Sprintf("/districts/%s/notifications/%d", c.String("district"), id), nil)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				return nil
			},
		},
	},
}

func parseIDArg(c *cli.Context) (int, error) {
	idStr := c.Args().Get(0)
	if len(idStr) == 0 {
		return 0, cli.NewExitError("Notification ID is required", 1)
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, cli.NewExitError(err.Error(), 1)
	}
	return id, nil
}

func printNotification(n *api.Notification) {
	fmt.Printf("ID:       %d\n", n.ID)
	fmt.Printf("Target:   %s\n", n.Target)
	fmt.Printf("Endpoint: %s\n", n.Endpoint)
}
