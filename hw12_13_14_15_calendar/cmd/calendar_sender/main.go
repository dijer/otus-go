package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	config "github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/config/senderconfig"
	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/notifications"
	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/rabbitmq"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/calendar_sender.toml", "Path to calendar_sender file")
}

func main() {
	flag.Parse()

	cfg, err := config.New(configFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	logg := logger.New(cfg.Logger.Level)

	rabbitURL := fmt.Sprintf("amqp://%s:%s@%s:%d", cfg.Rabbit.User, cfg.Rabbit.Password, cfg.Rabbit.Host, cfg.Rabbit.Port)
	rabbitClient, err := rabbitmq.New(rabbitURL, logg)
	if err != nil {
		logg.Error(err.Error())
		return
	}

	sender := notifications.NewSender(cfg, rabbitClient, logg)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err = sender.Run()
		if err != nil {
			logg.Error(err.Error())
		}
	}()

	<-sigs
}
