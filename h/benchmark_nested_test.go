package h

import (
	"bytes"
	"html/template"
	"io"
	"testing"
)

// Benchmarks for a realistic component tree of ~90 nested elements with
// dynamic text and attributes: a card grid with a header and nav, the shape a
// typical page section or component subtree takes in practice.

type GridCard struct {
	Title       string
	Description string
	ImageURL    string
	Tags        []string
	Link        string
}

type GridData struct {
	Heading  string
	Subtitle string
	Nav      []NavItem
	Cards    []GridCard
}

var gridData = GridData{
	Heading:  "Latest Articles",
	Subtitle: "Hand-picked reads from around the team",
	Nav: []NavItem{
		{"/all", "All"},
		{"/engineering", "Engineering"},
		{"/design", "Design"},
		{"/product", "Product"},
	},
	Cards: func() []GridCard {
		cards := make([]GridCard, 8)
		for i := range cards {
			cards[i] = GridCard{
				Title:       "Understanding Systems",
				Description: "A deep dive into how complex systems behave under load and what that means for design.",
				ImageURL:    "/img/articles/systems.png",
				Tags:        []string{"systems", "performance", "design"},
				Link:        "/articles/understanding-systems",
			}
		}
		return cards
	}(),
}

func renderCardGrid(w io.Writer, data GridData) {
	Render(w, func(h *B) {
		h.Div(Attrs("class", "container"), func(h *B) {
			h.Header(func(h *B) {
				h.H2(func(h *B) { h.Text(data.Heading) })
				h.P(Attrs("class", "subtitle"), func(h *B) { h.Text(data.Subtitle) })
			})
			h.Nav(func(h *B) {
				h.Ul(Attrs("class", "nav"), func(h *B) {
					for _, item := range data.Nav {
						h.Li(func(h *B) {
							h.A(Attrs("href", item.Href), func(h *B) { h.Text(item.Text) })
						})
					}
				})
			})
			h.Div(Attrs("class", "grid"), func(h *B) {
				for _, card := range data.Cards {
					h.Div(Attrs("class", "card"), func(h *B) {
						h.Img(Attrs("src", card.ImageURL, "alt", card.Title))
						h.Div(Attrs("class", "card-body"), func(h *B) {
							h.H3(func(h *B) { h.Text(card.Title) })
							h.P(func(h *B) { h.Text(card.Description) })
							h.Div(Attrs("class", "tags"), func(h *B) {
								for _, tag := range card.Tags {
									h.Span(Attrs("class", "tag"), func(h *B) { h.Text(tag) })
								}
							})
							h.A(Attrs("class", "btn", "href", card.Link), func(h *B) { h.Text("Read more") })
						})
					})
				}
			})
		})
	})
}

// renderCardGridTyped is renderCardGrid rewritten against the XxxE siblings,
// which take typed attributes and bodies instead of ...any.
func renderCardGridTyped(w io.Writer, data GridData) {
	Render(w, func(h *B) {
		h.DivE(Attrs("class", "container"), func(h *B) {
			h.HeaderE(nil, func(h *B) {
				h.H2E(nil, func(h *B) { h.Text(data.Heading) })
				h.PE(Attrs("class", "subtitle"), func(h *B) { h.Text(data.Subtitle) })
			})
			h.NavE(nil, func(h *B) {
				h.UlE(Attrs("class", "nav"), func(h *B) {
					for _, item := range data.Nav {
						h.LiE(nil, func(h *B) {
							h.AE(Attrs("href", item.Href), func(h *B) { h.Text(item.Text) })
						})
					}
				})
			})
			h.DivE(Attrs("class", "grid"), func(h *B) {
				for _, card := range data.Cards {
					h.DivE(Attrs("class", "card"), func(h *B) {
						h.ImgE(Attrs("src", card.ImageURL, "alt", card.Title))
						h.DivE(Attrs("class", "card-body"), func(h *B) {
							h.H3E(nil, func(h *B) { h.Text(card.Title) })
							h.PE(nil, func(h *B) { h.Text(card.Description) })
							h.DivE(Attrs("class", "tags"), func(h *B) {
								for _, tag := range card.Tags {
									h.SpanE(Attrs("class", "tag"), func(h *B) { h.Text(tag) })
								}
							})
							h.AE(Attrs("class", "btn", "href", card.Link), func(h *B) { h.Text("Read more") })
						})
					})
				}
			})
		})
	})
}

var cardGridTemplate = template.Must(template.New("grid").Parse(
	`<div class="container"><header><h2>{{.Heading}}</h2><p class="subtitle">{{.Subtitle}}</p></header><nav><ul class="nav">{{range .Nav}}<li><a href="{{.Href}}">{{.Text}}</a></li>{{end}}</ul></nav><div class="grid">{{range .Cards}}<div class="card"><img src="{{.ImageURL}}" alt="{{.Title}}"/><div class="card-body"><h3>{{.Title}}</h3><p>{{.Description}}</p><div class="tags">{{range .Tags}}<span class="tag">{{.}}</span>{{end}}</div><a class="btn" href="{{.Link}}">Read more</a></div></div>{{end}}</div></div>`))

// TestCardGridOutputsMatch guards benchmark fairness: all three renderers must
// produce byte-identical output.
func TestCardGridOutputsMatch(t *testing.T) {
	var got, typed, want bytes.Buffer
	renderCardGrid(&got, gridData)
	renderCardGridTyped(&typed, gridData)
	if err := cardGridTemplate.Execute(&want, gridData); err != nil {
		t.Fatal(err)
	}
	if got.String() != want.String() {
		t.Errorf("outputs differ:\nhtmlgen:  %s\ntemplate: %s", got.String(), want.String())
	}
	if typed.String() != want.String() {
		t.Errorf("outputs differ:\ntyped:    %s\ntemplate: %s", typed.String(), want.String())
	}
}

func BenchmarkCardGrid90Elements_HtmlGen(b *testing.B) {
	b.ReportAllocs()
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		renderCardGrid(&buf, gridData)
	}
}

func BenchmarkCardGrid90Elements_HtmlGenTyped(b *testing.B) {
	b.ReportAllocs()
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		renderCardGridTyped(&buf, gridData)
	}
}

func BenchmarkCardGrid90Elements_Template(b *testing.B) {
	b.ReportAllocs()
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		cardGridTemplate.Execute(&buf, gridData)
	}
}
