package main

import (
	"context"
	mongodb "github.com/nxtcoder19/nthreads-backend/package/mongo-db"
	"github.com/nxtcoder19/nthreads-backend/package/redis"
	"github.com/nxtcoder19/nthreads-backend/src/nthreads/app"
	nthreads "github.com/nxtcoder19/nthreads-backend/src/nthreads/domain"
	"os"
)

func main() {
	mongoUrl := os.Getenv("MONGO_URI")
	db := mongodb.NewDB("n-threads-db", mongoUrl)
	err := db.ConnectDB(context.TODO())
	if err != nil {
		panic(err)
	}

	threads := nthreads.NewNThreads(db)
	err = threads.Init(context.TODO())
	if err != nil {
		panic(err)
	}

	sessionCache := redis.NewRedisCache("localhost:6379")

	server := app.NewServer(threads, sessionCache)
	server.Init()
	err = server.Start(":3002")
	if err != nil {
		panic(err)
	}
}
