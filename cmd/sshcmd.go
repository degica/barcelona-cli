package cmd

import (
	"encoding/json"

	"github.com/degica/barcelona-cli/api"
	"github.com/urfave/cli"
)

var SSHCommand = cli.Command{
	Name:      "ssh",
	Usage:     "SSH into Barcelona container instance",
	ArgsUsage: "DISTRICT_NAME CONTAINER_INSTANCE_PRIVATE_IP",
	Action: func(c *cli.Context) error {
		districtName := c.Args().Get(0)
		if len(districtName) == 0 {
			return cli.NewExitError("district name is required", 1)
		}
		ip := c.Args().Get(1)
		if len(ip) == 0 {
			return cli.NewExitError("ip is required", 1)
		}

		resp, err := api.DefaultClient.Post("/districts/"+districtName+"/sign_public_key", nil)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		var districtResp api.DistrictResponse
		err = json.Unmarshal(resp, &districtResp)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		ssh := SSH{
			IP:          ip,
			BastionIP:   districtResp.District.BastionIP,
			Certificate: districtResp.Certificate,
		}
		if ssh.Run("") != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		return nil
	},
}
