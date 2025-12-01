package d

import (
	"strings"
	"testing"
	"time"
)

// Helper to build attribute and return name and value
func buildTestAttr(name string, options ...AttrMutator) (string, string) {
	attr := buildAttr(name, options...)
	return attr.name.String(), strings.Join(attr.statements, "; ")
}

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

			// Test AttrValueAppender
			var sb strings.Builder
			v.Append(&sb)
			if got := sb.String(); got != tt.expected {
				t.Errorf("Append() = %q, want %q", got, tt.expected)
			}

			// Test AttrMutator
			attrName, attrValue := buildTestAttr("test", v)
			if attrValue != tt.expected {
				t.Errorf("Modify() value = %q, want %q", attrValue, tt.expected)
			}
			if attrName != "test" {
				t.Errorf("Modify() name = %q, want %q", attrName, "test")
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
			var sb strings.Builder
			v.Append(&sb)
			if got := sb.String(); got != tt.expected {
				t.Errorf("Str(%q).Append() = %q, want %q", tt.input, got, tt.expected)
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
			var sb strings.Builder
			v.Append(&sb)
			if got := sb.String(); got != tt.expected {
				t.Errorf("JsonValue(%v).Append() = %q, want %q", tt.input, got, tt.expected)
			}

			// Also test AttrMutator path
			_, attrValue := buildTestAttr("test", v)
			if attrValue != tt.expected {
				t.Errorf("JsonValue(%v).Modify() = %q, want %q", tt.input, attrValue, tt.expected)
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
	v := JsonValue(make(chan int))
	var sb strings.Builder
	v.Append(&sb)
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
	action := OnSuccess("console.log('done')")
	var sb strings.Builder
	action.Append(&sb)
	expected := ".then(() => console.log('done'))"
	if got := sb.String(); got != expected {
		t.Errorf("OnSuccess() = %q, want %q", got, expected)
	}
}

func TestOnFailure(t *testing.T) {
	action := OnFailure("console.log(error)")
	var sb strings.Builder
	action.Append(&sb)
	expected := ".catch((error) => console.log(error))"
	if got := sb.String(); got != expected {
		t.Errorf("OnFailure() = %q, want %q", got, expected)
	}
}

func TestConsoleLog(t *testing.T) {
	tests := []struct {
		name     string
		values   []AttrValueAppender
		expected string
	}{
		{"single value", []AttrValueAppender{Raw("$foo")}, "console.log($foo)"},
		{"multiple values", []AttrValueAppender{Raw("$a"), Raw("$b")}, "console.log($a, $b)"},
		{"with string", []AttrValueAppender{Str("msg"), Raw("$x")}, `console.log("msg", $x)`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action := ConsoleLog(tt.values...)
			_, attrValue := buildTestAttr("test", action)
			if attrValue != tt.expected {
				t.Errorf("ConsoleLog() = %q, want %q", attrValue, tt.expected)
			}
		})
	}
}

func TestAnd(t *testing.T) {
	tests := []struct {
		name     string
		actions  []AttrValueAppender
		expected string
	}{
		{"two actions", []AttrValueAppender{Raw("$a"), Raw("$b")}, "$a && $b"},
		{"three actions", []AttrValueAppender{Raw("$a"), Raw("$b"), Raw("$c")}, "$a && $b && $c"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := And(tt.actions...)

			// Test AttrValueAppender
			var sb strings.Builder
			v.Append(&sb)
			if got := sb.String(); got != tt.expected {
				t.Errorf("And().Append() = %q, want %q", got, tt.expected)
			}

			// Test AttrMutator
			_, attrValue := buildTestAttr("test", v)
			if attrValue != tt.expected {
				t.Errorf("And().Modify() = %q, want %q", attrValue, tt.expected)
			}
		})
	}
}

// ============ modifiers.go tests ============

func TestPreventDefault(t *testing.T) {
	attrName, _ := buildTestAttr("data-on:click", PreventDefault())
	if !strings.Contains(attrName, "__prevent") {
		t.Errorf("PreventDefault() should add __prevent, got %q", attrName)
	}
}

