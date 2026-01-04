package entities

type StableLog struct {
	Id    []byte `gorm:"primaryKey;type:bytea"`
	Value uint64 `gorm:"type:numeric(20, 0)"`
}
