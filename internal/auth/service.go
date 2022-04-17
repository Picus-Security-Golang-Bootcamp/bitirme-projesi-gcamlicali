package auth

import (
	"github.com/gcamlicali/tradeshopExample/internal/api"
	"github.com/gcamlicali/tradeshopExample/internal/cart"
	httpErr "github.com/gcamlicali/tradeshopExample/internal/httpErrors"
	"github.com/gcamlicali/tradeshopExample/internal/models"
	"github.com/gcamlicali/tradeshopExample/pkg/config"
	jwtHelper "github.com/gcamlicali/tradeshopExample/pkg/jwt"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"time"
)

type authService struct {
	cfg   *config.Config
	repo  IAuthRepository
	cRepo cart.ICartRepository
}

type Service interface {
	SignIn(login *api.Login) (string, error)
	SignUp(login *api.User) (string, error)
	FillAdminData()
}

func NewAuthService(repo IAuthRepository, cRepo cart.ICartRepository, cfg *config.Config) Service {
	return &authService{repo: repo, cRepo: cRepo, cfg: cfg}
}

func (a *authService) SignIn(login *api.Login) (string, error) {

	//Find user by api response mail in DB
	user, err := a.repo.GetByMail(*login.Email)
	if err != nil {
		return "", httpErr.NewRestError(http.StatusBadRequest, "User get err", err.Error())
	}
	if user == nil {
		return "", httpErr.NewRestError(http.StatusBadRequest, "user not found", nil)
	}

	// Compare user apiModel password with Encrypted user password in Database
	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(*login.Password)); err != nil {
		return "", httpErr.NewRestError(http.StatusBadRequest, "Wrong Password", err.Error())
	}

	//Generate token for user
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

	//Encrypt the user password
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(*login.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", httpErr.NewRestError(http.StatusUnprocessableEntity, "encryption error", err.Error())
	}
	passBeforeReg := string(hashPassword)
	login.Password = &passBeforeReg

	//Create  api response based user
	createdUser, err := a.repo.Create(userApiToModel(login))
	if err != nil {
		return "", httpErr.NewRestError(http.StatusInternalServerError, "Can't create user", err.Error())
	}

	//Generate token for user
	jwtClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": createdUser.ID,
		"email":  createdUser.Mail,
		"iat":    time.Now().Unix(),
		"iss":    os.Getenv("ENV"),
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
		"roles":  createdUser.IsAdmin,
	})
	token := jwtHelper.GenerateToken(jwtClaims, a.cfg.JWTConfig.SecretKey)

	//Initial, create a new cart for user
	var cart = &models.Cart{}
	cart.UserID = createdUser.ID
	_, err = a.cRepo.Create(cart)
	if err != nil {
		return "", httpErr.NewRestError(http.StatusInternalServerError, "Can't create new cart for new user", err.Error())
	}

	return token, nil
}

func (a *authService) FillAdminData() {
	//Get Admin data from json file
	admin := models.GetAdmin()

	//Encrypt the admin password
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(*admin.Password), bcrypt.DefaultCost)
	passBeforeReg := string(hashPassword)
	admin.Password = &passBeforeReg

	//If admin existed ok, but doesn't exist create a new admin user
	isNew := a.repo.CheckAndCreateAdmin(&admin)

	//If new admin created, create a cart for admin
	if isNew {
		var cart = &models.Cart{}
		cart.UserID = admin.ID
		_, err := a.cRepo.Create(cart)

		if err != nil {
			log.Fatal("Can't create new cart for new user")
		}
	}
}
