package database

import (
	"fmt"
	"github.com/lopezator/migrator"
	"strings"
)

func migrations(schema, migrationsTable string) (*migrator.Migrator, error) {
	l := migrator.WithLogger(migrator.LoggerFunc(func(msg string, args ...interface{}) {
		if !strings.Contains(msg, "applied migration named") &&
			!strings.Contains(msg, "applying migration named") &&
			!strings.Contains(msg, "applied no tx migration named") &&
			!strings.Contains(msg, "applying no tx migration named") {
			fmt.Printf(msg+"\n", args)
		}
	}))

	return migrator.New(l,
		migrator.TableName(fmt.Sprintf("%s.%s", schema, migrationsTable)),
		migrator.Migrations(
			migrationInit(),
		),
	)
}
