package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Chinedum-Ogbuagu/guardianGo-backend.git/internal/auth"
	"github.com/Chinedum-Ogbuagu/guardianGo-backend.git/internal/child"
	"github.com/Chinedum-Ogbuagu/guardianGo-backend.git/internal/church"
	"github.com/Chinedum-Ogbuagu/guardianGo-backend.git/internal/dropoff"
	"github.com/Chinedum-Ogbuagu/guardianGo-backend.git/internal/guardian"
	"github.com/Chinedum-Ogbuagu/guardianGo-backend.git/internal/otp"
	"github.com/Chinedum-Ogbuagu/guardianGo-backend.git/internal/pickup"
	"github.com/Chinedum-Ogbuagu/guardianGo-backend.git/internal/security"
	"github.com/Chinedum-Ogbuagu/guardianGo-backend.git/internal/user"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


func seedChurch(db *gorm.DB) uuid.UUID {
    existing := church.Church{}
    if err := db.First(&existing).Error; err == nil {
        log.Println("✅ Church already exists, skipping seeding.")
        return existing.ID
    }

    churchID, err := uuid.NewV4()
    if err != nil {
        log.Fatalf("❌ Failed to generate UUID for church: %v", err)
    }

    log.Printf("Generated Church UUID: %s", churchID.String())

    newChurch := church.Church{
        ID:        churchID,
        Name:      "Living Word Church",
        Address:   "123 Grace Avenue",
        CreatedAt: time.Now(),
    }

    if err := db.Create(&newChurch).Error; err != nil {
        log.Fatalf("❌ Failed to seed church: %v", err)
    }

    log.Printf("✅ Seeded default church with ID: %s", churchID.String())
    return churchID
}

func seedSuperAdmin(db *gorm.DB, churchID uuid.UUID) {
    existing := user.User{}
    if err := db.Where("role = ?", user.RoleSuperAdmin).First(&existing).Error; err == nil {
        log.Println("✅ Super admin already exists, skipping seeding.")
        return
    }

    userID, err := uuid.NewV4()
    if err != nil {
        log.Fatalf("❌ Failed to generate UUID for super admin: %v", err)
    }

    newUser := user.User{
        ID:        userID,
        Name:      "Chinedum Ogbuagu",
        Phone:     "09058652947",
        Role:      user.RoleSuperAdmin,
        ChurchID:  &churchID,
        CreatedAt: time.Now(),
    }

    if err := db.Create(&newUser).Error; err != nil {
        log.Fatalf("❌ Failed to seed super admin: %v", err)
    }

    log.Printf("✅ Seeded super admin with ID: %s", userID.String())
}

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

 if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
        log.Fatalf("Failed to create uuid-ossp extension: %v", err)
    }
    fmt.Println("UUID extension enabled")

	fmt.Println("Running AutoMigrations...")
	if err := db.AutoMigrate(&church.Church{},
		&auth.AuthRequest{},
		&auth.AuthSession{},
		&guardian.Guardian{},
		&child.Child{},
		&dropoff.DropSession{},
		&dropoff.DropOff{},		
		&pickup.PickupSession{},
		&pickup.Pickup{},
		&user.User{},
		&otp.OTPRequest{},
		&otp.OTPToken{},); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	fmt.Println("Migrations completed!")
	
	churchID := seedChurch(db)
    seedSuperAdmin(db, churchID)

	r := gin.Default()
	
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, 
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	

	
	churchRepo := church.NewRepository()
	churchSvc := church.NewService(churchRepo)
	churchHandler := church.NewHandler(db, churchSvc)
	churchHandler.RegisterRoutes(r)
	
	guardianRepo := guardian.NewRepository()
	guardianSvc := guardian.NewService(guardianRepo)
	guardianHandler := guardian.NewHandler(db, guardianSvc)
	guardianHandler.RegisterRoutes(r)

	childRepo := child.NewRepository()
	childSvc := child.NewService(childRepo)
	childHandler := child.NewHandler(db, childSvc)
	childHandler.RegisterRoutes(r)

	dropoffRepo := dropoff.NewRepository()
	dropOffSvc := dropoff.NewService(dropoffRepo, guardianRepo, childRepo)
	dropoffHandler := dropoff.NewHandler(db, dropOffSvc)
	dropoffHandler.RegisterRoutes(r)

	pickupRepo := pickup.NewRepository()
	pickupSvc := pickup.NewService(pickupRepo, dropoffRepo)
	pickupHandler := pickup.NewHandler(db, pickupSvc, dropoffRepo)
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

	
	authSvc := auth.NewService( otpSvc, userSvc)
	authHandler := auth.NewHandler(db, authSvc)
	authHandler.RegisterRoutes(r)

	r.Run()
}