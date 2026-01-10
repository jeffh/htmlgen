package h

import (
	"bytes"
	"html/template"
	"io"
	"testing"
)

// Benchmarks comparing htmlgen declarative API vs html/template

// ============================================================================
// Simple Element Benchmarks
// ============================================================================

func BenchmarkSimpleDiv_HtmlGen(b *testing.B) {
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		Render(&buf, Div(Text("Hello, World!")))
	}
}

func BenchmarkSimpleDiv_Template(b *testing.B) {
	tmpl := template.Must(template.New("div").Parse(`<div>{{.}}</div>`))
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		tmpl.Execute(&buf, "Hello, World!")
	}
}

// ============================================================================
// Element with Attributes Benchmarks
// ============================================================================

func BenchmarkDivWithAttrs_HtmlGen(b *testing.B) {
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		Render(&buf, Div(
			Attrs("id", "main", "class", "container fluid", "data-value", "123"),
			Text("Content"),
		))
	}
}

func BenchmarkDivWithAttrs_Template(b *testing.B) {
	tmpl := template.Must(template.New("div").Parse(
		`<div id="{{.ID}}" class="{{.Class}}" data-value="{{.DataValue}}">{{.Content}}</div>`))
	data := struct {
		ID        string
		Class     string
		DataValue string
		Content   string
	}{"main", "container fluid", "123", "Content"}
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		tmpl.Execute(&buf, data)
	}
}

// ============================================================================
// Nested Elements Benchmarks
// ============================================================================

func BenchmarkNestedElements_HtmlGen(b *testing.B) {
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		Render(&buf, Div(
			Header(
				H1(Text("Title")),
			),
			Main(
				P(Text("Paragraph 1")),
				P(Text("Paragraph 2")),
			),
			Footer(
				Span(Text("Footer text")),
			),
		))
	}
}

func BenchmarkNestedElements_Template(b *testing.B) {
	tmpl := template.Must(template.New("nested").Parse(
		`<div><header><h1>{{.Title}}</h1></header><main><p>{{.P1}}</p><p>{{.P2}}</p></main><footer><span>{{.Footer}}</span></footer></div>`))
	data := struct {
		Title, P1, P2, Footer string
	}{"Title", "Paragraph 1", "Paragraph 2", "Footer text"}
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		tmpl.Execute(&buf, data)
	}
}

// ============================================================================
// List Benchmarks
// ============================================================================

func BenchmarkList10Items_HtmlGen(b *testing.B) {
	items := make([]TagArg, 10)
	for i := range 10 {
		items[i] = Li(Text("Item"))
	}
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		Render(&buf, Ul(items...))
	}
}

func BenchmarkList10Items_Template(b *testing.B) {
	tmpl := template.Must(template.New("list").Parse(
		`<ul>{{range .}}<li>{{.}}</li>{{end}}</ul>`))
	items := make([]string, 10)
	for i := range 10 {
		items[i] = "Item"
	}
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		tmpl.Execute(&buf, items)
	}
}

func BenchmarkList100Items_HtmlGen(b *testing.B) {
	items := make([]TagArg, 100)
	for i := range 100 {
		items[i] = Li(Text("Item"))
	}
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		Render(&buf, Ul(items...))
	}
}

func BenchmarkList100Items_Template(b *testing.B) {
	tmpl := template.Must(template.New("list").Parse(
		`<ul>{{range .}}<li>{{.}}</li>{{end}}</ul>`))
	items := make([]string, 100)
	for i := range 100 {
		items[i] = "Item"
	}
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		tmpl.Execute(&buf, items)
	}
}

// ============================================================================
// Table Benchmarks
// ============================================================================

type TableRow struct {
	Name  string
	Email string
	Age   int
}

func BenchmarkTable10Rows_HtmlGen(b *testing.B) {
	rows := make([]TableRow, 10)
	for i := range 10 {
		rows[i] = TableRow{"John Doe", "john@example.com", 30}
	}
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		tableRows := make([]TagArg, len(rows))
		for i, r := range rows {
			tableRows[i] = Tr(
				Td(Text(r.Name)),
				Td(Text(r.Email)),
				Td(Text("30")),
			)
		}
		Render(&buf, Table(
			Thead(Tr(
				Th(Text("Name")),
				Th(Text("Email")),
				Th(Text("Age")),
			)),
			Tbody(tableRows...),
		))
	}
}

