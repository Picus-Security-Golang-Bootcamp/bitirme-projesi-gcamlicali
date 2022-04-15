package cart

import (
	"github.com/gcamlicali/tradeshopExample/internal/api"
	httpErr "github.com/gcamlicali/tradeshopExample/internal/httpErrors"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
	"github.com/spf13/cast"
	"log"
	"net/http"
	"strconv"
	"time"
)

type cartHandler struct {
	service Service
}

func NewCartHandler(r *gin.RouterGroup, service Service) {
	h := &cartHandler{service: service}

	r.GET("/", h.get)
	r.POST("/:id", h.add)
	r.PUT("/:id", h.update)
	r.DELETE("/:id", h.delete)
}

func (ch *cartHandler) get(c *gin.Context) {
	userID, isExist := c.Get("userId")
	if !isExist {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "User not found", nil)))
		return
	}
	userid := cast.ToInt(userID)

	cart, err := ch.service.Get(userid)

	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	ExpireDate := cart.CreatedAt.Add(14 * 24 * time.Hour)

	log.Println("ExpireDate: ", ExpireDate)
	log.Println(time.Now().After(ExpireDate))
	c.JSON(http.StatusOK, CartToResponse(cart))
}

func (ch *cartHandler) add(c *gin.Context) {
	userid, isExist := c.Get("userId")
	if !isExist {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "User not found", nil)))
		return
	}
	userID := cast.ToInt(userid)

	paramID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "id is not integer", err)))
	}

	cart, err := ch.service.Add(userID, paramID)

	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, CartToResponse(cart))
}

func (ch *cartHandler) update(c *gin.Context) {
	userid, isExist := c.Get("userId")
	if !isExist {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "User not found", nil)))
		return
	}

	paramID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "id is not integer", err)))
	}
	reqQuantity := api.ItemQuantity{}
	if err := c.Bind(&reqQuantity); err != nil {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "check your request body", nil)))
		return
	}
	if err := reqQuantity.Validate(strfmt.NewFormats()); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	userID := cast.ToInt(userid)
	Quantity := int(*reqQuantity.Quantity)

	cart, err := ch.service.Update(userID, paramID, Quantity)
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, CartToResponse(cart))
}

func (ch *cartHandler) delete(c *gin.Context) {
	userid, isExist := c.Get("userId")
	if !isExist {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "User not found", nil)))
		return
	}

	paramID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "id is not integer", err)))
	}

	userID := cast.ToInt(userid)

	cart, err := ch.service.Delete(userID, paramID)

	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, CartToResponse(cart))
}
