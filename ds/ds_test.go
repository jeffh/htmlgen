package ds

import (
	"strings"
	"testing"
	"time"

	"github.com/jeffh/htmlgen/js"
)

// ============ values.go tests ============

func TestRaw(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple expression", "$foo", "$foo"},
		{"complex expression", "$a + $b", "$a + $b"},
		{"function call", "console.log($x)", "console.log($x)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Raw(tt.input)

			// Test via ToJS
			if got := ToJS(v.expr); got != tt.expected {
				t.Errorf("ToJS(Raw(%q)) = %q, want %q", tt.input, got, tt.expected)
			}

			// Test via OnClick attribute value
			attr := OnClick(v).Attribute()
			if attr.Value != tt.expected {
				t.Errorf("OnClick(Raw(%q)).Value = %q, want %q", tt.input, attr.Value, tt.expected)
			}
			if attr.Name != "data-on:click" {
				t.Errorf("OnClick(Raw(%q)).Name = %q, want %q", tt.input, attr.Name, "data-on:click")
			}
		})
	}
}

func TestStr(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple string", "hello", `"hello"`},
		{"string with quotes", `say "hi"`, `"say \"hi\""`},
		{"empty string", "", `""`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Str(tt.input)
			if got := ToJS(v.expr); got != tt.expected {
				t.Errorf("Str(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestJsonValue(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{"string", "hello", `"hello"`},
		{"number", 42, "42"},
		{"float", 3.14, "3.14"},
		{"bool true", true, "true"},
		{"bool false", false, "false"},
		{"nil", nil, "null"},
		{"slice", []int{1, 2, 3}, "[1,2,3]"},
		{"map", map[string]int{"a": 1}, `{"a":1}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := JsonValue(tt.input)
			if got := ToJS(v.expr); got != tt.expected {
				t.Errorf("JsonValue(%v) = %q, want %q", tt.input, got, tt.expected)
			}

			// Also test via OnClick attribute value
			attr := OnClick(v).Attribute()
			if attr.Value != tt.expected {
				t.Errorf("OnClick(JsonValue(%v)).Value = %q, want %q", tt.input, attr.Value, tt.expected)
			}
		})
	}
}

func TestJsonValuePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for unmarshalable value")
		}
	}()

	// channels cannot be marshaled to JSON
	JsonValue(make(chan int))
}

func TestNavigate(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		values   []any
		expected string
	}{
		{"simple path", "/home", nil, `window.location.href = "/home"`},
		{"path with format", "/users/%d", []any{42}, `window.location.href = "/users/42"`},
		{"path with string format", "/users/%s/edit", []any{"john"}, `window.location.href = "/users/john/edit"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Navigate(tt.path, tt.values...)
			if got != tt.expected {
				t.Errorf("Navigate(%q, %v) = %q, want %q", tt.path, tt.values, got, tt.expected)
			}
		})
	}
}

func TestOnSuccess(t *testing.T) {
	action := OnSuccess(Raw("console.log('done')"))
	var sb strings.Builder
	action.appendChain(&sb)
	expected := ".then(() => console.log('done'))"
	if got := sb.String(); got != expected {
		t.Errorf("OnSuccess() = %q, want %q", got, expected)
	}
}

func TestOnFailure(t *testing.T) {
	action := OnFailure(Raw("console.log(error)"))
	var sb strings.Builder
	action.appendChain(&sb)
	expected := ".catch((error) => console.log(error))"
	if got := sb.String(); got != expected {
		t.Errorf("OnFailure() = %q, want %q", got, expected)
	}
}

func TestConsoleLog(t *testing.T) {
	tests := []struct {
		name     string
		values   []js.Expr
		expected string
	}{
		{"single value", []js.Expr{js.Raw("$foo")}, "console.log($foo)"},
		{"multiple values", []js.Expr{js.Raw("$a"), js.Raw("$b")}, "console.log($a, $b)"},
		{"with string", []js.Expr{js.String("msg"), js.Raw("$x")}, `console.log("msg", $x)`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action := ConsoleLog(tt.values...)
			attr := OnClick(action).Attribute()
			if attr.Value != tt.expected {
				t.Errorf("ConsoleLog() = %q, want %q", attr.Value, tt.expected)
			}
		})
	}
}

func TestAnd(t *testing.T) {
	tests := []struct {
		name     string
		actions  []js.Expr
		expected string
	}{
		{"two actions", []js.Expr{js.Raw("$a"), js.Raw("$b")}, "($a && $b)"},
		{"three actions", []js.Expr{js.Raw("$a"), js.Raw("$b"), js.Raw("$c")}, "(($a && $b) && $c)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := And(tt.actions...)

			// Test via ToJS
			if got := ToJS(v); got != tt.expected {
				t.Errorf("And() = %q, want %q", got, tt.expected)
			}

			// Test as a Value via OnClick
			attr := OnClick(AndValue(tt.actions...)).Attribute()
			if attr.Value != tt.expected {
				t.Errorf("OnClick(AndValue(...)).Value = %q, want %q", attr.Value, tt.expected)
			}
		})
	}
}

