package setup

import (
	"github.com/hashicorp/go-hclog"
	"main/internal/database"
)

func ConfigureDatabaseAndRunMigrations(dbName string, logger hclog.Logger) *database.PostgresAccessor {
	db := &database.PostgresAccessor{
		Logger:       logger,
		DatabaseName: dbName,
	}
	db.RunMigrations()
	return db
}
