package main

import (
	"jwt/go/controllers"
	"jwt/go/initializers"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()

}

func main() {
	r := gin.Default()

	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)

	r.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Server is up and running"})
	})

	r.Run()

}
