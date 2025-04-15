package dropoff

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
	group := r.Group("/api/dropoff")
	{
		
		group.POST("/session", h.CreateDropSession)
		group.GET("/session/:id", h.GetDropSessionByID)
		group.GET("/session/code/:code", h.GetDropSessionByCode)
		
		// Drop Off endpoints
		group.POST("/session/:id/child", h.AddChildToSession)
		group.GET("/session/:id/children", h.GetDropOffsBySessionID)
		group.GET("/child/:id", h.GetDropOffByID)
	}
}

type CreateDropSessionRequest struct {
	GuardianID   uint     `json:"guardian_id" binding:"required"`
	ChurchID     uint     `json:"church_id" binding:"required"`
	Note         string   `json:"note"`
	ChildrenIDs  []uint   `json:"children_ids" binding:"required,min=1"`
	Classes      []string `json:"classes" binding:"required"`
	BagStatuses  []bool   `json:"bag_statuses" binding:"required"`
}

func (h *Handler) CreateDropSession(c *gin.Context) {
	var req CreateDropSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	dropSession := DropSession{
		GuardianID: req.GuardianID,
		ChurchID:   req.ChurchID,
		Note:       req.Note,
	}
	
	result, err := h.Service.CreateDropSession(h.DB, &dropSession, req.ChildrenIDs, req.Classes, req.BagStatuses)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) GetDropSessionByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}
	
	dropSession, err := h.Service.GetDropSessionByID(h.DB, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "drop session not found"})
		return
	}
	
	c.JSON(http.StatusOK, dropSession)
}

func (h *Handler) GetDropSessionByCode(c *gin.Context) {
	code := c.Param("code")
	
	dropSession, err := h.Service.GetDropSessionByCode(h.DB, code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "drop session not found"})
		return
	}
	
	c.JSON(http.StatusOK, dropSession)
}

type AddChildRequest struct {
	ChildID    uint   `json:"child_id" binding:"required"`
	Class      string `json:"class" binding:"required"`
	BagStatus  bool   `json:"bag_status"`
	Note       string `json:"note"`
}

func (h *Handler) AddChildToSession(c *gin.Context) {
	sessionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session ID"})
		return
	}
	
	var req AddChildRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	dropOff, err := h.Service.AddChildToSession(h.DB, uint(sessionID), req.ChildID, req.Class, req.BagStatus, req.Note)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, dropOff)
}

func (h *Handler) GetDropOffsBySessionID(c *gin.Context) {
	sessionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session ID"})
		return
	}
	
	dropOffs, err := h.Service.GetDropOffsBySessionID(h.DB, uint(sessionID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, dropOffs)
}

func (h *Handler) GetDropOffByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}
	
	dropOff, err := h.Service.GetDropOffByID(h.DB, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "drop-off not found"})
		return
	}
	
	c.JSON(http.StatusOK, dropOff)
}