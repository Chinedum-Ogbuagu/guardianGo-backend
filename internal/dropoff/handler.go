package dropoff

import (
	"net/http"
	"strconv"
	"time"

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
	group := r.Group("/api/dropoff")
	{
		group.POST("/session", h.CreateDropSession)
		group.GET("/session/:id", h.GetDropSessionByID)
		group.GET("/session/code/:code", h.GetDropSessionByCode)
		group.GET("/sessions", h.GetDropSessionsByDate)
		group.POST("/confirm/:id", h.ConfirmPickup) 
		group.PUT("/session/:id/image", h.UpdateDropSessionImage)


		group.GET("/session/:id/children", h.GetDropOffsBySessionID)
		group.GET("/child/:id", h.GetDropOffByID)
	}
}

type GuardianPayload struct {
	Name  string `json:"name" binding:"required"`
	Phone string `json:"phone" binding:"required"`
}

type ChildPayload struct {
	Name     string `json:"name" binding:"required"`
	Class    string `json:"class" binding:"required"`
	Bag      bool   `json:"bag"`
	Note     string `json:"note"`
}
type Pagination struct {
	Page     int `form:"page,default=0"`
	PageSize int `form:"page_size,default=0"`
}
type UpdateImageRequest struct {
	PhotoURL string `json:"photo_url" binding:"required"`
}


type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	TotalCount int64       `json:"total_count"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
}
type CreateDropSessionRequest struct {
	ChurchID  *uuid.UUID          `json:"church_id" binding:"required"`
	Note     string           `json:"note"`
	Guardian GuardianPayload  `json:"guardian" binding:"required"`
	Children []ChildPayload   `json:"children" binding:"required"`
}

func (h *Handler) CreateDropSession(c *gin.Context) {
	var req CreateDropSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.Service.CreateDropSession(
		h.DB,
		req,
	)
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
	dateParam := c.Query("date")
	if dateParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required date parameter (YYYY-MM-DD)"})
		return
	}

	parsedDate, err := time.Parse("2006-01-02", dateParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use<ctrl98>-MM-DD"})
		return
	}

	sessions, err := h.Service.GetDropSessionByCode(h.DB, parsedDate, code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "drop session not found"})
		return
	}

	c.JSON(http.StatusOK, sessions)
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

func (h *Handler) GetDropSessionsByDate(c *gin.Context) {
	var pagination Pagination
	if err := c.ShouldBindQuery(&pagination); err != nil {
		// If binding fails, it means pagination parameters might be missing
		// We will proceed with default/zero values for pagination,
		// which will fetch all records.
		pagination = Pagination{} // Initialize with default zero values
	}

	dateParam := c.Query("date")
	if dateParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required date parameter (YYYY-MM-DD)"})
		return
	}

	parsedDate, err := time.Parse("2006-01-02", dateParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use<ctrl98>-MM-DD"})
		return
	}

	sessions, totalCount, err := h.Service.GetDropSessionsByDate(h.DB, parsedDate, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Conditionally format the response based on whether pagination was used
	if pagination.Page > 0 && pagination.PageSize > 0 {
		c.JSON(http.StatusOK, PaginatedResponse{
			Data:       sessions,
			TotalCount: totalCount,
			Page:       pagination.Page,
			PageSize:   pagination.PageSize,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{"data": sessions})
	}
}
func (h *Handler) ConfirmPickup(c *gin.Context) {
	idParam := c.Param("id")
	sessionID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session ID"})
		return
	}

	err = h.Service.MarkDropSessionPickedUp(h.DB, uint(sessionID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to confirm pickup"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pickup confirmed successfully"})
}

func (h *Handler) UpdateDropSessionImage(c *gin.Context) {
	sessionID := c.Param("id")

	var req UpdateImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Service.UpdateDropSessionImageURL(h.DB, sessionID, req.PhotoURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Photo URL updated successfully"})
}