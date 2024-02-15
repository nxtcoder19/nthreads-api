package nthreads

import (
	"context"
	"errors"
	"github.com/google/uuid"
	mongodb "github.com/nxtcoder19/nthreads-backend/package/mongo-db"
	"github.com/nxtcoder19/nthreads-backend/src/entities"
	"time"
)

type Impl struct {
	db mongodb.DBRepo[*entities.Product]
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
	product, err := i.db.InsertRecord(ctx, ProductTable, record)
	if err != nil {
		return nil, err
	}
	//return &product, nil

	insertedProduct, ok := product.(*Product)
	if !ok {
		return nil, errors.New("inserted record is not of type *Product")
	}

	return insertedProduct, nil
}

//func (i *Impl) GetAllProducts(ctx context.Context) ([]*Product, error) {
//	products := make([]*Product, 0)
//	cursor, err := i.db.Find(ctx, ProductTable)
//	if err != nil {
//		return nil, err
//	}
//	defer func() {
//		if cer := cursor.Close(ctx); cer != nil {
//			fmt.Println(cer)
//		}
//	}()
//	for cursor.Next(ctx) {
//		var product Product
//		if err := cursor.Decode(&product); err != nil {
//			return nil, err
//		}
//		products = append(products, &product)
//	}
//	if err := cursor.Err(); err != nil {
//		return nil, err
//	}
//	return products, nil
//}

func (i *Impl) Init(ctx context.Context) error {
	err := i.db.CreateCollection(ctx, "products")
	if err != nil {
		// Ignore error
		return nil
	}
	return nil
}

func NewNThreads(db mongodb.DBRepo[*entities.Product]) NThreads {
	return &Impl{db: db}
}
