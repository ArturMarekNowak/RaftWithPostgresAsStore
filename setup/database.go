package setup

import (
	"github.com/hashicorp/go-hclog"
	databaseutil "main/internal/database"
)

func ConfigureDatabaseAndRunMigrations(databaseName string, logger hclog.Logger) *databaseutil.PostgresAccessor {
	database, err := databaseutil.NewPostgresAccessor(databaseName, logger)
	if err != nil {
		panic("Could not initialize the postgres database")
	}
	database.RunMigrations()
	return database
}
