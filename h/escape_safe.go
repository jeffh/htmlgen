//go:build purego

package h

import (
	"html/template"
	"io"
)

// writeEscapedString writes s to w with HTML escaping, avoiding allocations
// when no escaping is needed.
func writeEscapedString(w io.Writer, s string) error {
	// Fast path: check if escaping is needed
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '&', '<', '>', '"', '\'':
			// Slow path: use template.HTMLEscape which writes directly to w
			template.HTMLEscape(w, []byte(s))
			return nil
		}
	}
	// No escaping needed
	_, err := io.WriteString(w, s)
	return err
}
