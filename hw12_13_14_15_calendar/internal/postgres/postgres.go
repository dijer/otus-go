package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose"

	// init pg driver.
	_ "github.com/lib/pq"
)

type PgClient interface {
	Connect(ctx context.Context) (*sqlx.DB, error)
}

type pgClient struct {
	dsn string
	db  *sqlx.DB
}

func New(dsn string) PgClient {
	return &pgClient{
		dsn: dsn,
	}
}

func (p *pgClient) Connect(ctx context.Context) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", p.dsn)
	if err != nil {
		return nil, err
	}

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	p.db = db

	err = p.migrate()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (p *pgClient) migrate() error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return err
	}

	err = goose.Up(p.db.DB, "migrations")
	return err
}

func (p *pgClient) Close(_ context.Context) error {
	return p.db.Close()
}
