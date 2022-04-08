package category

import (
	"encoding/csv"
	httpErr "github.com/gcamlicali/tradeshopExample/internal/httpErrors"
	"github.com/gcamlicali/tradeshopExample/internal/models"
	"github.com/gcamlicali/tradeshopExample/pkg/config"
	mw "github.com/gcamlicali/tradeshopExample/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"log"
	"net/http"
)

type categoryHandler struct {
	repo *CategoryRepositoy
}

func NewCategoryHandler(r *gin.RouterGroup, repo *CategoryRepositoy, cfg *config.Config) {
	a := categoryHandler{repo: repo}

	r.GET("/", a.getAll)
	addRoute := r.Group("/add")
	addRoute.Use(mw.AuthMiddleware(cfg.JWTConfig.SecretKey))
	addRoute.POST("/bulkItems", a.addBulk)
}

func (h *categoryHandler) getAll(c *gin.Context) {

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
	reader := csv.NewReader(file)
	record, err := reader.ReadAll()
	if err != nil {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusInternalServerError, "Can not read csv file", nil)))
		return
	}

	for _, line := range record {
		catEntity := models.Category{}
		catEntity.Name = &line[0]
		_, err = h.repo.create(&catEntity)
		if err != nil {
			log.Println("DB yazma problemi: ", err)
		}
	}
	c.JSON(http.StatusCreated, "Categories uploaded and created")
	return
}
