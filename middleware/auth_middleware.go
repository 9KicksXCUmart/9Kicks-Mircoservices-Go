/*
Package middleware restricts access to certain routes based on the user's authentication status.
*/
package middleware

import (
	"9Kicks/model/auth"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie("jwt")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenString := cookie.Value

		token, err := jwt.ParseWithClaims(tokenString, &auth.Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return token, nil
		})
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*auth.Claims)
		if !ok || !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Set the user's email in the Gin context for later use
		c.Set("userEmail", claims.Email)

		c.Next()
	}
}
