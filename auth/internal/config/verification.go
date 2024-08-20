package config

import (
	"fmt"
	"net"
	"strings"

	"github.com/DevelopNaoki/gascloud/auth/internal/model"
)

func VerificationConfigs(config model.Config) (err error) {
	// APIConfig
	// - Address
	if config.API.Address == "" {
		config.API.Address = "0.0.0.0"
	} else if valid := isValidIP(config.API.Address); valid == "Invalid" {
		err = fmt.Errorf("api: invalid address")
	}
	// - Port
	if config.API.Port == 0 {
		config.API.Port = 80
	} else if config.API.Port < 1 || config.API.Port > 65535 {
		err = fmt.Errorf("%s, api: invalid port number", err.Error())
	}
	// - Prefix
	if config.API.Prefix == "" {
		config.API.Prefix = "/"
	} else if !strings.HasPrefix(config.API.Prefix, "/") || strings.HasSuffix(config.API.Prefix, "/") {
		err = fmt.Errorf("%s, api: invalid prefix", err.Error())
	}

	// DBConfig
	// - Driver
	switch config.DB.Driver {
	case "mysql", "mariadb":
		if config.DB.Port == 0 {
			config.DB.Port = 3306
		}
	case "postgresql":
		err = fmt.Errorf("%s, db: postgresql does not support", err.Error())
	default:
		err = fmt.Errorf("%s, db: required valid driver", err.Error())
	}
	// - Host
	if config.DB.Host == "" {
		err = fmt.Errorf("%s, db: required database host", err.Error())
	}
	// - Port
	if config.DB.Port < 1 && config.DB.Port > 65535 {
		err = fmt.Errorf("%s, db: invalid port number", err.Error())
	}
	// - DBName
	if config.DB.DBName == "" {
		err = fmt.Errorf("%s, db: required database", err.Error())
	}
	// - User
	if config.DB.User == "" {
		err = fmt.Errorf("%s, db: required user", err.Error())
	}
	// - Pass
	if config.DB.Pass == "" {
		err = fmt.Errorf("%s, db: required password", err.Error())
	}

	if err != nil {
		return err
	}

	return nil
}

func isValidIP(address string) string {
	ip := net.ParseIP(address)
	if ip == nil {
		return "Invalid"
	}

	if ip.To4() != nil {
		return "IPv4"
	}

	return "IPv6"
}
