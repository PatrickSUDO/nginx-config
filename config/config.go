package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	PathAccessControl []string `yaml:"path_access_control"`
	IPWhitelist       []string `yaml:"ip_whitelist"`
	ClientCerts       []string `yaml:"client_certs"`
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
