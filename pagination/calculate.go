package pagination

// HasMore returns true if there are more items after the current page.
func HasMore(offset, limit int, total int64) bool {
	return int64(offset+limit) < total
}

// HasMoreFromLen returns true if there are more items, using actual result length.
// This is useful when you fetch limit+1 items to check for more.
func HasMoreFromLen(offset, resultLen int, total int64) bool {
	return int64(offset+resultLen) < total
}
