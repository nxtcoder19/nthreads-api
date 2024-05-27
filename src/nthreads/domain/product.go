package domain

import (
	"context"
	"fmt"
	mongodb "github.com/nxtcoder19/nthreads-backend/package/mongo-db"
	"github.com/nxtcoder19/nthreads-backend/src/nthreads/entities"
)

func (i *Impl) CreateProduct(ctx context.Context, productIn *entities.Product) (*entities.Product, error) {
	id := i.db.NewId()
	nProduct := &entities.Product{
		Id:                  id,
		ProductCategoryName: productIn.ProductCategoryName,
		Name:                productIn.Name,
		Price:               productIn.Price,
		ImageUrl:            productIn.ImageUrl,
		Date:                productIn.Date,
		Description:         productIn.Description,
		Warranty:            productIn.Warranty,
		Place:               productIn.Place,
		AvailableColors:     productIn.AvailableColors,
		AvailableSizes:      productIn.AvailableSizes,
		Color:               productIn.Color,
		Size:                productIn.Size,
		ExtraImages:         productIn.ExtraImages,
		AvailableOffers:     productIn.AvailableOffers,
		QuestionsAnswers:    productIn.QuestionsAnswers,
		ReviewData:          productIn.ReviewData,
		Tags:                productIn.Tags,
	}
	_, err := i.db.InsertRecord(ctx, ProductTable, nProduct)
	if err != nil {
		return nil, err
	}
	return nProduct, nil
}

func (i *Impl) UpdateProduct(ctx context.Context, id string, productIn *entities.Product) (*entities.Product, error) {
	err := i.db.UpdateMany(
		ctx,
		ProductTable,
		mongodb.Filter{"id": id},
		mongodb.Filter{
			"productCategoryName": productIn.ProductCategoryName,
			"name":                productIn.Name,
			"price":               productIn.Price,
			"imageUrl":            productIn.ImageUrl,
			"date":                productIn.Date,
			"description":         productIn.Description,
			"warranty":            productIn.Warranty,
			"place":               productIn.Place,
			//"availableColors":     productIn.AvailableColors,
			//"availableSizes":      productIn.AvailableSizes,
			"color": productIn.Color,
			//"size":                productIn.Size,
			"extraImages": productIn.ExtraImages,
			//"availableOffers":     productIn.AvailableOffers,
			"questionsAnswers": productIn.QuestionsAnswers,
			//"reviewData":          productIn.ReviewData,
			//"tags":                productIn.Tags,
		},
	)

	var product entities.Product
	err = i.db.FindOne(ctx, ProductTable, &product, mongodb.Filter{"id": id})
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (i *Impl) GetProduct(ctx context.Context, id string) (*entities.Product, error) {
	var product entities.Product
	err := i.db.FindOne(ctx, ProductTable, &product, mongodb.Filter{"id": id})
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (i *Impl) GetProducts(ctx context.Context) ([]*entities.Product, error) {
	products := make([]*entities.Product, 0)
	cursor, err := i.db.Find(ctx, ProductTable, mongodb.Filter{})
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
	var product entities.Product
	err := i.db.FindOne(ctx, ProductTable, &product, mongodb.Filter{"id": id})
	if err != nil {
		return false, err
	}

	err = i.db.DeleteRecord(ctx, ProductTable, mongodb.Filter{"id": id})
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
