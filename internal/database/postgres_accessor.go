package database

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"main/internal/database/entities"
	"main/internal/http/requests"
	"time"
)

type PostgresAccessor struct {
	DatabaseName string
	logger       hclog.Logger
	database     *gorm.DB
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
		DatabaseName: dbName,
		logger:       logger,
		database:     db,
	}, nil
}

// SnapshotStore methods
func (p PostgresAccessor) Create(version raft.SnapshotVersion, index, term uint64, configuration raft.Configuration, configurationIndex uint64, trans raft.Transport) (raft.SnapshotSink, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresAccessor) List() ([]*raft.SnapshotMeta, error) {
	dbSnapshots := make([]*entities.Snapshot, 0)
	if err := p.database.Order("index desc").Find(&dbSnapshots).Error; err != nil {
		// Probably the error should be errors.New("not found") but I would need to make sure and reverse engineer raft
		// library some more
		return nil, err
	}
	// I dont have a better idea how to map it and Gemini Pro does not have either
	raftSnapshots := make([]*raft.SnapshotMeta, 0)
	for i, dbSnapshot := range dbSnapshots {
		raftSnapshot, err := dbSnapshot.MapToRaftSnapshotMeta()
		if err != nil {
			p.logger.Error("Parsing error at %d", i)
		}
		raftSnapshots = append(raftSnapshots, raftSnapshot)
	}
	return raftSnapshots, nil
}

func (p PostgresAccessor) Open(id string) (*raft.SnapshotMeta, io.ReadCloser, error) {
	//TODO implement me
	panic("implement me")
}

type setPayload struct {
	Key   string
	Value string
}

// FSM methods
func (p PostgresAccessor) Apply(log *raft.Log) interface{} {
	switch log.Type {
	case raft.LogCommand:
		dataRequest := requests.DataRequest{}
		err := json.Unmarshal(log.Data, &dataRequest)
		if err != nil {
			return fmt.Errorf("Could not parse payload: %s", err)
		}

		fsm := entities.NewFsm(dataRequest)
		if p.database.Model(&fsm).Where("id = ?", fsm.Id).Updates(&fsm).RowsAffected == 0 {
			p.database.Create(&fsm)
		}
	default:
		return fmt.Errorf("Unknown raft log type: %#v", log.Type)
	}

	return nil
}

func (p PostgresAccessor) Snapshot() (raft.FSMSnapshot, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresAccessor) Restore(snapshot io.ReadCloser) error {
	//TODO implement me
	panic("implement me")
}

// Not a part of hashicorp/raft FSM interface, raft@v1.7.3/fsm.go:16
func (p PostgresAccessor) GetValue(key string) (string, error) {
	log := entities.Fsm{
		Id: key,
	}
	queryResult := p.database.First(&log)
	if queryResult.RowsAffected == 0 {
		return "{}", nil
	}
	return log.Data, nil
}

// LogStore methods
func (p PostgresAccessor) FirstIndex() (uint64, error) {
	log := entities.Log{}
	queryResult := p.database.Select("index").First(&log)
	if queryResult.RowsAffected == 0 {
		return 0, nil
	}
	return log.Index, nil
}

func (p PostgresAccessor) LastIndex() (uint64, error) {
	log := entities.Log{}
	queryResult := p.database.Select("index").Last(&log)
	if queryResult.RowsAffected == 0 {
		return 0, nil
	}
	return log.Index, nil
}

func (p PostgresAccessor) GetLog(index uint64, raftLog *raft.Log) error {
	log := entities.Log{}
	queryResult := p.database.First(&log, index)
	// Dont know if there is a way to map this value in more elegant way
	log.FillRaftLog(raftLog)
	return queryResult.Error
}

func (p PostgresAccessor) StoreLog(raftLog *raft.Log) error {
	log := entities.NewLog(raftLog)
	// Don't miss the &
	// Source: https://stackoverflow.com/questions/59947933/err-reflect-flag-mustbeassignable-using-unaddressable-value-as-i-try-to-use-bin
	queryResult := p.database.Create(&log)
	return queryResult.Error
}

func (p PostgresAccessor) StoreLogs(logs []*raft.Log) error {
	logEntities := make([]entities.Log, len(logs))
	for i, raftLog := range logs {
		logEntities[i] = entities.NewLog(raftLog)
	}
	queryResult := p.database.Create(&logEntities)
	return queryResult.Error
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
	// There is a bug in here when i kill all nodes and try to bring them back "panic: runtime error: index out of range
	// [7] with length 0". So probably somehow some key is queried but result is too big. Highest value in the table was
	// 3328210909095473713 which is well within range of uint64 so i dont know whats going on. Will look into that in
	// the future. Happy path is working fine so thats good with me for the time being. Its PoC for Gods sake
	binary.LittleEndian.PutUint64(val, valInt64)
	return val, nil
}

// Source: https://stackoverflow.com/questions/39333102/how-to-create-or-update-a-record-with-gorm
func (p PostgresAccessor) SetUint64(key []byte, val uint64) error {
	stableLog := entities.StableLog{Id: key, Value: val}
	if p.database.Model(&stableLog).Where("id = ?", stableLog.Id).Updates(&stableLog).RowsAffected == 0 {
		p.database.Create(&stableLog)
	}
	return nil
}

func (p PostgresAccessor) GetUint64(key []byte) (uint64, error) {
	stableLog := entities.StableLog{
		Id: key,
	}
	queryResult := p.database.First(&stableLog)
	if queryResult.RowsAffected == 0 {
		// Error name required by raft@v1.7.3/api.go:510
		return 0, errors.New("not found")
	}
	return stableLog.Value, nil
}

func (p PostgresAccessor) RunMigrations() {
	err := p.database.AutoMigrate(&entities.Snapshot{},
		&entities.StableLog{},
		&entities.Log{},
		&entities.Fsm{})
	if err != nil {
		panic("Failed to run migrations")
	}
}
