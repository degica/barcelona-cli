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

var EndpointCommand = cli.Command{
	Name:  "endpoint",
	Usage: "Endpoint operations",
	Subcommands: []cli.Command{
		{
			Name:      "create",
			Usage:     "Create a new endpoint",
			ArgsUsage: "ENDPOINT_NAME",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "district, d",
					Value: "default",
					Usage: "AWS region",
				},
				cli.BoolFlag{
					Name:  "public",
					Usage: "Public facing endpoint",
				},
				cli.StringFlag{
					Name:  "certificate-arn",
					Usage: "ACM Certificate ARN",
				},
				cli.StringFlag{
					Name:  "ssl-policy",
					Usage: "HTTPS SSL Policy",
				},
			},
			Action: func(c *cli.Context) error {
				endpointName := c.Args().Get(0)
				if len(endpointName) == 0 {
					return cli.NewExitError("endpoint name is required", 1)
				}

				public := c.Bool("public")
				request := api.Endpoint{
					Name:          endpointName,
					Public:        &public,
					CertificateID: c.String("certificate-arn"),
					SslPolicy:     c.String("ssl-policy"),
				}

				b, err := json.Marshal(&request)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}

				resp, err := api.DefaultClient.Request("POST", fmt.Sprintf("/districts/%s/endpoints", c.String("district")), bytes.NewBuffer(b))
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				var eResp api.EndpointResponse
				err = json.Unmarshal(resp, &eResp)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				printEndpoint(eResp.Endpoint)

				return nil
			},
		},
		{
			Name:      "show",
			Usage:     "Show endpoint information",
			ArgsUsage: "ENDPOINT_NAME",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "district, d",
					Value: "default",
					Usage: "District name",
				},
			},
			Action: func(c *cli.Context) error {
				endpointName := c.Args().Get(0)
				if len(endpointName) == 0 {
					return cli.NewExitError("endpoint name is required", 1)
				}

				resp, err := api.DefaultClient.Request("GET", fmt.Sprintf("/districts/%s/endpoints/%s", c.String("district"), endpointName), nil)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				var eResp api.EndpointResponse
				err = json.Unmarshal(resp, &eResp)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				printEndpoint(eResp.Endpoint)

				return nil
			},
		},
		{
			Name:  "list",
			Usage: "List endpoints",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "district, d",
					Value: "default",
					Usage: "District name",
				},
			},
			Action: func(c *cli.Context) error {
				resp, err := api.DefaultClient.Request("GET", fmt.Sprintf("/districts/%s/endpoints", c.String("district")), nil)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				var eResp api.EndpointResponse
				err = json.Unmarshal(resp, &eResp)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				printEndpoints(eResp.Endpoints)

				return nil
			},
		},
		{
			Name:      "update",
			Usage:     "Update an endpoint",
			ArgsUsage: "ENDPOINT_NAME",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "district, d",
					Value: "default",
					Usage: "District name",
				},
				cli.StringFlag{
					Name:  "certificate-arn",
					Usage: "ACM Certificate ARN",
				},
				cli.StringFlag{

					Usage: "HTTPS SSL Policy",
				},
			},
			Action: func(c *cli.Context) error {
				endpointName := c.Args().Get(0)
				if len(endpointName) == 0 {
					return cli.NewExitError("endpoint name is required", 1)
				}

				request := api.Endpoint{
					CertificateID: c.String("certificate-arn"),
					SslPolicy:     c.String("ssl-policy"),
				}

				b, err := json.Marshal(&request)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}

				resp, err := api.DefaultClient.Request("PATCH", fmt.Sprintf("/districts/%s/endpoints/%s", c.String("district"), endpointName), bytes.NewBuffer(b))
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				var eResp api.EndpointResponse
				err = json.Unmarshal(resp, &eResp)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				printEndpoint(eResp.Endpoint)

				return nil
			},
		},
		{
			Name:      "delete",
			Usage:     "Delete an endpoint",
			ArgsUsage: "ENDPOINT_NAME",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "district, d",
					Value: "default",
					Usage: "District name",
				},
				cli.BoolFlag{
					Name: "no-confirmation",
				},
			},
			Action: func(c *cli.Context) error {
				endpointName := c.Args().Get(0)
				if len(endpointName) == 0 {
					return cli.NewExitError("endpoint name is required", 1)
				}

				fmt.Printf("You are attempting to delete /%s/endpoints/%s\n", c.String("district"), endpointName)
				if !c.Bool("no-confirmation") && !areYouSure("This operation cannot be undone. Are you sure?") {
					return nil
				}

				_, err := api.DefaultClient.Request("DELETE", fmt.Sprintf("/districts/%s/endpoints/%s", c.String("district"), endpointName), nil)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				return nil
			},
		},
	},
}

func printEndpoint(e *api.Endpoint) {
	fmt.Printf("Name: %s\n", e.Name)
	fmt.Printf("Public: %t\n", *e.Public)
	fmt.Printf("SSL Policy: %s\n", e.SslPolicy)
	fmt.Printf("Certificate ARN: %s\n", e.CertificateID)
	fmt.Printf("DNS Name: %s\n", e.DNSName)
}

func printEndpoints(es []*api.Endpoint) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "District", "Public", "SSL Policy", "Cert ID"})
	table.SetBorder(false)
	for _, e := range es {
		table.Append([]string{e.Name, e.District.Name, fmt.Sprintf("%t", *e.Public), e.SslPolicy, e.CertificateID})
	}
	table.Render()
}
