package cart

import (
	"github.com/gcamlicali/tradeshopExample/internal/cart_item"
	httpErr "github.com/gcamlicali/tradeshopExample/internal/httpErrors"
	"github.com/gcamlicali/tradeshopExample/internal/models"
	"github.com/gcamlicali/tradeshopExample/internal/product"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"log"
	"net/http"
)

type cartHandler struct {
	cirepo *cart_item.CartItemRepositoy
	crepo  *CartRepositoy
	prepo  *product.ProductRepositoy
}

func NewCartHandler(r *gin.RouterGroup, crepo *CartRepositoy, cirepo *cart_item.CartItemRepositoy, prepo *product.ProductRepositoy) {
	h := &cartHandler{crepo: crepo, cirepo: cirepo, prepo: prepo}

	r.GET("/", h.get)
	r.POST("/:id", h.addItem)
	//r.GET("/:id", h.getByID)
	//r.PUT("/:id", h.update)
	//r.DELETE("/:id", h.delete)
}

func (ch *cartHandler) get(c *gin.Context) {
	userID, isExist := c.Get("userId")
	if !isExist {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "User not found", nil)))
		return
	}
	userid := cast.ToInt(userID)
	cart, err := ch.crepo.GetByUserID(userid)

	if err != nil {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, err.Error(), nil)))
		return
	}

	c.JSON(http.StatusOK, CartToResponse(cart))
}

func (ch *cartHandler) addItem(c *gin.Context) {
	log.Println("add item a geldi")
	userID, isExist := c.Get("userId")
	if !isExist {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "User not found", nil)))
		return
	}
	userid := cast.ToInt(userID)
	cart, err := ch.crepo.GetByUserID(userid)

	if err != nil {
		log.Println("Cart cekmede hata")
		return
	}
	log.Println("cart cekildi")
	//cartItems, err := ch.cirepo.GetByCartID(int(cart.ID))
	paramid := c.Param("id")
	log.Println("ParamID okundu : ", paramid)
	product, err := ch.prepo.GetByID(paramid)
	if err != nil {
		log.Println("Product cekilemedi")
		return
	}
	log.Println("Product geldi adi: ", product.Name)
	cartItem := models.CartItem{}
	cartItem.Quantity = 1
	cartItem.Price = int(product.Price)
	cartItem.ProductID = int(product.ID)
	cartItem.Product = *product

	additem, err := ch.cirepo.Crate(&cartItem)
	if err != nil {
		log.Println("Cart item yaratilamadi")
		return
	}
	log.Println("Cart Item yaratildi productid: ", additem.ProductID)
	cart.CartItems = append(cart.CartItems, *additem)

	newCart, err := ch.crepo.Update(cart)
	if err != nil {
		log.Println("Cart update hata")
		return
	}
	log.Println("Cart olustu", newCart)
}
