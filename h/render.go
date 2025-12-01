package h

import "io"

// Render writes the HTML representation of the given Builder to w.
// Returns nil if b is nil.
func Render(w io.Writer, b Builder) error {
	if b == nil {
		return nil
	}
	return b.Build(NewWriter(w))
}

// RenderIndent writes the HTML representation of the given Builder to w
// with indentation for readability. The indent parameter specifies the string
// to use for each indentation level (e.g., "  " for two spaces or "\t" for tabs).
// Returns nil if b is nil.
func RenderIndent(w io.Writer, indent string, b Builder) error {
	if b == nil {
		return nil
	}
	writer := NewWriter(w)
	writer.SetIndent(indent)
	return b.Build(writer)
}
