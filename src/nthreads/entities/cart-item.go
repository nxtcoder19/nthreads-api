package entities

import mongodb "github.com/nxtcoder19/nthreads-backend/package/mongo-db"

type CartItem struct {
	Id       mongodb.ID `json:"id"`
	Email    string     `json:"email"`
	Product  Product    `json:"product"`
	Quantity int        `json:"quantity"`
}
