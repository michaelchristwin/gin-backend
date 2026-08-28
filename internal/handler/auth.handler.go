package handler

import (
	"gin-backend/internal/middleware"
	"gin-backend/internal/model"
	"gin-backend/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authservice service.AuthService
	userService service.UserService
}

func NewAuthHandler(userService service.UserService, authService service.AuthService) *AuthHandler {
	return &AuthHandler{authservice: authService, userService: userService}
}

func (h *AuthHandler) RegisterRoutes(r *gin.RouterGroup, authMiddleware *middleware.AuthMiddleware) {
	r.POST("/auth/register", h.RegisterUser)
	r.POST("/auth/login", h.Login)
	protected := r.Group("/")
	protected.Use(authMiddleware.RequireAuth())
	{
		protected.POST("/auth/logout", h.Logout)
		protected.POST("/auth/change-password", h.UpdatePassword)
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input model.RegisterRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(middleware.BadRequest("Bad request", err))
		return
	}
	session, err := h.authservice.Login(c.Request.Context(), input)
	if err != nil {
		c.Error(middleware.Internal(err))
		return
	}

	maxAge := int(time.Until(session.ExpiresAt).Seconds())

	c.SetCookie("session_id", string(session.ID), maxAge, "/", "localhost", true, true)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful",
		"user": gin.H{"id": session.UserID, "email": input.Email}})
}

func (h *AuthHandler) RegisterUser(c *gin.Context) {
	var input model.RegisterRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(middleware.BadRequest("Bad request", err))
		return
	}
	user, err := h.authservice.RegisterUser(c.Request.Context(), input)
	if err != nil {
		c.Error(middleware.Internal(err))
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	session := c.MustGet(middleware.SessionContextKey).(*model.Session)

	err := h.authservice.Logout(c.Request.Context(), session.ID)
	if err != nil {
		c.Error(middleware.Internal(err))
		return
	}
	c.SetCookie("sessionId", "", -1, "/", "", true, true)
	c.Status(http.StatusNoContent)
}

func (h *AuthHandler) UpdatePassword(c *gin.Context) {
	var input model.ChangePasswordRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(middleware.BadRequest("Bad request", err))
		return
	}
	userID := c.MustGet("user_id").(int64)
	if err := h.authservice.ChangePassword(c.Request.Context(), userID, input.CurrentPassword, input.NewPassword); err != nil {
		c.Error(middleware.Internal(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}
