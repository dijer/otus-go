package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/app"
	config "github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/config/calendar"
	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/logger"
	grpcserver "github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/server/grpc"
	httpserver "github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/server/http"
	factorystorage "github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/storage/factory"
)

var configFile string

const serverShutdownTimeout = 5 * time.Second

func init() {
	flag.StringVar(&configFile, "config", "./configs/config.toml", "Path to configuration file")
}

var wg sync.WaitGroup

func main() {
	flag.Parse()

	cfg, err := config.New(configFile)
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

	calendar := app.New(logg, storage)

	httpServer := httpserver.New(logg, calendar, cfg.HTTP, cfg.GRPC)
	grpcServer := grpcserver.New(logg, calendar, cfg.GRPC)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	logg.Info("calendar is running...")

	wg.Add(2)

	go func() {
		defer wg.Done()
		fmt.Printf("http runnin on: %v\n", cfg.HTTP.Port)
		if err := httpServer.Start(ctx); err != nil {
			logg.Error("failed to start http server: " + err.Error())
			cancel()
			os.Exit(1)
		}
	}()

	go func() {
		defer wg.Done()
		fmt.Printf("grpc runnin on: %v\n", cfg.GRPC.Port)
		if err := grpcServer.Start(ctx); err != nil {
			logg.Error("failed to start grpc server: " + err.Error())
			cancel()
			os.Exit(1)
		}
	}()

	go func() {
		<-ctx.Done()

		logg.Info("shut down servers")

		stopCtx, stopCancel := context.WithTimeout(context.Background(), serverShutdownTimeout)
		defer stopCancel()
		if err := httpServer.Stop(stopCtx); err != nil {
			logg.Error("Failed to stop HTTP server: " + err.Error())
		}

		if err := grpcServer.Stop(stopCtx); err != nil {
			logg.Error("Failed to stop gRPC server: " + err.Error())
		}

		logg.Info("Servers stopped.")
		cancel()
	}()

	wg.Wait()
}
