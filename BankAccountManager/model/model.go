package model

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type Account struct {
	Id          uint64 `gorm:"primaryKey"`
	Name        string
	Address     string
	PhoneNumber string
	Balance     float32
	Status      int
	Createtime  *timestamppb.Timestamp
}
type AccountORM struct {
	Id          uint64 `gorm:"primaryKey"`
	Name        string
	Address     string
	PhoneNumber string
	Balance     float32
	Status      int
	Createtime  *time.Time
}
type Transaction struct {
	From   uint64
	Amount float64
	To     uint64
}
