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

    r := gin.Default()

    r.POST("/api/register", auth.RegisterUser)
    r.POST("/api/login", auth.LoginUser)
    r.POST("/api/newtoken", auth.NewAccessToken)
    r.POST("/api/newbook", middlewares.AuthCheck(), books.CreateBooks)

    port := ":3000"

    if err := r.Run(port); err != nil {
        fmt.Println("There is a problem running the server : ", err.Error())
    }

    runserver := fmt.Sprintf("Server is running on port %s .", port)

    fmt.Println(runserver)
}