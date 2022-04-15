package product

import (
	"errors"
	"github.com/gcamlicali/tradeshopExample/internal/api"
	"github.com/gcamlicali/tradeshopExample/internal/category"
	httpErr "github.com/gcamlicali/tradeshopExample/internal/httpErrors"
	"github.com/gcamlicali/tradeshopExample/internal/models"
	csvRead "github.com/gcamlicali/tradeshopExample/pkg/csv"
	"gorm.io/gorm"
	"mime/multipart"
	"net/http"
	"strconv"
)

type productService struct {
	pRepo   *ProductRepositoy
	catRepo *category.CategoryRepositoy
}

type Service interface {
	AddBulk(file multipart.File) error
	AddSingle(product api.Product) (*models.Product, error)
	GetAll(pageIndex, pageSize int) (*[]models.Product, int, error)
	Delete(SKU int) error
	Update(SKU int, reqProduct *api.ProductUp) (*models.Product, error)
}

func NewProductService(pRepo *ProductRepositoy, catRepo *category.CategoryRepositoy) Service {
	return &productService{pRepo: pRepo, catRepo: catRepo}
}

func (p *productService) AddBulk(file multipart.File) error {
	record, err := csvRead.ReadFile(file)

	if err != nil {
		return httpErr.NewRestError(http.StatusInternalServerError, "Can not read csv file", err.Error())
	}

	// Duzelt
	for _, line := range record {
		proEntity := models.Product{}
		proEntity.CategoryName = line[0]
		_, err := p.catRepo.GetByName(proEntity.CategoryName)
		if err != nil {
			//c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusNotFound, "Category not found", proEntity.CategoryName)))
			continue
		}
		proEntity.Name = line[1]
		SKU, err := strconv.Atoi(line[2])
		if err != nil {
			//c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "SKU is not integer", proEntity.Name)))
			continue
		}
		proEntity.SKU = SKU
		proEntity.Description = line[3]
		price, err := strconv.Atoi(line[4])
		if err != nil {
			//c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "Price is not integer", proEntity.Name)))
			continue
		}
		proEntity.Price = int32(price)
		unitStock, err := strconv.Atoi(line[5])
		if err != nil {
			//c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "UnitStock is not integer", proEntity.Name)))
			continue
		}
		proEntity.UnitStock = int32(unitStock)

		_, err = p.pRepo.create(&proEntity)
		if err != nil {
			//c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, err.Error(), nil)))
			continue
		}
	}

	return nil
}

func (c productService) AddSingle(product api.Product) (*models.Product, error) {

	category, err := c.catRepo.GetByName(*product.Name)
	if err != nil {
		return nil, httpErr.NewRestError(http.StatusNotFound, "Category not found", err.Error())
	}

	prod := responseToProduct(&product)
	prod.CategoryName = *category.Name
	//log.Println("prod cat id: ", prod.CategoryID)
	NewProduct, err := c.pRepo.create(prod)
	if err != nil {
		return nil, httpErr.NewRestError(http.StatusInternalServerError, "Can not create new product", err.Error())
	}

	return NewProduct, nil
}

func (c productService) GetAll(pageIndex, pageSize int) (*[]models.Product, int, error) {

	categories, count, err := c.pRepo.getAll(pageIndex, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return categories, count, nil
}

func (c productService) Delete(SKU int) error {
	err := c.pRepo.deleteBySku(SKU)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return httpErr.NewRestError(http.StatusBadRequest, "Product not found", err.Error())
	}

	if err != nil {
		return httpErr.NewRestError(http.StatusInternalServerError, "Delete product error", err.Error())
	}

	return nil
}

func (c productService) Update(SKU int, reqProduct *api.ProductUp) (*models.Product, error) {
	product, err := c.pRepo.GetBySKU(SKU)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, httpErr.NewRestError(http.StatusBadRequest, "Product not found", err.Error())
	}
	if err != nil {
		return nil, httpErr.NewRestError(http.StatusInternalServerError, "Get product error", err.Error())
	}

	if reqProduct.Name != "" {
		product.Name = reqProduct.Name
	}
	if reqProduct.CategoryName != "" {
		//check category name of product
		cat, _ := c.catRepo.GetByName(reqProduct.CategoryName)
		if cat == nil {
			return nil, httpErr.NewRestError(http.StatusBadRequest, "Product category name not found", err.Error())
		}

		product.CategoryName = reqProduct.CategoryName
	}

	if reqProduct.Description != "" {
		product.Description = reqProduct.Description
	}
	if reqProduct.Price != 0 {
		product.Price = reqProduct.Price
	}
	if reqProduct.Sku != 0 {
		product.SKU = int(reqProduct.Sku)
	}
	if reqProduct.UnitStock != 0 {
		product.UnitStock = reqProduct.UnitStock
	}

	updatedProduct, err := c.pRepo.Update(product)
	if err != nil {
		return nil, httpErr.NewRestError(http.StatusInternalServerError, "Update product error", err.Error())
	}

	return updatedProduct, nil

}
