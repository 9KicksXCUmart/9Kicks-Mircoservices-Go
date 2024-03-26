package router

import (
	"9Kicks/controller"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/signup", controller.Signup)
	r.POST("/login", controller.Login)
	r.GET("/validate", controller.ValidateToken)

	return r
}
