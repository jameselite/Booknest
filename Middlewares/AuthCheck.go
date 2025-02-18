package middlewares

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtaccessSecret := os.Getenv("JWT_ACCESS")

		if jwtaccessSecret == "" {
			c.JSON(500, gin.H{"error": "JWT secret is not set in the environment variables."})
			c.Abort()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Token is missing."})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(401, gin.H{"error": "Token format is invalid. Use 'Bearer <token>'."})
			c.Abort()
			return
		}

		claims := jwt.MapClaims{}
		token, tokenerr := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtaccessSecret), nil
		})

		if tokenerr != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "Invalid or expired token."})
			fmt.Println("JWT Parsing Error:", tokenerr) // Log error for debugging
			c.Abort()
			return
		}

		userClaims := make(map[string]interface{})

		for key, value := range claims {
			userClaims[key] = value
		}

		c.Set("reqclaim", userClaims)
		c.Next()
	}
}
