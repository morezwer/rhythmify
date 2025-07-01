package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
}

// SuccessResponse sends a successful response
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponseWithCode sends an error response with custom error code
func ErrorResponseWithCode(c *gin.Context, statusCode int, error string, code string) {
	c.JSON(statusCode, ErrorResponse{
		Success: false,
		Error:   error,
		Code:    code,
	})
}

// BadRequest sends a 400 Bad Request response
func BadRequest(c *gin.Context, error string) {
	ErrorResponseWithCode(c, http.StatusBadRequest, error, "BAD_REQUEST")
}

// Unauthorized sends a 401 Unauthorized response
func Unauthorized(c *gin.Context, error string) {
	ErrorResponseWithCode(c, http.StatusUnauthorized, error, "UNAUTHORIZED")
}

// Forbidden sends a 403 Forbidden response
func Forbidden(c *gin.Context, error string) {
	ErrorResponseWithCode(c, http.StatusForbidden, error, "FORBIDDEN")
}

// NotFound sends a 404 Not Found response
func NotFound(c *gin.Context, error string) {
	ErrorResponseWithCode(c, http.StatusNotFound, error, "NOT_FOUND")
}

// Conflict sends a 409 Conflict response
func Conflict(c *gin.Context, error string) {
	ErrorResponseWithCode(c, http.StatusConflict, error, "CONFLICT")
}

// InternalServerError sends a 500 Internal Server Error response
func InternalServerError(c *gin.Context, error string) {
	ErrorResponseWithCode(c, http.StatusInternalServerError, error, "INTERNAL_SERVER_ERROR")
}

// Created sends a 201 Created response
func Created(c *gin.Context, message string, data interface{}) {
	SuccessResponse(c, http.StatusCreated, message, data)
}

// OK sends a 200 OK response
func OK(c *gin.Context, message string, data interface{}) {
	SuccessResponse(c, http.StatusOK, message, data)
}