# htmlgen

A Go library for programmatic HTML generation.

Experimental, API is subject to change.

## Installation

```bash
go get github.com/jeffh/htmlgen
```

## Overview

htmlgen provides three packages:

- **`h`** - Core HTML generation with both streaming and declarative APIs
- **`d`** - Datastar attribute helpers for building reactive web applications
- **`js`** - Type-safe JavaScript generation for event handler attributes

## Package `h` - HTML Generation

### Declarative Builder API

Build HTML trees using Go functions that mirror HTML elements:

```go
import "github.com/jeffh/htmlgen/h"

page := h.Html(
    h.Head(
        h.Title(h.Text("My Page")),
    ),
    h.Body(
        h.Div(h.Attrs("class", "container"),
            h.H1(h.Text("Hello, World!")),
            h.P(h.Text("Welcome to htmlgen.")),
            h.A(h.Attrs("href", "/about"), h.Text("About")),
        ),
    ),
)

// Render to any io.Writer
h.Render(os.Stdout, page)

// Or render with pretty-printed indentation (using two spaces)
h.RenderIndent(os.Stdout, "  ", page)
```

### Attributes

Create attributes using `Attrs()` with key-value pairs or `AttrsMap()` with a map:

```go
// Key-value pairs
attrs := h.Attrs("class", "btn", "id", "submit-btn", "disabled", "")

// From a map (keys are sorted for deterministic output)
attrs := h.AttrsMap(map[string]string{
    "class": "btn",
    "id":    "submit-btn",
})

// Modify attributes
attrs.Set("data-action", "submit")
attrs.SetDefault("type", "button")  // Only sets if not present
attrs.Delete("disabled")
value, ok := attrs.Get("class")
```

### Available Elements

All standard HTML5 elements are available as functions:

- **Document**: `Html`, `Head`, `Title`, `Meta`, `Link`, `Style`, `Script`, `Body`
- **Sections**: `Header`, `Footer`, `Main`, `Nav`, `Section`, `Article`, `Aside`
- **Headings**: `H1`, `H2`, `H3`, `H4`, `H5`, `H6`
- **Text**: `P`, `Span`, `Div`, `Pre`, `Code`, `Em`, `Strong`, `A`
- **Lists**: `Ul`, `Ol`, `Li`, `Dl`, `Dt`, `Dd`
- **Tables**: `Table`, `Thead`, `Tbody`, `Tfoot`, `Tr`, `Th`, `Td`
- **Forms**: `Form`, `Input`, `Button`, `Label`, `Select`, `Option`, `Textarea`, `Fieldset`
- **Media**: `Img`, `Video`, `Audio`, `Picture`, `Source`, `Canvas`, `Svg`
- **Helpers**: `Fragment`, `Text`, `Raw`, `CustomElement`

### Streaming Writer API

For lower-level control, use the Writer API directly:

```go
w := h.NewWriter(os.Stdout)
w.Doctype()
w.OpenTag("html", h.Attrs("lang", "en"))
w.OpenTag("body", nil)
w.Text("Hello, World!")
w.Close()  // Closes all open tags
```

### Pre-compiled Templates

For frequently rendered content, use `Compile` to pre-render HTML to bytes for faster subsequent renders:

```go
// Compile once at startup
header, err := h.Compile(h.Header(
    h.Nav(
        h.A(h.Attrs("href", "/"), h.Text("Home")),
        h.A(h.Attrs("href", "/about"), h.Text("About")),
    ),
))
if err != nil {
    // handle error
}

// Fast renders afterward - just writes pre-computed bytes
h.Render(w, header)

// Or use MustCompile to panic on error (for initialization code)
header := h.MustCompile(h.Header(
    h.Nav(
        h.A(h.Attrs("href", "/"), h.Text("Home")),
        h.A(h.Attrs("href", "/about"), h.Text("About")),
    ),
))
```

For templates with dynamic content, use `CompileParams` with parameter placeholders:

