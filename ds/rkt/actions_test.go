package rkt

import (
	"strings"
	"testing"

	"github.com/jeffh/htmlgen/ds"
	"github.com/jeffh/htmlgen/js"
)

func TestAction(t *testing.T) {
	tests := []struct {
		name     string
		value    ds.Value
		expected string
	}{
		{"no args", Action("increment"), "@increment()"},
		{"one arg", Action("add", js.Int(1)), "@add(1)"},
		{"several args", Action("move", js.Int(1), js.String("up")), `@move(1, "up")`},
		{"signal arg", Action("set", Sig("count").Ref()), "@set($$count)"},
		{"underscore name", Action("_private"), "@_private()"},
		{"dollar name", Action("$fn"), "@$fn()"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := js.ToJS(tt.value.Expr()); got != tt.expected {
				t.Errorf("%s = %q, want %q", tt.name, got, tt.expected)
			}
		})
	}
}

func TestDispatch(t *testing.T) {
	tests := []struct {
		name     string
		value    ds.Value
		expected string
	}{
		{"no args", Dispatch("increment"), `@dispatchRocket("increment")`},
		{"one arg", Dispatch("increment", js.Int(1)), `@dispatchRocket("increment", 1)`},
		{"several args", Dispatch("move", js.Int(1), js.String("up")), `@dispatchRocket("move", 1, "up")`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := js.ToJS(tt.value.Expr()); got != tt.expected {
				t.Errorf("%s = %q, want %q", tt.name, got, tt.expected)
			}
		})
	}
}

func TestActionRejectsInvalidNames(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty", ""},
		{"leading digit", "1up"},
		{"hyphen", "do-thing"},
		{"call syntax", "increment()"},
		{"space", "do thing"},
		{"at prefix", "@increment"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run("Action", func(t *testing.T) {
				defer expectPanicContaining(t, "invalid action name")
				_ = Action(tt.input)
			})
			t.Run("Dispatch", func(t *testing.T) {
				defer expectPanicContaining(t, "invalid action name")
				_ = Dispatch(tt.input)
			})
		})
	}
}

// expectPanicContaining is deferred by a test that expects a panic whose
// message mentions want.
func expectPanicContaining(t *testing.T, want string) {
	t.Helper()
	r := recover()
	if r == nil {
		t.Fatalf("expected a panic mentioning %q, got none", want)
	}
	msg, ok := r.(string)
	if !ok {
		t.Fatalf("panic value = %#v, want string", r)
	}
	if !strings.Contains(msg, want) {
		t.Errorf("panic message %q should mention %q", msg, want)
	}
}

func TestActionComposesWithDs(t *testing.T) {
	attr := ds.OnClick(Action("increment")).PreventDefault().Attribute()
	if attr.Name != "data-on:click__prevent" {
		t.Errorf("Name = %q, want %q", attr.Name, "data-on:click__prevent")
	}
	if attr.Value != "@increment()" {
		t.Errorf("Value = %q, want %q", attr.Value, "@increment()")
	}
}
