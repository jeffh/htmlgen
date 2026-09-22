package rkt

import (
	"math"
	"testing"
	"time"

	"github.com/jeffh/htmlgen/h"
)

// stringerProp is a fmt.Stringer used to exercise Prop's Stringer codec.
type stringerProp struct{ v string }

func (s stringerProp) String() string { return "str:" + s.v }

// textProp is an encoding.TextMarshaler used to exercise Prop's
// TextMarshaler codec.
type textProp struct{ v string }

func (t textProp) MarshalText() ([]byte, error) { return []byte("text:" + t.v), nil }

// jsonProp is a plain struct, which Prop falls back to JSON-encoding.
type jsonProp struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func TestBundleConstants(t *testing.T) {
	if Version != "1.0.4" {
		t.Errorf("Version = %q, want %q", Version, "1.0.4")
	}
	const wantBundle = "https://cdn.jsdelivr.net/gh/starfederation/datastar@v1.0.4/bundles/datastar-rocket.js"
	if BundleURL != wantBundle {
		t.Errorf("BundleURL = %q, want %q", BundleURL, wantBundle)
	}
	const wantAliased = "https://cdn.jsdelivr.net/gh/starfederation/datastar@v1.0.4/bundles/datastar-rocket-aliased.js"
	if AliasedBundleURL != wantAliased {
		t.Errorf("AliasedBundleURL = %q, want %q", AliasedBundleURL, wantAliased)
	}
}

func TestScriptAttrs(t *testing.T) {
	attrs := ScriptAttrs()
	if got, ok := attrs.Get("type"); !ok || got != "module" {
		t.Errorf("ScriptAttrs()[type] = %q, %v, want %q, true", got, ok, "module")
	}
	if got, ok := attrs.Get("src"); !ok || got != BundleURL {
		t.Errorf("ScriptAttrs()[src] = %q, %v, want %q, true", got, ok, BundleURL)
	}
}

