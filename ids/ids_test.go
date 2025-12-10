package ids_test

import (
	"testing"

	"github.com/doujins-org/ginapi/ids"
)

func TestFormatAndParse(t *testing.T) {
	tests := []struct {
		prefix string
		id     int64
	}{
		{"pay", 1},
		{"usr", 123},
		{"sub", 999999},
		{"gal", 0},
	}

	for _, tt := range tests {
		formatted := ids.Format(tt.prefix, tt.id)
		parsed, err := ids.Parse(tt.prefix, formatted)
		if err != nil {
			t.Errorf("Parse(%s, %s) error: %v", tt.prefix, formatted, err)
			continue
		}
		if parsed != tt.id {
			t.Errorf("Round trip failed: %d -> %s -> %d", tt.id, formatted, parsed)
		}
	}
}

func TestParseInvalidPrefix(t *testing.T) {
	_, err := ids.Parse("pay", "usr_123")
	if err == nil {
		t.Error("expected error for wrong prefix")
	}
}
