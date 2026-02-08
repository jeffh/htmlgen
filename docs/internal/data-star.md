# Datastar Reference Guide

This guide covers [Datastar](https://data-star.dev/), a lightweight hypermedia framework combining backend-driven reactivity with frontend interactivity. It extends HTML through `data-*` attributes without requiring npm or build tools.

## Table of Contents

- [Installation](#installation)
- [Core Concepts](#core-concepts)
- [Reactive Signals](#reactive-signals)
- [Datastar Expressions](#datastar-expressions)
- [Attributes Reference](#attributes-reference)
- [Actions Reference](#actions-reference)
- [Backend Requests](#backend-requests)
- [SSE Events](#sse-events)
- [Go SDK (datastar-go)](#go-sdk)
- [htmlgen/ds Package](#htmlgends-package)
- [Philosophy (The Tao of Datastar)](#philosophy)

## Installation

Include via CDN:

```html
<script type="module" src="https://cdn.jsdelivr.net/gh/starfederation/datastar@1.0.0-RC.7/bundles/datastar.js"></script>
```

Or install via package manager and self-host the bundle.

## Core Concepts

Datastar has two primary functions:

1. **Backend-Driven Updates**: Modify DOM and state through server-sent events
2. **Frontend Reactivity**: Build interactive UIs using `data-*` HTML attributes

### DOM Patching

The backend drives the frontend by patching HTML elements. Datastar uses morphing to update only modified DOM parts while preserving state.

**Response Types:**
- `text/html` - Elements morph into existing DOM based on IDs
- `text/event-stream` - SSE streaming for dynamic updates
- `application/json` - Signal patching

## Reactive Signals

Signals are reactive variables denoted with a `$` prefix that automatically track and propagate changes.

### Creating Signals

**Via `data-bind`** - Two-way binding on form elements:
```html
<input data-bind:username />
```
Creates `$username` linked to the input value.

**Via `data-signals`** - Explicit declaration:
```html
<div data-signals:count="0"></div>
<div data-signals="{foo: 1, bar: 'hello'}"></div>
```

**Nested signals** using dot notation:
```html
<div data-signals:form.email="''"></div>
```

### Signal Naming

- Hyphenated names auto-convert to camelCase: `data-bind:foo-bar` creates `$fooBar`
- Cannot begin with double underscores (`__`)
- Underscore-prefixed signals (`_local`) are excluded from backend requests by default

## Datastar Expressions

Expressions are evaluated strings within `data-*` attributes supporting JavaScript-like syntax.

### Available Variables

- `$signalName` - Reference signals with `$` prefix
- `el` - Current element reference (available in all expressions)
- `evt` - Event object (available in `data-on` handlers)
- `patch` - Signal patch details (available in `data-on-signal-patch`)

### Syntax

```html
<div data-text="$count"></div>
<div data-show="$items.length > 0"></div>
<button data-on:click="$count++; $total = $count * $price">Add</button>
```

Statements require semicolon separation (line breaks alone are insufficient).

## Attributes Reference

### Core Attributes

#### `data-attr`
Sets HTML attributes reactively.
```html
<div data-attr:aria-label="$foo"></div>
<div data-attr="{'title': $foo, 'disabled': $bar}"></div>
```

#### `data-bind`
Two-way data binding for form elements.
```html
<input data-bind:username />
<select data-bind:country>...</select>
```
**Modifiers:** `__case` (camel/kebab/snake/pascal)

#### `data-class`
Conditionally toggle CSS classes.
```html
<div data-class:font-bold="$isImportant"></div>
<div data-class="{active: $selected, hidden: $collapsed}"></div>
```

#### `data-computed`
Create read-only derived signals.
```html
<div data-computed:total="$price * $quantity"></div>
```

#### `data-effect`
Execute expressions on load and when dependencies change.
```html
<div data-effect="$total = $price * $quantity"></div>
```

#### `data-ignore`
Prevent Datastar from processing element and descendants.
```html
<div data-ignore>...</div>
```
**Modifiers:** `__self` (element only, not descendants)

#### `data-ignore-morph`
Skip DOM morphing for element and children.
```html
<div data-ignore-morph>Preserved content</div>
```

#### `data-indicator`
Track fetch request status (true during request).
```html
<button data-on:click="@get('/api')"
        data-indicator:loading
        data-attr:disabled="$loading">
  Submit
</button>
```

#### `data-init`
Run expressions when element loads into DOM.
```html
<div data-init="$count = 1"></div>
```
**Modifiers:** `__delay` (.500ms/.1s), `__viewtransition`

#### `data-json-signals`
Display JSON of signals for debugging.
```html
<pre data-json-signals></pre>
<pre data-json-signals="{include: /user/}"></pre>
```
**Modifiers:** `__terse` (compact format)

#### `data-on`
Attach event listeners.
```html
<button data-on:click="$count++">Click</button>
<form data-on:submit__prevent="@post('/submit')">...</form>
```
**Modifiers:**
- `__once` - Fire only once
- `__passive` - Passive event listener
- `__capture` - Capture phase
- `__prevent` - preventDefault()
- `__stop` - stopPropagation()
- `__debounce.Nms` - Debounce timing
- `__throttle.Nms` - Throttle timing
- `__delay.Nms` - Delay execution
- `__window` - Attach to window
- `__outside` - Trigger on outside events
- `__viewtransition` - Wrap in View Transition API

#### `data-on-intersect`
Execute when element enters/exits viewport.
```html
<div data-on-intersect="$seen = true"></div>
<div data-on-intersect__once__full="loadContent()"></div>
```
**Modifiers:** `__once`, `__exit`, `__half`, `__full`, `__threshold.N`

#### `data-on-interval`
Run expressions at regular intervals (default: 1s).
```html
<div data-on-interval="$seconds++"></div>
<div data-on-interval__duration.500ms="$tick++"></div>
```
**Modifiers:** `__duration.Nms` (.leading for immediate first fire)

#### `data-on-signal-patch`
Trigger when signals change.
```html
<div data-on-signal-patch="console.log('changed')"></div>
```
Use `data-on-signal-patch-filter="{include: /pattern/}"` to filter signals.

#### `data-preserve-attr`
Preserve attribute values during morphing.
```html
<details open data-preserve-attr="open">...</details>
```

#### `data-ref`
Create signal referencing the element.
```html
<input data-ref:myInput />
<!-- $myInput.value, $myInput.focus() -->
```

#### `data-show`
Conditionally show/hide elements.
```html
<div data-show="$isVisible">Content</div>
```

#### `data-signals`
Define signals with values.
```html
<div data-signals:count="0"></div>
<div data-signals="{user: {name: '', email: ''}}"></div>
```
**Modifiers:** `__ifmissing` (only set if absent), `__case`

#### `data-style`
Set inline CSS reactively.
```html
<div data-style:background-color="$color"></div>
<div data-style="{display: $hidden ? 'none' : 'block'}"></div>
```

#### `data-text`
Bind element text content.
```html
<span data-text="$message"></span>
<span data-text="`Count: ${$count}`"></span>
```

### Pro Attributes (Commercial License)

#### `data-animate`
Animate element attributes reactively.

#### `data-custom-validity`
Custom form validation. Empty = valid, non-empty = error message.
```html
<input data-bind:password
       data-custom-validity="$password.length >= 8 ? '' : 'Min 8 characters'" />
```

#### `data-on-raf`
Execute on every requestAnimationFrame.
```html
<div data-on-raf__throttle.16ms="updateAnimation()"></div>
```

#### `data-on-resize`
Trigger on element dimension changes.
```html
<div data-on-resize__debounce.200ms="$width = el.offsetWidth"></div>
```

#### `data-persist`
Persist signals to localStorage/sessionStorage.
```html
<div data-persist></div>
<div data-persist="{include: /settings/}"></div>
<div data-persist:mykey__session></div>
```

#### `data-query-string`
Sync query parameters with signals.
```html
<div data-query-string></div>
```
**Modifiers:** `__filter`, `__history`

#### `data-replace-url`
Update browser URL without reload.
```html
<div data-replace-url="`/items/${$page}`"></div>
```

#### `data-scroll-into-view`
Scroll element into viewport.
```html
<div data-scroll-into-view__smooth__vcenter></div>
```
**Modifiers:** `__smooth`, `__instant`, `__auto`, `__hstart`, `__hcenter`, `__hend`, `__hnearest`, `__vstart`, `__vcenter`, `__vend`, `__vnearest`, `__focus`

#### `data-view-transition`
Set view-transition-name for CSS animations.
```html
<div data-view-transition="$itemId"></div>
```

## Actions Reference

### Core Actions

#### `@peek(callable)`
Access signals without triggering re-evaluations.
```html
<div data-effect="console.log(@peek(() => $count))"></div>
```

#### `@setAll(value, filter?)`
Set all matching signals to a value.
```html
<button data-on:click="@setAll(false, {include: /selected/})">Clear All</button>
```

#### `@toggleAll(filter?)`
Toggle boolean state of matching signals.
```html
<button data-on:click="@toggleAll({include: /checkbox/})">Toggle All</button>
```

### Backend Actions

All accept a URI and optional configuration object.

| Action | HTTP Method |
|--------|------------|
| `@get(uri, options?)` | GET |
| `@post(uri, options?)` | POST |
| `@put(uri, options?)` | PUT |
| `@patch(uri, options?)` | PATCH |
| `@delete(uri, options?)` | DELETE |

**Options:**
- `contentType` - `'json'` (default) or `'form'`
- `filterSignals` - `{include: /regex/, exclude: /regex/}`
- `selector` - CSS selector for form elements
- `headers` - Custom HTTP headers object
- `openWhenHidden` - Keep connection alive when tab hidden
- `payload` - Override request body
- `retry` - `'auto'`, `'error'`, `'always'`, or `'never'`
- `retryInterval` - Wait between retries (ms)
- `retryScaler` - Exponential backoff multiplier
- `retryMaxWaitMs` - Maximum retry delay
- `retryMaxCount` - Maximum retry attempts
- `requestCancellation` - `'auto'`, `'disabled'`

Example:
```html
<button data-on:click="@post('/api/submit', {contentType: 'form'})">
  Submit
</button>
```

### Pro Actions (Commercial License)

#### `@clipboard(text, isBase64?)`
Copy text to clipboard.
```html
<button data-on:click="@clipboard($code)">Copy</button>
```

#### `@fit(v, oldMin, oldMax, newMin, newMax, clamp?, round?)`
Linear interpolation between ranges.
```html
<div data-computed:rgb="@fit($slider, 0, 100, 0, 255, true, true)"></div>
```

## Backend Requests

### Signal Transmission

By default, all signals (except underscore-prefixed locals) are sent with requests:
- **GET**: As `datastar` query parameter
- **Other methods**: As JSON body

### Response Handling

Responses are processed based on content type:

| Content-Type | Behavior |
|-------------|----------|
| `text/event-stream` | Process SSE events |
| `text/html` | Morph into DOM (use `datastar-selector` and `datastar-mode` headers) |
| `application/json` | Patch signals (RFC 7396 merge patch) |
| `text/javascript` | Execute client-side |

## SSE Events

Datastar processes `text/event-stream` responses with these event types:

### `datastar-patch-elements`
Modify DOM elements through morphing.

**Data fields:**
- `selector` - CSS selector (not required for `outer`/`replace` modes)
- `mode` - `outer` (default), `inner`, `replace`, `prepend`, `append`, `before`, `after`, `remove`
- `namespace` - `svg` or `mathml` for XML namespacing
- `useViewTransition` - Enable view transitions
- `elements` - HTML content

```
event: datastar-patch-elements
data: selector #container
data: mode inner
data: elements <div>Updated content</div>
```

### `datastar-patch-signals`
Update reactive signals.

**Data fields:**
- `signals` - Signal definitions (valid `data-signals` format)
- `onlyIfMissing` - Only update missing signals

```
event: datastar-patch-signals
data: signals {count: 42, message: "Hello"}
```

Set signal to `null` to remove it.

## Go SDK

Install the official Go SDK:

```bash
go get github.com/starfederation/datastar-go
```

Requires Go 1.24+.

### Reading Signals

```go
type Store struct {
    Count   int    `json:"count"`
    Message string `json:"message"`
}

func handler(w http.ResponseWriter, r *http.Request) {
    store := &Store{}
    if err := datastar.ReadSignals(r, store); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    // Use store.Count, store.Message
}
```

### SSE Responses

```go
func handler(w http.ResponseWriter, r *http.Request) {
    sse := datastar.NewSSE(w, r)

    // Update DOM
    sse.PatchElements(`<div id="output">Hello!</div>`)

    // Update signals
    sse.MarshalAndPatchSignals(map[string]any{
        "count": 42,
    })

    // Remove element
    sse.RemoveElement("#old-content")

    // Execute JavaScript
    sse.ExecuteScript("console.log('done')")

    // Redirect
    sse.Redirect("/new-page")
}
```

## htmlgen/ds Package

The `ds` package provides Go helpers for building Datastar attributes with htmlgen.

### Signals

```go
// Define signals
ds.Signal("count", 0)           // data-signals:count="0"
ds.Signals(map[string]any{      // data-signals="{...}"
    "foo": 1,
    "bar": "hello",
})
ds.Bind("username")             // data-bind="username"
ds.BindKey("email", ds.Case(ds.CamelCase)) // data-bind:email__case.camel

// Computed signals
ds.Computed("total", ds.Raw("$price * $qty"))
```

### Event Handlers

```go
ds.OnClick(ds.Get("/api"))                    // data-on:click="@get('/api')"
ds.OnSubmit(ds.PreventDefault(), ds.Post("/submit"))
ds.On("keydown", ds.Debounce(300*time.Millisecond), ds.Raw("search()"))
ds.OnIntersect(ds.Once(), ds.Raw("$seen = true"))
ds.OnInterval(ds.Duration(500*time.Millisecond), ds.Raw("$tick++"))
```

### Display Attributes

```go
ds.Show(ds.Raw("$isVisible"))    // data-show="$isVisible"
ds.Text(ds.Raw("$message"))      // data-text="$message"
ds.Class("active", ds.Raw("$selected"))
ds.Classes(map[string]string{"hidden": "$collapsed"})
ds.Style("color", ds.Raw("$textColor"))
ds.Attribute("title", ds.Raw("$tooltip"))
```

### Backend Actions

```go
// Simple requests
ds.Get("/api")
ds.Post("/submit")
ds.Put("/update")
ds.Delete("/remove")

// With options
ds.Get("/api", ds.RequestOptions(
    ds.ContentType("form"),
    ds.Headers(map[string]string{"X-Custom": "value"}),
    ds.FilterSignals(&ds.FilterOptions{IncludeReg: ptr("user")}),
    ds.RetryMaxCount(3),
))

// Dynamic paths
ds.GetDynamic(ds.Raw("`/api/items/${$id}`"))
```

### Modifiers

```go
ds.PreventDefault()
ds.StopPropagation()
ds.Once()
ds.Passive()
ds.Capture()
ds.Debounce(300 * time.Millisecond)
ds.Throttle(100 * time.Millisecond)
ds.Delay(500 * time.Millisecond)
ds.ViewTransition()
ds.Window()
ds.Outside()
ds.Case(ds.CamelCase)  // camel, kebab, snake, pascal
ds.IfMissing()
```

### DOM Control

```go
ds.Indicator("loading")         // data-indicator="loading"
ds.Ref("myElement")             // data-ref:myElement
ds.Init(ds.Raw("$count = 0"))
ds.Effect(ds.Raw("$total = $a + $b"))
ds.Ignore()
ds.IgnoreSelf()
ds.IgnoreMorph()
ds.PreserveAttr("open", "class")
ds.JsonSignalsDebug(nil)        // Debug display
```

### Values and Expressions

```go
ds.Raw("$foo + $bar")           // Raw JavaScript
ds.Str("hello")                 // JSON string: "hello"
ds.JsonValue(map[string]int{"a": 1})
ds.SetSignal("count", 42)       // $count = 42
ds.SetSignalExpr("msg", ds.Str("hi"))
ds.And(ds.Raw("$a"), ds.Raw("$b")) // $a && $b
ds.ConsoleLog(ds.Raw("$count"))
ds.Navigate("/page/%d", 5)      // window.location.href = "/page/5"
```

### Actions

```go
ds.Peek(ds.Raw("$count"))        // @peek(() => $count)
ds.SetAll(ds.Raw("false"), ds.FilterOptions{IncludeReg: ptr("selected")})
ds.ToggleAll(ds.FilterOptions{IncludeReg: ptr("checkbox")})
```

### Pro Features (Commercial License)

```go
// Attributes
ds.CustomValidity(ds.Raw("$pw === $confirm ? '' : 'Must match'"))
ds.Persist(nil)
ds.PersistKey("settings", ds.Session())
ds.QueryString(nil, ds.Filter(), ds.History())
ds.ReplaceURL(ds.Raw("`/page/${$num}`"))
ds.ScrollIntoView(ds.Smooth(), ds.VCenter())
ds.ViewTransitionName(ds.Raw("$itemId"))
ds.OnRAF(ds.Throttle(16*time.Millisecond), ds.Raw("animate()"))
ds.OnResize(ds.Debounce(200*time.Millisecond), ds.Raw("$w = el.offsetWidth"))

// Actions
ds.Clipboard(ds.Str("copied text"))
ds.ClipboardBase64(ds.Str("base64content"))
ds.Fit(ds.Raw("$v"), ds.Raw("0"), ds.Raw("100"), ds.Raw("0"), ds.Raw("255"))
ds.FitClamped(...)
ds.FitRounded(...)
ds.FitClampedRounded(...)
```

### Complete Example

```go
package main

import (
    "github.com/jeffh/htmlgen/h"
    "github.com/jeffh/htmlgen/ds"
    "time"
)

func CounterComponent() h.Builder {
    return h.Div(
        ds.Signal("count", 0),
        h.Button(
            h.Attrs("type", "button"),
            ds.OnClick(ds.Raw("$count++")),
            h.Text("Increment"),
        ),
        h.Span(
            ds.Text(ds.Raw("`Count: ${$count}`")),
        ),
    )
}

func SearchForm() h.Builder {
    return h.Form(
        ds.Signal("query", ""),
        ds.Indicator("searching"),
        h.Input(
            h.Attrs("type", "text", "placeholder", "Search..."),
            ds.Bind("query"),
            ds.OnInput(
                ds.Debounce(300*time.Millisecond),
                ds.Get("/search"),
            ),
        ),
        h.Div(
            h.Attrs("id", "results"),
            ds.Show(ds.Raw("!$searching")),
        ),
        h.Div(
            ds.Show(ds.Raw("$searching")),
            h.Text("Loading..."),
        ),
    )
}
```

## Philosophy

The Tao of Datastar emphasizes:

1. **Backend as Source of Truth**: Keep state server-side where it's protected from user interference.

2. **Use Defaults**: Start with recommended configuration and only deviate after careful consideration.

3. **Embrace Morphing**: Send HTML chunks from the backend; morphing updates only changed parts.

4. **Minimal Signals**: Use signals for UI interactions (toggles, form binding), not for caching state.

5. **CQRS Pattern**: Separate reads (long-lived SSE connections) from writes (short-lived requests).

6. **Show Loading States**: Use indicators during requests rather than optimistic updates.

7. **Semantic HTML**: Prioritize accessibility with proper ARIA attributes.

8. **Compression**: Use Brotli compression for efficient streaming (up to 200:1 ratios).

9. **Natural Navigation**: Let browsers handle history through anchor tags and standard page resources.
