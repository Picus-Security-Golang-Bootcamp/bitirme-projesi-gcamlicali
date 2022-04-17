package category

import (
	"github.com/gcamlicali/tradeshopExample/internal/api"
	httpErr "github.com/gcamlicali/tradeshopExample/internal/httpErrors"
	"github.com/gcamlicali/tradeshopExample/internal/models"
	csvRead "github.com/gcamlicali/tradeshopExample/pkg/csv"
	"mime/multipart"
	"net/http"
)

type categoryService struct {
	repo ICategoryRepository
}

type Service interface {
	Create(a *models.Category) (*models.Category, error)
	GetAll(pageIndex, pageSize int) (*[]models.Category, int, error)
	AddBulk(file multipart.File) error
	AddSingle(category api.Category) (*models.Category, error)
}

func NewCategoryService(repo ICategoryRepository) Service {
	return &categoryService{repo: repo}
}

func (c categoryService) Create(a *models.Category) (*models.Category, error) {
	NewCategory, err := c.repo.Create(a)
	if err != nil {
		return nil, httpErr.NewRestError(http.StatusInternalServerError, "Catagory create error", err.Error())
	}

	return NewCategory, nil
}

func (c categoryService) GetAll(pageIndex, pageSize int) (*[]models.Category, int, error) {

	categories, count, err := c.repo.GetAll(pageIndex, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return categories, count, nil
}

func (c categoryService) AddBulk(file multipart.File) error {

	record, err := csvRead.ReadFile(file)
	if err != nil {
		return httpErr.NewRestError(http.StatusInternalServerError, "Can not read csv file", err.Error())
	}

	for _, line := range record {
		catEntity := models.Category{}
		catEntity.Name = &line[0]
		_, err = c.Create(&catEntity)
		if err != nil {
			return httpErr.NewRestError(http.StatusBadRequest, "Category create error", err.Error())
		}
	}

	return nil
}

func (c categoryService) AddSingle(category api.Category) (*models.Category, error) {

	dbCat := models.Category{}
	dbCat.Name = category.Name

	createdCategory, err := c.Create(&dbCat)
	if err != nil {
		return nil, httpErr.NewRestError(http.StatusBadRequest, "Category create error", err.Error())
	}

	return createdCategory, nil
}
