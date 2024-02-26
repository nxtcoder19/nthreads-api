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
	InsertProduct(ctx context.Context, title string, price int) (*Product, error)
	//GetAllProducts(ctx context.Context) ([]*Product, error)
	SignUp(ctx context.Context, firstName string, lastName string, email string, password string) (*entities.User, error)
	Login(ctx context.Context, email string, password string) (res string, err error)
	UpdateUser(ctx context.Context, email string, firstName string, lastName string) (*entities.User, error)
	DeleteUser(ctx context.Context, email string) (string, error)
	CreateTodo(ctx context.Context, title string, description string) (*entities.Todo, error)
	UpdateTodo(ctx context.Context, id string, title string, description string) (*entities.Todo, error)
	GetTodo(ctx context.Context, id string) (*entities.Todo, error)
	GetTodos(ctx context.Context) ([]*entities.Todo, error)
	DeleteTodo(ctx context.Context, id string) (string, error)
}
