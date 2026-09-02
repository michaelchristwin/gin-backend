package handler

import (
	"errors"
	"fmt"
	"net/http"
)

type CustomError struct {
	BaseErr     error                  // Underlying Base error
	StatusCode  int                    // HTTP status code
	Message     string                 // Detailed error message
	UserMessage string                 // User-friendly error message
	ErrType     string                 // Error type
	ErrCode     string                 // Unique error code
	Retryable   bool                   // Retryable flag
	Metadata    map[string]interface{} // Additional metadata
}

type APIError interface {
	FromCustomError(customErr *CustomError) error
}

type APIErrorCreator interface {
	New() APIError
}

type ResponseWriter interface {
	WriteResponse(statusCode int, body interface{}) error
}

func New(statusCode int, message, userMessage, errType, errCode string, retryable bool) *CustomError {
	return &CustomError{
		BaseErr:     fmt.Errorf("error: %s", message),
		StatusCode:  statusCode,
		Message:     message,
		UserMessage: userMessage,
		ErrType:     errType,
		ErrCode:     errCode,
		Retryable:   retryable,
		Metadata:    make(map[string]interface{}),
	}
}

func NewNotFoundError(resource string) *CustomError {
	return New(
		http.StatusNotFound,
		fmt.Sprintf("%s not found", resource),
		"The requested resource could not be found.",
		"NOT_FOUND",
		"ERR_NOT_FOUND",
		false,
	)
}

func NewInternalServerError() *CustomError {
	return New(http.StatusInternalServerError,
		"something went wrong",
		"Something went wrong, please try again.",
		"INTERNAL_SERVER_ERROR",
		"INTERNAL",
		false)
}

func NewBadRequestError() *CustomError {

}

func NewFromError(err error, statusCode int, userMessage, errType, errCode string, retryable bool) *CustomError {
	return &CustomError{
		BaseErr:     err,
		StatusCode:  statusCode,
		Message:     err.Error(),
		UserMessage: userMessage,
		ErrType:     errType,
		ErrCode:     errCode,
		Retryable:   retryable,
		Metadata:    make(map[string]interface{}),
	}
}

func NewErrHandler(creator APIErrorCreator, writerFactory func() ResponseWriter) func(error) {
	return func(werr error) {
		var customErr *CustomError
		writer := writerFactory()

		// Check if the error is a CustomError
		if errors.As(werr, &customErr) {
			apiErr := creator.New()
			if err := apiErr.FromCustomError(customErr); err != nil {
				_ = writer.WriteResponse(http.StatusInternalServerError, NewParseErrorError(werr, err))
				return
			}

			// Write the transformed error response
			_ = writer.WriteResponse(customErr.StatusCode, apiErr)
			return
		}

		// Handle generic errors
		_ = writer.WriteResponse(http.StatusInternalServerError, NewInternalServerError())
	}
}
