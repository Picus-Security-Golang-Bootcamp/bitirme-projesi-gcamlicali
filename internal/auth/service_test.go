package auth

import (
	"github.com/gcamlicali/tradeshopExample/internal/api"
	"github.com/gcamlicali/tradeshopExample/internal/cart"
	"github.com/gcamlicali/tradeshopExample/internal/models"
	"github.com/gcamlicali/tradeshopExample/pkg/config"
	"github.com/go-openapi/errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

var (
	adminMail         = "admin@admin.com"
	adminPass         = "mockPass"
	adminName         = "adminFirst"
	adminLast         = "adminLast"
	userID            = uuid.New()
	AdminEncPass, err = bcrypt.GenerateFromPassword([]byte(adminPass), bcrypt.DefaultCost)
	AdminPassStrEnc   = string(AdminEncPass)
	AdmWrongPass      = "WrongPass"
	UserWrongMail     = "WrongMail"
	userMail          = "user@user.com"
	userPass          = "userPass"
	userName          = "userName"
	userLast          = "userLast"

	admin = models.User{
		Mail:      &adminMail,
		ID:        userID,
		IsAdmin:   true,
		Password:  &AdminPassStrEnc,
		FirstName: &adminName,
		LastName:  &adminLast,
		Mobile:    "333111",
	}
	user = models.User{
		Mail:      &userMail,
		Password:  &userPass,
		IsAdmin:   false,
		Mobile:    "2222444",
		FirstName: &userName,
		LastName:  &userLast,
	}

	conf = config.Config{
		JWTConfig: config.JWTConfig{
			SecretKey:   "Test",
			SessionTime: 30,
		},
	}

	logInUser = api.Login{
		Password: &adminPass,
		Email:    &adminMail,
	}
	logInUserU = api.Login{
		Password: &adminPass,
		Email:    &UserWrongMail,
	}
	logInUserW = api.Login{
		Password: &AdmWrongPass,
		Email:    &adminMail,
	}
	SignUpUser = api.User{
		Email:     &userMail,
		Password:  &userPass,
		IsAdmin:   false,
		Phone:     "2222444",
		FirstName: &userName,
		LastName:  &userLast,
	}
)

func Test_authService_SignIn(t *testing.T) {
	type fields struct {
		cfg   *config.Config
		repo  IAuthRepository
		cRepo cart.ICartRepository
	}
	type args struct {
		login *api.Login
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		//want    string
		wantErr bool
	}{
		{
			name: "authService_SignIn_ShouldSuccess",
			fields: fields{
				repo: &authMockRepo{
					Items: []models.User{
						admin,
					},
				},
				cRepo: &cartMockRepo{
					Items: []models.Cart{},
				},
				cfg: &conf,
			},
			args: args{
				login: &logInUser,
			},
			wantErr: false,
		},
		{
			name: "authService_SignIn_ErrorWrongPass_ShouldFail",
			fields: fields{
				repo: &authMockRepo{
					Items: []models.User{
						admin,
					},
				},
				cRepo: &cartMockRepo{
					Items: []models.Cart{},
				},
				cfg: &conf,
			},
			args: args{
				login: &logInUserW,
			},
			wantErr: true,
		},
		{
			name: "authService_SignIn_ErrorUnknownUser_ShouldFail",
			fields: fields{
				repo: &authMockRepo{
					Items: []models.User{
						admin,
					},
				},
				cRepo: &cartMockRepo{
					Items: []models.Cart{},
				},
				cfg: &conf,
			},
			args: args{
				login: &logInUserU,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &authService{
				cfg:   tt.fields.cfg,
				repo:  tt.fields.repo,
				cRepo: tt.fields.cRepo,
			}
			_, err := a.SignIn(tt.args.login)
			if (err != nil) != tt.wantErr {
				t.Errorf("SignIn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_authService_SignUp(t *testing.T) {
	type fields struct {
		cfg   *config.Config
		repo  IAuthRepository
		cRepo cart.ICartRepository
	}
	type args struct {
		login *api.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "authService_SignUp_ShouldSuccess",
			fields: fields{
				repo: &authMockRepo{
					Items: []models.User{
						admin,
					},
				},
				cRepo: &cartMockRepo{
					Items: []models.Cart{},
				},
				cfg: &conf,
			},
			args: args{
				login: &SignUpUser,
			},
			wantErr: false,
		},
		{
			name: "authService_SignUp_ErrorUserExist_ShouldFail",
			fields: fields{
				repo: &authMockRepo{
					Items: []models.User{
						admin,
						user,
					},
				},
				cRepo: &cartMockRepo{
					Items: []models.Cart{},
				},
				cfg: &conf,
			},
			args: args{
				login: &SignUpUser,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &authService{
				cfg:   tt.fields.cfg,
				repo:  tt.fields.repo,
				cRepo: tt.fields.cRepo,
			}
			_, err := a.SignUp(tt.args.login)
			if (err != nil) != tt.wantErr {
				t.Errorf("SignUp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

type cartMockRepo struct {
	Items []models.Cart
}
type authMockRepo struct {
	Items []models.User
}

func (c *cartMockRepo) Create(a *models.Cart) (*models.Cart, error) {
	c.Items = append(c.Items, *a)
	return a, nil
}
func (c *cartMockRepo) GetByUserID(userID uuid.UUID) (*models.Cart, error) {
	cart := models.Cart{}
	for _, item := range c.Items {
		if item.UserID == userID {
			cart = item
			return &cart, nil
		}
	}
	return nil, errors.New(400, "User Cart not found")
}
func (c *cartMockRepo) Update(a *models.Cart) (*models.Cart, error) {
	for i, item := range c.Items {
		if item.ID == a.ID {
			c.Items[i] = *a
			return a, nil
		}
	}
	return nil, errors.New(400, "Cart not found")
}

func (c *authMockRepo) Create(a *models.User) (*models.User, error) {
	for _, item := range c.Items {
		if item.Mail == a.Mail {
			return nil, errors.New(400, "User Already Exist")
		}
	}
	c.Items = append(c.Items, *a)
	return a, nil
}
func (c *authMockRepo) GetByMail(mail string) (*models.User, error) {
	user := models.User{}
	for _, item := range c.Items {
		if *item.Mail == mail {
			user = item
			return &user, nil
		}
	}
	return nil, errors.New(400, "User not found")
}
func (c *authMockRepo) CheckAndCreateAdmin(user *models.User) bool {
	admin := models.User{}
	for _, item := range c.Items {
		if item.Mail == admin.Mail {
			return true
		}
	}
	c.Create(user)
	c.Items = append(c.Items, *user)
	return false
}
