package raft

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/raft"
	"io"
	"sync"
)

type Fsm struct {
	mutex      sync.Mutex
	StateValue int
}

type event struct {
	Type  string
	Value int
}

// Apply applies a Raft log entry to the key-value store.
func (fsm *Fsm) Apply(logEntry *raft.Log) interface{} {
	var e event
	if err := json.Unmarshal(logEntry.Data, &e); err != nil {
		panic("Failed unmarshaling Raft log entry. This is a bug.")
	}

	switch e.Type {
	case "set":
		fsm.mutex.Lock()
		defer fsm.mutex.Unlock()
		fsm.StateValue = e.Value

		return nil
	default:
		panic(fmt.Sprintf("Unrecognized event type in Raft log entry: %s. This is a bug.", e.Type))
	}
}

func (fsm *Fsm) Snapshot() (raft.FSMSnapshot, error) {
	fsm.mutex.Lock()
	defer fsm.mutex.Unlock()

	return &fsmSnapshot{stateValue: fsm.StateValue}, nil
}

// Restore stores the key-value store to a previous state.
func (fsm *Fsm) Restore(serialized io.ReadCloser) error {
	var snapshot fsmSnapshot
	if err := json.NewDecoder(serialized).Decode(&snapshot); err != nil {
		return err
	}

	fsm.StateValue = snapshot.stateValue
	return nil
}
