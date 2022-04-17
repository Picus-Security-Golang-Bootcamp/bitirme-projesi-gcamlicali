package product

import (
	"github.com/gcamlicali/tradeshopExample/internal/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ProductRepositoy struct {
	db *gorm.DB
}

type IProductRepository interface {
	Create(a *models.Product) (*models.Product, error)
	GetAll(pageIndex, pageSize int) (*[]models.Product, int, error)
	GetByName(name string) (*[]models.Product, error)
	GetBySKU(sku int) (*models.Product, error)
	Update(a *models.Product) (*models.Product, error)
	Delete(sku int) error
}

func NewProductRepository(db *gorm.DB) *ProductRepositoy {
	return &ProductRepositoy{db: db}
}

func (r *ProductRepositoy) Create(a *models.Product) (*models.Product, error) {
	zap.L().Debug("product.repo.create", zap.Reflect("productBody", a))
	if err := r.db.Create(a).Error; err != nil {
		zap.L().Error("product.repo.Create failed to create product", zap.Error(err))
		return nil, err
	}
	return a, nil
}

func (r *ProductRepositoy) GetAll(pageIndex, pageSize int) (*[]models.Product, int, error) {
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

func (r *ProductRepositoy) GetByName(name string) (*[]models.Product, error) {
	zap.L().Debug("product.repo.getByName", zap.Reflect("name", name))

	var products = &[]models.Product{}

	err := r.db.Where("name ILIKE ? ", "%"+name+"%").
		Find(&products).Error
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepositoy) GetBySKU(sku int) (*models.Product, error) {
	zap.L().Debug("product.repo.getBySKU", zap.Reflect("SKU", sku))

	var product = &models.Product{}
	err := r.db.Where(&models.Product{SKU: sku}).First(&product).Error
	if err != nil {
		return nil, err
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

func (r *ProductRepositoy) Delete(sku int) error {
	zap.L().Debug("product.repo.deleteBySku", zap.Reflect("SKU", sku))

	product, err := r.GetBySKU(sku)
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
