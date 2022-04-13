package auth

import (
	"github.com/gcamlicali/tradeshopExample/internal/api"
	httpErr "github.com/gcamlicali/tradeshopExample/internal/httpErrors"
	"github.com/gcamlicali/tradeshopExample/pkg/config"
	jwtHelper "github.com/gcamlicali/tradeshopExample/pkg/jwt"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

type authService struct {
	cfg  *config.Config
	repo AuthRepositoy
}

type Service interface {
	SignIn(login *api.Login) (string, error)
	SignUp(login *api.User) (string, error)
}

func NewAuthService(repo AuthRepositoy, cfg *config.Config) Service {
	return &authService{repo: repo, cfg: cfg}
}

func (a *authService) SignIn(login *api.Login) (string, error) {
	user, err := a.repo.getByMail(*login.Email)
	if err != nil {
		//c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "user not found", nil)))
		return "", httpErr.NewRestError(http.StatusBadRequest, "user not found", nil)
	}
	if user == nil {
		//c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "user not found", nil)))
		return "", httpErr.NewRestError(http.StatusBadRequest, "user not found", nil)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(*login.Password)); err != nil {
		//c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusBadRequest, "Wrong Password", nil)))
		return "", httpErr.NewRestError(http.StatusBadRequest, "Wrong Password", nil)
	}
	jwtClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.ID,
		"email":  user.Mail,
		"iat":    time.Now().Unix(),
		"iss":    os.Getenv("ENV"),
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
		"roles":  user.IsAdmin,
	})

	token := jwtHelper.GenerateToken(jwtClaims, a.cfg.JWTConfig.SecretKey)

	return token, nil
}

func (a *authService) SignUp(login *api.User) (string, error) {

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(*login.Password), bcrypt.DefaultCost)
	if err != nil {
		//c.JSON(httpErr.ErrorResponse(httpErr.NewRestError(http.StatusUnprocessableEntity, "encryption error", nil)))
		return "", httpErr.NewRestError(http.StatusUnprocessableEntity, "encryption error", nil)
	}
	passBeforeReg := string(hashPassword)
	login.Password = &passBeforeReg

	createdUser, err := a.repo.create(userApiToModel(login))
	if err != nil {
		//c.JSON(httpErr.ErrorResponse(err))
		return "", err
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

	return token, nil
}
