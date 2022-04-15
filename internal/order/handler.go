package order

import (
	httpErr "github.com/gcamlicali/tradeshopExample/internal/httpErrors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/http"
)

type orderHandler struct {
	service Service
}

func NewOrderHandler(r *gin.RouterGroup, service Service) {
	h := &orderHandler{service: service}
	r.GET("/", h.getAll)
	r.POST("", h.add)
}

func (o *orderHandler) getAll(c *gin.Context) {
	userID, isExist := c.Get("userId")
	if !isExist {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "User not found", nil)))
		return
	}
	userid := cast.ToInt(userID)

	orders, err := o.service.GetAll(userid)
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}
	
	c.JSON(http.StatusOK, ordersToResponse(*orders))

}

func (o *orderHandler) add(c *gin.Context) {
	userID, isExist := c.Get("userId")
	if !isExist {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "User not found", nil)))
		return
	}
	userid := cast.ToInt(userID)

	order, err := o.service.Create(userid)

	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, OrderToResponse(order))

}
