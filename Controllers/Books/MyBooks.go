package books

import (
	"go_learn/database"
	"go_learn/models"

	"github.com/gin-gonic/gin"
)

func MyBooks(c *gin.Context) {

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

	var dbbooks []models.Book

	result := database.DB.Model(&models.Book{}).Preload("Author").Where("author_id = ?", authorID).Find(&dbbooks)

	if result.RowsAffected == 0 {
		c.JSON(404, gin.H{ "error" : "Books not found." })
		return
	}

	var res_book []BookResponse

	for _, book := range dbbooks {
		res_book = append(res_book, BookResponse{
			Title:       book.Title,
			Description: book.Description,
			Picture:     book.Picture,
			BookURL:     book.BookURL,
			BookSlug:    book.BookSlug,
			Author:      book.Author.Fullname,
			AuthorID:    book.AuthorID,
			CreatedAt:   book.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(200, res_book)
}