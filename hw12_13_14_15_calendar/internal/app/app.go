package app

import (
	"context"
	"fmt"

	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/storage/sql"
)

type App interface {
	AddEvent(ctx context.Context, event storage.Event) error
	UpdateEvent(ctx context.Context, event storage.Event) error
	DeleteEvent(ctx context.Context, id int32) error
	GetEventsList(ctx context.Context) ([]storage.Event, error)
}

type app struct {
	logger  logger.Logger
	storage storage.Storage
}

func New(logger *logger.Logger, cfg *config.Config) App {
	var storage storage.Storage
	if cfg.Storage.Storage == "sql" {
		storage = sqlstorage.New(cfg.Database)
		err := storage.(*sqlstorage.Storage).Connect(context.Background())
		if err != nil {
			panic(err)
		}
	} else {
		storage = memorystorage.New()
	}

	return &app{
		logger:  *logger,
		storage: storage,
	}
}

func (a app) AddEvent(ctx context.Context, event storage.Event) error {
	fmt.Println("app add event")

	if a.storage == nil {
		panic("app storage empty!")
	}

	return a.storage.AddEvent(ctx, event)
}

func (a app) UpdateEvent(ctx context.Context, event storage.Event) error {
	return a.storage.UpdateEvent(ctx, event)
}

func (a app) DeleteEvent(ctx context.Context, id int32) error {
	return a.storage.DeleteEvent(ctx, id)
}

func (a app) GetEventsList(ctx context.Context) ([]storage.Event, error) {
	return a.storage.GetEventsList(ctx)
}
