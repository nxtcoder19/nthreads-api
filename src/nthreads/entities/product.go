package entities

import mongo_db "github.com/nxtcoder19/nthreads-backend/package/mongo-db"

type Product struct {
	Id          mongo_db.ID `json:"id"`
	Name        string      `json:"name"`
	Price       string      `json:"price"`
	ImageUrl    string      `json:"image_url"`
	Date        string      `json:"date"`
	Description string      `json:"description"`
	Warranty    string      `json:"warranty"`
	Place       string      `json:"place"`
}
