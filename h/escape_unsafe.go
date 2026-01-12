//go:build !purego

package h

import (
	"io"
	"unsafe"
)

// writeEscapedString writes s to w with HTML escaping, avoiding allocations
// when no escaping is needed.
func writeEscapedString(w io.Writer, s string) error {
	L := len(s)
	if L > 0 {
		// This is only safe because writeHTMLEscape only reads from the input slice
		b := unsafe.Slice(unsafe.StringData(s), L)
		return writeHTMLEscape(w, b)
	}
	return nil
}
