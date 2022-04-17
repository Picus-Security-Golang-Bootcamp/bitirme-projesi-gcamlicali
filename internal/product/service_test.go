package product

import (
	"github.com/gcamlicali/tradeshopExample/internal/api"
	"github.com/gcamlicali/tradeshopExample/internal/category"
	"github.com/gcamlicali/tradeshopExample/internal/models"
	"github.com/go-openapi/errors"
	"github.com/google/uuid"
	"log"
	"strings"
	"testing"
)

var (
	categoryName = "CategoryName1"
	NExCatName   = "CategoryNameNotExisted"
	productName  = "productExample1"
	NExProName   = "ProductNameNotExisted"
	description  = "TestProduct"
	price        = 1000
	sku          = 1
	NExSku       = 99999
	sku64        = int64(sku)
	price32      = int32(price)
	unitStock    = int32(1000)

	apiProductName = "ApiProductName"
	apiSKU         = int64(2)
)

var product1 = models.Product{
	ID:           uuid.New(),
	CategoryName: categoryName,
	Name:         productName,
	Price:        price,
	SKU:          sku,
	Description:  description,
	UnitStock:    unitStock,
}

func Test_productService_AddSingle(t *testing.T) {

	type fields struct {
		pRepo   IProductRepository
		catRepo category.ICategoryRepository
	}
	type args struct {
		product api.Product
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Product
		wantErr bool
	}{
		{
			name: "productService_ProductAddSingle_ShouldSuccess",
			fields: fields{
				catRepo: &categoryMockRepo{
					Items: []models.Category{
						{
							ID:   uuid.New(),
							Name: &categoryName,
						},
					},
				},
				pRepo: &productMockRepo{
					Items: []models.Product{},
				},
			},
			args: args{
				product: api.Product{
					Name:         &apiProductName,
					CategoryName: &categoryName,
					Price:        &price32,
					UnitStock:    &unitStock,
					Description:  description,
					Sku:          &apiSKU,
				},
			},
			wantErr: false,
		},
		{
			name: "productService_ProductAddSingle_ErrorCategoryNotFound_ShouldFail",
			fields: fields{
				catRepo: &categoryMockRepo{
					Items: []models.Category{
						{
							ID:   uuid.New(),
							Name: &categoryName,
						},
					},
				},
				pRepo: &productMockRepo{
					Items: []models.Product{},
				},
			},
			args: args{
				product: api.Product{
					Name:         &apiProductName,
					CategoryName: &NExCatName,
					Price:        &price32,
					UnitStock:    &unitStock,
					Description:  description,
					Sku:          &apiSKU,
				},
			},
			wantErr: true,
		},
		{
			name: "productService_ProductAddSingle_ErrorProductSKUExist_ShouldFail",
			fields: fields{
				catRepo: &categoryMockRepo{
					Items: []models.Category{
						{
							ID:   uuid.New(),
							Name: &categoryName,
						},
					},
				},
				pRepo: &productMockRepo{
					Items: []models.Product{
						{
							Name:         product1.Name,
							Description:  product1.Description,
							CategoryName: product1.CategoryName,
							SKU:          product1.SKU,
							Price:        product1.Price,
							UnitStock:    product1.UnitStock,
							ID:           product1.ID,
						},
					},
				},
			},
			args: args{
				product: api.Product{
					Name:         &apiProductName,
					CategoryName: &categoryName,
					Price:        &price32,
					UnitStock:    &unitStock,
					Description:  description,
					Sku:          &sku64,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := productService{
				pRepo:   tt.fields.pRepo,
				catRepo: tt.fields.catRepo,
			}
			_, err := p.AddSingle(tt.args.product)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddSingle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_productService_GetAll(t *testing.T) {
	type fields struct {
		pRepo   IProductRepository
		catRepo category.ICategoryRepository
	}
	type args struct {
		pageIndex int
		pageSize  int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		//want    *[]models.Product
		//want1   int
		wantErr bool
	}{
		{
			name: "productService_ProductGetAll_ShouldSuccess",
			fields: fields{
				catRepo: &categoryMockRepo{
					Items: []models.Category{
						{
							ID:   uuid.New(),
							Name: &categoryName,
						},
					},
				},
				pRepo: &productMockRepo{
					Items: []models.Product{},
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
			p := productService{
				pRepo:   tt.fields.pRepo,
				catRepo: tt.fields.catRepo,
			}
			_, _, err := p.GetAll(tt.args.pageIndex, tt.args.pageSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_productService_Delete(t *testing.T) {
	type fields struct {
		pRepo   IProductRepository
		catRepo category.ICategoryRepository
	}
	type args struct {
		SKU int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "productService_ProductDelete_ShouldSuccess",
			fields: fields{
				catRepo: &categoryMockRepo{
					Items: []models.Category{
						{
							ID:   uuid.New(),
							Name: &categoryName,
						},
					},
				},
				pRepo: &productMockRepo{
					Items: []models.Product{
						{
							Name:         product1.Name,
							Description:  product1.Description,
							CategoryName: product1.CategoryName,
							SKU:          product1.SKU,
							Price:        product1.Price,
							UnitStock:    product1.UnitStock,
							ID:           product1.ID,
						},
					},
				},
			},
			args: args{
				SKU: sku,
			},
			wantErr: false,
		},
		{
			name: "productService_ProductDelete_ErrorProductNotFound_ShouldFail",
			fields: fields{
				catRepo: &categoryMockRepo{
					Items: []models.Category{
						{
							ID:   uuid.New(),
							Name: &categoryName,
						},
					},
				},
				pRepo: &productMockRepo{
					Items: []models.Product{
						{
							Name:         product1.Name,
							Description:  product1.Description,
							CategoryName: product1.CategoryName,
							SKU:          product1.SKU,
							Price:        product1.Price,
							UnitStock:    product1.UnitStock,
							ID:           product1.ID,
						},
					},
				},
			},
			args: args{
				SKU: NExSku,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := productService{
				pRepo:   tt.fields.pRepo,
				catRepo: tt.fields.catRepo,
			}
			if err := p.Delete(tt.args.SKU); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_productService_Update(t *testing.T) {
	type fields struct {
		pRepo   IProductRepository
		catRepo category.ICategoryRepository
	}
	type args struct {
		SKU        int
		reqProduct *api.ProductUp
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "productService_ProductUpdate_ShouldSuccess",
			fields: fields{
				catRepo: &categoryMockRepo{
					Items: []models.Category{
						{
							ID:   uuid.New(),
							Name: &categoryName,
						},
					},
				},
				pRepo: &productMockRepo{
					Items: []models.Product{
						{
							Name:         product1.Name,
							Description:  product1.Description,
							CategoryName: product1.CategoryName,
							SKU:          product1.SKU,
							Price:        product1.Price,
							UnitStock:    product1.UnitStock,
							ID:           product1.ID,
						},
					},
				},
			},
			args: args{
				SKU: sku,
				reqProduct: &api.ProductUp{
					Name:         apiProductName,
					CategoryName: categoryName,
					Price:        price32,
					UnitStock:    unitStock,
					Description:  description,
					Sku:          int64(NExSku),
				},
			},
			wantErr: false,
		},
		{
			name: "productService_ProductUpdate_ErrorProductNotFound_ShouldFail",
			fields: fields{
				catRepo: &categoryMockRepo{
					Items: []models.Category{
						{
							ID:   uuid.New(),
							Name: &categoryName,
						},
					},
				},
				pRepo: &productMockRepo{
					Items: []models.Product{
						{
							Name:         product1.Name,
							Description:  product1.Description,
							CategoryName: product1.CategoryName,
							SKU:          product1.SKU,
							Price:        product1.Price,
							UnitStock:    product1.UnitStock,
							ID:           product1.ID,
						},
					},
				},
			},
			args: args{
				SKU: NExSku,
				reqProduct: &api.ProductUp{
					Name:         apiProductName,
					CategoryName: categoryName,
					Price:        price32,
					UnitStock:    unitStock,
					Description:  description,
					Sku:          int64(NExSku),
				},
			},
			wantErr: true,
		},
		{
			name: "productService_ProductUpdate_ErrorCategoryNotFound_ShouldFail",
			fields: fields{
				catRepo: &categoryMockRepo{
					Items: []models.Category{
						{
							ID:   uuid.New(),
							Name: &categoryName,
						},
					},
				},
				pRepo: &productMockRepo{
					Items: []models.Product{
						{
							Name:         product1.Name,
							Description:  product1.Description,
							CategoryName: product1.CategoryName,
							SKU:          product1.SKU,
							Price:        product1.Price,
							UnitStock:    product1.UnitStock,
							ID:           product1.ID,
						},
					},
				},
			},
			args: args{
				SKU: sku,
				reqProduct: &api.ProductUp{
					Name:         apiProductName,
					CategoryName: NExCatName,
					Price:        price32,
					UnitStock:    unitStock,
					Description:  description,
					Sku:          int64(NExSku),
				},
			},
			wantErr: true,
		},
		{
			name: "productService_ProductUpdate_ErrorExistedSKU_ShouldFail",
			fields: fields{
				catRepo: &categoryMockRepo{
					Items: []models.Category{
						{
							ID:   uuid.New(),
							Name: &categoryName,
						},
					},
				},
				pRepo: &productMockRepo{
					Items: []models.Product{
						{
							Name:         product1.Name,
							Description:  product1.Description,
							CategoryName: product1.CategoryName,
							SKU:          product1.SKU,
							Price:        product1.Price,
							UnitStock:    product1.UnitStock,
							ID:           product1.ID,
						},
					},
				},
			},
			args: args{
				SKU: sku,
				reqProduct: &api.ProductUp{
					Name:         apiProductName,
					CategoryName: categoryName,
					Price:        price32,
					UnitStock:    unitStock,
					Description:  description,
					Sku:          int64(sku),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := productService{
				pRepo:   tt.fields.pRepo,
				catRepo: tt.fields.catRepo,
			}
			_, err := p.Update(tt.args.SKU, tt.args.reqProduct)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_productService_GetByName(t *testing.T) {
	type fields struct {
		pRepo   IProductRepository
		catRepo category.ICategoryRepository
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		//want    *[]models.Product
		wantErr bool
	}{
		{
			name: "productService_ProductGetByName_ShouldSuccess",
			fields: fields{
				catRepo: &categoryMockRepo{
					Items: []models.Category{
						{
							ID:   uuid.New(),
							Name: &categoryName,
						},
					},
				},
				pRepo: &productMockRepo{
					Items: []models.Product{
						{
							Name:         product1.Name,
							Description:  product1.Description,
							CategoryName: product1.CategoryName,
							SKU:          product1.SKU,
							Price:        product1.Price,
							UnitStock:    product1.UnitStock,
							ID:           product1.ID,
						},
					},
				},
			},
			args: args{
				name: productName,
			},
			wantErr: false,
		},
		{
			name: "productService_ProductGetByName_ErrorProductNotFound_ShouldFail",
			fields: fields{
				catRepo: &categoryMockRepo{
					Items: []models.Category{
						{
							ID:   uuid.New(),
							Name: &categoryName,
						},
					},
				},
				pRepo: &productMockRepo{
					Items: []models.Product{
						{
							Name:         product1.Name,
							Description:  product1.Description,
							CategoryName: product1.CategoryName,
							SKU:          product1.SKU,
							Price:        product1.Price,
							UnitStock:    product1.UnitStock,
							ID:           product1.ID,
						},
					},
				},
			},
			args: args{
				name: NExProName,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := productService{
				pRepo:   tt.fields.pRepo,
				catRepo: tt.fields.catRepo,
			}
			_, err := p.GetByName(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_productService_GetBySKU(t *testing.T) {
	type fields struct {
		pRepo   IProductRepository
		catRepo category.ICategoryRepository
	}
	type args struct {
		SKU int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		//want    *models.Product
		wantErr bool
	}{
		{
			name: "productService_ProductGetBySKU_ShouldSuccess",
			fields: fields{
				catRepo: &categoryMockRepo{
					Items: []models.Category{
						{
							ID:   uuid.New(),
							Name: &categoryName,
						},
					},
				},
				pRepo: &productMockRepo{
					Items: []models.Product{
						{
							Name:         product1.Name,
							Description:  product1.Description,
							CategoryName: product1.CategoryName,
							SKU:          product1.SKU,
							Price:        product1.Price,
							UnitStock:    product1.UnitStock,
							ID:           product1.ID,
						},
					},
				},
			},
			args: args{
				SKU: sku,
			},
			wantErr: false,
		},
		{
			name: "productService_ProductGetBySKU_ErrorProductNotFound_ShouldFail",
			fields: fields{
				catRepo: &categoryMockRepo{
					Items: []models.Category{
						{
							ID:   uuid.New(),
							Name: &categoryName,
						},
					},
				},
				pRepo: &productMockRepo{
					Items: []models.Product{
						{
							Name:         product1.Name,
							Description:  product1.Description,
							CategoryName: product1.CategoryName,
							SKU:          product1.SKU,
							Price:        product1.Price,
							UnitStock:    product1.UnitStock,
							ID:           product1.ID,
						},
					},
				},
			},
			args: args{
				SKU: NExSku,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := productService{
				pRepo:   tt.fields.pRepo,
				catRepo: tt.fields.catRepo,
			}
			_, err := p.GetBySKU(tt.args.SKU)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBySKU() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

type categoryMockRepo struct {
	Items []models.Category
}
type productMockRepo struct {
	Items []models.Product
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
		if *cat.Name == name {
			category = &cat
			return category, nil
		}
	}

	return nil, errors.New(404, "category not found")
}
func (c *categoryMockRepo) GetAll(pageIndex, pageSize int) (*[]models.Category, int, error) {

	return &c.Items, len(c.Items), nil
}

func (p *productMockRepo) Create(a *models.Product) (*models.Product, error) {
	for _, item := range p.Items {
		if item.SKU == a.SKU {
			return nil, errors.New(400, "Item should be unique on database")
		}
	}
	p.Items = append(p.Items, *a)
	return a, nil
}
func (p *productMockRepo) GetAll(pageIndex, pageSize int) (*[]models.Product, int, error) {
	log.Println("size: ", len(p.Items))
	return &p.Items, len(p.Items), nil
}
func (p *productMockRepo) GetByName(name string) (*[]models.Product, error) {
	products := []models.Product{}
	for _, item := range p.Items {
		if strings.Contains(item.Name, name) {
			products = append(products, item)
		}
	}
	if len(products) > 0 {
		return &products, nil
	} else {
		return nil, errors.New(400, "Product not found")
	}
}
func (p *productMockRepo) GetBySKU(SKU int) (*models.Product, error) {
	product := models.Product{}
	for _, item := range p.Items {
		if item.SKU == SKU {
			product = item
			return &product, nil
		}
	}

	return nil, errors.New(400, "Product not found")
}
func (p *productMockRepo) GetByCatName(catName string) (*[]models.Product, error) {
	products := []models.Product{}
	for _, item := range p.Items {
		if item.CategoryName == catName {
			products = append(products, item)
		}
	}
	if len(products) > 0 {
		return &products, nil
	} else {
		return nil, errors.New(400, "Product not found")
	}
}
func (p *productMockRepo) Update(a *models.Product) (*models.Product, error) {
	for i, item := range p.Items {
		if item.SKU == a.SKU {
			p.Items[i] = *a
			break
		}
	}
	return a, nil
}
func (p *productMockRepo) Delete(sku int) error {
	pro, err := p.GetBySKU(sku)
	if err != nil {
		return errors.New(400, "Product not found")
	}
	for i, item := range p.Items {
		if item.SKU == pro.SKU {
			p.Items = append(p.Items[:i], p.Items[i+1:]...)
			break
		}
	}
	return nil
}
