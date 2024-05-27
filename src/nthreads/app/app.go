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

func (s *ServerImpl) getSessionDataFromSessionId(ctx context.Context, sessionId string) (*entities.SessionData, error) {
	if sessionId == "" {
		return nil, errors.New("no session found")
	}
	sessionRawData, err := s.sessionCache.Get(ctx, sessionId)
	if err != nil {
		return nil, err
	}
	sessionData := &entities.SessionData{}
	err = json.Unmarshal(sessionRawData, sessionData)
	if err != nil {
		return nil, err
	}
	return sessionData, nil
}

func (s *ServerImpl) getContextWithSessionData(ctx *fiber.Ctx) (context.Context, error) {
	getSessionDataFromFiberContext := func(ctx *fiber.Ctx) (*entities.SessionData, error) {
		sessionId := ctx.Cookies("session_id", "")
		if sessionId == "" {
			sessionId = ctx.GetReqHeaders()["App-Sessionid"]
		}
		if sessionId == "" {
			return nil, errors.New("no session found")
		}
		return s.getSessionDataFromSessionId(ctx.Context(), sessionId)
	}
	sessionData, err := getSessionDataFromFiberContext(ctx)
	if err != nil {
		return nil, err
	}
	return context.WithValue(ctx.Context(), "user-session", sessionData), nil
}

