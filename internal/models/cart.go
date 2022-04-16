package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Cart struct {
	ID         uuid.UUID `gorm:"primary_key; type:uuid; default:uuid_generate_v4()"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	IsOrdered  bool
	UserID     uuid.UUID
	CartItems  []CartItem `gorm:"ForeignKey:CartID"`
	TotalPrice int
}

func (Cart) TableName() string {
	//default table name
	return "cart"
}
