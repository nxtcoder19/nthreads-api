package domain

import (
	"context"
	"fmt"
	mongo_db "github.com/nxtcoder19/nthreads-backend/package/mongo-db"
	"github.com/nxtcoder19/nthreads-backend/src/nthreads/entities"
)

func (i *Impl) AddItemToOrder(ctx context.Context, product *entities.Product) (*entities.OrderItem, error) {
	id := i.db.NewId()
	session := getSessionData(ctx)
	order := entities.OrderItem{
		Id:    id,
		Email: session.UserEmail,
		Product: entities.Product{
			Id:                  product.Id,
			ProductCategoryName: product.ProductCategoryName,
			Name:                product.Name,
			Price:               product.Price,
			ImageUrl:            product.ImageUrl,
			Date:                product.Date,
			Description:         product.Description,
			Warranty:            product.Warranty,
			Place:               product.Place,
		},
	}
	_, err := i.db.InsertRecord(ctx, OrderTable, order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (i *Impl) GetOrderItem(ctx context.Context, id string) (*entities.OrderItem, error) {
	var orderItems entities.OrderItem
	err := i.db.FindOne(ctx, OrderTable, &orderItems, mongo_db.Filter{"id": id})
	if err != nil {
		return nil, err
	}

	return &orderItems, nil
}

func (i *Impl) GetOrderItems(ctx context.Context) ([]*entities.OrderItem, error) {
	orderItems := make([]*entities.OrderItem, 0)
	session := getSessionData(ctx)
	cursor, err := i.db.Find(ctx, OrderTable, mongo_db.Filter{"email": session.UserEmail})
	if err != nil {
		return nil, err
	}
	defer func() {
		if cer := cursor.Close(ctx); cer != nil {
			fmt.Println(cer)
		}
	}()
	for cursor.Next(ctx) {
		var orderItem entities.OrderItem
		if err := cursor.Decode(&orderItem); err != nil {
			return nil, err
		}
		orderItems = append(orderItems, &orderItem)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return orderItems, nil
}
