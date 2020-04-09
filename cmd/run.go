package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
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
		cli.BoolFlag{
			Name:  "detach, D",
			Usage: "Detach mode",
		},
		cli.StringSliceFlag{
			Name:  "envvar, E",
			Usage: "Environment variable to pass to task",
		},
	},
	Action: func(c *cli.Context) error {
		envName := c.String("environment")
		heritageName := c.String("heritage-name")
		detach := c.Bool("detach")
		envVars := c.StringSlice("envvar")
		envVarMap, loadEnvVarMapErr := loadEnvVars(envName, heritageName)
		if loadEnvVarMapErr != nil {
			return cli.NewExitError(loadEnvVarMapErr.Error(), 1)
		}
		if len(envName) > 0 && len(heritageName) > 0 {
			return cli.NewExitError("environment and heritage-name are exclusive", 1)
		}
		if len(envVars) > 0 {
			varmap, err := checkEnvVars(envVars)
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
			for k, v := range varmap {
				envVarMap[k] = v
			}
		}
		if len(c.Args()) == 0 {
			return cli.NewExitError("Command is required", 1)
		}
		command := strings.Join(c.Args(), " ")
		params := map[string]interface{}{
			"interactive": !detach,
			"command":     command,
			"env_vars":    envVarMap,
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

		if detach {
			PrintOneoff(oneoff)
			return nil
		}

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

func loadEnvVars(envName string, heritageName string) (map[string]string, error) {
	result := make(map[string]string)
	if len(envName) > 0 {
		env, err := LoadEnvironment(envName)
		if err != nil {
			return nil, err
		}
		if env.RunEnv != nil {
			for k, v := range env.RunEnv.Vars {
				result[k] = v
			}
		}
	}
	return result, nil
}

func checkEnvVars(envvarSlice []string) (map[string]string, error) {
	var result = make(map[string]string)

	re := regexp.MustCompile(`^([A-Z_]+)=(.*)$`)
	for _, envvar := range envvarSlice {
		if !re.Match([]byte(envvar)) {
			return nil, errors.New(fmt.Sprintf("Env Variable  %s  is not valid. Name must have PASCAL_CASE=", envvar))
		}
		result[re.FindStringSubmatch(envvar)[1]] = re.FindStringSubmatch(envvar)[2]
	}
	return result, nil
}
