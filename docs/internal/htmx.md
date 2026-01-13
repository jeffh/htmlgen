# HTMX Reference

> A lightweight JavaScript library that extends HTML with custom attributes to enable AJAX, CSS transitions, WebSockets, and Server-Sent Events directly in markup.

**Version documented**: 2.0.8
**Last updated**: 2026-01-12
**Created with**: `/reference` command
**Updated**: Web Components section added

---

## Overview

HTMX is a dependency-free (~14kb min+gzip) JavaScript library that allows you to access modern browser features directly from HTML using attributes. Instead of writing JavaScript to make AJAX requests and manipulate the DOM, you declare behavior with HTML attributes like `hx-get`, `hx-post`, and `hx-swap`.

HTMX follows the HATEOAS (Hypermedia as the Engine of Application State) principle—servers return HTML fragments rather than JSON, keeping UI logic on the server. This approach often reduces code complexity significantly compared to JavaScript-heavy frameworks.

The library works with any backend technology and degrades gracefully when JavaScript is disabled (especially with `hx-boost`).

## Quick Start

```html
<!-- Include via CDN -->
<script src="https://cdn.jsdelivr.net/npm/htmx.org@2.0.8/dist/htmx.min.js"></script>
```

```html
<!-- Basic example: button that loads content via AJAX -->
<button hx-get="/api/data" hx-target="#result" hx-swap="innerHTML">
  Load Data
</button>
<div id="result"></div>
```

When clicked, this button issues a GET request to `/api/data` and replaces the inner HTML of `#result` with the response.

## Installation

### Requirements
- Any modern browser (IE11 compatible)
- No dependencies

### CDN
```html
<script src="https://cdn.jsdelivr.net/npm/htmx.org@2.0.8/dist/htmx.min.js"></script>
```

### NPM
```bash
npm install htmx.org@2.0.8
```

```javascript
// In your JavaScript entry point
import 'htmx.org';
```