func TestOnce(t *testing.T) {
	attrName, _ := buildTestAttr("data-on:click", Once())
	if !strings.Contains(attrName, "__once") {
		t.Errorf("Once() should add __once, got %q", attrName)
	}
}

func TestPassive(t *testing.T) {
	attrName, _ := buildTestAttr("data-on:scroll", Passive())
	if !strings.Contains(attrName, "__passive") {
		t.Errorf("Passive() should add __passive, got %q", attrName)
	}
}

func TestCapture(t *testing.T) {
	attrName, _ := buildTestAttr("data-on:click", Capture())
	if !strings.Contains(attrName, "__capture") {
		t.Errorf("Capture() should add __capture, got %q", attrName)
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
			attrName, _ := buildTestAttr("data-signals:", Case(tt.casing))
			if !strings.Contains(attrName, tt.expected) {
				t.Errorf("Case(%s) should add %s, got %q", tt.casing, tt.expected, attrName)
			}
		})
	}
}

func TestDelay(t *testing.T) {
	attrName, _ := buildTestAttr("data-on:click", Delay(500*time.Millisecond))
	if !strings.Contains(attrName, "__delay.500ms") {
		t.Errorf("Delay() should add __delay.500ms, got %q", attrName)
	}
}

func TestDebounce(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		opts     []TimingMutatorFunc
		expected string
	}{
		{"basic", 300 * time.Millisecond, nil, "__debounce.300ms"},
		{"with noleading", 500 * time.Millisecond, []TimingMutatorFunc{NoLeading()}, "__debounce.500ms.noleading"},
		{"with notrailing", 500 * time.Millisecond, []TimingMutatorFunc{NoTrailing()}, "__debounce.500ms.notrailing"},
		{"with both", 500 * time.Millisecond, []TimingMutatorFunc{NoLeading(), NoTrailing()}, "__debounce.500ms.noleading.notrailing"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attrName, _ := buildTestAttr("data-on:input", Debounce(tt.duration, tt.opts...))
			if !strings.Contains(attrName, tt.expected) {
				t.Errorf("Debounce() should contain %s, got %q", tt.expected, attrName)
			}
		})
	}
}

func TestThrottle(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		opts     []TimingMutatorFunc
		expected string
	}{
		{"basic", 100 * time.Millisecond, nil, "__throttle.100ms"},
		{"with leading", 200 * time.Millisecond, []TimingMutatorFunc{Leading()}, "__throttle.200ms.leading"},
		{"with trailing", 200 * time.Millisecond, []TimingMutatorFunc{Trailing()}, "__throttle.200ms.trailing"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attrName, _ := buildTestAttr("data-on:scroll", Throttle(tt.duration, tt.opts...))
			if !strings.Contains(attrName, tt.expected) {
				t.Errorf("Throttle() should contain %s, got %q", tt.expected, attrName)
			}
		})
	}
}

func TestViewTransition(t *testing.T) {
	attrName, _ := buildTestAttr("data-on:click", ViewTransition())
	if !strings.Contains(attrName, "__viewtransition") {
		t.Errorf("ViewTransition() should add __viewtransition, got %q", attrName)
	}
}

func TestWindow(t *testing.T) {
	attrName, _ := buildTestAttr("data-on:scroll", Window())
	if !strings.Contains(attrName, "__window") {
		t.Errorf("Window() should add __window, got %q", attrName)
	}
}

func TestOutside(t *testing.T) {
	attrName, _ := buildTestAttr("data-on:click", Outside())
	if !strings.Contains(attrName, "__outside") {
		t.Errorf("Outside() should add __outside, got %q", attrName)
	}
}

func TestStopPropagation(t *testing.T) {
	attrName, _ := buildTestAttr("data-on:click", StopPropagation())
	if !strings.Contains(attrName, "__stop") {
		t.Errorf("StopPropagation() should add __stop, got %q", attrName)
	}
}

func TestIfMissing(t *testing.T) {
	attrName, _ := buildTestAttr("data-signals:", IfMissing())
	if !strings.Contains(attrName, "__ifmissing") {
		t.Errorf("IfMissing() should add __ifmissing, got %q", attrName)
	}
}

