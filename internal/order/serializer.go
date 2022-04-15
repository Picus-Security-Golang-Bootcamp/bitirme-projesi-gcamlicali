package order

import (
	"github.com/gcamlicali/tradeshopExample/internal/api"
	"github.com/gcamlicali/tradeshopExample/internal/models"
)

func OrderToResponse(m *models.Order) *api.Order {
	return &api.Order{
		ID:         int64(m.ID),
		UserID:     int64(m.UserID),
		CartID:     int64(m.CartID),
		TotalPrice: m.TotalPrice,
	}
}

func ordersToResponse(ms []models.Order) []*api.Order {
	orders := make([]*api.Order, 0)

	for i, _ := range ms {

		orders = append(orders, OrderToResponse(&ms[i]))
	}

	return orders
}
