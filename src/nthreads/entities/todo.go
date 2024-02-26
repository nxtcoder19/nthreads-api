package entities

import mongo_db "github.com/nxtcoder19/nthreads-backend/package/mongo-db"

type Todo struct {
	Id          mongo_db.ID `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
}
