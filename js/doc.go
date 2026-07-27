// Package js provides type-safe JavaScript string generation for HTML event attributes.
//
// This package provides a fluent builder API for generating JavaScript code
// strings suitable for use in HTML event handler attributes like onclick,
// onchange, etc. It integrates seamlessly with the [github.com/jeffh/htmlgen/h] package.
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
//	button := h.RenderString(func(b *h.B) {
//	    b.Button(
//	        h.AttrsOf(js.OnClick(js.Alert(js.String("Hello!")).Stmt())),
//	        func(b *h.B) { b.Text("Say Hello") },
//	    )
//	})
//	// Output: <button onclick="alert(&quot;Hello!&quot;)">Say Hello</button>
//
// # Type System
//
// The package exposes two core types:
//
//   - [Expr] - JavaScript expressions that produce values; operators and
//     property/method access are exposed as methods on Expr (e.g. e.Add(other),
//     e.Eq(other), e.Prop("x"), e.Method("foo", arg)).
//   - [Stmt] - JavaScript statements (e.g. let, return, assignment).
//
// Both are concrete struct values, not interfaces. Construct them with the
// package constructors ([String], [Int], [Ident], [Raw], etc.) and then chain
// methods to build larger expressions or convert to statements.
//
// # Creating Values
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
// Property access and method calls are methods on Expr:
//
//	js.Ident("document").Prop("body")                       // document.body
//	js.Ident("console").Method("log", js.String("hello"))   // console.log("hello")
//	js.Ident("event").Prop("target").Prop("value")          // event.target.value
//
// For array/computed property access, use [Expr.Index]:
//
//	js.Ident("arr").Index(js.Int(0))         // arr[0]
//	js.Ident("obj").Index(js.String("key"))  // obj["key"]
//
// Optional chaining is supported with [Expr.OptionalProp] and [Expr.OptionalCall]:
//
//	js.Ident("user").OptionalProp("name")           // user?.name
//	js.Ident("obj").OptionalCall("method")          // obj?.method()
//
// # Operators
//
// Operators are methods on Expr:
//
//	js.Int(1).Add(js.Int(2))                      // (1 + 2)
//	js.Int(5).Sub(js.Int(3))                      // (5 - 3)
//	js.Ident("x").Eq(js.Int(5))                   // (x === 5)
//	js.Ident("x").Lt(js.Int(10))                  // (x < 10)
//	js.Ident("a").And(js.Ident("b"))              // (a && b)
//	js.Ident("x").Not()                           // !x
//	js.Ident("cond").Ternary(js.String("yes"), js.String("no"))
//	// (cond ? "yes" : "no")
//
//	js.Ident("x").NullishCoalesce(js.String("default"))
//	// (x ?? "default")
//
// # Statements
//
//	js.Let("x", js.Int(5))                         // let x = 5
//	js.Const("PI", js.Float(3.14))                 // const PI = 3.14
//	js.Ident("x").Assign(js.Int(10))               // x = 10
//	js.Ident("x").AddAssign(js.Int(1))             // x += 1
//	js.Ident("count").Incr()                       // count++
//	js.If(js.Ident("x").Eq(js.Int(0)),
//	    js.Return(js.Null()),
//	)                                              // if (x === 0) { return null }
//
// To use an expression as a statement, call [Expr.Stmt]:
//
//	js.ConsoleLog(js.String("hello")).Stmt()
//	// console.log("hello")
//
// # Event Handlers
//
// The [Handler] function combines statements into a handler string:
//
//	handler := js.Handler(
//	    js.PreventDefault().Stmt(),
//	    js.Let("value", js.EventValue()),
//	    js.ConsoleLog(js.Ident("value")).Stmt(),
//	)
//	// "event.preventDefault(); let value = event.target.value; console.log(value)"
//
// For convenience, [OnClick], [OnInput], etc. create [h.Attribute] values directly.
//
// # Built-in Helpers
//
//	js.ConsoleLog(js.String("message"))            // console.log("message")
//	js.GetElementById(js.String("myId"))           // document.getElementById("myId")
//	js.QuerySelector(js.String(".myClass"))        // document.querySelector(".myClass")
//	js.PreventDefault()                            // event.preventDefault()
//	js.EventValue()                                // event.target.value
//	js.Navigate(js.String("/home"))                // location.href = "/home"
//
// DOM helpers are methods on element expressions:
//
//	js.Ident("el").ClassListAdd(js.String("active"))
//	// el.classList.add("active")
//
//	js.Ident("el").SetStyle("backgroundColor", js.String("red"))
//	// el.style.backgroundColor = "red"
//
// # Arrow Functions and Async
//
//	js.ArrowFunc([]string{"x"}, js.Ident("x").Mul(js.Int(2)))
//	// x => (x * 2)
//
//	js.AsyncArrowFunc([]string{}, js.Fetch(js.String("/api")).Await())
//	// async () => await fetch("/api")
//
// # Promises
//
//	js.Fetch(js.String("/api")).
//	    Then(js.ArrowFunc([]string{"r"}, js.Ident("r").Method("json"))).
//	    Catch(js.ArrowFunc([]string{"err"}, js.ConsoleError(js.Ident("err"))))
//	// fetch("/api").then(r => r.json()).catch(err => console.error(err))
//
// # Template Literals
//
//	js.Template("Hello, ", js.Ident("name"), "!")
//	// `Hello, ${name}!`
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
// Use Raw sparingly, as it bypasses type safety.
package js
