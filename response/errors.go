package response

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Error represents a structured error response.
type Error struct {
	Object string    `json:"object"` // Always "error"
	Error  ErrorInfo `json:"error"`
}

// ErrorInfo contains error details.
type ErrorInfo struct {
	Type    string `json:"type"`            // error type category (see ErrorType* constants)
	Code    string `json:"code,omitempty"`  // machine-readable error code (see ErrorCode* constants)
	Message string `json:"message"`         // human-readable message
	Param   string `json:"param,omitempty"` // parameter that caused the error
}

// Error types - high-level categories for client-side error handling
const (
	ErrorTypeInvalidRequest = "invalid_request" // 400 - bad input, validation errors, missing params
	ErrorTypeAuthentication = "authentication"  // 401 - missing or invalid auth token
	ErrorTypeForbidden      = "forbidden"       // 403 - authenticated but not authorized
	ErrorTypeNotFound       = "not_found"       // 404 - resource doesn't exist
	ErrorTypeConflict       = "conflict"        // 409 - resource already exists, version conflict
	ErrorTypeRateLimit      = "rate_limit"      // 429 - too many requests
	ErrorTypeAPI            = "api"             // 5xx - server errors
)

// Error codes - specific machine-readable codes for programmatic handling
const (
	// Validation codes (used with ErrorTypeInvalidRequest)
	ErrorCodeInvalidParam  = "invalid_param"
	ErrorCodeMissingParam  = "missing_param"
	ErrorCodeInvalidFormat = "invalid_format"

	// Resource codes (used with ErrorTypeNotFound, ErrorTypeConflict)
	ErrorCodeResourceNotFound = "resource_not_found"
	ErrorCodeAlreadyExists    = "already_exists"

	// Auth codes (used with ErrorTypeAuthentication, ErrorTypeForbidden)
	ErrorCodeAuthRequired           = "auth_required"
	ErrorCodeInvalidToken           = "invalid_token"
	ErrorCodeTokenExpired           = "token_expired"
	ErrorCodeInsufficientPermission = "insufficient_permission"

	// Rate limit codes
	ErrorCodeRateLimitExceeded = "rate_limit_exceeded"

	// Server error codes (used with ErrorTypeAPI)
	ErrorCodeInternal           = "internal"
	ErrorCodeServiceUnavailable = "service_unavailable"
)

// sendError sends an error response with the given status and error info.
func sendError(c *gin.Context, status int, errType, code, message, param string) {
	c.JSON(status, Error{
		Object: "error",
		Error: ErrorInfo{
			Type:    errType,
			Code:    code,
			Message: message,
			Param:   param,
		},
	})
}

// BadRequest sends a 400 Bad Request error.
func BadRequest(c *gin.Context, message string) {
	sendError(c, http.StatusBadRequest, ErrorTypeInvalidRequest, "", message, "")
}

// BadRequestWithCode sends a 400 Bad Request error with a specific error code.
func BadRequestWithCode(c *gin.Context, code, message string) {
	sendError(c, http.StatusBadRequest, ErrorTypeInvalidRequest, code, message, "")
}

// BadRequestParam sends a 400 Bad Request error for a specific parameter.
func BadRequestParam(c *gin.Context, param, message string) {
	sendError(c, http.StatusBadRequest, ErrorTypeInvalidRequest, "", message, param)
}

// Unauthorized sends a 401 Unauthorized error.
func Unauthorized(c *gin.Context) {
	sendError(c, http.StatusUnauthorized, ErrorTypeAuthentication, "", "unauthorized", "")
}

// UnauthorizedWithMessage sends a 401 Unauthorized error with a custom message.
func UnauthorizedWithMessage(c *gin.Context, message string) {
	sendError(c, http.StatusUnauthorized, ErrorTypeAuthentication, "", message, "")
}

// Forbidden sends a 403 Forbidden error.
func Forbidden(c *gin.Context) {
	sendError(c, http.StatusForbidden, ErrorTypeForbidden, "", "forbidden", "")
}

// ForbiddenWithMessage sends a 403 Forbidden error with a custom message.
func ForbiddenWithMessage(c *gin.Context, message string) {
	sendError(c, http.StatusForbidden, ErrorTypeForbidden, "", message, "")
}

// NotFound sends a 404 Not Found error for an entity.
func NotFound(c *gin.Context, entity string) {
	sendError(c, http.StatusNotFound, ErrorTypeNotFound, "", fmt.Sprintf("%s not found", entity), "")
}

// NotFoundWithMessage sends a 404 Not Found error with a custom message.
func NotFoundWithMessage(c *gin.Context, message string) {
	sendError(c, http.StatusNotFound, ErrorTypeNotFound, "", message, "")
}

// Conflict sends a 409 Conflict error.
func Conflict(c *gin.Context, message string) {
	sendError(c, http.StatusConflict, ErrorTypeConflict, "", message, "")
}

// TooManyRequests sends a 429 Too Many Requests error.
func TooManyRequests(c *gin.Context, message string) {
	sendError(c, http.StatusTooManyRequests, ErrorTypeRateLimit, "", message, "")
}

// InternalError sends a 500 Internal Server Error.
func InternalError(c *gin.Context, message string) {
	sendError(c, http.StatusInternalServerError, ErrorTypeAPI, "", message, "")
}

// ServiceUnavailable sends a 503 Service Unavailable error.
func ServiceUnavailable(c *gin.Context, message string) {
	sendError(c, http.StatusServiceUnavailable, ErrorTypeAPI, "", message, "")
}

// NotImplemented sends a 501 Not Implemented error.
func NotImplemented(c *gin.Context, message string) {
	sendError(c, http.StatusNotImplemented, ErrorTypeAPI, "", message, "")
}

// UnprocessableEntity sends a 422 Unprocessable Entity error.
// Use for semantically invalid requests (e.g., valid JSON but business logic violation).
func UnprocessableEntity(c *gin.Context, message string) {
	sendError(c, http.StatusUnprocessableEntity, ErrorTypeInvalidRequest, "", message, "")
}

// UnsupportedMediaType sends a 415 Unsupported Media Type error.
func UnsupportedMediaType(c *gin.Context, message string) {
	sendError(c, http.StatusUnsupportedMediaType, ErrorTypeInvalidRequest, "", message, "")
}

// BadGateway sends a 502 Bad Gateway error.
// Use when an upstream service fails.
func BadGateway(c *gin.Context, message string) {
	sendError(c, http.StatusBadGateway, ErrorTypeAPI, "", message, "")
}
