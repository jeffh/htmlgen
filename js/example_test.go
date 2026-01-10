package js_test

import (
	"fmt"

	"github.com/jeffh/htmlgen/js"
)

func Example_basicClickHandler() {
	handler := js.Handler(
		js.ExprStmt(js.Alert(js.String("Hello, World!"))),
	)
	fmt.Println(handler)
	// Output: alert("Hello, World!")
}

func Example_preventDefault() {
	handler := js.Handler(
		js.ExprStmt(js.PreventDefault()),
		js.ExprStmt(js.ConsoleLog(js.String("Form submitted"))),
	)
	fmt.Println(handler)
	// Output: event.preventDefault(); console.log("Form submitted")
}

func Example_eventValue() {
	// Transform input to uppercase as the user types
	handler := js.Handler(
		js.Assign(
			js.Prop(js.EventTarget(), "value"),
			js.Method(js.EventValue(), "toUpperCase"),
		),
	)
	fmt.Println(handler)
	// Output: event.target.value = event.target.value.toUpperCase()
}

func Example_conditionalLogic() {
	// Show alert if input is empty
	handler := js.Handler(
		js.If(js.Eq(js.EventValue(), js.String("")),
			js.ExprStmt(js.Alert(js.String("Please enter a value"))),
			js.ReturnVoid(),
		),
	)
	fmt.Println(handler)
	// Output: if ((event.target.value === "")) { alert("Please enter a value"); return }
}

func Example_toggleClass() {
	// Toggle a CSS class on an element
	handler := js.Handler(
		js.ExprStmt(js.ClassListToggle(
			js.GetElementById(js.String("menu")),
			js.String("open"),
		)),
	)
	fmt.Println(handler)
	// Output: document.getElementById("menu").classList.toggle("open")
}

func Example_arrowFunction() {
	// Create an arrow function for setTimeout
	expr := js.SetTimeout(
		js.ArrowFunc(nil, js.ConsoleLog(js.String("Delayed!"))),
		js.Int(1000),
	)
	fmt.Println(js.ExprHandler(expr))
	// Output: setTimeout(() => console.log("Delayed!"), 1000)
}

func Example_ternaryOperator() {
	// Use ternary for conditional text
	expr := js.Ternary(
		js.Ident("isLoggedIn"),
		js.String("Logout"),
		js.String("Login"),
	)
	fmt.Println(js.ExprHandler(expr))
	// Output: (isLoggedIn ? "Logout" : "Login")
}

func Example_objectLiteral() {
	// Create a JavaScript object
	obj := js.Object(
		js.Pair("name", js.String("John")),
		js.Pair("age", js.Int(30)),
		js.Pair("active", js.Bool(true)),
	)
	fmt.Println(js.ExprHandler(obj))
	// Output: {"name": "John", "age": 30, "active": true}
}

func Example_arrayMethods() {
	// Chain array methods
	expr := js.Method(
		js.Method(js.Ident("items"), "filter", js.ArrowFunc([]string{"x"}, js.Prop(js.Ident("x"), "active"))),
		"map",
		js.ArrowFunc([]string{"x"}, js.Prop(js.Ident("x"), "name")),
	)
	fmt.Println(js.ExprHandler(expr))
	// Output: items.filter(x => x.active).map(x => x.name)
}

func Example_promiseChain() {
	// Fetch with promise chain
	expr := js.PromiseCatch(
		js.PromiseThen(
			js.Fetch(js.String("/api/users")),
			js.ArrowFunc([]string{"r"}, js.Method(js.Ident("r"), "json")),
		),
		js.ArrowFunc([]string{"e"}, js.ConsoleError(js.Ident("e"))),
	)
	fmt.Println(js.ExprHandler(expr))
	// Output: fetch("/api/users").then(r => r.json()).catch(e => console.error(e))
}

func Example_templateLiteral() {
	// Create a template literal with interpolation
	expr := js.Template("Hello, ", js.Ident("name"), "! You have ", js.Ident("count"), " messages.")
	fmt.Println(js.ExprHandler(expr))
	// Output: `Hello, ${name}! You have ${count} messages.`
}

func Example_rawEscapeHatch() {
	// Use Raw for custom JavaScript
	handler := js.Handler(
		js.ExprStmt(js.Raw("gtag('event', 'button_click')")),
	)
	fmt.Println(handler)
	// Output: gtag('event', 'button_click')
}

func Example_keyboardShortcut() {
	// Handle Ctrl+S keyboard shortcut
	handler := js.Handler(
		js.If(js.And(js.EventCtrlKey(), js.Eq(js.EventKey(), js.String("s"))),
			js.ExprStmt(js.PreventDefault()),
			js.ExprStmt(js.ConsoleLog(js.String("Save triggered"))),
		),
	)
	fmt.Println(handler)
	// Output: if ((event.ctrlKey && (event.key === "s"))) { event.preventDefault(); console.log("Save triggered") }
}

