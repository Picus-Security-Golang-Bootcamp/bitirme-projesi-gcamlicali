package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	CategoryName string // Category Category
	SKU          int    `gorm:"unique"`
	Name         string
	Description  string
	UnitStock    int32
	Price        int32
}

func (Product) TableName() string {
	//default table name
	return "products"
}
