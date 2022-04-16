package cart

import (
	"github.com/gcamlicali/tradeshopExample/internal/api"
	"github.com/gcamlicali/tradeshopExample/internal/cart_item"
	"github.com/gcamlicali/tradeshopExample/internal/models"
)

func CartToResponse(a *models.Cart) *api.Cart {
	items := make([]*api.CartItem, 0)
	for i, _ := range a.CartItems {

		items = append(items, cart_item.CartItemtoResponse(&a.CartItems[i]))
	}

	return &api.Cart{
		ID:        a.ID.String(),
		CartItems: items,
	}
}
