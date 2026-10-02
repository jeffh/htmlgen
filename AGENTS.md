# AGENTS.md

Guidance for AI coding agents working in `github.com/jeffh/htmlgen`.

**Source of truth for the HTML API:** [API_DESIGN.md](API_DESIGN.md).
Package-level detail lives in [CLAUDE.md](CLAUDE.md). Prefer those files over
duplicating long prose here.

## Commands

### Testing
```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./h
go test ./js
go test ./ds
go test ./ds/rkt
go test ./hx

# Run a specific test
go test -run TestStreamingElementsAndEscaping ./h
go test -run TestString ./js
go test -run TestRaw ./ds

# Run tests with coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...

# Run tests verbosely
go test -v ./...

# Run benchmarks
go test -bench=. ./h
go test -bench=. ./js
```

### Building
```bash
# Check compilation (no output if successful)
go build ./...

# Verify module dependencies
go mod tidy
go mod verify
```

### Linting
No specific linter is configured. Keep `go build ./... && go vet ./... && go test ./...` green.

## Architecture

Streaming HTML library. Package `h` has no tree of nodes: `Render` takes
`func(*B)` and streams immediately.

- **`h`** — `Render(w, func(*B))` creates a per-render `*h.B` bound to an
  `io.Writer`. Elements are methods on `*B`: containers are
  `Xxx(attrs Attributes, body Body)`, void elements are `Xxx(attrs Attributes)`.
  `nil` attrs and a `nil` body are valid. Control flow is native Go `if`/`for`.
- **`js`** — type-safe JavaScript expressions and statements for event handlers.
- **`ds`** — fluent Datastar attribute helpers (`OnClick()`, `Bind()`,
  `Signals()`, …).
- **`ds/rkt`** — Datastar Rocket helpers (private `$$` signals, template
  `For`/`If`, host props). Added in #63.
- **`hx`** — fluent HTMX attribute helpers.

`Attribute` and the fluent helpers in `ds`, `hx`, and `js` implement
`AttrBuilder`. Collect them with `h.AttrsOf(...)` or `Attributes.With(...)`.
Later same-name values override without changing position; zero attributes and
`nil` helpers are skipped; neither call mutates its inputs.

`B` buffers output (~4 KiB chunks), records a sticky first write error (later
output is a no-op; `Render` returns it), and is pooled via `builderPool`.
`Err()` is only current as of the last flush. Escaping uses a custom
`escapeTable` appended into the render buffer. `Text`/`Textf` escape;
`Raw`/`Rawf` write caller-sanitized content unchanged.

Do not invent APIs or refactor the HTML tag surface. See [API_DESIGN.md](API_DESIGN.md).

## Code style

### Go version
Requires Go 1.26.0 or later (`go.mod`).

### Imports
Standard library first, then third-party, then local. Alphabetical within each
group.

```go
import (
	"io"
	"strings"

	"github.com/jeffh/htmlgen/ds"
	"github.com/jeffh/htmlgen/h"
)
```

The only direct module dependency is `github.com/jeffh/gocheck` (tests).

### Naming
- **Exported**: PascalCase (`Render`, `AttrsOf`, `AttrBuilder`)
- **Unexported**: camelCase (`validName`, `writeAttrs`, `builderPool`)
- **Element methods on `*h.B`**: match HTML names (`Div`, `Span`, `A`)
- **Tests**: `Test...` / `Benchmark...`

### Types
- Element methods take concrete parameters, never `...any`
- `h.Body` is `func(*B)`
- Use interfaces for behavior (`AttrBuilder`, `Stmt`, `Expr`, `Callable`)

### Errors
- Return `error` last; return early on `if err != nil`
- `*h.B` uses a sticky write error — do not add per-call error returns on
  element methods
- Panic only for invalid attribute/tag names (written unescaped, so validated)
- In `defer` with `recover()`, re-panic with context

### Docs
Package comments start with `Package <name>`. Exported funcs start with the
function name. Document escaping, pooling, and other non-obvious behavior.
Mark deprecated items with `// Deprecated: Use XYZ instead.`

```go
// Render runs fn against a fresh *B writing to w and returns the first write
// error.
func Render(w io.Writer, fn func(*B)) error
```

### Functions
- Typed variadics (`...AttrBuilder` in `AttrsOf`, `...Stmt` in `js`), not `...any`
- Chain fluent methods where the existing helpers already do
- Exported functions first, then unexported helpers
- Keep functions focused

### Performance
- `sync.Pool` for hot objects (`builderPool`)
- Pre-allocate slices when the size is known
- `strings.Builder` with `Grow` when concatenating
- Minimize allocations; reuse buffers

### Tests
Table-driven (`name`/`Desc`, inputs, `expected`), `t.Run` subtests, helpers
such as `exprString()` / `stmtString()`. Error messages should include input,
got, and want.

### Layout
- `h/writer.go` (`*B` write/escape/indent), `h/render.go` (entry points + pool),
  `h/tags.go` (elements), `h/attrs.go`, `h/builder.go` (`Body` alias only)
- Tests alongside source (`*_test.go`); `doc.go` for package docs
- Unexported helpers stay in the same file as the related exported API

### Security
- Escape text and attribute values via `(*B).writeEscaped` / `appendEscaped`
  / `escapeTable`
- Expose unescaped output only through `Raw`/`Rawf`, with security warnings
- Attribute and custom tag names are written verbatim — `Attr`, `Attrs`,
  `AttrsMap`, `AttrIf`, `Set`, `SetDefault`, `El`, and `VoidEl` panic unless
  the name matches `[A-Za-z][A-Za-z0-9_.:-]*`. Constant-name tag methods skip
  the check. `Attribute` literals and `AttrBuilder` outputs are trusted; never
  derive their names from untrusted input
- JavaScript strings: `json.Marshal` or manual escaping

### Gotchas
- Invalid (including empty) names panic
- `B` tracks `atLineStart`, `openTags`, and optional `maxLineLen` wrapping
  when indenting
- Do not retain a `*B` after `Render` returns — it goes back to the pool
- Call `Flush()` when bytes must reach the client mid-render (SSE, long streams)

## Practices

1. Run tests before committing
2. Add table-driven tests for new behavior
3. Document exported APIs with godoc
4. Preserve backward compatibility
5. Escape by default; document raw/unsafe helpers
6. Prefer clarity over cleverness; run `gofmt`
