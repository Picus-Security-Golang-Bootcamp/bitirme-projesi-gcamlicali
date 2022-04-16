package cart

import (
	"github.com/gcamlicali/tradeshopExample/internal/cart_item"
	httpErr "github.com/gcamlicali/tradeshopExample/internal/httpErrors"
	"github.com/gcamlicali/tradeshopExample/internal/models"
	"github.com/gcamlicali/tradeshopExample/internal/product"
	"github.com/google/uuid"
	"net/http"
)

type cartService struct {
	crepo  *CartRepositoy
	cirepo *cart_item.CartItemRepositoy
	prepo  *product.ProductRepositoy
}

type Service interface {
	Get(userID uuid.UUID) (*models.Cart, error)
	Add(userID uuid.UUID, ProductID int) (*models.Cart, error)
	Update(userID uuid.UUID, ProductID int, Quantity int) (*models.Cart, error)
	Delete(userID uuid.UUID, ProductID int) (*models.Cart, error)
}

func NewCartService(crepo *CartRepositoy, cirepo *cart_item.CartItemRepositoy, prepo *product.ProductRepositoy) Service {
	return &cartService{crepo: crepo, cirepo: cirepo, prepo: prepo}
}

//Get all items from cart and list
func (c *cartService) Get(userID uuid.UUID) (*models.Cart, error) {

	cart, err := c.crepo.GetByUserID(userID)
	if err != nil {
		return nil, httpErr.NewRestError(http.StatusInternalServerError, "Cart get error", err.Error())
	}

	return cart, nil
}

//Add item to cart
func (c *cartService) Add(userID uuid.UUID, ProductSKU int) (*models.Cart, error) {

	cart, err := c.crepo.GetByUserID(userID) // duzelt kontrol edilecek bu error
	if err != nil {
		return nil, httpErr.NewRestError(http.StatusInternalServerError, "Cart get error", err.Error())
	}

	product, err := c.prepo.GetBySKU(ProductSKU) // duzelt kontrol edilecek bu error
	if err != nil {
		return nil, httpErr.NewRestError(http.StatusBadRequest, "Product not found", err.Error())
	}

	cartItem, err := c.cirepo.GetByCartAndProductSKU(cart.ID, ProductSKU)

	// If item exists in cart, increase item quantity by 1
	if cartItem != nil {

		cartItem.Quantity = cartItem.Quantity + 1
		_, err = c.cirepo.Update(cartItem)
		if err != nil {
			return nil, httpErr.NewRestError(http.StatusInternalServerError, "Cart Item update error", err.Error())
		}

		newCart, _ := c.crepo.GetByUserID(userID)
		if err != nil {
			return nil, httpErr.NewRestError(http.StatusInternalServerError, "Get Cart error", err.Error())
		}

		return newCart, nil

	} else {
		// If item does not exist in cart, create new item

		newCartItem := models.CartItem{
			Quantity:   1,
			Price:      int(product.Price),
			ProductSKU: product.SKU,
			Product:    *product,
		}

		addItem, err := c.cirepo.Crate(&newCartItem)
		if err != nil {
			return nil, httpErr.NewRestError(http.StatusInternalServerError, "Cart Item crate error", err.Error())
		}

		cart.CartItems = append(cart.CartItems, *addItem)

		newCart, err := c.crepo.Update(cart)
		if err != nil {
			return nil, httpErr.NewRestError(http.StatusInternalServerError, "Cart Item update error", err.Error())
		}

		return newCart, nil
	}
}

//Update quantity of given cart item
func (c *cartService) Update(userID uuid.UUID, ProductSKU int, Quantity int) (*models.Cart, error) {

	cart, err := c.crepo.GetByUserID(userID)
	if err != nil {
		return nil, httpErr.NewRestError(http.StatusInternalServerError, "Get Cart error", err.Error())
	}

	if err != nil {
		return nil, httpErr.NewRestError(http.StatusInternalServerError, "Cart error", err.Error())
	}

	cartItem, err := c.cirepo.GetByCartAndProductSKU(cart.ID, ProductSKU)
	if err != nil {
		return nil, httpErr.NewRestError(http.StatusBadRequest, "Product not found", err.Error())
	}
	// Duzelt Quantity control
	cartItem.Quantity = Quantity
	_, err = c.cirepo.Update(cartItem)
	if err != nil {
		return nil, httpErr.NewRestError(http.StatusInternalServerError, "Cart Item update error", err.Error())
	}

	newCart, err := c.crepo.GetByUserID(userID)
	if err != nil {
		return nil, httpErr.NewRestError(http.StatusInternalServerError, "Get Cart error", err.Error())
	}

	return newCart, nil
}

//Delete given item from cart
func (c *cartService) Delete(userID uuid.UUID, ProductSKU int) (*models.Cart, error) {
	cart, err := c.crepo.GetByUserID(userID)
	if err != nil {
		return nil, httpErr.NewRestError(http.StatusInternalServerError, "Get Cart error", err.Error())
	}

	cartItem, err := c.cirepo.GetByCartAndProductSKU(cart.ID, ProductSKU)
	if err != nil {
		return nil, httpErr.NewRestError(http.StatusBadRequest, "Product not found", err.Error())
	}

	err = c.cirepo.Delete(cartItem)
	if err != nil {
		return nil, httpErr.NewRestError(http.StatusInternalServerError, "Cart item Delete error", err.Error())
	}

	newCart, err := c.crepo.GetByUserID(userID) // duzelt kontrol edilecek bu error
	if err != nil {
		return nil, httpErr.NewRestError(http.StatusInternalServerError, "Get Cart error", err.Error())
	}

	return newCart, nil
}