// ============ modifiers.go tests ============

func TestPreventDefault(t *testing.T) {
	attr := OnClick().PreventDefault().Attribute()
	if !strings.Contains(attr.Name, "__prevent") {
		t.Errorf("PreventDefault() should add __prevent, got %q", attr.Name)
	}
}

func TestOnce(t *testing.T) {
	attr := OnClick().Once().Attribute()
	if !strings.Contains(attr.Name, "__once") {
		t.Errorf("Once() should add __once, got %q", attr.Name)
	}
}

func TestPassive(t *testing.T) {
	attr := On("scroll").Passive().Attribute()
	if !strings.Contains(attr.Name, "__passive") {
		t.Errorf("Passive() should add __passive, got %q", attr.Name)
	}
}

func TestCapture(t *testing.T) {
	attr := OnClick().Capture().Attribute()
	if !strings.Contains(attr.Name, "__capture") {
		t.Errorf("Capture() should add __capture, got %q", attr.Name)
	}
}

func TestCase(t *testing.T) {
	tests := []struct {
		casing   SignalCasing
		expected string
	}{
		{CamelCase, "__case.camel"},
		{KebabCase, "__case.kebab"},
		{SnakeCase, "__case.snake"},
		{PascalCase, "__case.pascal"},
	}

	for _, tt := range tests {
		t.Run(string(tt.casing), func(t *testing.T) {
			attr := Signals(map[string]any{}).Case(tt.casing).Attribute()
			if !strings.Contains(attr.Name, tt.expected) {
				t.Errorf("Case(%s) should add %s, got %q", tt.casing, tt.expected, attr.Name)
			}
		})
	}
}

func TestDelay(t *testing.T) {
	attr := OnClick().Delay(500 * time.Millisecond).Attribute()
	if !strings.Contains(attr.Name, "__delay.500ms") {
		t.Errorf("Delay() should add __delay.500ms, got %q", attr.Name)
	}
}

func TestDebounce(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		opts     []TimingOption
		expected string
	}{
		{"basic", 300 * time.Millisecond, nil, "__debounce.300ms"},
		{"with leading", 500 * time.Millisecond, []TimingOption{Leading()}, "__debounce.500ms.leading"},
		{"with notrailing", 500 * time.Millisecond, []TimingOption{NoTrailing()}, "__debounce.500ms.notrailing"},
		{"with both", 500 * time.Millisecond, []TimingOption{Leading(), NoTrailing()}, "__debounce.500ms.leading.notrailing"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attr := OnInput().Debounce(tt.duration, tt.opts...).Attribute()
			if !strings.Contains(attr.Name, tt.expected) {
				t.Errorf("Debounce() should contain %s, got %q", tt.expected, attr.Name)
			}
		})
	}
}

func TestThrottle(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		opts     []TimingOption
		expected string
	}{
		{"basic", 100 * time.Millisecond, nil, "__throttle.100ms"},
		{"with noleading", 200 * time.Millisecond, []TimingOption{NoLeading()}, "__throttle.200ms.noleading"},
		{"with trailing", 200 * time.Millisecond, []TimingOption{Trailing()}, "__throttle.200ms.trailing"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attr := On("scroll").Throttle(tt.duration, tt.opts...).Attribute()
			if !strings.Contains(attr.Name, tt.expected) {
				t.Errorf("Throttle() should contain %s, got %q", tt.expected, attr.Name)
			}
		})
	}
}

func TestViewTransition(t *testing.T) {
	attr := OnClick().ViewTransition().Attribute()
	if !strings.Contains(attr.Name, "__viewtransition") {
		t.Errorf("ViewTransition() should add __viewtransition, got %q", attr.Name)
	}
}

func TestWindow(t *testing.T) {
	attr := On("scroll").Window().Attribute()
	if !strings.Contains(attr.Name, "__window") {
		t.Errorf("Window() should add __window, got %q", attr.Name)
	}
}

func TestOutside(t *testing.T) {
	attr := OnClick().Outside().Attribute()
	if !strings.Contains(attr.Name, "__outside") {
		t.Errorf("Outside() should add __outside, got %q", attr.Name)
	}
}

func TestStopPropagation(t *testing.T) {
	attr := OnClick().StopPropagation().Attribute()
	if !strings.Contains(attr.Name, "__stop") {
		t.Errorf("StopPropagation() should add __stop, got %q", attr.Name)
	}
}

func TestIfMissing(t *testing.T) {
	attr := Signals(map[string]any{}).IfMissing().Attribute()
	if !strings.Contains(attr.Name, "__ifmissing") {
		t.Errorf("IfMissing() should add __ifmissing, got %q", attr.Name)
	}
}

func TestDocument(t *testing.T) {
	attr := On("visibilitychange").Document().Attribute()
	if !strings.Contains(attr.Name, "__document") {
		t.Errorf("Document() should add __document, got %q", attr.Name)
	}
}

