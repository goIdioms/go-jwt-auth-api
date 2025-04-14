package main

import (
	"context"
	"fmt"
	"log"
	"test/controllers"
	"test/database"
	"test/repository"
	"test/services"
	"time"

	"github.com/go-redis/redis/v8"
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
	redisclient *redis.Client
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

	redisclient = redis.NewClient(&redis.Options{
		Addr: config.RedisUri,
	})
	if _, err := redisclient.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	err = redisclient.Set(ctx, "test", "Redis and MongoDB", 0).Err()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to Redis successfully")

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

	// value, err := redisclient.Get(ctx, "test").Result()
	// if err == redis.Nil {
	// 	fmt.Println("key: test does not exist")
	// } else if err != nil {
	// 	panic(err)
	// }

	mongoDB := mongoclient.Database("golang_mongodb")
	database.UserCollection = mongoDB.Collection("users")

	authRepo := repository.NewAuthRepository(ctx)
	authService := services.NewAuthService(authRepo)
	authController := controllers.NewAuthController(authService)

	app := fiber.New()
	micro := fiber.New()
	SetupRoutes(micro, authController)

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