func ExampleOnClick() {
	attr := js.OnClick(js.ExprStmt(js.ConsoleLog(js.String("clicked"))))
	fmt.Printf("%s=%q", attr.Name, attr.Value)
	// Output: onclick="console.log(\"clicked\")"
}

func ExampleOnInput() {
	attr := js.OnInput(
		js.Assign(js.Ident("value"), js.EventValue()),
	)
	fmt.Printf("%s=%q", attr.Name, attr.Value)
	// Output: oninput="value = event.target.value"
}

func ExampleOnSubmit() {
	attr := js.OnSubmit(
		js.ExprStmt(js.PreventDefault()),
		js.ExprStmt(js.Method(js.This(), "submit")),
	)
	fmt.Printf("%s=%q", attr.Name, attr.Value)
	// Output: onsubmit="event.preventDefault(); this.submit()"
}

func ExampleRaw() {
	// Raw is the escape hatch for arbitrary JavaScript
	expr := js.Raw("customLibrary.doSomething()")
	fmt.Println(js.ExprHandler(expr))
	// Output: customLibrary.doSomething()
}

func ExampleString() {
	// Strings are JSON-encoded for safety
	expr := js.String(`He said "hello"`)
	fmt.Println(js.ExprHandler(expr))
	// Output: "He said \"hello\""
}

func ExampleJSON() {
	// JSON encodes any Go value
	expr := js.JSON(map[string]any{
		"users": []string{"alice", "bob"},
	})
	fmt.Println(js.ExprHandler(expr))
	// Output: {"users":["alice","bob"]}
}

func ExampleHandler() {
	// Combine multiple statements
	handler := js.Handler(
		js.Let("x", js.Int(1)),
		js.Incr(js.Ident("x")),
		js.ExprStmt(js.ConsoleLog(js.Ident("x"))),
	)
	fmt.Println(handler)
	// Output: let x = 1; x++; console.log(x)
}

func ExampleExprStmt() {
	// Convert an expression to a statement
	stmt := js.ExprStmt(js.ConsoleLog(js.String("hello")))
	fmt.Println(js.ToJSStmt(stmt))
	// Output: console.log("hello")
}

func ExampleProp() {
	// Access nested properties
	expr := js.Prop(js.Prop(js.Ident("window"), "location"), "href")
	fmt.Println(js.ExprHandler(expr))
	// Output: window.location.href
}

func ExampleMethod() {
	// Call a method with arguments
	expr := js.Method(js.Ident("arr"), "push", js.Int(1), js.Int(2))
	fmt.Println(js.ExprHandler(expr))
	// Output: arr.push(1, 2)
}

func ExampleArrowFunc() {
	// Single expression arrow function
	fn := js.ArrowFunc([]string{"a", "b"}, js.Add(js.Ident("a"), js.Ident("b")))
	fmt.Println(js.ExprHandler(fn))
	// Output: (a, b) => (a + b)
}

func ExampleArrowFuncStmts() {
	// Multi-statement arrow function
	fn := js.ArrowFuncStmts([]string{"x"},
		js.Let("result", js.Mul(js.Ident("x"), js.Int(2))),
		js.Return(js.Ident("result")),
	)
	fmt.Println(js.ExprHandler(fn))
	// Output: x => { let result = (x * 2); return result }
}

func ExampleIf() {
	stmt := js.If(js.Gt(js.Ident("x"), js.Int(0)),
		js.ExprStmt(js.ConsoleLog(js.String("positive"))),
	)
	fmt.Println(js.ToJSStmt(stmt))
	// Output: if ((x > 0)) { console.log("positive") }
}

func ExampleIfElse() {
	stmt := js.IfElse(
		js.Gt(js.Ident("x"), js.Int(0)),
		[]js.Stmt{js.Return(js.String("positive"))},
		[]js.Stmt{js.Return(js.String("non-positive"))},
	)
	fmt.Println(js.ToJSStmt(stmt))
	// Output: if ((x > 0)) { return "positive" } else { return "non-positive" }
}

func ExampleTernary() {
	expr := js.Ternary(
		js.Gt(js.Ident("age"), js.Int(18)),
		js.String("adult"),
		js.String("minor"),
	)
	fmt.Println(js.ExprHandler(expr))
	// Output: ((age > 18) ? "adult" : "minor")
}

func ExampleTemplate() {
	expr := js.Template("User: ", js.Ident("name"), " (", js.Ident("id"), ")")
	fmt.Println(js.ExprHandler(expr))
	// Output: `User: ${name} (${id})`
}

func ExampleAwait() {
	expr := js.Await(js.Method(js.Fetch(js.String("/api")), "json"))
	fmt.Println(js.ExprHandler(expr))
	// Output: await fetch("/api").json()
}

func ExampleAsyncArrowFunc() {
	fn := js.AsyncArrowFunc([]string{"url"},
		js.Await(js.Fetch(js.Ident("url"))),
	)
	fmt.Println(js.ExprHandler(fn))
	// Output: async url => await fetch(url)
}
