package app

import (
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/nxtcoder19/nthreads-backend/src/nthreads/domain"
)

type Server interface {
	Init()
	Start(port string) error
}

type ServerImpl struct {
	threads nthreads.NThreads
	app     *fiber.App
}

func (s *ServerImpl) Init() {
	app := fiber.New()
	s.app = app

	// allow cors
	app.Use(cors.New(cors.Config{
		AllowHeaders:     "*",
		AllowOrigins:     "*",
		AllowCredentials: true,
		ExposeHeaders:    "*",
		AllowOriginsFunc: func(origin string) bool {
			return true
		},
		AllowMethods: "*",
	}))

	app.Post("/api/product/create", func(ctx *fiber.Ctx) error {
		data := struct {
			Title string `json:"title"`
			Price int    `json:"price"`
		}{}
		body := ctx.Body()
		err := json.Unmarshal(body, &data)
		if err != nil {
			return err
		}
		product, err := s.threads.InsertProduct(context.TODO(), data.Title, data.Price)
		if err != nil {
			return err
		}
		return ctx.JSON(product)
	})

	//app.Get("/api/product/get", func(ctx *fiber.Ctx) error {
	//	products, err := s.threads.GetAllProducts(context.TODO())
	//	if err != nil {
	//		return err
	//	}
	//	return ctx.JSON(products)
	//})

}

func (s *ServerImpl) Start(port string) error {
	return s.app.Listen(port)
}

func NewServer(threads nthreads.NThreads) Server {
	return &ServerImpl{
		threads: threads,
	}
}
