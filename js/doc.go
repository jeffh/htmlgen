// Package js provides type-safe JavaScript string generation for HTML event attributes.
//
// This package provides a builder API for generating JavaScript code strings
// suitable for use in HTML event handler attributes like onclick, onchange, etc.
// It integrates seamlessly with the [github.com/jeffh/htmlgen/h] package.
//
// # Quick Start
//
// The most common use case is adding event handlers to HTML elements:
//
//	import (
//	    "github.com/jeffh/htmlgen/h"
//	    "github.com/jeffh/htmlgen/js"
//	)
//
//	// Simple click handler
//	button := h.Button(
//	    js.OnClick(js.ExprStmt(js.Alert(js.String("Hello!")))),
//	    h.Text("Say Hello"),
//	)
//	// Output: <button onclick="alert(&quot;Hello!&quot;)">Say Hello</button>
//
// # Type System
//
// The package uses three core interfaces to ensure type safety:
//
//   - [Expr] - JavaScript expressions that produce values (e.g., "1 + 2", "x.foo")
//   - [Stmt] - JavaScript statements that perform actions (e.g., "let x = 1", "x++")
//   - [Callable] - Expressions that can have properties accessed or methods called
//
// This type system prevents accidentally passing raw Go strings where JavaScript
// is expected. The only way to inject arbitrary JavaScript is through the explicit
// [Raw] escape hatch.
//
// # Creating Values
//
// Use these functions to create JavaScript literals safely:
//
//	js.String("hello")     // "hello" (JSON-escaped, prevents XSS)
//	js.Int(42)             // 42
//	js.Float(3.14)         // 3.14
//	js.Bool(true)          // true
//	js.Null()              // null
//	js.Undefined()         // undefined
//
// For complex values, use [JSON], [Array], or [Object]:
//
//	js.JSON(map[string]int{"a": 1})           // {"a":1}
//	js.Array(js.Int(1), js.Int(2), js.Int(3)) // [1, 2, 3]
//	js.Object(
//	    js.Pair("name", js.String("John")),
//	    js.Pair("age", js.Int(30)),
//	)                                         // {"name": "John", "age": 30}
//
// To reference JavaScript variables, use [Ident]:
//
//	js.Ident("myVariable")  // myVariable
//	js.Ident("window")      // window
//	js.This()               // this
//
// # Property and Method Access
//
// Access properties with [Prop] and call methods with [Method]:
//
//	js.Prop(js.Ident("document"), "body")
//	// document.body
//
//	js.Method(js.Ident("console"), "log", js.String("hello"))
//	// console.log("hello")
//
//	js.Prop(js.Prop(js.Ident("event"), "target"), "value")
//	// event.target.value
//
// For array/computed property access, use [Index]:
//
//	js.Index(js.Ident("arr"), js.Int(0))         // arr[0]
//	js.Index(js.Ident("obj"), js.String("key"))  // obj["key"]
//
// Optional chaining is supported with [OptionalProp] and [OptionalCall]:
//
//	js.OptionalProp(js.Ident("user"), "name")     // user?.name
//	js.OptionalCall(js.Ident("obj"), "method")    // obj?.method()
//
// # Operators
//
// All standard JavaScript operators are available:
//
//	// Arithmetic
//	js.Add(js.Int(1), js.Int(2))    // (1 + 2)
//	js.Sub(js.Int(5), js.Int(3))    // (5 - 3)
//	js.Mul(js.Int(4), js.Int(2))    // (4 * 2)
//	js.Div(js.Int(10), js.Int(2))   // (10 / 2)
//
//	// Comparison (strict by default)
//	js.Eq(js.Ident("x"), js.Int(5))       // (x === 5)
//	js.NotEq(js.Ident("x"), js.Null())    // (x !== null)
//	js.Lt(js.Ident("x"), js.Int(10))      // (x < 10)
//	js.Gt(js.Ident("x"), js.Int(0))       // (x > 0)
//
//	// Logical
//	js.And(js.Ident("a"), js.Ident("b"))  // (a && b)
//	js.Or(js.Ident("a"), js.Ident("b"))   // (a || b)
//	js.Not(js.Ident("x"))                 // !x
//
//	// Ternary
//	js.Ternary(js.Ident("cond"), js.String("yes"), js.String("no"))
//	// (cond ? "yes" : "no")
//
//	// Nullish coalescing
//	js.NullishCoalesce(js.Ident("x"), js.String("default"))
//	// (x ?? "default")
//
// # Statements
//
// Create JavaScript statements for use in handlers:
//
//	// Variable declarations
//	js.Let("x", js.Int(5))           // let x = 5
//	js.Const("PI", js.Float(3.14))   // const PI = 3.14
//
//	// Assignment
//	js.Assign(js.Ident("x"), js.Int(10))  // x = 10
//	js.AddAssign(js.Ident("x"), js.Int(1)) // x += 1
//
//	// Increment/decrement
//	js.Incr(js.Ident("count"))  // count++
//	js.Decr(js.Ident("count"))  // count--
//
//	// Conditionals
//	js.If(js.Eq(js.Ident("x"), js.Int(0)),
//	    js.Return(js.Null()),
//	)
//	// if (x === 0) { return null }
//
// To use an expression as a statement, wrap it with [ExprStmt]:
//
//	js.ExprStmt(js.ConsoleLog(js.String("hello")))
//	// console.log("hello")
//
// # Event Handlers
//
// The [Handler] function combines statements into a handler string:
//
//	handler := js.Handler(
//	    js.ExprStmt(js.PreventDefault()),
//	    js.Let("value", js.EventValue()),
//	    js.ExprStmt(js.ConsoleLog(js.Ident("value"))),
//	)
//	// "event.preventDefault(); let value = event.target.value; console.log(value)"
//
// For convenience, use the On* functions to create [h.Attribute] values directly:
//
//	js.OnClick(...)      // onclick="..."
//	js.OnInput(...)      // oninput="..."
//	js.OnChange(...)     // onchange="..."
//	js.OnSubmit(...)     // onsubmit="..."
//	js.OnKeyDown(...)    // onkeydown="..."
//	js.OnLoad(...)       // onload="..."
//
// For custom events, use [On]:
//
//	js.On("touchstart", js.ExprStmt(js.ConsoleLog(js.String("touched"))))
//	// ontouchstart="console.log(\"touched\")"
//
// # Built-in Helpers
//
// The package provides helpers for common JavaScript patterns:
//
// Console:
//
//	js.ConsoleLog(js.String("message"))   // console.log("message")
//	js.ConsoleError(js.String("error"))   // console.error("error")
//	js.ConsoleWarn(js.String("warning"))  // console.warn("warning")
//
// Document:
//
//	js.GetElementById(js.String("myId"))     // document.getElementById("myId")
//	js.QuerySelector(js.String(".myClass"))  // document.querySelector(".myClass")
//
// Event handling:
//
//	js.PreventDefault()      // event.preventDefault()
//	js.StopPropagation()     // event.stopPropagation()
//	js.EventTarget()         // event.target
//	js.EventValue()          // event.target.value
//	js.EventChecked()        // event.target.checked
//	js.EventKey()            // event.key
//
// Navigation:
//
//	js.Navigate(js.String("/home"))  // location.href = "/home"
//	js.Reload()                      // location.reload()
//	js.HistoryBack()                 // history.back()
//
// DOM manipulation:
//
//	js.ClassListAdd(js.Ident("el"), js.String("active"))
//	// el.classList.add("active")
//
//	js.ClassListToggle(js.Ident("el"), js.String("hidden"))
//	// el.classList.toggle("hidden")
//
//	js.SetStyle(js.Ident("el"), "backgroundColor", js.String("red"))
//	// el.style.backgroundColor = "red"
//
// # Arrow Functions
//
// Create arrow functions for callbacks:
//
//	// Expression body
//	js.ArrowFunc([]string{"x"}, js.Mul(js.Ident("x"), js.Int(2)))
//	// x => (x * 2)
//
//	// Statement body
//	js.ArrowFuncStmts([]string{"x"},
//	    js.Let("result", js.Mul(js.Ident("x"), js.Int(2))),
//	    js.Return(js.Ident("result")),
//	)
//	// x => { let result = (x * 2); return result }
//
//	// Async arrow functions
//	js.AsyncArrowFunc([]string{}, js.Await(js.Fetch(js.String("/api"))))
//	// async () => await fetch("/api")
//
// # Template Literals
//
// Create template literals with [Template]:
//
//	js.Template("Hello, ", js.Ident("name"), "!")
//	// `Hello, ${name}!`
//
// # Promises and Async
//
//	js.Await(js.Fetch(js.String("/api/data")))
//	// await fetch("/api/data")
//
//	js.PromiseThen(
//	    js.Fetch(js.String("/api")),
//	    js.ArrowFunc([]string{"r"}, js.Method(js.Ident("r"), "json")),
//	)
//	// fetch("/api").then(r => r.json())
//
// # Raw JavaScript Escape Hatch
//
// When you need to inject arbitrary JavaScript that isn't covered by the API,
// use [Raw]. This is the ONLY way to inject raw JavaScript and must be used
// explicitly:
//
//	js.Raw("myCustomFunction()")
//	js.Raw("window.gtag('event', 'click')")
//
// Use Raw sparingly, as it bypasses type safety. The API covers most common
// use cases, so prefer using the type-safe builders when possible.
//
// # Complete Examples
//
// Form submission with validation:
//
//	h.Form(
//	    js.OnSubmit(
//	        js.ExprStmt(js.PreventDefault()),
//	        js.Let("email", js.Prop(js.GetElementById(js.String("email")), "value")),
//	        js.If(js.Eq(js.Ident("email"), js.String("")),
//	            js.ExprStmt(js.Alert(js.String("Email is required"))),
//	            js.ReturnVoid(),
//	        ),
//	        js.ExprStmt(js.Method(js.This(), "submit")),
//	    ),
//	    // ... form fields
//	)
//
// Toggle visibility:
//
//	h.Button(
//	    js.OnClick(
//	        js.ExprStmt(js.ClassListToggle(
//	            js.GetElementById(js.String("panel")),
//	            js.String("hidden"),
//	        )),
//	    ),
//	    h.Text("Toggle Panel"),
//	)
//
// Live character count:
//
//	h.Textarea(
//	    h.Attrs("id", "message", "maxlength", "200"),
//	    js.OnInput(
//	        js.Assign(
//	            js.Prop(js.GetElementById(js.String("charCount")), "textContent"),
//	            js.Template(js.Prop(js.EventValue(), "length"), " / 200"),
//	        ),
//	    ),
//	)
//
// Keyboard shortcuts:
//
//	h.Body(
//	    js.OnKeyDown(
//	        js.If(js.And(js.EventCtrlKey(), js.Eq(js.EventKey(), js.String("s"))),
//	            js.ExprStmt(js.PreventDefault()),
//	            js.ExprStmt(js.Raw("saveDocument()")),
//	        ),
//	    ),
//	)
//
// Fetch with error handling:
//
//	h.Button(
//	    js.OnClick(
//	        js.ExprStmt(
//	            js.PromiseCatch(
//	                js.PromiseThen(
//	                    js.Fetch(js.String("/api/data")),
//	                    js.ArrowFunc([]string{"r"}, js.Method(js.Ident("r"), "json")),
//	                ),
//	                js.ArrowFunc([]string{"err"}, js.ConsoleError(js.Ident("err"))),
//	            ),
//	        ),
//	    ),
//	    h.Text("Load Data"),
//	)
package js