func TestScript(t *testing.T) {
	tests := []struct {
		name     string
		render   func(*h.B)
		expected string
	}{
		{
			"Script",
			Script,
			`<script type="module" src="https://cdn.jsdelivr.net/gh/starfederation/datastar@v1.0.4/bundles/datastar-rocket.js"></script>`,
		},
		{
			"AliasedScript",
			AliasedScript,
			`<script type="module" src="https://cdn.jsdelivr.net/gh/starfederation/datastar@v1.0.4/bundles/datastar-rocket-aliased.js"></script>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := h.RenderString(tt.render); got != tt.expected {
				t.Errorf("%s rendered %q, want %q", tt.name, got, tt.expected)
			}
		})
	}
}

func TestProp(t *testing.T) {
	stamp := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	offset := time.FixedZone("UTC-7", -7*3600)

	tests := []struct {
		name      string
		attr      h.Attribute
		wantName  string
		wantValue string
	}{
		{"bool true", Prop("isOpen", true), "is-open", "true"},
		{"bool false", Prop("isOpen", false), "is-open", "false"},
		{"string", Prop("label", "Count"), "label", "Count"},
		{"empty string", Prop("label", ""), "label", ""},
		{"int", Prop("step", 2), "step", "2"},
		{"negative int", Prop("step", -2), "step", "-2"},
		{"int64", Prop("step", int64(9007199254740993)), "step", "9007199254740993"},
		{"uint", Prop("step", uint(7)), "step", "7"},
		{"uint8", Prop("step", uint8(7)), "step", "7"},
		{"float64", Prop("ratio", 0.5), "ratio", "0.5"},
		{"float64 integral", Prop("ratio", 2.0), "ratio", "2"},
		{"float32", Prop("ratio", float32(1.5)), "ratio", "1.5"},
		{"date", Prop("startsAt", stamp), "starts-at", "2026-01-02T03:04:05.000Z"},
		{"date converted to UTC", Prop("startsAt", stamp.In(offset)), "starts-at", "2026-01-02T03:04:05.000Z"},
		{"bytes", Prop("blob", []byte("hi")), "blob", "aGk="},
		{"nil", Prop("label", nil), "label", ""},
		{"stringer", Prop("label", stringerProp{"x"}), "label", "str:x"},
		{"text marshaler", Prop("label", textProp{"x"}), "label", "text:x"},
		{"struct is JSON", Prop("config", jsonProp{Name: "a", Count: 1}), "config", `{"name":"a","count":1}`},
		{"pointer to struct is JSON", Prop("config", &jsonProp{Name: "a", Count: 1}), "config", `{"name":"a","count":1}`},
		{"map is JSON", Prop("config", map[string]int{"a": 1}), "config", `{"a":1}`},
		{"slice is JSON", Prop("items", []int{1, 2}), "items", "[1,2]"},
		{"kebab name unchanged", Prop("max-items", 5), "max-items", "5"},
		{"camel name kebabbed", Prop("maxItems", 5), "max-items", "5"},
		{"pascal name kebabbed", Prop("MaxItems", 5), "max-items", "5"},
		{"lowercase name unchanged", Prop("label", "x"), "label", "x"},
		{"digits keep their word", Prop("row2Items", 5), "row2-items", "5"},
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

func TestPropRejectsNonFiniteFloats(t *testing.T) {
	tests := []struct {
		name  string
		value float64
	}{
		{"NaN", math.NaN()},
		{"positive infinity", math.Inf(1)},
		{"negative infinity", math.Inf(-1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run("Prop", func(t *testing.T) {
				defer expectPanicContaining(t, "not representable as a number prop")
				_ = Prop("ratio", tt.value)
			})
			t.Run("NumberProp", func(t *testing.T) {
				defer expectPanicContaining(t, "not representable as a number prop")
				_ = NumberProp("ratio", tt.value)
			})
		})
	}
}

func TestTypedProps(t *testing.T) {
	stamp := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)

	tests := []struct {
		name      string
		attr      h.Attribute
		wantName  string
		wantValue string
	}{
		{"JSONProp string is quoted", JSONProp("label", "Count"), "label", `"Count"`},
		{"JSONProp slice", JSONProp("items", []int{1, 2}), "items", "[1,2]"},
		{"JSONProp nil", JSONProp("config", nil), "config", "null"},
		{"JSONProp kebabs the name", JSONProp("maxItems", 5), "max-items", "5"},
		{"DateProp", DateProp("startsAt", stamp), "starts-at", "2026-01-02T03:04:05.000Z"},
		{"BinProp", BinProp("blob", []byte{0xff, 0xfe}), "blob", "//4="},
		{"BinProp empty", BinProp("blob", nil), "blob", ""},
		{"BoolProp true", BoolProp("isOpen", true), "is-open", "true"},
		{"BoolProp false", BoolProp("isOpen", false), "is-open", "false"},
		{"NumberProp", NumberProp("step", 0.5), "step", "0.5"},
		{"NumberProp integral", NumberProp("step", 2), "step", "2"},
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

func TestProps(t *testing.T) {
	attrs := Props(map[string]any{
		"step":     2,
		"label":    "Count",
		"maxItems": 5,
		"isOpen":   true,
	})
	want := h.Attributes{
		{Name: "is-open", Value: "true"},
		{Name: "label", Value: "Count"},
		{Name: "max-items", Value: "5"},
		{Name: "step", Value: "2"},
	}
	if len(attrs) != len(want) {
		t.Fatalf("Props() = %#v, want %#v", attrs, want)
	}
	for i := range want {
		if attrs[i] != want[i] {
			t.Errorf("Props()[%d] = %#v, want %#v", i, attrs[i], want[i])
		}
	}
}

func TestPropsEmpty(t *testing.T) {
	if got := Props(nil); got != nil {
		t.Errorf("Props(nil) = %#v, want nil", got)
	}
}

func TestComponent(t *testing.T) {
	got := h.RenderString(func(b *h.B) {
		Component(b, "counter-widget", h.AttrsOf(Prop("step", 2)), func(b *h.B) {
			b.Text("hi")
		})
	})
	want := `<counter-widget step="2">hi</counter-widget>`
	if got != want {
		t.Errorf("Component() rendered %q, want %q", got, want)
	}
}

func TestComponentNilBody(t *testing.T) {
	got := h.RenderString(func(b *h.B) {
		Component(b, "counter-widget", nil, nil)
	})
	want := `<counter-widget></counter-widget>`
	if got != want {
		t.Errorf("Component() rendered %q, want %q", got, want)
	}
}

func TestComponentRejectsInvalidTags(t *testing.T) {
	tests := []struct {
		name string
		tag  string
	}{
		{"no hyphen", "widget"},
		{"uppercase", "Counter-Widget"},
		{"trailing uppercase", "counter-Widget"},
		{"leading hyphen", "-widget"},
		{"leading digit", "1-widget"},
		{"empty", ""},
		{"space", "counter widget"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer expectPanicContaining(t, "invalid custom element name")
			_ = h.RenderString(func(b *h.B) {
				Component(b, tt.tag, nil, nil)
			})
		})
	}
}
