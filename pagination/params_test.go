package pagination_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/doujins-org/ginapi/pagination"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestBind(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		wantLimit  int
		wantOffset int
		wantSort   string
	}{
		{
			name:       "all params",
			query:      "?limit=50&offset=100&sort=created_at",
			wantLimit:  50,
			wantOffset: 100,
			wantSort:   "created_at",
		},
		{
			name:       "no params",
			query:      "",
			wantLimit:  0,
			wantOffset: 0,
			wantSort:   "",
		},
		{
			name:       "sort_by alias",
			query:      "?sort_by=name",
			wantLimit:  0,
			wantOffset: 0,
			wantSort:   "name",
		},
		{
			name:       "invalid values default to zero",
			query:      "?limit=abc&offset=xyz",
			wantLimit:  0,
			wantOffset: 0,
			wantSort:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("GET", "/test"+tt.query, nil)

			params := pagination.Bind(c)

			if params.Limit != tt.wantLimit {
				t.Errorf("expected limit %d, got %d", tt.wantLimit, params.Limit)
			}
			if params.Offset != tt.wantOffset {
				t.Errorf("expected offset %d, got %d", tt.wantOffset, params.Offset)
			}
			if params.Sort != tt.wantSort {
				t.Errorf("expected sort '%s', got '%s'", tt.wantSort, params.Sort)
			}
		})
	}
}

func TestBindDefault(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test", nil)

	params := pagination.BindDefault(c)

	if params.Limit != pagination.DefaultLimit {
		t.Errorf("expected default limit %d, got %d", pagination.DefaultLimit, params.Limit)
	}
	if params.Offset != 0 {
		t.Errorf("expected offset 0, got %d", params.Offset)
	}
}

func TestBindWithDefaults(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test?limit=200", nil)

	params := pagination.BindWithDefaults(c, 10, 50)

	// Should be capped to max
	if params.Limit != 50 {
		t.Errorf("expected limit capped to 50, got %d", params.Limit)
	}
}

func TestParams_Normalize(t *testing.T) {
	tests := []struct {
		name         string
		input        pagination.Params
		defaultLimit int
		maxLimit     int
		wantLimit    int
		wantOffset   int
	}{
		{
			name:         "apply defaults",
			input:        pagination.Params{Limit: 0, Offset: 0},
			defaultLimit: 20,
			maxLimit:     100,
			wantLimit:    20,
			wantOffset:   0,
		},
		{
			name:         "cap to max",
			input:        pagination.Params{Limit: 500, Offset: 0},
			defaultLimit: 20,
			maxLimit:     100,
			wantLimit:    100,
			wantOffset:   0,
		},
		{
			name:         "fix negative offset",
			input:        pagination.Params{Limit: 10, Offset: -5},
			defaultLimit: 20,
			maxLimit:     100,
			wantLimit:    10,
			wantOffset:   0,
		},
		{
			name:         "keep valid values",
			input:        pagination.Params{Limit: 50, Offset: 100},
			defaultLimit: 20,
			maxLimit:     100,
			wantLimit:    50,
			wantOffset:   100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.input
			p.Normalize(tt.defaultLimit, tt.maxLimit)

			if p.Limit != tt.wantLimit {
				t.Errorf("expected limit %d, got %d", tt.wantLimit, p.Limit)
			}
			if p.Offset != tt.wantOffset {
				t.Errorf("expected offset %d, got %d", tt.wantOffset, p.Offset)
			}
		})
	}
}