func BenchmarkTable10Rows_Template(b *testing.B) {
	tmpl := template.Must(template.New("table").Parse(
		`<table><thead><tr><th>Name</th><th>Email</th><th>Age</th></tr></thead><tbody>{{range .}}<tr><td>{{.Name}}</td><td>{{.Email}}</td><td>{{.Age}}</td></tr>{{end}}</tbody></table>`))
	rows := make([]TableRow, 10)
	for i := range 10 {
		rows[i] = TableRow{"John Doe", "john@example.com", 30}
	}
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		tmpl.Execute(&buf, rows)
	}
}

func BenchmarkTable100Rows_HtmlGen(b *testing.B) {
	rows := make([]TableRow, 100)
	for i := range 100 {
		rows[i] = TableRow{"John Doe", "john@example.com", 30}
	}
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		tableRows := make([]TagArg, len(rows))
		for i, r := range rows {
			tableRows[i] = Tr(
				Td(Text(r.Name)),
				Td(Text(r.Email)),
				Td(Text("30")),
			)
		}
		Render(&buf, Table(
			Thead(Tr(
				Th(Text("Name")),
				Th(Text("Email")),
				Th(Text("Age")),
			)),
			Tbody(tableRows...),
		))
	}
}

func BenchmarkTable100Rows_Template(b *testing.B) {
	tmpl := template.Must(template.New("table").Parse(
		`<table><thead><tr><th>Name</th><th>Email</th><th>Age</th></tr></thead><tbody>{{range .}}<tr><td>{{.Name}}</td><td>{{.Email}}</td><td>{{.Age}}</td></tr>{{end}}</tbody></table>`))
	rows := make([]TableRow, 100)
	for i := range 100 {
		rows[i] = TableRow{"John Doe", "john@example.com", 30}
	}
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		tmpl.Execute(&buf, rows)
	}
}

// ============================================================================
// Full Page Document Benchmarks
// ============================================================================

type PageData struct {
	Title       string
	Description string
	NavItems    []NavItem
	Articles    []ArticleData
	FooterText  string
}

type NavItem struct {
	Href string
	Text string
}

type ArticleData struct {
	Title   string
	Content string
}

func BenchmarkFullPage_HtmlGen(b *testing.B) {
	data := PageData{
		Title:       "My Website",
		Description: "A sample website",
		NavItems: []NavItem{
			{"/", "Home"},
			{"/about", "About"},
			{"/contact", "Contact"},
		},
		Articles: []ArticleData{
			{"First Post", "This is the first post content."},
			{"Second Post", "This is the second post content."},
			{"Third Post", "This is the third post content."},
		},
		FooterText: "Copyright 2024",
	}
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		navItems := make([]TagArg, len(data.NavItems))
		for i, n := range data.NavItems {
			navItems[i] = Li(A(Attrs("href", n.Href), Text(n.Text)))
		}
		articles := make([]TagArg, len(data.Articles))
		for i, a := range data.Articles {
			articles[i] = Article(
				H2(Text(a.Title)),
				P(Text(a.Content)),
			)
		}
		Render(&buf, Html(
			Head(
				Meta(Attrs("charset", "utf-8")),
				Meta(Attrs("name", "viewport", "content", "width=device-width, initial-scale=1")),
				Meta(Attrs("name", "description", "content", data.Description)),
				Title(Text(data.Title)),
				Link(Attrs("rel", "stylesheet", "href", "/css/style.css")),
			),
			Body(
				Header(
					Nav(Ul(navItems...)),
				),
				Main(articles...),
				Footer(P(Text(data.FooterText))),
				Script(Attrs("src", "/js/app.js")),
			),
		))
	}
}

