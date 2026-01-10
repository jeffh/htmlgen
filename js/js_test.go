package js

import (
	"strings"
	"testing"
)

func exprString(e Expr) string {
	var sb strings.Builder
	e.js(&sb)
	return sb.String()
}

func stmtString(s Stmt) string {
	var sb strings.Builder
	s.stmt(&sb)
	return sb.String()
}

// === Value Tests ===

func TestString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", `"hello"`},
		{`say "hi"`, `"say \"hi\""`},
		{"line1\nline2", `"line1\nline2"`},
		{"", `""`},
		{"<script>alert('xss')</script>", `"\u003cscript\u003ealert('xss')\u003c/script\u003e"`},
	}
	for _, tt := range tests {
		got := exprString(String(tt.input))
		if got != tt.expected {
			t.Errorf("String(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestInt(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{0, "0"},
		{42, "42"},
		{-17, "-17"},
		{1000000, "1000000"},
	}
	for _, tt := range tests {
		got := exprString(Int(tt.input))
		if got != tt.expected {
			t.Errorf("Int(%d) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestFloat(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{3.14, "3.14"},
		{0.0, "0"},
		{-2.5, "-2.5"},
	}
	for _, tt := range tests {
		got := exprString(Float(tt.input))
		if got != tt.expected {
			t.Errorf("Float(%v) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestBool(t *testing.T) {
	if got := exprString(Bool(true)); got != "true" {
		t.Errorf("Bool(true) = %q, want %q", got, "true")
	}
	if got := exprString(Bool(false)); got != "false" {
		t.Errorf("Bool(false) = %q, want %q", got, "false")
	}
}

func TestNull(t *testing.T) {
	if got := exprString(Null()); got != "null" {
		t.Errorf("Null() = %q, want %q", got, "null")
	}
}

func TestUndefined(t *testing.T) {
	if got := exprString(Undefined()); got != "undefined" {
		t.Errorf("Undefined() = %q, want %q", got, "undefined")
	}
}

func TestJSON(t *testing.T) {
	tests := []struct {
		input    any
		expected string
	}{
		{42, "42"},
		{"hello", `"hello"`},
		{[]int{1, 2, 3}, "[1,2,3]"},
		{map[string]int{"a": 1}, `{"a":1}`},
		{true, "true"},
		{nil, "null"},
	}
	for _, tt := range tests {
		got := exprString(JSON(tt.input))
		if got != tt.expected {
			t.Errorf("JSON(%v) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestArray(t *testing.T) {
	got := exprString(Array(Int(1), Int(2), String("three")))
	expected := `[1, 2, "three"]`
	if got != expected {
		t.Errorf("Array() = %q, want %q", got, expected)
	}
}

func TestArrayEmpty(t *testing.T) {
	got := exprString(Array())
	expected := `[]`
	if got != expected {
		t.Errorf("Array() = %q, want %q", got, expected)
	}
}

func TestObject(t *testing.T) {
	got := exprString(Object(
		Pair("name", String("John")),
		Pair("age", Int(30)),
	))
	expected := `{"name": "John", "age": 30}`
	if got != expected {
		t.Errorf("Object() = %q, want %q", got, expected)
	}
}

func TestObjectEmpty(t *testing.T) {
	got := exprString(Object())
	expected := `{}`
	if got != expected {
		t.Errorf("Object() = %q, want %q", got, expected)
	}
}

func TestIdent(t *testing.T) {
	got := exprString(Ident("myVar"))
	if got != "myVar" {
		t.Errorf("Ident() = %q, want %q", got, "myVar")
	}
}

func TestThis(t *testing.T) {
	got := exprString(This())
	if got != "this" {
		t.Errorf("This() = %q, want %q", got, "this")
	}
}

// === Raw Tests ===

func TestRaw(t *testing.T) {
	got := exprString(Raw("customFn()"))
	if got != "customFn()" {
		t.Errorf("Raw() = %q, want %q", got, "customFn()")
	}
}

// === Access Tests ===

func TestProp(t *testing.T) {
	got := exprString(Prop(Ident("document"), "body"))
	if got != "document.body" {
		t.Errorf("Prop() = %q, want %q", got, "document.body")
	}
}

func TestPropChained(t *testing.T) {
	got := exprString(Prop(Prop(Ident("window"), "document"), "body"))
	if got != "window.document.body" {
		t.Errorf("Prop chained = %q, want %q", got, "window.document.body")
	}
}

func TestIndex(t *testing.T) {
	got := exprString(Index(Ident("arr"), Int(0)))
	if got != "arr[0]" {
		t.Errorf("Index() = %q, want %q", got, "arr[0]")
	}
}

func TestIndexWithString(t *testing.T) {
	got := exprString(Index(Ident("obj"), String("key")))
	expected := `obj["key"]`
	if got != expected {
		t.Errorf("Index() = %q, want %q", got, expected)
	}
}

func TestCall(t *testing.T) {
	got := exprString(Call(Ident("alert"), String("hello")))
	expected := `alert("hello")`
	if got != expected {
		t.Errorf("Call() = %q, want %q", got, expected)
	}
}

func TestCallNoArgs(t *testing.T) {
	got := exprString(Call(Ident("doSomething")))
	expected := `doSomething()`
	if got != expected {
		t.Errorf("Call() = %q, want %q", got, expected)
	}
}

func TestMethod(t *testing.T) {
	got := exprString(Method(Ident("console"), "log", String("msg")))
	expected := `console.log("msg")`
	if got != expected {
		t.Errorf("Method() = %q, want %q", got, expected)
	}
}

func TestMethodMultipleArgs(t *testing.T) {
	got := exprString(Method(Ident("console"), "log", String("a"), Int(1), Bool(true)))
	expected := `console.log("a", 1, true)`
	if got != expected {
		t.Errorf("Method() = %q, want %q", got, expected)
	}
}

func TestNew(t *testing.T) {
	got := exprString(New(Ident("Date")))
	expected := `new Date()`
	if got != expected {
		t.Errorf("New() = %q, want %q", got, expected)
	}
}

func TestNewWithArgs(t *testing.T) {
	got := exprString(New(Ident("Date"), Int(2024), Int(0), Int(1)))
	expected := `new Date(2024, 0, 1)`
	if got != expected {
		t.Errorf("New() = %q, want %q", got, expected)
	}
}

func TestOptionalProp(t *testing.T) {
	got := exprString(OptionalProp(Ident("obj"), "foo"))
	expected := `obj?.foo`
	if got != expected {
		t.Errorf("OptionalProp() = %q, want %q", got, expected)
	}
}

func TestOptionalCall(t *testing.T) {
	got := exprString(OptionalCall(Ident("obj"), "method", Int(1)))
	expected := `obj?.method(1)`
	if got != expected {
		t.Errorf("OptionalCall() = %q, want %q", got, expected)
	}
}

// === Operator Tests ===

func TestBinaryOps(t *testing.T) {
	tests := []struct {
		expr     Expr
		expected string
	}{
		{Add(Int(1), Int(2)), "(1 + 2)"},
		{Sub(Int(5), Int(3)), "(5 - 3)"},
		{Mul(Int(4), Int(2)), "(4 * 2)"},
		{Div(Int(10), Int(2)), "(10 / 2)"},
		{Mod(Int(10), Int(3)), "(10 % 3)"},
		{Eq(Ident("x"), Int(5)), "(x === 5)"},
		{NotEq(Ident("x"), Null()), "(x !== null)"},
		{LooseEq(Ident("x"), Int(5)), "(x == 5)"},
		{LooseNotEq(Ident("x"), Null()), "(x != null)"},
		{Lt(Ident("x"), Int(10)), "(x < 10)"},
		{LtEq(Ident("x"), Int(10)), "(x <= 10)"},
		{Gt(Ident("x"), Int(0)), "(x > 0)"},
		{GtEq(Ident("x"), Int(0)), "(x >= 0)"},
		{And(Bool(true), Bool(false)), "(true && false)"},
		{Or(Ident("a"), Ident("b")), "(a || b)"},
		{NullishCoalesce(Ident("x"), String("default")), `(x ?? "default")`},
		{Instanceof(Ident("obj"), Ident("Date")), "(obj instanceof Date)"},
		{In(String("key"), Ident("obj")), `("key" in obj)`},
	}
	for _, tt := range tests {
		got := exprString(tt.expr)
		if got != tt.expected {
			t.Errorf("got %q, want %q", got, tt.expected)
		}
	}
}

func TestUnaryOps(t *testing.T) {
	tests := []struct {
		expr     Expr
		expected string
	}{
		{Not(Ident("x")), "!x"},
		{Neg(Int(5)), "-5"},
		{Pos(Ident("x")), "+x"},
		{BitwiseNot(Int(5)), "~5"},
		{Typeof(Ident("x")), "typeof x"},
		{Void(Int(0)), "void 0"},
		{Delete(Prop(Ident("obj"), "key")), "delete obj.key"},
	}
	for _, tt := range tests {
		got := exprString(tt.expr)
		if got != tt.expected {
			t.Errorf("got %q, want %q", got, tt.expected)
		}
	}
}

func TestTernary(t *testing.T) {
	got := exprString(Ternary(Ident("cond"), String("yes"), String("no")))
	expected := `(cond ? "yes" : "no")`
	if got != expected {
		t.Errorf("Ternary() = %q, want %q", got, expected)
	}
}

func TestSpread(t *testing.T) {
	got := exprString(Spread(Ident("arr")))
	expected := `...arr`
	if got != expected {
		t.Errorf("Spread() = %q, want %q", got, expected)
	}
}

func TestGroup(t *testing.T) {
	got := exprString(Group(Add(Int(1), Int(2))))
	expected := `((1 + 2))`
	if got != expected {
		t.Errorf("Group() = %q, want %q", got, expected)
	}
}

func TestComma(t *testing.T) {
	got := exprString(Comma(Int(1), Int(2), Int(3)))
	expected := `(1, 2, 3)`
	if got != expected {
		t.Errorf("Comma() = %q, want %q", got, expected)
	}
}

// === Statement Tests ===

func TestAssign(t *testing.T) {
	got := stmtString(Assign(Ident("x"), Int(5)))
	if got != "x = 5" {
		t.Errorf("Assign() = %q, want %q", got, "x = 5")
	}
}

func TestCompoundAssign(t *testing.T) {
	tests := []struct {
		stmt     Stmt
		expected string
	}{
		{AddAssign(Ident("x"), Int(1)), "x += 1"},
		{SubAssign(Ident("x"), Int(1)), "x -= 1"},
		{MulAssign(Ident("x"), Int(2)), "x *= 2"},
		{DivAssign(Ident("x"), Int(2)), "x /= 2"},
		{ModAssign(Ident("x"), Int(3)), "x %= 3"},
		{AndAssign(Ident("x"), Ident("y")), "x &&= y"},
		{OrAssign(Ident("x"), Ident("y")), "x ||= y"},
		{NullishAssign(Ident("x"), String("default")), `x ??= "default"`},
	}
	for _, tt := range tests {
		got := stmtString(tt.stmt)
		if got != tt.expected {
			t.Errorf("got %q, want %q", got, tt.expected)
		}
	}
}

func TestLet(t *testing.T) {
	got := stmtString(Let("x", Int(5)))
	if got != "let x = 5" {
		t.Errorf("Let() = %q, want %q", got, "let x = 5")
	}
}

func TestLetDecl(t *testing.T) {
	got := stmtString(LetDecl("x"))
	if got != "let x" {
		t.Errorf("LetDecl() = %q, want %q", got, "let x")
	}
}

func TestConst(t *testing.T) {
	got := stmtString(Const("PI", Float(3.14159)))
	if got != "const PI = 3.14159" {
		t.Errorf("Const() = %q, want %q", got, "const PI = 3.14159")
	}
}

func TestVar(t *testing.T) {
	got := stmtString(Var("x", Int(5)))
	if got != "var x = 5" {
		t.Errorf("Var() = %q, want %q", got, "var x = 5")
	}
}

func TestIncr(t *testing.T) {
	got := stmtString(Incr(Ident("count")))
	if got != "count++" {
		t.Errorf("Incr() = %q, want %q", got, "count++")
	}
}

func TestDecr(t *testing.T) {
	got := stmtString(Decr(Ident("count")))
	if got != "count--" {
		t.Errorf("Decr() = %q, want %q", got, "count--")
	}
}

func TestPreIncr(t *testing.T) {
	got := exprString(PreIncr(Ident("x")))
	if got != "++x" {
		t.Errorf("PreIncr() = %q, want %q", got, "++x")
	}
}

func TestPreDecr(t *testing.T) {
	got := exprString(PreDecr(Ident("x")))
	if got != "--x" {
		t.Errorf("PreDecr() = %q, want %q", got, "--x")
	}
}

func TestReturn(t *testing.T) {
	got := stmtString(Return(Int(42)))
	if got != "return 42" {
		t.Errorf("Return() = %q, want %q", got, "return 42")
	}
}

func TestReturnVoid(t *testing.T) {
	got := stmtString(ReturnVoid())
	if got != "return" {
		t.Errorf("ReturnVoid() = %q, want %q", got, "return")
	}
}

func TestThrow(t *testing.T) {
	got := stmtString(Throw(New(Ident("Error"), String("oops"))))
	expected := `throw new Error("oops")`
	if got != expected {
		t.Errorf("Throw() = %q, want %q", got, expected)
	}
}

func TestBreak(t *testing.T) {
	got := stmtString(Break())
	if got != "break" {
		t.Errorf("Break() = %q, want %q", got, "break")
	}
}

func TestBreakLabel(t *testing.T) {
	got := stmtString(BreakLabel("outer"))
	if got != "break outer" {
		t.Errorf("BreakLabel() = %q, want %q", got, "break outer")
	}
}

func TestContinue(t *testing.T) {
	got := stmtString(Continue())
	if got != "continue" {
		t.Errorf("Continue() = %q, want %q", got, "continue")
	}
}

func TestIf(t *testing.T) {
	got := stmtString(If(Ident("cond"), Assign(Ident("x"), Int(1))))
	expected := "if (cond) { x = 1 }"
	if got != expected {
		t.Errorf("If() = %q, want %q", got, expected)
	}
}

func TestIfMultipleStmts(t *testing.T) {
	got := stmtString(If(Ident("cond"),
		Assign(Ident("x"), Int(1)),
		Incr(Ident("count")),
	))
	expected := "if (cond) { x = 1; count++ }"
	if got != expected {
		t.Errorf("If() = %q, want %q", got, expected)
	}
}

func TestIfElse(t *testing.T) {
	got := stmtString(IfElse(
		Ident("cond"),
		[]Stmt{Assign(Ident("x"), Int(1))},
		[]Stmt{Assign(Ident("x"), Int(0))},
	))
	expected := "if (cond) { x = 1 } else { x = 0 }"
	if got != expected {
		t.Errorf("IfElse() = %q, want %q", got, expected)
	}
}

func TestStmts(t *testing.T) {
	got := stmtString(Stmts(
		Let("x", Int(1)),
		Incr(Ident("x")),
	))
	expected := "let x = 1; x++"
	if got != expected {
		t.Errorf("Stmts() = %q, want %q", got, expected)
	}
}

func TestBlock(t *testing.T) {
	got := stmtString(Block(
		Let("x", Int(1)),
		Incr(Ident("x")),
	))
	expected := "{ let x = 1; x++ }"
	if got != expected {
		t.Errorf("Block() = %q, want %q", got, expected)
	}
}

func TestDebugger(t *testing.T) {
	got := stmtString(Debugger())
	if got != "debugger" {
		t.Errorf("Debugger() = %q, want %q", got, "debugger")
	}
}

func TestExprStmt(t *testing.T) {
	got := stmtString(ExprStmt(ConsoleLog(String("hello"))))
	expected := `console.log("hello")`
	if got != expected {
		t.Errorf("ExprStmt() = %q, want %q", got, expected)
	}
}

// === Handler Tests ===

func TestHandler(t *testing.T) {
	got := Handler(
		ExprStmt(PreventDefault()),
		Assign(Ident("x"), Int(5)),
	)
	expected := "event.preventDefault(); x = 5"
	if got != expected {
		t.Errorf("Handler() = %q, want %q", got, expected)
	}
}

func TestHandlerEmpty(t *testing.T) {
	got := Handler()
	if got != "" {
		t.Errorf("Handler() = %q, want empty", got)
	}
}

func TestExprHandler(t *testing.T) {
	got := ExprHandler(ConsoleLog(String("test")))
	expected := `console.log("test")`
	if got != expected {
		t.Errorf("ExprHandler() = %q, want %q", got, expected)
	}
}

func TestOnClick(t *testing.T) {
	attr := OnClick(ExprStmt(ConsoleLog(String("clicked"))))
	if attr.Name != "onclick" {
		t.Errorf("OnClick().Name = %q, want %q", attr.Name, "onclick")
	}
	expected := `console.log("clicked")`
	if attr.Value != expected {
		t.Errorf("OnClick().Value = %q, want %q", attr.Value, expected)
	}
}

func TestOnInput(t *testing.T) {
	attr := OnInput(Assign(Prop(EventTarget(), "value"), Method(EventValue(), "toUpperCase")))
	if attr.Name != "oninput" {
		t.Errorf("OnInput().Name = %q, want %q", attr.Name, "oninput")
	}
	expected := `event.target.value = event.target.value.toUpperCase()`
	if attr.Value != expected {
		t.Errorf("OnInput().Value = %q, want %q", attr.Value, expected)
	}
}

func TestOnSubmit(t *testing.T) {
	attr := OnSubmit(
		ExprStmt(PreventDefault()),
		ExprStmt(Method(This(), "submit")),
	)
	if attr.Name != "onsubmit" {
		t.Errorf("OnSubmit().Name = %q, want %q", attr.Name, "onsubmit")
	}
	expected := `event.preventDefault(); this.submit()`
	if attr.Value != expected {
		t.Errorf("OnSubmit().Value = %q, want %q", attr.Value, expected)
	}
}

func TestOn(t *testing.T) {
	attr := On("touchstart", ExprStmt(ConsoleLog(String("touched"))))
	if attr.Name != "ontouchstart" {
		t.Errorf("On().Name = %q, want %q", attr.Name, "ontouchstart")
	}
}

// === Builtin Tests ===

func TestConsoleLog(t *testing.T) {
	got := exprString(ConsoleLog(String("msg"), Int(42)))
	expected := `console.log("msg", 42)`
	if got != expected {
		t.Errorf("ConsoleLog() = %q, want %q", got, expected)
	}
}

func TestConsoleError(t *testing.T) {
	got := exprString(ConsoleError(String("error")))
	expected := `console.error("error")`
	if got != expected {
		t.Errorf("ConsoleError() = %q, want %q", got, expected)
	}
}

func TestGetElementById(t *testing.T) {
	got := exprString(GetElementById(String("myId")))
	expected := `document.getElementById("myId")`
	if got != expected {
		t.Errorf("GetElementById() = %q, want %q", got, expected)
	}
}

func TestQuerySelector(t *testing.T) {
	got := exprString(QuerySelector(String(".myClass")))
	expected := `document.querySelector(".myClass")`
	if got != expected {
		t.Errorf("QuerySelector() = %q, want %q", got, expected)
	}
}

func TestAlert(t *testing.T) {
	got := exprString(Alert(String("hello")))
	expected := `alert("hello")`
	if got != expected {
		t.Errorf("Alert() = %q, want %q", got, expected)
	}
}

func TestPreventDefault(t *testing.T) {
	got := exprString(PreventDefault())
	expected := "event.preventDefault()"
	if got != expected {
		t.Errorf("PreventDefault() = %q, want %q", got, expected)
	}
}

func TestStopPropagation(t *testing.T) {
	got := exprString(StopPropagation())
	expected := "event.stopPropagation()"
	if got != expected {
		t.Errorf("StopPropagation() = %q, want %q", got, expected)
	}
}

func TestEventValue(t *testing.T) {
	got := exprString(EventValue())
	expected := "event.target.value"
	if got != expected {
		t.Errorf("EventValue() = %q, want %q", got, expected)
	}
}

func TestEventChecked(t *testing.T) {
	got := exprString(EventChecked())
	expected := "event.target.checked"
	if got != expected {
		t.Errorf("EventChecked() = %q, want %q", got, expected)
	}
}

func TestNavigate(t *testing.T) {
	got := stmtString(Navigate(String("/home")))
	expected := `location.href = "/home"`
	if got != expected {
		t.Errorf("Navigate() = %q, want %q", got, expected)
	}
}

func TestReload(t *testing.T) {
	got := exprString(Reload())
	expected := "location.reload()"
	if got != expected {
		t.Errorf("Reload() = %q, want %q", got, expected)
	}
}

func TestHistoryBack(t *testing.T) {
	got := exprString(HistoryBack())
	expected := "history.back()"
	if got != expected {
		t.Errorf("HistoryBack() = %q, want %q", got, expected)
	}
}

func TestSetTimeout(t *testing.T) {
	got := exprString(SetTimeout(ArrowFunc(nil, ConsoleLog(String("delayed"))), Int(1000)))
	expected := `setTimeout(() => console.log("delayed"), 1000)`
	if got != expected {
		t.Errorf("SetTimeout() = %q, want %q", got, expected)
	}
}

func TestFetch(t *testing.T) {
	got := exprString(Fetch(String("/api/data")))
	expected := `fetch("/api/data")`
	if got != expected {
		t.Errorf("Fetch() = %q, want %q", got, expected)
	}
}

func TestClassListAdd(t *testing.T) {
	got := exprString(ClassListAdd(Ident("el"), String("active"), String("visible")))
	expected := `el.classList.add("active", "visible")`
	if got != expected {
		t.Errorf("ClassListAdd() = %q, want %q", got, expected)
	}
}

func TestSetStyle(t *testing.T) {
	got := stmtString(SetStyle(Ident("el"), "backgroundColor", String("red")))
	expected := `el.style.backgroundColor = "red"`
	if got != expected {
		t.Errorf("SetStyle() = %q, want %q", got, expected)
	}
}

func TestJSONStringify(t *testing.T) {
	got := exprString(JSONStringify(Ident("obj")))
	expected := "JSON.stringify(obj)"
	if got != expected {
		t.Errorf("JSONStringify() = %q, want %q", got, expected)
	}
}

func TestJSONParse(t *testing.T) {
	got := exprString(JSONParse(String(`{"a":1}`)))
	expected := `JSON.parse("{\"a\":1}")`
	if got != expected {
		t.Errorf("JSONParse() = %q, want %q", got, expected)
	}
}

// === Arrow Function Tests ===

func TestArrowFunc(t *testing.T) {
	got := exprString(ArrowFunc([]string{"x"}, Mul(Ident("x"), Int(2))))
	expected := "x => (x * 2)"
	if got != expected {
		t.Errorf("ArrowFunc() = %q, want %q", got, expected)
	}
}

func TestArrowFuncMultiParams(t *testing.T) {
	got := exprString(ArrowFunc([]string{"a", "b"}, Add(Ident("a"), Ident("b"))))
	expected := "(a, b) => (a + b)"
	if got != expected {
		t.Errorf("ArrowFunc() = %q, want %q", got, expected)
	}
}

func TestArrowFuncNoParams(t *testing.T) {
	got := exprString(ArrowFunc(nil, String("hello")))
	expected := `() => "hello"`
	if got != expected {
		t.Errorf("ArrowFunc() = %q, want %q", got, expected)
	}
}

func TestArrowFuncStmts(t *testing.T) {
	got := exprString(ArrowFuncStmts([]string{"x"},
		Let("result", Mul(Ident("x"), Int(2))),
		Return(Ident("result")),
	))
	expected := "x => { let result = (x * 2); return result }"
	if got != expected {
		t.Errorf("ArrowFuncStmts() = %q, want %q", got, expected)
	}
}

func TestFunc(t *testing.T) {
	got := exprString(Func([]string{"x", "y"}, Return(Add(Ident("x"), Ident("y")))))
	expected := "function(x, y) { return (x + y) }"
	if got != expected {
		t.Errorf("Func() = %q, want %q", got, expected)
	}
}

func TestIIFE(t *testing.T) {
	got := exprString(IIFE(ExprStmt(ConsoleLog(String("hello")))))
	expected := `(function() { console.log("hello") })()`
	if got != expected {
		t.Errorf("IIFE() = %q, want %q", got, expected)
	}
}

func TestTemplate(t *testing.T) {
	got := exprString(Template("Hello, ", Ident("name"), "!"))
	expected := "`Hello, ${name}!`"
	if got != expected {
		t.Errorf("Template() = %q, want %q", got, expected)
	}
}

func TestTemplateEscaping(t *testing.T) {
	got := exprString(Template("Price: $", Ident("amount")))
	expected := "`Price: \\$${amount}`"
	if got != expected {
		t.Errorf("Template() = %q, want %q", got, expected)
	}
}

func TestAwait(t *testing.T) {
	got := exprString(Await(Fetch(String("/api"))))
	expected := `await fetch("/api")`
	if got != expected {
		t.Errorf("Await() = %q, want %q", got, expected)
	}
}

func TestAsyncArrowFunc(t *testing.T) {
	got := exprString(AsyncArrowFunc(nil, Await(Fetch(String("/api")))))
	expected := `async () => await fetch("/api")`
	if got != expected {
		t.Errorf("AsyncArrowFunc() = %q, want %q", got, expected)
	}
}

func TestAsyncArrowFuncStmts(t *testing.T) {
	got := exprString(AsyncArrowFuncStmts([]string{"url"},
		Let("response", Await(Fetch(Ident("url")))),
		Return(Method(Ident("response"), "json")),
	))
	expected := "async url => { let response = await fetch(url); return response.json() }"
	if got != expected {
		t.Errorf("AsyncArrowFuncStmts() = %q, want %q", got, expected)
	}
}

func TestPromiseThen(t *testing.T) {
	got := exprString(PromiseThen(Fetch(String("/api")), ArrowFunc([]string{"r"}, Method(Ident("r"), "json"))))
	expected := `fetch("/api").then(r => r.json())`
	if got != expected {
		t.Errorf("PromiseThen() = %q, want %q", got, expected)
	}
}

// === ToJS Tests ===

func TestToJS(t *testing.T) {
	got := ToJS(Add(Int(1), Int(2)))
	expected := "(1 + 2)"
	if got != expected {
		t.Errorf("ToJS() = %q, want %q", got, expected)
	}
}

func TestToJSStmt(t *testing.T) {
	got := ToJSStmt(Let("x", Int(5)))
	expected := "let x = 5"
	if got != expected {
		t.Errorf("ToJSStmt() = %q, want %q", got, expected)
	}
}

// === Integration Tests ===

func TestComplexHandler(t *testing.T) {
	// A more complex real-world example
	handler := Handler(
		ExprStmt(PreventDefault()),
		Let("value", EventValue()),
		If(Eq(Ident("value"), String("")),
			Return(Null()),
		),
		ExprStmt(ClassListAdd(This(), String("submitted"))),
		ExprStmt(ConsoleLog(String("Submitted:"), Ident("value"))),
	)
	// Just check it doesn't panic and produces reasonable output
	if handler == "" {
		t.Error("Complex handler produced empty string")
	}
	if !strings.Contains(handler, "event.preventDefault()") {
		t.Error("Handler missing preventDefault")
	}
	if !strings.Contains(handler, "let value = event.target.value") {
		t.Error("Handler missing variable declaration")
	}
}
