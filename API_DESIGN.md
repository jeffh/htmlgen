# htmlgen streaming API design

Status: **shipped.** This began as the design spec for rewriting package `h`
from a tree-building `Builder` API to a **streaming, imperative** API inspired
by Clay (https://github.com/nicbarker/clay), and has been updated to describe
the API as it now stands. Native Go `if`/`for` are the control flow; a container
takes a **body closure** that runs immediately and streams children to the
writer.

Two revisions landed after the original rewrite and are folded in below:

- **Buffered output.** `B` accumulates into an internal buffer and writes to the
  `io.Writer` in ~4 KiB chunks instead of issuing a write per tag, attribute, and
  text run. `Flush` is exported for callers who need bytes delivered mid-render.
- **Typed element parameters.** Element methods originally took `...any` and
  unpacked it with a type switch. They now take concrete `Attributes` and `Body`
  parameters, so no argument is boxed and a capturing body closure stays on the
  stack. The variadic form is gone; see [Element methods](#element-methods).

## Goals / decisions (locked with the user)

1. **Element call style — "body receives handle".** Elements are **methods on
   `*h.B`** (a per-render streaming builder). A container takes its attributes
   and a trailing **body closure `func(b *h.B)`**. The body's `b` parameter is
   how nested elements are emitted. This keeps everything concurrency-safe (no
   package-level/ambient/goroutine-local state) — critical because the consumer
   (stdapps) renders concurrently across HTTP requests.

   ```go
   h.Render(w, func(b *h.B) {
       b.Div(h.Attrs("id", "container"), func(b *h.B) {
           for i := range 10 {
               b.Span(nil, func(b *h.B) { b.Textf("%d child element", i) })
           }
       })
   })
   ```

2. **Error handling — sticky error, returned by `Render`.** `*h.B` records the
   first write error and turns every subsequent element/text call into a no-op.
   User code never checks per-call errors inside loops. `Render` returns the
   sticky error at the end. No panics in normal control flow. (`RenderString`
   / `RenderBytes` still panic on error, preserving current behavior.)

3. **Replace the tree Builder API entirely.** Remove `Builder`, `TagArg`,
   `htmlTagBuilder`/`tagBuilder`/`fragmentBuilder`/`textBuilder`, and the
   combinators `If`, `When`, `Unless`, `First`, `ForEach`, `ForEach2`,
   `Fragment`. Native Go replaces them (`if`, `for`, plain sequential calls,
   helper funcs). Keep the attribute layer intact (below).

4. **Keep the attribute layer** so `ds`/`hx`/`js` keep working: `Attribute`,
   `Attributes` (+ its methods), `Attr`, `Attrs`, `AttrsMap`, `AttrIf`, and the
   `AttrBuilder` interface. These are values, not control flow, so they carry
   over as-is. Collecting a mix of them into one `Attributes` is what `AttrsOf`
   and `Attributes.With` do (below).

## Core type

```go
// B is a streaming HTML builder bound to an io.Writer. It tracks open tags,
// optional pretty-printing, and a sticky write error. Output accumulates in an
// internal buffer and is written to the io.Writer in chunks. Not safe for
// concurrent use by multiple goroutines; each Render gets its own *B.
type B struct {
    w           io.Writer
    buf         []byte // flushed to w at ~4 KiB
    scratch     []byte // one escaped chunk of an oversized value
    openTags    []string
    indent      string
    indentCache []string
    atLineStart bool
    maxLineLen  int
    err         error // sticky; first write error wins
}

// Err returns the first write error observed as of the last flush.
func (b *B) Err() error

// Flush writes buffered output to the underlying io.Writer. Render and
// RenderIndent flush before returning; call it explicitly when bytes must reach
// the client mid-render (server-sent events, long streaming responses).
func (b *B) Flush() error

// SetMaxLineLength sets the width at which attributes wrap when pretty-printing.
// 0 (the default) disables wrapping.
func (b *B) SetMaxLineLength(maxLen int)
```

The low-level writing/escaping/indent logic lives in `writer.go` (a 256-entry
escape table appended straight into the buffer, `writeAttrs`, indent cache,
open-tag stack). There is **no exported `Builder` interface** and no separate
exported `Writer` tree-walker; `Render` takes a `func(*B)`.

Because output is buffered, `Err` reflects only errors seen as of the last
flush — `Render`'s return value is always current. Values at or above the flush
threshold (a large `Raw`, `Text`, or attribute value) bypass the buffer and
stream through in bounded chunks, so a multi-megabyte value never costs a
multi-megabyte buffer.

## Element methods

Every HTML element in `tags.go` is a method on `*B` with the same Go name
(`Div`, `Span`, `A`, `H1`, `Table`, `Button`, `Input`, `Img`, …).

Signatures are concrete — one shape for containers, one for void elements:

```go
func (b *B) Div(attrs Attributes, body Body)   // container element
func (b *B) Img(attrs Attributes)              // void element: attributes only
```

- `attrs` is the element's attributes, or `nil` for none.
- `body` runs between the open and close tags, or `nil` for an empty element.
  Void elements take no body parameter at all, so passing one is a compile
  error rather than a silently ignored argument.

The `Body` alias keeps this readable and lets components declare a body
parameter:

```go
type Body = func(*B)
```

**On the earlier variadic design.** Elements originally took `args ...any` and
type-switched over `Attributes` / `Attribute` / `AttrBuilder` / `func(*B)` /
`nil`, so a `func(b *h.B){…}` literal could be passed directly alongside
attributes in any order. That reads well, but every call allocated a `[]any`,
boxed each argument into an interface, and forced capturing body closures onto
the heap. Typed siblings (`DivE`, `ImgE`, …) were added as a fast path, then the
variadic form was dropped and the typed signatures took the plain names —
carrying two spellings of every element was the worse trade. Removing the boxing
took the list and table benchmarks to zero allocations and roughly halved them.

The cost is that an `Attribute` or a `ds`/`hx` builder can no longer be handed
straight to an element; `AttrsOf` and `Attributes.With` collect them first (see
[Attribute layer](#attribute-layer)).

### Void / self-closing elements

`Area, Base, Br, Col, Embed, Hr, Img, Input, Link, Meta, Source, Track, Wbr`
render self-closing and take attributes only.

### Root document

```go
// Writes <!DOCTYPE html> then <html lang="en"> ... </html>, defaulting
// lang="en" unless attrs provides one.
func (b *B) Html(attrs Attributes, body Body)
func (b *B) Doctype()
```

`Html` appends its `lang` default with a full slice expression, so it can never
write into a caller slice's spare capacity.

### Text / raw content (leaf statements)

```go
func (b *B) Text(s string)                    // HTML-escaped
func (b *B) Textf(format string, a ...any)    // HTML-escaped, fmt.Sprintf
func (b *B) Raw(s string)                      // unescaped (caller-sanitized)
func (b *B) Rawf(format string, a ...any)      // unescaped, fmt.Sprintf
```

### Custom elements

```go
func (b *B) El(name string, attrs Attributes, body Body)  // arbitrary container tag
func (b *B) VoidEl(name string, attrs Attributes)         // arbitrary void tag
```

(Replaces `CustomElement`.) Tag names are written verbatim, never escaped, so
both panic unless `name` matches `[A-Za-z][A-Za-z0-9_.:-]*` — a name derived
from untrusted input could otherwise break out of its element context. The
constant-name methods (`Div`, `Span`, …) skip this check, so the common path
costs nothing at runtime.

## Render entry points

```go
// Render runs fn against a fresh *B writing to w, then returns the sticky error.
// It closes any tags fn left open (defensive), matching well-formed output.
func Render(w io.Writer, fn func(*B)) error

// RenderIndent is Render with pretty-printing using the given indent unit.
func RenderIndent(w io.Writer, indent string, fn func(*B)) error

// RenderString renders fn to a string. Panics on write error (writes to a
// strings.Builder, which never errors, so this is effectively panic-free).
func RenderString(fn func(*B)) string

// RenderBytes renders fn to a byte slice. Panics on write error.
func RenderBytes(fn func(*B)) []byte
```

`*B` and its buffer come from a `sync.Pool`, so a render that writes nothing
allocates nothing. A `B` is valid only for the duration of the `Render` call
that created it — callers must not retain it or call `Flush` on it afterwards.
An oversized buffer is dropped rather than carried back into the pool, so one
large render does not pin that memory for the process lifetime.

## Components / composition

A component is just a function taking `*B` (plus its own params). Composition is
a plain call — no wrapper types:

```go
func Card(b *h.B, title string, body h.Body) {
    b.Div(h.Attrs("class", "card"), func(b *h.B) {
        b.H2(nil, func(b *h.B) { b.Text(title) })
        body(b)
    })
}

h.Render(w, func(b *h.B) {
    Card(b, "Hello", func(b *h.B) {
        b.P(nil, func(b *h.B) { b.Text("world") })
    })
})
```

Conditionals and iteration use native Go:

```go
b.Ul(nil, func(b *h.B) {
    for _, it := range items {
        if it.Visible {
            b.Li(nil, func(b *h.B) { b.Text(it.Name) })
        }
    }
})
```

## Attribute layer

`Attribute`, `Attributes` (with `Get/Index/Set/SetDefault/Delete/Merge`), `Attr`,
`Attrs`, `AttrsMap`, `AttrIf`, `AttrBuilder` all carry over from the tree API.
The `isTagArg()` marker interface is gone.

Element methods take a single `Attributes`, so combining loose `Attribute`
values with the fluent `ds`/`hx`/`js` builders needs a collector:

```go
func AttrsOf(items ...AttrBuilder) Attributes
func (a Attributes) With(items ...AttrBuilder) Attributes
func (a Attribute) Attribute() Attribute  // Attribute satisfies AttrBuilder
```

```go
b.Div(h.AttrsOf(h.Attr("id", "app"), ds.Signal("count", 0)), body)
b.Button(h.Attrs("class", "btn").With(hx.Get("/api/data")), body)
```

Both carry the merge semantics the variadic type switch had: a later value
overrides an earlier one of the same name **without changing its position**,
zero attributes (an `AttrIf` whose condition was false) and `nil` builders are
skipped, and neither call mutates its inputs — a caller's `Attributes` is never
written through, including into its spare capacity.

Attribute names are written verbatim, so `Attr`, `AttrIf`, `Attrs`, `AttrsMap`,
`Set`, and `SetDefault` validate them against `[A-Za-z][A-Za-z0-9_.:-]*` and
panic otherwise. `Attribute` struct literals bypass validation and are trusted;
`AttrsOf`/`With` do not re-validate, since their inputs are already-constructed
`Attribute` values.

## Tests

`h`'s suite (`api_test.go`, `helpers_test.go`, `validate_test.go`,
`benchmark_test.go`, `benchmark_nested_test.go`) covers: escaping, attribute
merging/ordering, `AttrsOf`/`With` semantics, void elements, doctype/html
defaulting, name validation, indent + max-line wrapping, mid-render `Flush`,
large values bypassing the buffer, sticky-error propagation (a failing
io.Writer makes `Render` return the error and stop output), and pooling.

Two properties are worth guarding explicitly because they are easy to
regress:

- **Caller `Attributes` are never mutated**, including through spare capacity —
  see `TestCallerAttributesNotMutated`.
- **A capturing body closure does not allocate** — see
  `TestElementBodyDoesNotAllocate`, which is what makes the typed signatures
  worth the ergonomic cost.

The `_HtmlGen` / `_Template` benchmark pairs are checked to emit identical
output (`TestCardGridOutputsMatch`) so the published ratios measure generation
cost, not a difference in what is generated.

Keep `go build ./... && go vet ./... && go test ./...` green.
