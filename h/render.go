package h

import (
	"bytes"
	"io"
	"strings"
)

// Render writes the HTML representation of the given Builder to w.
// Returns nil if b is nil.
func Render(w io.Writer, b Builder) error {
	if b == nil {
		return nil
	}
	writer := getPooledWriter(w)
	err := b.Build(writer)
	putPooledWriter(writer)
	return err
}

// RenderIndent writes the HTML representation of the given Builder to w
// with indentation for readability. The indent parameter specifies the string
// to use for each indentation level (e.g., "  " for two spaces or "\t" for tabs).
// Returns nil if b is nil.
func RenderIndent(w io.Writer, indent string, b Builder) error {
	if b == nil {
		return nil
	}
	writer := getPooledWriter(w)
	writer.SetIndent(indent)
	err := b.Build(writer)
	putPooledWriter(writer)
	return err
}

// RenderString renders the Builder and returns the result as a string.
// Returns an empty string if b is nil. Panics if rendering fails.
//
//	html := h.RenderString(h.Div(h.Text("Hello")))
func RenderString(b Builder) string {
	if b == nil {
		return ""
	}
	var sb strings.Builder
	writer := getPooledWriter(&sb)
	err := b.Build(writer)
	putPooledWriter(writer)
	if err != nil {
		panic(err)
	}
	return sb.String()
}

// RenderBytes renders the Builder and returns the result as a byte slice.
// Returns nil if b is nil. Panics if rendering fails.
//
//	html := h.RenderBytes(h.Div(h.Text("Hello")))
func RenderBytes(b Builder) []byte {
	if b == nil {
		return nil
	}
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	writer := getPooledWriter(buf)
	err := b.Build(writer)
	putPooledWriter(writer)
	if err != nil {
		bufPool.Put(buf)
		panic(err)
	}
	result := make([]byte, buf.Len())
	copy(result, buf.Bytes())
	bufPool.Put(buf)
	return result
}
