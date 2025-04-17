package child

import (
	"fmt"
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
	group := r.Group("/api/children")
	{
		group.POST("/", h.FindOrCreateChild)
		group.GET("/guardian/:id", h.GetChildrenByGuardian)
	}
}

type ChildRequest struct {
	Name       string `json:"name" binding:"required"`
	GuardianID uint   `json:"guardian_id" binding:"required"`
}

func (h *Handler) FindOrCreateChild(c *gin.Context) {
	var req ChildRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	child, err := h.Service.FindOrCreateChild(h.DB, req.Name, req.GuardianID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create/find child"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"child": child})
}

func (h *Handler) GetChildrenByGuardian(c *gin.Context) {
	id := c.Param("id")

	var guardianID uint
	if _, err := fmt.Sscanf(id, "%d", &guardianID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid guardian ID"})
		return
	}

	children, err := h.Service.GetChildrenByGuardian(h.DB, guardianID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch children"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"children": children})
}
