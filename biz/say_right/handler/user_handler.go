package handler

import (
	"net/http"
	"strings"

	"api/biz/say_right/dal/model"
	"api/biz/say_right/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{
		svc: svc,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Auto-normalize email
	user.EmailNorm = strings.ToLower(strings.TrimSpace(user.Email))

	// Note: In a real application, you should handle password hashing, validation, etc.
	// This is just a generated example.

	if err := h.svc.CreateUser(c.Request.Context(), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}

	// Normalize email for lookup
	emailNorm := strings.ToLower(strings.TrimSpace(email))

	user, err := h.svc.GetUserByEmail(c.Request.Context(), emailNorm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
