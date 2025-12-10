// Package pagination provides utilities for handling pagination in HTTP requests.
package pagination

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// Default pagination values
const (
	DefaultLimit = 20
	MaxLimit     = 100
)

// Params holds pagination parameters from a request.
type Params struct {
	Limit  int
	Offset int
	Sort   string
}

// Normalize applies defaults and caps.
func (p *Params) Normalize(defaultLimit, maxLimit int) {
	if p.Limit <= 0 {
		p.Limit = defaultLimit
	}
	if p.Limit > maxLimit {
		p.Limit = maxLimit
	}
	if p.Offset < 0 {
		p.Offset = 0
	}
}

// Bind extracts pagination parameters from a Gin context.
// Supports: limit, offset, sort (or sort_by)
func Bind(c *gin.Context) Params {
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))

	sort := c.Query("sort")
	if sort == "" {
		sort = c.Query("sort_by")
	}

	return Params{
		Limit:  limit,
		Offset: offset,
		Sort:   sort,
	}
}

// BindWithDefaults extracts and normalizes pagination parameters.
func BindWithDefaults(c *gin.Context, defaultLimit, maxLimit int) Params {
	p := Bind(c)
	p.Normalize(defaultLimit, maxLimit)
	return p
}

// BindDefault extracts pagination with standard defaults (limit 20, max 100).
func BindDefault(c *gin.Context) Params {
	return BindWithDefaults(c, DefaultLimit, MaxLimit)
}
