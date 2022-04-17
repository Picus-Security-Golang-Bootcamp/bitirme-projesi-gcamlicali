package order

import (
	"github.com/gcamlicali/tradeshopExample/internal/cart"
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
	"time"
)

var (
	product1ID  = uuid.New()
	userID      = uuid.New()
	cartID      = uuid.New()
	cartItemID  = uuid.New()
	orderID     = uuid.New()
	NExUser     = uuid.New()
	NExOrder    = uuid.New()
	NExProSKU   = 9999999
	currentTime = time.Now()

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
		Quantity:   10,
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
	order1 = models.Order{
		ID:         orderID,
		UserID:     userID,
		Cart:       cart1,
		CartID:     cartID,
		TotalPrice: int32(cart1.TotalPrice),
		Status:     "Ordered",
		CreatedAt: time.Date(
			currentTime.Year(),
			currentTime.Month(),
			currentTime.Day()-1,
			currentTime.Hour(),
			currentTime.Minute(),
			currentTime.Second(),
			0,
			time.Local),
	}
)

func Test_orderService_GetAll(t *testing.T) {
	type fields struct {
		orRepo IOrderRepository
		cRepo  cart.ICartRepository
		ciRepo cart_item.ICartItemRepository
		pRepo  product.IProductRepository
	}
	type args struct {
		userID uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *[]models.Order
		wantErr bool
	}{
		{
			name: "orderService_OrderGetAll_ShouldSuccess",
			fields: fields{
				orRepo: &orderMockRepo{
					Items: []models.Order{
						order1,
					},
				},
				pRepo: &productMockRepo{
					Items: []models.Product{
						product1,
					},
				},
				ciRepo: &cartItemMockRepo{
					Items: []models.CartItem{
						cartItem1,
					},
				},
				cRepo: &cartMockRepo{
					Items: []models.Cart{
						cart1,
					},
				},
			},
			args: args{
				userID: userID,
			},
			want: &[]models.Order{
				order1,
			},
			wantErr: false,
		},
		{
			name: "orderService_OrderGetAll_ErrorUserNotFound_ShouldFail",
			fields: fields{
				orRepo: &orderMockRepo{
					Items: []models.Order{
						order1,
					},
				},
				pRepo: &productMockRepo{
					Items: []models.Product{
						product1,
					},
				},
				ciRepo: &cartItemMockRepo{
					Items: []models.CartItem{
						cartItem1,
					},
				},
				cRepo: &cartMockRepo{
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
			c := &orderService{
				orRepo: tt.fields.orRepo,
				cRepo:  tt.fields.cRepo,
				ciRepo: tt.fields.ciRepo,
				pRepo:  tt.fields.pRepo,
			}
			got, err := c.GetAll(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAll() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_orderService_Create(t *testing.T) {
	type fields struct {
		orRepo IOrderRepository
		cRepo  cart.ICartRepository
		ciRepo cart_item.ICartItemRepository
		pRepo  product.IProductRepository
	}
	type args struct {
		userID uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Order
		wantErr bool
	}{
		{
			name: "orderService_OrderCreate_ShouldSuccess",
			fields: fields{
				orRepo: &orderMockRepo{
					Items: []models.Order{},
				},
				pRepo: &productMockRepo{
					Items: []models.Product{
						product1,
					},
				},
				ciRepo: &cartItemMockRepo{
					Items: []models.CartItem{
						cartItem1,
					},
				},
				cRepo: &cartMockRepo{
					Items: []models.Cart{
						cart1,
					},
				},
			},
			args: args{
				userID: userID,
			},
			want:    &order1,
			wantErr: false,
		},
		{
			name: "orderService_OrderCreate_ErrorUserNotFound_ShouldFail",
			fields: fields{
				orRepo: &orderMockRepo{
					Items: []models.Order{},
				},
				pRepo: &productMockRepo{
					Items: []models.Product{
						product1,
					},
				},
				ciRepo: &cartItemMockRepo{
					Items: []models.CartItem{
						cartItem1,
					},
				},
				cRepo: &cartMockRepo{
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
		{
			name: "orderService_OrderCreate_ErrorNotEnoughStock_ShouldFail",
			fields: fields{
				orRepo: &orderMockRepo{
					Items: []models.Order{},
				},
				pRepo: &productMockRepo{
					Items: []models.Product{
						product2,
					},
				},
				ciRepo: &cartItemMockRepo{
					Items: []models.CartItem{
						cartItem2,
					},
				},
				cRepo: &cartMockRepo{
					Items: []models.Cart{
						cart1,
					},
				},
			},
			args: args{
				userID: userID,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &orderService{
				orRepo: tt.fields.orRepo,
				cRepo:  tt.fields.cRepo,
				ciRepo: tt.fields.ciRepo,
				pRepo:  tt.fields.pRepo,
			}
			got, err := c.Create(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_orderService_Cancel(t *testing.T) {
	type fields struct {
		orRepo IOrderRepository
		cRepo  cart.ICartRepository
		ciRepo cart_item.ICartItemRepository
		pRepo  product.IProductRepository
	}
	type args struct {
		userID  uuid.UUID
		orderID uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "orderService_OrderCancel_ShouldSuccess",
			fields: fields{
				orRepo: &orderMockRepo{
					Items: []models.Order{
						order1,
					},
				},
				pRepo: &productMockRepo{
					Items: []models.Product{
						product1,
					},
				},
				ciRepo: &cartItemMockRepo{
					Items: []models.CartItem{
						cartItem1,
					},
				},
				cRepo: &cartMockRepo{
					Items: []models.Cart{
						cart1,
					},
				},
			},
			args: args{
				userID:  userID,
				orderID: orderID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &orderService{
				orRepo: tt.fields.orRepo,
				cRepo:  tt.fields.cRepo,
				ciRepo: tt.fields.ciRepo,
				pRepo:  tt.fields.pRepo,
			}
			if err := c.Cancel(tt.args.userID, tt.args.orderID); (err != nil) != tt.wantErr {
				t.Errorf("Cancel() error = %v, wantErr %v", err, tt.wantErr)
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
type orderMockRepo struct {
	Items []models.Order
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

func (o *orderMockRepo) Create(a *models.Order) (*models.Order, error) {
	o.Items = append(o.Items, *a)
	return a, nil
}
func (o *orderMockRepo) GetByOrderAndUserID(userID uuid.UUID, orderID uuid.UUID) (*models.Order, error) {
	order := models.Order{}
	for i, item := range o.Items {
		if item.UserID == userID {
			if item.ID == orderID {
				order = o.Items[i]
				return &order, nil
			}
		}
	}
	return nil, gorm.ErrRecordNotFound
}
func (o *orderMockRepo) GetByUserID(userID uuid.UUID) (*[]models.Order, error) {
	orders := []models.Order{}
	for _, item := range o.Items {
		if item.UserID == userID {
			orders = append(orders, item)
		}
	}
	if len(orders) > 0 {
		return &orders, nil
	} else {
		return nil, errors.New(400, "User Orders not found")
	}
}
func (o *orderMockRepo) Update(a *models.Order) (*models.Order, error) {
	for i, item := range o.Items {
		if item.ID == a.ID {
			o.Items[i] = *a
			return a, nil
		}
	}
	return nil, errors.New(400, "Order not found")
}
