package router

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	route := gin.Default()

	// CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173", "https://9kicks.shop"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	route.Use(cors.New(config))
  
	log.Println("CORS enabled for frontend host: ", os.Getenv("FRONTEND_HOST"))

	// register all routes.
	AuthRegister(route)
	ReviewRegister(route)
	ProductRegister(route)

	return route
}
