package product

import (
	"github.com/gcamlicali/tradeshopExample/internal/api"
	category "github.com/gcamlicali/tradeshopExample/internal/category"
	httpErr "github.com/gcamlicali/tradeshopExample/internal/httpErrors"
	"github.com/gcamlicali/tradeshopExample/internal/models"
	"github.com/gcamlicali/tradeshopExample/pkg/csv"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
	"log"
	"net/http"
)

type productHandler struct {
	proRepo *ProductRepositoy
	catRepo *category.CategoryRepositoy
}

func NewProductHandler(r *gin.RouterGroup, proRepo *ProductRepositoy, catRepo *category.CategoryRepositoy) {
	h := &productHandler{proRepo: proRepo, catRepo: catRepo}

	r.GET("/", h.getAll)
	addRoute := r.Group("/add")
	addRoute.POST("/bulkItems", h.addBulk)
	addRoute.POST("/singleItem", h.addSingle)
}

func (p *productHandler) addBulk(c *gin.Context) {

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "Can not request body", nil)))
		return
	}

	defer file.Close()

	record, err := csvRead.ReadFile(file)

	if err != nil {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusInternalServerError, "Can not read csv file", nil)))
		return
	}

	for _, line := range record {
		proEntity := models.Product{}
		proEntity.CategoryName = line[0]
		log.Println("line okundu 0. column: ", line[0])
		_, err := p.catRepo.GetByName(proEntity.CategoryName)
		if err != nil {
			c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusNotFound, "Category not found", proEntity.CategoryName)))
			continue
		}
		proEntity.Name = line[1]
		proEntity.SKU = line[2]
		proEntity.Description = line[3]

		_, err = p.proRepo.create(&proEntity)
		if err != nil {
			c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, err.Error(), nil)))
			continue
		}
	}

	c.JSON(http.StatusCreated, "Products uploaded and created")
	return
}

func (p *productHandler) addSingle(c *gin.Context) {
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

	catName := productBody.CategoryName

	category, err := p.catRepo.GetByName(*catName)
	if err != nil {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusNotFound, "Category not found", nil)))
		return
	}

	prod := responseToProduct(productBody)
	prod.CategoryName = *category.Name
	//log.Println("prod cat id: ", prod.CategoryID)
	product, err := p.proRepo.create(prod)
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, ProductToResponse(product))
}
func (p *productHandler) getAll(c *gin.Context) {
	products, err := p.proRepo.getAll()
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, productsToResponse(*products))
}
