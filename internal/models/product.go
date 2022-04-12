package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	CategoryName string // Category Category
	SKU          string `gorm:"unique"`
	Name         string
	Description  string
}

func (Product) TableName() string {
	//default table name
	return "products"
}
