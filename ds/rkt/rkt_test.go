package rkt

import (
	"testing"

	"github.com/jeffh/htmlgen/ds"
	"github.com/jeffh/htmlgen/h"
)

func TestRenderComponentMarkup(t *testing.T) {
	count := Sig("count")

	got := h.RenderString(func(b *h.B) {
		Component(b, "counter-widget", h.AttrsOf(Prop("step", 2), Prop("label", "Count")), func(b *h.B) {
			b.Button(h.AttrsOf(
				ds.OnClick(Action("increment")),
				ds.Text(count.Value()),
			), nil)
			b.Ul(h.AttrsOf(Root(ds.BindKey("selected"))), func(b *h.B) {
				b.Template(h.AttrsOf(For("item", "i", SignalRef("items"))), func(b *h.B) {
					b.Li(h.AttrsOf(ds.Text(ds.Raw("item")), Ref("row").Case(ds.CamelCase)), nil)
				})
			})
			b.Template(h.AttrsOf(If(count.Eq(ds.Int(0)))), func(b *h.B) {
				b.Span(nil, func(b *h.B) { b.Text("none yet") })
			})
			b.Template(h.AttrsOf(ElseIf(count.Value())), func(b *h.B) {
				b.Span(nil, func(b *h.B) { b.Text("counting") })
			})
			b.Template(h.AttrsOf(Else()), func(b *h.B) {
				b.Span(nil, func(b *h.B) { b.Text("idle") })
			})
		})
	})

	want := `<counter-widget step="2" label="Count">` +
		`<button data-on:click="@increment()" data-text="$$count"></button>` +
		`<ul data-bind:selected__root>` +
		`<template data-for="item, i in $$items">` +
		`<li data-text="item" data-ref:row__case.camel></li>` +
		`</template>` +
		`</ul>` +
		`<template data-if="($$count === 0)"><span>none yet</span></template>` +
		`<template data-else-if="$$count"><span>counting</span></template>` +
		`<template data-else><span>idle</span></template>` +
		`</counter-widget>`

	if got != want {
		t.Errorf("rendered:\n%s\nwant:\n%s", got, want)
	}
}

func TestRenderDocumentHead(t *testing.T) {
	got := h.RenderString(func(b *h.B) {
		b.Head(nil, func(b *h.B) {
			Script(b)
		})
	})
	want := `<head><script type="module" src="https://cdn.jsdelivr.net/gh/starfederation/datastar@v1.0.4/bundles/datastar-rocket.js"></script></head>`
	if got != want {
		t.Errorf("rendered %q, want %q", got, want)
	}
}

func TestRenderDispatchEscapesQuotes(t *testing.T) {
	got := h.RenderString(func(b *h.B) {
		b.Button(h.AttrsOf(ds.OnClick(Dispatch("increment", ds.Int(1)))), nil)
	})
	want := `<button data-on:click="@dispatchRocket(&#34;increment&#34;, 1)"></button>`
	if got != want {
		t.Errorf("rendered %q, want %q", got, want)
	}
}
