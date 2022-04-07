package product

import (
	"github.com/gcamlicali/tradeshopExample/internal/api"
	"github.com/gcamlicali/tradeshopExample/internal/models"
	"gorm.io/gorm"
)

//Data Transfer Object
func ProductToResponse(p *models.Product) *api.Product {
	//code
	//code

	return &api.Product{
		Category:    nil,
		ID:          int64(p.ID),
		Sku:         &p.SKU,
		Name:        &p.Name,
		Description: p.Description,
	}
}

// return Objects
func productsToResponse(ps *[]models.Product) []*api.Product {
	products := make([]*api.Product, 0)
	for _, p := range *ps {
		products = append(products, ProductToResponse(&p))
	}
	return products
}

func responseToProduct(p *api.Product) *models.Product {
	return &models.Product{
		Model:       gorm.Model{ID: uint(p.ID)},
		CategoryID:  p.Category.ID,
		Name:        *p.Name,
		SKU:         *p.Sku,
		Description: p.Description,
	}
}
