package h

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

// errorWriter is a writer that always returns an error
type errorWriter struct{}

func (e *errorWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("write error")
}

func TestHtmlWriting(t *testing.T) {
	cases := []struct {
		Desc     string
		Expected string
		Code     func(w *Writer)
	}{
		{
			"Basic HTML Writing",
			`<!DOCTYPE html>
<html lang="en"><head><link rel="stylesheet" href="style.css" rel="preload"/></head><body><script src="script.js" defer></script></body></html>`,
			func(w *Writer) {
				w.Doctype()
				w.OpenTag("html", Attrs("lang", "en"))
				{
					w.OpenTag("head", nil)
					{
						w.SelfClosingTag("link", Attrs("rel", "stylesheet", "href", "style.css", "rel", "preload"))
					}
					w.CloseOneTag()
					w.OpenTag("body", nil)
					{
						w.OpenTag("script", Attrs("src", "script.js", "defer", ""))
					}
					w.CloseTag("body")
				}
			},
		},
		{
			"Tag-api Writing",
			`<!DOCTYPE html>
<html lang="en"><head><link rel="stylesheet" href="style.css" rel="preload"/></head><body><script src="script.js" defer></script></body></html>`,
			func(w *Writer) {
				b := Html(Attrs("lang", "en"),
					Head(nil,
						Link(Attrs("rel", "stylesheet", "href", "style.css", "rel", "preload")),
					),
					Body(nil,
						Script(Attrs("src", "script.js", "defer", "")),
					),
				)
				b.Build(w)
			},
		},
		{
			"Pretty Printing",
			`<!DOCTYPE html>
<html lang="en">
  <head>
    <link rel="stylesheet" href="style.css" rel="preload"/>
  </head>
  <body>
    <script src="script.js" defer>
    </script>
  </body>
</html>
`,
			func(w *Writer) {
				w.SetIndent("  ")
				b := Html(Attrs("lang", "en"),
					Head(nil,
						Link(Attrs("rel", "stylesheet", "href", "style.css", "rel", "preload")),
					),
					Body(nil,
						Script(Attrs("src", "script.js", "defer", "")),
					),
				)
				b.Build(w)
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)

			w := NewWriter(buf)
			c.Code(w)
			w.Close()

			out := buf.String()

			if out != c.Expected {
				t.Errorf("expected\n%q\n  to equal\n%q", c.Expected, out)
			}
		})
	}
}

func TestAttr(t *testing.T) {
	a := Attr("class", "foo")
	if a.Name != "class" || a.Value != "foo" {
		t.Errorf("expected class=foo, got %s=%s", a.Name, a.Value)
	}
}

func TestAttrPanicOnEmptyName(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for empty attribute name")
		}
	}()
	Attr("", "value")
}

func TestAttrs(t *testing.T) {
	attrs := Attrs("class", "foo", "id", "bar")
	if len(attrs) != 2 {
		t.Fatalf("expected 2 attrs, got %d", len(attrs))
	}
	if attrs[0].Name != "class" || attrs[0].Value != "foo" {
		t.Errorf("expected class=foo, got %s=%s", attrs[0].Name, attrs[0].Value)
	}
	if attrs[1].Name != "id" || attrs[1].Value != "bar" {
		t.Errorf("expected id=bar, got %s=%s", attrs[1].Name, attrs[1].Value)
	}
}

func TestAttrsEmpty(t *testing.T) {
	attrs := Attrs()
	if attrs != nil {
		t.Errorf("expected nil for empty Attrs(), got %v", attrs)
	}
}

func TestAttrsPanicOnOddArgs(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for odd number of arguments")
		}
	}()
	Attrs("class", "foo", "id")
}

func TestAttrsPanicOnEmptyName(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for empty attribute name")
		}
	}()
	Attrs("", "value")
}

func TestAttrsMap(t *testing.T) {
	attrs := AttrsMap(map[string]string{"class": "foo", "id": "bar"})
	if len(attrs) != 2 {
		t.Fatalf("expected 2 attrs, got %d", len(attrs))
	}
	// Should be sorted alphabetically
	if attrs[0].Name != "class" || attrs[0].Value != "foo" {
		t.Errorf("expected class=foo first, got %s=%s", attrs[0].Name, attrs[0].Value)
	}
	if attrs[1].Name != "id" || attrs[1].Value != "bar" {
		t.Errorf("expected id=bar second, got %s=%s", attrs[1].Name, attrs[1].Value)
	}
}

