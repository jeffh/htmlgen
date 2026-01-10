package js

import (
	"testing"
)

// ============================================================================
// Value Creation Benchmarks
// ============================================================================

func BenchmarkString(b *testing.B) {
	for b.Loop() {
		String("Hello, World!")
	}
}

func BenchmarkString_WithEscaping(b *testing.B) {
	for b.Loop() {
		String(`<script>alert("XSS")</script>`)
	}
}

func BenchmarkInt(b *testing.B) {
	for b.Loop() {
		Int(42)
	}
}

func BenchmarkInt64(b *testing.B) {
	for b.Loop() {
		Int64(9223372036854775807)
	}
}

func BenchmarkFloat(b *testing.B) {
	for b.Loop() {
		Float(3.14159265359)
	}
}

func BenchmarkBool(b *testing.B) {
	for b.Loop() {
		Bool(true)
	}
}

func BenchmarkJSON_Simple(b *testing.B) {
	for b.Loop() {
		JSON("hello")
	}
}

func BenchmarkJSON_Complex(b *testing.B) {
	data := map[string]any{
		"name":  "John",
		"age":   30,
		"items": []int{1, 2, 3},
	}
	for b.Loop() {
		JSON(data)
	}
}

func BenchmarkArray_Small(b *testing.B) {
	for b.Loop() {
		Array(Int(1), Int(2), Int(3))
	}
}

func BenchmarkArray_Medium(b *testing.B) {
	elements := make([]Expr, 10)
	for i := range 10 {
		elements[i] = Int(i)
	}
	for b.Loop() {
		Array(elements...)
	}
}

func BenchmarkObject_Small(b *testing.B) {
	for b.Loop() {
		Object(
			Pair("name", String("John")),
			Pair("age", Int(30)),
		)
	}
}

func BenchmarkObject_Medium(b *testing.B) {
	for b.Loop() {
		Object(
			Pair("name", String("John")),
			Pair("email", String("john@example.com")),
			Pair("age", Int(30)),
			Pair("active", Bool(true)),
			Pair("score", Float(95.5)),
		)
	}
}

func BenchmarkIdent(b *testing.B) {
	for b.Loop() {
		Ident("myVariable")
	}
}

// ============================================================================
// Expression Building Benchmarks
// ============================================================================

func BenchmarkProp(b *testing.B) {
	obj := Ident("document")
	for b.Loop() {
		Prop(obj, "body")
	}
}

func BenchmarkPropChain(b *testing.B) {
	for b.Loop() {
		Prop(Prop(Prop(Prop(Ident("window"), "document"), "body"), "style"), "display")
	}
}

func BenchmarkIndex(b *testing.B) {
	arr := Ident("arr")
	for b.Loop() {
		Index(arr, Int(0))
	}
}

func BenchmarkCall_NoArgs(b *testing.B) {
	fn := Ident("doSomething")
	for b.Loop() {
		Call(fn)
	}
}

func BenchmarkCall_WithArgs(b *testing.B) {
	fn := Ident("doSomething")
	for b.Loop() {
		Call(fn, String("arg1"), Int(42), Bool(true))
	}
}

func BenchmarkMethod(b *testing.B) {
	obj := Ident("console")
	for b.Loop() {
		Method(obj, "log", String("message"))
	}
}

func BenchmarkNew(b *testing.B) {
	for b.Loop() {
		New(Ident("Date"), Int(2024), Int(0), Int(1))
	}
}

func BenchmarkOptionalProp(b *testing.B) {
	obj := Ident("obj")
	for b.Loop() {
		OptionalProp(obj, "foo")
	}
}

func BenchmarkOptionalCall(b *testing.B) {
	obj := Ident("obj")
	for b.Loop() {
		OptionalCall(obj, "method", Int(1))
	}
}

// ============================================================================
// Operator Benchmarks
// ============================================================================

func BenchmarkBinaryOp_Simple(b *testing.B) {
	for b.Loop() {
		Add(Int(1), Int(2))
	}
}

func BenchmarkBinaryOp_Nested(b *testing.B) {
	for b.Loop() {
		Add(Mul(Int(2), Int(3)), Div(Int(10), Int(2)))
	}
}

func BenchmarkComparison(b *testing.B) {
	x := Ident("x")
	for b.Loop() {
		And(GtEq(x, Int(0)), Lt(x, Int(100)))
	}
}

func BenchmarkTernary(b *testing.B) {
	for b.Loop() {
		Ternary(Ident("cond"), String("yes"), String("no"))
	}
}

