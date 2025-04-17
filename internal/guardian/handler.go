package guardian

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	DB      *gorm.DB
	Service Service
}

func NewHandler(db *gorm.DB, s Service) *Handler {
	return &Handler{DB: db, Service: s}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/api/guardians")
	{
		group.POST("/", h.FindOrCreateGuardian)
	}
}

type GuardianRequest struct {
	Name  string `json:"name" binding:"required"`
	Phone string `json:"phone" binding:"required"`
}

func (h *Handler) FindOrCreateGuardian(c *gin.Context) {
	var req GuardianRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	guardian, err := h.Service.FindOrCreateGuardian(h.DB, req.Name, req.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create or find guardian"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"guardian": guardian})
}
