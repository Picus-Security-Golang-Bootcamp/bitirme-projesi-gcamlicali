package auth

import (
	"github.com/gcamlicali/tradeshopExample/internal/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthRepositoy struct {
	db *gorm.DB
}
type IAuthRepository interface {
	Create(a *models.User) (*models.User, error)
	GetByMail(mail string) (*models.User, error)
	CheckAndCreateAdmin(user *models.User) bool
}

func NewAuthRepository(db *gorm.DB) *AuthRepositoy {
	return &AuthRepositoy{db: db}
}

func (r *AuthRepositoy) Create(a *models.User) (*models.User, error) {
	zap.L().Debug("user.repo.create", zap.Reflect("userBody", a))
	if err := r.db.Create(a).Error; err != nil {
		zap.L().Error("user.repo.Create failed to create user", zap.Error(err))
		return nil, err
	}
	return a, nil
}

func (r *AuthRepositoy) GetByMail(mail string) (*models.User, error) {
	zap.L().Debug("User.repo.getByID", zap.Reflect("mail", mail))

	var user = &models.User{}
	if result := r.db.Where(models.User{Mail: &mail}).First(&user); result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (r *AuthRepositoy) CheckAndCreateAdmin(user *models.User) bool {
	zap.L().Debug("User.repo.crateAdmin", zap.Reflect("admin", user))

	if r.db.Where("mail = ?", user.Mail).Updates(&user).RowsAffected == 0 {
		_, err := r.Create(user)
		if err != nil {
			zap.L().Error("Create Admin Data Error : ", zap.Error(err))
		}

		return true
	}
	return false
}

func (r *AuthRepositoy) Migration() {
	r.db.AutoMigrate(&models.User{})
}
