package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/google/uuid"
	"github.com/nxtcoder19/nthreads-backend/package/errors"
	"github.com/nxtcoder19/nthreads-backend/package/redis"
	"github.com/nxtcoder19/nthreads-backend/src/nthreads/domain"
	"github.com/nxtcoder19/nthreads-backend/src/nthreads/entities"
)

type Server interface {
	Init()
	Start(port string) error
}

type ServerImpl struct {
	threads      domain.NThreads
	app          *fiber.App
	sessionCache redis.Cache
}

func (s *ServerImpl) CreateSession(ctx context.Context, email string) (string, error) {
	sessionId, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	sessionData := entities.SessionData{
		Email: email,
	}
	sessionDataBytes, err := json.Marshal(sessionData)
	if err != nil {
		return "", err
	}

	err = s.sessionCache.Set(ctx, sessionId.String(), sessionDataBytes, nil)
	if err != nil {
		return "", err
	}
	return sessionId.String(), nil
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

	app.Post("/api/auth/signup", func(ctx *fiber.Ctx) error {
		data := struct {
			FirstName string `json:"first_name"`
			Lastname  string `json:"last_name"`
			Email     string `json:"email"`
			Password  string `json:"password"`
		}{}
		body := ctx.Body()
		err := json.Unmarshal(body, &data)
		if err != nil {
			return err
		}
		fmt.Println("data", data)
		user, err := s.threads.SignUp(context.TODO(), data.FirstName, data.Lastname, data.Email, data.Password)
		if err != nil {
			return err
		}
		return ctx.JSON(user)
	})

	app.Get("/api/auth/login/:email/:password", func(ctx *fiber.Ctx) error {
		loginMessage, err := s.threads.Login(context.TODO(), ctx.Params("email"), ctx.Params("password"))
		if err != nil {
			return err
		}

		sessionId, err := s.CreateSession(ctx.Context(), ctx.Params("email"))
		if err != nil {
			return err
		}
		ctx.Set("App-Sessionid", sessionId)
		ctx.Cookie(&fiber.Cookie{
			Name:  "session_id",
			Value: sessionId,
		})

		return ctx.JSON(struct {
			Message string `json:"message"`
		}{
			Message: loginMessage,
		})
	})

	app.Post("/api/auth/updateUser", func(ctx *fiber.Ctx) error {
		data := struct {
			Email     string `json:"email"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
		}{}
		body := ctx.Body()
		err := json.Unmarshal(body, &data)
		if err != nil {
			return err
		}
		updatedUser, err := s.threads.UpdateUser(context.TODO(), data.Email, data.FirstName, data.LastName)
		if err != nil {
			return err
		}
		return ctx.JSON(updatedUser)
	})

	app.Get("/api/auth/deleteUser/:email", func(ctx *fiber.Ctx) error {
		deletedMessage, err := s.threads.DeleteUser(context.TODO(), ctx.Params("email"))
		if err != nil {
			return err
		}
		return ctx.JSON(struct {
			Message string `json:"message"`
		}{
			Message: deletedMessage,
		})
	})

	app.Delete("/api/auth/logout", func(ctx *fiber.Ctx) error {
		sessionId := ctx.Cookies("session_id", "")
		fmt.Println(sessionId)
		if sessionId == "" {
			return errors.New("no session found")
		}
		err := s.sessionCache.Del(ctx.Context(), sessionId)
		if err != nil {
			return err
		}
		ctx.ClearCookie("session_id")
		return ctx.JSON(struct {
			Message string `json:"message"`
		}{
			Message: "user logout",
		})
	})

	app.Post("/api/todo/create", func(ctx *fiber.Ctx) error {
		data := struct {
			Title       string `json:"title"`
			Description string `json:"description"`
		}{}
		body := ctx.Body()
		err := json.Unmarshal(body, &data)
		if err != nil {
			return err
		}
		fmt.Println("data", data)
		user, err := s.threads.CreateTodo(context.TODO(), data.Title, data.Description)
		if err != nil {
			return err
		}
		return ctx.JSON(user)
	})

	app.Post("/api/todo/update/:id", func(ctx *fiber.Ctx) error {
		data := struct {
			Title       string `json:"title"`
			Description string `json:"description"`
		}{}
		body := ctx.Body()
		err := json.Unmarshal(body, &data)
		if err != nil {
			return err
		}
		updatedUser, err := s.threads.UpdateTodo(context.TODO(), ctx.Params("id"), data.Title, data.Description)
		if err != nil {
			return err
		}
		return ctx.JSON(updatedUser)
	})

	app.Get("/api/todo/get/:id", func(ctx *fiber.Ctx) error {
		todo, err := s.threads.GetTodo(context.TODO(), ctx.Params("id"))
		if err != nil {
			return err
		}
		return ctx.JSON(todo)
	})

	app.Get("/api/todo/get", func(ctx *fiber.Ctx) error {
		todos, err := s.threads.GetTodos(context.TODO())
		if err != nil {
			return err
		}
		return ctx.JSON(todos)
	})

	app.Delete("/api/todo/delete/:id", func(ctx *fiber.Ctx) error {
		deletedMessage, err := s.threads.DeleteTodo(context.TODO(), ctx.Params("id"))
		if err != nil {
			return err
		}
		return ctx.JSON(struct {
			Message string `json:"message"`
		}{
			Message: deletedMessage,
		})
	})

}

func (s *ServerImpl) Start(port string) error {
	return s.app.Listen(port)
}

func NewServer(threads domain.NThreads, sessionCache redis.Cache) Server {
	return &ServerImpl{
		threads:      threads,
		sessionCache: sessionCache,
	}
}
