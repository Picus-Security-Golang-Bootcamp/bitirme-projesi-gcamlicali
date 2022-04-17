package category

import (
	"github.com/gcamlicali/tradeshopExample/internal/api"
	"github.com/gcamlicali/tradeshopExample/internal/models"
	"github.com/go-openapi/errors"
	"github.com/google/uuid"
	"testing"
)

var (
	categoryName = "CategoryName1"
	categoryID   = uuid.New()
	category1    = models.Category{
		ID:   categoryID,
		Name: &categoryName,
	}
)

func Test_categoryService_Create(t *testing.T) {
	categoryName := "CategoryName1"

	type fields struct {
		catRepo ICategoryRepository
	}
	type args struct {
		a *models.Category
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Category
		wantErr bool
	}{
		{
			name: "categoryService_CreateCatogory_ShouldSuccess",
			fields: fields{
				catRepo: &categoryMockRepo{
					Items: []models.Category{},
				},
			},
			args: args{
				a: &models.Category{
					ID:   uuid.New(),
					Name: &categoryName,
				},
			},
			wantErr: false,
		},
		{
			name: "categoryService_CreateCategory_Duplicate_ShouldFailed",
			fields: fields{
				catRepo: &categoryMockRepo{
					Items: []models.Category{
						{
							ID:   uuid.New(),
							Name: &categoryName,
						},
					},
				},
			},
			args: args{
				a: &models.Category{
					ID:   uuid.New(),
					Name: &categoryName,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := categoryService{
				repo: tt.fields.catRepo,
			}
			_, err := c.Create(tt.args.a)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func Test_categoryService_GetAll(t *testing.T) {
	categoryName := "CategoryName1"

	type fields struct {
		repo ICategoryRepository
	}
	type args struct {
		pageIndex int
		pageSize  int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		//want    *[]models.Category
		//want1   int
		wantErr bool
	}{
		{
			name: "categoryService_GetAll_ShouldSuccess",
			fields: fields{
				repo: &categoryMockRepo{
					Items: []models.Category{
						{
							ID:   uuid.New(),
							Name: &categoryName,
						},
					},
				},
			},
			args: args{
				pageIndex: 1,
				pageSize:  1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := categoryService{
				repo: tt.fields.repo,
			}
			_, _, err := c.GetAll(tt.args.pageIndex, tt.args.pageSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_categoryService_AddSingle(t *testing.T) {

	type fields struct {
		repo ICategoryRepository
	}
	type args struct {
		category api.Category
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "categoryService_CategoryAddSingle_ShouldSuccess",
			fields: fields{
				repo: &categoryMockRepo{
					Items: []models.Category{},
				},
			},

			args: args{
				category: api.Category{
					Name: &categoryName,
				},
			},
			wantErr: false,
		},
		{
			name: "categoryService_CategoryAddSingle_Duplicate_ShouldFail",
			fields: fields{
				repo: &categoryMockRepo{
					Items: []models.Category{
						category1,
					},
				},
			},
			args: args{
				category: api.Category{
					Name: &categoryName,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := categoryService{
				repo: tt.fields.repo,
			}
			_, err := c.AddSingle(tt.args.category)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddSingle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

type categoryMockRepo struct {
	Items []models.Category
}

func (c *categoryMockRepo) Create(a *models.Category) (*models.Category, error) {

	for _, item := range c.Items {
		if item.Name == a.Name {
			return nil, errors.New(400, "Item should be unique on database")
		}
	}
	c.Items = append(c.Items, *a)
	return a, nil
}
func (c *categoryMockRepo) GetByName(name string) (*models.Category, error) {
	category := &models.Category{}
	for _, cat := range c.Items {
		if cat.Name == &name {
			category = &cat
			break
		}
	}
	return category, nil

}
func (c *categoryMockRepo) GetAll(pageIndex, pageSize int) (*[]models.Category, int, error) {
	return &c.Items, 1, nil
}
