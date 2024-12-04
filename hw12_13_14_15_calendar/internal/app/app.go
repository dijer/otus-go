package app

import (
	"context"

	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type App interface {
	AddEvent(ctx context.Context, event storage.Event) error
	UpdateEvent(ctx context.Context, event storage.Event) error
	DeleteEvent(ctx context.Context, id int64) error
	GetEventsList(ctx context.Context) ([]storage.Event, error)
}

type Logger interface {
	Info(msg ...string)
	Error(msg ...string)
	Warn(msg ...string)
	Debug(msg ...string)
}

type Storage interface {
	AddEvent(ctx context.Context, event storage.Event) error
	UpdateEvent(ctx context.Context, event storage.Event) error
	DeleteEvent(ctx context.Context, id int64) error
	GetEventsList(ctx context.Context) ([]storage.Event, error)
}

type app struct {
	logger  Logger
	storage Storage
}

func New(logger Logger, storage Storage) App {
	return &app{
		logger:  logger,
		storage: storage,
	}
}

func (a *app) AddEvent(ctx context.Context, event storage.Event) error {
	return a.storage.AddEvent(ctx, event)
}

func (a *app) UpdateEvent(ctx context.Context, event storage.Event) error {
	return a.storage.UpdateEvent(ctx, event)
}

func (a *app) DeleteEvent(ctx context.Context, id int64) error {
	return a.storage.DeleteEvent(ctx, id)
}

func (a *app) GetEventsList(ctx context.Context) ([]storage.Event, error) {
	return a.storage.GetEventsList(ctx)
}
