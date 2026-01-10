package h

import "bytes"

// bakedBuilder stores pre-rendered HTML bytes for fast rendering.
type bakedBuilder struct {
	html []byte
}

func (b *bakedBuilder) isTagArg() {}

func (b *bakedBuilder) Build(w *Writer) error {
	_, err := w.w.Write(b.html)
	return err
}

// Bake pre-renders a Builder to bytes for faster subsequent renders.
// The resulting Builder writes the cached bytes directly without
// tree traversal, achieving performance comparable to html/template.
//
// Use Bake for static content that is rendered frequently:
//
//	// Bake once at startup
//	header := h.Bake(h.Header(h.H1(h.Text("My Site"))))
//
//	// Fast renders afterward
//	h.Render(w, header)
//
// Note: Baked builders ignore indentation settings since the output
// is pre-computed. For pretty-printed cached output, render with
// RenderIndent first, then bake the result using Raw.
func Bake(b Builder) Builder {
	if b == nil {
		return nil
	}
	var buf bytes.Buffer
	Render(&buf, b)
	return &bakedBuilder{html: buf.Bytes()}
}
