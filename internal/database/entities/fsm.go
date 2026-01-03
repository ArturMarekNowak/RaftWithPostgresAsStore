package entities

type Fsm struct {
	Id   string `gorm:"primaryKey"`
	Data string `gorm:"type:json"`
}
