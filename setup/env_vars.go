package setup

import "os"

func LoadArguments() (string, string, string, string) {
	var httpPort, raftPort, raftId, databaseName string
	for i, arg := range os.Args[1:] {
		if arg == "--raft-id" {
			raftId = os.Args[i+2]
			continue
		}

		if arg == "--http-port" {
			httpPort = os.Args[i+2]
			continue
		}

		if arg == "--raft-port" {
			raftPort = os.Args[i+2]
			continue
		}

		if arg == "--db-name" {
			databaseName = os.Args[i+2]
			continue
		}
	}

	if httpPort == "" || raftPort == "" || raftId == "" || databaseName == "" {
		panic("must provide --raft-id, --http-port, --raft-port, --db-name as arguments")
	}

	return httpPort, raftPort, raftId, databaseName
}
