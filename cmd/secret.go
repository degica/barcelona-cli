package cmd

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/degica/barcelona-cli/api"
	"github.com/urfave/cli"
)

type secretRequest struct {
	Type      string `json:"type"`
	Plaintext string `json:"plaintext"`
}

type secretResponse struct {
	Type           string `json:"type"`
	EncryptedValue string `json:"encrypted_value"`
}

var SecretCommand = cli.Command{
	Name:  "secret",
	Usage: "Manage secrets",
	Subcommands: []cli.Command{
		{
			Name:      "transit",
			Usage:     "Creates a transit secret",
			ArgsUsage: "DISTRICT_NAME",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "file, f",
				},
			},
			Action: func(c *cli.Context) error {
				districtName := c.Args().Get(0)
				if len(districtName) == 0 {
					return cli.NewExitError("district name is required", 1)
				}

				fileName := c.String("file")
				var buffer bytes.Buffer
				if len(fileName) > 0 {
					bs, err := ioutil.ReadFile(fileName)
					if err != nil {
						return err
					}
					buffer.WriteString(string(bs))
				} else {
					scanner := bufio.NewScanner(os.Stdin)

					for scanner.Scan() {
						buffer.WriteString(scanner.Text())
					}
					if err := scanner.Err(); err != nil {
						log.Println(err)
					}
				}
				sEnc := base64.StdEncoding.EncodeToString([]byte(buffer.Bytes()))

				req := secretRequest{
					Type:      "transit",
					Plaintext: sEnc,
				}

				j, err := json.Marshal(&req)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}

				resp, err := api.DefaultClient.Post("/districts/"+districtName+"/secrets", bytes.NewBuffer(j))
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}

				var rResp secretResponse
				err = json.Unmarshal(resp, &rResp)
				fmt.Printf("%s", rResp.EncryptedValue)

				return nil
			},
		},
	},
}
