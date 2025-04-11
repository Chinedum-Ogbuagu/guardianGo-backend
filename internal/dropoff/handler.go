package dropoff

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

func NewHandler (db *gorm.DB, s Service) *Handler {
	return &Handler{DB: db, Service: s}

}

func (h *Handler) RegisterRoutes (r *gin.Engine) {
	group := r.Group("/dropoff")
	{
		group.POST("/create", h.CreateDropOff)
		group.GET("/:id", h.GetDropOffByID)
	}
}

func (h *Handler) CreateDropOff (c *gin.Context) {
	var drop DropOff 
	if err := c.ShouldBindJSON(&drop); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return 
	}
	if err := h.Service.Create(h.DB, &drop); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, drop)
	
}

func (h *Handler) GetDropOffByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}
	drop, err := h.Service.GetByID(h.DB, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "drop-off not found"})
		return
	}
	c.JSON(http.StatusOK, drop)
}