package middlewares

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claim struct {
	Email string
	ID uint
}

func AuthCheck() gin.HandlerFunc {

	return func(c *gin.Context) {

		jwtaccessSecret := os.Getenv("JWT_ACCESS")

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(400, gin.H{ "error" : "Token is missing." })
			c.Abort()
			return
		}
		
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	
		if tokenString == authHeader {
			c.JSON(400, gin.H{ "error" : "Token format is invalid." })
			c.Abort()
			return
		}
	
		claims := jwt.MapClaims{}

		token, tokenerr := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtaccessSecret), nil
		})
	
		if tokenerr != nil || !token.Valid {
			c.JSON(400, gin.H{ "error" : tokenerr.Error() })
			c.Abort()
			return
		}
		reqclaim := Claim{ Email: claims["email"].(string), ID: claims["id"].(uint) }

		c.Set("reqclaim", reqclaim)
		c.Next()
	}
}