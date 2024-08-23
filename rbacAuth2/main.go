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
var MongoUri string = "mongodb://neutronroot:pass123@89.47.166.251:27017/ecomm?authSource=admin"
var userController *controllers.UserController
var productController *controllers.ProductController
var cartController *controllers.CartController
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
	userController = controllers.NewUserController(collection, ctx, redisClient)
	productController = controllers.NewProductController(collection, ctx)
	cartController = controllers.NewCartController(collection, ctx)
	middleware1 = middleware.NewMiddleware(ctx, redisClient)
}

func main() {
	app := fiber.New()
	app.Use(logger.New())
	app.Get("/", productController.AllProduct)      //产品展示页
	app.Get("/product", productController.FetchOne) //产品信息页
	app.Post("/signup", userController.CreateUser)
	app.Post("/login", userController.Login)

	app.Get("/cart", middleware1.UserMiddlewareHandler, cartController.AllfromCart) //产品结算页 用户可以增删查 改数量,前端localStorage，登录后同步到数据库 在支付的时候需要登录session
	app.Post("/cart", middleware1.UserMiddlewareHandler, cartController.AddtoCart)  //后端接收到购物车数据后，将其与当前登录的用户账户关联起来, 关联成功后，前端可以清除localStorage
	app.Get("/pay", middleware1.UserMiddlewareHandler, userController.OrderPay)
	app.Get("/userinfo", middleware1.UserMiddlewareHandler, userController.GetOneUser)

	app.Post("/addPermission", middleware1.AdminMiddlewareHandler, userController.AddPermission)
	app.Post("/adminTestRoute", middleware1.AdminMiddlewareHandler, userController.TestRoute)
	app.Get("/admin", middleware1.AdminMiddlewareHandler)                                                  //后台主页，展示销售数据,支付订单，未支付订单，数量和金钱，浏览数据统计
	app.Get("/admin/products", middleware1.AdminMiddlewareHandler, productController.AllProduct)           //展示后台产品数据
	app.Post("/admin/addproduct", middleware1.AdminMiddlewareHandler, productController.AddProduct)        //admin 添加产品
	app.Delete("/admin/delproduct/:id", middleware1.AdminMiddlewareHandler, productController.DelProduct)  //admin 删除产品
	app.Put("/admin/editproduct/:id", middleware1.AdminMiddlewareHandler, productController.UpdateProduct) //admin 编辑产品

	app.Get("/admin/users", middleware1.AdminMiddlewareHandler, userController.AllUsers)       //展示后台用户数据
	app.Delete("/admin/users/:id", middleware1.AdminMiddlewareHandler, userController.DelUser) //admin删除用户
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
// add_flag 添加权限的权限
// admin_flag  管理员权限
