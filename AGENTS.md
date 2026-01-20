# AGENTS.md

This file provides guidance for AI coding agents working in the htmlgen Go library repository.

## Commands

### Testing
```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./h
go test ./js
go test ./d

# Run a specific test
go test -run TestHtmlWriting ./h
go test -run TestString ./js
go test -run TestRaw ./d

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
No specific linter is configured. Code should compile without errors.

## Architecture Overview

This is a Go library (`github.com/jeffh/htmlgen`) with three main packages:

- **`h`** - Core HTML generation with streaming Writer API and declarative Builder API
- **`js`** - Type-safe JavaScript code generation for event handlers
- **`d`** - Datastar reactive attribute helpers
- **`hx`** - HTMX attribute helpers (not documented in CLAUDE.md but present)

## Code Style Guidelines

### Go Version
- Requires Go 1.25.5 or later (uses `iter.Seq` from standard library)

### Imports
- Standard library imports first (separated by blank line)
- Third-party imports second
- Local package imports last
- Sort alphabetically within each group
- Example:
  ```go
  import (
      "io"
      "strings"
      
      "github.com/hashicorp/golang-lru/v2"
      
      "github.com/jeffh/htmlgen/h"
  )
  ```

### Naming Conventions
- **Exported types/functions**: PascalCase (e.g., `Builder`, `Render`, `Html`)
- **Unexported types/functions**: camelCase (e.g., `tagBuilder`, `parseTagArgs`)
- **Interfaces**: Typically use nouns (e.g., `Builder`, `Callable`) or descriptive names
- **Functions that create builders**: Match HTML element names (e.g., `Div`, `Span`, `A`)
- **Attribute helpers**: Descriptive names (e.g., `Attrs`, `AttrsMap`, `Attr`)
- **Test functions**: Start with `Test` (e.g., `TestHtmlWriting`, `TestString`)
- **Benchmark functions**: Start with `Benchmark` (e.g., `BenchmarkRender`)

### Type Definitions
- Use marker interfaces for type constraints (e.g., `TagArg`, `Builder`)
- Implement marker methods as no-ops: `func (x) isTagArg() {}`
- Prefer struct types over type aliases for builders
- Use interfaces to define behavior (e.g., `Builder`, `Stmt`, `Expr`, `Callable`)

### Error Handling
- Return `error` as the last return value
- Use sentinel errors for known error conditions: `var ErrUnknownTagToClose = errors.New(...)`
- Return early on errors using `if err != nil { return err }`
- Check `nil` builders before rendering: `if b == nil { return nil }`
- Wrap errors with context using `fmt.Errorf("%w: %s", ErrSentinel, name)`
- In `defer` with `recover()`, re-panic with context: `panic(fmt.Sprintf("attr %q: %s", name, r))`

### Nil Handling
- Nil builders are explicitly allowed and skipped during rendering
- Check for nil before processing: `if child != nil { ... }`
- Return nil from helpers when appropriate (e.g., `When()`, `First()`)

### Comments and Documentation
- **Package comments**: Start with "Package <name>" and describe purpose (in doc.go or first file)
- **Exported functions**: Start with function name, describe what it does
- **Include examples** in comments using correct formatting
- Document **important implementation details** (e.g., escaping, pooling)
- Mark deprecated items: `// Deprecated: Use XYZ instead.`
- Example:
  ```go
  // Render writes the HTML representation of the given Builder to w.
  // Returns nil if b is nil.
  func Render(w io.Writer, b Builder) error
  ```

### Function Structure
- **Variadic parameters**: Use `...TagArg` or `...Stmt` for flexibility
- **Options pattern**: Use interfaces like `AttrMutator` with `Modify()` method
- **Builder pattern**: Chain method calls where appropriate
- **Function order**: Exported functions first, then unexported helpers
- Keep functions focused and single-purpose

### Performance Patterns
- **Object pooling**: Use `sync.Pool` for frequently allocated objects (e.g., `Writer`, `strings.Builder`)
- **Pre-allocate slices**: Use `make([]T, 0, capacity)` when size is known
- **String building**: Use `strings.Builder` for concatenation, call `Grow()` when size is known
- **Minimize allocations**: Reuse buffers, avoid unnecessary copying
- Example pooling pattern:
  ```go
  var writerPool = sync.Pool{
      New: func() any { return &Writer{openTags: make([]string, 0, 32)} },
  }
  ```

### Testing Conventions
- **Table-driven tests**: Use slice of structs with `Desc/name`, `Expected`, and inputs
- **Helper functions**: Create test helpers like `exprString()`, `stmtString()` for common operations
- **Subtests**: Use `t.Run()` for multiple test cases
- **Error messages**: Include input, got, and want values
- Example:
  ```go
  tests := []struct {
      name     string
      input    string
      expected string
  }{
      {"simple string", "hello", `"hello"`},
  }
  for _, tt := range tests {
      t.Run(tt.name, func(t *testing.T) {
          got := exprString(String(tt.input))
          if got != tt.expected {
              t.Errorf("String(%q) = %q, want %q", tt.input, got, tt.expected)
          }
      })
  }
  ```

### Code Organization
- **One type per file** for major types (e.g., `writer.go`, `builder.go`, `attrs.go`)
- **Group related functionality**: `tags.go` for all HTML element functions
- **Separate tests**: `*_test.go` files in the same package
- **Package documentation**: Use `doc.go` for package-level documentation
- **Internal helpers**: Keep unexported in the same file as related exported functions

### Special Patterns

#### Escaping and Security
- **HTML escaping**: Always escape attribute values and text content via `writeEscapedString()`
- **Raw content**: Only expose via explicit `Raw()` functions with security warnings
- **JavaScript escaping**: Use `json.Marshal()` or manual escaping for JS strings
- Document security implications in function comments

#### Builder Interface
- All builders implement: `Build(w *Writer) error`
- Tag arguments implement marker interface: `isTagArg()`
- Support nil builders by checking before calling `Build()`

#### Attributes API
- Attributes are `[]Attribute` with helper methods
- Merge later values over earlier (e.g., in `Merge()`, last value wins)
- Sort map keys for deterministic output
- Support both `Attrs("k", "v", ...)` and `AttrsMap(map[string]string)`

### Common Gotchas
- **Panic on empty attribute names**: `Attr()` and `Attrs()` panic if name is empty
- **Indentation tracking**: Writer tracks `atLineStart` to properly indent content
- **Tag stack**: Writer maintains `openTags` slice for proper closing
- **Line wrapping**: Writer can wrap attributes based on `maxLineLen` setting
- **Iterator consumption**: `ForEach()` builders consume iterators during `Build()`, not creation

### Dependencies
- Minimal external dependencies (only `github.com/jeffh/gocheck` for testing)
- Use standard library where possible (`html/template` for escaping, `encoding/json` for JSON)
- No linter or formatter dependencies (use `gofmt` by default)

## Best Practices

1. **Always run tests** before committing changes
2. **Add tests for new functionality** using table-driven approach
3. **Document exported functions** with clear godoc comments
4. **Preserve backward compatibility** when modifying exported APIs
5. **Use type safety** - leverage interfaces and marker types
6. **Handle errors properly** - return early, provide context
7. **Consider performance** - use pooling for hot paths
8. **Security first** - escape by default, document raw/unsafe functions
9. **Keep it simple** - favor clarity over cleverness
10. **Follow Go conventions** - run `gofmt`, use standard patterns
