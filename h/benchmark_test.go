package h

import (
	"io"
	"testing"
)

func BenchmarkRender(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		Render(io.Discard, func(h *B) {
			h.Div(Attrs("class", "container"), func(h *B) {
				h.H1(func(h *B) { h.Text("Title") })
				h.P(func(h *B) { h.Text("Content") })
			})
		})
	}
}

func BenchmarkRenderList(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		Render(io.Discard, func(h *B) {
			h.Ul(func(h *B) {
				for range 100 {
					h.Li(func(h *B) { h.Text("Item") })
				}
			})
		})
	}
}

func BenchmarkRenderEscaped(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		Render(io.Discard, func(h *B) {
			h.Div(func(h *B) {
				h.Text(`<script>alert("xss")</script>`)
			})
		})
	}
}
