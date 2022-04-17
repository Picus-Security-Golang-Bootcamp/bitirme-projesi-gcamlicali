package cart

import (
	"github.com/gcamlicali/tradeshopExample/internal/cart_item"
	"github.com/gcamlicali/tradeshopExample/internal/models"
	"github.com/gcamlicali/tradeshopExample/internal/product"
	"github.com/go-openapi/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"reflect"
	"strings"
	"testing"
)

var (
	product1ID = uuid.New()
	userID     = uuid.New()
	cartID     = uuid.New()
	cartItemID = uuid.New()
	NExUser    = uuid.New()
	NExProSKU  = 9999999

	product1 = models.Product{
		ID:           product1ID,
		CategoryName: "CategoryExample",
		Name:         "Product1",
		Price:        10,
		SKU:          1,
		Description:  "ExampleProduct1",
		UnitStock:    1,
	}
	product2 = models.Product{
		ID:           product1ID,
		CategoryName: "CategoryExample",
		Name:         "Product2",
		Price:        10,
		SKU:          2,
		Description:  "ExampleProduct2",
		UnitStock:    1,
	}
	cartItem1 = models.CartItem{
		ID:         cartItemID,
		CartID:     cartID,
		ProductSKU: product1.SKU,
		Price:      product1.Price,
		Product:    product1,
		Quantity:   1,
	}
	cartItem2 = models.CartItem{
		ID:         cartItemID,
		CartID:     cartID,
		ProductSKU: product2.SKU,
		Price:      product2.Price,
		Product:    product2,
		Quantity:   1,
	}
	cart1 = models.Cart{
		ID:         cartID,
		UserID:     userID,
		TotalPrice: 10,
		IsOrdered:  false,
	}
	cart1updated = models.Cart{
		ID:         cartID,
		UserID:     userID,
		TotalPrice: 30,
		IsOrdered:  false,
	}
)

