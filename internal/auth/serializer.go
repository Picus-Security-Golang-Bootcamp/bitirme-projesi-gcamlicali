package auth

import (
	"github.com/gcamlicali/tradeshopExample/internal/api"
	"github.com/gcamlicali/tradeshopExample/internal/models"
	"gorm.io/gorm"
)

func userApiToModel(a *api.User) *models.User {
	return &models.User{
		Model:     gorm.Model{},
		Mail:      a.Email,
		Password:  a.Password,
		FirstName: a.FirstName,
		LastName:  a.LastName,
		Mobile:    a.Phone,
		IsAdmin:   false,
	}
}
