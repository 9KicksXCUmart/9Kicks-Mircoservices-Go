package router

import (
	"9Kicks/controller"

	"github.com/gin-gonic/gin"
)

func ProductRegister(r *gin.Engine) {
	r.GET("/v1/products", controller.FilterProducts)
	r.POST("/v1/products/create", controller.PublishProduct)
	r.GET("/v1/products/:id", controller.GetProductDetailByID)
	r.GET("/v1/products/:id/stock", controller.GetStock)
	r.PUT("/v1/products/update-detail", controller.UpdateProductInfo)
	r.DELETE("/v1/products/:id", controller.DeleteProduct)
	r.PATCH("/v1/products/:id", controller.UpdateStock)
}
