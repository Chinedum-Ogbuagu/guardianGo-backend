package auth

import (
	"net/http"

	"github.com/Chinedum-Ogbuagu/guardianGo-backend.git/internal/otp"
	"github.com/Chinedum-Ogbuagu/guardianGo-backend.git/internal/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
	OTP otp.Service
	userRepo user.Repository
}

func NewHandler(db *gorm.DB, otpService otp.Service, userRepo user.Repository) *Handler {
	return &Handler{DB: db, OTP: otpService, userRepo: userRepo}
}


func (h *Handler) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/auth")
	{
		group.POST("/login", h.Login)
		group.POST("/register", h.Verify)
	}
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	

	_, err := h.userRepo.GetByPhone(h.DB, req.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}	

	_, err = h.OTP.SendOTP(h.DB, req.PhoneNumber, "login", 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}

func (h *Handler) Verify(c *gin.Context) {
	var req VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	valid, err := h.OTP.VerifyOTP(h.DB, req.PhoneNumber, req.Code)
	if err != nil || !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired code"})
		return
	}

	user, err := h.userRepo.GetByPhone(h.DB, req.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// In a real system, generate JWT or session
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user": gin.H{
			"id": user.ID,
			"name": user.Name,
			"role": user.Role,
			"church_id": user.ChurchID,
		},
	})
}
