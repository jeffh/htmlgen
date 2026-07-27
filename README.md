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
intermediate tree.

Every container element takes attributes and a body — `Div(attrs h.Attributes,
body h.Body)` — and every void element takes attributes alone — `Img(attrs
h.Attributes)`. `nil` is valid for either:

```go
import "github.com/jeffh/htmlgen/h"

err := h.Render(os.Stdout, func(b *h.B) {
    b.Html(nil, func(b *h.B) {
        b.Head(nil, func(b *h.B) {
            b.Title(nil, func(b *h.B) { b.Text("My Page") })
        })
        b.Body(nil, func(b *h.B) {
            b.Div(h.Attrs("class", "container"), func(b *h.B) {
                b.H1(nil, func(b *h.B) { b.Text("Hello, World!") })
                b.P(nil, func(b *h.B) { b.Text("Welcome to htmlgen.") })
                b.A(h.Attrs("href", "/about"), func(b *h.B) {
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
    b.Strong(nil, func(b *h.B) { b.Text("Rendered to a string") })
})
```

Concrete parameter types mean nothing is boxed into an interface, so a body
closure that captures a loop variable stays on the stack:

```go
b.Div(h.Attrs("class", "grid"), func(b *h.B) {
    for _, card := range cards {
        b.Div(h.Attrs("class", "card"), func(b *h.B) {
            b.H3(nil, func(b *h.B) { b.Text(card.Title) })
            b.Img(h.Attrs("src", card.Image))
        })
    }
})
```

Use `RenderIndent` for pretty-printed output and `RenderBytes` for an in-memory
byte slice. The first write error is sticky: later output calls become no-ops,
and `Render` returns that error.

Output is buffered and written to the `io.Writer` in ~4 KiB chunks. `Render`,
`RenderIndent`, `RenderString`, and `RenderBytes` all flush before returning, so
most code never thinks about it. Call `Flush` yourself when bytes must reach the
client before the render finishes — server-sent events, or a long response you
want to start streaming early:

```go
err := h.Render(w, func(b *h.B) {
    b.Div(h.Attrs("class", "feed"), func(b *h.B) {
        for event := range events {
            b.P(nil, func(b *h.B) { b.Text(event) })
            if err := b.Flush(); err != nil {
                return
            }
            if f, ok := w.(http.Flusher); ok {
                f.Flush()
            }
        }
    })
})
```

Because writes are buffered, `Err` reports only the errors seen as of the last
flush; the value returned by `Render` is always current. Calling `Flush` at loop
boundaries, as above, is also how you notice a disconnected client promptly
instead of at the end of the render.

Two consequences of buffering are worth knowing. If a body closure panics,
whatever is still buffered — up to the flush threshold — is discarded rather
than written. And the `*h.B` is valid only for the duration of the `Render` call
that created it, because it returns to a pool afterwards; do not retain it or
call `Flush` on it later. Values at or above the threshold (a large `Raw`, a
large `Text`, a large attribute value) bypass the buffer and stream straight to
the writer, so rendering a multi-megabyte value does not cost a multi-megabyte
buffer.

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

Element methods take an `Attributes` (or `nil`). Single `Attribute` values and
the fluent builders from `ds`, `hx`, and `js` implement `AttrBuilder`, so
`AttrsOf` collects them into one `Attributes`, and the `With` method adds them
to an existing set:

```go
b.Div(h.AttrsOf(h.Attr("id", "app"), ds.Signal("count", 0)), body)

b.Button(h.Attrs("class", "btn").With(
    hx.Get("/api/data"),
    hx.Target("#results"),
), body)
```

Later values override earlier values of the same name, keeping the earlier
position. Zero attributes — an `AttrIf` whose condition was false, or a `nil`
builder — are skipped, and neither `AttrsOf` nor `With` modifies its inputs.

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

Void HTML elements such as `Img`, `Input`, and `Br` are self-closing and take no
body parameter. `Html` emits the HTML5 doctype and defaults to `lang="en"`.

### Security

