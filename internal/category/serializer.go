package category

import (
	"github.com/gcamlicali/tradeshopExample/internal/api"
	"github.com/gcamlicali/tradeshopExample/internal/models"
)

func catModelToApi(a *models.Category) *api.Category {
	return &api.Category{
		ID:   a.ID.String(),
		Name: a.Name,
	}

}

func catsModelToApi(cs *[]models.Category) []*api.Category {
	categories := make([]*api.Category, 0)
	for _, c := range *cs {
		categories = append(categories, catModelToApi(&c))
	}
	return categories
}
