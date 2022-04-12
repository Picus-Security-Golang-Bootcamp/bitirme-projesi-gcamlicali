package auth

import (
	"github.com/gcamlicali/tradeshopExample/internal/api"
	httpErr "github.com/gcamlicali/tradeshopExample/internal/httpErrors"
	"github.com/gcamlicali/tradeshopExample/pkg/config"
	jwtHelper "github.com/gcamlicali/tradeshopExample/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
	"github.com/golang-jwt/jwt"
	"log"

	//"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

type authHandler struct {
	cfg  *config.Config
	repo *AuthRepositoy
}

func NewAuthHandler(r *gin.RouterGroup, repo *AuthRepositoy, cfg *config.Config) {
	a := authHandler{cfg: cfg, repo: repo}

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
	log.Println("USER MAILI: ", *req.Email)
	user, err := a.repo.getByMail(*req.Email)
	if err != nil {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "user not found", nil)))
		return
	}
	if user == nil {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "user not found", nil)))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(*req.Password)); err != nil {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "Wrong Password", nil)))
		return
	}
	jwtClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.ID,
		"email":  user.Mail,
		"iat":    time.Now().Unix(),
		"iss":    os.Getenv("ENV"),
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
		"roles":  user.IsAdmin,
	})
	log.Println("JWTClaims: ", jwtClaims)
	token := jwtHelper.GenerateToken(jwtClaims, a.cfg.JWTConfig.SecretKey)
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

	//Then encrypt the password
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(*reqUser.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusUnprocessableEntity, "encryption error", nil)))
		return
	}
	passBeforeReg := string(hashPassword)
	reqUser.Password = &passBeforeReg

	createdUser, err := a.repo.create(userApiToModel(&reqUser))
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	jwtClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": createdUser.ID,
		"email":  createdUser.Mail,
		"iat":    time.Now().Unix(),
		"iss":    os.Getenv("ENV"),
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
		"roles":  createdUser.IsAdmin,
	})
	token := jwtHelper.GenerateToken(jwtClaims, a.cfg.JWTConfig.SecretKey)
	c.JSON(http.StatusCreated, token)
}
