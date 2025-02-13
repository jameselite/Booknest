package auth

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func NewAccessToken(c *gin.Context) {

	jwtaccessSecret := os.Getenv("JWT_ACCESS")
	jwtrefreshSecret := os.Getenv("JWT_REFRESH")

	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.JSON(400, gin.H{ "error" : "Token is missing." })
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	if tokenString == authHeader {
		c.JSON(400, gin.H{ "error" : "Token format is invalid." })
		return
	}

	claims := jwt.MapClaims{}
	token, tokenerr := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtrefreshSecret), nil
	})

	if tokenerr != nil || !token.Valid {
		c.JSON(400, gin.H{ "error" : tokenerr.Error() })
		return		
	}

	accesstoken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": claims["email"],
		"id":    claims["id"],
		"exp":   time.Now().Add(time.Hour * 1).Unix(),
	})

	accesstokenString, erraccess := accesstoken.SignedString([]byte(jwtaccessSecret))
	if erraccess != nil {
		log.Printf("Failed to generate access token: %v", erraccess)
		c.JSON(500, gin.H{"error": "Failed to generate access token"})
		return
	}

	c.JSON(200, gin.H{ "accesstoken" : accesstokenString })
}