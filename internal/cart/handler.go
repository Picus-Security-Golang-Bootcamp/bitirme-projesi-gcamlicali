package cart

import (
	httpErr "github.com/gcamlicali/tradeshopExample/internal/httpErrors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/http"
)

type cartHandler struct {
	repo *CartRepositoy
}

func NewCartHandler(r *gin.RouterGroup, repo *CartRepositoy) {
	h := &cartHandler{repo: repo}

	r.GET("/", h.get)
	//r.POST("/create", h.create)
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
	id := cast.ToInt(userID)
	cart, err := ch.repo.GetByUserID(id)

	if err != nil {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, err.Error(), nil)))
		return
	}

	c.JSON(http.StatusOK, CartToResponse(cart))
}
