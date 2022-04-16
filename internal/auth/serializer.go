package auth

import (
	"github.com/gcamlicali/tradeshopExample/internal/api"
	"github.com/gcamlicali/tradeshopExample/internal/models"
)

func userApiToModel(a *api.User) *models.User {
	return &models.User{
		Mail:      a.Email,
		Password:  a.Password,
		FirstName: a.FirstName,
		LastName:  a.LastName,
		Mobile:    a.Phone,
		IsAdmin:   false,
	}
}
