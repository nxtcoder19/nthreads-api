package entities

import mongodb "github.com/nxtcoder19/nthreads-backend/package/mongo-db"

type ProductCategory struct {
	Id          mongodb.ID `json:"id"`
	Name        string     `json:"name"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	ImageUrl    string     `json:"imageUrl"`
}
