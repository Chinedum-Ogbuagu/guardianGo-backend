package guardian

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	DB      *gorm.DB
	Service Service
	Repo 	Repository
}

func NewHandler(db *gorm.DB, s Service) *Handler {
	return &Handler{DB: db, Service: s}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/api/guardians")
	{
		group.POST("/", h.FindOrCreateGuardian)
		group.GET("/with-children/:phone", h.GetGuardianWithChildren) 
	}
}

type GuardianRequest struct {
	Name  string `json:"name" binding:"required"`
	Phone string `json:"phone" binding:"required"`
}

func (h *Handler) FindOrCreateGuardian(c *gin.Context) {
	var req GuardianRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	guardian, err := h.Service.FindOrCreateGuardian(h.DB, req.Name, req.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create or find guardian"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"guardian": guardian})
}
func (h *Handler ) GetGuardianWithChildren(c *gin.Context) {
	phone := c.Param("phone")

	guardian, err := h.Service.FindGuardianByPhone(h.DB, phone)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Guardian not found"})
		return
	}

	children, err := h.Service.GetChildrenByGuardianPhone(h.DB, phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve children"})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"guardian": gin.H{
			"name":  guardian.Name,
			"phone": guardian.Phone,
		},
		"children": children,
	})
}