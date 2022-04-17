package category

import (
	"github.com/gcamlicali/tradeshopExample/internal/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CategoryRepositoy struct {
	db *gorm.DB
}

type ICategoryRepository interface {
	Create(a *models.Category) (*models.Category, error)
	GetByName(name string) (*models.Category, error)
	GetAll(pageIndex, pageSize int) (*[]models.Category, int, error)
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepositoy {
	return &CategoryRepositoy{db: db}
}

func (r *CategoryRepositoy) Create(a *models.Category) (*models.Category, error) {
	zap.L().Debug("category.repo.create", zap.Reflect("categoryBody", a))
	if err := r.db.Create(a).Error; err != nil {
		zap.L().Error("category.repo.Create failed to create category", zap.Error(err))
		return nil, err
	}
	return a, nil
}

func (r *CategoryRepositoy) GetByName(name string) (*models.Category, error) {
	zap.L().Debug("category.repo.getByName", zap.Reflect("name", name))
	var category = &models.Category{}
	if result := r.db.Where("Name=?", name).First(&category); result.Error != nil {
		return nil, result.Error
	}

	return category, nil
}

func (r *CategoryRepositoy) GetAll(pageIndex, pageSize int) (*[]models.Category, int, error) {
	zap.L().Debug("category.repo.getAll")

	var categories = &[]models.Category{}
	var junk = &[]models.Category{}
	var count int64
	if err := r.db.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&categories).Error; err != nil {
		return nil, 0, err
	}
	r.db.Find(&junk).Count(&count)
	junk = nil
	return categories, int(count), nil
}

func (r *CategoryRepositoy) Migration() {
	r.db.AutoMigrate(&models.Category{})
}
