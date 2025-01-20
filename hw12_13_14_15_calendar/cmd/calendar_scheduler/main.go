package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	config "github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/config/schedulerconfig"
	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/logger"
	scheduler "github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/notifications"
	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/postgres"
	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/rabbitmq"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/calendar_scheduler.toml", "Path to calendar_scheduler file")
}

func main() {
	flag.Parse()

	cfg, err := config.New(configFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	logg := logger.New(cfg.Logger.Level)

	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%d",
		cfg.Database.User, cfg.Database.Password, cfg.Database.DBName, cfg.Database.Host, cfg.Database.Port,
	)
	pgClient := postgres.New(dsn)
	db, err := pgClient.Connect(context.Background())
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	rabbitURL := fmt.Sprintf("amqp://%s:%s@%s:%d", cfg.Rabbit.User, cfg.Rabbit.Password, cfg.Rabbit.Host, cfg.Rabbit.Port)
	rabbitClient, err := rabbitmq.New(rabbitURL, logg)
	if err != nil {
		logg.Error(err.Error())
		return
	}
	defer rabbitClient.Close()

	notification := scheduler.NewScheduler(cfg, rabbitClient, db, logg)
	err = notification.Run()
	if err != nil {
		logg.Error(err.Error())
		return
	}
}
