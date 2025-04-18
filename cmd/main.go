package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/goIdioms/go-jwt-auth-api/pck/auth/controllers"
	"github.com/goIdioms/go-jwt-auth-api/pck/auth/repository"
	"github.com/goIdioms/go-jwt-auth-api/pck/auth/services"
	"github.com/goIdioms/go-jwt-auth-api/pck/cache"
	"github.com/goIdioms/go-jwt-auth-api/pck/database"
	"github.com/goIdioms/go-jwt-auth-api/pck/router"

	_ "github.com/goIdioms/go-jwt-auth-api/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	_ "github.com/lib/pq"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	app         *fiber.App
	ctx         context.Context
	mongoclient *mongo.Client
)

func init() {
	config, err := database.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load config: ", err)
	}

	ctx = context.Background()
	mongoconn := options.Client().ApplyURI(config.DBUri)
	mongoclient, err = mongo.Connect(ctx, mongoconn)
	if err != nil {
		panic(err)
	}

	if err := mongoclient.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("Connected to MongoDB successfully")

	app = fiber.New(fiber.Config{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	})

	app.Use(logger.New())
	app.Use(recover.New())
}

// @title JWT Authentication API
// @version 1.0
// @description API для аутентификации и авторизации с использованием JWT токенов
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /api
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Используйте формат: "Bearer {token}"
func main() {
	config, err := database.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load config: ", err)
	}
	defer mongoclient.Disconnect(ctx)

	mongoDB := mongoclient.Database("golang_mongodb")
	database.UserCollection = mongoDB.Collection("users")

	redisCache := cache.NewRedisCache(config.RedisUri)
	authRepo := repository.NewAuthRepository(ctx)
	authService := services.NewAuthService(ctx, authRepo, redisCache)
	authController := controllers.NewAuthController(authService, redisCache)

	app := fiber.New()
	micro := fiber.New()
	router.SetupRoutes(micro, authController)

	app.Mount("/api", micro)
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowMethods:     "GET, POST",
		AllowCredentials: true,
	}))

	// HealthCheck godoc
	// @Summary Проверка работоспособности API
	// @Description Проверяет, что API работает корректно
	// @Tags health
	// @Accept json
	// @Produce json
	// @Success 200 {object} map[string]interface{}
	// @Router /healthchecker [get]
	micro.Get("/healthchecker", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "JSON Web Token Authentication and Authorization in Golang",
		})
	})

	app.Get("/swagger/*", fiberSwagger.FiberWrapHandler())

	app.All("*", func(c *fiber.Ctx) error {
		path := c.Path()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": fmt.Sprintf("Path: %v does not exists on this server", path),
		})
	})

	log.Fatal(app.Listen(":" + config.Port))
}
