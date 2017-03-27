package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/degica/barcelona-cli/api"
	"github.com/urfave/cli"
)

var RunCommand = cli.Command{
	Name:      "run",
	Usage:     "Run command inside Barcelona environment",
	ArgsUsage: "COMMAND...",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "environment, e",
			Usage: "Environment of heritage",
		},
		cli.StringFlag{
			Name:  "heritage-name, H",
			Usage: "Heritage name",
		},
		cli.IntFlag{
			Name:  "memory, m",
			Usage: "Memory size in MB",
		},
		cli.StringFlag{
			Name:  "user, u",
			Usage: "User name",
		},
	},
	Action: func(c *cli.Context) error {
		envName := c.String("environment")
		heritageName := c.String("heritage-name")
		if len(envName) > 0 && len(heritageName) > 0 {
			return cli.NewExitError("environment and heritage-name are exclusive", 1)
		}
		if len(envName) > 0 {
			env, err := LoadEnvironment(c.String("environment"))
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
			heritageName = env.Name
		}
		if len(c.Args()) == 0 {
			return cli.NewExitError("Command is required", 1)
		}
		command := strings.Join(c.Args(), " ")
		params := map[string]interface{}{
			"interactive": true,
			"command":     command,
		}
		memory := c.Int("memory")
		if memory > 0 {
			params["memory"] = memory
		}

		user := c.String("user")
		if user != "" {
			params["user"] = user
		}
		j, err := json.Marshal(params)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		resp, err := api.DefaultClient.Post("/heritages/"+heritageName+"/oneoffs", bytes.NewBuffer(j))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		var respOneoff api.OneoffResponse
		err = json.Unmarshal(resp, &respOneoff)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		oneoff := respOneoff.Oneoff
		certificate := respOneoff.Certificate

		fmt.Println("Waiting for the process to start")

	LOOP:
		for {
			path := fmt.Sprintf("/oneoffs/%d", oneoff.ID)
			resp, err := api.DefaultClient.Get(path, nil)
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
			var respOneoff api.OneoffResponse
			err = json.Unmarshal(resp, &respOneoff)
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
			switch respOneoff.Oneoff.Status {
			case "RUNNING":
				break LOOP
			case "PENDING":
				time.Sleep(3 * time.Second)
			default:
				// INACTIVE or STOPPED
				return cli.NewExitError("Unexpected task status "+respOneoff.Oneoff.Status, 1)
			}
		}

		fmt.Println("Connecting to the process")

		var matchedCI *api.ContainerInstance
		for _, ci := range oneoff.District.ContainerInstances {
			if ci.ContainerInstanceArn == oneoff.ContainerInstanceARN {
				matchedCI = ci
				break
			}
		}

		ssh := SSH{
			IP:          matchedCI.PrivateIPAddress,
			BastionIP:   oneoff.District.BastionIP,
			Certificate: certificate,
		}
		if ssh.Run(oneoff.InteractiveRunCommand) != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		return nil
	},
}
