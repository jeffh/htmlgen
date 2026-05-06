// Pro Datastar attributes require a commercial license from https://data-star.dev/.
// These attributes provide additional functionality beyond the free tier.
package ds

import (
	"strings"

	"github.com/jeffh/htmlgen/h"
	"github.com/jeffh/htmlgen/js"
)

// Animate enables reactive animations on element attributes.
// Requires Datastar Pro.
func Animate(values ...Value) h.Attribute {
	stmts := make([]string, len(values))
	for i, v := range values {
		stmts[i] = js.ToJS(v.expr)
	}
	return h.Attr("data-animate", strings.Join(stmts, "; "))
}

// CustomValidity adds a custom validation message expression to a form input.
// Empty strings indicate valid; non-empty strings are shown as validation errors.
// Requires Datastar Pro.
func CustomValidity(expression Value) h.Attribute {
	return h.Attr("data-custom-validity", js.ToJS(expression.expr))
}

// OnRAF executes an expression on every requestAnimationFrame.
// Returns an EventBuilder that supports the standard event modifiers.
// Requires Datastar Pro.
func OnRAF(actions ...Value) *EventBuilder {
	return newEventBuilder("data-on-raf", actions)
}

// OnResize runs an expression when the element's dimensions change.
// Returns an EventBuilder that supports the standard event modifiers.
// Requires Datastar Pro.
func OnResize(actions ...Value) *EventBuilder {
	return newEventBuilder("data-on-resize", actions)
}

// PersistBuilder builds the data-persist attribute.
type PersistBuilder struct {
	*attrBase
}

// Session appends "__session" — use session storage instead of local storage.
func (b *PersistBuilder) Session() *PersistBuilder {
	b.name.WriteString("__session")
	return b
}

// Persist persists signals in local storage. Use Session() for session storage.
// Requires Datastar Pro.
func Persist(options *FilterOptions) *PersistBuilder {
	b := &PersistBuilder{attrBase: newAttr("data-persist")}
	if options != nil && (options.IncludeReg != nil || options.ExcludeReg != nil) {
		var sb strings.Builder
		options.appendJS(&sb)
		b.addStmt(sb.String())
	}
	return b
}

// PersistKey persists signals using a custom storage key.
// Requires Datastar Pro.
func PersistKey(key string) *PersistBuilder {
	return &PersistBuilder{attrBase: newAttr("data-persist:" + key)}
}

// QueryStringBuilder builds the data-query-string attribute.
type QueryStringBuilder struct {
	*attrBase
}

// Filter appends "__filter" — filters out empty values.
func (b *QueryStringBuilder) Filter() *QueryStringBuilder {
	b.name.WriteString("__filter")
	return b
}

// History appends "__history" — enables browser history integration.
func (b *QueryStringBuilder) History() *QueryStringBuilder {
	b.name.WriteString("__history")
	return b
}

// QueryString syncs query parameters to/from signal values.
// Requires Datastar Pro.
func QueryString(options *FilterOptions) *QueryStringBuilder {
	b := &QueryStringBuilder{attrBase: newAttr("data-query-string")}
	if options != nil && (options.IncludeReg != nil || options.ExcludeReg != nil) {
		var sb strings.Builder
		options.appendJS(&sb)
		b.addStmt(sb.String())
	}
	return b
}

// ReplaceURL replaces the browser URL without page reload.
// Requires Datastar Pro.
func ReplaceURL(expression Value) h.Attribute {
	return h.Attr("data-replace-url", js.ToJS(expression.expr))
}

// ScrollBuilder builds the data-scroll-into-view attribute.
type ScrollBuilder struct {
	*attrBase
}

// Smooth appends "__smooth".
func (b *ScrollBuilder) Smooth() *ScrollBuilder { b.name.WriteString("__smooth"); return b }

// Instant appends "__instant".
func (b *ScrollBuilder) Instant() *ScrollBuilder { b.name.WriteString("__instant"); return b }

// Auto appends "__auto" — browser default scrolling behavior.
func (b *ScrollBuilder) Auto() *ScrollBuilder { b.name.WriteString("__auto"); return b }

// HStart aligns to start of horizontal viewport.
func (b *ScrollBuilder) HStart() *ScrollBuilder { b.name.WriteString("__hstart"); return b }

// HCenter aligns to center of horizontal viewport.
func (b *ScrollBuilder) HCenter() *ScrollBuilder { b.name.WriteString("__hcenter"); return b }

// HEnd aligns to end of horizontal viewport.
func (b *ScrollBuilder) HEnd() *ScrollBuilder { b.name.WriteString("__hend"); return b }

// HNearest aligns to nearest edge of horizontal viewport.
func (b *ScrollBuilder) HNearest() *ScrollBuilder { b.name.WriteString("__hnearest"); return b }

// VStart aligns to start of vertical viewport.
func (b *ScrollBuilder) VStart() *ScrollBuilder { b.name.WriteString("__vstart"); return b }

// VCenter aligns to center of vertical viewport.
func (b *ScrollBuilder) VCenter() *ScrollBuilder { b.name.WriteString("__vcenter"); return b }

// VEnd aligns to end of vertical viewport.
func (b *ScrollBuilder) VEnd() *ScrollBuilder { b.name.WriteString("__vend")
	return b
}

// VNearest aligns to nearest edge of vertical viewport.
func (b *ScrollBuilder) VNearest() *ScrollBuilder {
	b.name.WriteString("__vnearest")
	return b
}

// Focus focuses the element after scrolling into view.
func (b *ScrollBuilder) Focus() *ScrollBuilder { b.name.WriteString("__focus"); return b }

// ScrollIntoView scrolls the element into the viewport.
// Requires Datastar Pro.
func ScrollIntoView() *ScrollBuilder {
	return &ScrollBuilder{attrBase: newAttr("data-scroll-into-view")}
}

// ViewTransitionName sets an explicit view-transition-name.
// Requires Datastar Pro.
func ViewTransitionName(expression Value) h.Attribute {
	return h.Attr("data-view-transition", js.ToJS(expression.expr))
}

// Pro Actions

// Clipboard copies text to the clipboard via @clipboard(text).
// Requires Datastar Pro.
func Clipboard(text Value) Value {
	return V(ActionClipboard(text.expr))
}

// ClipboardBase64 copies Base64-decoded text to the clipboard.
// Requires Datastar Pro.
func ClipboardBase64(text Value) Value {
	return V(ActionClipboardBase64(text.expr))
}

// Fit linearly interpolates a value from one range to another.
// Requires Datastar Pro.
func Fit(v, oldMin, oldMax, newMin, newMax Value) Value {
	return V(ActionFit(v.expr, oldMin.expr, oldMax.expr, newMin.expr, newMax.expr))
}

// FitClamped is like Fit but clamps the result to the target range.
// Requires Datastar Pro.
func FitClamped(v, oldMin, oldMax, newMin, newMax Value) Value {
	return V(ActionFitClamped(v.expr, oldMin.expr, oldMax.expr, newMin.expr, newMax.expr))
}

// FitRounded is like Fit but rounds the result to the nearest integer.
// Requires Datastar Pro.
func FitRounded(v, oldMin, oldMax, newMin, newMax Value) Value {
	return V(ActionFitRounded(v.expr, oldMin.expr, oldMax.expr, newMin.expr, newMax.expr))
}

// FitClampedRounded is like Fit with both clamping and rounding.
// Requires Datastar Pro.
func FitClampedRounded(v, oldMin, oldMax, newMin, newMax Value) Value {
	return V(ActionFitClampedRounded(v.expr, oldMin.expr, oldMax.expr, newMin.expr, newMax.expr))
}
