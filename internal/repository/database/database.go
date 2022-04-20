package database

import (
	"context"
	"github.com/IakimenkoD/xm-companies-service/internal/config"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Client struct {
	*sqlx.DB

	SchemaName string
}

func NewClient(cfg *config.Config) (*Client, error) {
	db, err := sqlx.Open("pgx", cfg.DB.URL)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	db.SetMaxIdleConns(cfg.DB.MaxIdleConns)

	return &Client{
		db,
		cfg.DB.SchemaName,
	}, nil
}

func (db *Client) Migrate() error {
	if _, err := db.Exec(`CREATE SCHEMA IF NOT EXISTS ` + db.SchemaName); err != nil {
		return errors.Wrap(err, "can't create schema")
	}
	m, err := migrations(db.SchemaName, "migrations")
	if err != nil {
		return errors.Wrap(err, "can't create a new migrator instance")
	}

	if err = m.Migrate(db.DB.DB); err != nil {
		return errors.Wrap(err, "can't migrate the db")
	}

	return nil
}

func (db *Client) StatusCheck(ctx context.Context) error {
	var ok bool
	return db.QueryRowContext(ctx, `SELECT true`).Scan(&ok)
}
