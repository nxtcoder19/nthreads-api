package domain

import (
	"context"
	mongodb "github.com/nxtcoder19/nthreads-backend/package/mongo-db"
)

type Impl struct {
	db mongodb.DBRepo
}

const (
	ProductTable string = "products"
	AuthTable    string = "auth"
	TodoTable    string = "todo"
)

func (i *Impl) Init(ctx context.Context) error {
	//err := i.db.CreateCollection(ctx, "products")
	err := i.db.CreateCollection(ctx, AuthTable)
	err = i.db.CreateCollection(ctx, TodoTable)
	if err != nil {
		// Ignore error
		return nil
	}
	return nil
}

func NewNThreads(db mongodb.DBRepo) NThreads {
	return &Impl{db: db}
}
