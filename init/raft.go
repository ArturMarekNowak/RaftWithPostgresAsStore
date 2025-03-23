package init

import (
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"log"
	raftutil "main/internal/raft"
	"net"
	"os"
	"path"
	"time"
)

func SetupRaft() *raft.Raft {
	fsm := &raftutil.Fsm{
		StateValue: 0,
	}

	dir := os.Getenv("BOLTDB_STORE_PATH")
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Fatal("Could not create data directory: %s", err)
	}

	store, err := raftboltdb.NewBoltStore(path.Join(dir, "bolt"))
	if err != nil {
		log.Fatal("Could not create bolt store: %s", err)
	}

	snapshots, err := raft.NewFileSnapshotStore(path.Join(dir, "snapshot"), 2, os.Stderr)
	if err != nil {
		log.Fatal("Could not create snapshot store: %s", err)
	}

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
	raftCfg := raft.DefaultConfig()
	raftCfg.LocalID = raft.ServerID(nodeId)

	r, err := raft.NewRaft(raftCfg, fsm, store, store, snapshots, transport)
	if err != nil {
		log.Fatal("Could not create raft instance: %s", err)
	}

	// Cluster consists of unjoined leaders. Picking a leader and
	// creating a real cluster is done manually after startup.
	r.BootstrapCluster(raft.Configuration{
		Servers: []raft.Server{
			{
				ID:      raft.ServerID(nodeId),
				Address: transport.LocalAddr(),
			},
		},
	})

	return r
}
