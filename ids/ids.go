// Package ids provides Stripe-style ID formatting utilities.
// IDs are formatted as prefix_base62encodedID (e.g., "pay_7B", "usr_1A3").
package ids

import (
	"fmt"
	"strconv"
	"strings"
)

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// Format creates a prefixed ID string from a numeric ID.
// Example: Format("pay", 123) returns "pay_1Z"
func Format(prefix string, id int64) string {
	return prefix + "_" + encodeBase62(id)
}

// Parse extracts the numeric ID from a prefixed ID string.
// Example: Parse("pay", "pay_1Z") returns 123, nil
func Parse(prefix, prefixedID string) (int64, error) {
	expectedPrefix := prefix + "_"
	if !strings.HasPrefix(prefixedID, expectedPrefix) {
		return 0, fmt.Errorf("invalid prefix: expected %s, got %s", expectedPrefix, prefixedID)
	}

	encoded := strings.TrimPrefix(prefixedID, expectedPrefix)
	return decodeBase62(encoded)
}

// MustParse is like Parse but panics on error.
func MustParse(prefix, prefixedID string) int64 {
	id, err := Parse(prefix, prefixedID)
	if err != nil {
		panic(err)
	}
	return id
}

// encodeBase62 encodes a number to base62.
func encodeBase62(num int64) string {
	if num == 0 {
		return "0"
	}

	var result []byte
	for num > 0 {
		result = append([]byte{base62Chars[num%62]}, result...)
		num /= 62
	}
	return string(result)
}

// decodeBase62 decodes a base62 string to a number.
func decodeBase62(s string) (int64, error) {
	var result int64
	for _, c := range s {
		idx := strings.IndexRune(base62Chars, c)
		if idx == -1 {
			return 0, fmt.Errorf("invalid base62 character: %c", c)
		}
		result = result*62 + int64(idx)
	}
	return result, nil
}

// Common ID prefixes used across projects
const (
	PrefixUser         = "usr"
	PrefixPayment      = "pay"
	PrefixSubscription = "sub"
	PrefixProduct      = "prod"
	PrefixPrice        = "price"
	PrefixGallery      = "gal"
	PrefixArtist       = "art"
	PrefixTag          = "tag"
	PrefixSeries       = "ser"
	PrefixCharacter    = "chr"
)

// Formatter provides type-safe ID formatting for a specific entity type.
type Formatter struct {
	Prefix string
}

// NewFormatter creates a new ID formatter with the given prefix.
func NewFormatter(prefix string) *Formatter {
	return &Formatter{Prefix: prefix}
}

// Format creates a prefixed ID string.
func (f *Formatter) Format(id int64) string {
	return Format(f.Prefix, id)
}

// Parse extracts the numeric ID from a prefixed string.
func (f *Formatter) Parse(prefixedID string) (int64, error) {
	return Parse(f.Prefix, prefixedID)
}

// FormatString creates a prefixed ID from a string ID (for UUIDs, etc).
func FormatString(prefix, id string) string {
	return prefix + "_" + id
}

// ParseString extracts the ID portion from a prefixed string ID.
func ParseString(prefix, prefixedID string) (string, error) {
	expectedPrefix := prefix + "_"
	if !strings.HasPrefix(prefixedID, expectedPrefix) {
		return "", fmt.Errorf("invalid prefix: expected %s", expectedPrefix)
	}
	return strings.TrimPrefix(prefixedID, expectedPrefix), nil
}

// FormatInt64 is an alias for Format (for clarity when working with int64).
func FormatInt64(prefix string, id int64) string {
	return Format(prefix, id)
}

// FormatInt formats an int as a prefixed ID.
func FormatInt(prefix string, id int) string {
	return Format(prefix, int64(id))
}

// ParseInt64 is an alias for Parse (for clarity).
func ParseInt64(prefix, prefixedID string) (int64, error) {
	return Parse(prefix, prefixedID)
}

// ParseIntParam parses a string that could be either a prefixed ID or raw numeric ID.
// This is useful for APIs that accept both formats.
func ParseIntParam(prefix, param string) (int64, error) {
	// Try prefixed format first
	if strings.HasPrefix(param, prefix+"_") {
		return Parse(prefix, param)
	}
	// Fall back to raw numeric
	return strconv.ParseInt(param, 10, 64)
}
