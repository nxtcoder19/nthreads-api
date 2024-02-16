package entities

import (
	mongo_db "github.com/nxtcoder19/nthreads-backend/package/mongo-db"
	"time"
)

type Product struct {
	mongo_db.BaseEntity `json:",inline" graphql:"noinput"`
	Title               string    `json:"title"`
	Price               int       `json:"price"`
	Time                time.Time `json:"time"`
}
