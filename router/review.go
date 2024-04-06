package router

import (
	"9Kicks/controller"

	"github.com/gin-gonic/gin"
)

func ReviewRegister(r *gin.Engine) {
	r.POST("/v1/reviews/create", controller.AddReview)
	r.GET("/v1/reviews", controller.GetReviewList)
	r.GET("/v1/user-reviews", controller.GetUserReviews)
}
