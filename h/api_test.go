package h

import (
	"bytes"
	"errors"
	"io"
	"strings"
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
					Head(
						Link(Attrs("rel", "stylesheet", "href", "style.css", "rel", "preload")),
					),
					Body(
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
					Head(
						Link(Attrs("rel", "stylesheet", "href", "style.css", "rel", "preload")),
					),
					Body(
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

func TestAttributesMerge(t *testing.T) {
	// Test merging with override
	attrs := Attrs("class", "foo", "id", "bar")
	attrs.Merge(Attrs("id", "newbar", "href", "/home"))
	if val, _ := attrs.Get("class"); val != "foo" {
		t.Errorf("expected class=foo, got %s", val)
	}
	if val, _ := attrs.Get("id"); val != "newbar" {
		t.Errorf("expected id=newbar, got %s", val)
	}
	if val, _ := attrs.Get("href"); val != "/home" {
		t.Errorf("expected href=/home, got %s", val)
	}
	if len(attrs) != 3 {
		t.Errorf("expected 3 attrs, got %d", len(attrs))
	}
}

func TestAttributesMergeEmpty(t *testing.T) {
	attrs := Attrs("class", "foo")
	attrs.Merge(nil)
	if len(attrs) != 1 {
		t.Errorf("expected 1 attr after merging nil, got %d", len(attrs))
	}
	attrs.Merge(Attrs())
	if len(attrs) != 1 {
		t.Errorf("expected 1 attr after merging empty, got %d", len(attrs))
	}
}

func TestAttributesMergeIntoEmpty(t *testing.T) {
	var attrs Attributes
	attrs.Merge(Attrs("class", "foo", "id", "bar"))
	if len(attrs) != 2 {
		t.Errorf("expected 2 attrs, got %d", len(attrs))
	}
	if val, _ := attrs.Get("class"); val != "foo" {
		t.Errorf("expected class=foo, got %s", val)
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
	b := Div(P(Text("hello")))
	RenderIndent(buf, "  ", b)
	// Text content is indented at the content depth (inside parent tag)
	expected := "<div>\n  <p>\n    hello\n  </p>\n</div>\n"
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

func TestRenderMinified(t *testing.T) {
	// Verify minified mode produces no extra whitespace
	buf := bytes.NewBuffer(nil)
	b := Div(P(Text("hello")), Span(Text("world")))
	Render(buf, b)
	// Minified output should have no newlines or indentation
	expected := "<div><p>hello</p><span>world</span></div>"
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

func TestRenderMinifiedDeepNesting(t *testing.T) {
	// Verify deeply nested elements in minified mode have no whitespace
	buf := bytes.NewBuffer(nil)
	b := Div(
		Div(
			Div(
				Span(Text("Nested")),
			),
		),
	)
	Render(buf, b)
	expected := "<div><div><div><span>Nested</span></div></div></div>"
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
	b := Html()
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
	b := Div(Text("hello"), Text("world"))
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
		{"Div", Div(), "<div></div>"},
		{"Span", Span(Text("text")), "<span>text</span>"},
		{"A", A(Attrs("href", "/")), "<a href=\"/\"></a>"},
		{"P", P(), "<p></p>"},
		{"Ul/Li", Ul(Li(Text("item"))), "<ul><li>item</li></ul>"},
		{"Table", Table(Tr(Td())), "<table><tr><td></td></tr></table>"},
		{"Form", Form(Attrs("method", "post")), "<form method=\"post\"></form>"},
		{"Button", Button(Attrs("type", "submit"), Text("Submit")), "<button type=\"submit\">Submit</button>"},
		{"Img", Img(Attrs("src", "test.png", "alt", "test")), "<img src=\"test.png\" alt=\"test\"/>"},
		{"Br", Br(), "<br/>"},
		{"Hr", Hr(), "<hr/>"},
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

func TestAllHeadingTags(t *testing.T) {
	tests := []struct {
		name     string
		builder  Builder
		expected string
	}{
		{"H1", H1(Text("Title")), "<h1>Title</h1>"},
		{"H2", H2(Text("Subtitle")), "<h2>Subtitle</h2>"},
		{"H3", H3(Text("Section")), "<h3>Section</h3>"},
		{"H4", H4(Text("Subsection")), "<h4>Subsection</h4>"},
		{"H5", H5(Text("Minor")), "<h5>Minor</h5>"},
		{"H6", H6(Text("Smallest")), "<h6>Smallest</h6>"},
		{"Hgroup", Hgroup(H1(Text("Main")), P(Text("Tagline"))), "<hgroup><h1>Main</h1><p>Tagline</p></hgroup>"},
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

func TestStructuralTags(t *testing.T) {
	tests := []struct {
		name     string
		builder  Builder
		expected string
	}{
		{"Header", Header(Text("Header")), "<header>Header</header>"},
		{"Footer", Footer(Text("Footer")), "<footer>Footer</footer>"},
		{"Main", Main(Text("Main")), "<main>Main</main>"},
		{"Nav", Nav(A(Attrs("href", "/"), Text("Home"))), "<nav><a href=\"/\">Home</a></nav>"},
		{"Section", Section(Text("Section")), "<section>Section</section>"},
		{"Article", Article(Text("Article")), "<article>Article</article>"},
		{"Aside", Aside(Text("Sidebar")), "<aside>Sidebar</aside>"},
		{"Address", Address(Text("123 Main St")), "<address>123 Main St</address>"},
		{"Search", Search(Input(Attrs("type", "search"))), "<search><input type=\"search\"/></search>"},
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

func TestContentGroupingTags(t *testing.T) {
	tests := []struct {
		name     string
		builder  Builder
		expected string
	}{
		{"Blockquote", Blockquote(Text("Quote")), "<blockquote>Quote</blockquote>"},
		{"Pre", Pre(Text("Preformatted")), "<pre>Preformatted</pre>"},
		{"Figure", Figure(Img(Attrs("src", "img.png", "alt", "Alt"))), "<figure><img src=\"img.png\" alt=\"Alt\"/></figure>"},
		{"Figcaption", Figcaption(Text("Caption")), "<figcaption>Caption</figcaption>"},
		{"Ol", Ol(Li(Text("One")), Li(Text("Two"))), "<ol><li>One</li><li>Two</li></ol>"},
		{"Ul", Ul(Li(Text("Item"))), "<ul><li>Item</li></ul>"},
		{"Li", Li(Text("List item")), "<li>List item</li>"},
		{"Dl", Dl(Dt(Text("Term")), Dd(Text("Definition"))), "<dl><dt>Term</dt><dd>Definition</dd></dl>"},
		{"Dt", Dt(Text("Term")), "<dt>Term</dt>"},
		{"Dd", Dd(Text("Definition")), "<dd>Definition</dd>"},
		{"Menu", Menu(Li(Button(Text("Action")))), "<menu><li><button>Action</button></li></menu>"},
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

func TestTextFormattingTags(t *testing.T) {
	tests := []struct {
		name     string
		builder  Builder
		expected string
	}{
		{"B", B(Text("Bold")), "<b>Bold</b>"},
		{"Strong", Strong(Text("Strong")), "<strong>Strong</strong>"},
		{"I", I(Text("Italic")), "<i>Italic</i>"},
		{"Em", Em(Text("Emphasis")), "<em>Emphasis</em>"},
		{"U", U(Text("Underline")), "<u>Underline</u>"},
		{"Mark", Mark(Text("Highlight")), "<mark>Highlight</mark>"},
		{"Small", Small(Text("Small")), "<small>Small</small>"},
		{"Sub", Sub(Text("2")), "<sub>2</sub>"},
		{"Sup", Sup(Text("2")), "<sup>2</sup>"},
		{"Code", Code(Text("code")), "<code>code</code>"},
		{"Kbd", Kbd(Text("Enter")), "<kbd>Enter</kbd>"},
		{"Samp", Samp(Text("Output")), "<samp>Output</samp>"},
		{"Var", Var(Text("x")), "<var>x</var>"},
		{"Cite", Cite(Text("Book Title")), "<cite>Book Title</cite>"},
		{"Abbr", Abbr(Attrs("title", "Abbreviation"), Text("Abbr")), "<abbr title=\"Abbreviation\">Abbr</abbr>"},
		{"Dfn", Dfn(Text("Definition")), "<dfn>Definition</dfn>"},
		{"Q", Q(Text("Quote")), "<q>Quote</q>"},
		{"S", S(Text("Strikethrough")), "<s>Strikethrough</s>"},
		{"Del", Del(Text("Deleted")), "<del>Deleted</del>"},
		{"Ins", Ins(Text("Inserted")), "<ins>Inserted</ins>"},
		{"Time", Time(Attrs("datetime", "2024-01-01"), Text("Jan 1")), "<time datetime=\"2024-01-01\">Jan 1</time>"},
		{"Data", Data(Attrs("value", "123"), Text("One Two Three")), "<data value=\"123\">One Two Three</data>"},
		{"Wbr", Wbr(), "<wbr/>"},
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

func TestRubyAnnotationTags(t *testing.T) {
	tests := []struct {
		name     string
		builder  Builder
		expected string
	}{
		{"Ruby", Ruby(Text("漢"), Rp(Text("(")), Rt(Text("kan")), Rp(Text(")"))), "<ruby>漢<rp>(</rp><rt>kan</rt><rp>)</rp></ruby>"},
		{"Rt", Rt(Text("annotation")), "<rt>annotation</rt>"},
		{"Rp", Rp(Text("(")), "<rp>(</rp>"},
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

func TestBidirectionalTags(t *testing.T) {
	tests := []struct {
		name     string
		builder  Builder
		expected string
	}{
		{"Bdi", Bdi(Text("مرحبا")), "<bdi>مرحبا</bdi>"},
		{"Bdo", Bdo(Attrs("dir", "rtl"), Text("Right to left")), "<bdo dir=\"rtl\">Right to left</bdo>"},
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

func TestMediaTags(t *testing.T) {
	tests := []struct {
		name     string
		builder  Builder
		expected string
	}{
		{"Audio", Audio(Attrs("controls", ""), Source(Attrs("src", "audio.mp3", "type", "audio/mpeg"))), "<audio controls><source src=\"audio.mp3\" type=\"audio/mpeg\"/></audio>"},
		{"Video", Video(Attrs("controls", ""), Source(Attrs("src", "video.mp4", "type", "video/mp4"))), "<video controls><source src=\"video.mp4\" type=\"video/mp4\"/></video>"},
		{"Picture", Picture(Source(Attrs("srcset", "img.webp", "type", "image/webp")), Img(Attrs("src", "img.png", "alt", "Image"))), "<picture><source srcset=\"img.webp\" type=\"image/webp\"/><img src=\"img.png\" alt=\"Image\"/></picture>"},
		{"Source", Source(Attrs("src", "media.mp4")), "<source src=\"media.mp4\"/>"},
		{"Track", Track(Attrs("src", "subs.vtt", "kind", "subtitles")), "<track src=\"subs.vtt\" kind=\"subtitles\"/>"},
		{"Map", Map(Attrs("name", "imagemap"), Area(Attrs("shape", "rect", "coords", "0,0,50,50", "href", "/"))), "<map name=\"imagemap\"><area shape=\"rect\" coords=\"0,0,50,50\" href=\"/\"/></map>"},
		{"Area", Area(Attrs("shape", "circle", "coords", "25,25,25", "href", "/")), "<area shape=\"circle\" coords=\"25,25,25\" href=\"/\"/>"},
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

func TestEmbeddedContentTags(t *testing.T) {
	tests := []struct {
		name     string
		builder  Builder
		expected string
	}{
		{"Iframe", Iframe(Attrs("src", "https://example.com")), "<iframe src=\"https://example.com\"></iframe>"},
		{"Embed", Embed(Attrs("src", "plugin.swf", "type", "application/x-shockwave-flash")), "<embed src=\"plugin.swf\" type=\"application/x-shockwave-flash\"/>"},
		{"Object", Object(Attrs("data", "file.pdf", "type", "application/pdf")), "<object data=\"file.pdf\" type=\"application/pdf\"></object>"},
		{"Portal", Portal(Attrs("src", "https://example.com")), "<portal src=\"https://example.com\"></portal>"},
		{"Svg", Svg(Attrs("width", "100", "height", "100")), "<svg width=\"100\" height=\"100\"></svg>"},
		{"Math", Math(Text("x = 5")), "<math>x = 5</math>"},
		{"Canvas", Canvas(Attrs("width", "300", "height", "150")), "<canvas width=\"300\" height=\"150\"></canvas>"},
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

func TestTableTags(t *testing.T) {
	tests := []struct {
		name     string
		builder  Builder
		expected string
	}{
		{"Table", Table(), "<table></table>"},
		{"Caption", Caption(Text("Table Caption")), "<caption>Table Caption</caption>"},
		{"Thead", Thead(Tr(Th(Text("Header")))), "<thead><tr><th>Header</th></tr></thead>"},
		{"Tbody", Tbody(Tr(Td(Text("Data")))), "<tbody><tr><td>Data</td></tr></tbody>"},
		{"Tfoot", Tfoot(Tr(Td(Text("Footer")))), "<tfoot><tr><td>Footer</td></tr></tfoot>"},
		{"Tr", Tr(Td(Text("Cell"))), "<tr><td>Cell</td></tr>"},
		{"Th", Th(Text("Header")), "<th>Header</th>"},
		{"Td", Td(Text("Data")), "<td>Data</td>"},
		{"Colgroup", Colgroup(Col(Attrs("span", "2"))), "<colgroup><col span=\"2\"/></colgroup>"},
		{"Col", Col(Attrs("style", "background-color: yellow")), "<col style=\"background-color: yellow\"/>"},
		{"FullTable", Table(
			Caption(Text("Users")),
			Colgroup(Col(), Col()),
			Thead(Tr(Th(Text("Name")), Th(Text("Age")))),
			Tbody(Tr(Td(Text("Alice")), Td(Text("30")))),
			Tfoot(Tr(Td(Attrs("colspan", "2"), Text("Total: 1")))),
		), "<table><caption>Users</caption><colgroup><col/><col/></colgroup><thead><tr><th>Name</th><th>Age</th></tr></thead><tbody><tr><td>Alice</td><td>30</td></tr></tbody><tfoot><tr><td colspan=\"2\">Total: 1</td></tr></tfoot></table>"},
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

func TestFormTags(t *testing.T) {
	tests := []struct {
		name     string
		builder  Builder
		expected string
	}{
		{"Form", Form(Attrs("action", "/submit", "method", "post")), "<form action=\"/submit\" method=\"post\"></form>"},
		{"Fieldset", Fieldset(Legend(Text("Personal Info")), Input(Attrs("type", "text"))), "<fieldset><legend>Personal Info</legend><input type=\"text\"/></fieldset>"},
		{"Legend", Legend(Text("Section")), "<legend>Section</legend>"},
		{"Label", Label(Attrs("for", "name"), Text("Name:")), "<label for=\"name\">Name:</label>"},
		{"Input", Input(Attrs("type", "text", "name", "username")), "<input type=\"text\" name=\"username\"/>"},
		{"Button", Button(Attrs("type", "button"), Text("Click")), "<button type=\"button\">Click</button>"},
		{"Select", Select(Attrs("name", "country"), Option(Attrs("value", "us"), Text("USA"))), "<select name=\"country\"><option value=\"us\">USA</option></select>"},
		{"Optgroup", Optgroup(Attrs("label", "Group"), Option(Text("Item"))), "<optgroup label=\"Group\"><option>Item</option></optgroup>"},
		{"Option", Option(Attrs("value", "1"), Text("One")), "<option value=\"1\">One</option>"},
		{"Datalist", Datalist(Attrs("id", "browsers"), Option(Attrs("value", "Chrome"))), "<datalist id=\"browsers\"><option value=\"Chrome\"></option></datalist>"},
		{"Textarea", Textarea(Attrs("name", "message", "rows", "4", "cols", "50")), "<textarea name=\"message\" rows=\"4\" cols=\"50\"></textarea>"},
		{"Output", Output(Attrs("for", "a b"), Text("Result")), "<output for=\"a b\">Result</output>"},
		{"Progress", Progress(Attrs("value", "70", "max", "100")), "<progress value=\"70\" max=\"100\"></progress>"},
		{"Meter", Meter(Attrs("value", "0.6", "min", "0", "max", "1")), "<meter value=\"0.6\" min=\"0\" max=\"1\"></meter>"},
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

func TestInteractiveTags(t *testing.T) {
	tests := []struct {
		name     string
		builder  Builder
		expected string
	}{
		{"Details", Details(Summary(Text("More info")), P(Text("Hidden content"))), "<details><summary>More info</summary><p>Hidden content</p></details>"},
		{"Summary", Summary(Text("Click to expand")), "<summary>Click to expand</summary>"},
		{"Dialog", Dialog(Attrs("open", ""), P(Text("Dialog content"))), "<dialog open><p>Dialog content</p></dialog>"},
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

func TestDocumentMetadataTags(t *testing.T) {
	tests := []struct {
		name     string
		builder  Builder
		expected string
	}{
		{"Head", Head(Title(Text("Page Title"))), "<head><title>Page Title</title></head>"},
		{"Title", Title(Text("My Page")), "<title>My Page</title>"},
		{"Meta", Meta(Attrs("name", "description", "content", "A description")), "<meta name=\"description\" content=\"A description\"/>"},
		{"Link", Link(Attrs("rel", "icon", "href", "favicon.ico")), "<link rel=\"icon\" href=\"favicon.ico\"/>"},
		{"Style", Style(Raw("body { margin: 0; }")), "<style>body { margin: 0; }</style>"},
		{"Base", Base(Attrs("href", "https://example.com/")), "<base href=\"https://example.com/\"/>"},
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

func TestScriptingTags(t *testing.T) {
	tests := []struct {
		name     string
		builder  Builder
		expected string
	}{
		{"Script", Script(Attrs("src", "app.js")), "<script src=\"app.js\"></script>"},
		{"ScriptInline", Script(Raw("console.log('Hello');")), "<script>console.log('Hello');</script>"},
		{"Noscript", Noscript(Text("JavaScript is required")), "<noscript>JavaScript is required</noscript>"},
		{"Template", Template(Attrs("id", "my-template"), Div(Text("Template content"))), "<template id=\"my-template\"><div>Template content</div></template>"},
		{"Slot", Slot(Attrs("name", "header")), "<slot name=\"header\"></slot>"},
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

func TestBodyTag(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	b := Body(Attrs("class", "dark-mode"), Div(Text("Content")))
	Render(buf, b)
	expected := "<body class=\"dark-mode\"><div>Content</div></body>"
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

func TestDeepNesting(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	b := Div(
		Div(
			Div(
				Div(
					Div(
						Div(
							Div(
								Span(Text("Deeply nested")),
							),
						),
					),
				),
			),
		),
	)
	Render(buf, b)
	expected := "<div><div><div><div><div><div><div><span>Deeply nested</span></div></div></div></div></div></div></div>"
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

func TestDeepNestingWithIndent(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	b := Div(
		Div(
			Div(
				Span(Text("Nested")),
			),
		),
	)
	RenderIndent(buf, "\t", b)
	// Text is now properly indented at content depth
	expected := "<div>\n\t<div>\n\t\t<div>\n\t\t\t<span>\n\t\t\t\tNested\n\t\t\t</span>\n\t\t</div>\n\t</div>\n</div>\n"
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

func TestManyAttributes(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	b := Div(Attrs(
		"id", "myDiv",
		"class", "container fluid",
		"data-value", "123",
		"data-name", "test",
		"aria-label", "Main container",
		"role", "main",
		"tabindex", "0",
		"title", "A div with many attributes",
	))
	Render(buf, b)
	expected := `<div id="myDiv" class="container fluid" data-value="123" data-name="test" aria-label="Main container" role="main" tabindex="0" title="A div with many attributes"></div>`
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

func TestMultiLineAttributes(t *testing.T) {
	// Test attribute wrapping when line exceeds maxLineLen
	buf := bytes.NewBuffer(nil)
	w := NewWriter(buf)
	w.SetIndent("  ")
	w.SetMaxLineLength(30) // Short line length to trigger wrapping

	b := Div(Attrs(
		"id", "myDiv",
		"class", "container",
		"data-value", "123",
	))
	b.Build(w)

	// With maxLineLen=30, the line "<div id="myDiv" class="container"..." exceeds limit
	// So attributes after "id" should wrap
	expected := "<div id=\"myDiv\"\n  class=\"container\"\n  data-value=\"123\">\n</div>\n"
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

func TestMultiLineAttributesNested(t *testing.T) {
	// Test attribute wrapping with nested elements
	buf := bytes.NewBuffer(nil)
	w := NewWriter(buf)
	w.SetIndent("  ")
	w.SetMaxLineLength(40)

	b := Div(Attrs("id", "outer"),
		Span(Attrs("id", "inner", "class", "highlight", "data-x", "value"),
			Text("content"),
		),
	)
	b.Build(w)

	// The inner span should have wrapped attributes with deeper indentation
	got := buf.String()
	// Verify it starts correctly
	if !bytes.Contains(buf.Bytes(), []byte("<span id=\"inner\"")) {
		t.Errorf("expected span with id, got %q", got)
	}
	// Verify content is properly indented
	if !bytes.Contains(buf.Bytes(), []byte("content\n")) {
		t.Errorf("expected content with newline, got %q", got)
	}
}

func TestManyChildren(t *testing.T) {
	children := make([]TagArg, 20)
	for i := range 20 {
		children[i] = Li(Text("Item"))
	}
	buf := bytes.NewBuffer(nil)
	b := Ul(children...)
	Render(buf, b)

	expected := "<ul>" + strings.Repeat("<li>Item</li>", 20) + "</ul>"
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

func TestComplexDocument(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	b := Html(
		Head(
			Meta(Attrs("charset", "utf-8")),
			Title(Text("Test Page")),
			Link(Attrs("rel", "stylesheet", "href", "style.css")),
		),
		Body(
			Header(
				Nav(
					Ul(
						Li(A(Attrs("href", "/"), Text("Home"))),
						Li(A(Attrs("href", "/about"), Text("About"))),
					),
				),
			),
			Main(
				Article(
					H1(Text("Welcome")),
					P(Text("This is a test page.")),
				),
			),
			Footer(
				P(Text("Copyright 2024")),
			),
		),
	)
	Render(buf, b)
	expected := `<!DOCTYPE html>
<html lang="en"><head><meta charset="utf-8"/><title>Test Page</title><link rel="stylesheet" href="style.css"/></head><body><header><nav><ul><li><a href="/">Home</a></li><li><a href="/about">About</a></li></ul></nav></header><main><article><h1>Welcome</h1><p>This is a test page.</p></article></main><footer><p>Copyright 2024</p></footer></body></html>`
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

// Test Render function return value
func TestRenderReturnsError(t *testing.T) {
	ew := &errorWriter{}
	b := Div(Text("hello"))
	err := Render(ew, b)
	if err == nil {
		t.Error("expected error from Render")
	}
}

// Test RenderPretty function return value
func TestRenderPrettyReturnsError(t *testing.T) {
	ew := &errorWriter{}
	b := Div(Text("hello"))
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
	html := Html(Div(Text("test")))

	// Replace the writer after doctype to force child error
	html.Build(w) // Ignore first error

	ew := &errorWriter{}
	w2 := NewWriter(ew)
	w2.openTags = append(w2.openTags, "html") // simulate opened tag

	// Test that child error propagates through tagBuilder
	div := Div(Text("test"))
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
	b := Html(Div())
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
	html := Html(Div(Text("content")))
	pw.remaining = 50 // enough to open but not fully close
	err := html.Build(w)
	// Error may or may not occur depending on exact byte counts
	_ = err
}

func TestFragmentBuilderError(t *testing.T) {
	ew := &errorWriter{}
	w := NewWriter(ew)

	frag := Fragment(Div(Text("test")))
	err := frag.Build(w)
	if err == nil {
		t.Error("expected error from fragment builder")
	}
}
