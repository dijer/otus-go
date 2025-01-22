package main

import (
	"context"
	"flag"
	"fmt"

	notificationcfg "github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/config/notificationconfig"
	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/logger"
	scheduler "github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/notifications"
	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/rabbitmq"
	factorystorage "github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/storage/factory"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/calendar_scheduler.toml", "Path to calendar_scheduler file")
}

func main() {
	flag.Parse()

	cfg, err := notificationcfg.New(configFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	logg := logger.New(cfg.Logger.Level)

	storage, err := factorystorage.New(factorystorage.Config{
		Database: factorystorage.DatabaseConf{
			Host:     cfg.Database.Host,
			User:     cfg.Database.User,
			Password: cfg.Database.Password,
			DBName:   cfg.Database.DBName,
			Port:     cfg.Database.Port,
		},
		Storage: factorystorage.StorageConf{
			Storage: cfg.Storage.Storage,
		},
	})
	if err != nil {
		logg.Error(err.Error())
		return
	}

	rabbitClient, err := rabbitmq.New(rabbitmq.Config{
		Port:     cfg.Rabbit.Port,
		Host:     cfg.Rabbit.Host,
		User:     cfg.Rabbit.User,
		Password: cfg.Rabbit.Password,
		Exchange: cfg.Rabbit.Exchange,
		Queue:    cfg.Rabbit.Queue,
	}, logg)
	if err != nil {
		logg.Error(err.Error())
		return
	}
	defer rabbitClient.Close()

	notification := scheduler.NewScheduler(cfg, rabbitClient, storage, logg)
	ctx := context.Background()
	err = notification.Run(ctx)
	if err != nil {
		logg.Error(err.Error())
		return
	}
}
