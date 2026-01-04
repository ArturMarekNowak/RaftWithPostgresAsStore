package main

import (
	"main/setup"
)

func main() {
	logger := setup.ConfigureLogger()
	httpPort, raftPort, raftId, dbName := setup.LoadArguments()
	database := setup.ConfigureDatabaseAndRunMigrations(dbName, logger)
	raft := setup.ConfigureRaft(raftPort, raftId, logger, database)
	setup.HttpServer(httpPort, raft, logger, database)
}
