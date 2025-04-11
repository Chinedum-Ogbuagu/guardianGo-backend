package pickup

import (
	"net/http"
	"strconv"

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
	group := r.Group("/pickup")
	{
		group.POST("/confirm", h.ConfirmPickup)
		group.GET("/:id", h.GetPickupByID)
	}
}

func (h *Handler) ConfirmPickup(c *gin.Context) {
	var pickup Pickup
	if err := c.ShouldBindJSON(&pickup); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Service.Create(h.DB, &pickup); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, pickup)
}

func (h *Handler) GetPickupByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pickup ID"})
		return
	}
	pickup, err := h.Service.GetByID(h.DB, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "pickup not found"})
		return
	}
	c.JSON(http.StatusOK, pickup)
}
