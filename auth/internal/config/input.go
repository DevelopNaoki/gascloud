package config

import (
	"github.com/DevelopNaoki/gascloud/auth/internal/model"
	"github.com/goccy/go-yaml"
	"os"
)

func LoadConfigFile(path string) (config model.Config, err error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}

	config = model.Config{
		API: model.APIConfig{
			Address: "0.0.0.0",
			Port:    80,
			Prefix:  "/",
			Expire:  2,
		},
	}

	err = yaml.Unmarshal(f, &config)
	if err != nil {
		return config, err
	}

	err = config.Verification()
	if err != nil {
		return config, err
	}

	return config, nil
}
