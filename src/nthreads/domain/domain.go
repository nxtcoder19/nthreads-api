package domain

import (
	"context"
	mongodb "github.com/nxtcoder19/nthreads-backend/package/mongo-db"
	"github.com/nxtcoder19/nthreads-backend/src/nthreads/entities"
)

type Impl struct {
	db mongodb.DBRepo
}

const (
	ProductTable         string = "products"
	AuthTable            string = "auth"
	TodoTable            string = "todo"
	ProductCategoryTable string = "product_category"
	CartTable            string = "cart_items"
	OrderTable           string = "orders"
	AddressTable         string = "Address"
)

func getSessionData(ctx context.Context) *entities.SessionData {
	return ctx.Value("user-session").(*entities.SessionData)
}

func (i *Impl) Init(ctx context.Context) error {
	//err := i.db.CreateCollection(ctx, "products")
	err := i.db.CreateCollection(ctx, AuthTable)
	err = i.db.CreateCollection(ctx, TodoTable)
	err = i.db.CreateCollection(ctx, ProductTable)
	err = i.db.CreateCollection(ctx, ProductCategoryTable)
	err = i.db.CreateCollection(ctx, CartTable)
	err = i.db.CreateCollection(ctx, OrderTable)
	err = i.db.CreateCollection(ctx, AddressTable)
	if err != nil {
		// Ignore error
		return nil
	}
	return nil
}

func NewNThreads(db mongodb.DBRepo) NThreads {
	return &Impl{db: db}
}
