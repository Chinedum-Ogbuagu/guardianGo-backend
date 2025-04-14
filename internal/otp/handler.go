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

type SendOTPRequest struct {
    PhoneNumber string `json:"phone_number" binding:"required"`
    Purpose     string `json:"purpose" binding:"required"`
    DropOffID   uint   `json:"drop_off_id" binding:"required"`
}

type VerifyOTPRequest struct {
    PhoneNumber string `json:"phone_number" binding:"required"`
    Code        string `json:"code" binding:"required"`
}

func NewHandler(db *gorm.DB, service Service) *Handler {
    return &Handler{DB: db, Service: service}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
    group := r.Group("/otp")
    {
        group.POST("/send", h.SendOTP)
        group.POST("/verify", h.VerifyOTP)
    }
}

func (h *Handler) SendOTP(c *gin.Context) {
    var req SendOTPRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    otpRequest, err := h.Service.SendOTP(h.DB, req.PhoneNumber, req.Purpose, req.DropOffID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "OTP sent successfully",
        "data": otpRequest,
    })
}

func (h *Handler) VerifyOTP(c *gin.Context) {
    var req VerifyOTPRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    isValid, err := h.Service.VerifyOTP(h.DB, req.PhoneNumber, req.Code)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    if !isValid {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "message": "Invalid OTP code",
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "OTP verified successfully",
    })
}