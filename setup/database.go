package setup

import (
	"github.com/hashicorp/go-hclog"
	"main/internal/database"
)

func ConfigureDatabaseAndRunMigrations(logger hclog.Logger) *database.PostgresAccessor {
	db := &database.PostgresAccessor{
		Logger: logger,
	}
	db.RunMigrations()
	return db
}