func TestEventCase(t *testing.T) {
	attr := On("myEvent").Case(KebabCase).Attribute()
	if !strings.Contains(attr.Name, "__case.kebab") {
		t.Errorf("Case() should add __case.kebab, got %q", attr.Name)
	}
}

func TestTerse(t *testing.T) {
	attr := JsonSignalsDebug(nil).Terse().Attribute()
	if !strings.Contains(attr.Name, "__terse") {
		t.Errorf("Terse() should add __terse, got %q", attr.Name)
	}
}

func TestHalf(t *testing.T) {
	attr := OnIntersect().Half().Attribute()
	if !strings.Contains(attr.Name, "__half") {
		t.Errorf("Half() should add __half, got %q", attr.Name)
	}
}

func TestFull(t *testing.T) {
	attr := OnIntersect().Full().Attribute()
	if !strings.Contains(attr.Name, "__full") {
		t.Errorf("Full() should add __full, got %q", attr.Name)
	}
}

func TestDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		opts     []TimingOption
		expected string
	}{
		{"basic", 500 * time.Millisecond, nil, "__duration.500ms"},
		{"with leading", 2 * time.Second, []TimingOption{DurationLeading()}, "__duration.2s.leading"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attr := OnInterval().Duration(tt.duration, tt.opts...).Attribute()
			if !strings.Contains(attr.Name, tt.expected) {
				t.Errorf("Duration() should contain %s, got %q", tt.expected, attr.Name)
			}
		})
	}
}

func TestExit(t *testing.T) {
	attr := OnIntersect().Exit().Attribute()
	if !strings.Contains(attr.Name, "__exit") {
		t.Errorf("Exit() should add __exit, got %q", attr.Name)
	}
}

func TestThreshold(t *testing.T) {
	tests := []struct {
		name     string
		percent  int
		expected string
	}{
		{"quarter", 25, "__threshold.25"},
		{"three quarters", 75, "__threshold.75"},
		{"full", 100, "__threshold.100"},
		{"zero", 0, "__threshold.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attr := OnIntersect().Threshold(tt.percent).Attribute()
			if !strings.Contains(attr.Name, tt.expected) {
				t.Errorf("Threshold(%v) should contain %s, got %q", tt.percent, tt.expected, attr.Name)
			}
		})
	}
}

func TestIntersectTiming(t *testing.T) {
	attr := OnIntersect().Debounce(300 * time.Millisecond).Attribute()
	if !strings.Contains(attr.Name, "__debounce.300ms") {
		t.Errorf("OnIntersect().Debounce() should add __debounce.300ms, got %q", attr.Name)
	}

	attr = OnIntersect().Throttle(100 * time.Millisecond).Delay(1 * time.Second).ViewTransition().Attribute()
	for _, want := range []string{"__throttle.100ms", "__delay.1s", "__viewtransition"} {
		if !strings.Contains(attr.Name, want) {
			t.Errorf("OnIntersect() modifiers should contain %s, got %q", want, attr.Name)
		}
	}
}

func TestIntervalViewTransition(t *testing.T) {
	attr := OnInterval().ViewTransition().Attribute()
	if !strings.Contains(attr.Name, "__viewtransition") {
		t.Errorf("OnInterval().ViewTransition() should add __viewtransition, got %q", attr.Name)
	}
}

func TestSignalPatchTiming(t *testing.T) {
	attr := OnSignalPatch(Raw("$x++")).Debounce(500 * time.Millisecond).Attribute()
	if !strings.Contains(attr.Name, "__debounce.500ms") {
		t.Errorf("OnSignalPatch().Debounce() should add __debounce.500ms, got %q", attr.Name)
	}
}

// ============ attrs.go tests ============

func TestSetSignalExpr(t *testing.T) {
	tests := []struct {
		name       string
		signalName string
		expr       js.Expr
		expected   string
	}{
		{"without $", "foo", js.Raw("1"), "$foo = 1"},
		{"with $", "$bar", js.Raw("true"), "$bar = true"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action := SetSignalExpr(tt.signalName, tt.expr)
			attr := OnClick(action).Attribute()
			if attr.Value != tt.expected {
				t.Errorf("SetSignalExpr() = %q, want %q", attr.Value, tt.expected)
			}
		})
	}
}

func TestSetSignal(t *testing.T) {
	tests := []struct {
		name       string
		signalName string
		value      any
		expected   string
	}{
		{"string value", "msg", "hello", `$msg = "hello"`},
		{"number value", "count", 42, "$count = 42"},
		{"bool value", "active", true, "$active = true"},
		{"raw expression", "expr", Raw("$a + $b"), "$expr = $a + $b"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action := SetSignal(tt.signalName, tt.value)
			attr := OnClick(action).Attribute()
			if attr.Value != tt.expected {
				t.Errorf("SetSignal() = %q, want %q", attr.Value, tt.expected)
			}
		})
	}
}

