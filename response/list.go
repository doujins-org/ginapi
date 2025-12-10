package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// List is a Stripe-style list response with offset/limit pagination.
type List[T any] struct {
	Object  string `json:"object"`   // Always "list"
	Data    []T    `json:"data"`     // The items
	Total   int64  `json:"total"`    // Total count across all pages
	Limit   int    `json:"limit"`    // Max items requested
	Offset  int    `json:"offset"`   // Items skipped
	HasMore bool   `json:"has_more"` // More items available
}

// NewList creates a List response with has_more calculated automatically.
func NewList[T any](data []T, total int64, limit, offset int) List[T] {
	// Ensure data is never nil (empty slice for clean JSON)
	if data == nil {
		data = []T{}
	}

	hasMore := int64(offset+len(data)) < total

	return List[T]{
		Object:  "list",
		Data:    data,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: hasMore,
	}
}

// ListResponse sends a Stripe-style list response.
func ListResponse[T any](c *gin.Context, data []T, total int64, limit, offset int) {
	c.JSON(http.StatusOK, NewList(data, total, limit, offset))
}
