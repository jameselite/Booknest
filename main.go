package main

import (
	"fmt"
	"go_learn/Controllers/Auth"
	books "go_learn/Controllers/Books"
	middlewares "go_learn/Middlewares"
	"go_learn/database"
	"go_learn/models"

	"github.com/gin-gonic/gin"
)

func main() {

    database.ConnectDatabase()

    database.DB.AutoMigrate(&models.User{}, &models.Book{})

    app := gin.Default()

    app.POST("/api/register", auth.RegisterUser)
    app.POST("/api/login", auth.LoginUser)
    app.POST("/api/newtoken", auth.NewAccessToken)
    app.POST("/api/newbook", middlewares.AuthCheck(), books.CreateBooks)
    app.GET("/api/books", books.GetAllBooks)
    app.GET("/api/books/:id", books.GetBook)
    app.PUT("/api/books/:id/update", middlewares.AuthCheck(), books.UpdateBook)
    app.DELETE("/api/books/:id/delete", middlewares.AuthCheck(), books.DeleteBook)
    app.POST("/api/mybooks", middlewares.AuthCheck(), books.MyBooks)
    app.GET("/api/profile/:id", books.UsersBooks)

    port := ":3000"

    if err := app.Run(port); err != nil {

        fmt.Println("There is a problem running the server : ", err.Error())
    }

    runserver := fmt.Sprintf("Server is running on port %s .", port)

    fmt.Println(runserver)
}