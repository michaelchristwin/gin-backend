package middleware

import (
	"gin-backend/internal/service"

	"github.com/gin-gonic/gin"
)

const (
	SessionContextKey = "session"
	UserContextKey    = "user"
)

type AuthMiddleware struct {
	sessionService service.AuthService
	userService    service.UserService
}

func NewAuthMiddleware(
	sessionService service.AuthService,
	userService service.UserService,
) *AuthMiddleware {
	return &AuthMiddleware{
		sessionService: sessionService,
		userService:    userService,
	}

}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("session_id")
		if err != nil || sessionID == "" {
			c.Error(Unauthorized("authentication required", err))
			c.Abort()
			return
		}

		session, err := m.sessionService.GetByID(
			c.Request.Context(),
			sessionID,
		)
		if err != nil {
			c.Error(Unauthorized("invalid session", err))
			c.Abort()
			return
		}

		user, err := m.userService.GetUser(
			c.Request.Context(),
			session.UserID,
		)
		if err != nil {
			c.Error(Internal(err))
			c.Abort()
			return
		}

		c.Set(SessionContextKey, session)
		c.Set(UserContextKey, user)

		c.Next()
	}
}
