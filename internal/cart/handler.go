package cart

import (
	"github.com/gcamlicali/tradeshopExample/internal/api"
	httpErr "github.com/gcamlicali/tradeshopExample/internal/httpErrors"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

type cartHandler struct {
	service Service
}

func NewCartHandler(r *gin.RouterGroup, service Service) {
	h := &cartHandler{service: service}

	r.GET("/", h.get)
	r.POST("/:SKU", h.add)
	r.PUT("/:SKU", h.update)
	r.DELETE("/:SKU", h.delete)
}

func (ch *cartHandler) get(c *gin.Context) {
	userid, isExist := c.Get("userId")
	if !isExist {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "User not found", nil)))
		return
	}
	//userid := cast.ToInt(userID)
	userID := userid.(uuid.UUID)

	cart, err := ch.service.Get(userID)
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, CartToResponse(cart))
}

func (ch *cartHandler) add(c *gin.Context) {
	userid, isExist := c.Get("userId")
	if !isExist {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "User not found", nil)))
		return
	}
	//userID := cast.ToString(userid)
	userID := userid.(uuid.UUID)
	paramID, err := strconv.Atoi(c.Param("SKU"))
	if err != nil {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "SKU is not integer", err.Error())))
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

	paramID, err := strconv.Atoi(c.Param("SKU"))
	if err != nil {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "SKU is not integer", err.Error())))
	}
	reqQuantity := api.ItemQuantity{}
	if err := c.Bind(&reqQuantity); err != nil {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "check your request body", err.Error())))
		return
	}
	if err := reqQuantity.Validate(strfmt.NewFormats()); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	//userID := cast.ToInt(userid)
	userID := userid.(uuid.UUID)
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

	paramID, err := strconv.Atoi(c.Param("SKU"))
	if err != nil {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "SKU is not integer", err.Error())))
	}

	//userID := cast.ToInt(userid)
	userID := userid.(uuid.UUID)
	cart, err := ch.service.Delete(userID, paramID)
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, CartToResponse(cart))
}
