package config

import (
	"os"

	"github.com/goccy/go-yaml"

	"github.com/DevelopNaoki/gascloud/auth/internal/model"
)

func ReadConfigFile(path string) (config model.Config, _ error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(f, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}