```go
// Define parameters
title := h.NewParam("title")
content := h.NewParam("content")

// Compile template with parameter slots
tmpl, err := h.CompileParams(h.Html(
    h.Head(h.Title(title)),
    h.Body(
        h.H1(title),
        h.Main(content),
    ),
))
if err != nil {
    // handle error
}

// Render with values
tmpl.Render(w,
    title.Value(h.Text("Welcome")),
    content.Value(h.P(h.Text("Hello, World!"))),
)

// Or create a reusable Builder
page := tmpl.With(
    title.Value(h.Text("Welcome")),
    content.Value(h.P(h.Text("Hello, World!"))),
)
h.Render(w, page)

// Or use MustCompileParams to panic on error (for initialization code)
tmpl := h.MustCompileParams(h.Html(
    h.Head(h.Title(title)),
    h.Body(
        h.H1(title),
        h.Main(content),
    ),
))
```

Compiled templates are ~8.0x faster than `html/template` for parameterized content.

## Package `d` - Datastar Integration

Build reactive attributes for [Datastar](https://data-star.dev/) applications:

### Signals

```go
import "github.com/jeffh/htmlgen/d"

// Define reactive signals
d.Signal("count", 0)           // data-signals:count="0"
d.Signal("name", "Alice")      // data-signals:name="\"Alice\""
d.Signals(map[string]any{      // data-signals="{...}"
    "x": 1,
    "y": 2,
})

// Two-way binding
d.Bind("username")             // data-bind="username"
```

### Event Handlers

```go
// Click events
d.OnClick(d.SetSignal("count", d.Raw("$count + 1")))

// Form events
d.OnSubmit(d.PreventDefault(), d.Post("/api/submit"))

// Other events
d.OnInput(d.Debounce(300*time.Millisecond), d.SetSignal("search", d.Raw("evt.target.value")))
d.OnChange(d.Get("/api/update"))
d.OnLoad(d.Get("/api/init"))
d.On("keydown", d.Raw("handleKey(evt)"))

// Intersection and interval observers
d.OnIntersect(d.Once(), d.Raw("$seen = true"))
d.OnInterval(d.Duration(1*time.Second), d.Raw("$tick++"))
```

### HTTP Actions

```go
d.Get("/api/data")
d.Post("/api/submit")
d.Put("/api/update")
d.Delete("/api/remove")

// With options
d.Post("/api/submit",
    d.ContentType("application/json"),
    d.Headers(map[string]string{"X-Custom": "value"}),
)
```

### Reactive Display

```go
d.Show(d.Raw("$isVisible"))                    // data-show="$isVisible"
d.Text(d.Raw("$message"))                      // data-text="$message"
d.Class("active", d.Raw("$isActive"))          // data-class:active="$isActive"
d.Style("color", d.Raw("$textColor"))          // data-style:color="$textColor"
d.Attribute("disabled", d.Raw("$isDisabled"))  // data-attr:disabled="$isDisabled"

// Multiple classes/styles/attrs at once
d.Classes(map[string]string{"hidden": "$foo", "bold": "$bar"})
d.Styles(map[string]string{"color": "$red ? 'red' : 'blue'"})
```

### Event Modifiers

```go
d.PreventDefault()
d.Debounce(300 * time.Millisecond)
d.Throttle(100 * time.Millisecond)
d.Delay(500 * time.Millisecond)
d.Once()
d.ViewTransition()
```

### Complete Example

```go
package main

import (
    "os"

    "github.com/jeffh/htmlgen/h"
    "github.com/jeffh/htmlgen/d"
)

func main() {
    page := h.Html(
        h.Head(
            h.Title(h.Text("Counter")),
            h.Script(h.Attrs("type", "module", "src", "https://cdn.jsdelivr.net/gh/starfederation/datastar@1.0.0-RC.7/bundles/datastar.js")),
        ),
        h.Body(
            h.Div(h.Attrs("id", "app"),
                h.Button(h.Attributes{
                    d.Signal("count", 0),
                    d.OnClick(d.SetSignal("count", d.Raw("$count + 1"))),
                },
                    h.Text("Count: "),
                    h.Span(h.Attributes{d.Text(d.Raw("$count"))}),
                ),
            ),
        ),
    )

    h.Render(os.Stdout, page)
}
```

### Datastar Pro

The `d` package also includes helpers for [Datastar Pro](https://data-star.dev/) features (requires commercial license):

- Animations: `Animate`
- Form validation: `CustomValidity`
- Observers: `OnRAF`, `OnResize`
- State persistence: `Persist`, `QueryString`
- URL management: `ReplaceURL`
- Scrolling: `ScrollIntoView`
- Transitions: `ViewTransitionName`
- Utility actions: `Clipboard`, `Fit`, `FitClamped`

## Benchmarks

htmlgen is benchmarked against Go's standard `html/template` package. Run benchmarks locally with:

```bash
go test -bench=. -benchmem ./h/

# Run with purego (no unsafe optimizations)
go test -bench=. -benchmem -tags=purego ./h/
```

The `purego` build tag disables unsafe pointer optimizations for environments that require pure Go code.

### Performance Comparison

| Scenario | htmlgen | htmlgen (purego) | html/template | Winner |
|----------|---------|------------------|---------------|--------|
| [Simple Div](h/benchmark_test.go#L16) | 151 ns | 152 ns | 521 ns | htmlgen ~3.5x faster |
| [Div with Attributes](h/benchmark_test.go#L37) | 303 ns | 349 ns | 2100 ns | htmlgen ~6.9x faster |
| [Nested Elements](h/benchmark_test.go#L68) | 1068 ns | 1095 ns | 2128 ns | htmlgen ~2.0x faster |
| [List (10 items)](h/benchmark_test.go#L104) | 887 ns | 1022 ns | 4762 ns | htmlgen ~5.4x faster |
| [List (100 items)](h/benchmark_test.go#L130) | 7.2 µs | 8.5 µs | 45.2 µs | htmlgen ~6.3x faster |
| [Table (10 rows)](h/benchmark_test.go#L166) | 7.2 µs | 7.7 µs | 17.1 µs | htmlgen ~2.4x faster |
| [Table (100 rows)](h/benchmark_test.go#L207) | 62.7 µs | 66.7 µs | 167.5 µs | htmlgen ~2.7x faster |
| [Full Page](h/benchmark_test.go#L270) | 4.9 µs | 5.2 µs | 11.0 µs | htmlgen ~2.3x faster |
| [Escaping](h/benchmark_test.go#L349) | 450 ns | 498 ns | 1437 ns | htmlgen ~3.2x faster |
| [Deep Nesting (10 levels)](h/benchmark_test.go#L376) | 1030 ns | 1030 ns | 530 ns | template ~1.9x faster |
| [Form](h/benchmark_test.go#L450) | 3.5 µs | 4.1 µs | 13.7 µs | htmlgen ~3.9x faster |
| [Pre-built Tree (static)](h/benchmark_test.go#L644) | 539 ns | 576 ns | 73 ns | template ~7.4x faster |
| [Compiled Tree (static)](h/benchmark_test.go#L675) | 19 ns | 19 ns | 73 ns | htmlgen ~3.8x faster |
| [Compiled Params](h/benchmark_test.go#L717) | 141 ns | 145 ns | 1122 ns | htmlgen ~8.0x faster |

*Benchmarks run on Apple M1 Ultra. Results may vary by hardware.*

### Key Insights

- **htmlgen is faster** for dynamic content generation with variable data structures
- **Compile** pre-renders static content for excellent performance (19 ns vs 539 ns)
- **CompileParams** is ~8.0x faster than html/template for parameterized content
- htmlgen excels at list/table generation where it can be 5-6x faster
- For attribute-heavy elements, htmlgen can be up to 7x faster
- **purego** adds ~5-15% overhead but remains significantly faster than html/template

### When to Use Each

| Use Case | Recommendation |
|----------|----------------|
| Dynamic lists/tables | htmlgen |
| Forms with many attributes | htmlgen |
| Full page generation with data | htmlgen |
| Static templates with no data | `Compile` |
| Parameterized templates | `CompileParams` |
| Component-based UI architecture | htmlgen |

## License

See LICENSE file for details.
