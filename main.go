package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Chinedum-Ogbuagu/guardianGo-backend.git/internal/church"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	

	fmt.Println("Database connection established")

	r := gin.Default()

	
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	repo := church.NewRepository()
	svc := church.NewService(repo)
	h := church.NewHandler(db,svc)
	h.RegisterRoutes(r)

	r.Run() 
}
