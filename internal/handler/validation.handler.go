package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func formatValidationError(err error) gin.H {
	if ve, ok := errors.AsType[validator.ValidationErrors](err); ok {
		fieldErrors := make([]gin.H, 0, len(ve))
		for _, fe := range ve {
			fieldErrors = append(fieldErrors, gin.H{
				"field":   fe.Field(),
				"message": validationMessage(fe),
			})
		}
		return gin.H{"error": "validation_failed", "fields": fieldErrors}
	}

	// malformed JSON body, wrong type, etc.
	return gin.H{"error": "invalid_request", "message": "Request body is malformed."}
}

func validationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fe.Field() + " is required."
	case "min":
		return fe.Field() + " must be at least " + fe.Param() + " characters."
	case "email":
		return fe.Field() + " must be a valid email address."
	case "nefield":
		return fe.Field() + " must be different from your current password."
	default:
		return fe.Field() + " is invalid."
	}
}
