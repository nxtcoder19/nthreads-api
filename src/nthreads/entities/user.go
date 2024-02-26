package entities

import (
	mongo_db "github.com/nxtcoder19/nthreads-backend/package/mongo-db"
)

type User struct {
	Id             mongo_db.ID `json:"id"`
	FirstName      string      `json:"first_name"`
	LastName       string      `json:"last_name"`
	Email          string      `json:"email"`
	Password       string      `json:"password"`
	VerifyPassword string      `json:"verify_password"`
}

type SessionData struct {
	Email string `json:"email"`
}
