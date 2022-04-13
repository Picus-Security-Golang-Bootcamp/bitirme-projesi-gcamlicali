package cart_item

import (
	"github.com/gcamlicali/tradeshopExample/internal/api"
	"github.com/gcamlicali/tradeshopExample/internal/models"
)

func CartItemtoResponse(ci *models.CartItem) *api.CartItem {
	return &api.CartItem{
		CartID:    int64(ci.CartID),
		ID:        int64(ci.ID),
		ProductID: int64(ci.ProductID),
		Quantity:  int32(ci.Quantity),
	}
}
