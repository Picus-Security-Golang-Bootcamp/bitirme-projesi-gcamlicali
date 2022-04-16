package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type CartItem struct {
	ID         uuid.UUID `gorm:"primary_key; type:uuid; default:uuid_generate_v4()"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	Quantity   int
	Product    Product `gorm:"ForeignKey:SKU;references:ProductSKU"`
	ProductSKU int
	CartID     uuid.UUID
	Price      int
}

func (CartItem) TableName() string {
	//default table name
	return "cart_item"
}
