package books

import (
	"go_learn/database"
	"go_learn/models"

	"github.com/gin-gonic/gin"
)

func UsersBooks(c *gin.Context) {
	
	userparam := c.Param("id")

	if userparam == "" {
		c.JSON(400, gin.H{ "error" : "User param is missing." })
		return
	}

	var dbbooks []models.Book

	result := database.DB.Model(&models.Book{}).Preload("Author").Where("author_id = ?", userparam).Find(&dbbooks)

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