### Direct Download
Download `htmx.min.js` from [htmx.org](https://htmx.org) and include in your project.

---

## Core Concepts

### AJAX Attributes

Any element can issue HTTP requests using these attributes:

| Attribute | HTTP Method |
|-----------|-------------|
| `hx-get` | GET |
| `hx-post` | POST |
| `hx-put` | PUT |
| `hx-patch` | PATCH |
| `hx-delete` | DELETE |

```html
<button hx-post="/items" hx-vals='{"name": "New Item"}'>Create</button>
```

### Targeting

`hx-target` specifies where the response content goes:

```html
<button hx-get="/content" hx-target="#container">Load</button>
<div id="container"><!-- Response goes here --></div>
```

**Target selectors:**
- `this` — The element itself
- `closest <selector>` — Nearest matching ancestor
- `next <selector>` — Next sibling matching selector
- `previous <selector>` — Previous sibling matching selector
- `find <selector>` — First descendant matching selector
- Any CSS selector — `#id`, `.class`, `element`

### Swapping

`hx-swap` controls how content is inserted:

| Value | Behavior |
|-------|----------|
| `innerHTML` | Replace inner content (default) |
| `outerHTML` | Replace entire element |
| `beforebegin` | Insert before element |
| `afterbegin` | Insert at start of element |
| `beforeend` | Insert at end of element |
| `afterend` | Insert after element |
| `delete` | Delete the target element |
| `none` | No swap (process headers only) |

**Swap modifiers:**
```html
<div hx-get="/data" hx-swap="innerHTML transition:true swap:500ms settle:100ms scroll:top">
```

- `transition:true` — Use View Transitions API
- `swap:<time>` — Delay before swap
- `settle:<time>` — Delay after swap for settling
- `scroll:top|bottom` — Scroll target after swap
- `show:top|bottom` — Scroll element into view
- `ignoreTitle:true` — Don't update page title

### Triggers

`hx-trigger` specifies what event initiates the request:

**Default triggers:**
- `input`, `textarea`, `select` → `change`
- `form` → `submit`
- Everything else → `click`

**Custom triggers:**
```html
<!-- Trigger on keyup with 500ms debounce -->
<input hx-get="/search" hx-trigger="keyup changed delay:500ms" hx-target="#results">

<!-- Trigger on page load -->
<div hx-get="/initial" hx-trigger="load"></div>

<!-- Poll every 2 seconds -->
<div hx-get="/status" hx-trigger="every 2s"></div>

<!-- Trigger when scrolled into view -->
<div hx-get="/lazy" hx-trigger="revealed"></div>
```

**Trigger modifiers:**
- `once` — Execute only once
- `changed` — Only if value changed
- `delay:<interval>` — Wait before executing (debounce)
- `throttle:<interval>` — Rate-limit execution
- `from:<selector>` — Listen on different element
- `target:<selector>` — Filter by event target
- `consume` — Consume the event (stopPropagation)
- `queue:first|last|all|none` — Queue behavior

---

## Usage

### Active Search

```html
<input type="search"
       hx-get="/search"
       hx-trigger="keyup changed delay:300ms"
       hx-target="#search-results"
       name="q">
<div id="search-results"></div>
```

### Form Submission

```html
<form hx-post="/submit" hx-target="#response" hx-swap="outerHTML">
  <input type="text" name="username" required>
  <input type="password" name="password" required>
  <button type="submit">Login</button>
</form>
<div id="response"></div>
```

### Infinite Scroll

```html
<div id="items">
  <!-- Existing items -->
  <div hx-get="/items?page=2"
       hx-trigger="revealed"
       hx-swap="afterend">
    Loading more...
  </div>
</div>
```

### Click to Edit

```html
<div hx-get="/contact/1/edit" hx-trigger="click" hx-swap="outerHTML">
  <p>John Doe</p>
  <p>john@example.com</p>
</div>
```

### Bulk Actions with Delete

```html
<button hx-delete="/items/5"
        hx-target="closest tr"
        hx-swap="outerHTML swap:500ms"
        hx-confirm="Are you sure?">
  Delete
</button>
```

### Loading Indicators

```html
<style>
  .htmx-indicator { opacity: 0; transition: opacity 200ms; }
  .htmx-request .htmx-indicator { opacity: 1; }
</style>

<button hx-get="/slow-endpoint" hx-target="#result">
  Load <span class="htmx-indicator">Loading...</span>
</button>
```

Or use `hx-indicator` to point to a separate element:
```html
<button hx-get="/slow" hx-indicator="#spinner">Load</button>
<span id="spinner" class="htmx-indicator">⏳</span>
```

### Disabling Elements During Request

```html
<button hx-post="/submit" hx-disabled-elt="this">
  Submit
</button>

<!-- Disable multiple elements -->
<form hx-post="/process" hx-disabled-elt="find input, find button">
  ...
</form>
```

---

## Attributes Reference

### Request Attributes

| Attribute | Description |
|-----------|-------------|
| `hx-get` | Issue GET request to URL |
| `hx-post` | Issue POST request to URL |
| `hx-put` | Issue PUT request to URL |
| `hx-patch` | Issue PATCH request to URL |
| `hx-delete` | Issue DELETE request to URL |

### Behavior Attributes

| Attribute | Description |
|-----------|-------------|
| `hx-trigger` | Event that triggers request |
| `hx-target` | Element to swap content into |
| `hx-swap` | How to swap content |
| `hx-swap-oob` | Mark content for out-of-band swap |
| `hx-select` | Select subset of response to swap |
| `hx-select-oob` | Select OOB content from response |
| `hx-boost` | Progressive enhancement for links/forms |
| `hx-push-url` | Push URL to browser history |
| `hx-replace-url` | Replace URL in browser history |

### Data Attributes

| Attribute | Description |
|-----------|-------------|
| `hx-include` | Include additional elements in request |
| `hx-params` | Filter request parameters |
| `hx-vals` | Add JSON values to request |
| `hx-headers` | Add custom headers |
| `hx-encoding` | Set encoding (e.g., `multipart/form-data`) |

### UX Attributes

| Attribute | Description |
|-----------|-------------|
| `hx-confirm` | Show confirmation dialog |
| `hx-prompt` | Show prompt dialog |
| `hx-indicator` | Element to show during request |
| `hx-disabled-elt` | Elements to disable during request |
| `hx-validate` | Force HTML5 validation |
| `hx-sync` | Coordinate requests between elements |

### Inheritance Attributes

| Attribute | Description |
|-----------|-------------|
| `hx-inherit` | Inherit attributes from ancestors |
| `hx-disinherit` | Disable specific attribute inheritance |
| `hx-ext` | Enable extensions |
| `hx-disable` | Disable htmx processing |

### History Attributes

| Attribute | Description |
|-----------|-------------|
| `hx-history` | Exclude from history cache |
| `hx-history-elt` | Element for history snapshots |

---

## Request Headers

HTMX sends these headers with requests:

| Header | Value |
|--------|-------|
| `HX-Request` | `true` |
| `HX-Current-URL` | Current browser URL |
| `HX-Target` | ID of target element |
| `HX-Trigger` | ID of triggering element |
| `HX-Trigger-Name` | Name of triggering element |
| `HX-Boosted` | `true` if boosted request |
| `HX-Prompt` | User's prompt response |

## Response Headers

Servers can control HTMX behavior with these headers:

| Header | Effect |
|--------|--------|
| `HX-Location` | Client-side redirect (no page reload) |
| `HX-Push-Url` | Push URL to history |
| `HX-Redirect` | Full page redirect |
| `HX-Refresh` | Full page refresh |
| `HX-Replace-Url` | Replace current URL in history |
| `HX-Reswap` | Override `hx-swap` value |
| `HX-Retarget` | Override `hx-target` selector |
| `HX-Reselect` | Override `hx-select` selector |
| `HX-Trigger` | Trigger client-side events |
| `HX-Trigger-After-Settle` | Trigger events after settle |
| `HX-Trigger-After-Swap` | Trigger events after swap |

---

## Out-of-Band Swaps

Update multiple elements from a single response:

```html
<!-- Server response -->
<div id="main-content">Main content here</div>
<div id="notification" hx-swap-oob="true">You have 3 new messages</div>
<div id="sidebar" hx-swap-oob="innerHTML">Updated sidebar</div>
```

Elements with `hx-swap-oob="true"` are swapped into their matching ID targets regardless of `hx-target`.

---

## Events Reference

### Request Lifecycle

| Event | When | Cancellable |
|-------|------|-------------|
| `htmx:configRequest` | Before request configured | No |
| `htmx:beforeRequest` | Before request sent | Yes |
| `htmx:beforeSend` | Immediately before transmission | No |
| `htmx:afterRequest` | After request completes | No |
| `htmx:afterOnLoad` | After onload, before swap | No |

### Swap Lifecycle

| Event | When | Cancellable |
|-------|------|-------------|
| `htmx:beforeSwap` | Before content swap | Yes |
| `htmx:afterSwap` | After content swap | No |
| `htmx:afterSettle` | After DOM settling | No |
| `htmx:load` | When new content loaded | No |

### Error Events

| Event | When |
|-------|------|
| `htmx:responseError` | HTTP error response |
| `htmx:sendError` | Network error |
| `htmx:timeout` | Request timeout |
| `htmx:swapError` | Error during swap |
| `htmx:targetError` | Invalid target selector |

### Validation Events

| Event | When |
|-------|------|
| `htmx:validation:validate` | Before validation |
| `htmx:validation:failed` | Validation failed |
| `htmx:validation:halted` | Request halted due to validation |

### Listening to Events

```html
<!-- Inline with hx-on -->
<button hx-post="/api" hx-on:htmx:before-request="console.log('Starting...')">
  Submit
</button>

<!-- Modify request parameters -->
<div hx-on:htmx:config-request="event.detail.parameters.extra = 'value'">
```

```javascript
// JavaScript
document.body.addEventListener('htmx:afterSwap', function(event) {
  console.log('Swapped:', event.detail.target);
});

// Or use htmx.on()
htmx.on('htmx:afterSwap', function(event) {
  console.log('Swapped:', event.detail.target);
});
```

---

## JavaScript API

### Core Functions

```javascript
// Execute AJAX request programmatically
htmx.ajax('GET', '/api/data', '#target');
htmx.ajax('POST', '/api/submit', {
  target: '#result',
  swap: 'innerHTML',
  values: { name: 'John' }
});

// Process dynamically added content
htmx.process(document.getElementById('new-content'));

// Find elements
htmx.find('#element');
htmx.findAll('.items');
htmx.closest(element, '.parent');

// DOM manipulation
htmx.addClass(element, 'active', 500); // with 500ms delay
htmx.removeClass(element, 'active');
htmx.toggleClass(element, 'visible');
htmx.remove(element, 1000); // remove after 1 second

// Custom swap
htmx.swap('#target', '<div>New content</div>', { swapStyle: 'innerHTML' });

// Get form values
const values = htmx.values(formElement);
```

### Event Handling

```javascript
// Register load handler for new content
htmx.onLoad(function(content) {
  // Initialize third-party libraries on new content
  initializeTooltips(content);
});

// Custom event listeners
htmx.on('#element', 'htmx:afterSwap', function(event) {
  console.log('Content swapped');
});

htmx.off('#element', 'htmx:afterSwap', handler);

// Trigger custom events
htmx.trigger('#element', 'custom-event', { data: 'value' });
```

### Debugging

```javascript
// Log all htmx events
htmx.logAll();

// Custom logger
htmx.logger = function(elt, event, data) {
  console.log('[HTMX]', event, elt, data);
};

// Disable logging
htmx.logNone();
```

---

## Configuration

### Via Meta Tag

```html
<meta name="htmx-config" content='{
  "defaultSwapStyle": "outerHTML",
  "historyCacheSize": 20,
  "timeout": 5000
}'>
```

### Via JavaScript

```javascript
htmx.config.defaultSwapStyle = 'outerHTML';
htmx.config.historyCacheSize = 20;
htmx.config.timeout = 5000;
```

### Key Options

| Option | Default | Description |
|--------|---------|-------------|
| `defaultSwapStyle` | `innerHTML` | Default swap strategy |
| `defaultSwapDelay` | `0` | Delay before swap (ms) |
| `defaultSettleDelay` | `20` | Delay after swap (ms) |
| `historyCacheSize` | `10` | History cache size |
| `historyEnabled` | `true` | Enable history support |
| `timeout` | `0` | Request timeout (ms, 0=none) |
| `withCredentials` | `false` | Include cookies in CORS |
| `scrollBehavior` | `instant` | Scroll animation |
| `globalViewTransitions` | `false` | Enable View Transitions API |
| `selfRequestsOnly` | `true` | Only allow same-origin requests |
| `allowEval` | `true` | Allow eval-based features |
| `allowScriptTags` | `true` | Process script tags in responses |
| `indicatorClass` | `htmx-indicator` | Loading indicator class |
| `requestClass` | `htmx-request` | Active request class |

---

## Extensions

Extensions add capabilities to HTMX. Enable with `hx-ext`:

```html
<body hx-ext="preload, response-targets">
```

### Core Extensions

| Extension | Purpose |
|-----------|---------|
| `head-support` | Merge head tags from responses |
| `idiomorph` | DOM morphing swap strategy |
| `preload` | Preload content on hover/focus |
| `response-targets` | Different targets by HTTP status |
| `sse` | Server-Sent Events |
| `ws` | WebSocket support |

### Installing Extensions

```html
<!-- Via CDN -->
<script src="https://cdn.jsdelivr.net/npm/htmx-ext-preload@2.1.0/preload.js"></script>

<!-- Enable on body or specific elements -->
<body hx-ext="preload">
```

### Example: Response Targets

```html
<body hx-ext="response-targets">
  <form hx-post="/submit"
        hx-target-200="#success"
        hx-target-422="#errors"
        hx-target-5xx="#server-error">
    ...
  </form>
</body>
```

### Example: Preload

```html
<body hx-ext="preload">
  <a href="/page" preload="mousedown">Hover to preload</a>
</body>
```

---

## Progressive Enhancement with hx-boost

Convert standard links and forms to AJAX:

```html
<body hx-boost="true">
  <!-- These work normally without JS, AJAX with JS -->
  <a href="/page">Link</a>
  <form action="/submit" method="post">...</form>
</body>
```

Boosted elements:
- Intercept navigation
- Load content via AJAX
- Update browser history
- Degrade gracefully without JavaScript

---

## Security

### Best Practices

1. **Escape all user content** — Server must escape HTML to prevent XSS
2. **Use CSRF tokens** — Include in headers or form values
3. **Validate URLs** — Use `htmx:validateUrl` event for custom validation

### CSRF Protection

```html
<!-- Via hx-headers -->
<body hx-headers='{"X-CSRF-Token": "token-value"}'>

<!-- Or via configRequest event -->
<script>
document.body.addEventListener('htmx:configRequest', function(event) {
  event.detail.headers['X-CSRF-Token'] = getCsrfToken();
});
</script>
```

### Security Configuration

```javascript
// Restrict to same-origin requests only (default)
htmx.config.selfRequestsOnly = true;

// Disable eval-based features for CSP compliance
htmx.config.allowEval = false;

// Disable script tag processing
htmx.config.allowScriptTags = false;
```

### Disable HTMX Processing

```html
<!-- Prevent htmx from processing user content -->
<div hx-disable>
  <!-- User-generated content here won't be processed by htmx -->
</div>
```

---

## Common Patterns

### Polling with Stop Condition

```html
<!-- Stop polling when response returns 286 status -->
<div hx-get="/status" hx-trigger="every 2s">
  Processing...
</div>
```

Server returns HTTP 286 to stop polling.

### Request Synchronization

```html
<!-- Abort previous requests when new one starts -->
<input hx-get="/search"
       hx-trigger="keyup changed delay:200ms"
       hx-sync="closest form:abort">
```

### Conditional Requests

```html
<!-- Only fire if condition is true -->
<button hx-get="/data" hx-trigger="click[ctrlKey]">
  Ctrl+Click to Load
</button>
```

### Custom Confirmation Dialog

```javascript
document.body.addEventListener('htmx:confirm', function(event) {
  event.preventDefault();
  showCustomDialog(event.detail.question).then(confirmed => {
    if (confirmed) event.detail.issueRequest();
  });
});
```

### Error Handling

```javascript
document.body.addEventListener('htmx:responseError', function(event) {
  const status = event.detail.xhr.status;
  if (status === 401) {
    window.location = '/login';
  } else if (status === 403) {
    showToast('Permission denied');
  }
});
```

---

## CSS Classes

HTMX applies these classes during operations:

| Class | When Applied |
|-------|--------------|
| `htmx-request` | During active request (on triggering element) |
| `htmx-indicator` | Always present on indicator elements |
| `htmx-added` | Briefly applied to new content before settling |
| `htmx-settling` | During settle phase |
| `htmx-swapping` | During swap phase |

### CSS Transitions

```html
<style>
  .item {
    opacity: 1;
    transition: opacity 500ms;
  }
  .item.htmx-swapping {
    opacity: 0;
  }
</style>

<div class="item" hx-get="/item/1" hx-swap="outerHTML swap:500ms">
  Content fades out before swap
</div>
```

---

## Web Components

HTMX can work with web components but requires manual setup for shadow DOM.

### Enabling HTMX in Shadow DOM

By default, HTMX doesn't see inside shadow DOM. Call `htmx.process()` on the shadow root:

```javascript
customElements.define('my-component', class MyComponent extends HTMLElement {
  connectedCallback() {
    const root = this.attachShadow({ mode: 'closed' })
    root.innerHTML = `
      <button hx-get="/endpoint" hx-target="next div">Click me!</button>
      <div></div>
    `
    htmx.process(root)  // Enable HTMX processing
  }
})
```

### Shadow DOM Selectors

Selectors are scoped to the current shadow DOM. Use these prefixes to escape:

| Prefix | Target |
|--------|--------|
| `host:` | The element hosting the shadow DOM |
| `global:` | Elements in the main document |

```html
<!-- Inside shadow DOM -->
<button hx-get="/data" hx-target="global:#main-content">Update Main</button>
```

### Limitations

- Form inputs inside shadow DOM aren't automatically included in requests
- Components must handle the `formdata` event or implement form element APIs
- Light DOM (no shadow DOM) works seamlessly with HTMX

---

## Debugging

### Enable Logging

```javascript
htmx.logAll();
```

### Browser Console

```javascript
// Monitor events on specific element
monitorEvents(document.getElementById('element'));
```

### Common Issues

**Request not firing:**
- Check for JavaScript errors in console
- Verify htmx is loaded
- Check trigger syntax

**Content not swapping:**
- Verify target selector exists
- Check server response content
- Look for `htmx:swapError` events

**History not working:**
- Ensure `hx-push-url` or `hx-boost` is set
- Check `htmx.config.historyEnabled`

---

## Sources

- [HTMX Official Documentation](https://htmx.org/docs/) - Retrieved 2026-01-12
- [HTMX Reference](https://htmx.org/reference/) - Retrieved 2026-01-12
- [HTMX Events Reference](https://htmx.org/events/) - Retrieved 2026-01-12
- [HTMX JavaScript API](https://htmx.org/api/) - Retrieved 2026-01-12
- [HTMX Extensions](https://htmx.org/extensions/) - Retrieved 2026-01-12
- [HTMX Web Components Example](https://htmx.org/examples/web-components/) - Retrieved 2026-01-12
- [HTMX GitHub Repository](https://github.com/bigskysoftware/htmx) - Retrieved 2026-01-12

---

*Reference created with `/reference` command*
