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

func ConfigureRaft(raftPort, raftId string, logger hclog.Logger, db *database.PostgresAccessor) *raft.Raft {
	tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:"+raftPort)
	if err != nil {
		logger.Error("Could not resolve address: %s", err)
	}

	transport, err := raft.NewTCPTransport("127.0.0.1:"+raftPort, tcpAddr, 10, time.Second*10, os.Stderr)
	if err != nil {
		log.Fatal("Could not create tcp transport: %s", err)
	}

	raftConfig := raft.DefaultConfig()
	raftConfig.LocalID = raft.ServerID(raftId)

	r, err := raft.NewRaft(raftConfig, db, db, db, db, transport)
	if err != nil {
		logger.Error("Could not create raft instance: %s", err)
	}

	r.BootstrapCluster(raft.Configuration{
		Servers: []raft.Server{
			{
				ID:      raft.ServerID(raftId),
				Address: transport.LocalAddr(),
			},
		}})

	return r
}
