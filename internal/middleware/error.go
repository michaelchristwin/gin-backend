package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ---- 1. Custom error type ----

// AppError lets handlers attach an HTTP status and a safe client-facing
// message, while still preserving the original error for logging.
type AppError struct {
	Code    int    // HTTP status code
	Message string // Safe message to show the client
	Err     error  // Underlying error (for logs, not shown to client)
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error { return e.Err }

// Helper constructors
func NewAppError(code int, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

func NotFound(message string) *AppError {
	return NewAppError(http.StatusNotFound, message, nil)
}

func BadRequest(message string, err error) *AppError {
	return NewAppError(http.StatusBadRequest, message, err)
}

func Internal(err error) *AppError {
	return NewAppError(http.StatusInternalServerError, "internal server error", err)
}
func InvalidParam(param string, err error) *AppError {
	return NewAppError(
		http.StatusBadRequest,
		fmt.Sprintf("invalid value for parameter %q", param),
		err,
	)
}

// MissingParam covers the case where a required param is empty/absent
// (e.g. c.Param returns "" or a required query param wasn't set).
func MissingParam(param string) *AppError {
	return NewAppError(
		http.StatusBadRequest,
		fmt.Sprintf("missing required parameter %q", param),
		nil,
	)
}

func Timeout(err error) *AppError {
	return NewAppError(
		http.StatusGatewayTimeout,
		"request timed out",
		err,
	)
}

// ---- 2. The global error-handling middleware ----

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Recover from panics so they become 500s, not crashes
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}
				log.Printf("panic recovered: %v", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "internal server error",
				})
			}
		}()

		c.Next() // run the rest of the chain (handlers may call c.Error(...))

		// If nothing errored, nothing to do
		if len(c.Errors) == 0 {
			return
		}

		// Use the last error added
		err := c.Errors.Last().Err

		if appErr, ok := errors.AsType[*AppError](err); ok {
			// Log full detail server-side
			log.Printf("handled error: %v", appErr)
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
			return
		}

		// Unknown/unexpected error type -> 500, don't leak details
		log.Printf("unhandled error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
