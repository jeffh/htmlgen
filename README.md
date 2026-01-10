# htmlgen

A Go library for programmatic HTML generation with first-class [Datastar](https://data-star.dev/) support.

## Installation

```bash
go get github.com/jeffh/htmlgen
```

## Overview

htmlgen provides two packages:

- **`h`** - Core HTML generation with both streaming and declarative APIs
- **`d`** - Datastar attribute helpers for building reactive web applications

## Package `h` - HTML Generation

### Declarative Builder API

Build HTML trees using Go functions that mirror HTML elements:

```go
import "github.com/jeffh/htmlgen/h"

page := h.Html(nil,
    h.Head(nil,
        h.Title(nil, h.Text("My Page")),
    ),
    h.Body(nil,
        h.Div(h.Attrs("class", "container"),
            h.H1(nil, h.Text("Hello, World!")),
            h.P(nil, h.Text("Welcome to htmlgen.")),
            h.A(h.Attrs("href", "/about"), h.Text("About")),
        ),
    ),
)

// Render to any io.Writer
h.Render(os.Stdout, page)

// Or render with pretty-printed indentation (using two spaces)
h.RenderPretty(os.Stdout, "  ", page)
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
    page := h.Html(nil,
        h.Head(nil,
            h.Title(nil, h.Text("Counter")),
            h.Script(h.Attrs("type", "module", "src", "https://cdn.jsdelivr.net/npm/@starfederation/datastar")),
        ),
        h.Body(nil,
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

## Platform Detection

The `h` package provides build-tag controlled constants for detecting runtime environment:

```go
if h.Server {
    // Server-side Go
}
if h.Client {
    // WebAssembly in browser
}
```

## Benchmarks

htmlgen is benchmarked against Go's standard `html/template` package. Run benchmarks locally with:

```bash
go test -bench=. -benchmem ./h/
```

### Performance Comparison

| Scenario | htmlgen | html/template | Winner |
|----------|---------|---------------|--------|
| Simple Div | 332 ns | 515 ns | htmlgen ~1.5x faster |
| Div with Attributes | 693 ns | 2105 ns | htmlgen ~3x faster |
| Nested Elements | 1890 ns | 2089 ns | htmlgen ~1.1x faster |
| List (10 items) | 1982 ns | 4732 ns | htmlgen ~2.4x faster |
| List (100 items) | 16.5 µs | 44.9 µs | htmlgen ~2.7x faster |
| Table (10 rows) | 11.7 µs | 16.9 µs | htmlgen ~1.4x faster |
| Table (100 rows) | 97.9 µs | 163 µs | htmlgen ~1.7x faster |
| Full Page | 8.3 µs | 10.8 µs | htmlgen ~1.3x faster |
| Escaping | 1081 ns | 1421 ns | htmlgen ~1.3x faster |
| Deep Nesting (static) | 1819 ns | 512 ns | template ~3.5x faster |
| Form | 6.6 µs | 13.5 µs | htmlgen ~2x faster |
| Pre-built Tree (static) | 1361 ns | 68 ns | template ~20x faster |

*Benchmarks run on Apple M1 Ultra. Results may vary by hardware.*

### Key Insights

- **htmlgen is faster** for dynamic content generation with variable data structures
- **html/template is faster** for static content with no runtime data substitution
- htmlgen excels at list/table generation where it can be 2-3x faster
- For attribute-heavy elements, htmlgen can be up to 3x faster
- The Writer API is faster than the Builder API when you don't need declarative syntax

### When to Use Each

| Use Case | Recommendation |
|----------|----------------|
| Dynamic lists/tables | htmlgen |
| Forms with many attributes | htmlgen |
| Full page generation with data | htmlgen |
| Static templates with no data | html/template |
| Reusing same template many times | html/template |
| Component-based UI architecture | htmlgen |

## License

See LICENSE file for details.
