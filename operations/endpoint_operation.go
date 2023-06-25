package operations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/degica/barcelona-cli/api"
	"github.com/degica/barcelona-cli/utils"
	"github.com/olekukonko/tablewriter"
)

type EndpointOperationApiClient interface {
	Request(method string, path string, body io.Reader) ([]byte, error)
	Read(secret bool) (string, error)
}

type EndpointOperation struct {
	districtName   string
	endpointName   string
	public         bool
	cert_arn       string
	policy         string
	noConfirmation bool
	op_type        OperationType
	client         EndpointOperationApiClient
}

func NewEndpointOperation(districtName string, endpointName string, public bool, cert_arn string, policy string, noConfirmation bool, op OperationType, client EndpointOperationApiClient) *EndpointOperation {
	return &EndpointOperation{
		districtName:   districtName,
		endpointName:   endpointName,
		public:         public,
		cert_arn:       cert_arn,
		policy:         policy,
		noConfirmation: noConfirmation,
		op_type:        op,
		client:         client,
	}
}

func (oper EndpointOperation) run() *runResult {
	operations := map[OperationType](func(oper EndpointOperation) *runResult){
		Create: endpoint_create,
		Update: endpoint_update,
		Delete: endpoint_delete,
		Show:   endpoint_show,
		List:   endpoint_list,
	}

	if function, ok := operations[oper.op_type]; ok {
		return function(oper)
	}

	return error_result("unknown operation")
}

func endpoint_create(oper EndpointOperation) *runResult {
	request := api.Endpoint{
		Name:          oper.endpointName,
		Public:        &oper.public,
		CertificateID: oper.cert_arn,
		SslPolicy:     oper.policy,
	}

	b, err := json.Marshal(&request)
	if err != nil {
		return error_result(err.Error())
	}

	resp, err := oper.client.Request("POST", fmt.Sprintf("/districts/%s/endpoints", oper.districtName), bytes.NewBuffer(b))
	if err != nil {
		return error_result(err.Error())
	}
	var eResp api.EndpointResponse
	err = json.Unmarshal(resp, &eResp)
	if err != nil {
		return error_result(err.Error())
	}
	printEndpoint(eResp.Endpoint)

	return ok_result()
}

func endpoint_show(oper EndpointOperation) *runResult {
	resp, err := oper.client.Request("GET", fmt.Sprintf("/districts/%s/endpoints/%s", oper.districtName, oper.endpointName), nil)
	if err != nil {
		return error_result(err.Error())
	}
	var eResp api.EndpointResponse
	err = json.Unmarshal(resp, &eResp)
	if err != nil {
		return error_result(err.Error())
	}
	printEndpoint(eResp.Endpoint)

	return ok_result()
}

func endpoint_list(oper EndpointOperation) *runResult {
	resp, err := oper.client.Request("GET", fmt.Sprintf("/districts/%s/endpoints", oper.districtName), nil)
	if err != nil {
		return error_result(err.Error())
	}
	var eResp api.EndpointResponse
	err = json.Unmarshal(resp, &eResp)
	if err != nil {
		return error_result(err.Error())
	}
	printEndpoints(eResp.Endpoints)

	return ok_result()
}

func endpoint_update(oper EndpointOperation) *runResult {
	request := api.Endpoint{
		CertificateID: oper.cert_arn,
		SslPolicy:     oper.policy,
	}

	b, err := json.Marshal(&request)
	if err != nil {
		return error_result(err.Error())
	}

	resp, err := oper.client.Request("PATCH", fmt.Sprintf("/districts/%s/endpoints/%s", oper.districtName, oper.endpointName), bytes.NewBuffer(b))
	if err != nil {
		return error_result(err.Error())
	}
	var eResp api.EndpointResponse
	err = json.Unmarshal(resp, &eResp)
	if err != nil {
		return error_result(err.Error())
	}
	printEndpoint(eResp.Endpoint)

	return ok_result()
}

func endpoint_delete(oper EndpointOperation) *runResult {
	fmt.Printf("You are attempting to delete /%s/endpoints/%s\n", oper.districtName, oper.endpointName)
	if !oper.noConfirmation && !utils.AreYouSure("This operation cannot be undone. Are you sure?", oper.client) {
		return nil
	}

	_, err := oper.client.Request("DELETE", fmt.Sprintf("/districts/%s/endpoints/%s", oper.districtName, oper.endpointName), nil)
	if err != nil {
		return error_result(err.Error())
	}
	return ok_result()
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
