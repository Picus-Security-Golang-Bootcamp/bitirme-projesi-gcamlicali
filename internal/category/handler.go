package category

import (
	"github.com/gcamlicali/tradeshopExample/internal/api"
	httpErr "github.com/gcamlicali/tradeshopExample/internal/httpErrors"
	"github.com/gcamlicali/tradeshopExample/pkg/config"
	mw "github.com/gcamlicali/tradeshopExample/pkg/middleware"
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
	addRoute := r.Group("/add")
	addRoute.Use(mw.AuthMiddleware(cfg.JWTConfig.SecretKey))
	addRoute.POST("/bulkItems", a.addBulk)
	addRoute.POST("/singleItem", a.addSingle)
}

func (h *categoryHandler) getAll(c *gin.Context) {

	categories, err := h.service.GetAll()
	if err != nil {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusInternalServerError, err.Error(), nil)))
		return
	}

	c.JSON(http.StatusOK, catsModelToApi(categories))
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
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "Can not request body", nil)))
		return
	}

	defer file.Close()

	err = h.service.AddBulk(file)
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}
	//record, err := csvRead.ReadFile(file)
	//if err != nil {
	//	c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusInternalServerError, "Can not read csv file", nil)))
	//	return
	//}
	//
	//for _, line := range record {
	//	catEntity := models.Category{}
	//	catEntity.Name = &line[0]
	//	_, err = h.service.Create(&catEntity)
	//	if err != nil {
	//		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, err.Error(), nil)))
	//	}
	//}
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
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "check your request body", nil)))
		return
	}

	if err := reqCategory.Validate(strfmt.NewFormats()); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	//dbCat := models.Category{}
	//dbCat.Name = reqCategory.Name
	//createdCategory, err := h.service.Create(&dbCat)
	//if err != nil {
	//	c.JSON(httpErr.ErrorResponse(err))
	//	return
	//}
	createdCategory, err := h.service.AddSingle(reqCategory)
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, createdCategory)

}
