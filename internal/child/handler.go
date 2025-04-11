package child

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
	Service Service
}

func NewHandler(db *gorm.DB, s Service) *Handler {
	return &Handler{DB: db, Service: s}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/child")
	{
		group.POST("/register", h.CreateChild)
		group.GET("/:id", h.GetChildByID)
	}
}

func (h *Handler) CreateChild(c *gin.Context) {
	var child Child
	if err := c.ShouldBindJSON(&child); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return 
	}
	if err := h.Service.Create(h.DB, &child); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return 
	}
	c.JSON(http.StatusCreated, child)

}

func (h *Handler) GetChildByID(c *gin.Context) {
	id,err := strconv.Atoi(c.Param("id"))
	if err !=nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : "invalid child ID"})
		return 
	}
	child, err := h.Service.GetByID(h.DB, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "child not found"})
		return 
	}
	c.JSON(http.StatusOK, child)
}