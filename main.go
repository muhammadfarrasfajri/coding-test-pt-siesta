package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/muhammadfarrasfajri/coding-test-pt-siesta/bootstrap"
	route "github.com/muhammadfarrasfajri/coding-test-pt-siesta/router"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Inisialisasi Database
	bootstrap.InitDatabase()

	// Inisialisasi Dependency Injection
	container := bootstrap.InitContainer()

	// Inisialisasi Gin Engine
	r := gin.Default()

	// Setup Routes
	route.SetupRouter(
		r,
		container.TransferController,
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	//Jalankan Server
	r.Run(":" + port)
}
