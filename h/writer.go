// Package h provides a low-level streaming HTML writer and a declarative builder API
// for programmatic HTML generation.
package h

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"sync"
)

// ErrUnknownTagToClose is returned when attempting to close a tag that was not opened.
var ErrUnknownTagToClose = errors.New("attempted to close tag not opened")

// writerPool pools Writer objects to reduce allocations.
var writerPool = sync.Pool{
	New: func() any {
		return &Writer{openTags: make([]string, 0, 32)}
	},
}

// getPooledWriter returns a Writer from the pool, configured to write to w.
func getPooledWriter(w io.Writer) *Writer {
	writer := writerPool.Get().(*Writer)
	writer.w = w
	writer.atLineStart = true
	return writer
}

// putPooledWriter returns a Writer to the pool after resetting it.
func putPooledWriter(w *Writer) {
	w.w = nil
	w.indent = ""
	w.indentCache = nil
	w.openTags = w.openTags[:0]
	w.atLineStart = false
	w.maxLineLen = 0
	writerPool.Put(w)
}

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

// NewWriter creates a new Writer that wraps the provided io.Writer.
// The Writer tracks open tags and provides methods for writing HTML elements.
func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w, openTags: make([]string, 0, 32), atLineStart: true}
}

// Writer is a low-level streaming HTML writer that wraps an io.Writer.
// It tracks open tags and provides methods for writing HTML elements,
// attributes, and content. Attribute values are automatically HTML-escaped.
type Writer struct {
	w           io.Writer
	indent      string
	indentCache []string // Cached indentation strings by depth
	openTags    []string
	atLineStart bool // Tracks if we're at the beginning of a line
	maxLineLen  int  // Max line length before wrapping attributes (0 = disabled)
}

// SetIndent sets the indentation prefix used for pretty-printing.
// When set to a non-empty string, each nested element will be indented
// by that prefix and newlines will be added after tags.
func (w *Writer) SetIndent(prefix string) {
	w.indent = prefix
}

// SetMaxLineLength sets the maximum line length before wrapping attributes
// to new lines. When set to 0 (default), attributes are never wrapped.
// When the combined tag + attributes would exceed this length, additional
// attributes are placed on new lines with extra indentation.
func (w *Writer) SetMaxLineLength(maxLen int) {
	w.maxLineLen = maxLen
}

func (w *Writer) isIndenting() bool { return len(w.indent) != 0 }

func (w *Writer) write(values ...string) error {
	for _, v := range values {
		_, err := io.WriteString(w.w, v)
		if err != nil {
			return err
		}
	}
	return nil
}

// Doctype writes the HTML5 doctype declaration (<!DOCTYPE html>).
func (w *Writer) Doctype() error { return w.write("<!DOCTYPE html>\n") }

func (w *Writer) writeIndentNewline() error {
	if w.isIndenting() {
		if err := w.write("\n"); err != nil {
			return err
		}
		w.atLineStart = true
	}
	return nil
}

func (w *Writer) writeIndent(modifier int) error {
	if !w.isIndenting() {
		return nil
	}
	depth := len(w.openTags) + modifier
	if depth <= 0 {
		return nil
	}
	// Grow cache if needed
	if depth > len(w.indentCache) {
		w.growIndentCache(depth)
	}
	_, err := io.WriteString(w.w, w.indentCache[depth-1])
	if err == nil {
		w.atLineStart = false
	}
	return err
}

func (w *Writer) growIndentCache(depth int) {
	for len(w.indentCache) < depth {
		var s string
		if len(w.indentCache) == 0 {
			s = w.indent
		} else {
			s = w.indentCache[len(w.indentCache)-1] + w.indent
		}
		w.indentCache = append(w.indentCache, s)
	}
}

// attrLen calculates the length of an attribute as it would be written.
// Returns the length of " name" or " name=\"value\"".
func attrLen(attr Attribute) int {
	if attr.Value == "" {
		return 1 + len(attr.Name) // " name"
	}
	return 1 + len(attr.Name) + 2 + len(attr.Value) + 1 // " name=\"value\""
}

