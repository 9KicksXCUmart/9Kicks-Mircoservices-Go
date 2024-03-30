package router

import (
	"github.com/gin-gonic/gin"
)

// SetupRouter 配置路由
func SetupRouter() *gin.Engine {
	route := gin.Default()

	// register all route.
	AuthRegister(route)

	return route
}
