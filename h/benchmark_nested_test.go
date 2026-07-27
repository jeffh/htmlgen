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
			h.Header(nil, func(h *B) {
				h.H2(nil, func(h *B) { h.Text(data.Heading) })
				h.P(Attrs("class", "subtitle"), func(h *B) { h.Text(data.Subtitle) })
			})
			h.Nav(nil, func(h *B) {
				h.Ul(Attrs("class", "nav"), func(h *B) {
					for _, item := range data.Nav {
						h.Li(nil, func(h *B) {
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
							h.H3(nil, func(h *B) { h.Text(card.Title) })
							h.P(nil, func(h *B) { h.Text(card.Description) })
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

var cardGridTemplate = template.Must(template.New("grid").Parse(
	`<div class="container"><header><h2>{{.Heading}}</h2><p class="subtitle">{{.Subtitle}}</p></header><nav><ul class="nav">{{range .Nav}}<li><a href="{{.Href}}">{{.Text}}</a></li>{{end}}</ul></nav><div class="grid">{{range .Cards}}<div class="card"><img src="{{.ImageURL}}" alt="{{.Title}}"/><div class="card-body"><h3>{{.Title}}</h3><p>{{.Description}}</p><div class="tags">{{range .Tags}}<span class="tag">{{.}}</span>{{end}}</div><a class="btn" href="{{.Link}}">Read more</a></div></div>{{end}}</div></div>`))

// TestCardGridOutputsMatch guards benchmark fairness: both renderers must
// produce byte-identical output.
func TestCardGridOutputsMatch(t *testing.T) {
	var got, want bytes.Buffer
	renderCardGrid(&got, gridData)
	if err := cardGridTemplate.Execute(&want, gridData); err != nil {
		t.Fatal(err)
	}
	if got.String() != want.String() {
		t.Errorf("outputs differ:\nhtmlgen:  %s\ntemplate: %s", got.String(), want.String())
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

func BenchmarkCardGrid90Elements_Template(b *testing.B) {
	b.ReportAllocs()
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		cardGridTemplate.Execute(&buf, gridData)
	}
}
