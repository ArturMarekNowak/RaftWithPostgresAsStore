package entities

import (
	"main/internal/http/requests"
)

type Fsm struct {
	Id   string `gorm:"primaryKey"`
	Data string `gorm:"type:json"`
}

func NewFsm(dataRequest requests.DataRequest) Fsm {
	return Fsm{
		Id:   dataRequest.Key,
		Data: dataRequest.Value,
	}
}
