package domain

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
)

func (i *Impl) InsertProduct(ctx context.Context, title string, price int) (*Product, error) {
	id := uuid.New()
	record := &Product{
		Id:    Id(id.String()),
		Title: title,
		Price: price,
		Time:  time.Now(),
	}
	_, err := i.db.InsertRecord(ctx, ProductTable, record)
	if err != nil {
		return nil, err
	}
	return record, nil

	//insertedProduct, ok := product.(*Product)
	//if !ok {
	//	return nil, errors.New("inserted record is not of type *Product")
	//}
	//
	//return insertedProduct, nil
}

func (i *Impl) GetAllProducts(ctx context.Context) ([]*Product, error) {
	products := make([]*Product, 0)
	cursor, err := i.db.Find(ctx, ProductTable)
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
