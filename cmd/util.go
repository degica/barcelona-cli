package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/degica/barcelona-cli/api"
	"github.com/ghodss/yaml"
)

func PrettyJSON(b []byte) string {
	var v interface{}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return ""
	}
	pretty, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return ""
	}
	return string(pretty)
}

func PrintHeritage(h *api.Heritage) {
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

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func ask(s string, required bool, secret bool) string {

	var response string
	var err error
	for {
		fmt.Printf("%s: ", s)

		if secret {
			bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
			if err != nil {
				continue
			}
			fmt.Printf("\n")
			response = string(bytePassword)
		} else {
			reader := bufio.NewReader(os.Stdin)
			response, err = reader.ReadString('\n')
			if err != nil {
				continue
			}
		}
		response = strings.TrimSpace(response)

		if len(response) == 0 && required {
			continue
		}
		break
	}

	return response
}

func areYouSure(message string) bool {
	for {
		res := ask(fmt.Sprintf("%s [y/n]", message), false, false)
		if res == "y" {
			return true
		} else if res == "n" {
			return false
		}
	}
}

type HeritageConfig struct {
	Environments map[string]*api.Heritage `yaml:"environments" json:"environments"`
}

func loadHeritageConfig() (*HeritageConfig, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	configFilePath := pwd + "/barcelona.yml"

	configFile, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	var config HeritageConfig
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func LoadEnvironment(env string) (*api.Heritage, error) {
	config, err := loadHeritageConfig()
	if err != nil {
		return nil, err
	}
	heritage := config.Environments[env]
	if heritage == nil {
		return nil, errors.New("environment is invalid")
	}
	return heritage, nil
}
