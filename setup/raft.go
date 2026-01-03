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
	advertiseAddressEnv := os.Getenv("RAFT_ADVERTISE_ADDRESS")
	advertiseAddress, err := net.ResolveTCPAddr("tcp", advertiseAddressEnv)
	if err != nil {
		log.Fatal("Could not resolve address: %s", err)
	}

	bindAddressEnv := os.Getenv("RAFT_BIND_ADDRESS")

	transport, err := raft.NewTCPTransport(bindAddressEnv, advertiseAddress, 10, time.Second*10, os.Stderr)
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
