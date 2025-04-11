package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Chinedum-Ogbuagu/guardianGo-backend.git/internal/child"
	"github.com/Chinedum-Ogbuagu/guardianGo-backend.git/internal/church"
	"github.com/Chinedum-Ogbuagu/guardianGo-backend.git/internal/dropoff"
	"github.com/Chinedum-Ogbuagu/guardianGo-backend.git/internal/parent"
	"github.com/Chinedum-Ogbuagu/guardianGo-backend.git/internal/pickup"
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

	
	churchRepo := church.NewRepository()
	churchSvc := church.NewService(churchRepo)
	churchHandler := church.NewHandler(db, churchSvc)
	churchHandler.RegisterRoutes(r)
	

	parentRepo := parent.NewRepository()
	parentSvc := parent.NewService(parentRepo)
	parentHandler := parent.NewHandler(db, parentSvc)
	parentHandler.RegisterRoutes(r)

	childRepo := child.NewRepository()
	childSvc := child.NewService(childRepo)
	childHandler := child.NewHandler(db, childSvc)
	childHandler.RegisterRoutes(r)

	dropoffRepo := dropoff.NewRepository()
	dropOffSvc := dropoff.NewService(dropoffRepo)
	dropoffHandler := dropoff.NewHandler(db, dropOffSvc)
	dropoffHandler.RegisterRoutes(r)

	pickupRepo := pickup.NewRepository()
	pickupSvc := pickup.NewService(pickupRepo)
	pickupHandler := pickup.NewHandler(db, pickupSvc)
	pickupHandler.RegisterRoutes(r)

	r.Run()
}