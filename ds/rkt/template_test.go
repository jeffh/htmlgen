package rkt

import (
	"testing"

	"github.com/jeffh/htmlgen/ds"
	"github.com/jeffh/htmlgen/h"
)

func TestFor(t *testing.T) {
	tests := []struct {
		name     string
		attr     h.Attribute
		expected string
	}{
		{"item and index", For("item", "index", SignalRef("items")), "item, index in $$items"},
		{"item only", For("item", "", SignalRef("items")), "item in $$items"},
		{"source only", For("", "", SignalRef("items")), "$$items"},
		{"global signal source", For("item", "", ds.SignalRef("items")), "item in $items"},
		{"expression source", For("row", "i", ds.Raw("$$rows.filter(r => r.on)")), "row, i in $$rows.filter(r => r.on)"},
		{"dollar identifiers", For("$item", "_i", SignalRef("items")), "$item, _i in $$items"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.attr.Name != "data-for" {
				t.Errorf("Name = %q, want %q", tt.attr.Name, "data-for")
			}
			if tt.attr.Value != tt.expected {
				t.Errorf("Value = %q, want %q", tt.attr.Value, tt.expected)
			}
		})
	}
}

func TestForRejectsInvalidIdentifiers(t *testing.T) {
	tests := []struct {
		name  string
		item  string
		index string
		want  string
	}{
		{"hyphen in item", "my-item", "", "invalid data-for item name"},
		{"leading digit in item", "1item", "", "invalid data-for item name"},
		{"space in item", "my item", "", "invalid data-for item name"},
		{"hyphen in index", "item", "my-index", "invalid data-for index name"},
		{"expression in index", "item", "i + 1", "invalid data-for index name"},
		{"index without item", "", "i", "index but no item"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer expectPanicContaining(t, tt.want)
			_ = For(tt.item, tt.index, SignalRef("items"))
		})
	}
}

func TestForRaw(t *testing.T) {
	attr := ForRaw("item, i in $$items.filter(x => x.done)")
	if attr.Name != "data-for" {
		t.Errorf("Name = %q, want %q", attr.Name, "data-for")
	}
	if want := "item, i in $$items.filter(x => x.done)"; attr.Value != want {
		t.Errorf("Value = %q, want %q", attr.Value, want)
	}
}

func TestConditionalTemplates(t *testing.T) {
	tests := []struct {
		name      string
		attr      h.Attribute
		wantName  string
		wantValue string
	}{
		{"If", If(Sig("open").Value()), "data-if", "$$open"},
		{"If negated", If(Sig("open").Not()), "data-if", "!$$open"},
		{"If comparison", If(Sig("count").Eq(ds.Int(0))), "data-if", "($$count === 0)"},
		{"ElseIf", ElseIf(Sig("busy").Value()), "data-else-if", "$$busy"},
		{"ElseIf global signal", ElseIf(ds.SignalRef("busy")), "data-else-if", "$busy"},
		{"Else", Else(), "data-else", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.attr.Name != tt.wantName {
				t.Errorf("Name = %q, want %q", tt.attr.Name, tt.wantName)
			}
			if tt.attr.Value != tt.wantValue {
				t.Errorf("Value = %q, want %q", tt.attr.Value, tt.wantValue)
			}
		})
	}
}

func TestRef(t *testing.T) {
	tests := []struct {
		name     string
		attr     h.Attribute
		expected string
	}{
		{"plain", Ref("input").Attribute(), "data-ref:input"},
		{"camel", Ref("my-input").Case(ds.CamelCase).Attribute(), "data-ref:my-input__case.camel"},
		{"kebab", Ref("myInput").Case(ds.KebabCase).Attribute(), "data-ref:myInput__case.kebab"},
		{"snake", Ref("myInput").Case(ds.SnakeCase).Attribute(), "data-ref:myInput__case.snake"},
		{"pascal", Ref("myInput").Case(ds.PascalCase).Attribute(), "data-ref:myInput__case.pascal"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.attr.Name != tt.expected {
				t.Errorf("Name = %q, want %q", tt.attr.Name, tt.expected)
			}
			if tt.attr.Value != "" {
				t.Errorf("Value = %q, want empty", tt.attr.Value)
			}
		})
	}
}

func TestRoot(t *testing.T) {
	tests := []struct {
		name      string
		attr      h.Attribute
		wantName  string
		wantValue string
	}{
		{"bind key", Root(ds.BindKey("name")), "data-bind:name__root", ""},
		{"bind value form", Root(ds.Bind("name")), "data-bind__root", "name"},
		{"signal", Root(ds.Signal("count", 0)), "data-signals:count__root", "0"},
		{"signals object", Root(ds.Signals(map[string]any{"count": 0})), "data-signals__root", `{"count":0}`},
		{"computed", Root(ds.Computed("total", ds.Raw("$price * $qty"))), "data-computed:total__root", "$price * $qty"},
		{"indicator", Root(ds.Indicator("loading")), "data-indicator__root", "loading"},
		{"indicator key", Root(ds.IndicatorKey("loading")), "data-indicator:loading__root", ""},
		{"plain attribute", Root(h.Attr("data-bind", "name")), "data-bind__root", "name"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.attr.Name != tt.wantName {
				t.Errorf("Name = %q, want %q", tt.attr.Name, tt.wantName)
			}
			if tt.attr.Value != tt.wantValue {
				t.Errorf("Value = %q, want %q", tt.attr.Value, tt.wantValue)
			}
		})
	}
}

func TestRootIgnoresEmptyBuilders(t *testing.T) {
	if got := Root(nil); got != (h.Attribute{}) {
		t.Errorf("Root(nil) = %#v, want zero Attribute", got)
	}
	if got := Root(h.Attribute{}); got != (h.Attribute{}) {
		t.Errorf("Root(zero) = %#v, want zero Attribute", got)
	}
}
