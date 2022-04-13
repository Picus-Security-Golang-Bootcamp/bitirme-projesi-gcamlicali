package models

import "gorm.io/gorm"

type CartItem struct {
	gorm.Model
	Quantity  int
	Product   Product `gorm:"ForeignKey:ID;references:ProductID"`
	ProductID int
	CartID    int
	Price     int
}

func (CartItem) TableName() string {
	//default table name
	return "cart_item"
}
