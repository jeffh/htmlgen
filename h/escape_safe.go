//go:build purego

package h

import (
	"io"
	"strings"
)

// writeEscapedString writes s to w with HTML escaping, avoiding allocations
// when no escaping is needed.
func writeEscapedString(w io.Writer, s string) error {
	if strings.ContainsAny(s, "&<>\"'") {
		// Slow path: use writeHTMLEscape which writes directly to w, allocating a byte slice
		return writeHTMLEscape(w, []byte(s))
	}
	// Fast path: No escaping needed
	_, err := io.WriteString(w, s)
	return err
}
