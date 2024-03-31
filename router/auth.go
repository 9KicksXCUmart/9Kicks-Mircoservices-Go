package router

import (
	"9Kicks/controller"

	"github.com/gin-gonic/gin"
)

func AuthRegister(r *gin.Engine) {
	r.POST("/v1/auth/signup", controller.Signup)
	r.POST("/v1/auth/verify-email", controller.VerifyEmail)
	r.POST("/v1/auth/resend-email", controller.ResendVerificationEmail)
	r.POST("/v1/auth/login", controller.Login)
	r.POST("/v1/auth/validate-token", controller.ValidateToken)
	r.POST("/v1/auth/forgot-password", controller.ForgotPassword)
	r.POST("/v1/auth/reset-password", controller.ResetPassword)
}
