package cart_item

import (
	"github.com/gcamlicali/tradeshopExample/internal/models"
	"gorm.io/gorm"
)

type CartItemRepositoy struct {
	db *gorm.DB
}

func NewCartItemRepository(db *gorm.DB) *CartItemRepositoy {
	return &CartItemRepositoy{db: db}
}

func (r *CartItemRepositoy) Migration() {
	r.db.AutoMigrate(&models.CartItem{})
}
