package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	IPFilter map[string][]string `yaml:"ipfilter"`
	Catchall map[string]struct {
		Port int `yaml:"port"`
	} `yaml:"catchall"`
	App map[string]struct {
		Catchall                   string   `yaml:"catchall"`
		FQDN                       []string `yaml:"fqdn"`
		RuntimePort                int      `yaml:"runtime_port"`
		PathBasedAccessRestriction map[string]struct {
			IPFilter string `yaml:"ipfilter"`
		} `yaml:"path_based_access_restriction"`
	} `yaml:"app"`
}

func LoadConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func LoadConfigFromString(yamlStr string) (*Config, error) {
	var config Config
	err := yaml.Unmarshal([]byte(yamlStr), &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
