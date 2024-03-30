package auth_model

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	jwt.RegisteredClaims
	Email  string `json:"email"`
	UserID string `json:"userID"`
}
