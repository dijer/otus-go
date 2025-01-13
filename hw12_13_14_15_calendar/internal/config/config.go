package config

import (
	"os"

	log "github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/logger"
	toml "github.com/pelletier/go-toml"
)

type Config struct {
	Logger   LoggerConf
	Storage  StorageConf
	Database DatabaseConf
	HTTP     HTTPServerConf
	GRPC     GRPCServerConf
}

type LoggerConf struct {
	Level log.LogLevel
}

type DatabaseConf struct {
	Host,
	User,
	Password,
	DBName,
	Migrate string
	Port int
}

type HTTPServerConf struct {
	Host string
	Port int
}

type GRPCServerConf struct {
	Host      string
	Port      int
	Transport string
}

type StorageConf struct {
	Storage string
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
