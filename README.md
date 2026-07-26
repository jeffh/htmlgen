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
```

```bash
# Run with purego (no unsafe optimizations)
go test -bench=. -benchmem -tags=purego ./h/
```

### Performance Comparison

| Scenario | htmlgen | htmlgen (purego) | html/template | Winner |
|----------|---------|------------------|---------------|--------|
| [Simple Div](h/benchmark_test.go#L20) | 79 ns | 80 ns | 489 ns | htmlgen ~6.2x faster |
| [Div with Attributes](h/benchmark_test.go#L45) | 248 ns | 294 ns | 1969 ns | htmlgen ~7.9x faster |
| [Nested Elements](h/benchmark_test.go#L80) | 535 ns | 548 ns | 1989 ns | htmlgen ~3.7x faster |
| [List (10 items)](h/benchmark_test.go#L120) | 940 ns | 1059 ns | 4489 ns | htmlgen ~4.8x faster |
| [List (100 items)](h/benchmark_test.go#L154) | 8.5 µs | 10.0 µs | 42.2 µs | htmlgen ~5.0x faster |
| [Table (10 rows)](h/benchmark_test.go#L232) | 4.0 µs | 4.4 µs | 15.4 µs | htmlgen ~3.8x faster |
| [Table (100 rows)](h/benchmark_test.go#L252) | 35.8 µs | 39.4 µs | 150.2 µs | htmlgen ~4.2x faster |
| [Full Page](h/benchmark_test.go#L352) | 3.4 µs | 3.6 µs | 10.3 µs | htmlgen ~3.1x faster |
| [Full Page (parallel)](h/benchmark_test.go#L580) | 1.5 µs | 1.5 µs | 5.5 µs | htmlgen ~3.7x faster |
| [Escaping](h/benchmark_test.go#L376) | 380 ns | 441 ns | 1276 ns | htmlgen ~3.4x faster |
| [Form](h/benchmark_test.go#L474) | 2.9 µs | 3.4 µs | 12.8 µs | htmlgen ~4.4x faster |
| [Deep Nesting (10 levels)](h/benchmark_test.go#L426) | 683 ns | 697 ns | 481 ns | template ~1.4x faster |

*Go 1.26.4, `darwin/arm64`, Apple M1 Ultra, `-count=8` medians via `benchstat`.
Results may vary by hardware and Go version.*

### Allocations

Streaming skips the intermediate node tree, so allocation counts track the data
being rendered rather than the shape of the document:

| Scenario | htmlgen | html/template |
|----------|---------|---------------|
| Simple Div | 0 B · 0 allocs | 240 B · 7 allocs |
| Div with Attributes | 120 B · 2 allocs | 576 B · 22 allocs |
| Nested Elements | 0 B · 0 allocs | 576 B · 22 allocs |
| List (100 items) | 2.4 KiB · 101 allocs | 12.7 KiB · 604 allocs |
| Table (100 rows) | 25.1 KiB · 402 allocs | 37.7 KiB · 1804 allocs |
| Full Page | 2.3 KiB · 41 allocs | 2.7 KiB · 122 allocs |
| Escaping | 56 B · 2 allocs | 736 B · 16 allocs |
| Form | 2.5 KiB · 37 allocs | 3.2 KiB · 148 allocs |

Static markup costs nothing: the simple-div and nested-element cases allocate
zero bytes, because a body closure that captures no variables is a static
function value.

### Streaming vs. the Previous Tree API

The streaming rewrite replaced an API that built a node tree before rendering.
Tree-API figures were re-measured from the last pre-rewrite commit (`95e46b0`)
on the same machine and Go version as everything above, with element
construction inside the timed loop on both sides:

| Scenario | Tree API | Streaming API | Change |
|----------|----------|---------------|--------|
| Simple Div | 150 ns · 3 allocs | 79 ns · 0 allocs | 1.9x faster |
| Div with Attributes | 301 ns · 5 allocs | 248 ns · 2 allocs | 1.2x faster |
| Nested Elements | 1004 ns · 20 allocs | 535 ns · 0 allocs | 1.9x faster |
| List (10 items) | 1695 ns · 35 allocs | 940 ns · 11 allocs | 1.8x faster |
| List (100 items) | 15.5 µs · 308 allocs | 8.5 µs · 101 allocs | 1.8x faster |
| Table (10 rows) | 6.7 µs · 130 allocs | 4.0 µs · 42 allocs | 1.6x faster |
| Table (100 rows) | 57.9 µs · 1123 allocs | 35.8 µs · 402 allocs | 1.6x faster |
| Full Page | 4.6 µs · 86 allocs | 3.4 µs · 41 allocs | 1.4x faster |
| Escaping | 449 ns · 5 allocs | 380 ns · 2 allocs | 1.2x faster |
| Deep Nesting (10 levels) | 1031 ns · 21 allocs | 683 ns · 10 allocs | 1.5x faster |
| Form | 3.4 µs · 56 allocs | 2.9 µs · 37 allocs | 1.1x faster |

Every scenario improved, by 1.1x to 1.9x, and allocations fell by 30-100%. The
gain is largest where the old API built the most nodes and smallest where the
work is dominated by escaping or attribute handling, which both APIs share.

One capability was lost. The tree API could pre-render a static subtree with
`Compile`, replaying it in ~18 ns with zero allocations — faster than anything
the streaming API can do, since streaming always re-walks the Go code. The
equivalent is now to render once with `RenderString` or `RenderBytes` and write
the cached result yourself.

### Entry-Point Overhead

| Operation | Cost |
|-----------|------|
| [`Render` call overhead (empty body)](h/benchmark_test.go#L531) | 13.1 ns · 0 B · 0 allocs |
| `RenderBytes`, small fragment | 349 ns · 80 B · 1 alloc |
| `RenderString`, small fragment | 433 ns · 280 B · 6 allocs |
| `RenderIndent`, small fragment | 502 ns · 16 B · 2 allocs |

`Render` itself is nearly free — the builder comes from a `sync.Pool`, so a
render that writes nothing allocates nothing. Writing to an `io.Writer` directly
is the cheapest path; `RenderBytes` adds one copy and `RenderString` adds
`strings.Builder` growth. Pretty-printing costs roughly 15% over compact output,
plus a small cached indent ladder.

### Key Insights

- **htmlgen is faster** for dynamic content, by 3-8x across most scenarios
- **Streaming beats the old tree API** on every benchmark, 1.1-1.9x, while
  cutting allocations 30-100%
- **Static markup is free**: body closures that capture nothing allocate nothing
- htmlgen excels at list and table generation, where it is ~4-5x faster
- For attribute-heavy elements, htmlgen is up to ~8x faster
- Concurrency helps, but sub-linearly: on 20 cores, `RunParallel` cuts per-render
  cost only ~2.3x (3.4 µs to 1.5 µs). The pooled builder is not the bottleneck —
  allocation and GC pressure from the rendered data is. The advantage over
  `html/template` holds at ~3.7x either way
- **purego** adds ~1-19% overhead but remains far faster than html/template

### When to Use Each

| Use Case | Recommendation |
|----------|----------------|
| Dynamic lists/tables | htmlgen |
| Forms with many attributes | htmlgen |
| Full page generation with data | htmlgen |
| Component-based UI architecture | htmlgen |
| Streaming to an `http.ResponseWriter` | htmlgen (`Render`) |
| Static markup rendered repeatedly | Render once, cache the bytes |
| Designer-edited templates, no recompile | `html/template` |

### Caveats

**Capturing closures allocate.** A body closure that captures a loop variable or
parameter is a heap allocation (~16 B each), which is why the deep-nesting case —
a recursive helper capturing its `depth` argument — is the one benchmark that
loses to `html/template`. Rendering many small elements in a loop pays this per
element; it is the reason the list and table benchmarks allocate roughly one
object per row. That is usually cheap next to the tree allocations it replaces,
but it does not disappear.

**`html/template` does more.** It escapes context-sensitively (HTML, CSS, JS,
URL) at execute time. htmlgen escapes text and attribute values only, and
`Raw`/`Rawf` are unescaped by contract. These numbers are a cost-of-output
comparison, not a claim of feature parity.

**Methodology.** The paired `_HtmlGen` / `_Template` benchmarks were checked to
emit byte-identical output before the ratios were recorded, and both write into
the same reusable `bytes.Buffer`, so they measure generation cost rather than
destination allocation.

For the tree-API comparison, the two List rows do not use that suite's published
numbers: its list benchmarks hoisted item construction out of the timed loop, so
they measured replaying a prebuilt tree rather than building and rendering one.
They were re-run with construction inside the loop to match the streaming
benchmarks. Every other tree-API row reproduced its previously published figure
to within a few percent.

### The `purego` build tag

The `purego` build tag disables the unsafe string-to-bytes conversion in the
escaper for environments that require pure Go. It costs about 4.5% geometric
mean across the suite, up to ~18% on attribute-heavy paths, and adds one
allocation per value that actually needs escaping — content with no escapable
characters takes the same fast path either way.

## License

See LICENSE file for details.
