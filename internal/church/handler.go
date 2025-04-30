package church

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

func NewHandler(db *gorm.DB, s Service) *Handler {
	return &Handler{DB: db, Service: s}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/church")
	{
		group.POST("", h.CreateChurch)
		group.GET(":id", h.GetChurchByID)
		group.GET("", h.GetAllChurches) // Added route for GetAllChurches
		group.PUT(":id", h.UpdateChurch)
	}
}

func (h *Handler) CreateChurch(c *gin.Context) {
	var church Church
	if err := c.ShouldBindJSON(&church); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Service.Create(h.DB, &church); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, church)
}

func (h *Handler) GetChurchByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.FromString(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}
	church, err := h.Service.GetByID(h.DB, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "church not found"})
		return
	}
	c.JSON(http.StatusOK, church)
}

func (h *Handler) GetAllChurches(c *gin.Context) {
	churches, err := h.Service.GetAll(h.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch churches"})
		return
	}
	c.JSON(http.StatusOK, churches)
}

func (h *Handler) UpdateChurch(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.FromString(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	var updatedChurch Church
	if err := c.ShouldBindJSON(&updatedChurch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedChurch.ID = id // Set the ID from the URL parameter

	if err := h.Service.Update(h.DB, &updatedChurch); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update church: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedChurch)
}
