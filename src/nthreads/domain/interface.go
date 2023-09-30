package nthreads

import (
	"context"
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
	GetAllProducts(ctx context.Context) ([]*Product, error)
}
