package api

import (
	"encoding/json"
	"net/http"
)

func (cli *Client) Login(endpoint string, githubToken string) (*User, error) {
	url := endpoint + pathPrefix + "/login"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-GitHub-Token", githubToken)
	b, err := cli.rawRequest(req)
	if err != nil {
		return nil, err
	}

	var userResp UserResponse
	err = json.Unmarshal(b, &userResp)
	if err != nil {
		return nil, err
	}
	return userResp.User, nil
}
