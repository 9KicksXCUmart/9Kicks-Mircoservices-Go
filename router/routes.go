package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	route := gin.Default()

	// CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173", "https://9kicks.shop"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Access-Control-Allow-Methods"}
	config.AllowCredentials = true
	route.Use(cors.New(config))

	// register all routes.
	AuthRegister(route)
	ReviewRegister(route)
	ProductRegister(route)

	return route
}
