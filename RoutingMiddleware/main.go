package main

import (
	"github.com/gofiber/fiber/v2"
)

// Custom middleware function
func Logger(c *fiber.Ctx) error {
	// Perform tasks before the route handling function
	println("Middleware: Request received")

	// Continue to the next middleware or route handling function
	return c.Next()
}

func someFunctionThatMayFail() error {
	return fiber.ErrBadRequest
}

func main() {

	app := fiber.New()
	// Apply the custom Logger middleware to all routes
	app.Use(Logger)

	// Define a route for the root URL
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Fiber!")
	})

	// Define a route for /about
	app.Get("/about", func(c *fiber.Ctx) error {
		return c.SendString("About Fiber")
	})

	//localhost:3000/users/47
	// Define a dynamic route that captures a user's ID
	app.Get("/users/:id", func(c *fiber.Ctx) error {
		// Get the user ID from the route parameters
		userID := c.Params("id")
		return c.SendString("User ID: " + userID)
	})

	// Send the JSON response
	app.Get("/json", func(c *fiber.Ctx) error {
		// Create a JSON response
		response := fiber.Map{
			"message": "Hello, Fiber!",
		}
		return c.JSON(response)
	})

	// Define an error handler 捕获恐慌并用错误消息进行响应。这允许您优雅地处理错误并向客户端发送有意义的错误响应。
	/*
		使用 app.Use 注册了一个中间件，它接受一个 fiber.Ctx 类型的参数 c，这个参数代表当前的请求上下文，并且返回一个 error 类型

		defer 关键字用于定义一个延迟函数，这个函数将在包裹它的函数（在这个例子中是中间件函数）返回之后执行。这对于执行清理操作或捕获 panic 非常有用

		在 defer 调用的匿名函数内部，使用 recover 函数捕获 panic

		如果捕获到 panic，使用 c.Status(500).SendString 设置响应的状态码为 500（内部服务器错误），并发送字符串 "Internal Server Error" 作为响应体

	*/
	app.Use(func(c *fiber.Ctx) error {
		defer func() {
			if err := recover(); err != nil {
				// Handle the error and respond with an error message
				c.Status(500).SendString("Internal Server Error")
			}
		}() // }是匿名函数的结束括号，紧随其后的是 defer 调用的结束括号()

		return c.Next() //在 defer 调用之前，中间件函数返回 c.Next()，这允许 Fiber 继续执行链中的下一个中间件或处理函数。由于 defer 会在当前函数退出后才执行，所以即使 c.Next() 执行后发生 panic，defer 中的 recover 也能捕获到它
	}) //结束中间件定义

	app.Get("/error", func(c *fiber.Ctx) error {
		// Simulate an error
		err := someFunctionThatMayFail()
		if err != nil {
			// Trigger a panic with the error
			panic(err)
		}
		return c.SendString("No error occurred")
	})

	// Start the Fiber application
	app.Listen(":3000")

}
