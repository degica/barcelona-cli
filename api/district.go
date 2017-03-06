package api

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func (cli *Client) ShowDistrict(name string) (*District, error) {
	resp, err := cli.Request("GET", "/districts/"+name, nil)
	if err != nil {
		return nil, err
	}

	var dResp DistrictResponse
	err = json.Unmarshal(resp, &dResp)
	if err != nil {
		return nil, err
	}
	return dResp.District, nil
}

func (cli *Client) ListDistricts() ([]*District, error) {
	resp, err := cli.Request("GET", "/districts", nil)
	if err != nil {
		return nil, err
	}
	var dResp DistrictResponse
	err = json.Unmarshal(resp, &dResp)
	if err != nil {
		return nil, err
	}
	return dResp.Districts, nil
}

func (cli *Client) CreateDistrict(req *DistrictRequest) (*District, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := cli.Request("POST", "/districts", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	var dResp DistrictResponse
	err = json.Unmarshal(resp, &dResp)
	if err != nil {
		return nil, err
	}

	return dResp.District, nil
}

func (cli *Client) UpdateDistrict(req *DistrictRequest) (*District, error) {
	name := req.Name
	req.Name = ""

	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := cli.Request("PATCH", "/districts/"+name, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	var dResp DistrictResponse
	err = json.Unmarshal(resp, &dResp)
	if err != nil {
		return nil, err
	}

	return dResp.District, nil
}

func (cli *Client) ApplyDistrict(name string) error {
	_, err := cli.Request("POST", "/districts/"+name+"/apply_stack", nil)
	if err != nil {
		return err
	}
	return nil
}

func (cli *Client) DeleteDistrict(name string) error {
	_, err := cli.Request("DELETE", fmt.Sprintf("/districts/%s", name), nil)
	if err != nil {
		return err
	}
	return nil
}

func (cli *Client) PutPlugin(districtName string, plugin *Plugin) (*Plugin, error) {
	b, err := json.Marshal(plugin)
	if err != nil {
		return nil, err
	}

	resp, err := cli.Request("PUT", fmt.Sprintf("/districts/%s/plugins/%s", districtName, plugin.Name), bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	var pResp PluginResponse
	err = json.Unmarshal(resp, &pResp)
	if err != nil {
		return nil, err
	}

	return pResp.Plugin, nil
}

func (cli *Client) DeletePlugin(districtName string, pluginName string) error {
	_, err := cli.Request("DELETE", fmt.Sprintf("/districts/%s/plugins/%s", districtName, pluginName), nil)
	if err != nil {
		return err
	}
	return nil
}
