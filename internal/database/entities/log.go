package entities

import (
	"github.com/hashicorp/raft"
	"time"
)

type Log struct {
	Index      uint64 `gorm:"primaryKey"`
	Term       uint64
	LogType    string `gorm:"type:varchar(50)"`
	Data       []byte
	Extensions []byte
	AppendedAt time.Time
}

// Source: https://tutorialedge.net/golang/go-constructors-tutorial/
func NewLog(log *raft.Log) Log {
	return Log{
		Index:      log.Index,
		Term:       log.Term,
		LogType:    string(log.Type),
		Data:       log.Data,
		Extensions: log.Extensions,
		AppendedAt: log.AppendedAt,
	}
}