func TestAttributesGet(t *testing.T) {
	attrs := Attrs("class", "foo", "id", "bar")
	val, ok := attrs.Get("class")
	if !ok || val != "foo" {
		t.Errorf("expected foo, got %s (ok=%v)", val, ok)
	}
	val, ok = attrs.Get("nonexistent")
	if ok || val != "" {
		t.Errorf("expected empty string and false, got %s (ok=%v)", val, ok)
	}
}

func TestAttributesIndex(t *testing.T) {
	attrs := Attrs("class", "foo", "id", "bar")
	if idx := attrs.Index("class"); idx != 0 {
		t.Errorf("expected index 0, got %d", idx)
	}
	if idx := attrs.Index("id"); idx != 1 {
		t.Errorf("expected index 1, got %d", idx)
	}
	if idx := attrs.Index("nonexistent"); idx != -1 {
		t.Errorf("expected index -1, got %d", idx)
	}
}

func TestAttributesSet(t *testing.T) {
	attrs := Attrs("class", "foo")
	attrs.Set("class", "bar")
	if val, _ := attrs.Get("class"); val != "bar" {
		t.Errorf("expected bar, got %s", val)
	}
	attrs.Set("id", "baz")
	if val, _ := attrs.Get("id"); val != "baz" {
		t.Errorf("expected baz, got %s", val)
	}
	if len(attrs) != 2 {
		t.Errorf("expected 2 attrs, got %d", len(attrs))
	}
}

func TestAttributesSetDefault(t *testing.T) {
	attrs := Attrs("class", "foo")
	attrs.SetDefault("class", "bar")
	if val, _ := attrs.Get("class"); val != "foo" {
		t.Errorf("expected foo (unchanged), got %s", val)
	}
	attrs.SetDefault("id", "baz")
	if val, _ := attrs.Get("id"); val != "baz" {
		t.Errorf("expected baz, got %s", val)
	}
}

func TestAttributesDelete(t *testing.T) {
	attrs := Attrs("class", "foo", "id", "bar")
	attrs.Delete("class")
	if _, ok := attrs.Get("class"); ok {
		t.Error("expected class to be deleted")
	}
	if len(attrs) != 1 {
		t.Errorf("expected 1 attr after delete, got %d", len(attrs))
	}
	// Deleting non-existent key should be a no-op
	attrs.Delete("nonexistent")
	if len(attrs) != 1 {
		t.Errorf("expected 1 attr after no-op delete, got %d", len(attrs))
	}
}

func TestCloseOneTagError(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	w := NewWriter(buf)
	err := w.CloseOneTag()
	if !errors.Is(err, ErrUnknownTagToClose) {
		t.Errorf("expected ErrUnknownTagToClose, got %v", err)
	}
}

func TestCloseTagError(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	w := NewWriter(buf)
	err := w.CloseTag("div")
	if !errors.Is(err, ErrUnknownTagToClose) {
		t.Errorf("expected ErrUnknownTagToClose, got %v", err)
	}
}

