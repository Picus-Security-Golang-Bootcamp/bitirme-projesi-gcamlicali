package product

import (
	"github.com/gcamlicali/tradeshopExample/internal/api"
	httpErr "github.com/gcamlicali/tradeshopExample/internal/httpErrors"
	"github.com/gcamlicali/tradeshopExample/pkg/config"
	mw "github.com/gcamlicali/tradeshopExample/pkg/middleware"
	"github.com/gcamlicali/tradeshopExample/pkg/pagination"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
	"github.com/spf13/cast"
	"net/http"
	"strconv"
)

type productHandler struct {
	service Service
}

func NewProductHandler(r *gin.RouterGroup, service Service, cfg *config.Config) {
	h := &productHandler{service: service}

	r.GET("/", h.getAll)

	signedRoute := r.Group("/signed")
	signedRoute.Use(mw.AuthMiddleware(cfg.JWTConfig.SecretKey))
	signedRoute.DELETE("/:SKU", h.delete)
	signedRoute.PUT("/:SKU", h.update)
	signedRoute.POST("/addBulk", h.addBulk)
	signedRoute.POST("/addSingle", h.addSingle)
}

func (p *productHandler) addBulk(c *gin.Context) {
	adminInterface, isExist := c.Get("isAdmin")
	if !isExist {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "Admin not found", nil)))
		return
	}

	isAdmin := cast.ToBool(adminInterface)
	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to use this endpoint!"})
		return
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "Can not request body", err.Error())))
		return
	}

	err = p.service.AddBulk(file)
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, "Products uploaded and created")
	return
}
func (p *productHandler) addSingle(c *gin.Context) {
	adminInterface, isExist := c.Get("isAdmin")
	if !isExist {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "Admin not found", nil)))
		return
	}

	isAdmin := cast.ToBool(adminInterface)
	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to use this endpoint!"})
		return
	}
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

	product, err := p.service.AddSingle(*productBody)
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, ProductToResponse(product))
}
func (p *productHandler) getAll(c *gin.Context) {
	pageIndex, pageSize := pagination.GetPaginationParametersFromRequest(c)

	products, count, err := p.service.GetAll(pageIndex, pageSize)
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	paginatedResult := pagination.NewFromGinRequest(c, count)
	paginatedResult.Items = productsToResponse(*products)

	c.JSON(http.StatusOK, paginatedResult)
}
func (p *productHandler) delete(c *gin.Context) {
	adminInterface, isExist := c.Get("isAdmin")
	if !isExist {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "Admin not found", nil)))
		return
	}

	isAdmin := cast.ToBool(adminInterface)
	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to use this endpoint!"})
		return
	}

	SKU, err := strconv.Atoi(c.Param("SKU"))
	if err != nil {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "SKU is not integer", err.Error())))
		return
	}

	err = p.service.Delete(SKU)
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, "Product delete succesful")

}
func (p *productHandler) update(c *gin.Context) {
	adminInterface, isExist := c.Get("isAdmin")
	if !isExist {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "Admin not found", nil)))
		return
	}

	isAdmin := cast.ToBool(adminInterface)
	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to use this endpoint!"})
		return
	}

	SKU, err := strconv.Atoi(c.Param("SKU"))
	if err != nil {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "SKU is not integer", err.Error())))
	}

	reqProduct := api.ProductUp{}
	if err := c.Bind(&reqProduct); err != nil {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "check your request body", err.Error())))
		return
	}

	updatedProduct, err := p.service.Update(SKU, &reqProduct)
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, ProductToResponse(updatedProduct))
}
