package handler

import (
	"gin-backend/internal/middleware"
	"gin-backend/internal/model"
	"gin-backend/internal/service"
	"net/http"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authservice service.AuthService
	userService service.UserService
}

func NewAuthHandler(userService service.UserService, authService service.AuthService) *AuthHandler {
	return &AuthHandler{authservice: authService, userService: userService}
}

func (h *AuthHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/auth/register", h.RegisterUser)
	r.POST("/auth/login", h.Login)
	r.POST("/auth/logout", h.Logout)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input model.UserRequest
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

	c.SetCookie("sessionId", string(session.ID), maxAge, "/", "", true, true)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful",
		"user": gin.H{"id": session.UserID, "email": input.Email}})
}

func (h *AuthHandler) RegisterUser(c *gin.Context) {
	var input model.UserRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(middleware.BadRequest("Bad request", err))
		return
	}
	passwordHash, err := argon2id.CreateHash(input.Password, argon2id.DefaultParams)
	if err != nil {
		c.Error(middleware.Internal(err))
		return
	}

	user, err := h.userService.RegisterUser(c.Request.Context(), input.Email, passwordHash)
	if err != nil {
		c.Error(middleware.Internal(err))
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var input model.DeleteSessionReq
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(middleware.BadRequest("Bad request", err))
		return
	}

	err := h.authservice.Logout(c.Request.Context(), input.SessionID)
	if err != nil {
		c.Error(middleware.Internal(err))
		return
	}
	c.SetCookie("sessionId", "", -1, "/", "", true, true)
	c.Status(http.StatusNoContent)
}