// writeAttrs writes attributes, wrapping to new lines if maxLineLen is exceeded.
// lineLen is the current line length before attributes.
// Returns the final line length.
func (w *Writer) writeAttrs(as Attributes, lineLen int) (int, error) {
	for _, attr := range as {
		aLen := attrLen(attr)
		wrapped := false

		// Check if we need to wrap
		if w.maxLineLen > 0 && lineLen+aLen > w.maxLineLen {
			// Wrap: newline + extra indent (one deeper than current tag)
			if _, err := io.WriteString(w.w, "\n"); err != nil {
				return lineLen, err
			}
			// Write indent at depth+1 (inside the tag for attributes)
			depth := len(w.openTags) + 1
			if depth > len(w.indentCache) {
				w.growIndentCache(depth)
			}
			indent := w.indentCache[depth-1]
			if _, err := io.WriteString(w.w, indent); err != nil {
				return lineLen, err
			}
			lineLen = len(indent)
			wrapped = true
		}

		// Write the attribute (skip leading space if we just wrapped)
		if !wrapped {
			if _, err := io.WriteString(w.w, " "); err != nil {
				return lineLen, err
			}
		}
		if _, err := io.WriteString(w.w, attr.Name); err != nil {
			return lineLen, err
		}
		if attr.Value != "" {
			if _, err := io.WriteString(w.w, "=\""); err != nil {
				return lineLen, err
			}
			if err := writeEscapedString(w.w, attr.Value); err != nil {
				return lineLen, err
			}
			if _, err := io.WriteString(w.w, "\""); err != nil {
				return lineLen, err
			}
		}
		// Update line length (subtract 1 if we skipped the space)
		if wrapped {
			lineLen += aLen - 1
		} else {
			lineLen += aLen
		}
	}
	return lineLen, nil
}

// SelfClosingTag writes a self-closing HTML tag with the given name and attributes.
// For example, SelfClosingTag("br", nil) writes "<br/>".
func (w *Writer) SelfClosingTag(name string, as Attributes) error {
	if err := w.writeIndent(0); err != nil {
		return err
	}

	// Calculate initial line length (indent + "<" + name)
	lineLen := 1 + len(name)
	if w.isIndenting() {
		depth := len(w.openTags)
		if depth > 0 && depth <= len(w.indentCache) {
			lineLen += len(w.indentCache[depth-1])
		}
	}

	if _, err := io.WriteString(w.w, "<"); err != nil {
		return err
	}
	if _, err := io.WriteString(w.w, name); err != nil {
		return err
	}

	if _, err := w.writeAttrs(as, lineLen); err != nil {
		return err
	}

	if _, err := io.WriteString(w.w, "/>"); err != nil {
		return err
	}
	if err := w.writeIndentNewline(); err != nil {
		return err
	}
	return nil
}

// OpenTag writes an opening HTML tag with the given name and attributes.
// The tag is added to the stack of open tags and must be closed with CloseTag,
// CloseOneTag, or Close. Attribute values are automatically HTML-escaped.
func (w *Writer) OpenTag(name string, as Attributes) error {
	if err := w.writeIndent(0); err != nil {
		return err
	}

	// Calculate initial line length (indent + "<" + name)
	lineLen := 1 + len(name)
	if w.isIndenting() {
		depth := len(w.openTags)
		if depth > 0 && depth <= len(w.indentCache) {
			lineLen += len(w.indentCache[depth-1])
		}
	}

	if _, err := io.WriteString(w.w, "<"); err != nil {
		return err
	}
	if _, err := io.WriteString(w.w, name); err != nil {
		return err
	}

	if _, err := w.writeAttrs(as, lineLen); err != nil {
		return err
	}

	if _, err := io.WriteString(w.w, ">"); err != nil {
		return err
	}
	if err := w.writeIndentNewline(); err != nil {
		return err
	}
	w.openTags = append(w.openTags, name)
	return nil
}

