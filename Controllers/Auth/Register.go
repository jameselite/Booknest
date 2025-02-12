package auth

import (
	"errors"
	"go_learn/database"
	"go_learn/helpers"
	"go_learn/models"
	"regexp"

	"github.com/gin-gonic/gin"
)

func validateRegister(usr models.User) error {
	if usr.Email == "" || usr.Fullname == "" || usr.Password == "" {
		return errors.New("requested data cannot be empty")
	}

	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)

	if !re.MatchString(usr.Email) {
		return errors.New("invalid email format")
	}

	if len(usr.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	return nil
}

func RegisterUser(c *gin.Context) {

	var user models.User

	if err := c.ShouldBindBodyWithJSON(&user); err != nil {
		c.JSON(400, gin.H{ "error" : err.Error() })
		return
	}

	if validaterr := validateRegister(user); validaterr != nil {
		c.JSON(400, gin.H{ "error" : validaterr.Error() })
		return
	}

	hashedpass := helpers.HashString(user.Password)

	user.Password = hashedpass

	database.DB.Create(&user)

	c.JSON(201, gin.H{ 
		"email" : user.Email,
		"fullname" : user.Fullname,
	})

}