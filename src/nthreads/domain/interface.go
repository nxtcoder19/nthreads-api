package domain

import (
	"context"
	"github.com/nxtcoder19/nthreads-backend/src/nthreads/entities"
	"time"
)

type Id string
type Product struct {
	Id    Id        `json:"id"`
	Title string    `json:"title"`
	Price int       `json:"price"`
	Time  time.Time `json:"time"`
}

type NThreads interface {
	Init(ctx context.Context) error
	//InsertProduct(ctx context.Context, title string, price int) (*Product, error)
	//GetAllProducts(ctx context.Context) ([]*Product, error)

	/// Auth
	SignUp(ctx context.Context, firstName string, lastName string, email string, password string) (*entities.User, error)
	Login(ctx context.Context, email string, password string) (res bool, err error)
	UpdateUser(ctx context.Context, email string, firstName string, lastName string) (*entities.User, error)
	DeleteUser(ctx context.Context, email string) (string, error)
	GetUser(ctx context.Context, email string) (*entities.User, error)

	/// Todo
	CreateTodo(ctx context.Context, title string, description string) (*entities.Todo, error)
	UpdateTodo(ctx context.Context, id string, title string, description string) (*entities.Todo, error)
	GetTodo(ctx context.Context, id string) (*entities.Todo, error)
	GetTodos(ctx context.Context) ([]*entities.Todo, error)
	DeleteTodo(ctx context.Context, id string) (string, error)

	/// Product
	CreateProduct(ctx context.Context, productIn *entities.Product) (*entities.Product, error)
	//UpdateProduct(ctx context.Context, id string, price string, imageUrl string, date string, warranty string, place string) (*entities.Product, error)
	UpdateProduct(ctx context.Context, id string, productIn *entities.Product) (*entities.Product, error)
	GetProduct(ctx context.Context, id string) (*entities.Product, error)
	GetProducts(ctx context.Context) ([]*entities.Product, error)
	DeleteProduct(ctx context.Context, id string) (bool, error)

	/// Product Category
	CreateProductCategory(ctx context.Context, name string, title string, description string, imageUrl string) (*entities.ProductCategory, error)
	UpdateProductCategory(ctx context.Context, id string, name string, title string, description string, imageUrl string) (*entities.ProductCategory, error)
	GetProductCategory(ctx context.Context, id string) (*entities.ProductCategory, error)
	GetProductCategories(ctx context.Context) ([]*entities.ProductCategory, error)
	DeleteProductCategory(ctx context.Context, id string) (bool, error)

	/// Cart Item
	AddItemToCart(ctx context.Context, product *entities.Product, quantity int) (*entities.CartItem, error)
	UpdateCartItem(ctx context.Context, id string, quantity int) (*entities.CartItem, error)
	RemoveCartItem(ctx context.Context, id string) (bool, error)
	GetCartItems(ctx context.Context) ([]*entities.CartItem, error)

	/// Order Item
	AddItemToOrder(ctx context.Context, product *entities.Product) (*entities.OrderItem, error)
	GetOrderItem(ctx context.Context, id string) (*entities.OrderItem, error)
	GetOrderItems(ctx context.Context) ([]*entities.OrderItem, error)

	/// Address
	CreateAddress(ctx context.Context, addressIn *entities.Address) (*entities.Address, error)
	UpdateAddress(ctx context.Context, id string, addressIn *entities.Address) (*entities.Address, error)
	GetAddress(ctx context.Context, id string) (*entities.Address, error)
	GetAddresses(ctx context.Context) ([]*entities.Address, error)
	DeleteAddress(ctx context.Context, id string) (bool, error)
}
