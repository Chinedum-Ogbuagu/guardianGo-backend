package guardian

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
	Service Service 
}

func NewHandler(db *gorm.DB, s Service)  *Handler {
	return &Handler{DB: db, Service: s}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/parent")
	{
		group.POST("/register", h.CreateParent)
		group.GET("/:phone_number", h.GetParentByPhone)
	}
}

func (h *Handler) CreateParent(c *gin.Context) {
	var parent Guardian 
	if  err := c.ShouldBindJSON(&parent); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
	}
	if err := h.Service.Create(h.DB, &parent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return 
	}
	c.JSON(http.StatusOK, parent)
}

func (h *Handler) GetParentByPhone(c *gin.Context) {
	phone := c.Param("phone_number")
	parent, err := h.Service.GetByPhoneNumber(h.DB, phone)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "parent not found"})
		return 
	}
	c.JSON(http.StatusOK, parent)

}