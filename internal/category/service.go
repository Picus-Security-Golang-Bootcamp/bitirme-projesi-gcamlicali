package category

import (
	"errors"
	"github.com/gcamlicali/tradeshopExample/internal/api"
	httpErr "github.com/gcamlicali/tradeshopExample/internal/httpErrors"
	"github.com/gcamlicali/tradeshopExample/internal/models"
	csvRead "github.com/gcamlicali/tradeshopExample/pkg/csv"
	"log"
	"mime/multipart"
	"net/http"
)

type categoryService struct {
	repo CategoryRepositoy
}

type Service interface {
	Create(a *models.Category) (*models.Category, error)
	GetByID(id string) (*models.Category, error)
	Update(a *models.Category) (*models.Category, error)
	Delete(id string) error
	GetAll(pageIndex, pageSize int) (*[]models.Category, int, error)
	AddBulk(file multipart.File) error
	AddSingle(category api.Category) (*models.Category, error)
}

func NewCategoryService(repo CategoryRepositoy) Service {
	//if repo == nil {
	//	return nil
	//}

	return &categoryService{repo: repo}
}

func (c categoryService) Create(a *models.Category) (*models.Category, error) {
	NewCategory, err := c.repo.Create(a)

	if err != nil {
		return nil, err
	}

	return NewCategory, nil
}

func (c categoryService) Update(ca *models.Category) (*models.Category, error) {
	//TODO implement me
	panic("implement me")
}

func (c categoryService) GetByID(id string) (*models.Category, error) {
	if len(id) == 0 {
		return nil, errors.New("Id cannot be nil or empty")
	}

	category, _ := c.repo.GetByID(id)
	return category, nil
}

func (c categoryService) Delete(id string) error {
	//TODO implement me
	panic("implement me")
}

func (c categoryService) GetAll(pageIndex, pageSize int) (*[]models.Category, int, error) {

	//categories := &[]models.Category{}
	categories, count, err := c.repo.GetAll(pageIndex, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return categories, count, nil
}

func (c categoryService) AddBulk(file multipart.File) error {

	record, err := csvRead.ReadFile(file)
	if err != nil {
		return httpErr.NewRestError(http.StatusInternalServerError, "Can not read csv file", nil)
	}

	for _, line := range record {
		catEntity := models.Category{}
		catEntity.Name = &line[0]
		_, err = c.Create(&catEntity)
		if err != nil {
			return httpErr.NewRestError(http.StatusBadRequest, err.Error(), nil)
		}
	}

	return nil

}

func (c categoryService) AddSingle(category api.Category) (*models.Category, error) {
	log.Println("Srvice Add single geldi")
	dbCat := models.Category{}
	dbCat.Name = category.Name
	createdCategory, err := c.Create(&dbCat)
	log.Println("Database service single geldi")
	if err != nil {
		log.Println("Database service single hata var")
		return nil, httpErr.NewRestError(http.StatusBadRequest, err.Error(), nil)
	}

	return createdCategory, nil
}
