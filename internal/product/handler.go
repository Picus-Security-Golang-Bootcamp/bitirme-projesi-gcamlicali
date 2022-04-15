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
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "Can not request body", nil)))
		return
	}

	err = p.service.AddBulk(file)
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	//c.JSON(http.StatusCreated, "Categories uploaded and created")
	//return
	//
	//record, err := csvRead.ReadFile(file)
	//
	//if err != nil {
	//	c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusInternalServerError, "Can not read csv file", nil)))
	//	return
	//}
	//
	//for _, line := range record {
	//	proEntity := models.Product{}
	//	proEntity.CategoryName = line[0]
	//	_, err := p.catRepo.GetByName(proEntity.CategoryName)
	//	if err != nil {
	//		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusNotFound, "Category not found", proEntity.CategoryName)))
	//		continue
	//	}
	//	proEntity.Name = line[1]
	//	SKU, err := strconv.Atoi(line[2])
	//	if err != nil {
	//		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "SKU is not integer", proEntity.Name)))
	//		continue
	//	}
	//	proEntity.SKU = SKU
	//	proEntity.Description = line[3]
	//	price, err := strconv.Atoi(line[4])
	//	if err != nil {
	//		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "Price is not integer", proEntity.Name)))
	//		continue
	//	}
	//	proEntity.Price = int32(price)
	//	unitStock, err := strconv.Atoi(line[5])
	//	if err != nil {
	//		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "UnitStock is not integer", proEntity.Name)))
	//		continue
	//	}
	//	proEntity.UnitStock = int32(unitStock)
	//
	//	_, err = p.proRepo.create(&proEntity)
	//	if err != nil {
	//		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, err.Error(), nil)))
	//		continue
	//	}
	//}

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
	//catName := productBody.CategoryName
	//
	//category, err := p.catRepo.GetByName(*catName)
	//if err != nil {
	//	c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusNotFound, "Category not found", nil)))
	//	return
	//}
	//
	//prod := responseToProduct(productBody)
	//prod.CategoryName = *category.Name
	////log.Println("prod cat id: ", prod.CategoryID)
	//product, err := p.proRepo.create(prod)
	//if err != nil {
	//	c.JSON(httpErr.ErrorResponse(err))
	//	return
	//}

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
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "SKU is not integer", err)))
		return
	}

	err = p.service.Delete(SKU)
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}
	//err = p.proRepo.deleteBySku(SKU)
	//
	//if errors.Is(err, gorm.ErrRecordNotFound) {
	//	c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "Product not found", err)))
	//	return
	//}
	//if err != nil {
	//	c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusInternalServerError, "Delete product error", err)))
	//	return
	//}

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
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "SKU is not integer", err)))
	}

	reqProduct := api.ProductUp{}
	if err := c.Bind(&reqProduct); err != nil {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "check your request body", nil)))
		return
	}

	updatedProduct, err := p.service.Update(SKU, &reqProduct)
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	//product, err := p.proRepo.GetBySKU(SKU)
	//if err != nil {
	//	c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusInternalServerError, "Get product error", err)))
	//}
	//
	//if reqCategory.Name != "" {
	//	product.Name = reqCategory.Name
	//}
	//if reqCategory.CategoryName != "" {
	//	product.CategoryName = reqCategory.CategoryName
	//}
	//
	//if reqCategory.Description != "" {
	//	product.Description = reqCategory.Description
	//}
	//if reqCategory.Price != 0 {
	//	product.Price = reqCategory.Price
	//}
	//log.Println("SKU: ", reqCategory.Sku)
	//if reqCategory.Sku != 0 {
	//	product.SKU = int(reqCategory.Sku)
	//}
	//if reqCategory.UnitStock != 0 {
	//	product.UnitStock = reqCategory.UnitStock
	//}
	//
	//updatedProduct, err := p.proRepo.Update(product)
	//if err != nil {
	//	c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusInternalServerError, "Update product error", err)))
	//}

	c.JSON(http.StatusOK, ProductToResponse(updatedProduct))
}
