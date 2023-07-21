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
	"github.com/degica/barcelona-cli/config"
	"github.com/degica/barcelona-cli/utils"
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
		cli.StringFlag{
			Name:  "b, branch",
			Usage: "Git branch name",
		},
	},
	Action: func(c *cli.Context) error {
		envName := c.String("environment")
		heritageName := c.String("heritage-name")
		branchName := c.String("branch")

		detach := c.Bool("detach")
		envVars := c.StringSlice("envvar")
		envVarMap, loadEnvVarMapErr := loadEnvVars(envName)

		if loadEnvVarMapErr != nil {
			return cli.NewExitError(loadEnvVarMapErr.Error(), 1)
		}

		if len(envName) > 0 && len(heritageName) > 0 {
			return cli.NewExitError("environment and heritage-name are exclusive", 1)
		}

		if len(branchName) > 0 {
			if len(envName) > 0 || len(heritageName) > 0 {
				return cli.NewExitError("environment, heritage-name and branch-name are exclusive", 1)
			}

			name, err := getHeritageName(branchName)
			if err != nil {
				return err
			}

			heritageName = name
		}

		if len(heritageName) == 0 {
			env, err := LoadEnvironment(envName)
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
			heritageName = env.Name
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
		err := connectToHeritage(params, heritageName, detach)

		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		return nil
	},
}

func loadEnvVars(envName string) (map[string]string, error) {
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

func connectToHeritage(params map[string]interface{}, heritageName string, detach bool) error {
	j, err := json.Marshal(params)

	if err != nil {
		return err
	}

	resp, err := api.DefaultClient.Post("/heritages/"+heritageName+"/oneoffs", bytes.NewBuffer(j))
	if err != nil {
		return err
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
		path := fmt.Sprintf("/districts/%s/heritages/%s/oneoffs/%d", oneoff.District.Name, heritageName, oneoff.ID)
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

	if matchedCI == nil {
		return cli.NewExitError("Failed to find the container. Maybe try again")
	}

	ssh := utils.NewSshCommand(
		matchedCI.PrivateIPAddress,
		oneoff.District.BastionIP,
		certificate,
		config.Get(),
		&utils.CommandRunner{},
	)

	if ssh.Run(oneoff.InteractiveRunCommand) != nil && err != nil {
		return err
	}

	return nil
}

func getHeritageName(branchName string) (string, error) {
	groupName, err := getGroupName()

	if err != nil {
		return "", err
	}

	review_apps, err := getReviewApps(groupName)

	if err != nil {
		return "", err
	}

	heritageName := ""
	for _, app := range review_apps {
		if app.Subject == branchName {
			heritageName = app.Heritage.Name
			break
		}
	}

	if heritageName == "" {
		return "", errors.New(fmt.Sprintf("No heritage found for branch: %s", branchName))
	}

	return heritageName, nil
}

func getGroupName() (string, error) {
	reviewDef, err := LoadReviewDefinition()
	if err != nil {
		return "", err
	}

	return reviewDef.GroupName, nil
}
