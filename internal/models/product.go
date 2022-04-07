package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	CategoryID  int64 // Category Category
	SKU         int64 `gorm:"unique"`
	Name        string
	Description string
}

func (Product) TableName() string {
	//default table name
	return "products"
}
