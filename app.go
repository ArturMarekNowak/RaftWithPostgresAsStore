package main

import (
	"main/init"
)

func main() {
	init.LoadEnvironmentalVariables()
	logger := init.ConfigureLogger()
	raft := init.SetupRaft()
	logger.Info("Node started")
}
