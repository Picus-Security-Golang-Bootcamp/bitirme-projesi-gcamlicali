package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	CartID     int
	UserID     int
	Status     string
	Cart       Cart
	TotalPrice int32
}

func (Order) TableName() string {
	//default table name
	return "order"
}
