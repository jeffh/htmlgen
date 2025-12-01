// Package h provides a low-level streaming HTML writer and a declarative builder API
// for programmatic HTML generation.
package h

import (
	"errors"
	"fmt"
	"html/template"
	"io"
)

// ErrUnknownTagToClose is returned when attempting to close a tag that was not opened.
var ErrUnknownTagToClose = errors.New("attempted to close tag not opened")

// NewWriter creates a new Writer that wraps the provided io.Writer.
// The Writer tracks open tags and provides methods for writing HTML elements.
func NewWriter(w io.Writer) *Writer {
	return &Writer{w, "", make([]string, 0, 16)}
}

// Writer is a low-level streaming HTML writer that wraps an io.Writer.
// It tracks open tags and provides methods for writing HTML elements,
// attributes, and content. Attribute values are automatically HTML-escaped.
type Writer struct {
	w        io.Writer
	indent   string
	openTags []string
}

// SetIndent sets the indentation prefix used for pretty-printing.
// When set to a non-empty string, each nested element will be indented
// by that prefix and newlines will be added after tags.
func (w *Writer) SetIndent(prefix string) {
	w.indent = prefix
}

func (w *Writer) isIndenting() bool { return len(w.indent) != 0 }

func (w *Writer) write(values ...string) error {
	for _, v := range values {
		_, err := w.w.Write([]byte(v))
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
		return w.write("\n")
	}
	return nil
}

func (w *Writer) writeIndent(modifier int) error {
	if w.isIndenting() {
		L := len(w.openTags) + modifier
		for i := 0; i < L; i++ {
			if err := w.write(w.indent); err != nil {
				return err
			}
		}
	}
	return nil
}

// SelfClosingTag writes a self-closing HTML tag with the given name and attributes.
// For example, SelfClosingTag("br", nil) writes "<br/>".
func (w *Writer) SelfClosingTag(name string, as Attributes) error {
	if err := w.writeIndent(0); err != nil {
		return err
	}
	if err := w.write("<", name); err != nil {
		return err
	}
	for _, attr := range as {
		val := attr.Value
		if val == "" {
			if err := w.write(" ", attr.Name); err != nil {
				return err
			}
		} else {
			if err := w.write(" ", attr.Name, "=\"", template.HTMLEscapeString(val), "\""); err != nil {
				return err
			}
		}
	}
	if err := w.write("/>"); err != nil {
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
	if err := w.write("<", name); err != nil {
		return err
	}
	for _, attr := range as {
		val := attr.Value
		if val == "" {
			if err := w.write(" ", attr.Name); err != nil {
				return err
			}
		} else {
			if err := w.write(" ", attr.Name, "=\"", template.HTMLEscapeString(val), "\""); err != nil {
				return err
			}
		}
	}
	if err := w.write(">"); err != nil {
		return err
	}
	if err := w.writeIndentNewline(); err != nil {
		return err
	}
	w.openTags = append(w.openTags, name)
	return nil
}

// Text writes HTML-escaped text content.
func (w *Writer) Text(txt string) error { return w.write(template.HTMLEscapeString(txt)) }

// Raw writes unescaped HTML content. Use with caution as this can introduce
// XSS vulnerabilities if the content is not properly sanitized.
func (w *Writer) Raw(unsafeHtml string) error { return w.write(unsafeHtml) }

// CloseTag closes the specified tag and all tags opened after it.
// Returns ErrUnknownTagToClose if no tags are open or the specified tag is not found.
func (w *Writer) CloseTag(name string) error {
	size := len(w.openTags)
	if size == 0 {
		return fmt.Errorf("%w: %s", ErrUnknownTagToClose, name)
	}
	if err := w.writeIndent(-1); err != nil {
		return err
	}
	for i := size - 1; i >= 0; i-- {
		if w.openTags[i] == name {
			for j := size - 1; j >= i; j-- {
				if err := w.write("</", w.openTags[j], ">"); err != nil {
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
	if err := w.writeIndent(-1); err != nil {
		return err
	}
	err := w.write("</", w.openTags[size-1], ">")
	if err != nil {
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
		if err := w.writeIndent(-1); err != nil {
			return err
		}
		if err := w.write("</", w.openTags[i], ">"); err != nil {
			return err
		}
		if err := w.writeIndentNewline(); err != nil {
			return err
		}
	}
	w.openTags = nil
	return nil
}
