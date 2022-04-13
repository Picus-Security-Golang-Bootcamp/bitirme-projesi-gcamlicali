package cart

import (
	"github.com/gcamlicali/tradeshopExample/internal/models"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CartRepositoy struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) *CartRepositoy {
	return &CartRepositoy{db: db}
}

func (r *CartRepositoy) Create(a *models.Cart) (*models.Cart, error) {
	zap.L().Debug("cart.repo.create", zap.Reflect("cartBody", a))
	if err := r.db.Create(a).Error; err != nil {
		zap.L().Error("cart.repo.Create failed to create cart", zap.Error(err))
		return nil, err
	}
	return a, nil
}

func (r *CartRepositoy) GetByUserID(userID int) (*models.Cart, error) {
	zap.L().Debug("cart.repo.GetByUserID", zap.Reflect("userID", userID))

	var cart = &models.Cart{}

	err := r.db.Table("cart").Preload("CartItems").Preload("CartItems.Product").Where(&models.Cart{UserID: userID, IsOrdered: false}).First(&cart).Error
	//err := r.db.Preload("CartItems").Where("user_id =?", userID).First(&cart).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		cart.UserID = userID
		newCart, err := r.Create(cart)
		if err != nil {
			return nil, err
		}
		return newCart, nil
	}
	if err != nil {
		return nil, err
	}
	if cart == nil {

	}
	return cart, nil
}

func (r *CartRepositoy) Update(a *models.Cart) (*models.Cart, error) {

	zap.L().Debug("cart.repo.update", zap.Reflect("cartBody", a))

	if result := r.db.Save(&a); result.Error != nil {
		return nil, result.Error
	}

	return a, nil
}

func (r *CartRepositoy) Migration() {
	r.db.AutoMigrate(&models.Cart{})
}
