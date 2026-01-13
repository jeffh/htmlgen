# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run all tests
go test ./...

# Run a specific test
go test -run TestName ./...

# Build/check compilation
go build ./...
```

## Architecture

This is a Go library (`github.com/jeffh/htmlgen`) for programmatic HTML generation with three packages:

### Package `h` - Core HTML Generation

**Writer API** (`h/writer.go`): Low-level streaming HTML writer wrapping `io.Writer`. Tracks open tags and provides:
- `Doctype()`, `OpenTag()`, `CloseTag()`, `SelfClosingTag()` for tag manipulation
- `Text()` (escaped) and `Raw()` (unescaped) for content
- `CloseOneTag()` closes most recent tag, `Close()` closes all remaining tags
- Attribute values are automatically HTML-escaped via `html/template`

**Builder API** (`h/builder.go`, `h/tags.go`): Declarative tree-building API where nodes implement `Builder` interface (`Build(w *Writer) error`). Use `Render(w, builder)` to write a builder tree to an io.Writer. All standard HTML5 elements have corresponding functions (e.g., `Div()`, `Span()`, `A()`) that take `Attributes` and child `Builder` elements.

**Attributes** (`h/attrs.go`): `Attributes` is a `[]Attribute` slice with `Get()`, `Set()`, `SetDefault()`, `Delete()` methods. Create via `Attrs("key", "value", ...)` or `AttrsMap(map[string]string{...})`.

**Helpers** (`h/helpers.go`): Control flow and iteration utilities for builder composition:
- `If(cond, ifTrue, ifElse)` - returns `ifTrue` if cond is true, else `ifElse`
- `When(cond, ifTrue)` - returns `ifTrue` if cond is true, else nil (skipped during render)
- `First(b...)` - returns first non-nil builder from arguments
- `ForEach[X](seq, fn)` - lazily maps `iter.Seq[X]` to builders during render
- `ForEach2[X,Y](seq, fn)` - lazily maps `iter.Seq2[X,Y]` (e.g., `slices.All()`, `maps.All()`) to builders

### Package `d` - Datastar Attribute Helpers

Provides helpers for building [Datastar](https://data-star.dev/) reactive attributes:

- **Signals**: `Signal()`, `Signals()`, `Bind()` - define reactive state
- **Events**: `OnClick()`, `OnSubmit()`, `OnInput()`, `OnChange()`, `OnLoad()`, `On()` - event handlers
- **Actions**: `Get()`, `Post()`, `Put()`, `Delete()` - HTTP request helpers
- **Modifiers**: `PreventDefault()`, `Debounce()`, `Throttle()`, `Delay()`, `Once()`, `ViewTransition()` - event modifiers
- **Values**: `Raw()`, `JsonValue()`, `Str()` - value builders for expressions

The `d` package uses a builder pattern with `AttrMutator` and `AttrValueAppender` interfaces to compose complex Datastar attributes.

### Package `js` - Type-Safe JavaScript Generation

Provides a type-safe builder API for generating JavaScript code strings for HTML event handler attributes (`onclick`, `onsubmit`, etc.). Integrates with the `h` package.

**Core Types** (`js/expr.go`, `js/stmt.go`):
- `Expr` - JavaScript expressions that produce values (e.g., `"1 + 2"`, `"x.foo"`)
- `Stmt` - JavaScript statements that perform actions (e.g., `"let x = 1"`, `"x++"`)
- `Callable` - Expressions that support property access and method calls

**Values** (`js/values.go`): Create literals with `String()`, `Int()`, `Float()`, `Bool()`, `Null()`, `Undefined()`, `JSON()`, `Array()`, `Object()`. Reference variables with `Ident()` or `This()`.

**Property/Method Access** (`js/access.go`): Use `Prop()` for property access, `Method()` for method calls, `Index()` for array/computed access, `OptionalProp()`/`OptionalCall()` for optional chaining.

**Operators** (`js/operators.go`): Arithmetic (`Add`, `Sub`, `Mul`, `Div`), comparison (`Eq`, `NotEq`, `Lt`, `Gt`), logical (`And`, `Or`, `Not`), ternary (`Ternary`), nullish coalescing (`NullishCoalesce`).

**Statements** (`js/stmt.go`): Variable declarations (`Let`, `Const`), assignment (`Assign`, `AddAssign`), increment/decrement (`Incr`, `Decr`), conditionals (`If`, `IfElse`), returns (`Return`, `ReturnVoid`). Wrap expressions as statements with `ExprStmt()`.

**Event Handlers** (`js/handler.go`): `Handler()` combines statements into a handler string. Convenience functions `OnClick()`, `OnInput()`, `OnSubmit()`, `OnChange()`, `OnKeyDown()`, `OnLoad()`, `On()` create `h.Attribute` values directly.

**Built-ins** (`js/builtins.go`): Helpers for console (`ConsoleLog`, `ConsoleError`), document (`GetElementById`, `QuerySelector`), events (`PreventDefault`, `StopPropagation`, `EventValue`, `EventTarget`), navigation (`Navigate`, `Reload`, `HistoryBack`), DOM manipulation (`ClassListAdd`, `ClassListToggle`, `SetStyle`).

**Functions** (`js/func.go`): Arrow functions with `ArrowFunc()` (expression body) and `ArrowFuncStmts()` (statement body). Async variants with `AsyncArrowFunc()`, `AsyncArrowFuncStmts()`. Await with `Await()`.

**Raw Escape Hatch**: `Raw()` is the only way to inject arbitrary JavaScript - use sparingly as it bypasses type safety.
