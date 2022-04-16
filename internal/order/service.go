package order

import (
	"database/sql"
	"errors"
	"github.com/gcamlicali/tradeshopExample/internal/cart"
	"github.com/gcamlicali/tradeshopExample/internal/cart_item"
	httpErr "github.com/gcamlicali/tradeshopExample/internal/httpErrors"
	"github.com/gcamlicali/tradeshopExample/internal/models"
	"github.com/gcamlicali/tradeshopExample/internal/product"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

var (
	OneDay    = 24
	ExpireDay = 14
)

type orderService struct {
	orRepo *OrderRepositoy
	cRepo  *cart.CartRepositoy
	ciRepo *cart_item.CartItemRepositoy
	pRepo  *product.ProductRepositoy
}

type Service interface {
	GetAll(userID uuid.UUID) (*[]models.Order, error)
	Create(userID uuid.UUID) (*models.Order, error)
	Cancel(userID uuid.UUID, orderID uuid.UUID) error
}

func NewOrderService(orRepo *OrderRepositoy, cRepo *cart.CartRepositoy, ciRepo *cart_item.CartItemRepositoy, pRepo *product.ProductRepositoy) Service {
	return &orderService{orRepo: orRepo, cRepo: cRepo, ciRepo: ciRepo, pRepo: pRepo}
}

func (c *orderService) GetAll(userID uuid.UUID) (*[]models.Order, error) {

	orders, err := c.orRepo.GetByUserID(userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, httpErr.NewRestError(http.StatusBadRequest, "Order not found", err.Error())
	}
	if err != nil {
		return nil, httpErr.NewRestError(http.StatusInternalServerError, "Can't get orders", err.Error())
	}
	return orders, nil
}

func (c *orderService) Create(userID uuid.UUID) (*models.Order, error) {

	cart, err := c.cRepo.GetByUserID(userID)
	if err != nil {
		return nil, httpErr.NewRestError(http.StatusInternalServerError, "Cart error", err.Error())
	}

	//Check cartItems quantity
	cartItems, err := c.ciRepo.GetByCartID(cart.ID)
	for _, cartItem := range cartItems {
		product, _ := c.pRepo.GetBySKU(cartItem.ProductSKU)
		if cartItem.Quantity > int(product.UnitStock) {
			log.Println("CartItemQuantity: ", cartItem.Quantity, " Product Quantity: ", product.UnitStock)
			return nil, httpErr.NewRestError(http.StatusBadRequest, "Not Enough Stock", cartItem.Product.Name)
		}
	}

	//Create a order of cart
	newOrder := models.Order{
		CartID:     cart.ID,
		UserID:     userID,
		Cart:       *cart,
		Status:     "Ordered",
		TotalPrice: 33,
	}
	order, err := c.orRepo.Create(&newOrder)
	if err != nil {
		return nil, httpErr.NewRestError(http.StatusInternalServerError, "Order create error", err.Error())
	}

	//Change ordered products quantity
	for _, cartItem := range cartItems {
		product, _ := c.pRepo.GetBySKU(cartItem.ProductSKU)
		product.UnitStock -= int32(cartItem.Quantity)
		_, err = c.pRepo.Update(product)
		if err != nil {
			return nil, httpErr.NewRestError(http.StatusInternalServerError, "Ordered Product quantity update error", err.Error())
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
		return nil, httpErr.NewRestError(http.StatusInternalServerError, "New cart create error after cart ordered", err.Error())
	}

	return order, nil
}

func (c *orderService) Cancel(userID uuid.UUID, orderID uuid.UUID) error {

	//Get given order by user and order ID
	order, err := c.orRepo.GetByOrderAndUserID(userID, orderID)
	if errors.Is(err, sql.ErrNoRows) {
		return httpErr.NewRestError(http.StatusNotFound, "Order not found", err.Error())
	}

	if err != nil {
		return httpErr.NewRestError(http.StatusInternalServerError, "Get order error", err.Error())
	}

	//Check order expire date
	orderExpireDate := order.CreatedAt.Add(time.Duration(ExpireDay*OneDay) * time.Hour)
	now := time.Now()

	if now.After(orderExpireDate) {
		return httpErr.NewRestError(http.StatusBadRequest, "Order cancel date expired", "You can not cancel your order!")
	}

	//Set order status cancelled
	order.Status = "Cancelled"
	_, err = c.orRepo.Update(order)
	if err != nil {
		return httpErr.NewRestError(http.StatusInternalServerError, "Order Update Error", err.Error())
	}

	//Give ordered product quantity back
	cartItems, err := c.ciRepo.GetByCartID(order.CartID)
	if err != nil {
		return httpErr.NewRestError(http.StatusInternalServerError, "Get cart items Error", err.Error())
	}

	for _, cartItem := range cartItems {
		product, _ := c.pRepo.GetBySKU(cartItem.ProductSKU)
		product.UnitStock += int32(cartItem.Quantity)
		_, err = c.pRepo.Update(product)
		if err != nil {
			return httpErr.NewRestError(http.StatusInternalServerError, "Ordered Product quantity update error", err.Error())
		}
	}

	return nil
}
