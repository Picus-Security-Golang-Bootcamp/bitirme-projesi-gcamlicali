package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Order struct {
	ID         uuid.UUID `gorm:"primary_key; type:uuid; default:uuid_generate_v4()"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	CartID     uuid.UUID
	UserID     uuid.UUID
	Status     string
	Cart       Cart
	TotalPrice int32
}

func (Order) TableName() string {
	//default table name
	return "order"
}
