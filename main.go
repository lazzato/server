package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lazzato/server/config"
	"github.com/lazzato/server/controllers"
	"github.com/lazzato/server/initializers"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDatabase()
	config.InitGoogleOAuth()
}

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	r.GET("/auth/google/login", controllers.GoogleAuth)
	r.GET("/auth/google/callback", controllers.GoogleAuthCallback)
	r.POST("/auth/refresh", controllers.RefreshAccessToken)

	auth := r.Group("/api")
	auth.Use(config.AuthMiddleware())
	auth.GET("/me", func(c *gin.Context) {
		userID := c.MustGet("userID").(uint)
		c.JSON(200, gin.H{
			"message": "Hello, World!",
			"userID":  userID,
		})
	})

	r.Run() // listen and serve on
}
