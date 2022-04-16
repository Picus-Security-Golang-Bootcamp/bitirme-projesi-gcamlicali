package order

import (
	"github.com/gcamlicali/tradeshopExample/internal/api"
	"github.com/gcamlicali/tradeshopExample/internal/models"
)

func OrderToResponse(m *models.Order) *api.Order {
	return &api.Order{
		ID:         m.ID.String(),
		UserID:     m.UserID.String(),
		CartID:     m.CartID.String(),
		Status:     m.Status,
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
