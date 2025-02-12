package main

import (
	"fmt"
	"go_learn/Controllers/Auth"
	"go_learn/database"
	"go_learn/models"

	"github.com/gin-gonic/gin"
)

func main() {

    database.ConnectDatabase()

    database.DB.AutoMigrate(&models.User{})

    r := gin.Default()

    r.POST("/api/register", auth.RegisterUser)
    r.POST("/api/login", auth.LoginUser)
    
    port := ":3000"

    if err := r.Run(port); err != nil {
        fmt.Println("There is a problem running the server : ", err.Error())
    }

    runserver := fmt.Sprintf("Server is running on port %s .", port)

    fmt.Println(runserver)
}