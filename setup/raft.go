package setup

import (
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	"log"
	"main/internal/database"
	"net"
	"os"
	"time"
)

func ConfigureRaft(logger hclog.Logger, db *database.PostgresAccessor) *raft.Raft {
	raftAddress := os.Getenv("RAFT_ADDRESS")
	tcpAddr, err := net.ResolveTCPAddr("tcp", raftAddress)
	if err != nil {
		log.Fatal("Could not resolve address: %s", err)
	}

	transport, err := raft.NewTCPTransport(raftAddress, tcpAddr, 10, time.Second*10, os.Stderr)
	if err != nil {
		log.Fatal("Could not create tcp transport: %s", err)
	}

	nodeId := os.Getenv("RAFT_NODE_ID")
	raftConfig := raft.DefaultConfig()
	raftConfig.LocalID = raft.ServerID(nodeId)

	r, err := raft.NewRaft(raftConfig, db, db, db, db, transport)
	if err != nil {
		log.Fatal("Could not create raft instance: %s", err)
	}

	r.BootstrapCluster(raft.Configuration{
		Servers: []raft.Server{
			{
				ID:      raft.ServerID(nodeId),
				Address: transport.LocalAddr(),
			},
		}})

	return r
}
