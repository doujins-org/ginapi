package pagination_test

import (
	"testing"

	"github.com/doujins-org/ginapi/pagination"
)

func TestHasMore(t *testing.T) {
	tests := []struct {
		name   string
		offset int
		limit  int
		total  int64
		want   bool
	}{
		{
			name:   "first page with more",
			offset: 0,
			limit:  20,
			total:  100,
			want:   true,
		},
		{
			name:   "last page",
			offset: 80,
			limit:  20,
			total:  100,
			want:   false,
		},
		{
			name:   "exact boundary",
			offset: 80,
			limit:  20,
			total:  100,
			want:   false,
		},
		{
			name:   "empty result",
			offset: 0,
			limit:  20,
			total:  0,
			want:   false,
		},
		{
			name:   "partial last page",
			offset: 90,
			limit:  20,
			total:  100,
			want:   false,
		},
		{
			name:   "middle page",
			offset: 40,
			limit:  20,
			total:  100,
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pagination.HasMore(tt.offset, tt.limit, tt.total)
			if got != tt.want {
				t.Errorf("HasMore(%d, %d, %d) = %v, want %v",
					tt.offset, tt.limit, tt.total, got, tt.want)
			}
		})
	}
}

func TestHasMoreFromLen(t *testing.T) {
	tests := []struct {
		name      string
		offset    int
		resultLen int
		total     int64
		want      bool
	}{
		{
			name:      "full page with more",
			offset:    0,
			resultLen: 20,
			total:     100,
			want:      true,
		},
		{
			name:      "partial page (last)",
			offset:    90,
			resultLen: 10,
			total:     100,
			want:      false,
		},
		{
			name:      "empty result",
			offset:    0,
			resultLen: 0,
			total:     0,
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pagination.HasMoreFromLen(tt.offset, tt.resultLen, tt.total)
			if got != tt.want {
				t.Errorf("HasMoreFromLen(%d, %d, %d) = %v, want %v",
					tt.offset, tt.resultLen, tt.total, got, tt.want)
			}
		})
	}
}
