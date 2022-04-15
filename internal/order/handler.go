package order

import (
	httpErr "github.com/gcamlicali/tradeshopExample/internal/httpErrors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/http"
	"strconv"
)

type orderHandler struct {
	service Service
}

func NewOrderHandler(r *gin.RouterGroup, service Service) {
	h := &orderHandler{service: service}
	r.GET("/", h.getAll)
	r.POST("/", h.add)
	r.PUT("/:id", h.cancel)
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

func (o *orderHandler) cancel(c *gin.Context) {
	userID, isExist := c.Get("userId")
	if !isExist {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "User not found", nil)))
		return
	}
	userid := cast.ToInt(userID)

	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "id is not integer", err)))
	}

	err = o.service.Cancel(userid, orderID)

	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, "Order Cancel Complete")
}
