package entities

type Fsm struct {
	Id   uint64 `gorm:"primaryKey"`
	Data string `gorm:"type:json"`
}
