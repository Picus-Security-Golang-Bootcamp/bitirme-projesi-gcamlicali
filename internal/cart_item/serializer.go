package cart_item

import (
	"github.com/gcamlicali/tradeshopExample/internal/api"
	"github.com/gcamlicali/tradeshopExample/internal/models"
	"github.com/gcamlicali/tradeshopExample/internal/product"
)

func CartItemtoResponse(ci *models.CartItem) *api.CartItem {

	product := product.ProductToResponse(&ci.Product)

	return &api.CartItem{
		Product:  product,
		Quantity: int32(ci.Quantity),
		Price:    int32(ci.Price),
	}
}