func BenchmarkFullPage_Template(b *testing.B) {
	tmpl := template.Must(template.New("page").Parse(`<!DOCTYPE html>
<html lang="en"><head><meta charset="utf-8"/><meta name="viewport" content="width=device-width, initial-scale=1"/><meta name="description" content="{{.Description}}"/><title>{{.Title}}</title><link rel="stylesheet" href="/css/style.css"/></head><body><header><nav><ul>{{range .NavItems}}<li><a href="{{.Href}}">{{.Text}}</a></li>{{end}}</ul></nav></header><main>{{range .Articles}}<article><h2>{{.Title}}</h2><p>{{.Content}}</p></article>{{end}}</main><footer><p>{{.FooterText}}</p></footer><script src="/js/app.js"></script></body></html>`))
	data := PageData{
		Title:       "My Website",
		Description: "A sample website",
		NavItems: []NavItem{
			{"/", "Home"},
			{"/about", "About"},
			{"/contact", "Contact"},
		},
		Articles: []ArticleData{
			{"First Post", "This is the first post content."},
			{"Second Post", "This is the second post content."},
			{"Third Post", "This is the third post content."},
		},
		FooterText: "Copyright 2024",
	}
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		tmpl.Execute(&buf, data)
	}
}

// ============================================================================
// Escaping Benchmarks
// ============================================================================

func BenchmarkEscaping_HtmlGen(b *testing.B) {
	unsafeContent := `<script>alert("XSS")</script> & "quotes" 'apostrophe'`
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		Render(&buf, Div(
			Attrs("data-unsafe", unsafeContent),
			Text(unsafeContent),
		))
	}
}

func BenchmarkEscaping_Template(b *testing.B) {
	tmpl := template.Must(template.New("escape").Parse(
		`<div data-unsafe="{{.}}">{{.}}</div>`))
	unsafeContent := `<script>alert("XSS")</script> & "quotes" 'apostrophe'`
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		tmpl.Execute(&buf, unsafeContent)
	}
}

// ============================================================================
// Deep Nesting Benchmarks
// ============================================================================

func BenchmarkDeepNesting_HtmlGen(b *testing.B) {
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		Render(&buf, Div(
			Div(
				Div(
					Div(
						Div(
							Div(
								Div(
									Div(
										Div(
											Div(Text("Deeply nested")),
										),
									),
								),
							),
						),
					),
				),
			),
		))
	}
}

func BenchmarkDeepNesting_Template(b *testing.B) {
	tmpl := template.Must(template.New("deep").Parse(
		`<div><div><div><div><div><div><div><div><div><div>{{.}}</div></div></div></div></div></div></div></div></div></div>`))
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		tmpl.Execute(&buf, "Deeply nested")
	}
}

// buildNestedDivs creates a nested div structure of the specified depth
func buildNestedDivs(depth int) Builder {
	if depth <= 0 {
		return Text("Nested")
	}
	return Div(buildNestedDivs(depth - 1))
}

func BenchmarkDeepNesting20_HtmlGen(b *testing.B) {
	tree := buildNestedDivs(20)
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		Render(&buf, tree)
	}
}

func BenchmarkDeepNesting50_HtmlGen(b *testing.B) {
	tree := buildNestedDivs(50)
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		Render(&buf, tree)
	}
}

// ============================================================================
// Form Benchmarks
// ============================================================================

type FormField struct {
	Name        string
	Label       string
	Type        string
	Placeholder string
	Required    bool
}

func BenchmarkForm_HtmlGen(b *testing.B) {
	fields := []FormField{
		{"username", "Username", "text", "Enter username", true},
		{"email", "Email", "email", "Enter email", true},
		{"password", "Password", "password", "Enter password", true},
		{"bio", "Biography", "text", "Tell us about yourself", false},
	}
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		formFields := make([]TagArg, len(fields))
		for i, f := range fields {
			attrs := Attrs("type", f.Type, "name", f.Name, "id", f.Name, "placeholder", f.Placeholder)
			if f.Required {
				attrs.Set("required", "")
			}
			formFields[i] = Div(
				Attrs("class", "form-group"),
				Label(Attrs("for", f.Name), Text(f.Label)),
				Input(attrs),
			)
		}
		Render(&buf, Form(
			append([]TagArg{Attrs("method", "post", "action", "/submit")}, formFields...)...,
		))
	}
}