func TestSelf(t *testing.T) {
	attrName, _ := buildTestAttr("data-ignore", Self())
	if !strings.Contains(attrName, "__self") {
		t.Errorf("Self() should add __self, got %q", attrName)
	}
}

func TestTerse(t *testing.T) {
	attrName, _ := buildTestAttr("data-json-signals", Terse())
	if !strings.Contains(attrName, "__terse") {
		t.Errorf("Terse() should add __terse, got %q", attrName)
	}
}

func TestHalf(t *testing.T) {
	attrName, _ := buildTestAttr("data-on-intersect", Half())
	if !strings.Contains(attrName, "__half") {
		t.Errorf("Half() should add __half, got %q", attrName)
	}
}

func TestFull(t *testing.T) {
	attrName, _ := buildTestAttr("data-on-intersect", Full())
	if !strings.Contains(attrName, "__full") {
		t.Errorf("Full() should add __full, got %q", attrName)
	}
}

func TestDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		opts     []DurationMutatorFunc
		expected string
	}{
		{"basic", 500 * time.Millisecond, nil, "__duration.500ms"},
		{"with leading", 2 * time.Second, []DurationMutatorFunc{DurationLeading()}, "__duration.2s.leading"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attrName, _ := buildTestAttr("data-on-interval", Duration(tt.duration, tt.opts...))
			if !strings.Contains(attrName, tt.expected) {
				t.Errorf("Duration() should contain %s, got %q", tt.expected, attrName)
			}
		})
	}
}

// ============ attrs.go tests ============

func TestSetSignalExpr(t *testing.T) {
	tests := []struct {
		name       string
		signalName string
		expr       AttrValueAppender
		expected   string
	}{
		{"without $", "foo", Raw("1"), "$foo = 1"},
		{"with $", "$bar", Raw("true"), "$bar = true"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action := SetSignalExpr(tt.signalName, tt.expr)
			_, attrValue := buildTestAttr("test", action)
			if attrValue != tt.expected {
				t.Errorf("SetSignalExpr() = %q, want %q", attrValue, tt.expected)
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
			_, attrValue := buildTestAttr("test", action)
			if attrValue != tt.expected {
				t.Errorf("SetSignal() = %q, want %q", attrValue, tt.expected)
			}
		})
	}
}

func TestOnSubmit(t *testing.T) {
	attr := OnSubmit(Raw("$submit()"))
	if attr.Name != "data-on:submit" {
		t.Errorf("OnSubmit().Name = %q, want %q", attr.Name, "data-on:submit")
	}
	if attr.Value != "$submit()" {
		t.Errorf("OnSubmit().Value = %q, want %q", attr.Value, "$submit()")
	}
}

func TestOnInput(t *testing.T) {
	attr := OnInput(Raw("$validate()"))
	if attr.Name != "data-on:input" {
		t.Errorf("OnInput().Name = %q, want %q", attr.Name, "data-on:input")
	}
}

func TestOnChange(t *testing.T) {
	attr := OnChange(Raw("$update()"))
	if attr.Name != "data-on:change" {
		t.Errorf("OnChange().Name = %q, want %q", attr.Name, "data-on:change")
	}
}

func TestOnClick(t *testing.T) {
	attr := OnClick(Raw("$count++"))
	if attr.Name != "data-on:click" {
		t.Errorf("OnClick().Name = %q, want %q", attr.Name, "data-on:click")
	}
	if attr.Value != "$count++" {
		t.Errorf("OnClick().Value = %q, want %q", attr.Value, "$count++")
	}
}

func TestOnLoad(t *testing.T) {
	attr := OnLoad(Raw("$init()"))
	if attr.Name != "data-on:load" {
		t.Errorf("OnLoad().Name = %q, want %q", attr.Name, "data-on:load")
	}
}

func TestOn(t *testing.T) {
	attr := On("keydown", Raw("$handleKey(event)"))
	if attr.Name != "data-on:keydown" {
		t.Errorf("On().Name = %q, want %q", attr.Name, "data-on:keydown")
	}
	if attr.Value != "$handleKey(event)" {
		t.Errorf("On().Value = %q, want %q", attr.Value, "$handleKey(event)")
	}
}

