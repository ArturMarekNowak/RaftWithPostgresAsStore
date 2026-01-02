package database

import (
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"main/internal/database/entities"
	"os"
)

type PostgresAccessor struct {
	Logger hclog.Logger
}

// SnapshotStore methods
func (p PostgresAccessor) Create(version raft.SnapshotVersion, index, term uint64, configuration raft.Configuration, configurationIndex uint64, trans raft.Transport) (raft.SnapshotSink, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresAccessor) List() ([]*raft.SnapshotMeta, error) {
	dbSnapshots := make([]*entities.Snapshot, 0)
	db := p.OpenConnection()
	_ = db.Order("index desc").Find(&dbSnapshots)
	// TODO extract it somewhere?
	raftSnapshots := make([]*raft.SnapshotMeta, 0)
	for i := 0; i < len(dbSnapshots); i++ {
		raftSnapshot, err := dbSnapshots[i].MapToRaftSnapshotMeta()
		if err != nil {
			p.Logger.Error("Parsing error at %d", i)
		}

		raftSnapshots = append(raftSnapshots, raftSnapshot)
	}
	return raftSnapshots, nil
}

func (p PostgresAccessor) Open(id string) (*raft.SnapshotMeta, io.ReadCloser, error) {
	//TODO implement me
	panic("implement me")
}

// FSM methods
func (p PostgresAccessor) Apply(log *raft.Log) interface{} {
	//TODO implement me
	panic("implement me")
}

func (p PostgresAccessor) Snapshot() (raft.FSMSnapshot, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresAccessor) Restore(snapshot io.ReadCloser) error {
	//TODO implement me
	panic("implement me")
}

// Not a part of hashicorp/raft
func (p PostgresAccessor) GetValue(key uint64) (string, error) {
	log := entities.Fsm{
		Id: key,
	}
	db := p.OpenConnection()
	queryResult := db.First(&log)
	if queryResult.RowsAffected == 0 {
		return "{}", nil
	}
	return log.Data, nil
}

// LogStore methods
func (p PostgresAccessor) FirstIndex() (uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresAccessor) LastIndex() (uint64, error) {
	log := entities.Log{}
	db := p.OpenConnection()
	queryResult := db.Last(&log)
	if queryResult.RowsAffected == 0 {
		return 0, nil
	}
	return log.Id, nil
}

func (p PostgresAccessor) GetLog(index uint64, log *raft.Log) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgresAccessor) StoreLog(log *raft.Log) error {
	//TODO implement me
	// Tu skończyłeś
	panic("implement me")
}

func (p PostgresAccessor) StoreLogs(logs []*raft.Log) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgresAccessor) DeleteRange(min, max uint64) error {
	//TODO implement me
	panic("implement me")
}

// StableStore methods
func (p PostgresAccessor) Set(key []byte, val []byte) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgresAccessor) Get(key []byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

// Source: https://stackoverflow.com/questions/39333102/how-to-create-or-update-a-record-with-gorm
func (p PostgresAccessor) SetUint64(key []byte, val uint64) error {
	stableLog := entities.StableLog{Id: key, Value: val}
	db := p.OpenConnection()
	if db.Model(&stableLog).Where("id = ?", stableLog.Id).Updates(&stableLog).RowsAffected == 0 {
		db.Create(&stableLog)
	}
	return nil
}

func (p PostgresAccessor) GetUint64(key []byte) (uint64, error) {
	stableLog := entities.StableLog{
		Id: key,
	}
	db := p.OpenConnection()
	queryResult := db.First(&stableLog)
	if queryResult.RowsAffected == 0 {
		return 0, nil
	}
	return stableLog.Value, nil
}

func (p PostgresAccessor) RunMigrations() {

	db := p.OpenConnection()

	err := db.AutoMigrate(&entities.Snapshot{},
		&entities.StableLog{},
		&entities.Log{},
		&entities.Fsm{})
	if err != nil {
		panic("Failed to run migrations")
	}
}

func (p PostgresAccessor) OpenConnection() *gorm.DB {
	connectionString := os.Getenv("CONNECTION_STRING")
	db, err := gorm.Open(postgres.Open(connectionString))
	if err != nil {
		p.Logger.Error("Failed to connect database")
	}
	return db
}
