package h

import (
	"bytes"
	"io"
	"strings"
	"sync"
)

var builderPool = sync.Pool{
	New: func() any {
		return &B{openTags: make([]string, 0, 32)}
	},
}

var bufferPool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

func getBuilder(w io.Writer, indent string) *B {
	b := builderPool.Get().(*B)
	b.w = w
	b.indent = indent
	b.atLineStart = true
	return b
}

func putBuilder(b *B) {
	b.w = nil
	b.openTags = b.openTags[:0]
	b.indent = ""
	b.indentCache = b.indentCache[:0]
	b.atLineStart = false
	b.maxLineLen = 0
	b.err = nil
	builderPool.Put(b)
}

// Render runs fn against a fresh B writing to w and returns the first write
// error. Any tags left open by fn are closed defensively.
func Render(w io.Writer, fn func(*B)) error {
	b := getBuilder(w, "")
	defer putBuilder(b)
	if fn != nil {
		fn(b)
	}
	b.closeAll()
	return b.err
}

// RenderIndent runs fn with pretty-printing using indent for each nesting level.
func RenderIndent(w io.Writer, indent string, fn func(*B)) error {
	b := getBuilder(w, indent)
	defer putBuilder(b)
	if fn != nil {
		fn(b)
	}
	b.closeAll()
	return b.err
}

// RenderString renders fn to a string and panics if rendering fails.
func RenderString(fn func(*B)) string {
	var result strings.Builder
	if err := Render(&result, fn); err != nil {
		panic(err)
	}
	return result.String()
}

// RenderBytes renders fn to a byte slice and panics if rendering fails.
func RenderBytes(fn func(*B)) []byte {
	buffer := bufferPool.Get().(*bytes.Buffer)
	buffer.Reset()
	defer bufferPool.Put(buffer)
	if err := Render(buffer, fn); err != nil {
		panic(err)
	}
	return append([]byte(nil), buffer.Bytes()...)
}
