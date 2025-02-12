package auth

import (
	"log"
	"os"
	"time"

	"go_learn/database"
	"go_learn/helpers"
	"go_learn/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func LoginUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	var dbUser models.User
	if err := database.DB.Where("email = ?", user.Email).First(&dbUser).Error; err != nil {
		log.Printf("Database error: %v", err)
		c.JSON(404, gin.H{"error": "Wrong email or password"})
		return
	}

	hashedInputPassword := helpers.HashString(user.Password)

	if hashedInputPassword != dbUser.Password {
		c.JSON(400, gin.H{"error": "Wrong email or password"})
		return
	}

	jwtaccessSecret := os.Getenv("JWT_ACCESS")
	jwtrefreshSecret := os.Getenv("JWT_REFRESH")

	accesstoken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": dbUser.Email,
		"id":    dbUser.ID,
		"exp":   time.Now().Add(time.Hour * 1).Unix(),
	})

	refreshtoken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": dbUser.Email,
		"id":    dbUser.ID,
		"exp":   time.Now().AddDate(0, 3, 0).Unix(),
	})

	accesstokenString, erraccess := accesstoken.SignedString([]byte(jwtaccessSecret))
	if erraccess != nil {
		log.Printf("Failed to generate access token: %v", erraccess)
		c.JSON(500, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshtokenString, errrefresh := refreshtoken.SignedString([]byte(jwtrefreshSecret))
	if errrefresh != nil {
		log.Printf("Failed to generate refresh token: %v", errrefresh)
		c.JSON(500, gin.H{"error": "Failed to generate refresh token"})
		return
	}
	c.JSON(200, gin.H{
	"accesstoken":  accesstokenString,
	"refreshtoken": refreshtokenString,
	})
}