func TestOnSubmit(t *testing.T) {
	attr := OnSubmit(Raw("$submit()")).Attribute()
	if attr.Name != "data-on:submit" {
		t.Errorf("OnSubmit().Name = %q, want %q", attr.Name, "data-on:submit")
	}
	if attr.Value != "$submit()" {
		t.Errorf("OnSubmit().Value = %q, want %q", attr.Value, "$submit()")
	}
}

func TestOnInput(t *testing.T) {
	attr := OnInput(Raw("$validate()")).Attribute()
	if attr.Name != "data-on:input" {
		t.Errorf("OnInput().Name = %q, want %q", attr.Name, "data-on:input")
	}
}

func TestOnChange(t *testing.T) {
	attr := OnChange(Raw("$update()")).Attribute()
	if attr.Name != "data-on:change" {
		t.Errorf("OnChange().Name = %q, want %q", attr.Name, "data-on:change")
	}
}

func TestOnClick(t *testing.T) {
	attr := OnClick(Raw("$count++")).Attribute()
	if attr.Name != "data-on:click" {
		t.Errorf("OnClick().Name = %q, want %q", attr.Name, "data-on:click")
	}
	if attr.Value != "$count++" {
		t.Errorf("OnClick().Value = %q, want %q", attr.Value, "$count++")
	}
}

func TestOnLoad(t *testing.T) {
	attr := OnLoad(Raw("$init()")).Attribute()
	if attr.Name != "data-init" {
		t.Errorf("OnLoad().Name = %q, want %q", attr.Name, "data-init")
	}
}

func TestOn(t *testing.T) {
	attr := On("keydown", Raw("$handleKey(event)")).Attribute()
	if attr.Name != "data-on:keydown" {
		t.Errorf("On().Name = %q, want %q", attr.Name, "data-on:keydown")
	}
	if attr.Value != "$handleKey(event)" {
		t.Errorf("On().Value = %q, want %q", attr.Value, "$handleKey(event)")
	}
}

func TestOnIntersect(t *testing.T) {
	attr := OnIntersect(Raw("$seen = true")).Once().Attribute()
	if !strings.HasPrefix(attr.Name, "data-on-intersect") {
		t.Errorf("OnIntersect().Name should start with data-on-intersect, got %q", attr.Name)
	}
	if !strings.Contains(attr.Name, "__once") {
		t.Errorf("OnIntersect().Name should contain __once, got %q", attr.Name)
	}
}

func TestOnInterval(t *testing.T) {
	attr := OnInterval(Raw("$tick()")).Duration(500 * time.Millisecond).Attribute()
	if !strings.HasPrefix(attr.Name, "data-on-interval") {
		t.Errorf("OnInterval().Name should start with data-on-interval, got %q", attr.Name)
	}
	if !strings.Contains(attr.Name, "__duration.500ms") {
		t.Errorf("OnInterval().Name should contain __duration.500ms, got %q", attr.Name)
	}
}

func TestOnSignalPatch(t *testing.T) {
	attr := OnSignalPatch(Raw("console.log('changed')")).Attribute()
	if attr.Name != "data-on-signal-patch" {
		t.Errorf("OnSignalPatch().Name = %q, want %q", attr.Name, "data-on-signal-patch")
	}
}

func TestOnSignalPatchFilter(t *testing.T) {
	pattern := "^user"
	attr := OnSignalPatchFilter(&FilterOptions{IncludeReg: &pattern})
	if attr.Name != "data-on-signal-patch-filter" {
		t.Errorf("OnSignalPatchFilter().Name = %q, want %q", attr.Name, "data-on-signal-patch-filter")
	}
	if !strings.Contains(attr.Value, "include: /^user/") {
		t.Errorf("OnSignalPatchFilter().Value should contain include regex, got %q", attr.Value)
	}

	// Test nil options
	attr = OnSignalPatchFilter(nil)
	if attr.Value != "" {
		t.Errorf("OnSignalPatchFilter(nil).Value = %q, want empty", attr.Value)
	}
}

func TestSignalExpr(t *testing.T) {
	attr := SignalExpr("count", Raw("0")).Attribute()
	if attr.Name != "data-signals:count" {
		t.Errorf("SignalExpr().Name = %q, want %q", attr.Name, "data-signals:count")
	}
	if attr.Value != "0" {
		t.Errorf("SignalExpr().Value = %q, want %q", attr.Value, "0")
	}
}

func TestSignal(t *testing.T) {
	attr := Signal("count", 0).Attribute()
	if attr.Name != "data-signals:count" {
		t.Errorf("Signal().Name = %q, want %q", attr.Name, "data-signals:count")
	}
	if attr.Value != "0" {
		t.Errorf("Signal().Value = %q, want %q", attr.Value, "0")
	}

	attr = Signal("name", "hello").Attribute()
	if attr.Value != `"hello"` {
		t.Errorf("Signal().Value = %q, want %q", attr.Value, `"hello"`)
	}
}

