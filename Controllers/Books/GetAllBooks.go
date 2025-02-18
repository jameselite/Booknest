package books

import (
	"go_learn/database"
	"go_learn/models"

	"github.com/gin-gonic/gin"
)

type BookResponse struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Picture     string `json:"picture"`
	BookURL     string `json:"book_url"`
	BookSlug    string `json:"book_slug"`
	Author      string `json:"author"`
	AuthorID    uint   `json:"author_id"`
	CreatedAt   string `json:"created_at"`
}

func GetAllBooks(c *gin.Context) {

	var allBooks []models.Book

	if err := database.DB.Preload("Author").Select("title", "description", "picture", "book_url", "book_slug", "author_id").Find(&allBooks).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var res_book []BookResponse

	for _, book := range allBooks {
		res_book = append(res_book, BookResponse{
			Title:       book.Title,
			Description: book.Description,
			Picture:     book.Picture,
			BookURL:     book.BookURL,
			BookSlug:    book.BookSlug,
			Author:      book.Author.Fullname,
			AuthorID:    book.AuthorID,
			CreatedAt:   book.CreatedAt.Format("2006-01-02 15:04:05"), // Formats timestamp to readable string
		})
	}

	c.JSON(200, res_book)
}