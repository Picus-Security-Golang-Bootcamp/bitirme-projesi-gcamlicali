package category

import (
	"github.com/gcamlicali/tradeshopExample/internal/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CategoryRepositoy struct {
	db *gorm.DB
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

func (r *CategoryRepositoy) GetByID(id string) (*models.Category, error) {
	zap.L().Debug("category.repo.getByID", zap.Reflect("id", id))

	var category = &models.Category{}
	if err := r.db.First(&category, id).Error; err != nil {
		return nil, err
	}

	return category, nil
}

func (r *CategoryRepositoy) GetByName(name string) (*models.Category, error) {
	zap.L().Debug("category.repo.getByName", zap.Reflect("name", name))
	var category = &models.Category{}
	if result := r.db.Where("Name=?", name).First(&category); result.Error != nil {
		return nil, result.Error
	}

	return category, nil
}

func (r *CategoryRepositoy) Update(a *models.Category) (*models.Category, error) {
	zap.L().Debug("category.repo.update", zap.Reflect("category", a))

	if result := r.db.Save(&a); result.Error != nil {
		return nil, result.Error
	}

	return a, nil
}

func (r *CategoryRepositoy) Delete(id string) error {
	zap.L().Debug("category.repo.delete", zap.Reflect("id", id))

	category, err := r.GetByID(id)
	if err != nil {
		return err
	}

	if result := r.db.Delete(&category); result.Error != nil {
		return result.Error
	}

	return nil
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
