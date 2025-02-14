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
	App map[string]AppConfig `yaml:"app"`
}

type AppConfig struct {
	Catchall                   string   `yaml:"catchall"`
	FQDN                       []string `yaml:"fqdn"`
	RuntimePort                int      `yaml:"runtime_port"`
	PathBasedAccessRestriction map[string]struct {
		IPFilter string `yaml:"ipfilter"`
	} `yaml:"path_based_access_restriction"`
}

func LoadConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return LoadConfigFromBytes(data)
}

func LoadConfigFromBytes(data []byte) (*Config, error) {
	var config Config
	err := yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
