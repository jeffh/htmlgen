package ds

import (
	"strings"
	"testing"

	"github.com/jeffh/htmlgen/js"
)

// ============ sig.go tests ============

func TestSigName(t *testing.T) {
	tests := []struct {
		name     string
		sig      Sig
		expected string
	}{
		{"plain", Sig("open"), "open"},
		{"dollar prefixed", Sig("$open"), "open"},
		{"nested path", Sig("form.email"), "form.email"},
		{"only dollar stripped once", Sig("$$open"), "$open"},
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
		value    Value
		expected string
	}{
		{"Value", Sig("open").Value(), "$open"},
		{"Value strips dollar", Sig("$open").Value(), "$open"},
		{"Not", Sig("open").Not(), "!$open"},
		{"Toggle", Sig("open").Toggle(), "($open = !$open)"},
		{"Toggle strips dollar", Sig("$open").Toggle(), "($open = !$open)"},
		{"Clear", Sig("q").Clear(), "($q = '')"},
		{"SetExpr", Sig("n").SetExpr(js.Int(1)), "($n = 1)"},
		{"SetExpr with signal", Sig("n").SetExpr(Sig("m").Ref()), "($n = $m)"},
		{"SetExpr with evt", Sig("q").SetExpr(EvtValue), "($q = evt.target.value)"},
		{"Eq", Sig("tab").Eq(js.String("plans")), `($tab === "plans")`},
		{"NotEq", Sig("tab").NotEq(js.String("plans")), `($tab !== "plans")`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToJS(tt.value.expr); got != tt.expected {
				t.Errorf("%s = %q, want %q", tt.name, got, tt.expected)
			}
		})
	}
}

func TestSigRef(t *testing.T) {
	if got := ToJS(Sig("open").Ref()); got != "$open" {
		t.Errorf(`Sig("open").Ref() = %q, want %q`, got, "$open")
	}
}

func TestSigSet(t *testing.T) {
	tests := []struct {
		name     string
		value    any
		expected string
	}{
		{"string is JSON encoded", "hello", `($n = "hello")`},
		{"int is JSON encoded", 3, "($n = 3)"},
		{"bool is JSON encoded", true, "($n = true)"},
		{"nil is JSON encoded", nil, "($n = null)"},
		{"slice is JSON encoded", []int{1, 2}, "($n = [1,2])"},
		{"js.Expr passes through", js.Ident("x"), "($n = x)"},
		{"Value passes through", Sig("m").Value(), "($n = $m)"},
		{"Value expression passes through", Sig("m").Not(), "($n = !$m)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToJS(Sig("n").Set(tt.value).expr)
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
		{"dollar prefixed base", Sig("$plan").Sub("q"), Sig("plan_q")},
		{"chained", Sig("plan").Sub("open").Sub("above"), Sig("plan_open_above")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("Sub = %q, want %q", string(tt.got), string(tt.expected))
			}
		})
	}
}

func TestSigInAttributes(t *testing.T) {
	s := Sig("plan")
	tests := []struct {
		name     string
		attrName string
		attrVal  string
		expName  string
		expVal   string
	}{
		{
			"Show",
			Show(s.Sub("open").Value()).Name, Show(s.Sub("open").Value()).Value,
			"data-show", "$plan_open",
		},
		{
			"Text",
			Text(s.Value()).Name, Text(s.Value()).Value,
			"data-text", "$plan",
		},
		{
			"OnClick toggle and clear",
			OnClick(s.Sub("open").Toggle(), s.Sub("q").Clear()).Attribute().Name,
			OnClick(s.Sub("open").Toggle(), s.Sub("q").Clear()).Attribute().Value,
			"data-on:click", "($plan_open = !$plan_open); ($plan_q = '')",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.attrName != tt.expName {
				t.Errorf("name = %q, want %q", tt.attrName, tt.expName)
			}
			if tt.attrVal != tt.expVal {
				t.Errorf("value = %q, want %q", tt.attrVal, tt.expVal)
			}
		})
	}
}

func TestSignalRefMatchesSig(t *testing.T) {
	for _, name := range []string{"open", "$open", "form.email"} {
		want := ToJS(Sig(name).Value().expr)
		if got := ToJS(SignalRef(name).expr); got != want {
			t.Errorf("SignalRef(%q) = %q, want %q (Sig parity)", name, got, want)
		}
	}
}

// ============ Do / combinator tests ============

