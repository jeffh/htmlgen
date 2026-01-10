package h

import (
	"bytes"
	"io"
)

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

// Param is a placeholder for dynamic content in a parameterized template.
// Use with BakeParams to create templates with variable content.
type Param struct {
	name string
}

func (p *Param) isTagArg() {}

// Build signals to BakeParams that a parameter slot was encountered.
// When used outside BakeParams, it renders nothing.
func (p *Param) Build(w *Writer) error {
	if bw, ok := w.w.(*bakeWriter); ok {
		bw.recordParam(p)
	}
	return nil
}

// Value binds a Builder to this parameter for rendering.
func (p *Param) Value(b Builder) ParamValue {
	return ParamValue{param: p, value: b}
}

// ParamValue binds a value to a parameter for rendering.
type ParamValue struct {
	param *Param
	value Builder
}

// NewParam creates a named placeholder for dynamic content in BakeParams.
//
//	title := h.NewParam("title")
//	tmpl := h.BakeParams(h.Div(h.H1(title)))
//	tmpl.Render(w, title.Value(h.Text("Hello")))
func NewParam(name string) *Param {
	return &Param{name: name}
}

// BakedTemplate is a pre-computed template with parameter placeholders.
// Created by BakeParams, it stores static HTML segments and parameter positions.
type BakedTemplate struct {
	segments [][]byte
	params   []*Param
}

// bakeWriter captures segments between parameters during BakeParams.
type bakeWriter struct {
	buf      bytes.Buffer
	segments [][]byte
	params   []*Param
}

func (bw *bakeWriter) Write(p []byte) (int, error) {
	return bw.buf.Write(p)
}

func (bw *bakeWriter) WriteString(s string) (int, error) {
	return bw.buf.WriteString(s)
}

func (bw *bakeWriter) recordParam(p *Param) {
	// Copy current buffer contents to a new segment
	segment := make([]byte, bw.buf.Len())
	copy(segment, bw.buf.Bytes())
	bw.segments = append(bw.segments, segment)
	bw.buf.Reset()
	bw.params = append(bw.params, p)
}

// BakeParams pre-computes a Builder with Slot placeholders for dynamic content.
// Static HTML is pre-rendered into segments, with slots marking where dynamic
// content will be inserted at render time.
//
//	title := h.NewSlot("title")
//	content := h.NewSlot("content")
//
//	tmpl := h.BakeParams(h.Div(
//	    h.H1(title),
//	    h.P(content),
//	))
//
//	// Render with values
//	tmpl.Render(w,
//	    title.Value(h.Text("Hello World")),
//	    content.Value(h.Text("Welcome!")),
//	)
func BakeParams(b Builder) *BakedTemplate {
	if b == nil {
		return &BakedTemplate{}
	}
	bw := &bakeWriter{}
	w := &Writer{w: bw, openTags: make([]string, 0, 32)}
	b.Build(w)

	// Add final segment
	segment := make([]byte, bw.buf.Len())
	copy(segment, bw.buf.Bytes())
	bw.segments = append(bw.segments, segment)

	return &BakedTemplate{
		segments: bw.segments,
		params:   bw.params,
	}
}

// With creates a Builder with parameter values bound for rendering.
func (t *BakedTemplate) With(values ...ParamValue) Builder {
	return &boundTemplate{template: t, values: values}
}

// Render writes the template directly with parameter values.
func (t *BakedTemplate) Render(w io.Writer, values ...ParamValue) error {
	return Render(w, t.With(values...))
}

// boundTemplate is a BakedTemplate with parameter values bound.
type boundTemplate struct {
	template *BakedTemplate
	values   []ParamValue
}

func (b *boundTemplate) isTagArg() {}

func (b *boundTemplate) Build(w *Writer) error {
	// Build value lookup map
	valueMap := make(map[string]Builder, len(b.values))
	for _, v := range b.values {
		valueMap[v.param.name] = v.value
	}

	// Write segments interleaved with parameter values
	for i, segment := range b.template.segments {
		if len(segment) > 0 {
			if _, err := w.w.Write(segment); err != nil {
				return err
			}
		}
		if i < len(b.template.params) {
			param := b.template.params[i]
			if value, ok := valueMap[param.name]; ok && value != nil {
				if err := value.Build(w); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
