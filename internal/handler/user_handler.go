package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"gin-backend/internal/middleware"
	"gin-backend/internal/model"
	"gin-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/users", h.CreateUser)
	r.GET("/users/:id", h.GetUser)
	r.DELETE("/users/:id", h.DeleteUser)
	r.GET("/users", h.ListUsers)
	r.PATCH("/users/:id", h.UpdateUser)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var input model.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.RegisterUser(c.Request.Context(), input.Name, input.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	idParam := c.Param("id")
	if idParam == "" {
		c.Error(middleware.MissingParam("id"))
		return
	}
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Error(middleware.InvalidParam("id", err))
		return
	}
	user, err := h.service.GetUser(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timed out"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	if idParam == "" {
		c.Error(middleware.MissingParam("id"))
		return
	}
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Error(middleware.InvalidParam("id", err))
		return
	}
	var input model.UpdateUserRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(middleware.BadRequest("invalid request body", err))
		return
	}

	user, err := h.service.UpdateUser(c.Request.Context(), &input, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	if idParam == "" {
		c.Error(middleware.MissingParam("id"))
		return
	}
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Error(middleware.InvalidParam("id", err))
		return
	}

	err = h.service.DeleteUser(c.Request.Context(), id)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			c.Error(middleware.Timeout(err))
			return
		}
		c.Error(middleware.Internal(err)) // don't expose err.Error() to the client
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	limit, err := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	if err != nil {
		c.Error(middleware.InvalidParam("limit", err))
		return
	}
	if limit < 1 || limit > 100 {
		c.Error(middleware.InvalidParam("limit", fmt.Errorf("must be between 1 and 100, got %d", limit)))
		return
	}
	offset, err := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)
	if err != nil {
		c.Error(middleware.InvalidParam("offset", err))
		return
	}
	if offset < 0 {
		c.Error(middleware.InvalidParam("offset", fmt.Errorf("must be >= 0, got %d", offset)))
		return
	}
	users, err := h.service.ListUsers(c.Request.Context(), limit, offset)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			c.Error(middleware.Timeout(err))
			return
		}
		c.Error(middleware.Internal(err))
		return
	}
	c.JSON(http.StatusOK, users)
}
