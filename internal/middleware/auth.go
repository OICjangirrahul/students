package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/OICjangirrahul/students/internal/config"
	"github.com/OICjangirrahul/students/internal/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// AuthMiddleware creates a middleware for JWT authentication
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, response.GeneralError(fmt.Errorf("authorization header is required")))
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(cfg.JWT.Secret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, response.GeneralError(fmt.Errorf("invalid token: %v", err)))
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Add claims to context
			c.Set("userID", claims["sub"])
			c.Set("email", claims["email"])
			c.Set("role", claims["role"])
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, response.GeneralError(fmt.Errorf("invalid token claims")))
			c.Abort()
			return
		}
	}
}
