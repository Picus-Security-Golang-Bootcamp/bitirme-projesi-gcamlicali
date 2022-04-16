package category

import (
	"errors"
	"github.com/gcamlicali/tradeshopExample/internal/api"
	httpErr "github.com/gcamlicali/tradeshopExample/internal/httpErrors"
	"github.com/gcamlicali/tradeshopExample/internal/models"
	csvRead "github.com/gcamlicali/tradeshopExample/pkg/csv"
	"gorm.io/gorm"
	"mime/multipart"
	"net/http"
)

type categoryService struct {
	repo *CategoryRepositoy
}

type Service interface {
	Create(a *models.Category) (*models.Category, error)
	GetByID(id int) (*models.Category, error)
	Update(catName string, ca *api.Category) (*models.Category, error)
	Delete(name string) error
	GetAll(pageIndex, pageSize int) (*[]models.Category, int, error)
	AddBulk(file multipart.File) error
	AddSingle(category api.Category) (*models.Category, error)
}

func NewCategoryService(repo *CategoryRepositoy) Service {
	return &categoryService{repo: repo}
}

func (c categoryService) Create(a *models.Category) (*models.Category, error) {
	NewCategory, err := c.repo.Create(a)
	if err != nil {
		return nil, httpErr.NewRestError(http.StatusInternalServerError, "Catagory create error", err.Error())
	}

	return NewCategory, nil
}

func (c categoryService) Update(catName string, ca *api.Category) (*models.Category, error) {
	//category, err := c.repo.GetByName(catName)
	//if errors.Is(err, gorm.ErrRecordNotFound) {
	//	return nil, httpErr.NewRestError(http.StatusBadRequest, "Category not found", err.Error())
	//}
	//if err != nil {
	//	return nil, httpErr.NewRestError(http.StatusInternalServerError, "Get category error", err.Error())
	//}
	//
	//category.Name = ca.Name
	//
	//updatedCategory, err := c.repo.Update(category)
	//if err != nil {
	//	return nil, httpErr.NewRestError(http.StatusInternalServerError, "Update category error", err.Error())
	//}
	//
	////Update Products with new category name
	//products, err := c.prepo.GetByCatName(catName)
	//for _, product := range *products {
	//	product.CategoryName = catName
	//	_, err := c.prepo.Update(&product)
	//	if err != nil {
	//		return nil, httpErr.NewRestError(http.StatusInternalServerError, "Update product error", err.Error())
	//	}
	//}

	//return updatedCategory, nil
	return nil, nil
}

func (c categoryService) GetByID(id int) (*models.Category, error) { //duzelt uuid olacak

	category, err := c.repo.GetByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, httpErr.NewRestError(http.StatusBadRequest, "Category not found", err.Error())
	}
	if err != nil {
		return nil, httpErr.NewRestError(http.StatusInternalServerError, "Catagory get by id error", err.Error())
	}
	return category, nil
}

func (c categoryService) Delete(name string) error {
	//category, err := c.repo.GetByName(name)
	//if errors.Is(err, gorm.ErrRecordNotFound) {
	//	return httpErr.NewRestError(http.StatusBadRequest, "Category not found", err.Error())
	//}
	//if err != nil {
	//	return httpErr.NewRestError(http.StatusInternalServerError, "Get category error", err.Error())
	//}
	//
	//err = c.repo.Delete(int(category.ID)) //duzelt uuid olacak
	//
	////Delete products include deleted category name
	//products, err := c.prepo.GetByCatName(name)
	//for _, product := range *products {
	//	err := c.prepo.Delete(&product)
	//	if err != nil {
	//		return httpErr.NewRestError(http.StatusInternalServerError, "Update product error", err.Error())
	//	}
	//}
	//return nil

	return nil
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