func TestOnIntersect(t *testing.T) {
	attr := OnIntersect(Once(), Raw("$seen = true"))
	if !strings.HasPrefix(attr.Name, "data-on-intersect") {
		t.Errorf("OnIntersect().Name should start with data-on-intersect, got %q", attr.Name)
	}
	if !strings.Contains(attr.Name, "__once") {
		t.Errorf("OnIntersect().Name should contain __once, got %q", attr.Name)
	}
}

func TestOnInterval(t *testing.T) {
	attr := OnInterval(Duration(500*time.Millisecond), Raw("$tick()"))
	if !strings.HasPrefix(attr.Name, "data-on-interval") {
		t.Errorf("OnInterval().Name should start with data-on-interval, got %q", attr.Name)
	}
	if !strings.Contains(attr.Name, "__duration.500ms") {
		t.Errorf("OnInterval().Name should contain __duration.500ms, got %q", attr.Name)
	}
}

func TestOnSignalPatch(t *testing.T) {
	attr := OnSignalPatch(Raw("console.log('changed')"))
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
	attr := SignalExpr("count", Raw("0"))
	// The expression is appended to the name, value is empty
	if attr.Name != "data-signals:count0" {
		t.Errorf("SignalExpr().Name = %q, want %q", attr.Name, "data-signals:count0")
	}
	if attr.Value != "" {
		t.Errorf("SignalExpr().Value = %q, want empty", attr.Value)
	}
}

func TestSignal(t *testing.T) {
	attr := Signal("count", 0)
	if attr.Name != "data-signals:count" {
		t.Errorf("Signal().Name = %q, want %q", attr.Name, "data-signals:count")
	}
	if attr.Value != "0" {
		t.Errorf("Signal().Value = %q, want %q", attr.Value, "0")
	}

	attr = Signal("name", "hello")
	if attr.Value != `"hello"` {
		t.Errorf("Signal().Value = %q, want %q", attr.Value, `"hello"`)
	}
}

func TestSignals(t *testing.T) {
	attr := Signals(map[string]any{"foo": 1, "bar": "hello"})
	if attr.Name != "data-signals" {
		t.Errorf("Signals().Name = %q, want %q", attr.Name, "data-signals")
	}
	// JSON encoding of map, order may vary
	if !strings.Contains(attr.Value, `"foo":1`) || !strings.Contains(attr.Value, `"bar":"hello"`) {
		t.Errorf("Signals().Value = %q, should contain foo:1 and bar:hello", attr.Value)
	}
}

func TestBind(t *testing.T) {
	attr := Bind("username")
	if attr.Name != "data-bind" {
		t.Errorf("Bind().Name = %q, want %q", attr.Name, "data-bind")
	}
	if attr.Value != "username" {
		t.Errorf("Bind().Value = %q, want %q", attr.Value, "username")
	}
}

