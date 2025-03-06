package model

import (
	"fmt"
	"net"
)

type Config struct {
	API   APIConfig   `yaml:"api"`
	Cache CacheConfig `yaml:"cache"`
	DB    DBConfig    `yaml:"database"`
}

type APIConfig struct {
	Address string `yaml:"bind-address"`
	Port    int    `yaml:"bind-port"`
	Prefix  string `yaml:"prefix"`
	Expire  int    `yaml:"session-expire"` // Hour
}

type CacheConfig struct {
	Driver string `yaml:"driver"`
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
}

type DBConfig struct {
	Driver string `yaml:"driver"`
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	DBName string `yaml:"dbname"`
	User   string `yaml:"user"`
	Passwd string `yaml:"password"`
}

func (config *Config) Verification() (err error) {
	err = config.API.Verification()
	if err != nil {
		return fmt.Errorf("api: %s", err.Error())
	}
	err = config.Cache.Verification()
	if err != nil {
		return fmt.Errorf("cache: %s", err.Error())
	}
	err = config.DB.Verification()
	if err != nil {
		return fmt.Errorf("db: %s", err.Error())
	}
	return nil
}

func (api *APIConfig) Verification() error {
	if valid := isValidIP(api.Address); valid == "Invalid" {
		return fmt.Errorf("invalid address")
	}
	if api.Port < 1 || api.Port > 65535 {
		return fmt.Errorf("invalid port number")
	}
	return nil
}

func (cache *CacheConfig) Verification() error {
	switch cache.Driver {
	case "memcache", "memcached":
		if cache.Port == 0 {
			cache.Port = 11211
		}
	case "redis":
		return fmt.Errorf("redis does not support")
	default:
		return fmt.Errorf("required valid driver")
	}
	if cache.Host == "" {
		return fmt.Errorf("required cache host")
	}
	if cache.Port < 1 && cache.Port > 65535 {
		return fmt.Errorf("invalid port number")
	}
	return nil
}

func (db *DBConfig) Verification() error {
	switch db.Driver {
	case "mysql", "mariadb":
		if db.Port == 0 {
			db.Port = 3306
		}
	case "postgres", "postgresql":
		return fmt.Errorf("postgresql does not support")
	default:
		return fmt.Errorf("required valid driver")
	}
	if db.Host == "" {
		return fmt.Errorf("required database host")
	}
	if db.Port < 1 && db.Port > 65535 {
		return fmt.Errorf("invalid port number")
	}
	if db.DBName == "" {
		return fmt.Errorf("required database")
	}
	if db.User == "" {
		return fmt.Errorf("required user")
	}
	if db.Passwd == "" {
		return fmt.Errorf("required password")
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