func (s *ServerImpl) CreateSession(ctx context.Context, email string, name string, id string) (string, error) {
	sessionId, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	sessionData := entities.SessionData{
		UserEmail: email,
		UserName:  name,
		UserId:    id,
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

	//app.Post("/api/product/create", func(ctx *fiber.Ctx) error {
	//	data := struct {
	//		Title string `json:"title"`
	//		Price int    `json:"price"`
	//	}{}
	//	body := ctx.Body()
	//	err := json.Unmarshal(body, &data)
	//	if err != nil {
	//		return err
	//	}
	//	product, err := s.threads.InsertProduct(context.TODO(), data.Title, data.Price)
	//	if err != nil {
	//		return err
	//	}
	//	return ctx.JSON(product)
	//})

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

		user, err := s.threads.GetUser(ctx.Context(), ctx.Params("email"))
		if err != nil {
			return err
		}

		if !loginMessage {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errorMessage": "Login failed",
			})
		}

		sessionId, err := s.CreateSession(ctx.Context(), ctx.Params("email"), user.FirstName, string(user.Id))
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
			Message: "Login Successfull",
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

	/// Products
	app.Post("/api/product/create", func(ctx *fiber.Ctx) error {
		var data *entities.Product
		if err := ctx.BodyParser(&data); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
		}
		product, err := s.threads.CreateProduct(context.TODO(), data)
		if err != nil {
			return err
		}
		return ctx.JSON(product)
	})

	app.Put("/api/product/update/:id", func(ctx *fiber.Ctx) error {
		//data := struct {
		//	Price    string `json:"price"`
		//	ImageUrl string `json:"image_url"`
		//	Date     string `json:"date"`
		//	Warranty string `json:"warranty"`
		//	Place    string `json:"place"`
		//}{}
		//body := ctx.Body()
		//err := json.Unmarshal(body, &data)
		//if err != nil {
		//	return err
		//}
		var data *entities.Product
		if err := ctx.BodyParser(&data); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
		}
		UpdatedProduct, err := s.threads.UpdateProduct(context.TODO(), ctx.Params("id"), data)
		if err != nil {
			return err
		}
		return ctx.JSON(UpdatedProduct)
	})

	app.Get("/api/product/get/:id", func(ctx *fiber.Ctx) error {
		product, err := s.threads.GetProduct(context.TODO(), ctx.Params("id"))
		if err != nil {
			return err
		}
		return ctx.JSON(product)
	})

	app.Get("/api/product/get", func(ctx *fiber.Ctx) error {
		products, err := s.threads.GetProducts(context.TODO())
		if err != nil {
			return err
		}
		return ctx.JSON(products)
	})

	app.Delete("/api/product/delete/:id", func(ctx *fiber.Ctx) error {
		productStatus, err := s.threads.DeleteProduct(context.TODO(), ctx.Params("id"))
		if err != nil {
			return err
		}

		if productStatus == false {
			return ctx.SendStatus(400)
		}

		return ctx.JSON(struct {
			Status bool `json:"status"`
		}{
			Status: productStatus,
		})
	})

	/// Products category
	app.Post("/api/product-category/create", func(ctx *fiber.Ctx) error {
		data := struct {
			Name        string `json:"name"`
			Title       string `json:"title"`
			Description string `json:"description"`
			ImageUrl    string `json:"imageUrl"`
		}{}
		body := ctx.Body()
		err := json.Unmarshal(body, &data)
		if err != nil {
			return err
		}
		//fmt.Println("data", data)
		productCategory, err := s.threads.CreateProductCategory(context.TODO(), data.Name, data.Title, data.Description, data.ImageUrl)
		if err != nil {
			return err
		}
		return ctx.JSON(productCategory)
	})

	app.Put("/api/product-category/update/:id", func(ctx *fiber.Ctx) error {
		data := struct {
			Name        string `json:"name"`
			Title       string `json:"title"`
			Description string `json:"description"`
			ImageUrl    string `json:"imageUrl"`
		}{}
		body := ctx.Body()
		err := json.Unmarshal(body, &data)
		if err != nil {
			return err
		}
		UpdatedProductCategory, err := s.threads.UpdateProductCategory(context.TODO(), ctx.Params("id"), data.Name, data.Title, data.Description, data.ImageUrl)
		if err != nil {
			return err
		}
		return ctx.JSON(UpdatedProductCategory)
	})

	app.Get("/api/product-category/get/:id", func(ctx *fiber.Ctx) error {
		product, err := s.threads.GetProductCategory(context.TODO(), ctx.Params("id"))
		if err != nil {
			return err
		}
		return ctx.JSON(product)
	})

	app.Get("/api/product-category/get", func(ctx *fiber.Ctx) error {
		productCategories, err := s.threads.GetProductCategories(context.TODO())
		if err != nil {
			return err
		}
		return ctx.JSON(productCategories)
	})

	app.Delete("/api/product-category/delete/:id", func(ctx *fiber.Ctx) error {
		productCategoryStatus, err := s.threads.DeleteProductCategory(context.TODO(), ctx.Params("id"))
		if err != nil {
			return err
		}

		if productCategoryStatus == false {
			return ctx.SendStatus(400)
		}

		return ctx.JSON(struct {
			Status bool `json:"status"`
		}{
			Status: productCategoryStatus,
		})
	})

	app.Use("/api", func(ctx *fiber.Ctx) error {
		sessionCtx, err := s.getContextWithSessionData(ctx)
		if err != nil {
			return err
		}
		ctx.SetUserContext(sessionCtx)
		return ctx.Next()
	})

	/// Auth
	app.Get("/api/.secret/session-data", func(ctx *fiber.Ctx) error {
		sessionData := ctx.UserContext().Value("user-session").(*entities.SessionData)
		return ctx.JSON(sessionData)
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

	/// Cart items
	app.Post("/api/cart/create", func(ctx *fiber.Ctx) error {
		var data *entities.CartItem
		if err := ctx.BodyParser(&data); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
		}
		cartItem, err := s.threads.AddItemToCart(ctx.UserContext(), &data.Product, data.Quantity)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add item to cart"})
		}
		return ctx.JSON(cartItem)
	})

	app.Get("/api/cart/get", func(ctx *fiber.Ctx) error {
		cartItems, err := s.threads.GetCartItems(ctx.UserContext())
		if err != nil {
			return err
		}
		return ctx.JSON(cartItems)
	})

	app.Delete("/api/cart/delete/:id", func(ctx *fiber.Ctx) error {
		cartItemStatus, err := s.threads.RemoveCartItem(ctx.UserContext(), ctx.Params("id"))
		if err != nil {
			return err
		}

		if cartItemStatus == false {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to remove item from cart"})
		}

		return ctx.JSON(struct {
			Status bool `json:"status"`
		}{
			Status: cartItemStatus,
		})
	})

	/// Order items
	app.Post("/api/order/create", func(ctx *fiber.Ctx) error {
		var data *entities.OrderItem
		if err := ctx.BodyParser(&data); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
		}
		orderItem, err := s.threads.AddItemToOrder(ctx.UserContext(), &data.Product)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add item to order"})
		}
		return ctx.JSON(orderItem)
	})

	app.Get("/api/order/get/:id", func(ctx *fiber.Ctx) error {
		orderItem, err := s.threads.GetOrderItem(ctx.UserContext(), ctx.Params("id"))
		if err != nil {
			return err
		}
		return ctx.JSON(orderItem)
	})

	app.Get("/api/order/get", func(ctx *fiber.Ctx) error {
		orderItems, err := s.threads.GetOrderItems(ctx.UserContext())
		if err != nil {
			return err
		}
		return ctx.JSON(orderItems)
	})

	/// Address
	app.Post("/api/address/create", func(ctx *fiber.Ctx) error {
		var data *entities.Address
		if err := ctx.BodyParser(&data); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
		}
		orderItem, err := s.threads.CreateAddress(ctx.UserContext(), data)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add address"})
		}
		return ctx.JSON(orderItem)
	})

	app.Put("/api/address/update/:id", func(ctx *fiber.Ctx) error {
		var data *entities.Address
		if err := ctx.BodyParser(&data); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
		}
		UpdatedProduct, err := s.threads.UpdateAddress(ctx.UserContext(), ctx.Params("id"), data)
		if err != nil {
			return err
		}
		return ctx.JSON(UpdatedProduct)
	})

	app.Get("/api/address/get/:id", func(ctx *fiber.Ctx) error {
		product, err := s.threads.GetAddress(ctx.UserContext(), ctx.Params("id"))
		if err != nil {
			return err
		}
		return ctx.JSON(product)
	})

	app.Get("/api/address/get", func(ctx *fiber.Ctx) error {
		products, err := s.threads.GetAddresses(ctx.UserContext())
		if err != nil {
			return err
		}
		return ctx.JSON(products)
	})

	app.Delete("/api/address/delete/:id", func(ctx *fiber.Ctx) error {
		productStatus, err := s.threads.DeleteAddress(ctx.UserContext(), ctx.Params("id"))
		if err != nil {
			return err
		}

		if productStatus == false {
			return ctx.SendStatus(400)
		}

		return ctx.JSON(struct {
			Status bool `json:"status"`
		}{
			Status: productStatus,
		})
	})

	/// Todo
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
