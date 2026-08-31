package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		logger.Info("request completed", "method", c.Request.Method,
			"path", c.Request.URL.Path, "status", c.Writer.Status(), "duaration", time.Since(start).String(), "client_ip", c.ClientIP())
	}
}
