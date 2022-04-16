package cart

import (
	"github.com/gcamlicali/tradeshopExample/internal/cart_item"
	httpErr "github.com/gcamlicali/tradeshopExample/internal/httpErrors"
	"github.com/gcamlicali/tradeshopExample/internal/models"
	"github.com/gcamlicali/tradeshopExample/internal/product"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
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

	if err != nil {
		if err == gorm.ErrRecordNotFound {

		} else {
			return nil, httpErr.NewRestError(http.StatusInternalServerError, "Cart get error", err.Error())
		}
	}
	if cartItem != nil {
		log.Println("cartitem bulundu: ", cartItem.ProductSKU)
	}

	// If item exists in cart, increase item quantity by 1
	if cartItem != nil {

		cartItem.Quantity = cartItem.Quantity + 1
		cartItem.Price = cartItem.Quantity * product.Price
		_, err = c.cirepo.Update(cartItem)
		if err != nil {
			return nil, httpErr.NewRestError(http.StatusInternalServerError, "Cart Item update error", err.Error())
		}

		cart.TotalPrice = c.calculateCartPrice(cart)
		log.Println("cartTotalPrice: ", cart.TotalPrice)

		_, err := c.crepo.Update(cart)
		if err != nil {
			return nil, httpErr.NewRestError(http.StatusInternalServerError, "Cart update error", err.Error())
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
			Price:      product.Price,
			ProductSKU: product.SKU,
			Product:    *product,
		}

		addItem, err := c.cirepo.Crate(&newCartItem)
		if err != nil {
			return nil, httpErr.NewRestError(http.StatusInternalServerError, "Cart Item crate error", err.Error())
		}

		cart.CartItems = append(cart.CartItems, *addItem)

		cart.TotalPrice = c.calculateCartPrice(cart) + newCartItem.Price
		log.Println("cartTotalPrice: ", cart.TotalPrice)
		newCart, err := c.crepo.Update(cart)
		if err != nil {
			return nil, httpErr.NewRestError(http.StatusInternalServerError, "Cart update error", err.Error())
		}

		return newCart, nil
	}
}

//Update quantity of given cart item
func (c *cartService) Update(userID uuid.UUID, ProductSKU int, Quantity int) (*models.Cart, error) {
	// Get user cart
	cart, err := c.crepo.GetByUserID(userID)
	if err != nil {
		return nil, httpErr.NewRestError(http.StatusInternalServerError, "Get Cart error", err.Error())
	}

	//Get cart_item by SKU in cart
	cartItem, err := c.cirepo.GetByCartAndProductSKU(cart.ID, ProductSKU)
	if err != nil {
		return nil, httpErr.NewRestError(http.StatusBadRequest, "Product not found", err.Error())
	}

	productPrice := cartItem.Price / cartItem.Quantity

	// Duzelt Quantity control
	cartItem.Quantity = Quantity
	cartItem.Price = Quantity * productPrice
	_, err = c.cirepo.Update(cartItem)
	if err != nil {
		return nil, httpErr.NewRestError(http.StatusInternalServerError, "Cart Item update error", err.Error())
	}

	// Update Cart Total Price
	cart.TotalPrice = c.calculateCartPrice(cart)
	_, err = c.crepo.Update(cart)
	if err != nil {
		return nil, httpErr.NewRestError(http.StatusInternalServerError, "Cart update error", err.Error())
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

	// Update Cart Total Price
	cart.TotalPrice = c.calculateCartPrice(cart)
	_, err = c.crepo.Update(cart)
	if err != nil {
		return nil, httpErr.NewRestError(http.StatusInternalServerError, "Cart update error", err.Error())
	}

	newCart, err := c.crepo.GetByUserID(userID) // duzelt kontrol edilecek bu error
	if err != nil {
		return nil, httpErr.NewRestError(http.StatusInternalServerError, "Get Cart error", err.Error())
	}

	return newCart, nil
}

func (c *cartService) calculateCartPrice(cart *models.Cart) int {
	cartItems, _ := c.cirepo.GetByCartID(cart.ID)

	var totalPrice int

	for _, cartItem := range cartItems {
		log.Println("CartItemPrice: ", cartItem.Price)
		totalPrice += cartItem.Price
	}
	log.Println("TotalPrice: ", totalPrice)
	return totalPrice
}
