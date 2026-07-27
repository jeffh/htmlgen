// Package h provides streaming, imperative HTML generation.
package h

import (
	"fmt"
	"io"
)

// B is a streaming HTML builder bound to an io.Writer. It tracks open tags,
// optional pretty-printing, and a sticky write error. A B is not safe for
// concurrent use by multiple goroutines; each Render gets its own B.
type B struct {
	w           io.Writer
	openTags    []string
	indent      string
	indentCache []string
	atLineStart bool
	maxLineLen  int
	err         error
}

// Err returns the first write error encountered, if any.
func (b *B) Err() error {
	return b.err
}

// SetMaxLineLength sets the maximum line length before attributes wrap onto
// new lines when pretty-printing. 0 (the default) disables wrapping. Call it
// at the top of a RenderIndent body; it has no visible effect without
// indentation.
func (b *B) SetMaxLineLength(maxLen int) {
	b.maxLineLen = maxLen
}

func (b *B) setErr(err error) {
	if b.err == nil && err != nil {
		b.err = err
	}
}

func (b *B) writeString(value string) {
	if b.err != nil {
		return
	}
	_, err := io.WriteString(b.w, value)
	b.setErr(err)
}

func (b *B) isIndenting() bool {
	return b.indent != ""
}

func (b *B) writeIndentNewline() {
	if !b.isIndenting() || b.err != nil {
		return
	}
	b.writeString("\n")
	if b.err == nil {
		b.atLineStart = true
	}
}

func (b *B) writeIndent(modifier int) {
	if !b.isIndenting() || b.err != nil {
		return
	}
	depth := len(b.openTags) + modifier
	if depth <= 0 {
		return
	}
	if depth > len(b.indentCache) {
		b.growIndentCache(depth)
	}
	b.writeString(b.indentCache[depth-1])
	if b.err == nil {
		b.atLineStart = false
	}
}

func (b *B) growIndentCache(depth int) {
	for len(b.indentCache) < depth {
		value := b.indent
		if len(b.indentCache) > 0 {
			value = b.indentCache[len(b.indentCache)-1] + b.indent
		}
		b.indentCache = append(b.indentCache, value)
	}
}

func attrLen(attr Attribute) int {
	if attr.Value == "" {
		return 1 + len(attr.Name)
	}
	return 1 + len(attr.Name) + 2 + len(attr.Value) + 1
}

func (b *B) writeAttrs(attrs Attributes, lineLen int) {
	for _, attr := range attrs {
		if b.err != nil {
			return
		}
		if attr.Name == "" {
			continue
		}

		length := attrLen(attr)
		wrapped := false
		if b.maxLineLen > 0 && b.isIndenting() && lineLen+length > b.maxLineLen {
			b.writeString("\n")
			if b.err != nil {
				return
			}
			depth := len(b.openTags) + 1
			if depth > len(b.indentCache) {
				b.growIndentCache(depth)
			}
			indent := b.indentCache[depth-1]
			b.writeString(indent)
			lineLen = len(indent)
			wrapped = true
		}

		if !wrapped {
			b.writeString(" ")
		}
		b.writeString(attr.Name)
		if attr.Value != "" {
			b.writeString("=\"")
			if b.err == nil {
				b.setErr(writeEscapedString(b.w, attr.Value))
			}
			b.writeString("\"")
		}
		if wrapped {
			lineLen += length - 1
		} else {
			lineLen += length
		}
	}
}

func (b *B) openTag(name string, attrs Attributes) {
	if b.err != nil {
		return
	}
	b.writeIndent(0)

	lineLen := 1 + len(name)
	if b.isIndenting() {
		depth := len(b.openTags)
		if depth > 0 && depth <= len(b.indentCache) {
			lineLen += len(b.indentCache[depth-1])
		}
	}

	b.writeString("<")
	b.writeString(name)
	b.writeAttrs(attrs, lineLen)
	b.writeString(">")
	b.writeIndentNewline()
	if b.err == nil {
		b.openTags = append(b.openTags, name)
	}
}

