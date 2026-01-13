package h

import (
	"bytes"
	"slices"
	"testing"
)

func TestIf(t *testing.T) {
	tests := []struct {
		name     string
		cond     bool
		ifTrue   Builder
		ifElse   Builder
		expected string
	}{
		{"true condition", true, Text("yes"), Text("no"), "yes"},
		{"false condition", false, Text("yes"), Text("no"), "no"},
		{"true with nil else", true, Text("yes"), nil, "yes"},
		{"false with nil else", false, Text("yes"), nil, ""},
		{"true with nil true", true, nil, Text("no"), ""},
		{"both nil true", true, nil, nil, ""},
		{"both nil false", false, nil, nil, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			Render(buf, If(tt.cond, tt.ifTrue, tt.ifElse))
			if buf.String() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, buf.String())
			}
		})
	}
}

func TestIfInContext(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	b := Div(
		If(true,
			Span(Text("visible")),
			Span(Text("hidden")),
		),
	)
	Render(buf, b)
	expected := "<div><span>visible</span></div>"
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

func TestWhen(t *testing.T) {
	tests := []struct {
		name     string
		cond     bool
		ifTrue   Builder
		expected string
	}{
		{"true condition", true, Text("shown"), "shown"},
		{"false condition", false, Text("hidden"), ""},
		{"true with nil", true, nil, ""},
		{"false with nil", false, nil, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			Render(buf, When(tt.cond, tt.ifTrue))
			if buf.String() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, buf.String())
			}
		})
	}
}

func TestWhenInContext(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	b := Div(
		When(true, Span(Text("visible"))),
		When(false, Span(Text("hidden"))),
	)
	Render(buf, b)
	expected := "<div><span>visible</span></div>"
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

func TestFirst(t *testing.T) {
	tests := []struct {
		name     string
		builders []Builder
		expected string
	}{
		{"first non-nil", []Builder{Text("a"), Text("b")}, "a"},
		{"skip leading nils", []Builder{nil, nil, Text("c")}, "c"},
		{"all nil", []Builder{nil, nil, nil}, ""},
		{"empty", []Builder{}, ""},
		{"single", []Builder{Text("only")}, "only"},
		{"single nil", []Builder{nil}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			Render(buf, First(tt.builders...))
			if buf.String() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, buf.String())
			}
		})
	}
}

func TestFirstInContext(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	b := Div(
		First(nil, nil, Span(Text("fallback"))),
	)
	Render(buf, b)
	expected := "<div><span>fallback</span></div>"
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

func TestForEach(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	items := []string{"a", "b", "c"}
	b := Ul(
		ForEach(slices.Values(items), func(s string) Builder {
			return Li(Text(s))
		}),
	)
	Render(buf, b)
	expected := "<ul><li>a</li><li>b</li><li>c</li></ul>"
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

func TestForEachEmpty(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	var items []string
	b := Ul(
		ForEach(slices.Values(items), func(s string) Builder {
			return Li(Text(s))
		}),
	)
	Render(buf, b)
	expected := "<ul></ul>"
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

func TestForEachWithNilReturns(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	items := []string{"a", "", "c"} // empty string returns nil
	b := Ul(
		ForEach(slices.Values(items), func(s string) Builder {
			if s == "" {
				return nil
			}
			return Li(Text(s))
		}),
	)
	Render(buf, b)
	expected := "<ul><li>a</li><li>c</li></ul>"
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

func TestForEach2(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	items := []string{"a", "b", "c"}
	b := Ol(
		ForEach2(slices.All(items), func(i int, s string) Builder {
			return Li(Text(s))
		}),
	)
	Render(buf, b)
	expected := "<ol><li>a</li><li>b</li><li>c</li></ol>"
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

func TestForEach2Empty(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	var items []string
	b := Ol(
		ForEach2(slices.All(items), func(i int, s string) Builder {
			return Li(Text(s))
		}),
	)
	Render(buf, b)
	expected := "<ol></ol>"
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

func TestForEach2WithNilReturns(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	items := []string{"a", "", "c"}
	b := Ol(
		ForEach2(slices.All(items), func(i int, s string) Builder {
			if s == "" {
				return nil
			}
			return Li(Text(s))
		}),
	)
	Render(buf, b)
	expected := "<ol><li>a</li><li>c</li></ol>"
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

// Error propagation tests - using errorWriter from api_test.go
func TestForEachError(t *testing.T) {
	ew := &errorWriter{}
	w := NewWriter(ew)
	items := []string{"a", "b"}
	b := ForEach(slices.Values(items), func(s string) Builder {
		return Text(s)
	})
	err := b.Build(w)
	if err == nil {
		t.Error("expected error from ForEach build")
	}
}

func TestForEach2Error(t *testing.T) {
	ew := &errorWriter{}
	w := NewWriter(ew)
	items := []string{"a", "b"}
	b := ForEach2(slices.All(items), func(i int, s string) Builder {
		return Text(s)
	})
	err := b.Build(w)
	if err == nil {
		t.Error("expected error from ForEach2 build")
	}
}

// Verify builders implement TagArg interface
func TestBuildersAreTagArg(t *testing.T) {
	var _ TagArg = ForEach(slices.Values([]string{}), func(s string) Builder { return nil })
	var _ TagArg = ForEach2(slices.All([]string{}), func(i int, s string) Builder { return nil })
}

// Verify lazy evaluation - iterator is not consumed until Build
func TestForEachLazyEvaluation(t *testing.T) {
	callCount := 0
	seq := func(yield func(int) bool) {
		for i := range 3 {
			callCount++
			if !yield(i) {
				return
			}
		}
	}

	// Create the ForEach builder - should not consume iterator yet
	b := ForEach(seq, func(i int) Builder {
		return Text("x")
	})

	if callCount != 0 {
		t.Errorf("expected iterator not consumed yet, got %d calls", callCount)
	}

	// Now render - should consume iterator
	buf := bytes.NewBuffer(nil)
	Render(buf, b)

	if callCount != 3 {
		t.Errorf("expected 3 iterator calls, got %d", callCount)
	}
}
