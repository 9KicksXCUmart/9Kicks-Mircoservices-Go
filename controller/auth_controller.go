/*
Package controller implements the functions for handling requests to the authentication endpoints.
*/
package controller

import (
	"9Kicks/config"
	. "9Kicks/model/auth"
	"9Kicks/service/auth"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var secretKey = config.GetJWTSecrets().JWTUserSecret

func Signup(c *gin.Context) {
	var user UserSignUpForm
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	exists, _ := auth.CheckEmailExists(user.Email)
	if exists {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": "Email already exists",
		})
		return
	}

	verificationToken, success := auth.CreateUser(user.Email, user.FirstName, user.LastName, user.Password)
	if !success {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create user",
		})
		return
	}

	// Send verification email
	err := auth.SendEmailTo(user.Email, verificationToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to send verification email",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "A verification email has been sent to your email address. Please verify your email to complete registration",
	})
}

func Login(c *gin.Context) {
	var user UserLoginForm

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	exists, _ := auth.CheckEmailExists(user.Email)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "This email is not registered",
		})
		return
	}

	userProfile, err := auth.GetUserProfileByEmail(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve user profile",
		})
		return
	}

	// Check if the password is correct
	isValidPassword := auth.IsValidPassword(userProfile.Password, user.Password)
	if !isValidPassword {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Wrong password",
		})
		return
	}

	// Check if the email is verified
	if !userProfile.IsVerified {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Email not verified",
		})
		return
	}

	parts := strings.Split(userProfile.PK, "#")
	tokenString, _, err := auth.GenerateJWT(secretKey, user.Email, parts[1])

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login successful",
		"token":   tokenString,
	})
}

func ValidateToken(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid Authorization header format",
		})
		return
	}

	tokenString := parts[1]

	claims, err := auth.DecodeJWT(secretKey, tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Token is valid",
		"data": gin.H{
			"email":   claims.Email,
			"user_id": claims.UserID,
		},
	})
}

func ResendVerificationEmail(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Missing email",
		})
		return
	}

	userProfile, err := auth.GetUserProfileByEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve user profile",
		})
		return
	}
	userId := userProfile.PK

	verificationToken, _, err := auth.UpdateVerificationToken(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update verification token",
		})
		return
	}

	err = auth.SendEmailTo(email, verificationToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to send verification email",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "A verification email has been sent to your email address. Please verify your email to complete registration",
	})
}

func VerifyEmail(c *gin.Context) {
	// Get the verification token and email from the request parameters
	var form EmailVerificationForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	token := form.Token
	email := form.Email
	if token == "" || email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing token or email"})
		return
	}

	userProfile, err := auth.GetUserProfileByEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve user profile",
		})
		return
	}

	if userProfile.IsVerified {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": "Email already verified",
		})
		return
	}

	storedToken := userProfile.VerificationToken
	tokenExpirationTime := userProfile.TokenExpiry

	if storedToken != token {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid verification token",
		})
		return
	}

	// Check if the token has expired
	if time.Now().Unix() > tokenExpirationTime {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Verification token has expired",
		})
		return
	}

	// Update the user profile to set isVerified to true
	err = auth.VerifyUserEmail(userProfile.PK)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update user profile",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Email verified successfully",
	})
}

func ForgotPassword(c *gin.Context) {
	email := c.Query("email")
	exists, _ := auth.CheckEmailExists(email)
	if !exists {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": "This email is not registered",
		})
		return
	}

	userProfile, err := auth.GetUserProfileByEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve user profile",
		})
		return
	}

	userId := userProfile.PK
	userName := userProfile.FirstName
	verificationToken, _, err := auth.UpdateVerificationToken(userId)
	success := auth.SendResetPasswordEmail(email, userName, verificationToken)
	if !success {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to reset password",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "A password reset email has been sent to your email address. Please follow the instructions in the email to reset your password",
	})
}

func ResetPassword(c *gin.Context) {
	var resetForm ResetPasswordForm
	if err := c.ShouldBindJSON(&resetForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	userProfile, err := auth.GetUserProfileByEmail(resetForm.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve user profile",
		})
		return
	}

	storedToken := userProfile.VerificationToken
	tokenExpirationTime := userProfile.TokenExpiry

	if storedToken != resetForm.Token {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid verification token",
		})
		return
	}

	// Check if the token has expired
	if time.Now().Unix() > tokenExpirationTime {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Verification token has expired",
		})
		return
	}

	success := auth.UpdatePassword(userProfile.PK, resetForm.NewPassword)
	if !success {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update password",
		})
		return
	}

	if !userProfile.IsVerified {
		// Update the user profile to set isVerified to true
		err = auth.VerifyUserEmail(userProfile.PK)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to update user profile",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Password reset successful",
	})
}
