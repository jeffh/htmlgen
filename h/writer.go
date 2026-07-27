// Package h provides streaming, imperative HTML generation.
package h

import (
	"fmt"
	"io"
)

// flushThreshold is the buffered byte count that triggers a write to the
// underlying io.Writer.
const flushThreshold = 4 << 10

// B is a streaming HTML builder bound to an io.Writer. It tracks open tags,
// optional pretty-printing, and a sticky write error. Output accumulates in an
// internal buffer and is written to the io.Writer in chunks. A B is not safe
// for concurrent use by multiple goroutines; each Render gets its own B.
type B struct {
	w           io.Writer
	buf         []byte
	scratch     []byte
	openTags    []string
	indent      string
	indentCache []string
	atLineStart bool
	maxLineLen  int
	err         error
}

// Err returns the first write error encountered, if any. Because output is
// buffered, Err reflects only the errors observed as of the last flush; a
// failing writer may not surface until the buffer fills or Flush is called. In
// a long streaming loop, call Flush at loop boundaries to observe write errors
// promptly. Render and RenderIndent flush before returning, so their return
// value always reports the write error.
func (b *B) Err() error {
	return b.err
}

// Flush writes any buffered output to the underlying io.Writer and returns the
// sticky error. Render and RenderIndent flush on completion; call Flush
// explicitly when bytes must reach the client mid-render, as with server-sent
// events or a long streaming response. After a failed flush all further output
// is discarded.
//
// A B is valid only for the duration of the Render call that created it, since
// it is returned to a pool afterwards. Do not retain it, and do not call Flush
// once that Render has returned.
func (b *B) Flush() error {
	if b.err != nil || len(b.buf) == 0 {
		return b.err
	}
	_, err := b.w.Write(b.buf)
	b.buf = b.buf[:0]
	b.setErr(err)
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

func (b *B) maybeFlush() {
	if len(b.buf) >= flushThreshold {
		b.Flush()
	}
}

func (b *B) writeString(value string) {
	if b.err != nil {
		return
	}
	b.buf = append(b.buf, value...)
	b.maybeFlush()
}

func (b *B) writeEscaped(value string) {
	if b.err != nil {
		return
	}
	if len(value) >= flushThreshold {
		b.writeEscapedLarge(value)
		return
	}
	b.buf = appendEscaped(b.buf, value)
	b.maybeFlush()
}

// writeEscapedLarge escapes value a chunk of input at a time, writing each
// chunk straight through, so a multi-megabyte value costs a bounded amount of
// memory instead of an escaped copy of the whole input. scratch holds one
// escaped chunk, at most 5x the chunk size because that is the longest
// replacement.
func (b *B) writeEscapedLarge(value string) {
	if b.Flush() != nil {
		return
	}
	for len(value) > 0 {
		chunk := value
		if len(chunk) > flushThreshold {
			chunk = chunk[:flushThreshold]
		}
		value = value[len(chunk):]
		b.scratch = appendEscaped(b.scratch[:0], chunk)
		if _, err := b.w.Write(b.scratch); err != nil {
			b.setErr(err)
			return
		}
	}
}

// escapeTable maps each byte to its HTML replacement; an empty string means the
// byte is written unchanged.
var escapeTable = func() (table [256]string) {
	table['\000'] = "\uFFFD"
	table['"'] = "&#34;"
	table['\''] = "&#39;"
	table['&'] = "&amp;"
	table['<'] = "&lt;"
	table['>'] = "&gt;"
	return
}()

func appendEscaped(dst []byte, value string) []byte {
	last := 0
	for i := 0; i < len(value); i++ {
		escaped := escapeTable[value[i]]
		if escaped == "" {
			continue
		}
		dst = append(dst, value[last:i]...)
		dst = append(dst, escaped...)
		last = i + 1
	}
	return append(dst, value[last:]...)
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
			b.buf = append(b.buf, ' ')
		}
		b.buf = append(b.buf, attr.Name...)
		if attr.Value != "" {
			b.buf = append(b.buf, '=', '"')
			if len(attr.Value) >= flushThreshold {
				b.writeEscapedLarge(attr.Value)
				if b.err != nil {
					return
				}
			} else {
				b.buf = appendEscaped(b.buf, attr.Value)
			}
			b.buf = append(b.buf, '"')
		}
		b.maybeFlush()
		if wrapped {
			lineLen += length - 1
		} else {
			lineLen += length
		}
	}
}

