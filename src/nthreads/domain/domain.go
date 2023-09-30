package nthreads

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	mongodb "github.com/nxtcoder19/nthreads-backend/package/mongo-db"
	"time"
)

type Impl struct {
	db mongodb.DBInterface
}

const (
	ProductTable string = "products"
)

func (i *Impl) InsertProduct(ctx context.Context, title string, price int) (*Product, error) {
	id := uuid.New()
	record := &Product{
		Id:    Id(id.String()),
		Title: title,
		Price: price,
		Time:  time.Now(),
	}
	err := i.db.InsertRecord(ctx, ProductTable, record)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (i *Impl) GetAllProducts(ctx context.Context) ([]*Product, error) {
	products := make([]*Product, 0)
	cursor, err := i.db.GetAllRecords(ctx, ProductTable)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cer := cursor.Close(ctx); cer != nil {
			fmt.Println(cer)
		}
	}()
	for cursor.Next(ctx) {
		var product Product
		if err := cursor.Decode(&product); err != nil {
			return nil, err
		}
		products = append(products, &product)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

func (i *Impl) Init(ctx context.Context) error {
	err := i.db.CreateCollection(ctx, "products")
	if err != nil {
		// Ignore error
		return nil
	}
	return nil
}

func NewNThreads(db mongodb.DBInterface) NThreads {
	return &Impl{db: db}
}
