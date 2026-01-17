package h

import (
	"bytes"
	"io"
	"slices"
	"sync"
)

// bufPool pools bytes.Buffer objects for Compile operations.
var bufPool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

// compiledBuilder stores pre-rendered HTML bytes for fast rendering.
type compiledBuilder struct {
	html []byte
}

func (b *compiledBuilder) isTagArg() {}

func (b *compiledBuilder) Build(w *Writer) error {
	_, err := w.w.Write(b.html)
	return err
}

// Compile pre-renders a Builder to bytes for faster subsequent renders.
// The resulting Builder writes the cached bytes directly without
// tree traversal, achieving performance comparable to html/template.
//
// Use Compile for static content that is rendered frequently:
//
//	// Compile once at startup
//	header, err := h.Compile(h.Header(h.H1(h.Text("My Site"))))
//	if err != nil {
//		// handle error
//	}
//
//	// Fast renders afterward
//	h.Render(w, header)
//
// Note: Compiled builders ignore indentation settings since the output
// is pre-computed. For pretty-printed cached output, render with
// RenderIndent first, then compile the result using Raw.
func Compile(b Builder) (Builder, error) {
	if b == nil {
		return nil, nil
	}
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	if err := Render(buf, b); err != nil {
		bufPool.Put(buf)
		return nil, err
	}
	result := &compiledBuilder{html: slices.Clone(buf.Bytes())}
	bufPool.Put(buf)
	return result, nil
}

// MustCompile is like Compile but panics if compilation fails.
// It is intended for use during initialization when errors should be caught immediately.
//
//	// Compile once at startup
//	header := h.MustCompile(h.Header(h.H1(h.Text("My Site"))))
func MustCompile(b Builder) Builder {
	compiled, err := Compile(b)
	if err != nil {
		panic("htmlgen: Compile failed: " + err.Error())
	}
	return compiled
}

// Param is a placeholder for dynamic content in a parameterized template.
// Use with CompileParams to create templates with variable content.
type Param struct {
	name string
}

func (p *Param) isTagArg() {}

// Build signals to CompileParams that a parameter slot was encountered.
// When used outside CompileParams, it renders nothing.
func (p *Param) Build(w *Writer) error {
	if cw, ok := w.w.(*compileWriter); ok {
		cw.recordParam(p)
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

// NewParam creates a named placeholder for dynamic content in CompileParams.
//
//	title := h.NewParam("title")
//	tmpl, err := h.CompileParams(h.Div(h.H1(title)))
//	if err != nil {
//		// handle error
//	}
//	tmpl.Render(w, title.Value(h.Text("Hello")))
func NewParam(name string) *Param {
	return &Param{name: name}
}

// CompiledTemplate is a pre-computed template with parameter placeholders.
// Created by CompileParams, it stores static HTML segments and parameter positions.
type CompiledTemplate struct {
	segments [][]byte
	params   []*Param
}

// compileWriter captures segments between parameters during CompileParams.
type compileWriter struct {
	buf      bytes.Buffer
	segments [][]byte
	params   []*Param
}

func (cw *compileWriter) Write(p []byte) (int, error) {
	return cw.buf.Write(p)
}

func (cw *compileWriter) WriteString(s string) (int, error) {
	return cw.buf.WriteString(s)
}

func (cw *compileWriter) recordParam(p *Param) {
	// Copy current buffer contents to a new segment
	segment := make([]byte, cw.buf.Len())
	copy(segment, cw.buf.Bytes())
	cw.segments = append(cw.segments, segment)
	cw.buf.Reset()
	cw.params = append(cw.params, p)
}

// CompileParams pre-computes a Builder with Param placeholders for dynamic content.
// Static HTML is pre-rendered into segments, with params marking where dynamic
// content will be inserted at render time.
//
//	title := h.NewParam("title")
//	content := h.NewParam("content")
//
//	tmpl, err := h.CompileParams(h.Div(
//	    h.H1(title),
//	    h.P(content),
//	))
//	if err != nil {
//		// handle error
//	}
//
//	// Render with values
//	tmpl.Render(w,
//	    title.Value(h.Text("Hello World")),
//	    content.Value(h.Text("Welcome!")),
//	)
func CompileParams(b Builder) (*CompiledTemplate, error) {
	if b == nil {
		return &CompiledTemplate{}, nil
	}
	cw := &compileWriter{}
	w := &Writer{w: cw, openTags: make([]string, 0, 32)}
	if err := b.Build(w); err != nil {
		return nil, err
	}

	// Add final segment
	segment := make([]byte, cw.buf.Len())
	copy(segment, cw.buf.Bytes())
	cw.segments = append(cw.segments, segment)

	return &CompiledTemplate{
		segments: cw.segments,
		params:   cw.params,
	}, nil
}

// MustCompileParams is like CompileParams but panics if compilation fails.
// It is intended for use during initialization when errors should be caught immediately.
//
//	title := h.NewParam("title")
//	content := h.NewParam("content")
//
//	tmpl := h.MustCompileParams(h.Div(
//	    h.H1(title),
//	    h.P(content),
//	))
func MustCompileParams(b Builder) *CompiledTemplate {
	tmpl, err := CompileParams(b)
	if err != nil {
		panic("htmlgen: CompileParams failed: " + err.Error())
	}
	return tmpl
}

// With creates a Builder with parameter values bound for rendering.
func (t *CompiledTemplate) With(values ...ParamValue) Builder {
	return &boundTemplate{template: t, values: values}
}

// Render writes the template directly with parameter values.
func (t *CompiledTemplate) Render(w io.Writer, values ...ParamValue) error {
	return Render(w, t.With(values...))
}

// boundTemplate is a CompiledTemplate with parameter values bound.
type boundTemplate struct {
	template *CompiledTemplate
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
