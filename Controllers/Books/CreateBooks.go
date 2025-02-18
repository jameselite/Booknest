package books

import (
	"errors"
	"go_learn/database"
	"go_learn/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

func ValidateCreate(book models.Book) error {

	if (book.Title == "" || book.Description == "" || book.BookURL == "" || book.Picture == "") {
		return errors.New("requested data can not be empty")
	}else {
		return nil
	}
}

func CreateBooks(c *gin.Context) {

	var book models.Book

	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if validerr := ValidateCreate(book); validerr != nil {
		c.JSON(400, gin.H{ "error" : validerr.Error() })
		return
	}
	
	reqclaimRaw, exists := c.Get("reqclaim")
	if !exists {
		c.JSON(401, gin.H{"error": "User authentication required."})
		return
	}

	reqclaim, ok := reqclaimRaw.(map[string]interface{})
	if !ok {
		c.JSON(401, gin.H{"error": "Invalid user claims format."})
		return
	}

	idFloat, idExists := reqclaim["id"].(float64)
	if !idExists {
		c.JSON(403, gin.H{"error": "User ID is missing or invalid."})
		return
	}
	authorID := uint(idFloat)

	_, emailExists := reqclaim["email"].(string)
	if !emailExists {
		c.JSON(403, gin.H{"error": "User email is missing or invalid."})
		return
	}

	var user models.User
	if err := database.DB.First(&user, authorID).Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found."})
		return
	}

	baseSlug := slug.Make(book.Title)
	
	var maxID uint

	if maxerr := database.DB.Model(&models.Book{}).Select("MAX(id)").Scan(&maxID).Error; maxerr != nil {
		c.JSON(500, gin.H{"error": maxerr.Error()})
		return
	}
		
	book.BookSlug = baseSlug + "-" + strconv.Itoa(int(maxID))
	book.AuthorID = authorID

	if err := database.DB.Table("books").Create(&book).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, book)
}