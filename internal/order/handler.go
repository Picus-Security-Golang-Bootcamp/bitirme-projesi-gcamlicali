package order

import (
	httpErr "github.com/gcamlicali/tradeshopExample/internal/httpErrors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
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
	userid, isExist := c.Get("userId")
	if !isExist {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "User not found", nil)))
		return
	}
	//userid := cast.ToInt(userID)
	userID := userid.(uuid.UUID)
	orders, err := o.service.GetAll(userID)
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, ordersToResponse(*orders))

}

func (o *orderHandler) add(c *gin.Context) {
	userid, isExist := c.Get("userId")
	if !isExist {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "User not found", nil)))
		return
	}
	//userid := cast.ToInt(userID)
	userID := userid.(uuid.UUID)
	order, err := o.service.Create(userID)

	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, OrderToResponse(order))

}

func (o *orderHandler) cancel(c *gin.Context) {
	userid, isExist := c.Get("userId")
	if !isExist {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "User not found", nil)))
		return
	}
	//userID := cast.ToInt(userID)
	userID := userid.(uuid.UUID)
	orderID, err := uuid.FromBytes([]byte((c.Param("id"))))

	err = o.service.Cancel(userID, orderID)

	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, "Order Cancel Complete")
}
