package auth

import (
	"github.com/gcamlicali/tradeshopExample/internal/models"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthRepositoy struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepositoy {
	return &AuthRepositoy{db: db}
}

func (r *AuthRepositoy) create(a *models.User) (*models.User, error) {
	zap.L().Debug("user.repo.create", zap.Reflect("userBody", a))
	if err := r.db.Create(a).Error; err != nil {
		zap.L().Error("user.repo.Create failed to create user", zap.Error(err))
		return nil, err
	}
	return a, nil
}

func (r *AuthRepositoy) getByID(id string) (*models.User, error) {
	zap.L().Debug("user.repo.getByID", zap.Reflect("id", id))

	var user = &models.User{}
	if result := r.db.First(&user, id); result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (r *AuthRepositoy) getByMail(mail string) (*models.User, error) {
	zap.L().Debug("User.repo.getByID", zap.Reflect("mail", mail))

	var user = &models.User{}
	if result := r.db.Where(models.User{Mail: &mail}).First(&user); result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (r *AuthRepositoy) FillAdminData() {
	admin := models.GetAdmin()
	//Then encrypt the password
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(*admin.Password), bcrypt.DefaultCost)
	passBeforeReg := string(hashPassword)
	admin.Password = &passBeforeReg
	if r.db.Where("mail = ?", admin.Mail).Updates(&admin).RowsAffected == 0 {
		_, err := r.create(&admin)

		if err != nil {
			zap.L().Error("Create Admin Data Error : ", zap.Error(err))
		}
	}
}

func (r *AuthRepositoy) Migration() {
	r.db.AutoMigrate(&models.User{})
}
