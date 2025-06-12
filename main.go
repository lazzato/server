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

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "https://dashboard.thebkht.com")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

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
	auth.GET("/me", controllers.GetMeHandler)

	r.Run() // listen and serve on
}
