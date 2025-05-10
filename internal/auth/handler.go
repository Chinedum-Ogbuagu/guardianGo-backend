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
		group.POST("/logout", h.Logout)
	}
}

type RequestOTPBody struct {
	Phone     string `json:"phone" binding:"required"`
	Name      string `json:"name" binding:"required"`
	DropOffID uint   `json:"drop_off_id"` // optional
}

func (h *Handler) RequestOTP(c *gin.Context) {
	var req RequestOTPBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Service.RequestOTP(h.DB, req.Phone, req.Name, req.DropOffID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Login failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}

type VerifyOTPBody struct {
	Phone   string `json:"phone" binding:"required"`
	Code    string `json:"code" binding:"required"`
	Name    string `json:"name"`    // optional
	Purpose string `json:"purpose"` // optional
}

func (h *Handler) VerifyOTP(c *gin.Context) {
	var req VerifyOTPBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := h.Service.VerifyOTPAndLogin(h.DB, req.Phone, req.Code, req.Name, req.Purpose)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Set HttpOnly secure cookie with JWT token
	c.SetCookie(
		"auth_token",
		token,
		3600*24, // 1 day
		"/",
		"",    // domain
		false,  // secure (HTTPS)
		false,  // HttpOnly
	)

	c.JSON(http.StatusOK, gin.H{"message": "login successful", "user": user})
}

func (h *Handler) Logout(c *gin.Context) {
	// Invalidate the cookie
	c.SetCookie("auth_token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}