func BenchmarkUnaryOp(b *testing.B) {
	x := Ident("x")
	for b.Loop() {
		Not(x)
	}
}

func BenchmarkNullishCoalesce(b *testing.B) {
	x := Ident("x")
	for b.Loop() {
		NullishCoalesce(x, String("default"))
	}
}

// ============================================================================
// Statement Benchmarks
// ============================================================================

func BenchmarkAssign(b *testing.B) {
	x := Ident("x")
	for b.Loop() {
		Assign(x, Int(5))
	}
}

func BenchmarkLet(b *testing.B) {
	for b.Loop() {
		Let("x", Int(5))
	}
}

func BenchmarkConst(b *testing.B) {
	for b.Loop() {
		Const("PI", Float(3.14159))
	}
}

func BenchmarkIf(b *testing.B) {
	for b.Loop() {
		If(Ident("cond"), Assign(Ident("x"), Int(1)))
	}
}

func BenchmarkIfElse(b *testing.B) {
	for b.Loop() {
		IfElse(
			Ident("cond"),
			[]Stmt{Assign(Ident("x"), Int(1))},
			[]Stmt{Assign(Ident("x"), Int(0))},
		)
	}
}

func BenchmarkBlock(b *testing.B) {
	for b.Loop() {
		Block(
			Let("x", Int(1)),
			Incr(Ident("x")),
			Return(Ident("x")),
		)
	}
}

func BenchmarkStmts(b *testing.B) {
	for b.Loop() {
		Stmts(
			Let("a", Int(1)),
			Let("b", Int(2)),
			Assign(Ident("c"), Add(Ident("a"), Ident("b"))),
		)
	}
}

// ============================================================================
// Function Expression Benchmarks
// ============================================================================

func BenchmarkArrowFunc_NoParams(b *testing.B) {
	for b.Loop() {
		ArrowFunc(nil, String("hello"))
	}
}

func BenchmarkArrowFunc_OneParam(b *testing.B) {
	for b.Loop() {
		ArrowFunc([]string{"x"}, Mul(Ident("x"), Int(2)))
	}
}

func BenchmarkArrowFunc_MultiParams(b *testing.B) {
	for b.Loop() {
		ArrowFunc([]string{"a", "b", "c"}, Add(Add(Ident("a"), Ident("b")), Ident("c")))
	}
}

func BenchmarkArrowFuncStmts(b *testing.B) {
	for b.Loop() {
		ArrowFuncStmts([]string{"x"},
			Let("result", Mul(Ident("x"), Int(2))),
			Return(Ident("result")),
		)
	}
}

func BenchmarkFunc(b *testing.B) {
	for b.Loop() {
		Func([]string{"x", "y"}, Return(Add(Ident("x"), Ident("y"))))
	}
}

func BenchmarkIIFE(b *testing.B) {
	for b.Loop() {
		IIFE(
			Let("x", Int(1)),
			Return(Ident("x")),
		)
	}
}

func BenchmarkTemplate_Simple(b *testing.B) {
	for b.Loop() {
		Template("Hello, ", Ident("name"), "!")
	}
}

func BenchmarkTemplate_Complex(b *testing.B) {
	for b.Loop() {
		Template("User: ", Ident("user"), " (", Ident("role"), ") - Score: ", Ident("score"))
	}
}

func BenchmarkAwait(b *testing.B) {
	fetch := Fetch(String("/api"))
	for b.Loop() {
		Await(fetch)
	}
}

func BenchmarkAsyncArrowFunc(b *testing.B) {
	for b.Loop() {
		AsyncArrowFunc(nil, Await(Fetch(String("/api"))))
	}
}

func BenchmarkAsyncArrowFuncStmts(b *testing.B) {
	for b.Loop() {
		AsyncArrowFuncStmts([]string{"url"},
			Let("response", Await(Fetch(Ident("url")))),
			Return(Method(Ident("response"), "json")),
		)
	}
}

// ============================================================================
// Handler Generation Benchmarks (End-to-End)
// ============================================================================

func BenchmarkHandler_Empty(b *testing.B) {
	for b.Loop() {
		Handler()
	}
}

func BenchmarkHandler_Simple(b *testing.B) {
	for b.Loop() {
		Handler(ExprStmt(ConsoleLog(String("clicked"))))
	}
}

func BenchmarkHandler_TwoStatements(b *testing.B) {
	for b.Loop() {
		Handler(
			ExprStmt(PreventDefault()),
			ExprStmt(ConsoleLog(String("clicked"))),
		)
	}
}

