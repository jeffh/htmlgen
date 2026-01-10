// Package js provides type-safe JavaScript string generation for HTML event attributes.
//
// This package provides a builder API for generating JavaScript code strings
// suitable for use in HTML event handler attributes like onclick, onchange, etc.
//
// # Type Safety
//
// The package uses distinct types to prevent accidental mixing of Go strings
// with JavaScript expressions:
//   - [Expr]: JavaScript expressions that produce values
//   - [Stmt]: JavaScript statements that perform actions
//   - [Callable]: Expressions that can have properties accessed or methods called
//
// # Raw Escape Hatch
//
// The only way to inject arbitrary JavaScript is through the [Raw] function:
//
//	js.Raw("customFunction()")
//
// All other constructors produce properly escaped/encoded output.
//
// # Value Constructors
//
// Safe constructors for JavaScript literals:
//   - [String]: JSON-escaped string literal
//   - [Int], [Int64], [Float]: Number literals
//   - [Bool]: Boolean literal
//   - [Null], [Undefined]: Null and undefined literals
//   - [JSON]: Any Go value serialized as JSON
//   - [Array], [Object]: Array and object literals
//   - [Ident]: Variable/identifier reference
//   - [This]: The this keyword
//
// # Expression Builders
//
// Build complex expressions:
//   - [Prop]: Property access (obj.name)
//   - [Index]: Index access (obj[idx])
//   - [Call]: Function call
//   - [Method]: Method call
//   - [New]: New expression
//   - [OptionalProp], [OptionalCall]: Optional chaining (?.)
//   - Operators: [Add], [Sub], [Mul], [Div], [Eq], [Lt], [And], [Or], [Not], [Ternary], etc.
//
// # Statement Builders
//
// Build JavaScript statements:
//   - [Assign]: Assignment
//   - [Let], [Const], [Var]: Variable declarations
//   - [Incr], [Decr]: Increment/decrement
//   - [If], [IfElse]: Conditional
//   - [Return], [Throw]: Control flow
//   - [Stmts]: Combine multiple statements
//
// # Built-in Helpers
//
// Common JavaScript patterns:
//   - Console: [ConsoleLog], [ConsoleError], [ConsoleWarn]
//   - Document: [GetElementById], [QuerySelector]
//   - Event: [PreventDefault], [StopPropagation], [EventTarget], [EventValue]
//   - Navigation: [Navigate], [Reload], [HistoryBack]
//   - Functions: [ArrowFunc], [ArrowFuncStmts], [Func], [IIFE]
//   - Async: [Await], [AsyncArrowFunc], [PromiseThen]
//
// # Handler Attributes
//
// Create HTML event attributes that return [h.Attribute]:
//
//	js.OnClick(js.ExprStmt(js.ConsoleLog(js.String("clicked"))))
//	// => onclick="console.log(\"clicked\")"
//
// Event attribute functions include [OnClick], [OnInput], [OnChange],
// [OnSubmit], [OnKeyDown], [OnLoad], and many more.
//
// # Usage Example
//
//	import (
//	    "github.com/jeffh/htmlgen/h"
//	    "github.com/jeffh/htmlgen/js"
//	)
//
//	button := h.Button(
//	    js.OnClick(
//	        js.ExprStmt(js.PreventDefault()),
//	        js.Assign(
//	            js.Prop(js.GetElementById(js.String("output")), "textContent"),
//	            js.String("Clicked!"),
//	        ),
//	    ),
//	    h.Text("Click me"),
//	)
//
// This generates:
//
//	<button onclick="event.preventDefault(); document.getElementById(&quot;output&quot;).textContent = &quot;Clicked!&quot;">Click me</button>
//
// # Integration with h package
//
// The event handler functions (OnClick, OnInput, etc.) return [h.Attribute] values
// that can be passed directly to HTML element functions:
//
//	h.Input(
//	    h.Attrs("type", "text"),
//	    js.OnInput(
//	        js.Assign(
//	            js.Prop(js.EventTarget(), "value"),
//	            js.Method(js.EventValue(), "toUpperCase"),
//	        ),
//	    ),
//	)
package js
