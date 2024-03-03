package domain

import (
	"context"
	"fmt"
	mongo_db "github.com/nxtcoder19/nthreads-backend/package/mongo-db"
	"github.com/nxtcoder19/nthreads-backend/src/nthreads/entities"
)

func (i *Impl) CreateProduct(ctx context.Context, name string, price string, imageUrl string, date string, description string, warranty string, place string, extraImages []string) (*entities.Product, error) {
	id := i.db.NewId()
	product := entities.Product{
		Id:          id,
		Name:        name,
		Price:       price,
		ImageUrl:    imageUrl,
		Date:        date,
		Description: description,
		Warranty:    warranty,
		Place:       place,
		ExtraImages: extraImages,
	}
	_, err := i.db.InsertRecord(ctx, ProductTable, product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (i *Impl) UpdateProduct(ctx context.Context, id string, price string, imageUrl string, date string, warranty string, place string) (*entities.Product, error) {
	err := i.db.UpdateMany(
		ctx,
		ProductTable,
		mongo_db.Filter{"id": id},
		mongo_db.Filter{
			"price":     price,
			"image_url": imageUrl,
			"date":      date,
			"warranty":  warranty,
			"place":     place,
		},
	)

	var product entities.Product
	err = i.db.FindOne(ctx, ProductTable, &product, mongo_db.Filter{"id": id})
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (i *Impl) GetProduct(ctx context.Context, id string) (*entities.Product, error) {
	var product entities.Product
	err := i.db.FindOne(ctx, ProductTable, &product, mongo_db.Filter{"id": id})
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (i *Impl) GetProducts(ctx context.Context) ([]*entities.Product, error) {
	products := make([]*entities.Product, 0)
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
		var product entities.Product
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

func (i *Impl) DeleteProduct(ctx context.Context, id string) (bool, error) {
	var todo entities.Product
	err := i.db.FindOne(ctx, ProductTable, &todo, mongo_db.Filter{"id": id})
	if err != nil {
		return false, err
	}

	err = i.db.DeleteRecord(ctx, ProductTable, mongo_db.Filter{"id": id})
	if err != nil {
		return false, err
	}
	return true, nil
}

//func (i *Impl) InsertProduct(ctx context.Context, title string, price int) (*Product, error) {
//	id := uuid.New()
//	record := &Product{
//		Id:    Id(id.String()),
//		Title: title,
//		Price: price,
//		Time:  time.Now(),
//	}
//	_, err := i.db.InsertRecord(ctx, ProductTable, record)
//	if err != nil {
//		return nil, err
//	}
//	return record, nil
//
//	//insertedProduct, ok := product.(*Product)
//	//if !ok {
//	//	return nil, errors.New("inserted record is not of type *Product")
//	//}
//	//
//	//return insertedProduct, nil
//}
//
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