// openTag writes an opening tag and pushes close onto the open-tag stack. open
// is the tag's literal prefix ("<div") and close its literal end tag
// ("</div>"), both precomputed so writing a tag is a plain append.
func (b *B) openTag(open, close string, attrs Attributes) {
	if b.err != nil {
		return
	}
	lineLen := len(open)
	if b.isIndenting() {
		b.writeIndent(0)
		depth := len(b.openTags)
		if depth > 0 && depth <= len(b.indentCache) {
			lineLen += len(b.indentCache[depth-1])
		}
	}

	b.buf = append(b.buf, open...)
	if len(attrs) > 0 {
		b.writeAttrs(attrs, lineLen)
	}
	b.buf = append(b.buf, '>')
	b.maybeFlush()
	b.writeIndentNewline()
	if b.err == nil {
		b.openTags = append(b.openTags, close)
	}
}

func (b *B) voidTag(open string, attrs Attributes) {
	if b.err != nil {
		return
	}
	lineLen := len(open)
	if b.isIndenting() {
		b.writeIndent(0)
		depth := len(b.openTags)
		if depth > 0 && depth <= len(b.indentCache) {
			lineLen += len(b.indentCache[depth-1])
		}
	}

	b.buf = append(b.buf, open...)
	if len(attrs) > 0 {
		b.writeAttrs(attrs, lineLen)
	}
	b.buf = append(b.buf, '/', '>')
	b.maybeFlush()
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
	b.buf = append(b.buf, b.openTags[size-1]...)
	b.maybeFlush()
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

// element writes a container element. attrs and body arrive with concrete
// types, so no argument is boxed and a capturing body closure stays on the
// stack.
func (b *B) element(open, close string, attrs Attributes, body Body) {
	if b.err != nil {
		return
	}
	b.openTag(open, close, attrs)
	if b.err != nil {
		return
	}
	if body != nil {
		body(b)
	}
	b.closeOneTag()
}

func (b *B) voidElement(open string, attrs Attributes) {
	if b.err != nil {
		return
	}
	b.voidTag(open, attrs)
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
	b.writeEscaped(value)
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
	if len(value) >= flushThreshold {
		// Written straight through so one large value cannot inflate the
		// pooled buffer past its working size.
		if b.Flush() == nil {
			b.setErr(b.writeStringDirect(value))
		}
	} else {
		b.writeString(value)
	}
	if b.isIndenting() && b.err == nil && value != "" {
		b.atLineStart = value[len(value)-1] == '\n'
	}
}

// writeStringDirect writes value straight to the underlying writer, bypassing
// the buffer. When w lacks io.StringWriter, the []byte conversion happens one
// scratch-sized chunk at a time so a multi-megabyte value never allocates a
// full copy.
func (b *B) writeStringDirect(value string) error {
	if sw, ok := b.w.(io.StringWriter); ok {
		_, err := sw.WriteString(value)
		return err
	}
	for len(value) > 0 {
		chunk := value
		if len(chunk) > flushThreshold {
			chunk = chunk[:flushThreshold]
		}
		b.scratch = append(b.scratch[:0], chunk...)
		if _, err := b.w.Write(b.scratch); err != nil {
			return err
		}
		value = value[len(chunk):]
	}
	return nil
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
func (b *B) El(name string, attrs Attributes, body Body) {
	validateTagName(name)
	b.element("<"+name, "</"+name+">", attrs, body)
}

// VoidEl writes an arbitrary self-closing element.
// Panics if name is not a valid element name (see El).
func (b *B) VoidEl(name string, attrs Attributes) {
	validateTagName(name)
	b.voidElement("<"+name, attrs)
}
