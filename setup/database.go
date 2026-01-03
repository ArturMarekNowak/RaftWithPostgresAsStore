package setup

import (
	"github.com/hashicorp/go-hclog"
	"main/internal/database"
)

func ConfigureDatabaseAndRunMigrations(dbName string, logger hclog.Logger) *database.PostgresAccessor {
	db, err := database.NewPostgresAccessor(dbName, logger)
	if err != nil {
		panic("Couldnt initialize db")
	}
	db.RunMigrations()
	return db
}
