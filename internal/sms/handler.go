package sms

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
	group := r.Group("/api/sms")
	{
		group.POST("/send", h.SendSMS)
	}
}

type SendSMSRequest struct {
	PhoneNumber string `json:"phone" binding:"required"`
	Message     string `json:"message" binding:"required"`
}

func (h *Handler) SendSMS(c *gin.Context) {
	var req SendSMSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Service.SendSMS(h.DB, req.PhoneNumber, req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SMS sent successfully"})
}
