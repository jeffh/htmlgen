# CLAUDE.md

This file provides package-level detail for Claude Code (claude.ai/code).
Commands, code style, and the streaming HTML API live in [AGENTS.md](AGENTS.md)
and [API_DESIGN.md](API_DESIGN.md). Do not invent APIs or refactor the HTML tag
surface.

```bash
go test ./...
go build ./...
go vet ./...
```

## Architecture

`github.com/jeffh/htmlgen` is a streaming HTML library with five packages:
`h`, `ds`, `hx`, `js`, and `ds/rkt`.

### Package `h` - Core HTML Generation

See [API_DESIGN.md](API_DESIGN.md). `Render` creates a per-render `*h.B`.
Containers are `Xxx(attrs Attributes, body Body)`; void elements are
`Xxx(attrs Attributes)`. `nil` attrs and a `nil` body are valid. Collect
`AttrBuilder`s with `AttrsOf` / `Attributes.With`.

### Package `ds` - Datastar Attribute Helpers

Provides helpers for building [Datastar](https://data-star.dev/) reactive attributes:

- **Signals**: `Signal()`, `Signals()`, `Bind()` - define reactive state
- **Typed signals**: `Sig("name")` - a signal name with methods `Ref()`, `Value()`, `Not()`, `Set()`, `SetExpr()`, `Toggle()`, `Clear()`, `Eq()`, `NotEq()`, `Sub()` (derived `name_suffix` signals)
- **Events**: `OnClick()`, `OnSubmit()`, `OnInput()`, `OnChange()`, `On()` - event handlers; `Init()` - run on element load (data-init)
- **Actions**: `Get()`, `Post()`, `Put()`, `Delete()` - HTTP request helpers
- **Modifiers**: `PreventDefault()`, `Debounce()`, `Throttle()`, `Delay()`, `Once()`, `ViewTransition()` - chain on event builders
- **Values**: `Raw()`, `JsonValue()`, `Str()` - value builders for expressions
- **Composition**: `Do(stmts...)` bridges typed `js.Stmt`s into a Value (statement positions only: `data-on:*`, `data-init`, `data-effect`); Value methods `Not()`, `And()`, `Or()`, `Ternary()` combine expressions; `Confirm(msg, then...)` guards actions behind a `confirm()` dialog
- **Expression-valued maps**: `ClassesExpr()`, `StylesExpr()`, `AttrsExpr()` emit `data-class`/`data-style`/`data-attr` object literals with expression values and sorted keys. Prefer these over `Classes()`/`Styles()`/`Attrs()`, which JSON-encode values into always-truthy string literals
- **Scope identifiers**: `Evt`, `El`, `EvtTarget`, `EvtValue`, `EvtKey` - Datastar expressions expose `evt`/`el`, not the `event` of legacy inline handlers

The `ds` package composes attributes with fluent builders: `OnClick()`/`On()`/`Bind()`/`Signals()` and friends return builder structs whose modifier methods (`.Outside()`, `.PreventDefault()`, `.Debounce()`, ...) append to the attribute name and whose `Attribute()` method produces the final `h.Attribute`.

As of Datastar v1.0.4, Rocket moved out of Datastar Pro into its own free
bundle; it is supported by the `ds/rkt` subpackage below, not by `ds/pro.go`.

### Package `ds/rkt` - Datastar Rocket Helpers

`github.com/jeffh/htmlgen/ds/rkt` provides Go helpers for [Rocket](https://data-star.dev/), Datastar's beta web-component API (`rocket('my-tag', {props, setup, render})`). Rocket layers custom elements — with props, private signals, and local actions wired into Datastar — on top of the core library; it ships in its own bundle, `datastar-rocket.js`, loaded in place of `datastar.js`. This package does not generate the `rocket(...)` JavaScript definition itself; it helps build the template markup, expressions, and host elements that use a Rocket component. See `docs/internal/data-star.md` for the full template conventions.

- **Private signals**: `Sig("name")` - like `ds.Sig` but for the component-private `$$name` signal, with the same methods: `Name()`, `Ref()`, `Value()`, `Not()`, `Set()`, `SetExpr()`, `Toggle()`, `Clear()`, `Eq()`, `NotEq()`, `Sub()`; `SignalRef(name)` for a bare `$$name` expression
- **Local actions**: `Action(name, args...)` - emits `@name(...)`, a call to an action registered in the component's `setup({action})`; `Dispatch(name, args...)` - emits `@dispatchRocket("name", ...)`
- **Template control flow**: `For(item, index, source)`, `ForRaw(expr)` - `data-for` loops over `$$`-scoped sources; `If(cond)`, `ElseIf(cond)`, `Else()` - `data-if`/`data-else-if`/`data-else` on sibling `<template>` elements; `Ref(name).Case(...)` - `data-ref:name` with `__case` casing; `Root(attrBuilder)` - appends `__root` to escape the component's private scope on `data-signals`/`data-bind`/`data-computed`/`data-indicator`
- **Bundle**: `Version`, `BundleURL`, `AliasedBundleURL` (CSP-safe build), `ScriptAttrs()`, `Script(b)`, `AliasedScript(b)` - load `datastar-rocket.js` (a superset of `datastar.js`) instead of the base bundle
- **Host element props**: `Prop(name, value)`, `JSONProp()`, `DateProp()`, `BinProp()`, `BoolProp()`, `NumberProp()`, `Props(map)` - kebab-case attribute names with codec encodings (bool, number, string, ISO date, JSON, base64 binary), sorted for `Props`
- **Components**: `Component(b, tag, attrs, body)` - writes a custom element, validating that `tag` is lowercase and contains a hyphen

### Package `hx` - HTMX Attribute Helpers

Fluent HTMX attribute helpers (`Get()`, `Target()`, `Swap()`, `Trigger()`, …).
See `hx/doc.go`.

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