func (b *B) voidTag(name string, attrs Attributes) {
	if b.err != nil {
		return
	}
	b.writeIndent(0)

	lineLen := 1 + len(name)
	if b.isIndenting() {
		depth := len(b.openTags)
		if depth > 0 && depth <= len(b.indentCache) {
			lineLen += len(b.indentCache[depth-1])
		}
	}

	b.writeString("<")
	b.writeString(name)
	b.writeAttrs(attrs, lineLen)
	b.writeString("/>")
	b.writeIndentNewline()
}

func (b *B) closeOneTag() {
	if b.err != nil {
		return
	}
	size := len(b.openTags)
	if size == 0 {
		return
	}
	if b.isIndenting() && !b.atLineStart {
		b.writeString("\n")
		if b.err != nil {
			return
		}
		b.atLineStart = true
	}
	b.writeIndent(-1)
	b.writeString("</")
	b.writeString(b.openTags[size-1])
	b.writeString(">")
	b.writeIndentNewline()
	if b.err == nil {
		b.openTags = b.openTags[:size-1]
	}
}

func (b *B) closeAll() {
	for b.err == nil && len(b.openTags) > 0 {
		b.closeOneTag()
	}
}

func (b *B) element(name string, args ...any) {
	if b.err != nil {
		return
	}
	attrs, body := parseArgs(name, args)
	b.openTag(name, attrs)
	if b.err != nil {
		return
	}
	if body != nil {
		body(b)
	}
	b.closeOneTag()
}

func (b *B) voidElement(name string, args ...any) {
	if b.err != nil {
		return
	}
	attrs, _ := parseArgs(name, args)
	b.voidTag(name, attrs)
}

// Doctype writes the HTML5 doctype declaration.
func (b *B) Doctype() {
	b.writeString("<!DOCTYPE html>\n")
}

// Text writes HTML-escaped text.
func (b *B) Text(value string) {
	if b.err != nil {
		return
	}
	if b.isIndenting() && b.atLineStart {
		b.writeIndent(0)
	}
	if b.err == nil {
		b.setErr(writeEscapedString(b.w, value))
	}
	if b.isIndenting() && b.err == nil {
		b.atLineStart = false
		b.writeIndentNewline()
	}
}

// Textf writes HTML-escaped formatted text.
func (b *B) Textf(format string, args ...any) {
	if b.err == nil {
		b.Text(fmt.Sprintf(format, args...))
	}
}

// Raw writes unescaped HTML. The caller must ensure value is safe.
func (b *B) Raw(value string) {
	if b.err != nil {
		return
	}
	b.writeString(value)
	if b.isIndenting() && b.err == nil && value != "" {
		b.atLineStart = value[len(value)-1] == '\n'
	}
}

// Rawf writes unescaped formatted HTML. The caller must ensure the result is safe.
func (b *B) Rawf(format string, args ...any) {
	if b.err == nil {
		b.Raw(fmt.Sprintf(format, args...))
	}
}

// El writes an arbitrary container element.
// Panics if name is not an ASCII letter followed by ASCII letters, digits,
// '_', '.', ':', or '-'; tag names must never come from untrusted input.
func (b *B) El(name string, args ...any) {
	validateTagName(name)
	b.element(name, args...)
}

// VoidEl writes an arbitrary self-closing element.
// Panics if name is not a valid element name (see El).
func (b *B) VoidEl(name string, args ...any) {
	validateTagName(name)
	b.voidElement(name, args...)
}

// copied from text/template.HTMLEscape so escaping errors can be returned
var (
	htmlQuot = []byte("&#34;")
	htmlApos = []byte("&#39;")
	htmlAmp  = []byte("&amp;")
	htmlLt   = []byte("&lt;")
	htmlGt   = []byte("&gt;")
	htmlNull = []byte("\uFFFD")
)

func writeHTMLEscape(w io.Writer, value []byte) error {
	last := 0
	for i, char := range value {
		var escaped []byte
		switch char {
		case '\000':
			escaped = htmlNull
		case '"':
			escaped = htmlQuot
		case '\'':
			escaped = htmlApos
		case '&':
			escaped = htmlAmp
		case '<':
			escaped = htmlLt
		case '>':
			escaped = htmlGt
		default:
			continue
		}
		if _, err := w.Write(value[last:i]); err != nil {
			return err
		}
		if _, err := w.Write(escaped); err != nil {
			return err
		}
		last = i + 1
	}
	_, err := w.Write(value[last:])
	return err
}
