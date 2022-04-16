package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Product struct {
	ID           uuid.UUID `gorm:"primary_key; type:uuid; default:uuid_generate_v4()"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	CategoryName string
	SKU          int `gorm:"unique"`
	Name         string
	Description  string
	UnitStock    int32
	Price        int32
}

func (Product) TableName() string {
	//default table name
	return "products"
}
