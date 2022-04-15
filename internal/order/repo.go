package order

import (
	"github.com/gcamlicali/tradeshopExample/internal/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type OrderRepositoy struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepositoy {
	return &OrderRepositoy{db: db}
}

func (r *OrderRepositoy) Create(a *models.Order) (*models.Order, error) {
	zap.L().Debug("order.repo.create", zap.Reflect("orderBody", a))
	if err := r.db.Create(a).Error; err != nil {
		zap.L().Error("order.repo.Create failed to create order", zap.Error(err))
		return nil, err
	}
	return a, nil
}

func (r *OrderRepositoy) GetByUserID(userID int) (*[]models.Order, error) {
	zap.L().Debug("order.repo.GetByUserID", zap.Reflect("userID", userID))
	var orders []models.Order //duzelt dikkat
	err := r.db.Where(&models.Order{UserID: userID}).Find(&orders).Error
	if err != nil {
		zap.L().Error("order.repo.GetByUserID failed to get Orders", zap.Error(err))
		return nil, err
	}
	return &orders, nil

}

func (r *OrderRepositoy) Update(a *models.Order) (*models.Order, error) {

	zap.L().Debug("order.repo.update", zap.Reflect("orderBody", a))

	if result := r.db.Save(&a); result.Error != nil {
		return nil, result.Error
	}

	return a, nil
}
func (r *OrderRepositoy) Migration() {
	r.db.AutoMigrate(&models.Order{})
}
