package factorystorage

import (
	"context"
	"fmt"

	config "github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/config/calendar"
	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/postgres"
	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/storage/sql"
)

func New(cfg *config.Config) (storage.Storage, error) {
	var storage storage.Storage
	if cfg.Storage.Storage == "sql" {
		dsn := fmt.Sprintf(
			"user=%s password=%s dbname=%s host=%s port=%d",
			cfg.Database.User, cfg.Database.Password, cfg.Database.DBName, cfg.Database.Host, cfg.Database.Port,
		)
		pgClient := postgres.New(dsn)
		db, err := pgClient.Connect(context.Background())
		if err != nil {
			return nil, err
		}

		storage = sqlstorage.New(cfg.Database, db)
	} else {
		storage = memorystorage.New()
	}

	return storage, nil
}
