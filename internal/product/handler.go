package product

import (
	"github.com/gcamlicali/tradeshopExample/internal/api"
	httpErr "github.com/gcamlicali/tradeshopExample/internal/httpErrors"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
	"net/http"
)

type productHandler struct {
	repo *ProductRepositoy
}

func NewProductHandler(r *gin.RouterGroup, repo *ProductRepositoy) {
	h := &productHandler{repo: repo}

	r.GET("/", h.getAll)
}

func (p *productHandler) create(c *gin.Context) {
	productBody := &api.Product{}
	if err := c.Bind(&productBody); err != nil {
		c.JSON(httpErr.ErrorResponse(httpErr.CannotBindGivenData))
		return
	}
	// Validating all required areas
	if err := productBody.Validate(strfmt.NewFormats()); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	product, err := p.repo.create(responseToProduct(productBody))
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}
 
	c.JSON(http.StatusOK, ProductToResponse(product))
}
func (p *productHandler) getAll(c *gin.Context) {
	products, err := p.repo.getAll()
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}
	c.JSON(http.StatusOK, productsToResponse(products))
}
