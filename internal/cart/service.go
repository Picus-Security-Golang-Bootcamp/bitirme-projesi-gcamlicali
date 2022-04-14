package cart

import (
	"github.com/gcamlicali/tradeshopExample/internal/cart_item"
	httpErr "github.com/gcamlicali/tradeshopExample/internal/httpErrors"
	"github.com/gcamlicali/tradeshopExample/internal/models"
	"github.com/gcamlicali/tradeshopExample/internal/product"
	"net/http"
)

type cartService struct {
	crepo  *CartRepositoy
	cirepo *cart_item.CartItemRepositoy
	prepo  *product.ProductRepositoy
}

type Service interface {
	Get(userID int) (*models.Cart, error)
	Add(userID int, ProductID int) (*models.Cart, error)
	Update(userID int, ProductID int, Quantity int) (*models.Cart, error)
	Delete(userID int, ProductID int) (*models.Cart, error)
}

func NewCartService(crepo *CartRepositoy, cirepo *cart_item.CartItemRepositoy, prepo *product.ProductRepositoy) Service {
	return &cartService{crepo: crepo, cirepo: cirepo, prepo: prepo}
}

func (c *cartService) Get(userID int) (*models.Cart, error) {

	cart, err := c.crepo.GetByUserID(userID)

	if err != nil {
		return nil, httpErr.NewRestError(http.StatusBadRequest, err.Error(), err)
	}

	return cart, nil
}

func (c *cartService) Add(userID int, ProductID int) (*models.Cart, error) {
	cart, err := c.crepo.GetByUserID(userID)

	if err != nil {
		return nil, httpErr.NewRestError(http.StatusInternalServerError, "Cart error", err)
	}

	product, err := c.prepo.GetByID(ProductID)
	if err != nil {
		return nil, httpErr.NewRestError(http.StatusBadRequest, "Product not found", err)
	}

	cartItem, err := c.cirepo.GetByCartAndProductID(int(cart.ID), ProductID)

	// If item exists in cart, increase item quantity by 1
	if cartItem != nil {

		cartItem.Quantity = cartItem.Quantity + 1
		_, err = c.cirepo.Update(cartItem)
		if err != nil {
			return nil, httpErr.NewRestError(http.StatusInternalServerError, "Cart Item update error", err)
		}

		newCart, _ := c.crepo.GetByUserID(userID)

		return newCart, nil

	} else {
		// If item does not exist in cart, create new item

		newCartItem := models.CartItem{
			Quantity:  1,
			Price:     int(product.Price),
			ProductID: int(product.ID),
			Product:   *product,
		}

		addItem, err := c.cirepo.Crate(&newCartItem)
		if err != nil {
			return nil, httpErr.NewRestError(http.StatusInternalServerError, "Cart Item crate error", err)
		}

		cart.CartItems = append(cart.CartItems, *addItem)

		newCart, err := c.crepo.Update(cart)
		if err != nil {
			return nil, httpErr.NewRestError(http.StatusInternalServerError, "Cart Item update error", err)
		}

		return newCart, nil
	}
}

func (c *cartService) Update(userID int, ProductID int, Quantity int) (*models.Cart, error) {
	cart, err := c.crepo.GetByUserID(userID)

	if err != nil {
		return nil, httpErr.NewRestError(http.StatusInternalServerError, "Cart error", err)
	}

	cartItem, err := c.cirepo.GetByCartAndProductID(int(cart.ID), ProductID)
	if err != nil {
		return nil, httpErr.NewRestError(http.StatusBadRequest, "Product not found", nil)
	}
	// Duzelt Quantity control
	cartItem.Quantity = Quantity
	_, err = c.cirepo.Update(cartItem)
	if err != nil {
		return nil, httpErr.NewRestError(http.StatusInternalServerError, "Cart Item update error", err)
	}

	newCart, _ := c.crepo.GetByUserID(userID)

	return newCart, nil
}

func (c *cartService) Delete(userID int, ProductID int) (*models.Cart, error) {
	cart, err := c.crepo.GetByUserID(userID)

	if err != nil {
		return nil, httpErr.NewRestError(http.StatusInternalServerError, "Cart error", err)
	}

	cartItem, err := c.cirepo.GetByCartAndProductID(int(cart.ID), ProductID)
	if err != nil {
		return nil, httpErr.NewRestError(http.StatusBadRequest, "Product not found", nil)
	}

	err = c.cirepo.Delete(cartItem)
	if err != nil {
		return nil, httpErr.NewRestError(http.StatusInternalServerError, "Cart item Delete error", err)
	}

	newCart, _ := c.crepo.GetByUserID(userID)

	return newCart, nil
}
