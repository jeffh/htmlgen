package rkt

import (
	"strings"
	"testing"

	"github.com/jeffh/htmlgen/ds"
	"github.com/jeffh/htmlgen/js"
)

func TestSigName(t *testing.T) {
	tests := []struct {
		name     string
		sig      Sig
		expected string
	}{
		{"plain", Sig("count"), "count"},
		{"double dollar prefixed", Sig("$$count"), "count"},
		{"single dollar prefixed", Sig("$count"), "count"},
		{"nested path", Sig("form.email"), "form.email"},
		{"only the prefix is stripped", Sig("$$$count"), "$count"},
		{"empty", Sig(""), ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sig.Name(); got != tt.expected {
				t.Errorf("Sig(%q).Name() = %q, want %q", string(tt.sig), got, tt.expected)
			}
		})
	}
}

func TestSigExpressions(t *testing.T) {
	tests := []struct {
		name     string
		value    ds.Value
		expected string
	}{
		{"Value", Sig("count").Value(), "$$count"},
		{"Value strips double dollar", Sig("$$count").Value(), "$$count"},
		{"Value strips single dollar", Sig("$count").Value(), "$$count"},
		{"Not", Sig("open").Not(), "!$$open"},
		{"Toggle", Sig("open").Toggle(), "($$open = !$$open)"},
		{"Toggle strips prefix", Sig("$$open").Toggle(), "($$open = !$$open)"},
		{"Clear", Sig("q").Clear(), "($$q = '')"},
		{"SetExpr", Sig("count").SetExpr(js.Int(1)), "($$count = 1)"},
		{"SetExpr with private signal", Sig("n").SetExpr(Sig("m").Ref()), "($$n = $$m)"},
		{"SetExpr with global signal", Sig("n").SetExpr(ds.Sig("m").Ref()), "($$n = $m)"},
		{"SetExpr with evt", Sig("q").SetExpr(ds.EvtValue), "($$q = evt.target.value)"},
		{"Eq", Sig("tab").Eq(js.String("plans")), `($$tab === "plans")`},
		{"NotEq", Sig("tab").NotEq(js.String("plans")), `($$tab !== "plans")`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := js.ToJS(tt.value.Expr()); got != tt.expected {
				t.Errorf("%s = %q, want %q", tt.name, got, tt.expected)
			}
		})
	}
}

func TestSigRef(t *testing.T) {
	if got := js.ToJS(Sig("count").Ref()); got != "$$count" {
		t.Errorf(`Sig("count").Ref() = %q, want %q`, got, "$$count")
	}
}

func TestSigRefRejectsEmptyName(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal(`Sig("").Ref() did not panic`)
		}
	}()
	_ = Sig("").Ref()
}

func TestSigSet(t *testing.T) {
	tests := []struct {
		name     string
		value    any
		expected string
	}{
		{"string is JSON encoded", "hello", `($$n = "hello")`},
		{"int is JSON encoded", 0, "($$n = 0)"},
		{"bool is JSON encoded", true, "($$n = true)"},
		{"nil is JSON encoded", nil, "($$n = null)"},
		{"slice is JSON encoded", []int{1, 2}, "($$n = [1,2])"},
		{"js.Expr passes through", js.Ident("x"), "($$n = x)"},
		{"ds.Value passes through", Sig("m").Value(), "($$n = $$m)"},
		{"ds.Value expression passes through", Sig("m").Not(), "($$n = !$$m)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := js.ToJS(Sig("n").Set(tt.value).Expr())
			if got != tt.expected {
				t.Errorf("Sig(\"n\").Set(%#v) = %q, want %q", tt.value, got, tt.expected)
			}
		})
	}
}

func TestSigSetRejectsStmt(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal(`Sig("n").Set(js.Stmt) did not panic`)
		}
		msg, ok := r.(string)
		if !ok {
			t.Fatalf("panic value = %#v, want string", r)
		}
		for _, want := range []string{"js.Stmt", "SetExpr", "ds.Do"} {
			if !strings.Contains(msg, want) {
				t.Errorf("panic message %q should mention %q", msg, want)
			}
		}
	}()
	_ = Sig("n").Set(js.Ident("x").Incr())
}

func TestSigSub(t *testing.T) {
	tests := []struct {
		name     string
		got      Sig
		expected Sig
	}{
		{"suffix", Sig("plan").Sub("open"), Sig("plan_open")},
		{"suffix with leading underscore", Sig("plan").Sub("_open"), Sig("plan_open")},
		{"suffix on prefixed signal", Sig("$$plan").Sub("open"), Sig("plan_open")},
		{"chained", Sig("plan").Sub("list").Sub("active"), Sig("plan_list_active")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("Sub = %q, want %q", string(tt.got), string(tt.expected))
			}
			if want := "$$" + string(tt.expected); js.ToJS(tt.got.Ref()) != want {
				t.Errorf("Sub().Ref() = %q, want %q", js.ToJS(tt.got.Ref()), want)
			}
		})
	}
}

func TestSignalRef(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"plain", "items", "$$items"},
		{"double dollar prefixed", "$$items", "$$items"},
		{"single dollar prefixed", "$items", "$$items"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := js.ToJS(SignalRef(tt.input).Expr()); got != tt.expected {
				t.Errorf("SignalRef(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
