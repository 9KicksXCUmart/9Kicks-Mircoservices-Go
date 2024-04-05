package router

import (
	"9Kicks/controller"

	"github.com/gin-gonic/gin"
)

func ProductRegister(r *gin.Engine) {
	r.GET("/v1/products", controller.FilterProducts)
	r.POST("/v1/products/create", controller.PublishProduct)
	r.GET("/v1/products/:id", controller.GetProductDetailByID)
}
