package storage

import (
	"context"
	"errors"
	"time"
)

type Storage interface {
	AddEvent(ctx context.Context, event Event) error
	UpdateEvent(ctx context.Context, event Event) error
	DeleteEvent(ctx context.Context, id int64) error
	GetEventsList(ctx context.Context) ([]Event, error)
}

type Event struct {
	ID          int64
	Title       string
	StartTime   time.Time
	EndTime     time.Time
	Description *string
}

var (
	ErrEventNotFound = errors.New("event not found")
	ErrDateBusy      = errors.New("date is busy by another event")
)
