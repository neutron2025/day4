package main

import (
	. "gin-demo/src"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	v1 := router.Group("/v1")
	AddUserRouter(v1)
	// r := gin.Default()
	// r.GET("/hello", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message":  "hello world",
	// 		"message2": "Success",
	// 	})
	// })

	// r.POST("/ping/:id", func(ctx *gin.Context) {
	// 	id := ctx.Param("id")
	// 	ctx.JSON(200, gin.H{
	// 		"id": id,
	// 	})
	// })

	router.Run(":8000")
}
