package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	DB      *gorm.DB
	Service Service
}

func NewHandler(db *gorm.DB, service Service) *Handler {
	return &Handler{DB: db, Service: service}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/api/auth")
	{
		group.POST("/request-otp", h.RequestOTP)
		group.POST("/verify-otp", h.VerifyOTP)
	}
}

type RequestOTPBody struct {
	Phone string `json:"phone" binding:"required"`
	Name  string `json:"name" binding:"required"`
	DropOffID uint   `json:"drop_off_id"` // optional, based on use case
}

func (h *Handler) RequestOTP(c *gin.Context) {
	var req RequestOTPBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Service.RequestOTP(h.DB, req.Phone, req.Name, req.DropOffID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}

type VerifyOTPBody struct {
	Phone string `json:"phone" binding:"required"`
	Code  string `json:"code" binding:"required"`
	Name string `json:"name"` // optional, based on use case
	Purpose string `json:"purpose"` // optional, based on use case
}

func (h *Handler) VerifyOTP(c *gin.Context) {
	var req VerifyOTPBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Service.VerifyOTPAndLogin(h.DB, req.Phone, req.Code, req.Name, req.Purpose)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
