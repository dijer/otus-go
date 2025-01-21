package memorystorage

import (
	"context"
	"sync"

	storage "github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type Event = storage.Event

type Storage struct {
	storage.Storage
	events map[int32]Event
	mu     sync.RWMutex
	id     int32
}

func New() *Storage {
	return &Storage{
		events: make(map[int32]Event),
		id:     0,
	}
}

func (s *Storage) AddEvent(_ context.Context, event Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, e := range s.events {
		if event.StartTime.Before(e.EndTime) && event.EndTime.After(e.StartTime) {
			return storage.ErrDateBusy
		}
	}

	s.id++
	event.ID = s.id
	s.events[event.ID] = event

	return nil
}

func (s *Storage) UpdateEvent(_ context.Context, event Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.events[event.ID]; !exists {
		return storage.ErrEventNotFound
	}

	for _, e := range s.events {
		if event.StartTime.Before(e.EndTime) && event.EndTime.After(e.StartTime) && event.ID != e.ID {
			return storage.ErrDateBusy
		}
	}

	s.events[event.ID] = event

	return nil
}

func (s *Storage) DeleteEvent(_ context.Context, id int32) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.events[id]; !exists {
		return storage.ErrEventNotFound
	}

	delete(s.events, id)

	return nil
}

func (s *Storage) GetEventsList(_ context.Context) ([]Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	events := make([]Event, 0, len(s.events))

	for _, event := range s.events {
		events = append(events, event)
	}

	return events, nil
}
