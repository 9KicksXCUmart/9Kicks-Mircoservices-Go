package router

import (
	"9Kicks/controller"

	"github.com/gin-gonic/gin"
)

func ProductRegister(r *gin.Engine) {
	r.GET("/v1/products", controller.SearchProducts)
}
