package models

import (
	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model
	IsOrdered bool
	UserID    int
	CartItems []CartItem `gorm:"ForeignKey:CartID"`
}

func (Cart) TableName() string {
	//default table name
	return "cart"
}
