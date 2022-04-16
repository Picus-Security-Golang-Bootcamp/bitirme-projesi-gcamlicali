package cart_item

import (
	"github.com/gcamlicali/tradeshopExample/internal/models"
	"github.com/google/uuid"
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
		zap.L().Error("cartitem.repo.Create failed to create CartItem", zap.Error(err))
		return nil, err
	}
	return a, nil

}
func (ci *CartItemRepositoy) GetByCartID(cartID uuid.UUID) ([]models.CartItem, error) {
	zap.L().Debug("cartitem.repo.getByCartID", zap.Reflect("CartID", cartID))
	var cartItems = []models.CartItem{}
	err := ci.db.Where(&models.CartItem{CartID: cartID}).Find(&cartItems).Error
	if err != nil {
		zap.L().Error("cartitem.repo.GetByCartID failed to get CartItems", zap.Error(err))
		return nil, err
	}
	return cartItems, nil

}

func (ci *CartItemRepositoy) GetByCartAndProductSKU(cartID uuid.UUID, productSKU int) (*models.CartItem, error) {
	zap.L().Debug("cartitem.repo.getByCartID", zap.Reflect("CartID", cartID))
	cartItem := models.CartItem{}
	err := ci.db.Where(&models.CartItem{CartID: cartID, ProductSKU: productSKU}).First(&cartItem).Error
	if err != nil {
		zap.L().Error("cartitem.repo.GetByProductID failed to get CartItems", zap.Error(err))
		return nil, err
	}
	return &cartItem, nil

}

func (ci *CartItemRepositoy) Update(a *models.CartItem) (*models.CartItem, error) {
	zap.L().Debug("cartitem.repo.update", zap.Reflect("cartBody", a))
	if err := ci.db.Save(a).Error; err != nil {
		zap.L().Error("cartitem.repo.Update failed to update CartItem", zap.Error(err))
		return nil, err
	}
	return a, nil
}

func (ci *CartItemRepositoy) Delete(a *models.CartItem) error {
	zap.L().Debug("cartitem.repo.delete", zap.Reflect("cartBody", a))

	if err := ci.db.Delete(&a); err.Error != nil {
		zap.L().Error("cartitem.repo.Delete failed to delete CartItem", zap.Error(err.Error))
		return err.Error
	}

	return nil
}

//func (ci *CartItemRepositoy) DeleteBy(a *models.CartItem) error {
//	zap.L().Debug("cartitem.repo.delete", zap.Reflect("cartBody", a))
//
//	if err := ci.db.Delete(&a); err.Error != nil {
//		zap.L().Error("cartitem.repo.Delete failed to delete CartItem", zap.Error(err.Error))
//		return err.Error
//	}
//
//	return nil
//}

func (ci *CartItemRepositoy) Migration() {
	ci.db.AutoMigrate(&models.CartItem{})
}
