package domain

import (
	"context"
	"fmt"
	mongo_db "github.com/nxtcoder19/nthreads-backend/package/mongo-db"
	"github.com/nxtcoder19/nthreads-backend/src/nthreads/entities"
)

func (i *Impl) AddItemToCart(ctx context.Context, product *entities.Product, quantity int) (*entities.CartItem, error) {
	id := i.db.NewId()
	session := getSessionData(ctx)
	cartItem := entities.CartItem{
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
		Quantity: quantity,
	}
	_, err := i.db.InsertRecord(ctx, CartTable, cartItem)
	if err != nil {
		return nil, err
	}
	return &cartItem, nil
}

func (i *Impl) UpdateCartItem(ctx context.Context, id string, quantity int) (*entities.CartItem, error) {
	//TODO implement me
	panic("implement me")
}

func (i *Impl) RemoveCartItem(ctx context.Context, id string) (bool, error) {
	var cartItem entities.CartItem
	session := getSessionData(ctx)
	err := i.db.FindOne(ctx, CartTable, &cartItem, mongo_db.Filter{"id": id, "email": session.UserEmail})
	if err != nil {
		return false, err
	}

	err = i.db.DeleteRecord(ctx, CartTable, mongo_db.Filter{"id": id})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (i *Impl) GetCartItems(ctx context.Context) ([]*entities.CartItem, error) {
	cartItems := make([]*entities.CartItem, 0)
	session := getSessionData(ctx)
	cursor, err := i.db.Find(ctx, CartTable, mongo_db.Filter{"email": session.UserEmail})
	if err != nil {
		return nil, err
	}
	defer func() {
		if cer := cursor.Close(ctx); cer != nil {
			fmt.Println(cer)
		}
	}()
	for cursor.Next(ctx) {
		var cartItem entities.CartItem
		if err := cursor.Decode(&cartItem); err != nil {
			return nil, err
		}
		cartItems = append(cartItems, &cartItem)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return cartItems, nil
}
