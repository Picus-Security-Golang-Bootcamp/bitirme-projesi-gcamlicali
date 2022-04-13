package models

import "gorm.io/gorm"

type CartItem struct {
	gorm.Model
	Quantity  int
	ProductID int
	CartID    int
}

func (CartItem) TableName() string {
	//default table name
	return "cart_item"
}
