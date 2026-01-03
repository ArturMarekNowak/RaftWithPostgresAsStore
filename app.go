package main

import (
	"main/setup"
)

func main() {
	logger := setup.ConfigureLogger()
	httpPort, raftPort, raftId, dbName := setup.LoadArguments()
	db := setup.ConfigureDatabaseAndRunMigrations(dbName, logger)
	raft := setup.ConfigureRaft(raftPort, raftId, logger, db)
	setup.HttpServer(httpPort, raft, logger, db)
}
