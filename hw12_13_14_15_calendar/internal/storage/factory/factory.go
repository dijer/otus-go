package factorystorage

import (
	"context"

	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/storage/sql"
)

func New(cfg *config.Config) (storage.Storage, error) {
	var storage storage.Storage
	if cfg.Storage.Storage == "sql" {
		storage = sqlstorage.New(cfg.Database)
		err := storage.(*sqlstorage.Storage).Connect(context.Background())
		if err != nil {
			return nil, err
		}
	} else {
		storage = memorystorage.New()
	}

	return storage, nil
}
