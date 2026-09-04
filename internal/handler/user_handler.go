package handler

import (
	"net/http"

	"errors"

	"github.com/gin-gonic/gin"
	"github.com/iamrichmon/subscription-api/internal/service"
	"github.com/iamrichmon/subscription-api/internal/utils"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{userService: service}
}

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=12"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var req CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.Register(req.Name, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, utils.ErrEmailTaken) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, utils.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.userService.Login(req.Email, req.Password)

	if err != nil {
		if errors.Is(err, utils.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": utils.ErrInvalidCredentials.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.ErrInternalServerError.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.ErrInternalServerError.Error()})
		return
	}
	user, err := h.userService.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.ErrInternalServerError.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}
