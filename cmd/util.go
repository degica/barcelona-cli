package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"

	"github.com/degica/barcelona-cli/api"
)

var HeritageConfigFilePath string

func PrintOneoff(o *api.Oneoff) {
	fmt.Printf("Task ARN: %s\n", o.TaskARN)
}

type HeritageConfig struct {
	Environments map[string]*api.Heritage `yaml:"environments" json:"environments"`
	Review       *api.ReviewAppDefinition `yaml:"review" json:"review"`
}

func loadHeritageConfig() (*HeritageConfig, error) {
	configFile, err := ioutil.ReadFile(HeritageConfigFilePath)
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

func LoadReviewDefinition() (*api.ReviewAppDefinition, error) {
	config, err := loadHeritageConfig()
	if err != nil {
		return nil, err
	}
	review := config.Review
	if review == nil {
		return nil, errors.New("reviewapp is invalid")
	}
	return review, nil
}
