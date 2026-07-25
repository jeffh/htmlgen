# htmlgen

**Experimental**

[![Go Reference](https://pkg.go.dev/badge/github.com/jeffh/htmlgen.svg)](https://pkg.go.dev/github.com/jeffh/htmlgen)

A Go library for programmatic HTML generation.

Experimental, API is subject to change.

## Installation

```bash
go get github.com/jeffh/htmlgen
```

## Overview

htmlgen provides four packages:

- **`h`** - Streaming HTML generation
- **`ds`** - Datastar attribute helpers for building reactive web applications
- **`hx`** - HTMX attribute helpers
- **`js`** - Type-safe JavaScript generation for event handler attributes

## Package `h` - HTML Generation

### Streaming API

Render HTML imperatively through a per-render `*h.B`. Container bodies execute
immediately, so normal Go control flow and component functions work without an
intermediate tree:

```go
import "github.com/jeffh/htmlgen/h"

err := h.Render(os.Stdout, func(b *h.B) {
    b.Html(func(b *h.B) {
        b.Head(func(b *h.B) {
            b.Title(func(b *h.B) { b.Text("My Page") })
        })
        b.Body(func(b *h.B) {
            b.Div(h.Attrs("class", "container"), func(b *h.B) {
                b.H1(func(b *h.B) { b.Text("Hello, World!") })
                b.P(func(b *h.B) { b.Text("Welcome to htmlgen.") })
                b.A(h.Attr("href", "/about"), func(b *h.B) {
                    b.Text("About")
                })
            })
        })
    })
})
if err != nil {
    panic(err)
}

html := h.RenderString(func(b *h.B) {
    b.Strong(func(b *h.B) { b.Text("Rendered to a string") })
})
```

Use `RenderIndent` for pretty-printed output and `RenderBytes` for an in-memory
byte slice. The first write error is sticky: later output calls become no-ops,
and `Render` returns that error.

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

Element arguments may be `Attributes`, an `Attribute`, an `AttrBuilder` from a
companion package, a trailing `func(*h.B)` body, or `nil`. Later attribute
values override earlier values.

### Available Elements

All standard HTML5 elements are available as methods on `*h.B`:

- **Document**: `Html`, `Head`, `Title`, `Meta`, `Link`, `Style`, `Script`, `Body`
- **Sections**: `Header`, `Footer`, `Main`, `Nav`, `Section`, `Article`, `Aside`
- **Headings**: `H1`, `H2`, `H3`, `H4`, `H5`, `H6`
- **Text**: `P`, `Span`, `Div`, `Pre`, `Code`, `Em`, `Strong`, `A`
- **Lists**: `Ul`, `Ol`, `Li`, `Dl`, `Dt`, `Dd`
- **Tables**: `Table`, `Thead`, `Tbody`, `Tfoot`, `Tr`, `Th`, `Td`
- **Forms**: `Form`, `Input`, `Button`, `Label`, `Select`, `Option`, `Textarea`, `Fieldset`
- **Media**: `Img`, `Video`, `Audio`, `Picture`, `Source`, `Canvas`, `Svg`
- **Content and custom tags**: `Text`, `Textf`, `Raw`, `Rawf`, `El`, `VoidEl`

Void HTML elements such as `Img`, `Input`, and `Br` are self-closing and do not
run body closures. `Html` emits the HTML5 doctype and defaults to `lang="en"`.

## Package `ds` - Datastar Integration

Build reactive attributes for [Datastar](https://data-star.dev/) applications:

### Signals

```go
import "github.com/jeffh/htmlgen/ds"

// Define reactive signals
ds.Signal("count", 0)           // data-signals:count="0"
ds.Signal("name", "Alice")      // data-signals:name="\"Alice\""
ds.Signals(map[string]any{      // data-signals="{...}"
    "x": 1,
    "y": 2,
})

// Two-way binding
ds.Bind("username")             // data-bind="username"
```

### Event Handlers

Event-handler builders are fluent: actions are passed at construction, modifiers chain after.

```go
// Click events
ds.OnClick(ds.SetSignal("count", ds.Raw("$count + 1")))

// Form events
ds.OnSubmit(ds.Post("/api/submit")).PreventDefault()

// Other events
ds.OnInput(ds.SetSignal("search", ds.Raw("evt.target.value"))).Debounce(300*time.Millisecond)
ds.OnChange(ds.Get("/api/update"))
ds.Init(ds.Get("/api/init"))  // data-init: runs when the element loads
ds.On("keydown", ds.Raw("handleKey(evt)"))

// Intersection and interval observers
ds.OnIntersect(ds.Raw("$seen = true")).Once()
ds.OnInterval(ds.Raw("$tick++")).Duration(1*time.Second)
```

### HTTP Actions

```go
ds.Get("/api/data")
ds.Post("/api/submit")
ds.Put("/api/update")
ds.Delete("/api/remove")

// With options
ds.PostWithOptions("/api/submit",
    ds.RequestOptions().
        ContentType("json").
        Headers(map[string]string{"X-Custom": "value"}),
)
```

### Reactive Display

```go
ds.Show(ds.Raw("$isVisible"))                    // data-show="$isVisible"
ds.Text(ds.Raw("$message"))                      // data-text="$message"
ds.Class("active", ds.Raw("$isActive"))          // data-class:active="$isActive"
ds.Style("color", ds.Raw("$textColor"))          // data-style:color="$textColor"
ds.Attribute("disabled", ds.Raw("$isDisabled"))  // data-attr:disabled="$isDisabled"

// Multiple classes/styles/attrs at once
ds.Classes(map[string]string{"hidden": "$foo", "bold": "$bar"})
ds.Styles(map[string]string{"color": "$red ? 'red' : 'blue'"})
```

### Event Modifiers

Modifiers are methods on event builders, not standalone functions:

```go
ds.OnClick(ds.Raw("$x++")).PreventDefault()
ds.OnInput(ds.Raw("$y++")).Debounce(300 * time.Millisecond)
ds.On("scroll", ds.Raw("$z++")).Throttle(100 * time.Millisecond)
ds.OnClick(ds.Raw("$a++")).Delay(500 * time.Millisecond)
ds.OnClick(ds.Raw("$b++")).Once()
ds.OnClick(ds.Raw("$c++")).ViewTransition()
```

### Complete Example

```go
package main

import (
    "os"

    "github.com/jeffh/htmlgen/h"
    "github.com/jeffh/htmlgen/ds"
)

func main() {
    if err := h.Render(os.Stdout, func(b *h.B) {
        b.Html(func(b *h.B) {
            b.Head(func(b *h.B) {
                b.Title(func(b *h.B) { b.Text("Counter") })
                b.Script(h.Attrs("type", "module", "src", "https://cdn.jsdelivr.net/gh/starfederation/datastar@v1.0.2/bundles/datastar.js"))
            })
            b.Body(func(b *h.B) {
                b.Div(
                    h.Attr("id", "app"),
                    ds.Signal("count", 0),
                    func(b *h.B) {
                        b.Button(
                            ds.OnClick(ds.SetSignal("count", ds.Raw("$count + 1"))),
                            func(b *h.B) { b.Text("Count: ") },
                        )
                        b.Span(ds.Text(ds.Raw("$count")))
                    },
                )
            })
        })
    }); err != nil {
        panic(err)
    }
}
```

### Datastar Pro

The `ds` package also includes helpers for [Datastar Pro](https://data-star.dev/) features (requires commercial license):

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

The benchmark suite covers nested output, iteration, escaping, and allocation
counts for the pooled streaming builder. Results vary by Go version and
hardware; run the suite locally for current measurements.

## License

See LICENSE file for details.
