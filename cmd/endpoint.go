package cmd

import (
	"github.com/degica/barcelona-cli/api"
	"github.com/degica/barcelona-cli/operations"
	"github.com/degica/barcelona-cli/utils"
	"github.com/urfave/cli"
)

func getEndpointSubcommands() []cli.Command {
	districtFlag := cli.StringFlag{
		Name:  "district, d",
		Value: "default",
		Usage: "AWS region",
	}

	publicFlag := cli.BoolFlag{
		Name:  "public",
		Usage: "Public facing endpoint",
	}

	certificateFlag := cli.StringFlag{
		Name:  "certificate-arn",
		Usage: "ACM Certificate ARN",
	}

	sslPolicyFlag := cli.StringFlag{
		Name:  "ssl-policy",
		Usage: "HTTPS SSL Policy",
	}

	noConfirmFlag := cli.BoolFlag{
		Name: "no-confirmation",
	}

	endpointSubcommands := map[operations.OperationType]cli.Command{
		operations.Create: cli.Command{
			Name:      "create",
			Usage:     "Create a new endpoint",
			ArgsUsage: "ENDPOINT_NAME",
			Flags: []cli.Flag{
				districtFlag,
				publicFlag,
				certificateFlag,
				sslPolicyFlag,
			},
		},

		operations.Update: cli.Command{
			Name:      "update",
			Usage:     "Update an endpoint",
			ArgsUsage: "ENDPOINT_NAME",
			Flags: []cli.Flag{
				districtFlag,
				certificateFlag,
				sslPolicyFlag,
			},
		},

		operations.Delete: cli.Command{
			Name:      "delete",
			Usage:     "Delete an endpoint",
			ArgsUsage: "ENDPOINT_NAME",
			Flags: []cli.Flag{
				districtFlag,
				noConfirmFlag,
			},
		},

		operations.Show: cli.Command{
			Name:      "show",
			Usage:     "Show endpoint information",
			ArgsUsage: "ENDPOINT_NAME",
			Flags: []cli.Flag{
				districtFlag,
			},
		},

		operations.List: cli.Command{
			Name:  "list",
			Usage: "List endpoints",
			Flags: []cli.Flag{
				districtFlag,
			},
		},
	}

	var array []cli.Command

	for opname, command := range endpointSubcommands {

		captured_opname := opname
		captured_command := command

		subcommand := cli.Command{
			Name:      command.Name,
			Usage:     command.Usage,
			ArgsUsage: command.ArgsUsage,
			Flags:     command.Flags,
			Action: func(c *cli.Context) error {

				endpointName := c.Args().Get(0)
				if captured_command.ArgsUsage == "ENDPOINT_NAME" {
					if len(endpointName) == 0 {
						return cli.NewExitError("endpoint name is required", 1)
					}

					if len(c.Args()) != 1 {
						return cli.NewExitError("please place options and flags before the endpoint name.", 1)
					}
				}

				public := c.Bool("public")
				cert_arn := c.String("certificate-arn")
				policy := c.String("ssl-policy")
				districtName := c.String("district")
				noconfirm := c.Bool("no-confirmation")

				client := struct {
					*api.Client
					utils.UserInputReader
				}{
					api.DefaultClient,
					utils.NewStdinInputReader(),
				}

				oper := operations.NewEndpointOperation(
					districtName,
					endpointName,
					public,
					cert_arn,
					policy,
					noconfirm,
					captured_opname,
					client,
				)
				return operations.Execute(oper)
			},
		}
		array = append(array, subcommand)
	}

	return array
}

var EndpointCommand = cli.Command{
	Name:        "endpoint",
	Usage:       "Endpoint operations",
	Subcommands: getEndpointSubcommands(),
}
