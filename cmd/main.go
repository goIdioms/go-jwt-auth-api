package main

import (
	"context"
	"fmt"
	"log"
	"test/pck/auth/controllers"
	"test/pck/auth/repository"
	"test/pck/auth/services"
	"test/pck/cache"
	"test/pck/database"
	"test/pck/router"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	_ "github.com/lib/pq"
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
	authController := controllers.NewAuthController(authService)

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

	micro.Get("/healthchecker", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "JSON Web Token Authentication and Authorization in Golang",
		})
	})

	app.All("*", func(c *fiber.Ctx) error {
		path := c.Path()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": fmt.Sprintf("Path: %v does not exists on this server", path),
		})
	})

	log.Fatal(app.Listen(":" + config.Port))
}
