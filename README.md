# ginapi

Shared Gin utilities for building consistent REST APIs across multiple Go projects.

## Why?

When you have multiple services (doujins, hentai0, etc.) that need consistent API behavior, this library provides:

- **Response formatting** - Stripe-style JSON envelopes so all APIs return the same structure
- **Pagination** - Offset/limit binding with sensible defaults
- **Language detection** - Middleware that detects user language from query params, URL path, cookies, or Accept-Language headers
- **Language redirects** - Helpers for routes without a language prefix (e.g., `/galleries`) that need to 302 redirect to a concrete language (e.g., `/en/galleries`) based on user preference

Instead of copy-pasting this boilerplate into each project, import it once.

## Install

```bash
go get github.com/doujins-org/ginapi
```

## Response Format

```json
{"object": "artist", "id": "123", "name": "..."}

{"object": "list", "data": [...], "total": 100, "limit": 20, "offset": 0, "has_more": true}

{"object": "error", "error": {"type": "not_found_error", "message": "..."}}
```

## Pagination

```go
params := pagination.BindDefault(c)  // ?limit=20&offset=0
response.ListResponse(c, items, total, params.Limit, params.Offset)
```

## Language Middleware

Detects language from: query param → URL path → cookie → Accept-Language → default.

```go
router.Use(middleware.Language(middleware.LanguageConfig{
    Supported: []string{"en", "ja", "ko"},
    Default:   "en",
}))

lang := middleware.GetLanguage(c)
```

## Language Redirect (NoRoute)

Redirects `/galleries` → `/en/galleries` based on user preference.

```go
r.NoRoute(func(c *gin.Context) {
    if middleware.HandleLanguageRedirect(c, middleware.LanguageRedirectConfig{
        Supported: []string{"en", "ja", "ko"},
        Default:   "en",
    }) {
        return
    }
    serveIndexHTML(c)
})
```

## Reference

| Function | Description |
|----------|-------------|
| `GetLanguage(c)` | Get detected language from gin context |
| `LanguageFromContext(ctx)` | Get language from request context |
| `SetLanguageCookie(c, lang)` | Set 1-year language cookie |
| `ExtractLanguageFromPath(path)` | Extract lang prefix from URL |
| `ParseAcceptLanguage(header, supported)` | Parse Accept-Language header |
