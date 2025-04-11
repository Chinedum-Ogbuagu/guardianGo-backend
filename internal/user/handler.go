package user

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
	group := r.Group("/user")
	{
		group.POST("/register", h.CreateUser)
		group.GET("/:phone_number", h.GetUserByPhone)
		group.GET("", h.ListUsersByChurch)
	}
}

func (h *Handler) CreateUser(c *gin.Context) {
	var u User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Service.Create(h.DB, &u); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, u)
}

func (h *Handler) GetUserByPhone(c *gin.Context) {
	phone := c.Param("phone_number")
	u, err := h.Service.GetByPhone(h.DB, phone)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, u)
}

func (h *Handler) ListUsersByChurch(c *gin.Context) {
	churchID, err := strconv.Atoi(c.Query("church_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid church_id"})
		return
	}
	users, err := h.Service.ListByChurch(h.DB, uint(churchID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}
