package api

import (
	"bytes"
	"encoding/json"
)

func (cli *Client) CreateHeritage(districtName string, h *Heritage) (*Heritage, error) {
	j, err := json.Marshal(h)
	if err != nil {
		return nil, err
	}

	resp, err := cli.Post("/districts/"+districtName+"/heritages", bytes.NewBuffer(j))
	if err != nil {
		return nil, err
	}
	var hResp HeritageResponse
	err = json.Unmarshal(resp, &hResp)
	if err != nil {
		return nil, err
	}

	return hResp.Heritage, nil
}
