# ginapi

Shared Gin API utilities for consistent REST APIs across projects. Uses Stripe-style response formats with offset/limit pagination.

## Why?

Eliminates copy-paste boilerplate for:
- Response envelopes (consistent JSON structure)
- Error handling (typed errors with proper HTTP status codes)
- Pagination (offset/limit binding and response formatting)
- ID formatting (Stripe-style `gal_7B` prefixes)

## Installation

```bash
go get github.com/doujins-org/ginapi
```

## Response Format

All responses follow Stripe conventions:

```json
// Single object
{"object": "artist", "id": "art_7B", "name": "Takehiko Inoue"}

// List with pagination
{
  "object": "list",
  "data": [{"object": "artist", ...}, ...],
  "total": 1523,
  "limit": 20,
  "offset": 0,
  "has_more": true
}

// Error
{
  "object": "error",
  "error": {
    "type": "not_found_error",
    "message": "artist not found"
  }
}
```

## Quick Start

```go
import (
    "github.com/doujins-org/ginapi/response"
    "github.com/doujins-org/ginapi/pagination"
)

func ListArtists(c *gin.Context) {
    params := pagination.BindDefault(c)  // ?limit=20&offset=0

    artists, total, err := artistRepo.List(ctx, params.Limit, params.Offset)
    if err != nil {
        response.InternalError(c, "failed to list artists")
        return
    }

    response.ListResponse(c, artists, total, params.Limit, params.Offset)
}

func GetArtist(c *gin.Context) {
    artist, err := artistRepo.GetByID(ctx, c.Param("id"))
    if err != nil {
        response.NotFound(c, "artist")
        return
    }

    response.Object(c, artist)
}
```

## Packages

| Package | Purpose |
|---------|---------|
| `response` | Stripe-style JSON responses and errors |
| `pagination` | Offset/limit param binding |
| `ids` | Stripe-style ID formatting (`gal_7B`) |
| `middleware` | Language detection from query/path/headers |

## ID Formatting

```go
import "github.com/doujins-org/ginapi/ids"

formatted := ids.Format("gal", 123)     // "gal_1Z"
parsed, _ := ids.Parse("gal", "gal_1Z") // 123
```

## Pagination

```go
import "github.com/doujins-org/ginapi/pagination"

// Bind with defaults (limit=20, max=100)
params := pagination.BindDefault(c)  // ?limit=20&offset=0

// Custom defaults
params := pagination.BindWithDefaults(c, 50, 200)  // default=50, max=200

// Use in queries
artists, total, err := repo.List(ctx, params.Limit, params.Offset)

// Return paginated response
response.ListResponse(c, artists, total, params.Limit, params.Offset)
```

## Language Detection

```go
import "github.com/doujins-org/ginapi/middleware"

// Add middleware to router
router.Use(middleware.Language(middleware.LanguageConfig{
    Supported:       []string{"en", "ja", "ko", "zh"},
    Default:         "en",
    QueryParam:      "lang",        // ?lang=ja
    ExtractFromPath: true,          // /ja/galleries
}))

// In handlers:
lang := middleware.GetLanguage(c)  // "ja"

// Or from context (for use in services):
lang := middleware.LanguageFromContext(ctx)
```
