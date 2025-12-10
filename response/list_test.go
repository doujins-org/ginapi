package response_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/doujins-org/ginapi/response"
)

func TestNewList(t *testing.T) {
	tests := []struct {
		name       string
		data       []string
		total      int64
		limit      int
		offset     int
		wantMore   bool
		wantLen    int
	}{
		{
			name:     "first page with more",
			data:     []string{"a", "b", "c"},
			total:    10,
			limit:    3,
			offset:   0,
			wantMore: true,
			wantLen:  3,
		},
		{
			name:     "last page",
			data:     []string{"j"},
			total:    10,
			limit:    3,
			offset:   9,
			wantMore: false,
			wantLen:  1,
		},
		{
			name:     "empty list",
			data:     nil,
			total:    0,
			limit:    20,
			offset:   0,
			wantMore: false,
			wantLen:  0,
		},
		{
			name:     "exact page boundary",
			data:     []string{"a", "b", "c"},
			total:    3,
			limit:    3,
			offset:   0,
			wantMore: false,
			wantLen:  3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := response.NewList(tt.data, tt.total, tt.limit, tt.offset)

			if list.Object != "list" {
				t.Errorf("expected object 'list', got '%s'", list.Object)
			}
			if len(list.Data) != tt.wantLen {
				t.Errorf("expected %d items, got %d", tt.wantLen, len(list.Data))
			}
			if list.Total != tt.total {
				t.Errorf("expected total %d, got %d", tt.total, list.Total)
			}
			if list.Limit != tt.limit {
				t.Errorf("expected limit %d, got %d", tt.limit, list.Limit)
			}
			if list.Offset != tt.offset {
				t.Errorf("expected offset %d, got %d", tt.offset, list.Offset)
			}
			if list.HasMore != tt.wantMore {
				t.Errorf("expected has_more %v, got %v", tt.wantMore, list.HasMore)
			}
		})
	}
}

func TestNewList_NilDataBecomesEmptySlice(t *testing.T) {
	list := response.NewList[string](nil, 0, 20, 0)

	if list.Data == nil {
		t.Error("expected data to be empty slice, not nil")
	}

	// Verify JSON serialization produces [] not null
	b, err := json.Marshal(list)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var parsed map[string]any
	if err := json.Unmarshal(b, &parsed); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	data, ok := parsed["data"].([]any)
	if !ok {
		t.Errorf("expected data to be array, got %T", parsed["data"])
	}
	if len(data) != 0 {
		t.Errorf("expected empty array, got %v", data)
	}
}

func TestListResponse(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	items := []string{"item1", "item2"}
	response.ListResponse(c, items, 100, 20, 0)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var result response.List[string]
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if result.Object != "list" {
		t.Errorf("expected object 'list', got '%s'", result.Object)
	}
	if len(result.Data) != 2 {
		t.Errorf("expected 2 items, got %d", len(result.Data))
	}
	if result.Total != 100 {
		t.Errorf("expected total 100, got %d", result.Total)
	}
	if !result.HasMore {
		t.Error("expected has_more to be true")
	}
}
