package auth

import (
	"github.com/gcamlicali/tradeshopExample/internal/api"
	httpErr "github.com/gcamlicali/tradeshopExample/internal/httpErrors"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"

	"net/http"
)

type authHandler struct {
	service Service
}

func NewAuthHandler(r *gin.RouterGroup, service Service) {
	a := authHandler{service: service}

	r.POST("/signin", a.signin)
	r.POST("/signup", a.signup)
}

func (a *authHandler) signin(c *gin.Context) {
	req := api.Login{}

	if err := c.Bind(&req); err != nil {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "check your request body", nil)))
		return
	}

	if err := req.Validate(strfmt.NewFormats()); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}
	token, err := a.service.SignIn(&req)
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, token)
}

func (a *authHandler) signup(c *gin.Context) {

	var reqUser api.User

	if err := c.Bind(&reqUser); err != nil {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "check your request body", nil)))
		return
	}

	if err := reqUser.Validate(strfmt.NewFormats()); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}
	token, err := a.service.SignUp(&reqUser)
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, token)
}
