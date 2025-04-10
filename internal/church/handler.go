package church

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

func ( h *Handler) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/church")
	{
		group.POST("", h.CreateChurch)
		group.GET(":id", h.GetChurchByID)
	}
}

func (h *Handler) CreateChurch(c *gin.Context) {
	var church Church 
	if err := c.ShouldBindJSON(&church); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return 
	}
	if err := h.Service.Create(h.DB, &church); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error" : err.Error()})
		return
	}
	c.JSON(http.StatusCreated, church)

}

func (h *Handler) GetChurchByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}
	church, err := h.Service.GetByID(h.DB, uint(id))
	if err != nil {
	c.JSON(http.StatusNotFound, gin.H{"error": "church not found"})
	return
}
	c.JSON(http.StatusOK, church)

}