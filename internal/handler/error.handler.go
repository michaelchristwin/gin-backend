package handler

import (
	"errors"
	"gin-backend/internal/domain"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type apiError struct {
	status  int
	code    string
	message string
}

var errorMap = map[error]apiError{
	domain.ErrEmailTaken: {
		status:  http.StatusConflict,
		code:    "email_taken",
		message: "An account with this email already exists.",
	},
	domain.ErrPasswordPwned: {
		status:  http.StatusBadRequest,
		code:    "password_compromised",
		message: "This password has appeared in a known data breach. Please choose a different password.",
	},
	domain.ErrUserNotFound: {
		status:  http.StatusNotFound,
		code:    "user_not_found",
		message: "User not found.",
	},
	domain.ErrInvalidCredentials: {
		status:  http.StatusUnauthorized,
		code:    "invalid_credentials",
		message: "Invalid email or password.",
	},
}

func respondWithError(c *gin.Context, err error) {
	for domainErr, apiErr := range errorMap {
		if errors.Is(err, domainErr) {
			c.JSON(apiErr.status, gin.H{
				"error":   apiErr.code,
				"message": apiErr.message,
			})
			return
		}
	}

	// fallback: unknown error, don't leak internals
	log.Printf("unhandled error: %v", err)
	c.JSON(http.StatusInternalServerError, gin.H{
		"error":   "internal_error",
		"message": "Something went wrong. Please try again later.",
	})
}
