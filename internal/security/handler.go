package security

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
	Service Service
}

func NewHandler(db *gorm.DB, s Service) *Handler {
	return &Handler{
		DB: db,
		Service: s,
	}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/security")
	{
		group.POST("/flag-issue", h.FlagIssue)
		group.GET("/logs", h.ListFlags)
	}
}

func (h *Handler) FlagIssue(c *gin.Context) {
	var flag SecurityFlag
	if err := c.ShouldBindJSON(&flag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Service.FlagIssue(h.DB, &flag); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, flag)
}

func  (h *Handler) ListFlags(c *gin.Context) {
	flags, err := h.Service.GetAllFlags(h.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, flags)
}