// Text writes HTML-escaped text content.
// When indentation is enabled, text is indented at the current content depth
// and followed by a newline.
func (w *Writer) Text(txt string) error {
	if w.isIndenting() && w.atLineStart {
		if err := w.writeIndent(0); err != nil {
			return err
		}
	}
	if err := writeEscapedString(w.w, txt); err != nil {
		return err
	}
	if w.isIndenting() {
		w.atLineStart = false
		return w.writeIndentNewline()
	}
	return nil
}

// Raw writes unescaped HTML content. Use with caution as this can introduce
// XSS vulnerabilities if the content is not properly sanitized.
// When indentation is enabled, tracks whether content ends with newline.
func (w *Writer) Raw(unsafeHtml string) error {
	if err := w.write(unsafeHtml); err != nil {
		return err
	}
	if w.isIndenting() && len(unsafeHtml) > 0 {
		w.atLineStart = unsafeHtml[len(unsafeHtml)-1] == '\n'
	}
	return nil
}

// CloseTag closes the specified tag and all tags opened after it.
// Returns ErrUnknownTagToClose if no tags are open or the specified tag is not found.
func (w *Writer) CloseTag(name string) error {
	size := len(w.openTags)
	if size == 0 {
		return fmt.Errorf("%w: %s", ErrUnknownTagToClose, name)
	}
	// Ensure we're on a new line before closing tag
	if w.isIndenting() && !w.atLineStart {
		if err := w.write("\n"); err != nil {
			return err
		}
		w.atLineStart = true
	}
	if err := w.writeIndent(-1); err != nil {
		return err
	}
	for i := size - 1; i >= 0; i-- {
		if w.openTags[i] == name {
			for j := size - 1; j >= i; j-- {
				if _, err := io.WriteString(w.w, "</"); err != nil {
					return err
				}
				if _, err := io.WriteString(w.w, w.openTags[j]); err != nil {
					return err
				}
				if _, err := io.WriteString(w.w, ">"); err != nil {
					return err
				}
			}
			if err := w.writeIndentNewline(); err != nil {
				return err
			}
			w.openTags = w.openTags[:i]
			break
		}
	}
	return nil
}

// CloseOneTag closes the most recently opened tag.
// Returns ErrUnknownTagToClose if no tags are open.
func (w *Writer) CloseOneTag() error {
	size := len(w.openTags)
	if size == 0 {
		return ErrUnknownTagToClose
	}
	// Ensure we're on a new line before closing tag
	if w.isIndenting() && !w.atLineStart {
		if err := w.write("\n"); err != nil {
			return err
		}
		w.atLineStart = true
	}
	if err := w.writeIndent(-1); err != nil {
		return err
	}
	if _, err := io.WriteString(w.w, "</"); err != nil {
		return err
	}
	if _, err := io.WriteString(w.w, w.openTags[size-1]); err != nil {
		return err
	}
	if _, err := io.WriteString(w.w, ">"); err != nil {
		return err
	}
	if err := w.writeIndentNewline(); err != nil {
		return err
	}
	w.openTags = w.openTags[:size-1]
	return nil
}

// Close closes all remaining open tags in reverse order (most recent first).
func (w *Writer) Close() error {
	for i := len(w.openTags) - 1; i >= 0; i-- {
		// Ensure we're on a new line before closing tag
		if w.isIndenting() && !w.atLineStart {
			if err := w.write("\n"); err != nil {
				return err
			}
			w.atLineStart = true
		}
		if err := w.writeIndent(-1); err != nil {
			return err
		}
		if _, err := io.WriteString(w.w, "</"); err != nil {
			return err
		}
		if _, err := io.WriteString(w.w, w.openTags[i]); err != nil {
			return err
		}
		if _, err := io.WriteString(w.w, ">"); err != nil {
			return err
		}
		if err := w.writeIndentNewline(); err != nil {
			return err
		}
	}
	w.openTags = nil
	return nil
}
