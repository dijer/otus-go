package config

import (
	"os"

	log "github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/logger"
	toml "github.com/pelletier/go-toml"
)

type Config struct {
	Logger  LoggerConf
	Storage struct {
		Storage string
	}
	Database DatabaseConf
	Server   ServerConf
}

type LoggerConf struct {
	Level log.LogLevel
}

type DatabaseConf struct {
	Host, User, Password, DBName, Migrate string
	Port                                  int64
}

type ServerConf struct {
	Host string
	Port int64
}

func NewConfig(configFile string) (*Config, error) {
	var config Config

	file, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	err = toml.NewDecoder(file).Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
