package product

import (
	"github.com/gcamlicali/tradeshopExample/internal/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ProductRepositoy struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepositoy {
	return &ProductRepositoy{db: db}
}

func (r *ProductRepositoy) create(a *models.Product) (*models.Product, error) {
	zap.L().Debug("product.repo.create", zap.Reflect("productBody", a))
	if err := r.db.Create(a).Error; err != nil {
		zap.L().Error("product.repo.Create failed to create product", zap.Error(err))
		return nil, err
	}
	return a, nil
}

func (r *ProductRepositoy) getAll(pageIndex, pageSize int) (*[]models.Product, int, error) {
	zap.L().Debug("product.repo.getAll")

	var ps = &[]models.Product{}
	var junk = &[]models.Product{}
	var count int64

	if err := r.db.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&ps).Error; err != nil {
		zap.L().Error("product.repo.getAll failed to get products", zap.Error(err))
		return nil, 0, err
	}
	r.db.Find(&junk).Count(&count)
	junk = nil
	return ps, int(count), nil
}

func (r *ProductRepositoy) GetByID(id int) (*models.Product, error) {
	zap.L().Debug("product.repo.getByID", zap.Reflect("id", id))

	var product = &models.Product{}
	if result := r.db.First(&product, id); result.Error != nil {
		return nil, result.Error
	}

	return product, nil
}

func (r *ProductRepositoy) Update(a *models.Product) (*models.Product, error) {
	zap.L().Debug("product.repo.update", zap.Reflect("product", a))

	if result := r.db.Save(&a); result.Error != nil {
		return nil, result.Error
	}

	return a, nil
}

func (r *ProductRepositoy) delete(id int) error {
	zap.L().Debug("product.repo.delete", zap.Reflect("id", id))

	product, err := r.GetByID(id)
	if err != nil {
		return err
	}

	if result := r.db.Delete(&product); result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *ProductRepositoy) Migration() {
	r.db.AutoMigrate(&models.Product{})
}
