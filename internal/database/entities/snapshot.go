package entities

import (
	"encoding/json"
	"github.com/hashicorp/raft"
)

type Snapshot struct {
	Index              uint64 `gorm:"primaryKey"`
	ConfigurationIndex uint64 `gorm:"primaryKey"`
	Term               uint64 `gorm:"type:int"`
	Version            int    `gorm:"type:int"`
	Configuration      string `gorm:"type:json"`
	Transport          string `gorm:"type:json"`
}

func (s Snapshot) MapToRaftSnapshotMeta() (*raft.SnapshotMeta, error) {
	var raftConfig raft.Configuration
	config := []byte(s.Configuration)
	if err := json.Unmarshal(config, &raftConfig); err != nil {
		return nil, err
	}

	return &raft.SnapshotMeta{
		Version:            raft.SnapshotVersion(s.Version),
		ID:                 "", // Opaque to the store
		Index:              s.Index,
		Term:               s.Term,
		Peers:              nil, // Deprecated
		Configuration:      raftConfig,
		ConfigurationIndex: s.ConfigurationIndex,
		Size:               0,
	}, nil
}
