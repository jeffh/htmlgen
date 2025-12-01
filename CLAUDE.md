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

This is a Go library (`github.com/jeffh/htmlgen`) for programmatic HTML generation with two packages:

### Package `h` - Core HTML Generation

**Writer API** (`h/api.go`): Low-level streaming HTML writer wrapping `io.Writer`. Tracks open tags and provides:
- `Doctype()`, `OpenTag()`, `CloseTag()`, `SelfClosingTag()` for tag manipulation
- `Text()` (escaped) and `Raw()` (unescaped) for content
- `End()` closes most recent tag, `Close()` closes all remaining tags
- Attribute values are automatically HTML-escaped via `html/template`

**Builder API** (`h/builder.go`, `h/tags.go`): Declarative tree-building API where nodes implement `Builder` interface (`Build(w *Writer) error`). Use `Render(w, builder)` to write a builder tree to an io.Writer. All standard HTML5 elements have corresponding functions (e.g., `Div()`, `Span()`, `A()`) that take `Attributes` and child `Builder` elements.

**Attributes** (`h/attrs.go`): `Attributes` is a `[]Attribute` slice with `Get()`, `Set()`, `SetDefault()`, `Delete()` methods. Create via `Attrs("key", "value", ...)` or `AttrsMap(map[string]string{...})`.

**Platform Detection** (`h/dsl.go`, `h/dsl_js.go`): Build-tag controlled constants `Server`/`Client` indicate runtime environment (server-side Go vs WebAssembly).

### Package `d` - Datastar Attribute Helpers

Provides helpers for building [Datastar](https://data-star.dev/) reactive attributes:

- **Signals**: `Signal()`, `Signals()`, `Bind()` - define reactive state
- **Events**: `OnClick()`, `OnSubmit()`, `OnInput()`, `OnChange()`, `OnLoad()`, `On()` - event handlers
- **Actions**: `Get()`, `Post()`, `Put()`, `Delete()` - HTTP request helpers
- **Modifiers**: `PreventDefault()`, `Debounce()`, `Throttle()`, `Delay()`, `Once()`, `ViewTransition()` - event modifiers
- **Values**: `Raw()`, `JsonValue()`, `Str()` - value builders for expressions

The `d` package uses a builder pattern with `AttrMutator` and `AttrValueAppender` interfaces to compose complex Datastar attributes.
