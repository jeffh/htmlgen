package h

import (
	"bytes"
	"html/template"
	"io"
	"testing"
)

// Benchmarks comparing the streaming htmlgen API against html/template.
//
// Each pair renders byte-identical (or as close as the two APIs allow) output
// into the same reusable bytes.Buffer so the comparison measures generation
// cost rather than allocation of the destination.

// ============================================================================
// Simple Element
// ============================================================================

func BenchmarkSimpleDiv_HtmlGen(b *testing.B) {
	b.ReportAllocs()
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		Render(&buf, func(h *B) {
			h.Div(func(h *B) { h.Text("Hello, World!") })
		})
	}
}

func BenchmarkSimpleDiv_Template(b *testing.B) {
	b.ReportAllocs()
	tmpl := template.Must(template.New("div").Parse(`<div>{{.}}</div>`))
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		tmpl.Execute(&buf, "Hello, World!")
	}
}

// ============================================================================
// Element with Attributes
// ============================================================================

func BenchmarkDivWithAttrs_HtmlGen(b *testing.B) {
	b.ReportAllocs()
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		Render(&buf, func(h *B) {
			h.Div(
				Attrs("id", "main", "class", "container fluid", "data-value", "123"),
				func(h *B) { h.Text("Content") },
			)
		})
	}
}

func BenchmarkDivWithAttrs_Template(b *testing.B) {
	b.ReportAllocs()
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
// Nested Elements
// ============================================================================

func BenchmarkNestedElements_HtmlGen(b *testing.B) {
	b.ReportAllocs()
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		Render(&buf, func(h *B) {
			h.Div(func(h *B) {
				h.Header(func(h *B) {
					h.H1(func(h *B) { h.Text("Title") })
				})
				h.Main(func(h *B) {
					h.P(func(h *B) { h.Text("Paragraph 1") })
					h.P(func(h *B) { h.Text("Paragraph 2") })
				})
				h.Footer(func(h *B) {
					h.Span(func(h *B) { h.Text("Footer text") })
				})
			})
		})
	}
}

func BenchmarkNestedElements_Template(b *testing.B) {
	b.ReportAllocs()
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
// Lists
// ============================================================================

func BenchmarkList10Items_HtmlGen(b *testing.B) {
	b.ReportAllocs()
	items := make([]string, 10)
	for i := range items {
		items[i] = "Item"
	}
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		Render(&buf, func(h *B) {
			h.Ul(func(h *B) {
				for _, item := range items {
					h.Li(func(h *B) { h.Text(item) })
				}
			})
		})
	}
}

func BenchmarkList10Items_Template(b *testing.B) {
	b.ReportAllocs()
	tmpl := template.Must(template.New("list").Parse(
		`<ul>{{range .}}<li>{{.}}</li>{{end}}</ul>`))
	items := make([]string, 10)
	for i := range items {
		items[i] = "Item"
	}
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		tmpl.Execute(&buf, items)
	}
}

func BenchmarkList100Items_HtmlGen(b *testing.B) {
	b.ReportAllocs()
	items := make([]string, 100)
	for i := range items {
		items[i] = "Item"
	}
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		Render(&buf, func(h *B) {
			h.Ul(func(h *B) {
				for _, item := range items {
					h.Li(func(h *B) { h.Text(item) })
				}
			})
		})
	}
}

func BenchmarkList100Items_Template(b *testing.B) {
	b.ReportAllocs()
	tmpl := template.Must(template.New("list").Parse(
		`<ul>{{range .}}<li>{{.}}</li>{{end}}</ul>`))
	items := make([]string, 100)
	for i := range items {
		items[i] = "Item"
	}
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		tmpl.Execute(&buf, items)
	}
}

// ============================================================================
// Tables
// ============================================================================

type TableRow struct {
	Name  string
	Email string
	Age   string
}

func makeRows(n int) []TableRow {
	rows := make([]TableRow, n)
	for i := range rows {
		rows[i] = TableRow{"John Doe", "john@example.com", "30"}
	}
	return rows
}

func renderTable(w io.Writer, rows []TableRow) {
	Render(w, func(h *B) {
		h.Table(func(h *B) {
			h.Thead(func(h *B) {
				h.Tr(func(h *B) {
					h.Th(func(h *B) { h.Text("Name") })
					h.Th(func(h *B) { h.Text("Email") })
					h.Th(func(h *B) { h.Text("Age") })
				})
			})
			h.Tbody(func(h *B) {
				for _, row := range rows {
					h.Tr(func(h *B) {
						h.Td(func(h *B) { h.Text(row.Name) })
						h.Td(func(h *B) { h.Text(row.Email) })
						h.Td(func(h *B) { h.Text(row.Age) })
					})
				}
			})
		})
	})
}

