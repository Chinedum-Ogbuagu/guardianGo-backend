package pickup

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/Chinedum-Ogbuagu/guardianGo-backend.git/internal/dropoff"
)

type Handler struct {
	DB      *gorm.DB
	Service Service
	DropRepo dropoff.Repository
}

func NewHandler(db *gorm.DB, s Service, dropRepo dropoff.Repository) *Handler {
	return &Handler{DB: db, Service: s, DropRepo: dropRepo}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/api/pickup")
	{
		group.POST("/session", h.ConfirmPickup)
		group.POST("/confirm/:dropSessionID", h.ConfirmPickup)
		group.GET("/session/:dropSessionID", h.GetPickupSession)
	}
}
type ConfirmPickupRequest struct {
	DropSessionID uint   `json:"drop_session_id"`
	GuardianID    uint   `json:"guardian_id"`
	VerifiedByID  uint   `json:"verified_by_id"`
	Notes         string `json:"notes"`
}

func (h *Handler) ConfirmPickup(c *gin.Context) {
	var req ConfirmPickupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session, err := h.Service.ConfirmPickupSession(h.DB, req.DropSessionID, req.GuardianID, req.VerifiedByID, req.Notes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, session)
}

func (h *Handler) GetPickupSession(c *gin.Context) {
	dropSessionID, err := strconv.ParseUint(c.Param("dropSessionID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	session, err := h.Service.GetPickupSessionByDropSessionID(h.DB, uint(dropSessionID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pickup session not found"})
		return
	}

	c.JSON(http.StatusOK, session)
}
