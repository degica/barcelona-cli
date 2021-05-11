package cmd

import (
	"bytes"
	"encoding/json"

	"github.com/degica/barcelona-cli/api"
	"github.com/urfave/cli"
)

var DeployCommand = cli.Command{
	Name:  "deploy",
	Usage: "Deploy a Barcelona heritage",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "environment, e",
			Usage: "Environment of heritage",
		},
		cli.StringFlag{
			Name:  "tag, t",
			Usage: "Tag of docker image",
		},
		cli.StringFlag{
			Name:  "heritage-token",
			Usage: "Heritage token",
		},
		cli.BoolTFlag{
			Name:  "quiet, q",
			Usage: "Do not print output if successful. By default it is true",
		},
	},
	Action: func(c *cli.Context) error {
		env := c.String("environment")
		tag := c.String("tag")
		token := c.String("heritage-token")
		quiet := c.Bool("quiet")

		var heritage *api.Heritage
		var err error
		if len(token) > 0 {
			heritage, err = doDeployWithHeritageToken(env, tag, token)
		} else {
			heritage, err = doDeploy(env, tag)
		}
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		if !quiet {
			heritage.Print()
		}

		return nil
	},
}

func doDeploy(env, tag string) (*api.Heritage, error) {
	h, err := LoadEnvironment(env)
	if err != nil {
		return nil, err
	}

	h.FillinDefaults()
	h.ImageTag = tag

	j, err := json.Marshal(h)
	if err != nil {
		return nil, err
	}

	resp, err := api.DefaultClient.Patch("/heritages/"+h.Name, bytes.NewBuffer(j))
	if err != nil {
		return nil, err
	}

	var hResp api.HeritageResponse
	err = json.Unmarshal(resp, &hResp)
	if err != nil {
		return nil, err
	}
	return hResp.Heritage, nil
}

func doDeployWithHeritageToken(env, tag, token string) (*api.Heritage, error) {
	h, err := LoadEnvironment(env)
	if err != nil {
		return nil, err
	}

	h.ImageTag = tag
	h.FillinDefaults()

	j, err := json.Marshal(h)
	if err != nil {
		return nil, err
	}

	resp, err := api.DefaultClient.Post("/heritages/"+h.Name+"/trigger/"+token, bytes.NewBuffer(j))
	if err != nil {
		return nil, err
	}

	var hResp api.HeritageResponse
	err = json.Unmarshal(resp, &hResp)
	if err != nil {
		return nil, err
	}
	return hResp.Heritage, nil
}
