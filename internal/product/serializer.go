package product

import (
	"github.com/gcamlicali/tradeshopExample/internal/api"
	"github.com/gcamlicali/tradeshopExample/internal/models"
	"gorm.io/gorm"
	"log"
)

//Data Transfer Object
func ProductToResponse(p *models.Product) *api.Product {
	//code
	//code
	log.Println("CategoryName: ", &p.CategoryName)
	log.Println("ID: ", int64(p.ID))
	log.Println("Sku: ", &p.SKU)
	log.Println("Desc: ", p.Description)
	return &api.Product{

		CategoryName: &p.CategoryName,
		ID:           int64(p.ID),
		Sku:          &p.SKU,
		Name:         &p.Name,
		Description:  p.Description,
	}
}

// return Objects
func productsToResponse(ps []models.Product) []*api.Product {
	products := make([]*api.Product, 0)

	for i, _ := range ps {

		products = append(products, ProductToResponse(&ps[i]))
	}

	return products
}

func responseToProduct(p *api.Product) *models.Product {
	return &models.Product{
		Model:        gorm.Model{ID: uint(p.ID)},
		CategoryName: *p.CategoryName,
		Name:         *p.Name,
		SKU:          *p.Sku,
		Description:  p.Description,
	}
}
