package router

import (
	"9Kicks/controller"
	"github.com/gin-gonic/gin"
)

func ReviewRegister(r *gin.Engine) {
	r.POST("/v1/review/add-review", controller.AddReview)
	r.GET("/v1/review/get-review", controller.GetReviewList)
}
