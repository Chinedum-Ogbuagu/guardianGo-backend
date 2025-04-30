package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
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
	group := r.Group("/api/users")
	{
		group.POST("/register", h.RegisterUser)
		group.GET("/by-phone/:phone", h.GetUserByPhone)
	}
}

type RegisterUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Role     Role   `json:"role" binding:"required"`
	ChurchID uuid.UUID   `json:"church_id" binding:"required"`
}

func (h *Handler) RegisterUser(c *gin.Context) {
	var req RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Service.RegisterUser(h.DB, req.Name, req.Phone, req.Role, req.ChurchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) GetUserByPhone(c *gin.Context) {
	phone := c.Param("phone")
	user, err := h.Service.GetUserByPhone(h.DB, phone)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}
