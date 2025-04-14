package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Chinedum-Ogbuagu/guardianGo-backend.git/internal/child"
	"github.com/Chinedum-Ogbuagu/guardianGo-backend.git/internal/church"
	"github.com/Chinedum-Ogbuagu/guardianGo-backend.git/internal/dropoff"
	"github.com/Chinedum-Ogbuagu/guardianGo-backend.git/internal/guardian"
	"github.com/Chinedum-Ogbuagu/guardianGo-backend.git/internal/otp"
	"github.com/Chinedum-Ogbuagu/guardianGo-backend.git/internal/pickup"
	"github.com/Chinedum-Ogbuagu/guardianGo-backend.git/internal/security"
	"github.com/Chinedum-Ogbuagu/guardianGo-backend.git/internal/user"
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

	db.AutoMigrate(
	&church.Church{},
	&guardian.Guardian{},
	&child.Child{},
	&dropoff.DropOff{},
	&pickup.Pickup{},
	&user.User{},
	&otp.OTPRequest{},
	)

	r := gin.Default()
	
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	
	churchRepo := church.NewRepository()
	churchSvc := church.NewService(churchRepo)
	churchHandler := church.NewHandler(db, churchSvc)
	churchHandler.RegisterRoutes(r)
	
	parentRepo := guardian.NewRepository()
	parentSvc := guardian.NewService(parentRepo)
	parentHandler := guardian.NewHandler(db, parentSvc)
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

	userRepo := user.NewRepository()
	userSvc := user.NewService(userRepo)
	userHandler := user.NewHandler(db, userSvc)
	userHandler.RegisterRoutes(r)


	secRepo := security.NewRepository()
	secSvc := security.NewService(secRepo)
	secHandler := security.NewHandler(db, secSvc)
	secHandler.RegisterRoutes(r)

	otpRepo := otp.NewRepository()
	otpSvc := otp.NewService(otpRepo)
	otpHandler := otp.NewHandler(db, otpSvc)
	otpHandler.RegisterRoutes(r)

	r.Run()
}