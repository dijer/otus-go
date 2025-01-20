package senderconfig

import (
	"os"

	log "github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/logger"

	toml "github.com/pelletier/go-toml"
)

type SenderConfig struct {
	Rabbit RabbitConf
	Logger LoggerConf
}

type RabbitConf struct {
	Port int
	Host,
	User,
	Password,
	Queue,
	Exchange string
}

type LoggerConf struct {
	Level log.LogLevel
}

func New(configFile string) (*SenderConfig, error) {
	var config SenderConfig

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