var tableTemplate = template.Must(template.New("table").Parse(
	`<table><thead><tr><th>Name</th><th>Email</th><th>Age</th></tr></thead><tbody>{{range .}}<tr><td>{{.Name}}</td><td>{{.Email}}</td><td>{{.Age}}</td></tr>{{end}}</tbody></table>`))

func BenchmarkTable10Rows_HtmlGen(b *testing.B) {
	b.ReportAllocs()
	rows := makeRows(10)
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		renderTable(&buf, rows)
	}
}

func BenchmarkTable10Rows_Template(b *testing.B) {
	b.ReportAllocs()
	rows := makeRows(10)
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		tableTemplate.Execute(&buf, rows)
	}
}

func BenchmarkTable100Rows_HtmlGen(b *testing.B) {
	b.ReportAllocs()
	rows := makeRows(100)
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		renderTable(&buf, rows)
	}
}

func BenchmarkTable100Rows_Template(b *testing.B) {
	b.ReportAllocs()
	rows := makeRows(100)
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		tableTemplate.Execute(&buf, rows)
	}
}

// ============================================================================
// Full Page Document
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

var pageData = PageData{
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

func renderPage(w io.Writer, data PageData) {
	Render(w, func(h *B) {
		h.Html(func(h *B) {
			h.Head(func(h *B) {
				h.Meta(Attrs("charset", "utf-8"))
				h.Meta(Attrs("name", "viewport", "content", "width=device-width, initial-scale=1"))
				h.Meta(Attrs("name", "description", "content", data.Description))
				h.Title(func(h *B) { h.Text(data.Title) })
				h.Link(Attrs("rel", "stylesheet", "href", "/css/style.css"))
			})
			h.Body(func(h *B) {
				h.Header(func(h *B) {
					h.Nav(func(h *B) {
						h.Ul(func(h *B) {
							for _, item := range data.NavItems {
								h.Li(func(h *B) {
									h.A(Attrs("href", item.Href), func(h *B) { h.Text(item.Text) })
								})
							}
						})
					})
				})
				h.Main(func(h *B) {
					for _, article := range data.Articles {
						h.Article(func(h *B) {
							h.H2(func(h *B) { h.Text(article.Title) })
							h.P(func(h *B) { h.Text(article.Content) })
						})
					}
				})
				h.Footer(func(h *B) {
					h.P(func(h *B) { h.Text(data.FooterText) })
				})
				h.Script(Attrs("src", "/js/app.js"))
			})
		})
	})
}

var pageTemplate = template.Must(template.New("page").Parse(`<!DOCTYPE html>
<html lang="en"><head><meta charset="utf-8"/><meta name="viewport" content="width=device-width, initial-scale=1"/><meta name="description" content="{{.Description}}"/><title>{{.Title}}</title><link rel="stylesheet" href="/css/style.css"/></head><body><header><nav><ul>{{range .NavItems}}<li><a href="{{.Href}}">{{.Text}}</a></li>{{end}}</ul></nav></header><main>{{range .Articles}}<article><h2>{{.Title}}</h2><p>{{.Content}}</p></article>{{end}}</main><footer><p>{{.FooterText}}</p></footer><script src="/js/app.js"></script></body></html>`))

func BenchmarkFullPage_HtmlGen(b *testing.B) {
	b.ReportAllocs()
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		renderPage(&buf, pageData)
	}
}

func BenchmarkFullPage_Template(b *testing.B) {
	b.ReportAllocs()
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		pageTemplate.Execute(&buf, pageData)
	}
}

// ============================================================================
// Escaping
// ============================================================================

const unsafeContent = `<script>alert("XSS")</script> & "quotes" 'apostrophe'`

func BenchmarkEscaping_HtmlGen(b *testing.B) {
	b.ReportAllocs()
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		Render(&buf, func(h *B) {
			h.Div(
				Attrs("data-unsafe", unsafeContent),
				func(h *B) { h.Text(unsafeContent) },
			)
		})
	}
}

func BenchmarkEscaping_Template(b *testing.B) {
	b.ReportAllocs()
	tmpl := template.Must(template.New("escape").Parse(
		`<div data-unsafe="{{.}}">{{.}}</div>`))
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		tmpl.Execute(&buf, unsafeContent)
	}
}

// BenchmarkTextPlain measures the no-escaping fast path in isolation.
func BenchmarkTextPlain_HtmlGen(b *testing.B) {
	b.ReportAllocs()
	plain := "The quick brown fox jumps over the lazy dog"
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		Render(&buf, func(h *B) {
			h.Div(func(h *B) { h.Text(plain) })
		})
	}
}

// ============================================================================
// Deep Nesting
// ============================================================================

