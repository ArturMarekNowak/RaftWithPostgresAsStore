package database

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"main/internal/database/entities"
	"time"
)

type PostgresAccessor struct {
	Logger       hclog.Logger
	DatabaseName string
	db           *gorm.DB
}

func NewPostgresAccessor(dbName string, logger hclog.Logger) (*PostgresAccessor, error) {
	connectionString := fmt.Sprintf("postgresql://postgres:admin@127.0.0.1:5432/%s?sslmode=disable", dbName)

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{
		PrepareStmt: true, // Cache prepared statements for performance
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// Configure the connection pool
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	return &PostgresAccessor{
		Logger:       logger,
		DatabaseName: dbName,
		db:           db,
	}, nil
}

// SnapshotStore methods
func (p PostgresAccessor) Create(version raft.SnapshotVersion, index, term uint64, configuration raft.Configuration, configurationIndex uint64, trans raft.Transport) (raft.SnapshotSink, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresAccessor) List() ([]*raft.SnapshotMeta, error) {
	dbSnapshots := make([]*entities.Snapshot, 0)
	_ = p.db.Order("index desc").Find(&dbSnapshots)
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
	queryResult := p.db.First(&log)
	if queryResult.RowsAffected == 0 {
		return "{}", nil
	}
	return log.Data, nil
}

// LogStore methods
func (p PostgresAccessor) FirstIndex() (uint64, error) {
	log := entities.Log{}
	// TODO Select only Id
	queryResult := p.db.First(&log)
	if queryResult.RowsAffected == 0 {
		return 0, nil
	}
	return log.Index, nil
}

func (p PostgresAccessor) LastIndex() (uint64, error) {
	log := entities.Log{}
	// TODO Select only Id
	queryResult := p.db.Last(&log)
	if queryResult.RowsAffected == 0 {
		return 0, nil
	}
	return log.Index, nil
}

func (p PostgresAccessor) GetLog(index uint64, raftLog *raft.Log) error {
	log := entities.Log{}
	queryResult := p.db.First(&log, index)
	// TODO Not the most elegant, is it?
	log.FillRaftLog(raftLog)
	return queryResult.Error
}

func (p PostgresAccessor) StoreLog(raftLog *raft.Log) error {
	log := entities.NewLog(raftLog)
	// Don't miss the &
	// Source: https://stackoverflow.com/questions/59947933/err-reflect-flag-mustbeassignable-using-unaddressable-value-as-i-try-to-use-bin
	queryResult := p.db.Create(&log)
	return queryResult.Error
}

func (p PostgresAccessor) StoreLogs(logs []*raft.Log) error {
	// TODO Not the best idea to open connections so many times but we will stick with this for the time being
	for i := 0; i < len(logs); i++ {
		err := p.StoreLog(logs[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (p PostgresAccessor) DeleteRange(min, max uint64) error {
	//TODO implement me
	panic("implement me")
}

// StableStore methods
func (p PostgresAccessor) Set(key []byte, val []byte) error {
	num := binary.LittleEndian.Uint64(val)
	return p.SetUint64(key, num)
}

func (p PostgresAccessor) Get(key []byte) ([]byte, error) {
	valInt64, err := p.GetUint64(key)
	if err != nil {
		// stable.go:11
		return []byte{}, errors.New("not found")
	}
	val := []byte{}
	binary.LittleEndian.PutUint64(val, valInt64)
	return val, nil
}

// Source: https://stackoverflow.com/questions/39333102/how-to-create-or-update-a-record-with-gorm
func (p PostgresAccessor) SetUint64(key []byte, val uint64) error {
	stableLog := entities.StableLog{Id: key, Value: val}
	if p.db.Model(&stableLog).Where("id = ?", stableLog.Id).Updates(&stableLog).RowsAffected == 0 {
		p.db.Create(&stableLog)
	}
	return nil
}

func (p PostgresAccessor) GetUint64(key []byte) (uint64, error) {
	stableLog := entities.StableLog{
		Id: key,
	}
	queryResult := p.db.First(&stableLog)
	if queryResult.RowsAffected == 0 {
		// Error name required by raft@v1.7.3/api.go:510
		return 0, errors.New("not found")
	}
	return stableLog.Value, nil
}

func (p PostgresAccessor) RunMigrations() {
	err := p.db.AutoMigrate(&entities.Snapshot{},
		&entities.StableLog{},
		&entities.Log{},
		&entities.Fsm{})
	if err != nil {
		panic("Failed to run migrations")
	}
}
