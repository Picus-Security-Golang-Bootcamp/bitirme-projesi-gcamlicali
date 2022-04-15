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
	//log.Println("CategoryName: ", &p.CategoryName)
	//log.Println("ID: ", int64(p.ID))
	//log.Println("Sku: ", &p.SKU)
	//log.Println("Desc: ", p.Description)
	int64Sku := int64(p.SKU)
	return &api.Product{

		CategoryName: &p.CategoryName,
		Sku:          &int64Sku,
		Name:         &p.Name,
		Description:  p.Description,
		Price:        &p.Price,
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
		Model:        gorm.Model{},
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
