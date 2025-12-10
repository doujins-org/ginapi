package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/doujins-org/ginapi/middleware"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestLanguageFromQueryParam(t *testing.T) {
	router := gin.New()
	router.Use(middleware.Language(middleware.LanguageConfig{
		Supported: []string{"en", "ja", "ko"},
		Default:   "en",
	}))
	router.GET("/test", func(c *gin.Context) {
		lang := middleware.GetLanguage(c)
		c.String(http.StatusOK, lang)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test?lang=ja", nil)
	router.ServeHTTP(w, req)

	if w.Body.String() != "ja" {
		t.Errorf("expected 'ja', got '%s'", w.Body.String())
	}
	if w.Header().Get("Content-Language") != "ja" {
		t.Errorf("expected Content-Language 'ja', got '%s'", w.Header().Get("Content-Language"))
	}
}

func TestLanguageFromPath(t *testing.T) {
	router := gin.New()
	router.Use(middleware.Language(middleware.LanguageConfig{
		Supported:       []string{"en", "ja", "ko"},
		Default:         "en",
		ExtractFromPath: true,
	}))
	router.GET("/ko/*path", func(c *gin.Context) {
		lang := middleware.GetLanguage(c)
		c.String(http.StatusOK, lang)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ko/galleries", nil)
	router.ServeHTTP(w, req)

	if w.Body.String() != "ko" {
		t.Errorf("expected 'ko', got '%s'", w.Body.String())
	}
}

func TestLanguageFromAcceptHeader(t *testing.T) {
	router := gin.New()
	router.Use(middleware.Language(middleware.LanguageConfig{
		Supported: []string{"en", "ja", "ko"},
		Default:   "en",
	}))
	router.GET("/test", func(c *gin.Context) {
		lang := middleware.GetLanguage(c)
		c.String(http.StatusOK, lang)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Accept-Language", "ko-KR,ko;q=0.9,en;q=0.8")
	router.ServeHTTP(w, req)

	if w.Body.String() != "ko" {
		t.Errorf("expected 'ko', got '%s'", w.Body.String())
	}
}

func TestLanguageDefault(t *testing.T) {
	router := gin.New()
	router.Use(middleware.Language(middleware.LanguageConfig{
		Supported: []string{"en", "ja"},
		Default:   "en",
	}))
	router.GET("/test", func(c *gin.Context) {
		lang := middleware.GetLanguage(c)
		c.String(http.StatusOK, lang)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Accept-Language", "fr,de")
	router.ServeHTTP(w, req)

	if w.Body.String() != "en" {
		t.Errorf("expected 'en', got '%s'", w.Body.String())
	}
}

func TestLanguageContext(t *testing.T) {
	ctx := context.Background()
	ctx = middleware.WithLanguage(ctx, "JA")

	lang := middleware.LanguageFromContext(ctx)
	if lang != "ja" {
		t.Errorf("expected 'ja', got '%s'", lang)
	}
}

func TestLanguageContextEmpty(t *testing.T) {
	lang := middleware.LanguageFromContext(nil)
	if lang != "" {
		t.Errorf("expected empty string, got '%s'", lang)
	}

	lang = middleware.LanguageFromContext(context.Background())
	if lang != "" {
		t.Errorf("expected empty string, got '%s'", lang)
	}
}
