package entities

import (
	"github.com/hashicorp/raft"
	"time"
)

type Log struct {
	Index      uint64 `gorm:"primaryKey"`
	Term       uint64
	LogType    uint8
	Data       []byte
	Extensions []byte
	AppendedAt time.Time
}

func (l Log) FillRaftLog(raftLog *raft.Log) {
	raftLog.Index = l.Index
	raftLog.Term = l.Term
	raftLog.Type = raft.LogType(l.LogType)
	raftLog.Data = l.Data
	raftLog.Extensions = l.Extensions
	raftLog.AppendedAt = l.AppendedAt
}

// Source: https://tutorialedge.net/golang/go-constructors-tutorial/
func NewLog(log *raft.Log) Log {
	return Log{
		Index:      log.Index,
		Term:       log.Term,
		LogType:    uint8(log.Type),
		Data:       log.Data,
		Extensions: log.Extensions,
		AppendedAt: log.AppendedAt,
	}
}