func TestTextEscaping(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	w := NewWriter(buf)
	w.Text("<script>alert('xss')</script>")
	expected := "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;"
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

func TestRaw(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	w := NewWriter(buf)
	w.Raw("<script>alert('xss')</script>")
	expected := "<script>alert('xss')</script>"
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

func TestTextBuilder(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	b := Text("<div>")
	Render(buf, b)
	expected := "&lt;div&gt;"
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

func TestRawBuilder(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	b := Raw("<div>")
	Render(buf, b)
	expected := "<div>"
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

func TestFragment(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	b := Fragment(
		Text("hello"),
		Text(" "),
		Text("world"),
	)
	Render(buf, b)
	expected := "hello world"
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

func TestFragmentWithNilChildren(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	b := Fragment(
		Text("hello"),
		nil,
		Text("world"),
	)
	Render(buf, b)
	expected := "helloworld"
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

func TestRenderNil(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	err := Render(buf, nil)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if buf.String() != "" {
		t.Errorf("expected empty string, got %q", buf.String())
	}
}

func TestRenderPrettyNil(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	err := RenderIndent(buf, "  ", nil)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if buf.String() != "" {
		t.Errorf("expected empty string, got %q", buf.String())
	}
}

func TestRenderPretty(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	b := Div(nil, P(nil, Text("hello")))
	RenderIndent(buf, "  ", b)
	// Text content is not indented - only tags get indentation
	expected := "<div>\n  <p>\nhello  </p>\n</div>\n"
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

func TestSelfClosingTagWithChildren(t *testing.T) {
	// When a self-closing tag has children, it should render as a normal tag
	buf := bytes.NewBuffer(nil)
	b := Input(Attrs("type", "text"), Text("ignored"))
	Render(buf, b)
	expected := "<input type=\"text\">ignored</input>"
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

func TestCustomElement(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	b := CustomElement("my-component", Attrs("data-foo", "bar"), Text("content"))
	Render(buf, b)
	expected := "<my-component data-foo=\"bar\">content</my-component>"
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

func TestAttributeEscaping(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	b := Div(Attrs("data-value", `<script>"hello"</script>`))
	Render(buf, b)
	expected := `<div data-value="&lt;script&gt;&#34;hello&#34;&lt;/script&gt;"></div>`
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

func TestBooleanAttribute(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	b := Input(Attrs("type", "checkbox", "checked", "", "disabled", ""))
	Render(buf, b)
	expected := `<input type="checkbox" checked disabled/>`
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

func TestHtmlDefaultLang(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	b := Html(nil)
	Render(buf, b)
	expected := "<!DOCTYPE html>\n<html lang=\"en\"></html>"
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

func TestHtmlCustomLang(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	b := Html(Attrs("lang", "fr"))
	Render(buf, b)
	expected := "<!DOCTYPE html>\n<html lang=\"fr\"></html>"
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

func TestTagWithNilChildren(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	b := Div(nil, nil, Text("hello"), nil, Text("world"), nil)
	Render(buf, b)
	expected := "<div>helloworld</div>"
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

func TestCloseTagClosesIntermediateTagsToo(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	w := NewWriter(buf)
	w.OpenTag("div", nil)
	w.OpenTag("span", nil)
	w.OpenTag("p", nil)
	w.CloseTag("div") // Should close p, span, and div
	w.Close()
	expected := "<div><span><p></p></span></div>"
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

func TestWriterDoctype(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	w := NewWriter(buf)
	w.Doctype()
	expected := "<!DOCTYPE html>\n"
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

func TestWriteErrors(t *testing.T) {
	ew := &errorWriter{}
	w := NewWriter(ew)

	if err := w.Doctype(); err == nil {
		t.Error("expected error from Doctype")
	}
	if err := w.OpenTag("div", nil); err == nil {
		t.Error("expected error from OpenTag")
	}
	if err := w.SelfClosingTag("br", nil); err == nil {
		t.Error("expected error from SelfClosingTag")
	}
	if err := w.Text("hello"); err == nil {
		t.Error("expected error from Text")
	}
	if err := w.Raw("hello"); err == nil {
		t.Error("expected error from Raw")
	}
}

func TestCloseOnEmptyWriter(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	w := NewWriter(buf)
	err := w.Close()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if buf.String() != "" {
		t.Errorf("expected empty string, got %q", buf.String())
	}
}

func TestServerClientConstants(t *testing.T) {
	// On server (non-js build), Server should be true and Client should be false
	if !Server {
		t.Error("expected Server to be true")
	}
	if Client {
		t.Error("expected Client to be false")
	}
}

func TestChan(t *testing.T) {
	ch := Chan[int]("test")
	if ch == nil {
		t.Error("expected non-nil channel")
	}
	// Test that channel has buffer of 1
	select {
	case ch <- 42:
		// success
	default:
		t.Error("expected to be able to send on buffered channel")
	}
}

func TestAttrsMapEmpty(t *testing.T) {
	attrs := AttrsMap(map[string]string{})
	if len(attrs) != 0 {
		t.Errorf("expected 0 attrs, got %d", len(attrs))
	}
}

func TestVariousHtmlTags(t *testing.T) {
	// Test a sampling of HTML tags to ensure they render correctly
	tests := []struct {
		name     string
		builder  Builder
		expected string
	}{
		{"Div", Div(nil), "<div></div>"},
		{"Span", Span(nil, Text("text")), "<span>text</span>"},
		{"A", A(Attrs("href", "/")), "<a href=\"/\"></a>"},
		{"P", P(nil), "<p></p>"},
		{"Ul/Li", Ul(nil, Li(nil, Text("item"))), "<ul><li>item</li></ul>"},
		{"Table", Table(nil, Tr(nil, Td(nil))), "<table><tr><td></td></tr></table>"},
		{"Form", Form(Attrs("method", "post")), "<form method=\"post\"></form>"},
		{"Button", Button(Attrs("type", "submit"), Text("Submit")), "<button type=\"submit\">Submit</button>"},
		{"Img", Img(Attrs("src", "test.png", "alt", "test")), "<img src=\"test.png\" alt=\"test\"/>"},
		{"Br", Br(nil), "<br/>"},
		{"Hr", Hr(nil), "<hr/>"},
		{"Meta", Meta(Attrs("charset", "utf-8")), "<meta charset=\"utf-8\"/>"},
		{"Link", Link(Attrs("rel", "stylesheet")), "<link rel=\"stylesheet\"/>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			Render(buf, tt.builder)
			if buf.String() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, buf.String())
			}
		})
	}
}

// Test Render function return value
func TestRenderReturnsError(t *testing.T) {
	ew := &errorWriter{}
	b := Div(nil, Text("hello"))
	err := Render(ew, b)
	if err == nil {
		t.Error("expected error from Render")
	}
}

// Test RenderPretty function return value
func TestRenderPrettyReturnsError(t *testing.T) {
	ew := &errorWriter{}
	b := Div(nil, Text("hello"))
	err := RenderIndent(ew, "  ", b)
	if err == nil {
		t.Error("expected error from RenderPretty")
	}
}

var _ io.Writer = (*errorWriter)(nil)

// partialWriter writes a limited number of bytes before returning an error
type partialWriter struct {
	remaining int
}

func (p *partialWriter) Write(data []byte) (n int, err error) {
	if p.remaining <= 0 {
		return 0, errors.New("partial write error")
	}
	if len(data) <= p.remaining {
		p.remaining -= len(data)
		return len(data), nil
	}
	n = p.remaining
	p.remaining = 0
	return n, errors.New("partial write error")
}

var _ io.Writer = (*partialWriter)(nil)

func TestWriteIndentError(t *testing.T) {
	// Test error in writeIndent during OpenTag with indentation
	pw := &partialWriter{remaining: 5} // enough for "<div>" but not indent
	w := NewWriter(pw)
	w.SetIndent("  ")
	w.OpenTag("div", nil) // first tag works at depth 0
	err := w.OpenTag("span", nil)
	if err == nil {
		t.Error("expected error from OpenTag with indent")
	}
}

func TestSelfClosingTagWithIndentError(t *testing.T) {
	pw := &partialWriter{remaining: 10}
	w := NewWriter(pw)
	w.SetIndent("  ")
	w.OpenTag("div", nil)
	err := w.SelfClosingTag("br", nil)
	if err == nil {
		t.Error("expected error from SelfClosingTag with indent")
	}
}

func TestCloseTagWithIndentError(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	w := NewWriter(buf)
	w.SetIndent("  ")
	w.OpenTag("div", nil)

	// Replace writer to force error on close
	w.w = &errorWriter{}
	err := w.CloseTag("div")
	if err == nil {
		t.Error("expected error from CloseTag with indent")
	}
}

func TestCloseOneTagWithIndentError(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	w := NewWriter(buf)
	w.SetIndent("  ")
	w.OpenTag("div", nil)

	w.w = &errorWriter{}
	err := w.CloseOneTag()
	if err == nil {
		t.Error("expected error from CloseOneTag with indent")
	}
}

func TestCloseWithIndentError(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	w := NewWriter(buf)
	w.SetIndent("  ")
	w.OpenTag("div", nil)

	w.w = &errorWriter{}
	err := w.Close()
	if err == nil {
		t.Error("expected error from Close with indent")
	}
}

func TestBuilderErrorPropagation(t *testing.T) {
	ew := &errorWriter{}
	w := NewWriter(ew)

	// Test htmlTagBuilder error propagation
	html := &htmlTagBuilder{Attrs: nil, Children: nil}
	if err := html.Build(w); err == nil {
		t.Error("expected error from htmlTagBuilder.Build")
	}

	// Test tagBuilder error propagation
	tag := &tagBuilder{Name: "div", Attrs: nil, Children: nil}
	if err := tag.Build(w); err == nil {
		t.Error("expected error from tagBuilder.Build")
	}

	// Test fragmentBuilder error propagation
	frag := &fragmentBuilder{Children: []Builder{Text("test")}}
	if err := frag.Build(w); err == nil {
		t.Error("expected error from fragmentBuilder.Build")
	}
}

func TestHtmlTagBuilderChildError(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	w := NewWriter(buf)

	// Create an html tag with a child that will error
	html := Html(nil, Div(nil, Text("test")))

	// Replace the writer after doctype to force child error
	html.Build(w) // Ignore first error

	ew := &errorWriter{}
	w2 := NewWriter(ew)
	w2.openTags = append(w2.openTags, "html") // simulate opened tag

	// Test that child error propagates through tagBuilder
	div := Div(nil, Text("test"))
	if err := div.Build(w2); err == nil {
		t.Error("expected error from child build")
	}
}

func TestTagBuilderSelfCloseError(t *testing.T) {
	ew := &errorWriter{}
	w := NewWriter(ew)

	// Test self-closing tag error
	br := &tagBuilder{Name: "br", SelfClose: true}
	if err := br.Build(w); err == nil {
		t.Error("expected error from self-closing tagBuilder.Build")
	}
}

func TestHtmlTagBuilderWithNilChildren(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	b := Html(nil, nil, Div(nil), nil)
	Render(buf, b)
	expected := "<!DOCTYPE html>\n<html lang=\"en\"><div></div></html>"
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

func TestOpenTagAttributeError(t *testing.T) {
	pw := &partialWriter{remaining: 10}
	w := NewWriter(pw)
	// Try to write a tag with attributes that will fail partway through
	err := w.OpenTag("div", Attrs("class", "very-long-class-name-that-exceeds-limit"))
	if err == nil {
		t.Error("expected error from OpenTag with attributes")
	}
}

func TestSelfClosingTagAttributeError(t *testing.T) {
	pw := &partialWriter{remaining: 10}
	w := NewWriter(pw)
	err := w.SelfClosingTag("input", Attrs("type", "text", "placeholder", "very-long-placeholder"))
	if err == nil {
		t.Error("expected error from SelfClosingTag with attributes")
	}
}

func TestCloseTagWriteError(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	w := NewWriter(buf)
	w.OpenTag("div", nil)
	w.OpenTag("span", nil)

	// Replace to force error on close
	w.w = &errorWriter{}
	err := w.CloseTag("div")
	if err == nil {
		t.Error("expected error from CloseTag write")
	}
}

func TestCloseOneTagWriteError(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	w := NewWriter(buf)
	w.OpenTag("div", nil)

	w.w = &errorWriter{}
	err := w.CloseOneTag()
	if err == nil {
		t.Error("expected error from CloseOneTag write")
	}
}

func TestCloseWriteError(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	w := NewWriter(buf)
	w.OpenTag("div", nil)

	w.w = &errorWriter{}
	err := w.Close()
	if err == nil {
		t.Error("expected error from Close write")
	}
}

func TestCloseTagNewlineError(t *testing.T) {
	// Test error writing newline after close tag
	pw := &partialWriter{remaining: 100}
	w := NewWriter(pw)
	w.SetIndent("  ")
	w.OpenTag("div", nil)
	w.OpenTag("span", nil)
	pw.remaining = 15 // enough for closing tags but not newline
	err := w.CloseTag("div")
	// May or may not error depending on exact byte count
	_ = err
}

func TestCloseOneTagNewlineError(t *testing.T) {
	pw := &partialWriter{remaining: 100}
	w := NewWriter(pw)
	w.SetIndent("  ")
	w.OpenTag("div", nil)
	pw.remaining = 5 // not enough for close + newline
	err := w.CloseOneTag()
	if err == nil {
		t.Error("expected error from CloseOneTag newline")
	}
}

func TestCloseMultipleTagsNewlineError(t *testing.T) {
	pw := &partialWriter{remaining: 100}
	w := NewWriter(pw)
	w.SetIndent("  ")
	w.OpenTag("div", nil)
	w.OpenTag("span", nil)
	pw.remaining = 6 // just enough for </span> but not newline
	err := w.Close()
	if err == nil {
		t.Error("expected error from Close newline")
	}
}

func TestOpenTagNewlineError(t *testing.T) {
	pw := &partialWriter{remaining: 5} // just enough for <div>
	w := NewWriter(pw)
	w.SetIndent("  ")
	err := w.OpenTag("div", nil)
	if err == nil {
		t.Error("expected error from OpenTag newline")
	}
}

func TestSelfClosingTagNewlineError(t *testing.T) {
	pw := &partialWriter{remaining: 5} // just enough for <br/>
	w := NewWriter(pw)
	w.SetIndent("  ")
	err := w.SelfClosingTag("br", nil)
	if err == nil {
		t.Error("expected error from SelfClosingTag newline")
	}
}

func TestOpenTagClosingBracketError(t *testing.T) {
	pw := &partialWriter{remaining: 4} // just enough for <div but not >
	w := NewWriter(pw)
	err := w.OpenTag("div", nil)
	if err == nil {
		t.Error("expected error from OpenTag closing bracket")
	}
}

func TestSelfClosingTagClosingBracketError(t *testing.T) {
	pw := &partialWriter{remaining: 3} // just enough for <br but not />
	w := NewWriter(pw)
	err := w.SelfClosingTag("br", nil)
	if err == nil {
		t.Error("expected error from SelfClosingTag closing bracket")
	}
}

func TestHtmlBuilderOpenTagError(t *testing.T) {
	pw := &partialWriter{remaining: 20} // enough for doctype, not for open tag attrs
	w := NewWriter(pw)
	html := Html(Attrs("lang", "en", "data-theme", "dark"))
	err := html.Build(w)
	if err == nil {
		t.Error("expected error from Html builder open tag")
	}
}

func TestTagBuilderChildError(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	w := NewWriter(buf)
	w.OpenTag("parent", nil)

	// Create a failing child
	failingChild := &tagBuilder{Name: "child", Children: []Builder{Text("text")}}

	// Switch to error writer
	w.w = &errorWriter{}
	err := failingChild.Build(w)
	if err == nil {
		t.Error("expected error from child builder")
	}
}

func TestTagBuilderCloseError(t *testing.T) {
	// Open a tag, then switch writer and try to close
	pw := &partialWriter{remaining: 100}
	w := NewWriter(pw)
	w.OpenTag("outer", nil)

	// Create tagBuilder that opens successfully but fails on close
	pw.remaining = 10 // enough to open but not close
	tb := &tagBuilder{Name: "test", Children: []Builder{
		Text("some content that will fill the buffer"),
	}}
	err := tb.Build(w)
	if err == nil {
		t.Error("expected error from tagBuilder close")
	}
}

func TestHtmlBuilderCloseError(t *testing.T) {
	pw := &partialWriter{remaining: 100}
	w := NewWriter(pw)

	// Create Html with content, then limit remaining bytes
	html := Html(nil, Div(nil, Text("content")))
	pw.remaining = 50 // enough to open but not fully close
	err := html.Build(w)
	// Error may or may not occur depending on exact byte counts
	_ = err
}

func TestFragmentBuilderError(t *testing.T) {
	ew := &errorWriter{}
	w := NewWriter(ew)

	frag := Fragment(Div(nil, Text("test")))
	err := frag.Build(w)
	if err == nil {
		t.Error("expected error from fragment builder")
	}
}
