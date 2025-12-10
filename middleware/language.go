// Package middleware provides Gin middleware for common API concerns.
package middleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type contextKey string

const languageKey contextKey = "ginapi_language"

const (
	// LanguageCookieName is the default cookie name for language preference
	LanguageCookieName = "lang"
	// LanguageCookieMaxAge is 1 year in seconds
	LanguageCookieMaxAge = 31536000
)

// LanguageConfig configures the language detection middleware.
type LanguageConfig struct {
	// Supported languages (e.g., []string{"en", "ja", "ko", "zh"})
	Supported []string
	// Default language if none detected (defaults to "en")
	Default string
	// QueryParam to check for language override (defaults to "lang")
	QueryParam string
	// CookieName to check for language preference (defaults to "lang")
	CookieName string
}

// Language returns middleware that detects user language from:
// 1. Query parameter (?lang=ja) - for API routes
// 2. URL path prefix (/ja/...) - for frontend routes
// 3. Cookie (user's saved preference)
// 4. Accept-Language header with q-value parsing
// 5. Default language
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

	cookieName := cfg.CookieName
	if cookieName == "" {
		cookieName = "lang"
	}

	return func(c *gin.Context) {
		lang := resolveLanguage(c, supportedMap, defaultLang, queryParam, cookieName)

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
func resolveLanguage(c *gin.Context, supported map[string]struct{}, fallback, queryParam, cookieName string) string {
	if c == nil || c.Request == nil {
		return fallback
	}

	// 1. Check query parameter (for API routes like /api/v1/videos?lang=ja)
	if lang := strings.ToLower(strings.TrimSpace(c.Query(queryParam))); lang != "" {
		if _, ok := supported[lang]; ok {
			return lang
		}
	}

	// 2. Check URL path prefix (for frontend routes like /ja/videos)
	if lang := extractLanguageFromPath(c.Request.URL.Path); lang != "" {
		if _, ok := supported[lang]; ok {
			return lang
		}
	}

	// 3. Check cookie (user's saved preference)
	if cookieName != "" {
		if lang, err := c.Cookie(cookieName); err == nil && lang != "" {
			lang = strings.ToLower(strings.TrimSpace(lang))
			if _, ok := supported[lang]; ok {
				return lang
			}
		}
	}

	// 4. Check Accept-Language header
	if header := c.GetHeader("Accept-Language"); header != "" {
		if lang := ParseAcceptLanguage(header, supported); lang != "" {
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

// ParseAcceptLanguage parses the Accept-Language header and returns the best
// supported language based on q-values. Exported for use by redirect middleware.
func ParseAcceptLanguage(header string, supported map[string]struct{}) string {
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

// BuildSupportedMap creates a map of supported languages for fast lookup.
// Useful for redirect middleware that needs to check language validity.
func BuildSupportedMap(languages []string) map[string]struct{} {
	m := make(map[string]struct{}, len(languages))
	for _, lang := range languages {
		m[strings.ToLower(lang)] = struct{}{}
	}
	return m
}

// ExtractLanguageFromPath extracts a 2-3 character language code from URL path prefix.
// e.g., "/ja/galleries" -> "ja", "/galleries" -> ""
// Exported for use by redirect middleware.
func ExtractLanguageFromPath(path string) string {
	return extractLanguageFromPath(path)
}

// LanguageRedirectConfig configures the language redirect middleware.
type LanguageRedirectConfig struct {
	// Supported languages (e.g., []string{"en", "ja", "ko", "zh"})
	Supported []string
	// Default language if none detected (defaults to "en")
	Default string
	// SkipPaths are path prefixes that should not be redirected (e.g., "/api/", "/auth/")
	// By default, /api/, /auth/, /oauth/, /.well-known/ are skipped.
	SkipPaths []string
	// SkipExtensions are file extensions that should not be redirected
	// By default, common static file extensions are skipped.
	SkipExtensions []string
}

// LanguageRedirect returns middleware that handles language URL routing for frontend routes.
//
// Behavior:
// 1. If URL has a valid language prefix (e.g., /en/videos, /ja/videos):
//   - Set the language cookie for future visits
//   - Continue to serve the page
//
// 2. If URL has NO language prefix (e.g., /videos):
//   - Detect language from: cookie → Accept-Language → default
//   - 302 redirect to the language-prefixed URL (e.g., /en/videos)
//
// This ensures all frontend URLs have a language prefix.
// API routes (/api/*) and static files are automatically skipped.
func LanguageRedirect(cfg LanguageRedirectConfig) gin.HandlerFunc {
	if len(cfg.Supported) == 0 {
		// No languages configured - pass through without redirects
		return func(c *gin.Context) {
			c.Next()
		}
	}

	// Build supported language map for fast lookup
	supportedMap := BuildSupportedMap(cfg.Supported)
	defaultLang := strings.ToLower(strings.TrimSpace(cfg.Default))
	if defaultLang == "" {
		defaultLang = "en"
	}

	// Default skip paths
	skipPaths := cfg.SkipPaths
	if len(skipPaths) == 0 {
		skipPaths = []string{"/api/", "/auth/", "/oauth/", "/.well-known/"}
	}

	// Default skip extensions
	skipExtensions := cfg.SkipExtensions
	if len(skipExtensions) == 0 {
		skipExtensions = []string{
			".js", ".css", ".png", ".jpg", ".jpeg", ".gif", ".svg", ".ico",
			".woff", ".woff2", ".ttf", ".eot", ".map", ".webp", ".mp4", ".webm",
		}
	}

	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// Skip paths that shouldn't be redirected
		if shouldSkipRedirect(path, skipPaths, skipExtensions) {
			c.Next()
			return
		}

		// Check if URL already has a language prefix
		langFromPath := extractLanguageFromPath(path)
		if langFromPath != "" {
			// Validate it's a supported language
			if _, ok := supportedMap[langFromPath]; ok {
				// Valid language prefix - set cookie and continue
				setLanguageCookie(c, langFromPath)
				c.Next()
				return
			}
			// Invalid language prefix - fall through to redirect
		}

		// No valid language prefix - determine preferred language and redirect
		preferredLang := detectPreferredLanguage(c, supportedMap, defaultLang)

		// Build redirect URL with language prefix
		redirectURL := "/" + preferredLang + path
		if c.Request.URL.RawQuery != "" {
			redirectURL += "?" + c.Request.URL.RawQuery
		}

		c.Redirect(http.StatusFound, redirectURL)
		c.Abort()
	}
}

// shouldSkipRedirect returns true for paths that should not be redirected.
func shouldSkipRedirect(path string, skipPaths, skipExtensions []string) bool {
	// Check skip paths
	for _, prefix := range skipPaths {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}

	// Skip common health check paths
	if path == "/health" || path == "/healthz" || path == "/readyz" {
		return true
	}

	// Skip runtime config
	if path == "/runtime-config.js" {
		return true
	}

	// Check file extensions
	lowPath := strings.ToLower(path)
	for _, ext := range skipExtensions {
		if strings.HasSuffix(lowPath, ext) {
			return true
		}
	}

	return false
}

// detectPreferredLanguage determines user's preferred language for redirect.
// Priority: cookie → Accept-Language → default
func detectPreferredLanguage(c *gin.Context, supportedMap map[string]struct{}, defaultLang string) string {
	// 1. Check cookie (user's saved preference)
	if lang, err := c.Cookie(LanguageCookieName); err == nil && lang != "" {
		lang = strings.ToLower(strings.TrimSpace(lang))
		if _, ok := supportedMap[lang]; ok {
			return lang
		}
	}

	// 2. Check Accept-Language header
	if header := c.GetHeader("Accept-Language"); header != "" {
		if lang := ParseAcceptLanguage(header, supportedMap); lang != "" {
			return lang
		}
	}

	// 3. Fall back to default
	return defaultLang
}

// setLanguageCookie sets the language preference cookie.
func setLanguageCookie(c *gin.Context, lang string) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(LanguageCookieName, lang, LanguageCookieMaxAge, "/", "", false, false)
}
