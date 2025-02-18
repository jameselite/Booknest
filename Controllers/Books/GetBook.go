package books

import (
	"go_learn/database"
	"go_learn/models"

	"github.com/gin-gonic/gin"
)

func GetBook(c *gin.Context) {

	bookparam := c.Param("id")

	if bookparam == "" {
		c.JSON(400, gin.H{ "error": "Empty or invalid url param."})
		return
	}

	var book models.Book

	if err := database.DB.Model(&models.Book{}).Preload("Author").Where("book_slug = ?", bookparam).First(&book).Error; err != nil {
		c.JSON(404, gin.H{ "error": "Book not found" })
		return
	}

	res_book := BookResponse{
		Title: book.Title,
		Description: book.Description,
		Picture: book.Picture,
		BookURL: book.BookURL,
		BookSlug: book.BookSlug,
		Author: book.Author.Fullname,
		AuthorID: book.AuthorID,
		CreatedAt: book.CreatedAt.String(),
	}

	c.JSON(200, res_book)
}