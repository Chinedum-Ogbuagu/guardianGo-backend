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
	group := r.Group("/api/pickup")
	{
		// PickupSession endpoints
		group.POST("/session", h.CreatePickupSession)
		group.GET("/session/:id", h.GetPickupSessionByID)
		group.GET("/session/drop/:dropSessionID", h.GetPickupSessionByDropSessionID)
		group.POST("/validate", h.ValidatePickupCode)
		
		// Pickup-specific endpoints
		group.GET("/child/:childID/drop-session/:dropSessionID", h.GetPickupByChildAndDropSession)
	}
}

type CreatePickupSessionRequest struct {
	DropSessionID uint   `json:"drop_session_id" binding:"required"`
	GuardianID    uint   `json:"guardian_id" binding:"required"`
	VerifiedByID  uint   `json:"verified_by_id" binding:"required"`
	Notes         string `json:"notes"`
}

func (h *Handler) CreatePickupSession(c *gin.Context) {
	var req CreatePickupSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	result, err := h.Service.CreatePickupSession(
		h.DB, 
		req.DropSessionID, 
		req.GuardianID, 
		req.VerifiedByID, 
		req.Notes,
	)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) GetPickupSessionByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}
	
	pickupSession, err := h.Service.GetPickupSessionByID(h.DB, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "pickup session not found"})
		return
	}
	
	c.JSON(http.StatusOK, pickupSession)
}

func (h *Handler) GetPickupSessionByDropSessionID(c *gin.Context) {
	dropSessionID, err := strconv.ParseUint(c.Param("dropSessionID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid drop session ID"})
		return
	}
	
	pickupSession, err := h.Service.GetPickupSessionByDropSessionID(h.DB, uint(dropSessionID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "pickup session not found"})
		return
	}
	
	c.JSON(http.StatusOK, pickupSession)
}

type ValidatePickupCodeRequest struct {
	DropSessionID uint   `json:"drop_session_id" binding:"required"`
	Code          string `json:"code" binding:"required"`
}

func (h *Handler) ValidatePickupCode(c *gin.Context) {
	var req ValidatePickupCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	isValid, err := h.Service.ValidatePickupCode(h.DB, req.DropSessionID, req.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "valid": false})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"valid": isValid})
}

func (h *Handler) GetPickupByChildAndDropSession(c *gin.Context) {
	childID, err := strconv.ParseUint(c.Param("childID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid child ID"})
		return
	}
	
	dropSessionID, err := strconv.ParseUint(c.Param("dropSessionID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid drop session ID"})
		return
	}
	
	pickup, err := h.Service.GetPickupByChildAndDropSessionID(h.DB, uint(childID), uint(dropSessionID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "pickup not found"})
		return
	}
	
	c.JSON(http.StatusOK, pickup)
}