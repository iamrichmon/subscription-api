package handler

import (
	"net/http"

	"errors"

	"github.com/gin-gonic/gin"
	"github.com/iamrichmon/subscription-api/internal/service"
	"github.com/iamrichmon/subscription-api/internal/utils"
)

type AdminHandler struct {
	adminService *service.AdminService
}

func NewAdminHandler(service *service.AdminService) *AdminHandler {
	return &AdminHandler{adminService: service}
}

type CreateAdminRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=12"`
}

type AdminLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AdminHandler) Register(c *gin.Context) {
	var req CreateAdminRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	admin, err := h.adminService.RegisterAdmin(req.Name, req.Email, req.Password)
	if err != nil {
		status := http.StatusInternalServerError
		switch {
		case errors.Is(err, utils.ErrEmailTaken):
			status = http.StatusConflict
		case errors.Is(err, utils.ErrInvalidCredentials):
			status = http.StatusUnauthorized
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, admin)
}

func (h *AdminHandler) Login(c *gin.Context) {
	var req AdminLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.adminService.Login(req.Email, req.Password)

	if err != nil {
		status := http.StatusInternalServerError
		switch {
		case errors.Is(err, utils.ErrInvalidCredentials):
			status = http.StatusUnauthorized
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *AdminHandler) GetAllUsers(c *gin.Context) {
	users, err := h.adminService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.ErrInternalServerError.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}