func TestSignals(t *testing.T) {
	attr := Signals(map[string]any{"foo": 1, "bar": "hello"}).Attribute()
	if attr.Name != "data-signals" {
		t.Errorf("Signals().Name = %q, want %q", attr.Name, "data-signals")
	}
	// JSON encoding of map, order may vary
	if !strings.Contains(attr.Value, `"foo":1`) || !strings.Contains(attr.Value, `"bar":"hello"`) {
		t.Errorf("Signals().Value = %q, should contain foo:1 and bar:hello", attr.Value)
	}
}

func TestBind(t *testing.T) {
	attr := Bind("username").Attribute()
	if attr.Name != "data-bind" {
		t.Errorf("Bind().Name = %q, want %q", attr.Name, "data-bind")
	}
	if attr.Value != "username" {
		t.Errorf("Bind().Value = %q, want %q", attr.Value, "username")
	}
}

func TestBindProp(t *testing.T) {
	attr := BindKey("is-checked").Prop("checked").Attribute()
	if attr.Name != "data-bind:is-checked__prop.checked" {
		t.Errorf("BindKey().Prop().Name = %q, want %q", attr.Name, "data-bind:is-checked__prop.checked")
	}
}

func TestBindEvent(t *testing.T) {
	attr := BindKey("query").Event("input", "change").Attribute()
	if attr.Name != "data-bind:query__event.input.change" {
		t.Errorf("BindKey().Event().Name = %q, want %q", attr.Name, "data-bind:query__event.input.change")
	}
}

func TestClass(t *testing.T) {
	attr := Class("active", Raw("$isActive")).Attribute()
	if attr.Name != "data-class:active" {
		t.Errorf("Class().Name = %q, want %q", attr.Name, "data-class:active")
	}
	if attr.Value != "$isActive" {
		t.Errorf("Class().Value = %q, want %q", attr.Value, "$isActive")
	}

	attr = Class("font-bold", Raw("$bold")).Case(KebabCase).Attribute()
	if attr.Name != "data-class:font-bold__case.kebab" {
		t.Errorf("Class().Case().Name = %q, want %q", attr.Name, "data-class:font-bold__case.kebab")
	}
}

func TestText(t *testing.T) {
	attr := Text(Raw("$message"))
	if attr.Name != "data-text" {
		t.Errorf("Text().Name = %q, want %q", attr.Name, "data-text")
	}
	if attr.Value != "$message" {
		t.Errorf("Text().Value = %q, want %q", attr.Value, "$message")
	}
}

func TestShow(t *testing.T) {
	attr := Show(Raw("$visible"))
	if attr.Name != "data-show" {
		t.Errorf("Show().Name = %q, want %q", attr.Name, "data-show")
	}
	if attr.Value != "$visible" {
		t.Errorf("Show().Value = %q, want %q", attr.Value, "$visible")
	}
}

func TestHide(t *testing.T) {
	attr := Hide()
	if attr.Name != "style" {
		t.Errorf("Hide().Name = %q, want %q", attr.Name, "style")
	}
	if attr.Value != "display: none" {
		t.Errorf("Hide().Value = %q, want %q", attr.Value, "display: none")
	}
}

func TestAttribute(t *testing.T) {
	attr := Attribute("title", Raw("$tooltip"))
	if attr.Name != "data-attr:title" {
		t.Errorf("Attribute().Name = %q, want %q", attr.Name, "data-attr:title")
	}
	if attr.Value != "$tooltip" {
		t.Errorf("Attribute().Value = %q, want %q", attr.Value, "$tooltip")
	}
}

func TestIndicator(t *testing.T) {
	attr := Indicator("loading")
	if attr.Name != "data-indicator" {
		t.Errorf("Indicator().Name = %q, want %q", attr.Name, "data-indicator")
	}
	if attr.Value != "loading" {
		t.Errorf("Indicator().Value = %q, want %q", attr.Value, "loading")
	}

	// With $ prefix
	attr = Indicator("$loading")
	if attr.Value != "loading" {
		t.Errorf("Indicator($loading).Value = %q, want %q", attr.Value, "loading")
	}
}

func TestIgnore(t *testing.T) {
	attr := Ignore()
	if attr.Name != "data-ignore" {
		t.Errorf("Ignore().Name = %q, want %q", attr.Name, "data-ignore")
	}
	if attr.Value != "" {
		t.Errorf("Ignore().Value = %q, want empty", attr.Value)
	}
}

func TestEffect(t *testing.T) {
	attr := Effect(Raw("console.log($count)"))
	if attr.Name != "data-effect" {
		t.Errorf("Effect().Name = %q, want %q", attr.Name, "data-effect")
	}
}

func TestPeek(t *testing.T) {
	v := Peek(Raw("$foo"))
	got := ToJS(v.expr)
	expected := "@peek(() => $foo)"
	if got != expected {
		t.Errorf("Peek() = %q, want %q", got, expected)
	}
}

