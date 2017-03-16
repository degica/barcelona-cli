package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/degica/barcelona-cli/api"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

var DistrictCommand = cli.Command{
	Name:  "district",
	Usage: "District operations",
	Subcommands: []cli.Command{
		{
			Name:      "create",
			Usage:     "Create a new district",
			ArgsUsage: "DISTRICT_NAME",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "region",
					Value: "us-east-1",
					Usage: "AWS region",
				},
				cli.StringFlag{
					Name:  "nat-type",
					Value: "instance",
					Usage: "NAT type",
				},
				cli.StringFlag{
					Name:  "cluster-instance-type",
					Value: "t2.small",
					Usage: "Cluster Instance Type",
				},
			},
			Action: func(c *cli.Context) error {
				districtName := c.Args().Get(0)
				if len(districtName) == 0 {
					return cli.NewExitError("district name is required", 1)
				}

				size := 1
				request := api.DistrictRequest{
					Name:                districtName,
					Region:              c.String("region"),
					NatType:             c.String("nat-type"),
					ClusterSize:         &size,
					ClusterInstanceType: c.String("cluster-instance-type"),
					ClusterBackend:      "autoscaling",
				}
				request.AwsAccessKeyId = ask("AWS Access Key ID", true, false)
				request.AwsSecretAccessKey = ask("AWS Secret Access Key", true, true)

				district, err := api.DefaultClient.CreateDistrict(&request)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}

				printDistrict(district)

				return nil
			},
		},
		{
			Name:  "list",
			Usage: "List Districts",
			Action: func(c *cli.Context) error {
				districts, err := api.DefaultClient.ListDistricts()
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				printDistricts(districts)

				return nil
			},
		},
		{
			Name:      "show",
			Usage:     "Show District Information",
			ArgsUsage: "DISTRICT_NAME",
			Action: func(c *cli.Context) error {
				districtName := c.Args().Get(0)
				if len(districtName) == 0 {
					return cli.NewExitError("district name is required", 1)
				}

				district, err := api.DefaultClient.ShowDistrict(districtName)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				printDistrict(district)

				return nil
			},
		},
		{
			Name:      "update",
			Usage:     "Update District Information",
			ArgsUsage: "DISTRICT_NAME",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "nat-type",
					Usage: "NAT type",
				},
				cli.StringFlag{
					Name:  "cluster-instance-type",
					Usage: "Cluster Instance Type",
				},
				cli.IntFlag{
					Name:  "cluster-size",
					Value: -1,
					Usage: "Cluster Instance Type",
				},
				cli.BoolFlag{
					Name:  "apply",
					Usage: "Apply immediately",
				},
			},
			Action: func(c *cli.Context) error {
				districtName := c.Args().Get(0)
				if len(districtName) == 0 {
					return cli.NewExitError("district name is required", 1)
				}

				request := api.DistrictRequest{
					Name:                districtName,
					NatType:             c.String("nat-type"),
					ClusterInstanceType: c.String("cluster-instance-type"),
				}
				if size := c.Int("cluster-size"); size >= 0 {
					request.ClusterSize = &size
				}

				district, err := api.DefaultClient.UpdateDistrict(&request)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				printDistrict(district)

				err = applyOrNotice(districtName, c.Bool("apply"))
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}

				return nil
			},
		},
		{
			Name:      "apply",
			Usage:     "apply district stack",
			ArgsUsage: "DISTRICT_NAME",
			Action: func(c *cli.Context) error {
				districtName := c.Args().Get(0)
				if len(districtName) == 0 {
					return cli.NewExitError("district name is required", 1)
				}

				err := applyOrNotice(districtName, true)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}

				return nil
			},
		},
		{
			Name:      "delete",
			Usage:     "Delete a district",
			ArgsUsage: "DISTRICT_NAME",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "no-confirmation",
				},
			},
			Action: func(c *cli.Context) error {
				districtName := c.Args().Get(0)
				if len(districtName) == 0 {
					return cli.NewExitError("district name is required", 1)
				}

				fmt.Printf("You are attempting to delete %s\n", districtName)
				if !c.Bool("no-confirmation") && !areYouSure("This operation cannot be undone. Are you sure?") {
					return nil
				}

				err := api.DefaultClient.DeleteDistrict(districtName)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				return nil
			},
		},
		{
			Name:      "put-plugin",
			Usage:     "Add or Update plugin configuration",
			ArgsUsage: "DISTRICT_NAME PLUGIN_NAME",
			Flags: []cli.Flag{
				cli.StringSliceFlag{
					Name:  "attribute, a",
					Usage: "ATTR_NAME=VALUE",
				},
				cli.BoolFlag{
					Name:  "apply",
					Usage: "Apply immediately",
				},
			},
			Action: func(c *cli.Context) error {
				districtName := c.Args().Get(0)
				if len(districtName) == 0 {
					return cli.NewExitError("district name is required", 1)
				}

				pluginName := c.Args().Get(1)
				if len(districtName) == 0 {
					return cli.NewExitError("plugin name is required", 1)
				}

				req := api.Plugin{
					Name: pluginName,
				}
				req.Attributes = make(map[string]string)
				attrs := c.StringSlice("attribute")
				for _, s := range attrs {
					ss := strings.SplitN(s, "=", 2)
					req.Attributes[ss[0]] = ss[1]
				}

				plugin, err := api.DefaultClient.PutPlugin(districtName, &req)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				printPlugin(plugin)

				err = applyOrNotice(districtName, c.Bool("apply"))
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}

				return nil
			},
		},
		{
			Name:      "delete-plugin",
			Usage:     "Delete a plugin",
			ArgsUsage: "DISTRICT_NAME PLUGIN_NAME",
			Action: func(c *cli.Context) error {
				districtName := c.Args().Get(0)
				if len(districtName) == 0 {
					return cli.NewExitError("district name is required", 1)
				}

				pluginName := c.Args().Get(1)
				if len(districtName) == 0 {
					return cli.NewExitError("plugin name is required", 1)
				}

				err := api.DefaultClient.DeletePlugin(districtName, pluginName)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				return nil
			},
		},
	},
}