Text and attribute *values* are always HTML-escaped (`Raw`/`Rawf` are the
explicit opt-outs). Attribute *names* and custom tag names are written
verbatim, so `Attr`, `Attrs`, `AttrsMap`, `AttrIf`, `Set`, `SetDefault`, `El`,
and `VoidEl` validate names against `[A-Za-z][A-Za-z0-9_.:-]*` and panic on
anything else, preventing injection through a name that breaks out of its
attribute or element context. `Attribute` struct literals bypass this
validation and are trusted — never build one with a name derived from
untrusted input.

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
        b.Html(nil, func(b *h.B) {
            b.Head(nil, func(b *h.B) {
                b.Title(nil, func(b *h.B) { b.Text("Counter") })
                b.Script(h.Attrs("type", "module", "src", "https://cdn.jsdelivr.net/gh/starfederation/datastar@v1.0.2/bundles/datastar.js"), nil)
            })
            b.Body(nil, func(b *h.B) {
                b.Div(
                    h.AttrsOf(h.Attr("id", "app"), ds.Signal("count", 0)),
                    func(b *h.B) {
                        b.Button(
                            h.AttrsOf(ds.OnClick(ds.SetSignal("count", ds.Raw("$count + 1")))),
                            func(b *h.B) { b.Text("Count: ") },
                        )
                        b.Span(h.AttrsOf(ds.Text(ds.Raw("$count"))), nil)
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

### Performance Comparison

| Scenario | htmlgen | html/template | Winner |
|----------|---------|---------------|--------|
| [Card Grid (~90 elements)](h/benchmark_nested_test.go#L106) | 5.8 µs | 54.4 µs | htmlgen ~9.4x faster |
| [Simple Div](h/benchmark_test.go#L20) | 41 ns | 491 ns | htmlgen ~11.9x faster |
| [Div with Attributes](h/benchmark_test.go#L45) | 127 ns | 1975 ns | htmlgen ~15.5x faster |
| [Nested Elements](h/benchmark_test.go#L80) | 172 ns | 1997 ns | htmlgen ~11.6x faster |
| [List (10 items)](h/benchmark_test.go#L120) | 258 ns | 4.5 µs | htmlgen ~17.5x faster |
| [List (100 items)](h/benchmark_test.go#L154) | 2.2 µs | 41.8 µs | htmlgen ~18.6x faster |
| [Table (10 rows)](h/benchmark_test.go#L232) | 1.1 µs | 15.3 µs | htmlgen ~14.2x faster |
| [Table (100 rows)](h/benchmark_test.go#L252) | 9.4 µs | 150.7 µs | htmlgen ~16.1x faster |
| [Full Page](h/benchmark_test.go#L352) | 1.1 µs | 10.1 µs | htmlgen ~9.6x faster |
| [Full Page (parallel)](h/benchmark_test.go#L580) | 382 ns | 6.5 µs | htmlgen ~16.9x faster |
| [Escaping](h/benchmark_test.go#L376) | 227 ns | 1322 ns | htmlgen ~5.8x faster |
| [Form](h/benchmark_test.go#L474) | 1.3 µs | 13.0 µs | htmlgen ~9.7x faster |
| [Deep Nesting (10 levels)](h/benchmark_test.go#L426) | 178 ns | 481 ns | htmlgen ~2.7x faster |

*Go 1.26.4, `darwin/arm64`, Apple M1 Ultra, `-count=8` medians via `benchstat`.
Results may vary by hardware and Go version.*

### Allocations

Streaming skips the intermediate node tree, and concrete element parameters mean
nothing is boxed, so what is left to allocate is whatever the caller itself
builds — mostly `Attrs` slices:

| Scenario | htmlgen | html/template |
|----------|---------|---------------|
| Card Grid (~90 elements) | 2.8 KiB · 72 allocs | 14.4 KiB · 655 allocs |
| Simple Div | 0 B · 0 allocs | 240 B · 7 allocs |
| Div with Attributes | 96 B · 1 alloc | 576 B · 22 allocs |
| Nested Elements | 0 B · 0 allocs | 576 B · 22 allocs |
| List (100 items) | 0 B · 0 allocs | 12.7 KiB · 604 allocs |
| Table (100 rows) | 0 B · 0 allocs | 37.7 KiB · 1804 allocs |
| Full Page | 352 B · 8 allocs | 2.7 KiB · 122 allocs |
| Escaping | 32 B · 1 alloc | 736 B · 16 allocs |
| Form | 1.6 KiB · 16 allocs | 3.2 KiB · 148 allocs |

Static markup costs nothing: the simple-div and nested-element cases allocate
zero bytes, because a body closure that captures no variables is a static
function value. The list and table benchmarks are also allocation-free — their
bodies capture loop variables, but a `Body` parameter keeps the closure on the
stack.

### Buffered Output

Element methods append to an internal buffer that flushes to the `io.Writer`
every ~4 KiB, instead of issuing an `io.Writer` call per tag, attribute, and
text run. A single small element used to cost a dozen writes; it now costs one
`append` each and shares a flush with its neighbors. That roughly halved every
fine-grained benchmark when it landed — simple div 82 ns to 43 ns, nested
elements 532 ns to 195 ns, measured before element methods dropped `...any` —
and it also removed the previous `unsafe` string-to-bytes conversion
in the escaper, since escaping now appends into the buffer directly. The escape
path itself is a 256-entry lookup table rather than a byte switch. Values at or
above the threshold skip the buffer and stream through in bounded chunks, so
buffering never trades memory for speed on large content.

`Flush` is exported for callers who need bytes delivered before the render ends;
see [Streaming API](#streaming-api).

### Streaming vs. the Previous Tree API

The streaming rewrite replaced an API that built a node tree before rendering.
Tree-API figures were re-measured from the last pre-rewrite commit (`95e46b0`)
on the same machine and Go version as everything above, with element
construction inside the timed loop on both sides:

| Scenario | Tree API | Streaming API | Change |
|----------|----------|---------------|--------|
| Simple Div | 150 ns · 3 allocs | 41 ns · 0 allocs | 3.7x faster |
| Div with Attributes | 301 ns · 5 allocs | 127 ns · 1 alloc | 2.4x faster |
| Nested Elements | 1004 ns · 20 allocs | 172 ns · 0 allocs | 5.8x faster |
| List (10 items) | 1695 ns · 35 allocs | 258 ns · 0 allocs | 6.6x faster |
| List (100 items) | 15.5 µs · 308 allocs | 2.2 µs · 0 allocs | 6.9x faster |
| Table (10 rows) | 6.7 µs · 130 allocs | 1.1 µs · 0 allocs | 6.2x faster |
| Table (100 rows) | 57.9 µs · 1123 allocs | 9.4 µs · 0 allocs | 6.2x faster |
| Full Page | 4.6 µs · 86 allocs | 1.1 µs · 8 allocs | 4.4x faster |
| Escaping | 449 ns · 5 allocs | 227 ns · 1 alloc | 2.0x faster |
| Deep Nesting (10 levels) | 1031 ns · 21 allocs | 178 ns · 0 allocs | 5.8x faster |
| Form | 3.4 µs · 56 allocs | 1.3 µs · 16 allocs | 2.5x faster |

Every scenario improved, by 2.0x to 6.9x, and allocations fell by 71-100%. The
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
| [`Render` call overhead (empty body)](h/benchmark_test.go#L531) | 14.4 ns · 0 B · 0 allocs |
| `RenderBytes`, small fragment | 145 ns · 80 B · 1 alloc |
| `RenderString`, small fragment | 153 ns · 112 B · 2 allocs |
| `RenderIndent`, small fragment | 212 ns · 16 B · 2 allocs |

`Render` itself is nearly free — the builder comes from a `sync.Pool`, so a
render that writes nothing allocates nothing. Buffering also makes the in-memory
entry points cheap: a small fragment reaches `strings.Builder` or `bytes.Buffer`
as one write, so neither has to grow. Pretty-printing costs roughly 40% over
compact output, plus a small cached indent ladder.

### Key Insights

- **htmlgen is faster** for dynamic content, by 6-19x across most scenarios
- **Nothing is boxed**: element methods take concrete `Attributes` and `Body`
  parameters, so a capturing body closure stays on the stack and the only
  allocations left are the `Attrs` slices the caller builds
- **Buffering dominates the small cases**: writing whole tags into a 4 KiB
  buffer instead of per-fragment `io.Writer` calls roughly halved the cost of
  every fine-grained benchmark
- **Streaming beats the old tree API** on every benchmark, 2.0-6.9x, while
  cutting allocations 71-100%
- **Static markup is free**: body closures that capture nothing allocate nothing,
  and closures that capture only loop variables stay on the stack
- htmlgen excels at list and table generation, where it is ~14-19x faster and
  allocates nothing
- For attribute-heavy elements, htmlgen is up to ~16x faster
- Concurrency scales well now that rendering allocates so little: on 20 cores,
  `RunParallel` cuts per-render cost ~2.8x (1.05 µs to 382 ns), up from ~1.3x
  when every element call boxed its arguments

### When to Use Each

| Use Case | Recommendation |
|----------|----------------|
| Dynamic lists/tables | htmlgen |
| Forms with many attributes | htmlgen |
| Full page generation with data | htmlgen |
| Component-based UI architecture | htmlgen |
| Streaming to an `http.ResponseWriter` | htmlgen (`Render`, plus `Flush` for SSE) |
| Static markup rendered repeatedly | Render once, cache the bytes |
| Designer-edited templates, no recompile | `html/template` |

### Caveats

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

## License

See LICENSE file for details.
