# htmlgen streaming API redesign

Status: design spec for a full rewrite of package `h` from a tree-building
`Builder` API to a **streaming, imperative** API inspired by Clay
(https://github.com/nicbarker/clay). Native Go `if`/`for` become the control
flow; a container is a function whose **body closure** runs immediately and
streams children to the writer.

## Goals / decisions (locked with the user)

1. **Element call style — "body receives handle".** Elements are **methods on
   `*h.B`** (a per-render streaming builder). A container takes attribute args
   and a trailing **body closure `func(b *h.B)`**. The body's `b` parameter is
   how nested elements are emitted. This keeps everything concurrency-safe (no
   package-level/ambient/goroutine-local state) — critical because the consumer
   (stdapps) renders concurrently across HTTP requests.

   ```go
   h.Render(w, func(b *h.B) {
       b.Div(h.Attrs("id", "container"), func(b *h.B) {
           for i := range 10 {
               b.Span(func(b *h.B) { b.Textf("%d child element", i) })
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

4. **Keep the attribute layer** so `ds`/`hx`/`js` keep working unchanged:
   `Attribute`, `Attributes` (+ its methods), `Attr`, `Attrs`, `AttrsMap`,
   `AttrIf`, and the `AttrBuilder` interface. These are values passed as
   element args; they are not control flow, so they carry over as-is.

## Core type

```go
// B is a streaming HTML builder bound to an io.Writer. It tracks open tags,
// optional pretty-printing, and a sticky write error. Not safe for concurrent
// use by multiple goroutines; each Render gets its own *B.
type B struct {
    w           io.Writer
    openTags    []string
    indent      string
    indentCache []string
    atLineStart bool
    maxLineLen  int
    err         error // sticky; first write error wins
}

// Err returns the first write error encountered, if any.
func (b *B) Err() error { return b.err }
```

Reuse the existing low-level writing/escaping/indent logic from `writer.go`
(escaping tables, `writeAttrs`, indent cache, open-tag stack). The current
`Writer` type may be folded into `B` or kept private as an internal helper —
implementer's choice — but there must be **no exported `Builder` interface** and
no separate exported `Writer` tree-walker. `Render` no longer takes a `Builder`.

## Element methods

Every HTML element in the current `tags.go` becomes a method on `*B` with the
same Go name (`Div`, `Span`, `A`, `H1`, `Table`, `Button`, `Input`, `Img`, …).

Signature (variadic, accepts the union below):

```go
func (b *B) Div(args ...any)   // normal (container) element
func (b *B) Img(args ...any)   // void element: attrs only, body ignored/omitted
```

`args` are interpreted by type:

- `Attributes`, `Attribute`, `AttrBuilder` → merged into the element's
  attributes (later values override earlier; same merge semantics as today's
  `parseTagArgs`).
- `func(b *B)` (and the alias `Body`) → the element's body closure, run between
  the open and close tags. At most one; last wins. Void elements have no body.
- `nil` → ignored (so `b.Div(cond && x)`-style and optional args are ergonomic).

`any` is used (rather than a sealed interface) specifically so a
`func(b *h.B){…}` literal can be passed directly, matching the target syntax.
Unknown arg types should `panic` with a clear message (a programming error
surfaced at dev time), consistent with how `Attrs` panics today.

Provide a `Body` alias for readability and for typed component helpers:

```go
type Body = func(*B)
```

### Void / self-closing elements

`Area, Base, Br, Col, Embed, Hr, Img, Input, Link, Meta, Source, Track, Wbr`
(the current `stag` set) render self-closing and accept attributes only.

### Root document

```go
func (b *B) Html(args ...any)   // writes <!DOCTYPE html> then <html lang="en"> ... </html>,
                                 // defaulting lang="en" if not provided (as today)
func (b *B) Doctype()
```

### Text / raw content (leaf statements)

```go
func (b *B) Text(s string)                    // HTML-escaped
func (b *B) Textf(format string, a ...any)    // HTML-escaped, fmt.Sprintf
func (b *B) Raw(s string)                      // unescaped (caller-sanitized)
func (b *B) Rawf(format string, a ...any)      // unescaped, fmt.Sprintf
```

### Custom elements

```go
func (b *B) El(name string, args ...any)       // arbitrary container tag
func (b *B) VoidEl(name string, args ...any)    // arbitrary void tag
```

(Replaces `CustomElement`.)

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

Keep the `sync.Pool` reuse of `*B`/buffers for allocation efficiency (as the
current writer pool does). `SetMaxLineLength`-style knobs may remain as internal
options configured by the Render* variants; no need to expose new public knobs
beyond what exists, but preserve indent + max-line-length behavior for tests.

## Components / composition

A component is just a function taking `*B` (plus its own params). Composition is
a plain call — no wrapper types:

```go
func Card(b *h.B, title string, body h.Body) {
    b.Div(h.Attrs("class", "card"), func(b *h.B) {
        b.H2(func(b *h.B) { b.Text(title) })
        body(b)
    })
}

h.Render(w, func(b *h.B) {
    Card(b, "Hello", func(b *h.B) {
        b.P(func(b *h.B) { b.Text("world") })
    })
})
```

Conditionals and iteration use native Go:

```go
b.Ul(func(b *h.B) {
    for _, it := range items {
        if it.Visible {
            b.Li(func(b *h.B) { b.Text(it.Name) })
        }
    }
})
```

## Attribute layer (unchanged, keep exactly)

`Attribute`, `Attributes` (with `Get/Index/Set/SetDefault/Delete/Merge`), `Attr`,
`Attrs`, `AttrsMap`, `AttrIf`, `AttrBuilder`. The only change: `AttrBuilder` and
`Attribute`/`Attributes` no longer implement an `isTagArg()` marker (that
interface is gone); element methods accept them via the `any` type switch. If a
marker method (`isTagArg`) is referenced anywhere, remove it. `ds`/`hx` builders
implement `AttrBuilder` (their `Attribute() h.Attribute` method) and must
continue to compile and be accepted as element args.

## Tests

Rewrite `h`'s test suite (`api_test.go`, `helpers_test.go`, `benchmark_test.go`)
to exercise the streaming API. Preserve coverage of: escaping, attribute
merging/ordering, void elements, doctype/html defaulting, indent + max-line
wrapping, sticky-error propagation (a failing io.Writer makes Render return the
error and stops output), and pooling. Update `ds`/`hx`/`js` `doc.go` example
comments to the new API. Keep `go build ./... && go vet ./... && go test ./...`
green.