func TestComputed(t *testing.T) {
	attr := Computed("total", Raw("$price * $qty")).Attribute()
	if attr.Name != "data-computed:total" {
		t.Errorf("Computed().Name = %q, want %q", attr.Name, "data-computed:total")
	}
	if attr.Value != "$price * $qty" {
		t.Errorf("Computed().Value = %q, want %q", attr.Value, "$price * $qty")
	}
}

func TestComputedExpr(t *testing.T) {
	attr := Computed("total", Raw("$price * $qty")).Case(CamelCase).Attribute()
	if !strings.HasPrefix(attr.Name, "data-computed:") {
		t.Errorf("ComputedExpr().Name should start with data-computed:, got %q", attr.Name)
	}
	if !strings.Contains(attr.Name, "total") {
		t.Errorf("ComputedExpr().Name should contain total, got %q", attr.Name)
	}
}

func TestInit(t *testing.T) {
	attr := Init(Raw("$count = 1")).Attribute()
	if attr.Name != "data-init" {
		t.Errorf("Init().Name = %q, want %q", attr.Name, "data-init")
	}
	if attr.Value != "$count = 1" {
		t.Errorf("Init().Value = %q, want %q", attr.Value, "$count = 1")
	}

	attr = Init(Raw("$count = 1")).Delay(500 * time.Millisecond).ViewTransition().Attribute()
	if !strings.Contains(attr.Name, "__delay.500ms") || !strings.Contains(attr.Name, "__viewtransition") {
		t.Errorf("Init() modifiers should contain __delay.500ms and __viewtransition, got %q", attr.Name)
	}
}

func TestRef(t *testing.T) {
	attr := Ref("myElement").Attribute()
	if attr.Name != "data-ref:myElement" {
		t.Errorf("Ref().Name = %q, want %q", attr.Name, "data-ref:myElement")
	}
}

func TestStyle(t *testing.T) {
	attr := Style("background-color", Raw("$bgColor"))
	if attr.Name != "data-style:background-color" {
		t.Errorf("Style().Name = %q, want %q", attr.Name, "data-style:background-color")
	}
}

func TestStyles(t *testing.T) {
	attr := Styles(map[string]string{"color": "$red ? 'red' : 'blue'"})
	if attr.Name != "data-style" {
		t.Errorf("Styles().Name = %q, want %q", attr.Name, "data-style")
	}
}

func TestAttrs(t *testing.T) {
	attr := Attrs(map[string]string{"title": "$foo"})
	if attr.Name != "data-attr" {
		t.Errorf("Attrs().Name = %q, want %q", attr.Name, "data-attr")
	}
}

func TestClasses(t *testing.T) {
	attr := Classes(map[string]string{"hidden": "$foo"})
	if attr.Name != "data-class" {
		t.Errorf("Classes().Name = %q, want %q", attr.Name, "data-class")
	}
}

func TestIgnoreMorph(t *testing.T) {
	attr := IgnoreMorph()
	if attr.Name != "data-ignore-morph" {
		t.Errorf("IgnoreMorph().Name = %q, want %q", attr.Name, "data-ignore-morph")
	}
	if attr.Value != "" {
		t.Errorf("IgnoreMorph().Value = %q, want empty", attr.Value)
	}
}

func TestPreserveAttr(t *testing.T) {
	attr := PreserveAttr("open", "class")
	if attr.Name != "data-preserve-attr" {
		t.Errorf("PreserveAttr().Name = %q, want %q", attr.Name, "data-preserve-attr")
	}
	if attr.Value != "open class" {
		t.Errorf("PreserveAttr().Value = %q, want %q", attr.Value, "open class")
	}
}

func TestJsonSignalsDebug(t *testing.T) {
	// Without options
	attr := JsonSignalsDebug(nil).Attribute()
	if attr.Name != "data-json-signals" {
		t.Errorf("JsonSignalsDebug(nil).Name = %q, want %q", attr.Name, "data-json-signals")
	}

	// With options
	pattern := "user"
	attr = JsonSignalsDebug(&FilterOptions{IncludeReg: &pattern}).Attribute()
	if !strings.Contains(attr.Value, "include: /user/") {
		t.Errorf("JsonSignalsDebug().Value should contain include regex, got %q", attr.Value)
	}

	// With modifiers only
	attr = JsonSignalsDebug(nil).Terse().Attribute()
	if !strings.Contains(attr.Name, "__terse") {
		t.Errorf("JsonSignalsDebug(nil).Terse().Name should contain __terse, got %q", attr.Name)
	}
}

func TestBindKey(t *testing.T) {
	attr := BindKey("foo").Case(CamelCase).Attribute()
	if !strings.HasPrefix(attr.Name, "data-bind:") {
		t.Errorf("BindKey().Name should start with data-bind:, got %q", attr.Name)
	}
	if !strings.Contains(attr.Name, "foo") {
		t.Errorf("BindKey().Name should contain foo, got %q", attr.Name)
	}
}

