package order

import (
	"github.com/gcamlicali/tradeshopExample/internal/cart"
	"github.com/gcamlicali/tradeshopExample/internal/cart_item"
	httpErr "github.com/gcamlicali/tradeshopExample/internal/httpErrors"
	"github.com/gcamlicali/tradeshopExample/internal/models"
	"github.com/gcamlicali/tradeshopExample/internal/product"
	"log"
	"net/http"
)

type orderService struct {
	orRepo *OrderRepositoy
	cRepo  *cart.CartRepositoy
	ciRepo *cart_item.CartItemRepositoy
	pRepo  *product.ProductRepositoy
}

type Service interface {
	GetAll(userID int) (*[]models.Order, error)
	Create(userID int) (*models.Order, error)
}

func NewOrderService(orRepo *OrderRepositoy, cRepo *cart.CartRepositoy, ciRepo *cart_item.CartItemRepositoy, pRepo *product.ProductRepositoy) Service {
	return &orderService{orRepo: orRepo, cRepo: cRepo, ciRepo: ciRepo, pRepo: pRepo}
}

func (c *orderService) GetAll(userID int) (*[]models.Order, error) {

	orders, err := c.orRepo.GetByUserID(userID)
	if err != nil {
		return nil, httpErr.NewRestError(http.StatusBadRequest, err.Error(), err)
	}
	return orders, nil
}

func (c *orderService) Create(userID int) (*models.Order, error) {

	cart, err := c.cRepo.GetByUserID(userID)
	if err != nil {
		return nil, httpErr.NewRestError(http.StatusInternalServerError, "Cart error", err)
	}

	//Check cartItems quantity
	cartItems, err := c.ciRepo.GetByCartID(int(cart.ID))
	for _, cartItem := range cartItems {
		product, _ := c.pRepo.GetByID(cartItem.ProductID)
		if cartItem.Quantity > int(product.UnitStock) {
			log.Println("CartItemQuantity: ", cartItem.Quantity, " Product Quantity: ", product.UnitStock)
			return nil, httpErr.NewRestError(http.StatusBadRequest, "Not Enough Stock", cartItem.Product.Name)
		}
	}

	//Create a order of cart
	newOrder := models.Order{
		CartID: int(cart.ID),
		UserID: userID,
		Cart:   *cart,
	}
	order, err := c.orRepo.Create(&newOrder)
	if err != nil {
		return nil, httpErr.NewRestError(http.StatusInternalServerError, "Order create error", err)
	}

	//Change ordered products quantity
	for _, cartItem := range cartItems {
		product, _ := c.pRepo.GetByID(cartItem.ProductID)
		product.UnitStock -= int32(cartItem.Quantity)
		_, err = c.pRepo.Update(product)
		if err != nil {
			return nil, httpErr.NewRestError(http.StatusInternalServerError, "Ordered Product quantity update error", err)
		}
	}

	//Change current cart status after order operation
	cart.IsOrdered = true
	c.cRepo.Update(cart)

	//Create a new cart for user, current cart is ordered
	newCart := models.Cart{
		UserID: userID,
	}
	_, err = c.cRepo.Create(&newCart)
	if err != nil {
		return nil, httpErr.NewRestError(http.StatusInternalServerError, "New cart create error after cart ordered", err)
	}

	return order, nil
}
