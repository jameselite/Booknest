package books

import (
	"go_learn/database"
	"go_learn/models"
	"github.com/gin-gonic/gin"
)

func UpdateBook(c *gin.Context) {

	bookSlug := c.Param("id")

	if bookSlug == "" {
		c.JSON(400, gin.H{"error": "Book slug cannot be empty."})
		return
	}

	var dbBook models.Book

	if err := database.DB.Where("book_slug = ?", bookSlug).Preload("Author").First(&dbBook).Error; err != nil {
		c.JSON(404, gin.H{"error": "Book not found."})
		return
	}

	claimsRaw, exists := c.Get("reqclaim")

	if !exists {
		c.JSON(400, gin.H{"error": "Unauthorized request."})
		return
	}

	claims, ok := claimsRaw.(map[string]interface{})

	if !ok {
		c.JSON(400, gin.H{"error": "Invalid claims format."})
		return
	}

	userIDFloat, ok := claims["id"].(float64)

	if !ok {
		c.JSON(400, gin.H{"error": "Invalid user ID."})
		return
	}
	userID := uint(userIDFloat)

	if dbBook.AuthorID != userID {
		c.JSON(403, gin.H{"error": "Access denied."})
		return
	}

	var input map[string]interface{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if len(input) == 0 {
		c.JSON(400, gin.H{"error": "At least one field must be updated."})
		return
	}

	if err := database.DB.Model(&dbBook).Updates(input).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update book."})
		return
	}

	res_book := BookResponse{

		Title: dbBook.Title,
		Description: dbBook.Description,
		Picture: dbBook.Picture,
		BookURL: dbBook.BookURL,
		BookSlug: dbBook.BookSlug,
		Author: dbBook.Author.Fullname,
		AuthorID: dbBook.Author.ID,
		CreatedAt: dbBook.CreatedAt.String(),
	}

	c.JSON(200, res_book)
}