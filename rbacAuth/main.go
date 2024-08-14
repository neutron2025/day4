package main

import (
	"blog-auth-server/controllers"
	"blog-auth-server/middleware"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/net/context"
)

var ctx context.Context
var err error
var client *mongo.Client
var MongoUri string = "mongodb://neutronroot:pass123@89.47.166.251:27017/dbname?authSource=admin"
var userController *controllers.UserController
var middleware1 *middleware.Middleware

func init() {
	log.Println()
	ctx = context.Background()
	client, err = mongo.Connect(ctx,
		options.Client().ApplyURI(MongoUri))
	if err = client.Ping(context.TODO(),
		readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	collection := client.Database("auth-server").Collection("users")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "89.47.166.251:6379",
		Password: "redis123",
		DB:       0})
	status := redisClient.Ping(ctx)
	fmt.Println(status)
	userController = controllers.NewUserController(collection,
		ctx, redisClient)
	middleware1 = middleware.NewMiddleware(ctx, redisClient)
}

func main() {
	app := fiber.New()
	app.Use(logger.New())
	app.Post("/signup", userController.CreateUser)
	app.Post("/login", userController.Login)
	app.Post("/addPermission", userController.AddPermission)
	app.Post("/adminTestRoute", middleware1.AdminMiddlewareHandler, userController.TestRoute)
	err := app.Listen(":3000")
	if err != nil {
		log.Fatal("Error in running the server")
		return
	}
	log.Println("Server is running")
}

// /addPermission
// {
//   "username": "root",
//   "permission": {
//     "entry": 5,
//     "add_flag": true,
//     "admin_flag": true
//   }
// }

// /adminTestRoute
// Authorization: Bearer <your.jwt.token.here>
// Entry: <entry-integer-value>

// entry 级别
// add_flag 添加权限
// admin_flag  管理员权限
