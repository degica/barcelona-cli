package operations

import (
	"encoding/json"
	"github.com/degica/barcelona-cli/api"
	"github.com/degica/barcelona-cli/utils"
	"io"
)

type SshcmdOperationApiClient interface {
	Post(path string, body io.Reader) ([]byte, error)
}

type SshcmdOperation struct {
	client        SshcmdOperationApiClient
	districtName  string
	ip            string
	config        utils.SshConfig
	commandRunner utils.SshCommandRunner
}

func NewSshcmdOperation(
	client SshcmdOperationApiClient,
	districtName string,
	ip string,
	config utils.SshConfig,
	commandRunner utils.SshCommandRunner) *SshcmdOperation {
	return &SshcmdOperation{
		client:        client,
		districtName:  districtName,
		ip:            ip,
		config:        config,
		commandRunner: commandRunner,
	}
}

func (oper SshcmdOperation) run() *runResult {
	if len(oper.districtName) == 0 {
		return error_result("district name is required")
	}
	if len(oper.ip) == 0 {
		return error_result("ip is required")
	}

	resp, err := oper.client.Post("/districts/"+oper.districtName+"/sign_public_key", nil)
	if err != nil {
		return error_result(err.Error())
	}

	var districtResp api.DistrictResponse
	err = json.Unmarshal(resp, &districtResp)
	if err != nil {
		return error_result(err.Error())
	}

	ssh := utils.NewSshCommand(
		oper.ip,
		districtResp.District.BastionIP,
		districtResp.Certificate,
		oper.config,
		oper.commandRunner,
	)

	if ssh.Run("") != nil {
		return error_result(err.Error())
	}
	return ok_result()
}
