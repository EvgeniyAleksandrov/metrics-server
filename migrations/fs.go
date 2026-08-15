package migrations

import "embed"

//go:embed *.sql
var migrationsFS embed.FS

func GetMigrationFS() embed.FS {
	return migrationsFS
}
