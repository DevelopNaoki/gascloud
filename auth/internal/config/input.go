package config

import (
	"os"

	"github.com/goccy/go-yaml"

	"github.com/DevelopNaoki/gascloud/auth/internal/model"
)

func LoadConfigFile(path string) (config model.Config, _ error) {
	// Open
	f, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}

	// Parse
	err = yaml.Unmarshal(f, &config)
	if err != nil {
		return config, err
	}

	// Verfication
	err = VerificationConfigs(config)
	if err != nil {
		return config, err
	}
	
	return config, nil
}