func TestDo(t *testing.T) {
	tests := []struct {
		name     string
		stmts    []js.Stmt
		expected string
	}{
		{"empty", nil, ""},
		{"single", []js.Stmt{js.Ident("x").Incr()}, "x++"},
		{
			"multiple joined with semicolons",
			[]js.Stmt{js.Let("n", js.Int(1)), js.Ident("n").Incr()},
			"let n = 1; n++",
		},
		{
			"try/catch",
			[]js.Stmt{js.Try(js.ExprStmt(js.Ident("el").Method("focus")))},
			"try { el.focus() } catch {}",
		},
		{
			"if statement",
			[]js.Stmt{js.If(Sig("open").Ref(), js.Ident("x").Incr())},
			"if ($open) { x++ }",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToJS(Do(tt.stmts...).expr); got != tt.expected {
				t.Errorf("Do(...) = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestDoInEventHandler(t *testing.T) {
	attr := OnClick(Do(js.Let("n", js.Int(1)), js.Ident("n").Incr())).Attribute()
	if attr.Name != "data-on:click" {
		t.Errorf("name = %q, want %q", attr.Name, "data-on:click")
	}
	if attr.Value != "let n = 1; n++" {
		t.Errorf("value = %q, want %q", attr.Value, "let n = 1; n++")
	}
}

func TestValueCombinators(t *testing.T) {
	a := Sig("a").Value()
	b := Sig("b").Value()

	tests := []struct {
		name     string
		value    Value
		expected string
	}{
		{"Not", a.Not(), "!$a"},
		{"Not twice", a.Not().Not(), "!!$a"},
		{"And", a.And(b), "($a && $b)"},
		{"Or", a.Or(b), "($a || $b)"},
		{"And then Not", a.And(b).Not(), "!($a && $b)"},
		{"Or of Nots", a.Not().Or(b.Not()), "(!$a || !$b)"},
		{"Ternary", a.Ternary(Str("yes"), Str("no")), `($a ? "yes" : "no")`},
		{"Ternary nested", a.Ternary(b.Ternary(Str("x"), Str("y")), Str("z")), `($a ? ($b ? "x" : "y") : "z")`},
		{"chained And", a.And(b).And(Sig("c").Value()), "(($a && $b) && $c)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToJS(tt.value.expr); got != tt.expected {
				t.Errorf("%s = %q, want %q", tt.name, got, tt.expected)
			}
		})
	}
}

func TestValueCombinatorsAreImmutable(t *testing.T) {
	a := Sig("a").Value()
	_ = a.Not()
	_ = a.And(Sig("b").Value())
	if got := ToJS(a.expr); got != "$a" {
		t.Errorf("original Value mutated: %q, want %q", got, "$a")
	}
}

// ============ Confirm tests ============

func TestConfirm(t *testing.T) {
	tests := []struct {
		name     string
		value    Value
		expected string
	}{
		{"no actions", Confirm("Sure?"), `confirm("Sure?")`},
		{
			"one action",
			Confirm("Delete this plan?", Delete("/plans/1")),
			`(confirm("Delete this plan?") && @delete("/plans/1"))`,
		},
		{
			"two actions fold left",
			Confirm("Ok?", Sig("a").Value(), Sig("b").Value()),
			`((confirm("Ok?") && $a) && $b)`,
		},
		{
			"message is escaped",
			Confirm(`say "hi"`),
			`confirm("say \"hi\"")`,
		},
		{
			"with signal assignment",
			Confirm("Reset?", Sig("q").Clear()),
			`(confirm("Reset?") && ($q = ''))`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToJS(tt.value.expr); got != tt.expected {
				t.Errorf("%s = %q, want %q", tt.name, got, tt.expected)
			}
		})
	}
}

func TestConfirmInEventHandler(t *testing.T) {
	attr := OnClick(Confirm("Delete?", Delete("/x"))).PreventDefault().Attribute()
	if attr.Name != "data-on:click__prevent" {
		t.Errorf("name = %q, want %q", attr.Name, "data-on:click__prevent")
	}
	want := `(confirm("Delete?") && @delete("/x"))`
	if attr.Value != want {
		t.Errorf("value = %q, want %q", attr.Value, want)
	}
}

// ============ ClassesExpr / StylesExpr / AttrsExpr tests ============

func TestClassesExpr(t *testing.T) {
	tests := []struct {
		name     string
		classes  map[string]Value
		expected string
	}{
		{"nil", nil, "{}"},
		{"empty", map[string]Value{}, "{}"},
		{
			"single expression value",
			map[string]Value{"hidden": Sig("open").Not()},
			`{"hidden": !$open}`,
		},
		{
			"keys are sorted",
			map[string]Value{
				"hidden": Sig("open").Not(),
				"active": Sig("open").Value(),
				"busy":   Sig("loading").Value(),
				"zebra":  Sig("z").Value(),
				"alpha":  Sig("a").Value(),
			},
			`{"active": $open, "alpha": $a, "busy": $loading, "hidden": !$open, "zebra": $z}`,
		},
		{
			"tailwind-style keys with dashes and slashes",
			map[string]Value{"bg-red-500/50": Sig("hot").Value()},
			`{"bg-red-500/50": $hot}`,
		},
		{
			"comparison values",
			map[string]Value{"selected": Sig("tab").Eq(js.String("plans"))},
			`{"selected": ($tab === "plans")}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attr := ClassesExpr(tt.classes)
			if attr.Name != "data-class" {
				t.Errorf("name = %q, want %q", attr.Name, "data-class")
			}
			if attr.Value != tt.expected {
				t.Errorf("value = %q, want %q", attr.Value, tt.expected)
			}
		})
	}
}

func TestClassesExprIsDeterministic(t *testing.T) {
	classes := map[string]Value{
		"a": Sig("a").Value(),
		"b": Sig("b").Value(),
		"c": Sig("c").Value(),
		"d": Sig("d").Value(),
		"e": Sig("e").Value(),
		"f": Sig("f").Value(),
		"g": Sig("g").Value(),
		"h": Sig("h").Value(),
	}
	want := ClassesExpr(classes).Value
	for i := 0; i < 100; i++ {
		if got := ClassesExpr(classes).Value; got != want {
			t.Fatalf("iteration %d: ClassesExpr = %q, want %q", i, got, want)
		}
	}
	if want != `{"a": $a, "b": $b, "c": $c, "d": $d, "e": $e, "f": $f, "g": $g, "h": $h}` {
		t.Errorf("unexpected sorted output: %q", want)
	}
}

func TestStylesExpr(t *testing.T) {
	tests := []struct {
		name     string
		styles   map[string]Value
		expected string
	}{
		{"nil", nil, "{}"},
		{"single", map[string]Value{"width": Sig("pct").Value()}, `{"width": $pct}`},
		{
			"keys are sorted",
			map[string]Value{"width": Sig("w").Value(), "height": Sig("h").Value()},
			`{"height": $h, "width": $w}`,
		},
		{
			"template expression",
			map[string]Value{"width": V(js.Template("", Sig("pct").Ref(), "%"))},
			"{\"width\": `${$pct}%`}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attr := StylesExpr(tt.styles)
			if attr.Name != "data-style" {
				t.Errorf("name = %q, want %q", attr.Name, "data-style")
			}
			if attr.Value != tt.expected {
				t.Errorf("value = %q, want %q", attr.Value, tt.expected)
			}
		})
	}
}

func TestAttrsExpr(t *testing.T) {
	tests := []struct {
		name     string
		attrs    map[string]Value
		expected string
	}{
		{"nil", nil, "{}"},
		{"single", map[string]Value{"disabled": Sig("busy").Value()}, `{"disabled": $busy}`},
		{
			"keys are sorted",
			map[string]Value{
				"disabled":      Sig("busy").Value(),
				"aria-expanded": Sig("open").Value(),
			},
			`{"aria-expanded": $open, "disabled": $busy}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attr := AttrsExpr(tt.attrs)
			if attr.Name != "data-attr" {
				t.Errorf("name = %q, want %q", attr.Name, "data-attr")
			}
			if attr.Value != tt.expected {
				t.Errorf("value = %q, want %q", attr.Value, tt.expected)
			}
		})
	}
}

