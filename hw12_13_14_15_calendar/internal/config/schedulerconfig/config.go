package schedulerconfig

import (
	"os"

	log "github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/logger"

	toml "github.com/pelletier/go-toml"
)

type SchedulerConfig struct {
	Rabbit    RabbitConf
	Logger    LoggerConf
	Database  DatabaseConf
	Scheduler SchedulerConf
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

func New(configFile string) (*SchedulerConfig, error) {
	var config SchedulerConfig

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
