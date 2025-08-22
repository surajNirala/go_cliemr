package utils

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`    // machine-readable error code
	Message string `json:"message"` // user-friendly message
	Details any    `json:"details,omitempty"`
}

// SuccessResponse for uniform success messages
type SuccessResponse struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func RespondError(c *gin.Context, status string, code int, message string, details ...any) {
	var detail any
	if len(details) > 0 {
		detail = details[0] // take the first element if provided
	}
	c.JSON(code, ErrorResponse{
		Status:  status,
		Code:    code,
		Message: message,
		Details: detail,
	})
}

func ValidationError(errs validator.ValidationErrors) string {
	var errMsg []string
	for _, err := range errs {
		switch err.Tag() {
		case "required":
			errMsg = append(errMsg, fmt.Sprintf("field %s is required", err.Field()))
		case "email":
			errMsg = append(errMsg, fmt.Sprintf("field %s must be a valid email", err.Field()))
		default:
			errMsg = append(errMsg, fmt.Sprintf("field %s is invalid", err.Field()))
		}
	}
	return strings.Join(errMsg, ", ")
}

func RespondSuccess(c *gin.Context, status string, code int, message string, data any) {
	c.JSON(code, SuccessResponse{
		Status:  status,
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func RequestAll(c *gin.Context) {
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		RespondError(c, "error", 500, "Failed to read request body", err.Error())
		c.Abort()
		return
	}
	// Restore body so Gin/other handlers can still read it
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	// Now you can log/inspect the raw body
	fmt.Println("Raw body:", bodyBytes)
	RespondSuccess(c, "success", 200, "Request received!!", string(bodyBytes))
}
