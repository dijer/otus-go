package sqlstorage

import (
	"context"

	config "github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/config/calendar"
	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/storage"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	storage.Storage
	config config.DatabaseConf
	db     *sqlx.DB
}

type Event = storage.Event

func New(config config.DatabaseConf, db *sqlx.DB) *Storage {
	return &Storage{
		config: config,
		db:     db,
	}
}

func (s *Storage) AddEvent(_ context.Context, event Event) error {
	var exists bool
	sql := `select exists (select 1 from events where start_time < $1 and end_time > $2)`
	err := s.db.QueryRow(sql, event.EndTime, event.StartTime).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return storage.ErrDateBusy
	}

	_, err = s.db.Exec(
		`insert into events (title, start_time, end_time, owner, description) values ($1, $2, $3, $4, $5)`,
		event.Title, event.StartTime, event.EndTime, event.Owner, event.Description,
	)
	return err
}

func (s *Storage) UpdateEvent(_ context.Context, event Event) error {
	res, err := s.db.Exec(
		`update events set title=$1, start_time=$2, end_time=$3, description=$4, owner=$5 where id=$6`,
		event.Title, event.StartTime, event.EndTime, event.Description, event.Owner, event.ID,
	)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return storage.ErrEventNotFound
	}

	return nil
}

func (s *Storage) DeleteEvent(_ context.Context, id int32) error {
	res, err := s.db.Exec(
		`delete from events where id=$1`,
		id,
	)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return storage.ErrEventNotFound
	}

	return err
}

func (s *Storage) GetEventsList(_ context.Context) ([]Event, error) {
	var events []Event
	err := s.db.Select(&events, `select id, title, owner, start_time, end_time, description from events`)
	return events, err
}