func renderNested(h *B, depth int) {
	if depth <= 0 {
		h.Text("Nested")
		return
	}
	h.Div(func(h *B) { renderNested(h, depth-1) })
}

func BenchmarkDeepNesting10_HtmlGen(b *testing.B) {
	b.ReportAllocs()
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		Render(&buf, func(h *B) { renderNested(h, 10) })
	}
}

func BenchmarkDeepNesting10_Template(b *testing.B) {
	b.ReportAllocs()
	tmpl := template.Must(template.New("deep").Parse(
		`<div><div><div><div><div><div><div><div><div><div>{{.}}</div></div></div></div></div></div></div></div></div></div>`))
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		tmpl.Execute(&buf, "Nested")
	}
}

func BenchmarkDeepNesting50_HtmlGen(b *testing.B) {
	b.ReportAllocs()
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		Render(&buf, func(h *B) { renderNested(h, 50) })
	}
}

// ============================================================================
// Forms
// ============================================================================

type FormField struct {
	Name        string
	Label       string
	Type        string
	Placeholder string
	Required    bool
}

var formFields = []FormField{
	{"username", "Username", "text", "Enter username", true},
	{"email", "Email", "email", "Enter email", true},
	{"password", "Password", "password", "Enter password", true},
	{"bio", "Biography", "text", "Tell us about yourself", false},
}

func BenchmarkForm_HtmlGen(b *testing.B) {
	b.ReportAllocs()
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		Render(&buf, func(h *B) {
			h.Form(Attrs("method", "post", "action", "/submit"), func(h *B) {
				for _, field := range formFields {
					h.Div(Attrs("class", "form-group"), func(h *B) {
						h.Label(Attrs("for", field.Name), func(h *B) { h.Text(field.Label) })
						attrs := Attrs("type", field.Type, "name", field.Name, "id", field.Name, "placeholder", field.Placeholder)
						if field.Required {
							attrs.Set("required", "")
						}
						h.Input(attrs)
					})
				}
			})
		})
	}
}

func BenchmarkForm_Template(b *testing.B) {
	b.ReportAllocs()
	tmpl := template.Must(template.New("form").Parse(
		`<form method="post" action="/submit">{{range .}}<div class="form-group"><label for="{{.Name}}">{{.Label}}</label><input type="{{.Type}}" name="{{.Name}}" id="{{.Name}}" placeholder="{{.Placeholder}}"{{if .Required}} required{{end}}/></div>{{end}}</form>`))
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		tmpl.Execute(&buf, formFields)
	}
}

// ============================================================================
// Streaming to io.Discard (generation cost with no buffer growth)
// ============================================================================

func BenchmarkDiscard_HtmlGen_FullPage(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		renderPage(io.Discard, pageData)
	}
}

func BenchmarkDiscard_Template_FullPage(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		pageTemplate.Execute(io.Discard, pageData)
	}
}

// ============================================================================
// Render entry points
// ============================================================================

// BenchmarkRenderOverhead measures the fixed cost of a Render call: builder
// pool checkout, the body call, and return.
func BenchmarkRenderOverhead(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		Render(io.Discard, func(h *B) {})
	}
}

func BenchmarkRenderString_Fragment(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		RenderString(func(h *B) {
			h.Div(func(h *B) {
				h.Header(func(h *B) { h.H1(func(h *B) { h.Text("Title") }) })
				h.Main(func(h *B) { h.P(func(h *B) { h.Text("Content") }) })
			})
		})
	}
}

func BenchmarkRenderBytes_Fragment(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		RenderBytes(func(h *B) {
			h.Div(func(h *B) {
				h.Header(func(h *B) { h.H1(func(h *B) { h.Text("Title") }) })
				h.Main(func(h *B) { h.P(func(h *B) { h.Text("Content") }) })
			})
		})
	}
}

func BenchmarkRenderIndent_Fragment(b *testing.B) {
	b.ReportAllocs()
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		RenderIndent(&buf, "  ", func(h *B) {
			h.Div(func(h *B) {
				h.Header(func(h *B) { h.H1(func(h *B) { h.Text("Title") }) })
				h.Main(func(h *B) { h.P(func(h *B) { h.Text("Content") }) })
			})
		})
	}
}

// ============================================================================
// Concurrency (the builder pool under parallel renders)
// ============================================================================

func BenchmarkParallel_FullPage_HtmlGen(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var buf bytes.Buffer
		for pb.Next() {
			buf.Reset()
			renderPage(&buf, pageData)
		}
	})
}

func BenchmarkParallel_FullPage_Template(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var buf bytes.Buffer
		for pb.Next() {
			buf.Reset()
			pageTemplate.Execute(&buf, pageData)
		}
	})
}
