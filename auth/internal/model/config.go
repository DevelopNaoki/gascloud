package model

type Config struct {
	API APIConfig `yaml:"api"`
	DB  DBConfig  `yaml:"database"`
}

type APIConfig struct {
	Address string `yaml:"bind-address"`
	Port    int    `yaml:"bind-port"`
	Prefix  string `yaml:"prefix"`
}

type DBConfig struct {
	Driver string `yaml:"driver"`
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	DBName string `yaml:"dbname"`
	User   string `yaml:"user"`
	Pass   string `yaml:"pass"`
}
