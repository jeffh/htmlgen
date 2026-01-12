//go:build !purego

package h

import (
	"html/template"
	"io"
	"unsafe"
)

// writeEscapedString writes s to w with HTML escaping, avoiding allocations
// when no escaping is needed.
func writeEscapedString(w io.Writer, s string) error {
	if len(s) > 0 {
		// This is only safe because HTMLEscape only reads from the input slice
		b := unsafe.Slice(unsafe.StringData(s), len(s))
		template.HTMLEscape(w, b)
	}
	return nil
}
