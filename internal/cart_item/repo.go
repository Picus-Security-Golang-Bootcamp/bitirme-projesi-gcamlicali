package cart_item

import (
	"github.com/gcamlicali/tradeshopExample/internal/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CartItemRepositoy struct {
	db *gorm.DB
}

func NewCartItemRepository(db *gorm.DB) *CartItemRepositoy {
	return &CartItemRepositoy{db: db}
}

func (ci *CartItemRepositoy) Crate(a *models.CartItem) (*models.CartItem, error) {
	zap.L().Debug("cartitem.repo.create", zap.Reflect("cartBody", a))
	if err := ci.db.Create(a).Error; err != nil {
		zap.L().Error("cart.repo.Create failed to create cart", zap.Error(err))
		return nil, err
	}
	return a, nil

}
func (ci *CartItemRepositoy) GetByCartID(cartID int) ([]models.CartItem, error) {
	zap.L().Debug("cartitem.repo.getByCartID", zap.Reflect("CartID", cartID))
	var cartItems = []models.CartItem{}
	err := ci.db.Where(&models.CartItem{CartID: cartID}).Find(&cartItems).Error
	if err != nil {
		zap.L().Error("cart.repo.GetByCartID failed to get CartItems", zap.Error(err))
		return nil, err
	}
	return cartItems, nil

}

func (ci *CartItemRepositoy) Migration() {
	ci.db.AutoMigrate(&models.CartItem{})
}
