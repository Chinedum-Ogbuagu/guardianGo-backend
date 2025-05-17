package otp

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
	return &Handler{
		DB:      db,
		Service: service,
	}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/api/otp")
	{
		group.POST("/send", h.SendOTP)
		group.POST("/verify", h.VerifyOTP)

	}
}

type SendOTPRequest struct {
	PhoneNumber string `json:"phone" binding:"required"`
	Purpose     string `json:"purpose" binding:"required"`
	DropOffID   uint   `json:"drop_off_id"` // optional, based on use case
}

func (h *Handler) SendOTP(c *gin.Context) {
	var req SendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	otp, err := h.Service.SendOTP(h.DB, req.PhoneNumber, req.Purpose, req.DropOffID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent", "pin_id": otp.PinID})
}

type VerifyOTPRequest struct {
	Phone   string `json:"phone" binding:"required"`
	Code    string `json:"code" binding:"required"`
	Purpose string `json:"purpose"` // optional, based on use case
}

func (h *Handler) VerifyOTP(c *gin.Context) {
	var req VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	valid, err := h.Service.VerifyOTP(h.DB, req.Phone, req.Code, req.Purpose)
	if err != nil || !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP verified"})
}
