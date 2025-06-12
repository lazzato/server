package main

import (
	"github.com/lazzato/server/initializers"
	"github.com/lazzato/server/models"
)

func init() {
	// Load environment variables
	initializers.LoadEnvVariables()
	initializers.ConnectToDatabase()
}

func main() {
	initializers.DB.AutoMigrate(
		&models.User{},
		&models.Restaurant{},
	
	)
}