func TestClass(t *testing.T) {
	attr := Class("active", Raw("$isActive"))
	if attr.Name != "data-classactive" {
		t.Errorf("Class().Name = %q, want %q", attr.Name, "data-classactive")
	}
	if attr.Value != "$isActive" {
		t.Errorf("Class().Value = %q, want %q", attr.Value, "$isActive")
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
	var sb strings.Builder
	v.Append(&sb)
	expected := "@peek(() => $foo)"
	if got := sb.String(); got != expected {
		t.Errorf("Peek() = %q, want %q", got, expected)
	}
}

func TestComputed(t *testing.T) {
	attr := Computed("total", Raw("$price * $qty"))
	// The expression is appended to the name, value is empty
	if attr.Name != "data-computed:total$price * $qty" {
		t.Errorf("Computed().Name = %q, want %q", attr.Name, "data-computed:total$price * $qty")
	}
	if attr.Value != "" {
		t.Errorf("Computed().Value = %q, want empty", attr.Value)
	}
}

func TestComputedExpr(t *testing.T) {
	attr := ComputedExpr("total", Case(CamelCase), Raw("$price * $qty"))
	if !strings.HasPrefix(attr.Name, "data-computed:") {
		t.Errorf("ComputedExpr().Name should start with data-computed:, got %q", attr.Name)
	}
	if !strings.Contains(attr.Name, "total") {
		t.Errorf("ComputedExpr().Name should contain total, got %q", attr.Name)
	}
}

func TestInit(t *testing.T) {
	attr := Init(Raw("$count = 1"))
	if attr.Name != "data-init" {
		t.Errorf("Init().Name = %q, want %q", attr.Name, "data-init")
	}
	if attr.Value != "$count = 1" {
		t.Errorf("Init().Value = %q, want %q", attr.Value, "$count = 1")
	}
}

func TestRef(t *testing.T) {
	attr := Ref("myElement")
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
	attr := JsonSignalsDebug(nil)
	if attr.Name != "data-json-signals" {
		t.Errorf("JsonSignalsDebug(nil).Name = %q, want %q", attr.Name, "data-json-signals")
	}

	// With options
	pattern := "user"
	attr = JsonSignalsDebug(&FilterOptions{IncludeReg: &pattern})
	if !strings.Contains(attr.Value, "include: /user/") {
		t.Errorf("JsonSignalsDebug().Value should contain include regex, got %q", attr.Value)
	}

	// With modifiers only
	attr = JsonSignalsDebug(nil, Terse())
	if !strings.Contains(attr.Name, "__terse") {
		t.Errorf("JsonSignalsDebug(nil, Terse()).Name should contain __terse, got %q", attr.Name)
	}
}

func TestBindKey(t *testing.T) {
	attr := BindKey("foo", Case(CamelCase))
	if !strings.HasPrefix(attr.Name, "data-bind:") {
		t.Errorf("BindKey().Name should start with data-bind:, got %q", attr.Name)
	}
	if !strings.Contains(attr.Name, "foo") {
		t.Errorf("BindKey().Name should contain foo, got %q", attr.Name)
	}
}

func TestIndicatorKey(t *testing.T) {
	attr := IndicatorKey("fetching", Case(CamelCase))
	if !strings.HasPrefix(attr.Name, "data-indicator:") {
		t.Errorf("IndicatorKey().Name should start with data-indicator:, got %q", attr.Name)
	}
	if !strings.Contains(attr.Name, "fetching") {
		t.Errorf("IndicatorKey().Name should contain fetching, got %q", attr.Name)
	}

	// With $ prefix
	attr = IndicatorKey("$fetching")
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

func TestFilterOptionsValue(t *testing.T) {
	include := "^user"
	opts := &FilterOptions{IncludeReg: &include}
	v := FilterOptionsValue(opts)

	var sb strings.Builder
	v.Append(&sb)
	if !strings.Contains(sb.String(), "include: /^user/") {
		t.Errorf("FilterOptionsValue().Append() = %q, should contain include regex", sb.String())
	}
}

func TestSetAll(t *testing.T) {
	v := SetAll(JsonValue(true), FilterOptions{})
	var sb strings.Builder
	v.Append(&sb)
	expected := "@setAll(true)"
	if got := sb.String(); got != expected {
		t.Errorf("SetAll() = %q, want %q", got, expected)
	}

	// With filter
	include := "^check"
	v = SetAll(JsonValue(false), FilterOptions{IncludeReg: &include})
	sb.Reset()
	v.Append(&sb)
	if !strings.Contains(sb.String(), "@setAll(false, {include: /^check/})") {
		t.Errorf("SetAll() with filter = %q, should contain filter options", sb.String())
	}
}

func TestToggleAll(t *testing.T) {
	v := ToggleAll(FilterOptions{})
	var sb strings.Builder
	v.Append(&sb)
	expected := "@toggleAll()"
	if got := sb.String(); got != expected {
		t.Errorf("ToggleAll() = %q, want %q", got, expected)
	}

	// With filter
	include := "^check"
	v = ToggleAll(FilterOptions{IncludeReg: &include})
	sb.Reset()
	v.Append(&sb)
	if !strings.Contains(sb.String(), "@toggleAll({include: /^check/})") {
		t.Errorf("ToggleAll() with filter = %q, should contain filter options", sb.String())
	}
}

// ============ actions.go tests ============

func TestGet(t *testing.T) {
	v := Get("/api/users")

	var sb strings.Builder
	v.Append(&sb)
	expected := `@get("/api/users")`
	if got := sb.String(); got != expected {
		t.Errorf("Get() = %q, want %q", got, expected)
	}

	// Test as AttrMutator
	_, attrValue := buildTestAttr("test", v)
	if attrValue != expected {
		t.Errorf("Get().Modify() = %q, want %q", attrValue, expected)
	}
}

func TestPost(t *testing.T) {
	v := Post("/api/users")
	var sb strings.Builder
	v.Append(&sb)
	if !strings.HasPrefix(sb.String(), `@post("/api/users")`) {
		t.Errorf("Post() = %q, should start with @post", sb.String())
	}
}

func TestPut(t *testing.T) {
	v := Put("/api/users/1")
	var sb strings.Builder
	v.Append(&sb)
	if !strings.HasPrefix(sb.String(), `@put("/api/users/1")`) {
		t.Errorf("Put() = %q, should start with @put", sb.String())
	}
}

func TestDelete(t *testing.T) {
	v := Delete("/api/users/1")
	var sb strings.Builder
	v.Append(&sb)
	if !strings.HasPrefix(sb.String(), `@delete("/api/users/1")`) {
		t.Errorf("Delete() = %q, should start with @delete", sb.String())
	}
}

func TestPatch(t *testing.T) {
	v := Patch("/api/users/1")
	var sb strings.Builder
	v.Append(&sb)
	if !strings.HasPrefix(sb.String(), `@patch("/api/users/1")`) {
		t.Errorf("Patch() = %q, should start with @patch", sb.String())
	}
}

func TestGetDynamic(t *testing.T) {
	v := GetDynamic(Raw("`/api/users/${$userId}`"))
	var sb strings.Builder
	v.Append(&sb)
	expected := "@get(`/api/users/${$userId}`)"
	if got := sb.String(); got != expected {
		t.Errorf("GetDynamic() = %q, want %q", got, expected)
	}

	// Test as AttrMutator
	_, attrValue := buildTestAttr("test", v)
	if attrValue != expected {
		t.Errorf("GetDynamic().Modify() = %q, want %q", attrValue, expected)
	}
}

func TestPostDynamic(t *testing.T) {
	v := PostDynamic(Raw("`/api/${$resource}`"))
	var sb strings.Builder
	v.Append(&sb)
	if !strings.HasPrefix(sb.String(), "@post(") {
		t.Errorf("PostDynamic() = %q, should start with @post(", sb.String())
	}
}

func TestPutDynamic(t *testing.T) {
	v := PutDynamic(Raw("`/api/${$resource}`"))
	var sb strings.Builder
	v.Append(&sb)
	if !strings.HasPrefix(sb.String(), "@put(") {
		t.Errorf("PutDynamic() = %q, should start with @put(", sb.String())
	}
}

func TestDeleteDynamic(t *testing.T) {
	v := DeleteDynamic(Raw("`/api/${$resource}`"))
	var sb strings.Builder
	v.Append(&sb)
	if !strings.HasPrefix(sb.String(), "@delete(") {
		t.Errorf("DeleteDynamic() = %q, should start with @delete(", sb.String())
	}
}

func TestPatchDynamic(t *testing.T) {
	v := PatchDynamic(Raw("`/api/${$resource}`"))
	var sb strings.Builder
	v.Append(&sb)
	if !strings.HasPrefix(sb.String(), "@patch(") {
		t.Errorf("PatchDynamic() = %q, should start with @patch(", sb.String())
	}
}

func TestRequestWithOptions(t *testing.T) {
	v := Get("/api/data", OnSuccess("$done = true"), OnFailure("$error = true"))
	var sb strings.Builder
	v.Append(&sb)
	got := sb.String()
	if !strings.Contains(got, ".then(") {
		t.Errorf("Get with OnSuccess should contain .then(), got %q", got)
	}
	if !strings.Contains(got, ".catch(") {
		t.Errorf("Get with OnFailure should contain .catch(), got %q", got)
	}
}

func TestRequestOptions(t *testing.T) {
	tests := []struct {
		name     string
		opts     []RequestOption
		expected string
	}{
		{"empty", nil, ""},
		{"content type", []RequestOption{ContentType("form")}, `, {contentType: "form"}`},
		{"selector", []RequestOption{Selector("#myForm")}, `, {selector: "#myForm"}`},
		{"open when hidden", []RequestOption{OpenWhenHidden(true)}, `, {openWhenHidden: true}`},
		{"open when hidden false", []RequestOption{OpenWhenHidden(false)}, `, {openWhenHidden: false}`},
		{"retry interval", []RequestOption{RetryInterval(2000)}, `, {retryInterval: 2000}`},
		{"retry scaler", []RequestOption{RetryScaler(1.5)}, `, {retryScaler: 1.5}`},
		{"retry max wait", []RequestOption{RetryMaxWaitMs(60000)}, `, {retryMaxWaitMs: 60000}`},
		{"retry max count", []RequestOption{RetryMaxCount(5)}, `, {retryMaxCount: 5}`},
		{"request cancellation", []RequestOption{RequestCancellation("disabled")}, `, {requestCancellation: "disabled"}`},
		{"multiple options", []RequestOption{ContentType("json"), OpenWhenHidden(true)}, `, {contentType: "json", openWhenHidden: true}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := RequestOptions(tt.opts...)
			var sb strings.Builder
			opt.Append(&sb)
			if got := sb.String(); got != tt.expected {
				t.Errorf("RequestOptions() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestHeaders(t *testing.T) {
	opt := Headers(map[string]string{"X-Custom": "value"})
	var sb strings.Builder
	requestOptionFunc(func(sb *strings.Builder) {
		opt.appendOption(sb)
	}).appendOption(&sb)
	if !strings.Contains(sb.String(), `"X-Custom"`) || !strings.Contains(sb.String(), `"value"`) {
		t.Errorf("Headers() = %q, should contain X-Custom and value", sb.String())
	}
}

func TestFilterSignals(t *testing.T) {
	include := "^user"
	opt := FilterSignals(&FilterOptions{IncludeReg: &include})
	var sb strings.Builder
	requestOptionFunc(func(sb *strings.Builder) {
		opt.appendOption(sb)
	}).appendOption(&sb)
	if !strings.Contains(sb.String(), "filterSignals:") {
		t.Errorf("FilterSignals() = %q, should contain filterSignals:", sb.String())
	}
}

// ============ builders.go tests ============

func TestValueMutatorFunc(t *testing.T) {
	v := ValueMutatorFunc(func(sb *strings.Builder) {
		sb.WriteString("custom value")
	})

	// Test AttrValueAppender
	var sb strings.Builder
	v.Append(&sb)
	if got := sb.String(); got != "custom value" {
		t.Errorf("ValueMutatorFunc().Append() = %q, want %q", got, "custom value")
	}

	// Test AttrMutator
	_, attrValue := buildTestAttr("test", v)
	if attrValue != "custom value" {
		t.Errorf("ValueMutatorFunc().Modify() = %q, want %q", attrValue, "custom value")
	}
}

func TestBuildAttr(t *testing.T) {
	attr := buildAttr("data-test", Raw("$foo"), Raw("$bar"))
	if attr.name.String() != "data-test" {
		t.Errorf("buildAttr().name = %q, want %q", attr.name.String(), "data-test")
	}
	if len(attr.statements) != 2 {
		t.Errorf("buildAttr().statements length = %d, want 2", len(attr.statements))
	}
}

func TestExprAttr(t *testing.T) {
	attr := exprAttr("data-on:click", PreventDefault(), Raw("$count++"))
	if !strings.HasPrefix(attr.Name, "data-on:click") {
		t.Errorf("exprAttr().Name = %q, should start with data-on:click", attr.Name)
	}
	if attr.Value != "$count++" {
		t.Errorf("exprAttr().Value = %q, want %q", attr.Value, "$count++")
	}
}

func TestMultipleStatements(t *testing.T) {
	attr := exprAttr("data-on:click", Raw("$a = 1"), Raw("$b = 2"))
	if attr.Value != "$a = 1; $b = 2" {
		t.Errorf("multiple statements = %q, want %q", attr.Value, "$a = 1; $b = 2")
	}
}