func Test_cartService_Get(t *testing.T) {
	type fields struct {
		crepo  ICartRepository
		cirepo cart_item.ICartItemRepository
		prepo  product.IProductRepository
	}
	type args struct {
		userID uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Cart
		wantErr bool
	}{
		{
			name: "cartService_cartGet_ShouldSuccess",
			fields: fields{
				prepo: &productMockRepo{
					Items: []models.Product{
						product1,
					},
				},
				cirepo: &cartItemMockRepo{
					Items: []models.CartItem{
						cartItem1,
					},
				},
				crepo: &cartMockRepo{
					Items: []models.Cart{
						cart1,
					},
				},
			},
			args: args{
				userID: userID,
			},
			want:    &cart1,
			wantErr: false,
		},
		{
			name: "cartService_cartGet_ErrorUserNotExisted_ShouldFail",
			fields: fields{
				prepo: &productMockRepo{
					Items: []models.Product{
						product1,
					},
				},
				cirepo: &cartItemMockRepo{
					Items: []models.CartItem{
						cartItem1,
					},
				},
				crepo: &cartMockRepo{
					Items: []models.Cart{
						cart1,
					},
				},
			},
			args: args{
				userID: NExUser,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cartService{
				crepo:  tt.fields.crepo,
				cirepo: tt.fields.cirepo,
				prepo:  tt.fields.prepo,
			}
			got, err := c.Get(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cartService_Add(t *testing.T) {
	type fields struct {
		crepo  ICartRepository
		cirepo cart_item.ICartItemRepository
		prepo  product.IProductRepository
	}
	type args struct {
		userID     uuid.UUID
		ProductSKU int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		//want    *models.Cart
		wantErr bool
	}{
		{
			name: "cartService_cartAdd_ShouldSuccess",
			fields: fields{
				prepo: &productMockRepo{
					Items: []models.Product{
						product1,
					},
				},
				cirepo: &cartItemMockRepo{
					Items: []models.CartItem{},
				},
				crepo: &cartMockRepo{
					Items: []models.Cart{
						cart1,
					},
				},
			},
			args: args{
				userID:     userID,
				ProductSKU: product1.SKU,
			},
			wantErr: false,
		},
		{
			name: "cartService_cartAdd_ErrorUserNotFound_ShouldFail",
			fields: fields{
				prepo: &productMockRepo{
					Items: []models.Product{
						product1,
					},
				},
				cirepo: &cartItemMockRepo{
					Items: []models.CartItem{},
				},
				crepo: &cartMockRepo{
					Items: []models.Cart{
						cart1,
					},
				},
			},
			args: args{
				userID:     NExUser,
				ProductSKU: product1.SKU,
			},
			wantErr: true,
		},
		{
			name: "cartService_cartAdd_ErrorProductNotFound_ShouldFail",
			fields: fields{
				prepo: &productMockRepo{
					Items: []models.Product{
						product1,
					},
				},
				cirepo: &cartItemMockRepo{
					Items: []models.CartItem{},
				},
				crepo: &cartMockRepo{
					Items: []models.Cart{
						cart1,
					},
				},
			},
			args: args{
				userID:     userID,
				ProductSKU: NExProSKU,
			},
			wantErr: true,
		},
		{
			name: "cartService_cartAdd_ErrorProductNotFound_ShouldFail",
			fields: fields{
				prepo: &productMockRepo{
					Items: []models.Product{
						product1,
					},
				},
				cirepo: &cartItemMockRepo{
					Items: []models.CartItem{},
				},
				crepo: &cartMockRepo{
					Items: []models.Cart{
						cart1,
					},
				},
			},
			args: args{
				userID:     userID,
				ProductSKU: NExProSKU,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cartService{
				crepo:  tt.fields.crepo,
				cirepo: tt.fields.cirepo,
				prepo:  tt.fields.prepo,
			}
			_, err := c.Add(tt.args.userID, tt.args.ProductSKU)
			if (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_cartService_Update(t *testing.T) {
	type fields struct {
		crepo  ICartRepository
		cirepo cart_item.ICartItemRepository
		prepo  product.IProductRepository
	}
	type args struct {
		userID     uuid.UUID
		ProductSKU int
		Quantity   int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Cart
		wantErr bool
	}{
		{
			name: "cartService_cartUpdate_ShouldSuccess",
			fields: fields{
				prepo: &productMockRepo{
					Items: []models.Product{
						product1,
					},
				},
				cirepo: &cartItemMockRepo{
					Items: []models.CartItem{
						cartItem1,
					},
				},
				crepo: &cartMockRepo{
					Items: []models.Cart{
						cart1,
					},
				},
			},
			args: args{
				userID:     userID,
				ProductSKU: product1.SKU,
				Quantity:   3,
			},
			want:    &cart1updated,
			wantErr: false,
		},
		{
			name: "cartService_cartUpdate_ErrorUserNotFound_ShouldFail",
			fields: fields{
				prepo: &productMockRepo{
					Items: []models.Product{
						product1,
					},
				},
				cirepo: &cartItemMockRepo{
					Items: []models.CartItem{
						cartItem1,
					},
				},
				crepo: &cartMockRepo{
					Items: []models.Cart{
						cart1,
					},
				},
			},
			args: args{
				userID:     NExUser,
				ProductSKU: product1.SKU,
				Quantity:   3,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "cartService_cartUpdate_ErrorSKUNotFound_ShouldFail",
			fields: fields{
				prepo: &productMockRepo{
					Items: []models.Product{
						product1,
					},
				},
				cirepo: &cartItemMockRepo{
					Items: []models.CartItem{
						cartItem1,
					},
				},
				crepo: &cartMockRepo{
					Items: []models.Cart{
						cart1,
					},
				},
			},
			args: args{
				userID:     userID,
				ProductSKU: NExProSKU,
				Quantity:   3,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cartService{
				crepo:  tt.fields.crepo,
				cirepo: tt.fields.cirepo,
				prepo:  tt.fields.prepo,
			}
			got, err := c.Update(tt.args.userID, tt.args.ProductSKU, tt.args.Quantity)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Update() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cartService_Delete(t *testing.T) {
	type fields struct {
		crepo  ICartRepository
		cirepo cart_item.ICartItemRepository
		prepo  product.IProductRepository
	}
	type args struct {
		userID     uuid.UUID
		ProductSKU int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Cart
		wantErr bool
	}{
		{
			name: "cartService_cartDelete_ShouldSuccess",
			fields: fields{
				prepo: &productMockRepo{
					Items: []models.Product{
						product1,
						product2,
					},
				},
				cirepo: &cartItemMockRepo{
					Items: []models.CartItem{
						cartItem1,
						cartItem2,
					},
				},
				crepo: &cartMockRepo{
					Items: []models.Cart{
						cart1,
					},
				},
			},
			args: args{
				userID:     userID,
				ProductSKU: product2.SKU,
			},
			want:    &cart1,
			wantErr: false,
		},
		{
			name: "cartService_cartDelete_ErrorUserNotFound_ShouldFail",
			fields: fields{
				prepo: &productMockRepo{
					Items: []models.Product{
						product1,
						product2,
					},
				},
				cirepo: &cartItemMockRepo{
					Items: []models.CartItem{
						cartItem1,
						cartItem2,
					},
				},
				crepo: &cartMockRepo{
					Items: []models.Cart{
						cart1,
					},
				},
			},
			args: args{
				userID:     NExUser,
				ProductSKU: product2.SKU,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "cartService_cartDelete_ErrorUserNotFound_ShouldFail",
			fields: fields{
				prepo: &productMockRepo{
					Items: []models.Product{
						product1,
						product2,
					},
				},
				cirepo: &cartItemMockRepo{
					Items: []models.CartItem{
						cartItem1,
						cartItem2,
					},
				},
				crepo: &cartMockRepo{
					Items: []models.Cart{
						cart1,
					},
				},
			},
			args: args{
				userID:     userID,
				ProductSKU: NExProSKU,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cartService{
				crepo:  tt.fields.crepo,
				cirepo: tt.fields.cirepo,
				prepo:  tt.fields.prepo,
			}
			got, err := c.Delete(tt.args.userID, tt.args.ProductSKU)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delete() got = %v, want %v", got, tt.want)
			}
		})
	}
}

type productMockRepo struct {
	Items []models.Product
}
type cartItemMockRepo struct {
	Items []models.CartItem
}
type cartMockRepo struct {
	Items []models.Cart
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

func (ci *cartItemMockRepo) Crate(a *models.CartItem) (*models.CartItem, error) {
	ci.Items = append(ci.Items, *a)
	return a, nil
}
func (ci *cartItemMockRepo) GetByCartID(cartID uuid.UUID) (*[]models.CartItem, error) {
	cartItems := []models.CartItem{}
	for _, item := range ci.Items {
		if item.CartID == cartID {
			cartItems = append(cartItems, item)
		}
	}
	return &cartItems, nil
}
func (ci *cartItemMockRepo) GetByCartAndProductSKU(cartID uuid.UUID, productSKU int) (*models.CartItem, error) {
	cartItem := models.CartItem{}
	for i, item := range ci.Items {
		if item.CartID == cartID {
			if item.ProductSKU == productSKU {
				cartItem = ci.Items[i]
				return &cartItem, nil
			}
		}
	}
	return nil, gorm.ErrRecordNotFound
}
func (ci *cartItemMockRepo) Update(a *models.CartItem) (*models.CartItem, error) {
	for i, item := range ci.Items {
		if item.ProductSKU == a.ProductSKU {
			ci.Items[i] = *a
			return &ci.Items[i], nil
		}
	}
	return nil, errors.New(400, "Cart not found")
}
func (ci *cartItemMockRepo) Delete(a *models.CartItem) error {

	cartItem, err := ci.GetByCartAndProductSKU(a.CartID, a.ProductSKU)
	if err != nil {
		return errors.New(400, "Product not found")
	}

	for i, item := range ci.Items {
		if item.ProductSKU == cartItem.ProductSKU {
			ci.Items = append(ci.Items[:i], ci.Items[i+1:]...)
			break
		}
	}
	return nil
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
