package router

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	route := gin.Default()

	// register all routes.
	AuthRegister(route)
	ReviewRegister(route)

	return route
}
