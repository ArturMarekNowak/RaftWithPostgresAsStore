package entities

type Log struct {
	Id   uint64 `gorm:"primaryKey"`
	Data string `gorm:"type:json"`
}