func TestIndicatorKey(t *testing.T) {
	attr := IndicatorKey("fetching").Case(CamelCase).Attribute()
	if !strings.HasPrefix(attr.Name, "data-indicator:") {
		t.Errorf("IndicatorKey().Name should start with data-indicator:, got %q", attr.Name)
	}
	if !strings.Contains(attr.Name, "fetching") {
		t.Errorf("IndicatorKey().Name should contain fetching, got %q", attr.Name)
	}

	// With $ prefix
	attr = IndicatorKey("$fetching").Attribute()
	if !strings.Contains(attr.Name, "fetching") {
		t.Errorf("IndicatorKey($fetching).Name should contain fetching (without $), got %q", attr.Name)
	}
}

func TestIgnoreSelf(t *testing.T) {
	attr := IgnoreSelf()
	if attr.Name != "data-ignore__self" {
		t.Errorf("IgnoreSelf().Name = %q, want %q", attr.Name, "data-ignore__self")
	}
}

func TestFilterOptions(t *testing.T) {
	include := "^user"
	exclude := "password$"

	tests := []struct {
		name     string
		opts     *FilterOptions
		expected string
	}{
		{"include only", &FilterOptions{IncludeReg: &include}, "{include: /^user/}"},
		{"exclude only", &FilterOptions{ExcludeReg: &exclude}, "{exclude: /password$/}"},
		{"both", &FilterOptions{IncludeReg: &include, ExcludeReg: &exclude}, "{include: /^user/, exclude: /password$/}"},
		{"neither", &FilterOptions{}, "{}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sb strings.Builder
			tt.opts.appendJS(&sb)
			if got := sb.String(); got != tt.expected {
				t.Errorf("appendJS() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestSetAll(t *testing.T) {
	v := SetAll(JsonValue(true), nil)
	got := ToJS(v.expr)
	expected := "@setAll(true)"
	if got != expected {
		t.Errorf("SetAll() = %q, want %q", got, expected)
	}

	// With filter
	include := "^check"
	v = SetAll(JsonValue(false), &FilterOptions{IncludeReg: &include})
	got = ToJS(v.expr)
	if !strings.Contains(got, "@setAll(false, {include: /^check/})") {
		t.Errorf("SetAll() with filter = %q, should contain filter options", got)
	}
}

func TestToggleAll(t *testing.T) {
	v := ToggleAll(nil)
	got := ToJS(v.expr)
	expected := "@toggleAll()"
	if got != expected {
		t.Errorf("ToggleAll() = %q, want %q", got, expected)
	}

	// With filter
	include := "^check"
	v = ToggleAll(&FilterOptions{IncludeReg: &include})
	got = ToJS(v.expr)
	if !strings.Contains(got, "@toggleAll({include: /^check/})") {
		t.Errorf("ToggleAll() with filter = %q, should contain filter options", got)
	}
}

// ============ actions.go tests ============

func TestGet(t *testing.T) {
	v := Get("/api/users")

	got := ToJS(v.expr)
	expected := `@get("/api/users")`
	if got != expected {
		t.Errorf("Get() = %q, want %q", got, expected)
	}

	// Test as event-handler value
	attr := OnClick(v).Attribute()
	if attr.Value != expected {
		t.Errorf("OnClick(Get()).Value = %q, want %q", attr.Value, expected)
	}
}

func TestPost(t *testing.T) {
	v := Post("/api/users")
	got := ToJS(v.expr)
	if !strings.HasPrefix(got, `@post("/api/users")`) {
		t.Errorf("Post() = %q, should start with @post", got)
	}
}

func TestPut(t *testing.T) {
	v := Put("/api/users/1")
	got := ToJS(v.expr)
	if !strings.HasPrefix(got, `@put("/api/users/1")`) {
		t.Errorf("Put() = %q, should start with @put", got)
	}
}

func TestDelete(t *testing.T) {
	v := Delete("/api/users/1")
	got := ToJS(v.expr)
	if !strings.HasPrefix(got, `@delete("/api/users/1")`) {
		t.Errorf("Delete() = %q, should start with @delete", got)
	}
}

func TestPatch(t *testing.T) {
	v := Patch("/api/users/1")
	got := ToJS(v.expr)
	if !strings.HasPrefix(got, `@patch("/api/users/1")`) {
		t.Errorf("Patch() = %q, should start with @patch", got)
	}
}

func TestGetDynamic(t *testing.T) {
	v := GetDynamic(Raw("`/api/users/${$userId}`"))
	got := ToJS(v.expr)
	expected := "@get(`/api/users/${$userId}`)"
	if got != expected {
		t.Errorf("GetDynamic() = %q, want %q", got, expected)
	}

	// Test as event-handler value
	attr := OnClick(v).Attribute()
	if attr.Value != expected {
		t.Errorf("OnClick(GetDynamic()).Value = %q, want %q", attr.Value, expected)
	}
}

func TestPostDynamic(t *testing.T) {
	v := PostDynamic(Raw("`/api/${$resource}`"))
	got := ToJS(v.expr)
	if !strings.HasPrefix(got, "@post(") {
		t.Errorf("PostDynamic() = %q, should start with @post(", got)
	}
}

func TestPutDynamic(t *testing.T) {
	v := PutDynamic(Raw("`/api/${$resource}`"))
	got := ToJS(v.expr)
	if !strings.HasPrefix(got, "@put(") {
		t.Errorf("PutDynamic() = %q, should start with @put(", got)
	}
}

func TestDeleteDynamic(t *testing.T) {
	v := DeleteDynamic(Raw("`/api/${$resource}`"))
	got := ToJS(v.expr)
	if !strings.HasPrefix(got, "@delete(") {
		t.Errorf("DeleteDynamic() = %q, should start with @delete(", got)
	}
}

func TestPatchDynamic(t *testing.T) {
	v := PatchDynamic(Raw("`/api/${$resource}`"))
	got := ToJS(v.expr)
	if !strings.HasPrefix(got, "@patch(") {
		t.Errorf("PatchDynamic() = %q, should start with @patch(", got)
	}
}

func TestRequestWithOptions(t *testing.T) {
	v := Get("/api/data", OnSuccess(Raw("$done = true")), OnFailure(Raw("$error = true")))
	got := ToJS(v.expr)
	if !strings.Contains(got, ".then(") {
		t.Errorf("Get with OnSuccess should contain .then(), got %q", got)
	}
	if !strings.Contains(got, ".catch(") {
		t.Errorf("Get with OnFailure should contain .catch(), got %q", got)
	}
}

func TestRequestOptionsBuilder(t *testing.T) {
	tests := []struct {
		name     string
		builder  RequestOptionsBuilder
		expected string
	}{
		{"content type", RequestOptions().ContentType("form"), `contentType: "form"`},
		{"selector", RequestOptions().Selector("#myForm"), `selector: "#myForm"`},
		{"open when hidden", RequestOptions().OpenWhenHidden(true), `openWhenHidden: true`},
		{"open when hidden false", RequestOptions().OpenWhenHidden(false), `openWhenHidden: false`},
		{"retry interval", RequestOptions().RetryInterval(2000), `retryInterval: 2000`},
		{"retry scaler", RequestOptions().RetryScaler(1.5), `retryScaler: 1.5`},
		{"retry max wait", RequestOptions().RetryMaxWait(60000), `retryMaxWait: 60000`},
		{"retry max count", RequestOptions().RetryMaxCount(5), `retryMaxCount: 5`},
		{"request cancellation", RequestOptions().RequestCancellation("disabled"), `requestCancellation: "disabled"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := GetWithOptions("/api", tt.builder)
			got := ToJS(v.expr)
			if !strings.Contains(got, tt.expected) {
				t.Errorf("RequestOptions with %s = %q, should contain %q", tt.name, got, tt.expected)
			}
		})
	}
}

func TestHeaders(t *testing.T) {
	v := GetWithOptions("/api", RequestOptions().Headers(map[string]string{"X-Custom": "value"}))
	got := ToJS(v.expr)
	if !strings.Contains(got, `"X-Custom"`) || !strings.Contains(got, `"value"`) {
		t.Errorf("Headers() = %q, should contain X-Custom and value", got)
	}
}

func TestFilterSignals(t *testing.T) {
	include := "^user"
	v := GetWithOptions("/api", RequestOptions().FilterSignals(&FilterOptions{IncludeReg: &include}))
	got := ToJS(v.expr)
	if !strings.Contains(got, "filterSignals:") {
		t.Errorf("FilterSignals() = %q, should contain filterSignals:", got)
	}
}

func TestRetry(t *testing.T) {
	tests := []struct {
		name     string
		mode     string
		expected string
	}{
		{"auto", "auto", `retry: "auto"`},
		{"error", "error", `retry: "error"`},
		{"always", "always", `retry: "always"`},
		{"never", "never", `retry: "never"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := GetWithOptions("/api", RequestOptions().Retry(tt.mode))
			got := ToJS(v.expr)
			if !strings.Contains(got, tt.expected) {
				t.Errorf("Retry(%q) = %q, should contain %q", tt.mode, got, tt.expected)
			}
		})
	}
}

func TestPayload(t *testing.T) {
	tests := []struct {
		name     string
		data     any
		expected string
	}{
		{"map", map[string]any{"name": "John"}, `payload: {"name":"John"}`},
		{"simple value", 42, `payload: 42`},
		{"string", "hello", `payload: "hello"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := PostWithOptions("/api", RequestOptions().Payload(tt.data))
			got := ToJS(v.expr)
			if !strings.Contains(got, tt.expected) {
				t.Errorf("Payload(%v) = %q, should contain %q", tt.data, got, tt.expected)
			}
		})
	}
}

// ============ builders.go tests ============

func TestMultipleStatements(t *testing.T) {
	attr := OnClick(Raw("$a = 1"), Raw("$b = 2")).Attribute()
	if attr.Value != "$a = 1; $b = 2" {
		t.Errorf("multiple statements = %q, want %q", attr.Value, "$a = 1; $b = 2")
	}
}
