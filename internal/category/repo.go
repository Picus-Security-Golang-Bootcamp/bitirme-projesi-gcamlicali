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
	if result := r.db.First(&category, id); result.Error != nil {
		return nil, result.Error
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

func (r *CategoryRepositoy) GetAll() (*[]models.Category, error) {
	zap.L().Debug("category.repo.getAll")

	var categories = &[]models.Category{}
	if result := r.db.Find(&categories); result.Error != nil {
		return nil, result.Error
	}

	return categories, nil
}

func (r *CategoryRepositoy) Migration() {
	r.db.AutoMigrate(&models.Category{})
}
