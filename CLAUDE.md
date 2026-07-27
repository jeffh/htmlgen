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

**Streaming API** (`h/writer.go`, `h/tags.go`, `h/tags_typed.go`, `h/render.go`): `Render` creates a per-render `*h.B` bound to an `io.Writer`. All standard HTML5 elements are methods on `B` (for example, `b.Div`, `b.Span`, and `b.A`). Container elements accept attributes plus a trailing `func(*h.B)` body, which runs immediately and streams its children. Void elements write self-closing tags. `Text`/`Textf` escape content, while `Raw`/`Rawf` write caller-sanitized content unchanged. `El`, `VoidEl`, `ElE`, and `VoidElE` support custom tags; they panic unless the tag name matches `[A-Za-z][A-Za-z0-9_.:-]*` (attribute names are validated the same way in `Attr`/`Attrs`/`AttrsMap`/`AttrIf`/`Set`/`SetDefault`, since names are written unescaped; `Attribute` struct literals bypass validation and are trusted).

**Typed fast paths** (`h/tags_typed.go`): every element method `Xxx` has a sibling `XxxE(attrs Attributes, body Body)` (void elements: `XxxE(attrs Attributes)`), plus `ElE`, `VoidElE`, and `HtmlE`. They take concrete parameters instead of `...any`, so arguments are never boxed and a capturing body closure stays stack-allocated. `nil` attrs and a `nil` body are valid. Prefer them in hot paths; the variadic forms remain the ergonomic default.

`B` buffers output internally and writes to the `io.Writer` in ~4 KiB chunks. `Render`/`RenderIndent` flush before returning; call `Flush()` explicitly when bytes must reach the client mid-render (server-sent events, long streaming responses). A `Raw` value at or above the flush threshold bypasses the buffer.

`B` keeps the first write error as a sticky error; later output calls become no-ops, and `Render` returns that error. Because output is buffered, `Err()` reflects only errors observed as of the last flush. `RenderIndent` enables pretty-printing. `RenderString` and `RenderBytes` are in-memory convenience entry points. Each render owns its `B`, so callers use native Go `if` and `for` statements safely without ambient package state.

**Attributes** (`h/attrs.go`): `Attributes` is a `[]Attribute` slice with `Get()`, `Set()`, `SetDefault()`, `Delete()` methods. Create via `Attrs("key", "value", ...)` or `AttrsMap(map[string]string{...})`.

Element arguments accept `Attributes`, `Attribute`, any `AttrBuilder`, a body closure, or `nil`. Later attributes override earlier values without changing their position. Companion-package attribute builders from `ds`, `hx`, and `js` can be passed directly to element methods.

### Package `ds` - Datastar Attribute Helpers

