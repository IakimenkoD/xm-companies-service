package database

import (
	"database/sql"
	"github.com/lopezator/migrator"
	"github.com/pkg/errors"
)

func migrationInit() *migrator.Migration {
	return &migrator.Migration{
		Name: "init",
		Func: func(tx *sql.Tx) error {
			qs := []string{
				`CREATE TABLE IF NOT EXISTS xm.migrations (` +
					`id BIGSERIAL PRIMARY KEY` +
					`, version VARCHAR NOT NULL UNIQUE)`,

				`CREATE TABLE IF NOT EXISTS xm.companies (` +
					`id BIGSERIAL PRIMARY KEY` +
					`, name VARCHAR NOT NULL UNIQUE` +
					`, code VARCHAR NOT NULL UNIQUE` +
					`, country VARCHAR NOT NULL` +
					`, website VARCHAR NOT NULL` +
					`, phone VARCHAR(50)` +
					`, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()` +
					`, updated_at TIMESTAMPTZ` +
					`)`,
			}
			for k, query := range qs {
				if _, err := tx.Exec(query); err != nil {
					return errors.Wrapf(err, "applying 202105051514_init migration #%d", k)
				}
			}
			return nil
		},
	}
}

/* ROLLBACK SQL
DROP TABLE IF EXISTS xm.companies;
DROP TABLE IF EXISTS facerec.persons;
*/