// ============ Datastar scope identifier tests ============

func TestDatastarScopeIdentifiers(t *testing.T) {
	tests := []struct {
		name     string
		expr     js.Expr
		expected string
	}{
		{"Evt", Evt, "evt"},
		{"El", El, "el"},
		{"EvtTarget", EvtTarget, "evt.target"},
		{"EvtValue", EvtValue, "evt.target.value"},
		{"EvtKey", EvtKey, "evt.key"},
		{"El method call", El.Method("focus"), "el.focus()"},
		{"El property", El.Prop("dataset").Prop("id"), "el.dataset.id"},
		{"legacy Event still emits event", Event, "event"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToJS(tt.expr); got != tt.expected {
				t.Errorf("%s = %q, want %q", tt.name, got, tt.expected)
			}
		})
	}
}

func TestEvtInEventHandler(t *testing.T) {
	attr := On("keydown", Sig("k").SetExpr(EvtKey)).Attribute()
	if attr.Name != "data-on:keydown" {
		t.Errorf("name = %q, want %q", attr.Name, "data-on:keydown")
	}
	if attr.Value != "($k = evt.key)" {
		t.Errorf("value = %q, want %q", attr.Value, "($k = evt.key)")
	}
}

// ============ Headers determinism tests ============

func TestHeadersSorted(t *testing.T) {
	headers := map[string]string{
		"X-Csrf":        "tok",
		"Accept":        "text/event-stream",
		"X-Request-Id":  "1",
		"Authorization": "Bearer x",
		"Content-Type":  "application/json",
	}
	want := `@get("/api", {headers: {"Accept": "text/event-stream", "Authorization": "Bearer x", ` +
		`"Content-Type": "application/json", "X-Csrf": "tok", "X-Request-Id": "1"}})`

	for i := 0; i < 100; i++ {
		got := ToJS(GetWithOptions("/api", RequestOptions().Headers(headers)).expr)
		if got != want {
			t.Fatalf("iteration %d: Headers = %q, want %q", i, got, want)
		}
	}
}

func TestHeadersEmpty(t *testing.T) {
	got := ToJS(GetWithOptions("/api", RequestOptions().Headers(nil)).expr)
	want := `@get("/api", {headers: {}})`
	if got != want {
		t.Errorf("Headers(nil) = %q, want %q", got, want)
	}
}

func TestSigEmptyNamePanics(t *testing.T) {
	for _, name := range []string{"", "$"} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("Sig(%q).Ref() did not panic", name)
				}
			}()
			Sig(name).Ref()
		}()
	}
}
