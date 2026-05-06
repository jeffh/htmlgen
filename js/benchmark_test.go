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
		obj.Prop("body")
	}
}

func BenchmarkPropChain(b *testing.B) {
	for b.Loop() {
		Ident("window").Prop("document").Prop("body").Prop("style").Prop("display")
	}
}

func BenchmarkIndex(b *testing.B) {
	arr := Ident("arr")
	for b.Loop() {
		arr.Index(Int(0))
	}
}

func BenchmarkCall_NoArgs(b *testing.B) {
	fn := Ident("doSomething")
	for b.Loop() {
		fn.Call()
	}
}

func BenchmarkCall_WithArgs(b *testing.B) {
	fn := Ident("doSomething")
	for b.Loop() {
		fn.Call(String("arg1"), Int(42), Bool(true))
	}
}

func BenchmarkMethod(b *testing.B) {
	obj := Ident("console")
	for b.Loop() {
		obj.Method("log", String("message"))
	}
}

func BenchmarkNew(b *testing.B) {
	for b.Loop() {
		Ident("Date").New(Int(2024), Int(0), Int(1))
	}
}

func BenchmarkOptionalProp(b *testing.B) {
	obj := Ident("obj")
	for b.Loop() {
		obj.OptionalProp("foo")
	}
}

func BenchmarkOptionalCall(b *testing.B) {
	obj := Ident("obj")
	for b.Loop() {
		obj.OptionalCall("method", Int(1))
	}
}

// ============================================================================
// Operator Benchmarks
// ============================================================================

func BenchmarkBinaryOp_Simple(b *testing.B) {
	for b.Loop() {
		Int(1).Add(Int(2))
	}
}

func BenchmarkBinaryOp_Nested(b *testing.B) {
	for b.Loop() {
		Int(2).Mul(Int(3)).Add(Int(10).Div(Int(2)))
	}
}

func BenchmarkComparison(b *testing.B) {
	x := Ident("x")
	for b.Loop() {
		x.GtEq(Int(0)).And(x.Lt(Int(100)))
	}
}

func BenchmarkTernary(b *testing.B) {
	for b.Loop() {
		Ident("cond").Ternary(String("yes"), String("no"))
	}
}

func BenchmarkUnaryOp(b *testing.B) {
	x := Ident("x")
	for b.Loop() {
		x.Not()
	}
}

func BenchmarkNullishCoalesce(b *testing.B) {
	x := Ident("x")
	for b.Loop() {
		x.NullishCoalesce(String("default"))
	}
}

// ============================================================================
// Statement Benchmarks
// ============================================================================

func BenchmarkAssign(b *testing.B) {
	x := Ident("x")
	for b.Loop() {
		x.Assign(Int(5))
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
		If(Ident("cond"), Ident("x").Assign(Int(1)))
	}
}

func BenchmarkIfElse(b *testing.B) {
	for b.Loop() {
		IfElse(
			Ident("cond"),
			[]Stmt{Ident("x").Assign(Int(1))},
			[]Stmt{Ident("x").Assign(Int(0))},
		)
	}
}

func BenchmarkBlock(b *testing.B) {
	for b.Loop() {
		Block(
			Let("x", Int(1)),
			Ident("x").Incr(),
			Return(Ident("x")),
		)
	}
}

func BenchmarkStmts(b *testing.B) {
	for b.Loop() {
		Stmts(
			Let("a", Int(1)),
			Let("b", Int(2)),
			Ident("c").Assign(Ident("a").Add(Ident("b"))),
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
		ArrowFunc([]string{"x"}, Ident("x").Mul(Int(2)))
	}
}

func BenchmarkArrowFunc_MultiParams(b *testing.B) {
	for b.Loop() {
		ArrowFunc([]string{"a", "b", "c"}, Ident("a").Add(Ident("b")).Add(Ident("c")))
	}
}

func BenchmarkArrowFuncStmts(b *testing.B) {
	for b.Loop() {
		ArrowFuncStmts([]string{"x"},
			Let("result", Ident("x").Mul(Int(2))),
			Return(Ident("result")),
		)
	}
}

func BenchmarkFunc(b *testing.B) {
	for b.Loop() {
		Func([]string{"x", "y"}, Return(Ident("x").Add(Ident("y"))))
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
		fetch.Await()
	}
}

func BenchmarkAsyncArrowFunc(b *testing.B) {
	for b.Loop() {
		AsyncArrowFunc(nil, Fetch(String("/api")).Await())
	}
}

func BenchmarkAsyncArrowFuncStmts(b *testing.B) {
	for b.Loop() {
		AsyncArrowFuncStmts([]string{"url"},
			Let("response", Fetch(Ident("url")).Await()),
			Return(Ident("response").Method("json")),
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
			If(Ident("value").Eq(String("")),
				Return(Null()),
			),
			ExprStmt(This().ClassListAdd(String("submitted"))),
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
	expr := Ident("x").Mul(Int(2)).Add(Int(1))
	for b.Loop() {
		ToJS(expr)
	}
}

func BenchmarkToJSStmt(b *testing.B) {
	stmt := Let("x", Int(1).Add(Int(2)))
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
			ExprStmt(This().Method("submit")),
		)
	}
}

func BenchmarkOnInput(b *testing.B) {
	for b.Loop() {
		OnInput(EventTarget().Prop("value").Assign(EventValue().Method("toUpperCase")))
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
		el.ClassListAdd(String("active"), String("visible"))
	}
}

func BenchmarkClassListToggle(b *testing.B) {
	el := Ident("el")
	for b.Loop() {
		el.ClassListToggle(String("active"))
	}
}

func BenchmarkClassListToggle_WithForce(b *testing.B) {
	el := Ident("el")
	for b.Loop() {
		el.ClassListToggle(String("active"), Bool(true))
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
		el.SetStyle("backgroundColor", String("red"))
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
			If(Ident("value").Eq(String("")),
				Return(Null()),
			),
			ExprStmt(This().ClassListAdd(String("submitted"))),
		)
	}
}

func BenchmarkAllocations_ArrowFunc(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		ArrowFunc([]string{"x"}, Ident("x").Mul(Int(2)))
	}
}

func BenchmarkAllocations_Template(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		Template("Hello, ", Ident("name"), "!")
	}
}
