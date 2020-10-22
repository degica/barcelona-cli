package api

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func (h *Heritage) Print() {
	fmt.Printf("Name:          %s\n", h.Name)
	fmt.Printf("Image Name:    %s\n", h.ImageName)
	fmt.Printf("Image Tag :    %s\n", h.ImageTag)
	fmt.Printf("Version:       %d\n", h.Version)
	if h.BeforeDeploy != nil {
		fmt.Printf("Before Deploy: %s\n", *h.BeforeDeploy)
	} else {
		fmt.Printf("Before Deploy: None\n")
	}
	fmt.Printf("Token:         %s\n", h.Token)
	fmt.Printf("Scheduled Tasks:\n")
	for _, task := range h.ScheduledTasks {
		fmt.Printf("%-20s %s\n", task.Schedule, task.Command)
	}

	fmt.Printf("Environment Variables\n")
	for name, value := range h.EnvVars {
		fmt.Printf("  %s: %s\n", name, value)
	}
}
