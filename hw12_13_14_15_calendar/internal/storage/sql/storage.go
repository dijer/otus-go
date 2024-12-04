package sqlstorage

import (
	"context"
	"fmt"

	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/storage"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose"
)

type Storage struct {
	config config.DatabaseConf
	db     *sqlx.DB
}

type Event = storage.Event

func New(config config.DatabaseConf) *Storage {
	return &Storage{
		config: config,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%d",
		s.config.User, s.config.Password, s.config.DBName, s.config.Host, s.config.Port,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return err
	}

	err = s.Migrate()
	if err != nil {
		return err
	}

	err = db.PingContext(ctx)
	if err != nil {
		return err
	}

	s.db = db

	return nil
}

func (s *Storage) Close(_ context.Context) error {
	return s.db.Close()
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
		`insert into events (title, start_time, end_time) values ($1, $2, $3)`,
		event.Title, event.StartTime, event.EndTime,
	)
	return err
}

func (s *Storage) UpdateEvent(_ context.Context, event Event) error {
	res, err := s.db.Exec(
		`update events set title=$1, start_time=$2, end_time=$3 where id=$4`,
		event.Title, event.StartTime, event.EndTime, event.ID,
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

func (s *Storage) DeleteEvent(_ context.Context, id int64) error {
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
	err := s.db.Select(&events, `select * from events`)
	return events, err
}

func (s *Storage) Migrate() error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return err
	}

	err = goose.Up(s.db.DB, s.config.Migrate)
	return err
}
