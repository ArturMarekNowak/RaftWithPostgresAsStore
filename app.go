package main

import (
	"main/setup"
)

func main() {
	logger := setup.ConfigureLogger()
	setup.LoadEnvironmentalVariables(logger)
	db := setup.ConfigureDatabaseAndRunMigrations(logger)
	raft := setup.ConfigureRaft(logger, db)
	setup.HttpServer(raft, logger, db)
}
