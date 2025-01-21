package storage

import (
	"context"
	"errors"
	"time"
)

type Storage interface {
	AddEvent(ctx context.Context, event Event) error
	UpdateEvent(ctx context.Context, event Event) error
	DeleteEvent(ctx context.Context, id int32) error
	GetEventsList(ctx context.Context) ([]Event, error)
	GetNotifications(ctx context.Context) ([]Event, error)
	CleanupOldEvents(ctx context.Context) error
	SendNotification(ctx context.Context, id int32) error
}

type Event struct {
	ID          int32     `db:"id"`
	Title       string    `db:"title"`
	Owner       int32     `db:"owner"`
	StartTime   time.Time `db:"start_time"`
	EndTime     time.Time `db:"end_time"`
	Description *string   `db:"description"`
}

var (
	ErrEventNotFound = errors.New("event not found")
	ErrDateBusy      = errors.New("date is busy by another event")
)
