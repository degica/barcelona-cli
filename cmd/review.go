package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/degica/barcelona-cli/api"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

type DeployReviewApp struct {
	Request  *api.ReviewAppRequest
	Response *api.ReviewAppResponse
	Token    string
}

func (com *DeployReviewApp) Execute() error {
	j, err := json.Marshal(com.Request)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	var resp []byte
	if len(com.Token) == 0 {
		resp, err = api.DefaultClient.Post("/review_groups/"+com.Request.GroupName+"/apps", bytes.NewBuffer(j))
	} else {
		resp, err = api.DefaultClient.Post("/review_groups/"+com.Request.GroupName+"/ci/apps/"+com.Token, bytes.NewBuffer(j))
	}

	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	var rResp api.ReviewAppResponse
	err = json.Unmarshal(resp, &rResp)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	fmt.Printf("Domain: %s\n", rResp.ReviewApp.Domain)

	return nil
}

var ReviewCommand = cli.Command{
	Name:  "review",
	Usage: "Review Apps",
	Subcommands: []cli.Command{
		{
			Name:      "deploy",
			ArgsUsage: "SUBJECT",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "tag, t",
					Usage: "Tag of docker image",
				},
				cli.StringFlag{
					Name:  "token",
					Usage: "review group token",
				},
				cli.StringFlag{
					Name: "retention, r",
				},
			},
			Action: func(c *cli.Context) error {
				subject := c.Args().Get(0)
				tag := c.String("tag")
				token := c.String("token")
				retention := c.String("retention")
				var retentionSec int

				if len(retention) == 0 {
					retentionSec = 24 * 3600
				} else {
					d, err := time.ParseDuration(retention)
					if err != nil {
						return cli.NewExitError(err.Error(), 1)
					}
					retentionSec = int(d.Seconds())
				}

				reviewDef, err := LoadReviewDefinition()
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}

				com := DeployReviewApp{
					Request: &api.ReviewAppRequest{
						ReviewAppDefinition: reviewDef,
						Subject:             subject,
						Retention:           retentionSec,
						ImageTag:            tag,
					},
					Token: token,
				}

				return com.Execute()
			},
		},
		{
			Name:      "delete",
			ArgsUsage: "REVIEWAPP_NAME",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "group",
					Usage: "Review group name",
				},
				cli.StringFlag{
					Name:  "token",
					Usage: "review group token",
				},
			},
			Action: func(c *cli.Context) error {
				name := c.Args().Get(0)
				token := c.String("token")
				groupName := c.String("group")

				var err error
				if len(token) == 0 {
					_, err = api.DefaultClient.Delete("/review_groups/"+groupName+"/apps/"+name, nil)
				} else {
					_, err = api.DefaultClient.Delete("/review_groups/"+groupName+"/ci/apps/"+token+"/"+name, nil)

				}
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}

				return nil
			},
		},
		{
			Name: "list",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "group",
					Usage: "Review group name",
				},
			},
			Action: func(c *cli.Context) error {
				groupName := c.String("group")

				if len(groupName) == 0 {
					reviewDef, err := LoadReviewDefinition()
					if err != nil {
						return cli.NewExitError(err.Error(), 1)
					}
					groupName = reviewDef.GroupName
				}

				review_apps, err := getReviewApps(groupName)

				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}

				renderApps(review_apps)

				return nil
			},
		},
		{
			Name:  "run",
			Usage: "Run command inside Barcelona environment by branch name",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "B, branch",
					Usage: "Git branch name",
				},
			},
			Action: func(c *cli.Context) error {
				branchName := c.String("branch")
				groupName, err := getGroupName()

				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}

				review_apps, err := getReviewApps(groupName)

				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}

				heritageName := getHeritageName(branchName, review_apps)

				if heritageName == "" {
					return cli.NewExitError("heritage is not found", 1)
				}

				command := strings.Join(c.Args(), " ")
				params := map[string]interface{}{
					"interactive": true,
					"command":     command,
					"env_vars":    make(map[string]string),
				}

				err = connectToHeritage(params, heritageName, false)

				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}

				return nil
			},
		},
		ReviewGroupCommand,
	},
}

func getReviewApps(groupName string) ([]*api.ReviewApp, error) {
	resp, err := api.DefaultClient.Get("/review_groups/"+groupName+"/apps", nil)
	if err != nil {
		return nil, err
	}

	var appResp api.ReviewAppResponse
	err = json.Unmarshal(resp, &appResp)
	if err != nil {
		return nil, err
	}

	return appResp.ReviewApps, nil
}

func getHeritageName(branchName string, review_apps []*api.ReviewApp) string {
	heritageName := ""
	for _, app := range review_apps {
		if app.Subject == branchName {
			heritageName = app.Heritage.Name
			break
		}
	}

	return heritageName
}
func getGroupName() (string, error) {
	reviewDef, err := LoadReviewDefinition()
	if err != nil {
		return "", err
	}

	return reviewDef.GroupName, nil
}

func renderApps(apps []*api.ReviewApp) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Subject", "Domain", ""})
	table.SetBorder(false)
	for _, app := range apps {
		table.Append([]string{app.Subject, app.Domain, app.Heritage.Name})
	}
	table.Render()
}