func BenchmarkForm_Template(b *testing.B) {
	tmpl := template.Must(template.New("form").Parse(
		`<form method="post" action="/submit">{{range .}}<div class="form-group"><label for="{{.Name}}">{{.Label}}</label><input type="{{.Type}}" name="{{.Name}}" id="{{.Name}}" placeholder="{{.Placeholder}}"{{if .Required}} required{{end}}/></div>{{end}}</form>`))
	fields := []FormField{
		{"username", "Username", "text", "Enter username", true},
		{"email", "Email", "email", "Enter email", true},
		{"password", "Password", "password", "Enter password", true},
		{"bio", "Biography", "text", "Tell us about yourself", false},
	}
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		tmpl.Execute(&buf, fields)
	}
}

// ============================================================================
// Writer API vs Builder API Benchmarks
// ============================================================================

func BenchmarkWriterAPI_Simple(b *testing.B) {
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		w := NewWriter(&buf)
		w.OpenTag("div", Attrs("class", "container"))
		w.Text("Hello, World!")
		w.CloseTag("div")
	}
}

func BenchmarkBuilderAPI_Simple(b *testing.B) {
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		Render(&buf, Div(Attrs("class", "container"), Text("Hello, World!")))
	}
}

func BenchmarkWriterAPI_Complex(b *testing.B) {
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		w := NewWriter(&buf)
		w.OpenTag("div", Attrs("class", "container"))
		w.OpenTag("header", nil)
		w.OpenTag("h1", nil)
		w.Text("Title")
		w.CloseTag("h1")
		w.CloseTag("header")
		w.OpenTag("main", nil)
		w.OpenTag("p", nil)
		w.Text("Paragraph 1")
		w.CloseTag("p")
		w.OpenTag("p", nil)
		w.Text("Paragraph 2")
		w.CloseTag("p")
		w.CloseTag("main")
		w.CloseTag("div")
	}
}

func BenchmarkBuilderAPI_Complex(b *testing.B) {
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		Render(&buf, Div(
			Attrs("class", "container"),
			Header(H1(Text("Title"))),
			Main(
				P(Text("Paragraph 1")),
				P(Text("Paragraph 2")),
			),
		))
	}
}

// ============================================================================
// Allocation Benchmarks (using ReportAllocs)
// ============================================================================

func BenchmarkAllocations_HtmlGen_SimpleDiv(b *testing.B) {
	b.ReportAllocs()
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		Render(&buf, Div(Text("Hello")))
	}
}

func BenchmarkAllocations_Template_SimpleDiv(b *testing.B) {
	b.ReportAllocs()
	tmpl := template.Must(template.New("div").Parse(`<div>{{.}}</div>`))
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		tmpl.Execute(&buf, "Hello")
	}
}

func BenchmarkAllocations_HtmlGen_FullPage(b *testing.B) {
	b.ReportAllocs()
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		Render(&buf, Html(
			Head(
				Meta(Attrs("charset", "utf-8")),
				Title(Text("Test")),
			),
			Body(
				Header(H1(Text("Title"))),
				Main(P(Text("Content"))),
				Footer(P(Text("Footer"))),
			),
		))
	}
}

func BenchmarkAllocations_Template_FullPage(b *testing.B) {
	b.ReportAllocs()
	tmpl := template.Must(template.New("page").Parse(
		`<!DOCTYPE html>
<html lang="en"><head><meta charset="utf-8"/><title>{{.Title}}</title></head><body><header><h1>{{.Header}}</h1></header><main><p>{{.Content}}</p></main><footer><p>{{.Footer}}</p></footer></body></html>`))
	data := struct{ Title, Header, Content, Footer string }{"Test", "Title", "Content", "Footer"}
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		tmpl.Execute(&buf, data)
	}
}

// ============================================================================
// Discard Writer Benchmarks (pure generation speed, no I/O)
// ============================================================================

func BenchmarkDiscard_HtmlGen_FullPage(b *testing.B) {
	for b.Loop() {
		Render(io.Discard, Html(
			Head(
				Meta(Attrs("charset", "utf-8")),
				Title(Text("Test")),
			),
			Body(
				Header(H1(Text("Title"))),
				Main(P(Text("Content"))),
				Footer(P(Text("Footer"))),
			),
		))
	}
}

