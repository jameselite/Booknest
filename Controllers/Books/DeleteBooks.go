package books

import (
	"go_learn/database"
	"go_learn/models"

	"github.com/gin-gonic/gin"
)

func DeleteBook(c *gin.Context) {

	bookparam := c.Param("id")

	if bookparam == "" {
		c.JSON(400, gin.H{ "error" : "Book slug not found." })
	}

	var bookdb models.Book

	if err := database.DB.Model(&models.Book{}).Preload("Author").Where("book_slug = ?", bookparam).First(&bookdb).Error; err != nil {
		c.JSON(404, gin.H{ "error" : "Book not found." })
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

	if bookdb.AuthorID != authorID {
		c.JSON(400, gin.H{ "error" : "Access denied." })
		return
	}

	if errdel := database.DB.Delete(&bookdb).Error; errdel != nil {
		c.JSON(400, gin.H{ "error" : errdel })
		return
	}

	c.JSON(200, gin.H{ "message" : "Book deleted successfully." })
}