Provides helpers for building [Datastar](https://data-star.dev/) reactive attributes:

- **Signals**: `Signal()`, `Signals()`, `Bind()` - define reactive state
- **Typed signals**: `Sig("name")` - a signal name with methods `Ref()`, `Value()`, `Not()`, `Set()`, `SetExpr()`, `Toggle()`, `Clear()`, `Eq()`, `NotEq()`, `Sub()` (derived `name_suffix` signals)
- **Events**: `OnClick()`, `OnSubmit()`, `OnInput()`, `OnChange()`, `On()` - event handlers; `Init()` - run on element load (data-init)
- **Actions**: `Get()`, `Post()`, `Put()`, `Delete()` - HTTP request helpers
- **Modifiers**: `PreventDefault()`, `Debounce()`, `Throttle()`, `Delay()`, `Once()`, `ViewTransition()` - event modifiers
- **Values**: `Raw()`, `JsonValue()`, `Str()` - value builders for expressions
- **Composition**: `Do(stmts...)` bridges typed `js.Stmt`s into a Value (statement positions only: `data-on:*`, `data-init`, `data-effect`); Value methods `Not()`, `And()`, `Or()`, `Ternary()` combine expressions; `Confirm(msg, then...)` guards actions behind a `confirm()` dialog
- **Expression-valued maps**: `ClassesExpr()`, `StylesExpr()`, `AttrsExpr()` emit `data-class`/`data-style`/`data-attr` object literals with expression values and sorted keys. Prefer these over `Classes()`/`Styles()`/`Attrs()`, which JSON-encode values into always-truthy string literals
- **Scope identifiers**: `Evt`, `El`, `EvtTarget`, `EvtValue`, `EvtKey` - Datastar expressions expose `evt`/`el`, not the `event` of legacy inline handlers (`js.Event` and the deprecated `Event`/`EventTarget`/`EventValue` re-exports)

The `ds` package composes attributes with fluent builders: `OnClick()`/`On()`/`Bind()`/`Signals()` and friends return builder structs whose modifier methods (`.Outside()`, `.PreventDefault()`, `.Debounce()`, ...) append to the attribute name and whose `Attribute()` method produces the final `h.Attribute`.

### Package `js` - Type-Safe JavaScript Generation

Provides a type-safe builder API for generating JavaScript code strings for HTML event handler attributes (`onclick`, `onsubmit`, etc.). Integrates with the `h` package.

**Core Types** (`js/expr.go`, `js/stmt.go`):
- `Expr` - JavaScript expressions that produce values (e.g., `"1 + 2"`, `"x.foo"`); property access, method calls and operators are methods on `Expr`
- `Stmt` - JavaScript statements that perform actions (e.g., `"let x = 1"`, `"x++"`)

**Values** (`js/values.go`): Create literals with `String()`, `Int()`, `Float()`, `Bool()`, `Null()`, `Undefined()`, `JSON()`, `Array()`, `Object()`. Reference variables with `Ident()` or `This()`. `Regex(pattern, flags)` emits a `/pattern/flags` literal verbatim (no escaping).

**Property/Method Access** (`js/access.go`): Use `Prop()` for property access, `Method()` for method calls, `Index()` for array/computed access, `OptionalProp()`/`OptionalCall()` for optional chaining.

**Operators** (`js/operators.go`): Arithmetic (`Add`, `Sub`, `Mul`, `Div`), comparison (`Eq`, `NotEq`, `Lt`, `Gt`), logical (`And`, `Or`, `Not`), ternary (`Ternary`), nullish coalescing (`NullishCoalesce`).

**Statements** (`js/stmt.go`): Variable declarations (`Let`, `Const`), assignment (`Assign`, `AddAssign`), increment/decrement (`Incr`, `Decr`), conditionals (`If`, `IfElse`), returns (`Return`, `ReturnVoid`), error handling (`TryCatch(body, errName, catchBody)`, and `Try(body...)` for a best-effort `try { ... } catch {}`). Wrap expressions as statements with `ExprStmt()`.

**Event Handlers** (`js/handler.go`): `Handler()` combines statements into a handler string. Convenience functions `OnClick()`, `OnInput()`, `OnSubmit()`, `OnChange()`, `OnKeyDown()`, `OnLoad()`, `On()` create `h.Attribute` values directly.

**Built-ins** (`js/builtins.go`): Helpers for console (`ConsoleLog`, `ConsoleError`), document (`GetElementById`, `QuerySelector`), events (`PreventDefault`, `StopPropagation`, `EventValue`, `EventTarget`), navigation (`Navigate`, `Reload`, `HistoryBack`), DOM manipulation (`ClassListAdd`, `ClassListToggle`, `SetStyle`).

**Functions** (`js/func.go`): Arrow functions with `ArrowFunc()` (expression body) and `ArrowFuncStmts()` (statement body). Async variants with `AsyncArrowFunc()`, `AsyncArrowFuncStmts()`. Await with `Await()`.

**Raw Escape Hatch**: `Raw()` is the only way to inject arbitrary JavaScript - use sparingly as it bypasses type safety.