func applyOrNotice(districtName string, apply bool) error {
	if apply {
		err := api.DefaultClient.ApplyDistrict(districtName)
		if err != nil {
			return err
		}
		fmt.Println("Applying network stack")
	} else {
		fmt.Println("The change has not been applied to the hosts.")
		fmt.Println("Run `bcn district apply` to apply the change")
	}

	return nil
}

func printPlugin(p *api.Plugin) {
	fmt.Printf("Name %s\n", p.Name)
	for k, v := range p.Attributes {
		fmt.Printf("%s: %s\n", k, v)
	}
}

func printDistrict(d *api.District) {
	fmt.Printf("Name: %s\n", d.Name)
	fmt.Printf("Region: %s\n", d.Region)
	fmt.Printf("Cluster Backend: %s\n", d.ClusterBackend)
	fmt.Printf("Cluster Instance Type: %s\n", d.ClusterInstanceType)
	fmt.Printf("Cluster Size: %d\n", d.ClusterSize)
	fmt.Printf("S3 Bucket Name: %s\n", d.S3BucketName)
	fmt.Printf("Stack Name: %s\n", d.StackName)
	fmt.Printf("Stack Status: %s\n", d.StackStatus)
	fmt.Printf("NAT Type: %s\n", d.NatType)
	fmt.Printf("CIDR Block: %s\n", d.CidrBlock)
	fmt.Printf("AWS Access Key ID: %s\n", d.AwsAccessKeyId)
	fmt.Printf("AWS Role: %s\n", d.AwsRole)
	fmt.Printf("Container Instances:\n")
	for _, ci := range d.ContainerInstances {
		fmt.Printf("  %s %s %s\n", ci.EC2InstanceID, ci.PrivateIPAddress, ci.Status)
	}

	fmt.Printf("Heritages:\n")
	for _, h := range d.Heritages {
		fmt.Printf("  %s\n", h.Name)
	}

	fmt.Printf("Notifications:\n")
	for _, n := range d.Notifications {
		fmt.Printf("  %d %s %s\n", n.ID, n.Target, n.Endpoint)
	}

	fmt.Printf("Plugins:\n")
	for _, plugin := range d.Plugins {
		attrs := ""
		for k, v := range plugin.Attributes {
			attrs += fmt.Sprintf("%s=%s ", k, v)
		}
		fmt.Printf("  %s: %s\n", plugin.Name, attrs)
	}
}

func printDistricts(ds []*api.District) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Region", "Instance Type", "Cluster Size", "AWS Role", "Access Key ID"})
	table.SetBorder(false)
	for _, d := range ds {
		table.Append([]string{d.Name, d.Region, d.ClusterInstanceType, fmt.Sprintf("%d", d.ClusterSize), d.AwsRole, d.AwsAccessKeyId})
	}
	table.Render()
}
