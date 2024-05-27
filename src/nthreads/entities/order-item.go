package entities

import mongodb "github.com/nxtcoder19/nthreads-backend/package/mongo-db"

type OrderItem struct {
	Id      mongodb.ID `json:"id"`
	Email   string     `json:"email"`
	Product Product    `json:"product"`
}