func BenchmarkHandler_Complex(b *testing.B) {
	for b.Loop() {
		Handler(
			ExprStmt(PreventDefault()),
			Let("value", EventValue()),
			If(Eq(Ident("value"), String("")),
				Return(Null()),
			),
			ExprStmt(ClassListAdd(This(), String("submitted"))),
			ExprStmt(ConsoleLog(String("Submitted:"), Ident("value"))),
		)
	}
}

func BenchmarkExprHandler(b *testing.B) {
	for b.Loop() {
		ExprHandler(ConsoleLog(String("test")))
	}
}

func BenchmarkToJS(b *testing.B) {
	expr := Add(Mul(Ident("x"), Int(2)), Int(1))
	for b.Loop() {
		ToJS(expr)
	}
}

func BenchmarkToJSStmt(b *testing.B) {
	stmt := Let("x", Add(Int(1), Int(2)))
	for b.Loop() {
		ToJSStmt(stmt)
	}
}

// ============================================================================
// Event Handler Attribute Benchmarks
// ============================================================================

func BenchmarkOnClick(b *testing.B) {
	for b.Loop() {
		OnClick(ExprStmt(ConsoleLog(String("clicked"))))
	}
}

func BenchmarkOnSubmit(b *testing.B) {
	for b.Loop() {
		OnSubmit(
			ExprStmt(PreventDefault()),
			ExprStmt(Method(This(), "submit")),
		)
	}
}

func BenchmarkOnInput(b *testing.B) {
	for b.Loop() {
		OnInput(Assign(Prop(EventTarget(), "value"), Method(EventValue(), "toUpperCase")))
	}
}

// ============================================================================
// Builtin Helper Benchmarks
// ============================================================================

func BenchmarkConsoleLog(b *testing.B) {
	for b.Loop() {
		ConsoleLog(String("message"), Int(42))
	}
}

func BenchmarkGetElementById(b *testing.B) {
	for b.Loop() {
		GetElementById(String("myId"))
	}
}

func BenchmarkQuerySelector(b *testing.B) {
	for b.Loop() {
		QuerySelector(String(".myClass"))
	}
}

func BenchmarkFetch_NoOptions(b *testing.B) {
	for b.Loop() {
		Fetch(String("/api/data"))
	}
}

func BenchmarkFetch_WithOptions(b *testing.B) {
	for b.Loop() {
		Fetch(String("/api/data"), Object(
			Pair("method", String("POST")),
			Pair("body", JSONStringify(Ident("data"))),
		))
	}
}

func BenchmarkClassListAdd(b *testing.B) {
	el := Ident("el")
	for b.Loop() {
		ClassListAdd(el, String("active"), String("visible"))
	}
}

func BenchmarkClassListToggle(b *testing.B) {
	el := Ident("el")
	for b.Loop() {
		ClassListToggle(el, String("active"))
	}
}

func BenchmarkClassListToggle_WithForce(b *testing.B) {
	el := Ident("el")
	for b.Loop() {
		ClassListToggle(el, String("active"), Bool(true))
	}
}

func BenchmarkJSONStringify(b *testing.B) {
	obj := Ident("obj")
	for b.Loop() {
		JSONStringify(obj)
	}
}

func BenchmarkJSONParse(b *testing.B) {
	for b.Loop() {
		JSONParse(String(`{"a":1}`))
	}
}

func BenchmarkSetStyle(b *testing.B) {
	el := Ident("el")
	for b.Loop() {
		SetStyle(el, "backgroundColor", String("red"))
	}
}

// ============================================================================
// Allocation Benchmarks (using ReportAllocs)
// ============================================================================

func BenchmarkAllocations_String(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		String("Hello, World!")
	}
}

func BenchmarkAllocations_Object(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		Object(
			Pair("name", String("John")),
			Pair("age", Int(30)),
		)
	}
}

func BenchmarkAllocations_Handler(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		Handler(
			ExprStmt(PreventDefault()),
			ExprStmt(ConsoleLog(String("clicked"))),
		)
	}
}

func BenchmarkAllocations_ComplexHandler(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		Handler(
			ExprStmt(PreventDefault()),
			Let("value", EventValue()),
			If(Eq(Ident("value"), String("")),
				Return(Null()),
			),
			ExprStmt(ClassListAdd(This(), String("submitted"))),
		)
	}
}

func BenchmarkAllocations_ArrowFunc(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		ArrowFunc([]string{"x"}, Mul(Ident("x"), Int(2)))
	}
}

func BenchmarkAllocations_Template(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		Template("Hello, ", Ident("name"), "!")
	}
}
