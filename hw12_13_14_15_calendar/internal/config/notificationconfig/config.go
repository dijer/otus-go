package notificationcfg

import (
	"os"

	log "github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/logger"

	toml "github.com/pelletier/go-toml"
)

type NotificationConfig struct {
	Rabbit    RabbitConf
	Logger    LoggerConf
	Database  DatabaseConf
	Scheduler SchedulerConf
	Storage   StorageConf
}

type RabbitConf struct {
	Port int
	Host,
	User,
	Password,
	Exchange,
	Queue string
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

type SchedulerConf struct {
	Interval,
	Cleanup string
}

type StorageConf struct {
	Storage string
}

func New(configFile string) (*NotificationConfig, error) {
	var config NotificationConfig

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
