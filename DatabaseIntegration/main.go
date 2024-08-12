package main

import (
	"github.com/gofiber/fiber/v2"

	"gorm.io/gorm"

	"gorm.io/driver/mysql"
)

// Define a struct that corresponds to the database table
type Album struct {
	ID     uint
	Title  string
	Artist string
	Price  float32
}

func main() {

	app := fiber.New()

	// Define the database configuration
	dsn := "admin:123456@tcp(127.0.0.1:3306)/recordings?parseTime=true"
	// Initialize the database connection
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the database")
	}

	// Define your routes and handle database operations

	// Insert a new user record
	// INSERT INTO `albums` (`title`,`artist`,`price`) VALUES ('John Doe','3sds0',33.2)
	newUser := Album{Title: "John Doe", Artist: "3sds0", Price: 33.2}
	result := db.Create(&newUser)
	if result != nil {
		println("faild to insert")
	}
	// ...

	app.Listen(":3000")

}
