package domain

import (
	"context"
	"fmt"
	mongodb "github.com/nxtcoder19/nthreads-backend/package/mongo-db"
	"github.com/nxtcoder19/nthreads-backend/src/nthreads/entities"
)

func (i *Impl) CreateProductCategory(ctx context.Context, name string, title string, description string, imageUrl string) (*entities.ProductCategory, error) {
	id := i.db.NewId()
	productCategory := entities.ProductCategory{
		Id:          id,
		Name:        name,
		Title:       title,
		Description: description,
		ImageUrl:    imageUrl,
	}
	_, err := i.db.InsertRecord(ctx, ProductCategoryTable, productCategory)
	if err != nil {
		return nil, err
	}
	return &productCategory, nil
}

func (i *Impl) UpdateProductCategory(ctx context.Context, id string, name string, title string, description string, imageUrl string) (*entities.ProductCategory, error) {
	err := i.db.UpdateMany(
		ctx,
		ProductCategoryTable,
		mongodb.Filter{"id": id},
		mongodb.Filter{
			"name":        name,
			"title":       title,
			"description": description,
			"image":       imageUrl,
		},
	)

	var productCategory entities.ProductCategory
	err = i.db.FindOne(ctx, ProductCategoryTable, &productCategory, mongodb.Filter{"id": id})
	if err != nil {
		return nil, err
	}

	return &productCategory, nil
}

func (i *Impl) GetProductCategory(ctx context.Context, id string) (*entities.ProductCategory, error) {
	var productCategory entities.ProductCategory
	err := i.db.FindOne(ctx, ProductCategoryTable, &productCategory, mongodb.Filter{"id": id})
	if err != nil {
		return nil, err
	}

	return &productCategory, nil
}

func (i *Impl) GetProductCategories(ctx context.Context) ([]*entities.ProductCategory, error) {
	productCategories := make([]*entities.ProductCategory, 0)
	cursor, err := i.db.Find(ctx, ProductCategoryTable, mongodb.Filter{})
	if err != nil {
		return nil, err
	}
	defer func() {
		if cer := cursor.Close(ctx); cer != nil {
			fmt.Println(cer)
		}
	}()
	for cursor.Next(ctx) {
		var productCategory entities.ProductCategory
		if err := cursor.Decode(&productCategory); err != nil {
			return nil, err
		}
		productCategories = append(productCategories, &productCategory)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return productCategories, nil
}

func (i *Impl) DeleteProductCategory(ctx context.Context, id string) (bool, error) {
	var productCategory entities.ProductCategory
	err := i.db.FindOne(ctx, ProductCategoryTable, &productCategory, mongodb.Filter{"id": id})
	if err != nil {
		return false, err
	}

	err = i.db.DeleteRecord(ctx, ProductCategoryTable, mongodb.Filter{"id": id})
	if err != nil {
		return false, err
	}
	return true, nil
}