func BenchmarkDiscard_Template_FullPage(b *testing.B) {
	tmpl := template.Must(template.New("page").Parse(
		`<!DOCTYPE html>
<html lang="en"><head><meta charset="utf-8"/><title>{{.Title}}</title></head><body><header><h1>{{.Header}}</h1></header><main><p>{{.Content}}</p></main><footer><p>{{.Footer}}</p></footer></body></html>`))
	data := struct{ Title, Header, Content, Footer string }{"Test", "Title", "Content", "Footer"}
	for b.Loop() {
		tmpl.Execute(io.Discard, data)
	}
}

// ============================================================================
// Pre-built Tree Benchmarks (reusing the same tree structure)
// ============================================================================

func BenchmarkPrebuiltTree_HtmlGen(b *testing.B) {
	tree := Div(
		Header(H1(Text("Title"))),
		Main(
			P(Text("Paragraph 1")),
			P(Text("Paragraph 2")),
			P(Text("Paragraph 3")),
		),
		Footer(Span(Text("Footer"))),
	)
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		Render(&buf, tree)
	}
}

func BenchmarkPrebuiltTree_Template(b *testing.B) {
	tmpl := template.Must(template.New("tree").Parse(
		`<div><header><h1>Title</h1></header><main><p>Paragraph 1</p><p>Paragraph 2</p><p>Paragraph 3</p></main><footer><span>Footer</span></footer></div>`))
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		tmpl.Execute(&buf, nil)
	}
}

// ============================================================================
// Compiled Builder Benchmarks (pre-computed HTML bytes)
// ============================================================================

func BenchmarkCompiledTree_HtmlGen(b *testing.B) {
	tree := Div(
		Header(H1(Text("Title"))),
		Main(
			P(Text("Paragraph 1")),
			P(Text("Paragraph 2")),
			P(Text("Paragraph 3")),
		),
		Footer(Span(Text("Footer"))),
	)
	compiled := Compile(tree)
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		Render(&buf, compiled)
	}
}

func BenchmarkCompiledDeepNesting_HtmlGen(b *testing.B) {
	tree := buildNestedDivs(10)
	compiled := Compile(tree)
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		Render(&buf, compiled)
	}
}

func BenchmarkCompiledDeepNesting50_HtmlGen(b *testing.B) {
	tree := buildNestedDivs(50)
	compiled := Compile(tree)
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		Render(&buf, compiled)
	}
}

// ============================================================================
// Parameterized Compile Benchmarks
// ============================================================================

func BenchmarkCompiledParams_HtmlGen(b *testing.B) {
	title := NewParam("title")
	content := NewParam("content")
	tmpl := CompileParams(Div(
		H1(title),
		P(content),
	))
	titleVal := title.Value(Text("Hello World"))
	contentVal := content.Value(Text("Welcome to my site"))

	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		tmpl.Render(&buf, titleVal, contentVal)
	}
}

func BenchmarkCompiledParams_Template(b *testing.B) {
	tmpl := template.Must(template.New("params").Parse(
		`<div><h1>{{.Title}}</h1><p>{{.Content}}</p></div>`))
	data := struct{ Title, Content string }{"Hello World", "Welcome to my site"}
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		tmpl.Execute(&buf, data)
	}
}

func BenchmarkCompiledParamsComplex_HtmlGen(b *testing.B) {
	title := NewParam("title")
	nav := NewParam("nav")
	content := NewParam("content")
	footer := NewParam("footer")

	tmpl := CompileParams(Html(
		Head(
			Meta(Attrs("charset", "utf-8")),
			Title(title),
		),
		Body(
			Header(nav),
			Main(content),
			Footer(footer),
		),
	))

	titleVal := title.Value(Text("My Page"))
	navVal := nav.Value(Ul(Li(A(Attrs("href", "/"), Text("Home")))))
	contentVal := content.Value(P(Text("Page content here")))
	footerVal := footer.Value(Text("Copyright 2024"))

	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		tmpl.Render(&buf, titleVal, navVal, contentVal, footerVal)
	}
}

func BenchmarkDynamicEquivalent_HtmlGen(b *testing.B) {
	// Equivalent to CompiledParams but rebuilding each time
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		Render(&buf, Div(
			H1(Text("Hello World")),
			P(Text("Welcome to my site")),
		))
	}
}
