// Package middleware provides Gin middleware for common API concerns.
package middleware

import (
	"context"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type contextKey string

const languageKey contextKey = "ginapi_language"

// LanguageConfig configures the language detection middleware.
type LanguageConfig struct {
	// Supported languages (e.g., []string{"en", "ja", "ko", "zh"})
	Supported []string
	// Default language if none detected (defaults to "en")
	Default string
	// QueryParam to check for language override (defaults to "lang")
	QueryParam string
	// ExtractFromPath enables checking URL path prefix (e.g., /ja/...)
	ExtractFromPath bool
}

// Language returns middleware that detects user language from:
// 1. Query parameter (?lang=ja)
// 2. URL path prefix (/ja/...) if ExtractFromPath is true
// 3. Accept-Language header with q-value parsing
//
// The detected language is stored in both gin context and request context,
// and the Content-Language header is set on the response.
func Language(cfg LanguageConfig) gin.HandlerFunc {
	// Build supported language map for fast lookup
	supportedMap := make(map[string]struct{}, len(cfg.Supported))
	for _, lang := range cfg.Supported {
		supportedMap[strings.ToLower(lang)] = struct{}{}
	}
	if len(supportedMap) == 0 {
		supportedMap["en"] = struct{}{}
	}

	// Normalize defaults
	defaultLang := strings.ToLower(strings.TrimSpace(cfg.Default))
	if defaultLang == "" {
		defaultLang = "en"
	}

	queryParam := cfg.QueryParam
	if queryParam == "" {
		queryParam = "lang"
	}

	return func(c *gin.Context) {
		lang := resolveLanguage(c, supportedMap, defaultLang, queryParam, cfg.ExtractFromPath)

		// Store in gin context
		c.Set("language", lang)

		// Store in request context
		ctx := WithLanguage(c.Request.Context(), lang)
		c.Request = c.Request.WithContext(ctx)

		// Set response header
		c.Header("Content-Language", lang)

		c.Next()
	}
}

// resolveLanguage determines the best language from available sources.
func resolveLanguage(c *gin.Context, supported map[string]struct{}, fallback, queryParam string, extractFromPath bool) string {
	if c == nil || c.Request == nil {
		return fallback
	}

	// 1. Check query parameter
	if lang := strings.ToLower(strings.TrimSpace(c.Query(queryParam))); lang != "" {
		if _, ok := supported[lang]; ok {
			return lang
		}
	}

	// 2. Check URL path prefix
	if extractFromPath {
		if lang := extractLanguageFromPath(c.Request.URL.Path); lang != "" {
			if _, ok := supported[lang]; ok {
				return lang
			}
		}
	}

	// 3. Check Accept-Language header
	if header := c.GetHeader("Accept-Language"); header != "" {
		if lang := parseAcceptLanguage(header, supported); lang != "" {
			return lang
		}
	}

	return fallback
}

// extractLanguageFromPath extracts a 2-3 character language code from URL path prefix.
// e.g., "/ja/galleries" -> "ja"
func extractLanguageFromPath(path string) string {
	trimmed := strings.TrimPrefix(path, "/")
	parts := strings.SplitN(trimmed, "/", 2)
	if len(parts) == 0 {
		return ""
	}
	first := parts[0]
	if len(first) == 2 || len(first) == 3 {
		return strings.ToLower(first)
	}
	return ""
}

// parseAcceptLanguage parses the Accept-Language header and returns the best
// supported language based on q-values.
func parseAcceptLanguage(header string, supported map[string]struct{}) string {
	type candidate struct {
		lang string
		q    float64
	}

	var candidates []candidate
	parts := strings.Split(header, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		lang := part
		q := 1.0

		// Parse q-value if present (e.g., "en-US;q=0.9")
		if idx := strings.Index(part, ";"); idx >= 0 {
			lang = part[:idx]
			param := strings.TrimSpace(part[idx+1:])
			if strings.HasPrefix(param, "q=") {
				if v, err := strconv.ParseFloat(strings.TrimPrefix(param, "q="), 64); err == nil {
					q = v
				}
			}
		}

		// Extract base language (e.g., "en-US" -> "en")
		if hyphen := strings.Index(lang, "-"); hyphen >= 0 {
			lang = lang[:hyphen]
		}
		lang = strings.ToLower(strings.TrimSpace(lang))

		if lang == "" {
			continue
		}

		candidates = append(candidates, candidate{lang: lang, q: q})
	}

	// Find highest q-value language that's supported
	var bestLang string
	var bestQ float64 = -1

	for _, c := range candidates {
		if _, ok := supported[c.lang]; ok && c.q > bestQ {
			bestLang = c.lang
			bestQ = c.q
		}
	}

	return bestLang
}

// GetLanguage retrieves the detected language from the gin context.
// Falls back to request context, then returns "en" as ultimate fallback.
func GetLanguage(c *gin.Context) string {
	if c == nil {
		return "en"
	}

	// First try gin context (faster)
	if lang, exists := c.Get("language"); exists {
		if s, ok := lang.(string); ok && s != "" {
			return s
		}
	}

	// Fallback to request context
	if lang := LanguageFromContext(c.Request.Context()); lang != "" {
		return lang
	}

	return "en"
}

// WithLanguage stores a language code in the context.
func WithLanguage(ctx context.Context, lang string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if lang == "" {
		return ctx
	}
	return context.WithValue(ctx, languageKey, strings.ToLower(lang))
}

// LanguageFromContext retrieves the language from a context.
// Returns empty string if not set.
func LanguageFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v, ok := ctx.Value(languageKey).(string); ok {
		return v
	}
	return ""
}
