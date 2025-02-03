package model

import (
	"fmt"
	"net"
	"strings"
)

type Config struct {
	API APIConfig `yaml:"api"`
	DB  DBConfig  `yaml:"database"`
}

type APIConfig struct {
	Address string `yaml:"bind-address"`
	Port    int    `yaml:"bind-port"`
	Prefix  string `yaml:"prefix"`
	Expire  int    `yaml:"session-expire"` // Hour
}

type DBConfig struct {
	Driver string `yaml:"driver"`
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	DBName string `yaml:"dbname"`
	User   string `yaml:"user"`
	Passwd string `yaml:"password"`
}

func (config *Config) Verification() error {
	// APIConfig
	// - Address
	if valid := isValidIP(config.API.Address); valid == "Invalid" {
		return fmt.Errorf("api: invalid address")
	}
	// - Port
	if config.API.Port < 1 || config.API.Port > 65535 {
		return fmt.Errorf("api: invalid port number")
	}
	// - Prefix
	if !strings.HasPrefix(config.API.Prefix, "/") || strings.HasSuffix(config.API.Prefix, "/") {
		return fmt.Errorf("api: invalid prefix")
	}

	// DBConfig
	// - Driver
	switch config.DB.Driver {
	case "mysql", "mariadb":
		if config.DB.Port == 0 {
			config.DB.Port = 3306
		}
	case "postgres", "postgresql":
		return fmt.Errorf("db: postgresql does not support")
	default:
		return fmt.Errorf("db: required valid driver")
	}
	// - Host
	if config.DB.Host == "" {
		return fmt.Errorf("db: required database host")
	}
	// - Port
	if config.DB.Port < 1 && config.DB.Port > 65535 {
		return fmt.Errorf("db: invalid port number")
	}
	// - DBName
	if config.DB.DBName == "" {
		return fmt.Errorf("db: required database")
	}
	// - User
	if config.DB.User == "" {
		return fmt.Errorf("db: required user")
	}
	// - Pass
	if config.DB.Passwd == "" {
		return fmt.Errorf("db: required password")
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
