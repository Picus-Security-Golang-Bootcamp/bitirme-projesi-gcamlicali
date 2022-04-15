package category

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
)

type categoryHandler struct {
	service Service
}

func NewCategoryHandler(r *gin.RouterGroup, service Service, cfg *config.Config) {
	a := categoryHandler{service: service}

	r.GET("/", a.getAll)

	signedRoute := r.Group("/signed")
	signedRoute.Use(mw.AuthMiddleware(cfg.JWTConfig.SecretKey))
	signedRoute.DELETE("/", a.delete)
	signedRoute.PUT("/", a.update)
	signedRoute.POST("/addBulk", a.addBulk)
	signedRoute.POST("/addSingle", a.addSingle)
}

func (h *categoryHandler) getAll(c *gin.Context) {

	pageIndex, pageSize := pagination.GetPaginationParametersFromRequest(c)

	categories, count, err := h.service.GetAll(pageIndex, pageSize)
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}
	paginatedResult := pagination.NewFromGinRequest(c, count)
	paginatedResult.Items = catsModelToApi(categories)

	c.JSON(http.StatusOK, paginatedResult)
}

func (h *categoryHandler) addBulk(c *gin.Context) {

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
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "Request File download error", err.Error())))
		return
	}

	defer file.Close()

	err = h.service.AddBulk(file)
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, "Categories uploaded and created")
	return
}

func (h *categoryHandler) addSingle(c *gin.Context) {
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

	reqCategory := api.Category{}
	if err := c.Bind(&reqCategory); err != nil {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "check your request body", err.Error())))
		return
	}

	if err := reqCategory.Validate(strfmt.NewFormats()); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	createdCategory, err := h.service.AddSingle(reqCategory)
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, createdCategory)
}

func (h *categoryHandler) delete(c *gin.Context) {
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

	catName := c.Param("Name")

	err := h.service.Delete(catName)
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, "Category and its items Delete Succesful")

}

func (h *categoryHandler) update(c *gin.Context) {
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

	catName := c.Param("Name")

	reqCategory := api.Category{}
	if err := c.Bind(&reqCategory); err != nil {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "check your request body", err.Error())))
		return
	}

	if err := reqCategory.Validate(strfmt.NewFormats()); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	updatedCategory, err := h.service.Update(catName, &reqCategory)
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, catModelToApi(updatedCategory))
}
