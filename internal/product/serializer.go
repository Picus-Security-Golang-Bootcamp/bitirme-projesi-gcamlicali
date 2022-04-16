package product

import (
	"github.com/gcamlicali/tradeshopExample/internal/api"
	"github.com/gcamlicali/tradeshopExample/internal/models"
)

//Data Transfer Object
func ProductToResponse(p *models.Product) *api.Product {
	int64Sku := int64(p.SKU)
	int32Price := int32(p.Price)
	return &api.Product{

		CategoryName: &p.CategoryName,
		Sku:          &int64Sku,
		Name:         &p.Name,
		Description:  p.Description,
		Price:        &int32Price,
		UnitStock:    &p.UnitStock,
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
		CategoryName: *p.CategoryName,
		Name:         *p.Name,
		SKU:          int(*p.Sku),
		Description:  p.Description,
		UnitStock:    *p.UnitStock,
	}
}

func responseToProductUp(p *api.ProductUp) *models.Product {
	return &models.Product{
		CategoryName: p.CategoryName,
		Name:         p.Name,
		SKU:          int(p.Sku),
		Description:  p.Description,
		UnitStock:    p.UnitStock,
	}
}
