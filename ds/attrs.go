package ds

import (
	"strings"

	"github.com/jeffh/htmlgen/h"
	"github.com/jeffh/htmlgen/js"
)

// SetSignalExpr returns a Value that assigns a signal to a JavaScript expression.
// The signal name is automatically prefixed with "$".
func SetSignalExpr(signalName string, expression js.Expr) Value {
	var sb strings.Builder
	sb.Grow(len(signalName) + 10)
	if strings.HasPrefix(signalName, "$") {
		sb.WriteString(signalName)
	} else {
		sb.WriteString("$")
		sb.WriteString(signalName)
	}
	sb.WriteString(" = ")
	sb.WriteString(js.ToJS(expression))
	return Value{expr: js.Raw(sb.String())}
}

// SetSignal returns a Value that assigns a signal to a value.
// Use SetSignalExpr if you need to set the signal to a complex expression.
// The signal name is automatically prefixed with "$".
func SetSignal(signalName string, jsValue any) Value {
	switch v := jsValue.(type) {
	case js.Expr:
		return SetSignalExpr(signalName, v)
	case Value:
		return SetSignalExpr(signalName, v.expr)
	default:
		return SetSignalExpr(signalName, js.JSON(jsValue))
	}
}

// OnSubmit creates a data-on:submit event handler.
func OnSubmit(actions ...Value) *EventBuilder {
	return newEventBuilder("data-on:submit", actions)
}

// OnInput creates a data-on:input event handler.
func OnInput(actions ...Value) *EventBuilder {
	return newEventBuilder("data-on:input", actions)
}

// OnChange creates a data-on:change event handler.
func OnChange(actions ...Value) *EventBuilder {
	return newEventBuilder("data-on:change", actions)
}

// OnClick creates a data-on:click event handler.
func OnClick(actions ...Value) *EventBuilder {
	return newEventBuilder("data-on:click", actions)
}

// OnLoad creates a data-on:load event handler.
func OnLoad(actions ...Value) *EventBuilder {
	return newEventBuilder("data-on:load", actions)
}

// On creates a custom data-on:<eventName> event handler.
func On(eventName string, actions ...Value) *EventBuilder {
	return newEventBuilder("data-on:"+eventName, actions)
}

// OnIntersect creates a data-on-intersect handler triggered when the element
// intersects the viewport.
func OnIntersect(actions ...Value) *IntersectBuilder {
	return newIntersectBuilder(actions)
}

// OnInterval creates a data-on-interval handler executed at regular intervals.
func OnInterval(actions ...Value) *IntervalBuilder {
	return newIntervalBuilder(actions)
}

// OnSignalPatch creates a data-on-signal-patch handler that runs whenever a
// signal is updated.
func OnSignalPatch(actions ...Value) *SignalPatchBuilder {
	return newSignalPatchBuilder(actions)
}

// OnSignalPatchFilter sets the filter for signals that trigger
// OnSignalPatch handlers.
func OnSignalPatchFilter(options *FilterOptions) h.Attribute {
	if options == nil {
		return h.Attr("data-on-signal-patch-filter", "")
	}
	var sb strings.Builder
	options.appendJS(&sb)
	return h.Attr("data-on-signal-patch-filter", sb.String())
}

// SignalExpr declares a signal initialized from an expression.
// The signal name is appended to "data-signals:". Modifiers (Case) may be chained.
func SignalExpr(name string, defaultExpression Value) *NamedBuilder {
	b := &NamedBuilder{attrBase: newAttr("data-signals:" + name)}
	b.addStmt(js.ToJS(defaultExpression.expr))
	return b
}

// Signal declares a signal with a JSON-encoded default value.
// The signal name is appended to "data-signals:". Modifiers (Case) may be chained.
func Signal(name string, defaultJsValue any) *NamedBuilder {
	b := &NamedBuilder{attrBase: newAttr("data-signals:" + name)}
	b.addStmt(js.ToJS(js.JSON(defaultJsValue)))
	return b
}

// Signals declares multiple signals using object syntax.
// Modifiers (Case, IfMissing, Terse) may be chained.
func Signals(signals map[string]any) *SignalsBuilder {
	b := &SignalsBuilder{attrBase: newAttr("data-signals")}
	b.addStmt(js.ToJS(js.JSON(signals)))
	return b
}

// Bind creates a two-way data binding for the named signal as the value form.
// Modifiers (Case) may be chained.
func Bind(signalName string) *NamedBuilder {
	b := &NamedBuilder{attrBase: newAttr("data-bind")}
	b.addStmt(signalName)
	return b
}

// BindKey creates a two-way data binding using key syntax (signal in attribute name).
func BindKey(signalName string) *NamedBuilder {
	return &NamedBuilder{attrBase: newAttr("data-bind:" + signalName)}
}

// Class binds a CSS class to a JavaScript expression.
//
//	ds.Class("hidden", ds.Raw("$collapsed"))  =>  data-class:hidden="$collapsed"
func Class(clsName string, value Value) h.Attribute {
	return h.Attr("data-class:"+clsName, js.ToJS(value.expr))
}

// Text binds the element's text content to a JavaScript expression.
func Text(value Value) h.Attribute {
	return h.Attr("data-text", js.ToJS(value.expr))
}

// Show conditionally shows or hides the element based on an expression.
func Show(value Value) h.Attribute {
	return h.Attr("data-show", js.ToJS(value.expr))
}

// Hide returns a style attribute that hides the element. For reactive hiding
// use Show with a negated condition instead.
func Hide() h.Attribute {
	return h.Attr("style", "display: none")
}

// Attribute reactively sets a single HTML attribute.
//
//	ds.Attribute("title", ds.Raw("$foo"))  =>  data-attr:title="$foo"
func Attribute(name string, value Value) h.Attribute {
	return h.Attr("data-attr:"+name, js.ToJS(value.expr))
}

// Indicator declares a fetch-indicator signal.
func Indicator(signalName string) h.Attribute {
	signalName = strings.TrimLeft(signalName, "$")
	return h.Attr("data-indicator", signalName)
}

// IndicatorKey creates a fetch-indicator signal using key syntax.
func IndicatorKey(signalName string) *NamedBuilder {
	signalName = strings.TrimLeft(signalName, "$")
	return &NamedBuilder{attrBase: newAttr("data-indicator:" + signalName)}
}

// Ignore marks an element to be ignored by Datastar.
func Ignore() h.Attribute {
	return h.Attr("data-ignore", "")
}

// IgnoreSelf ignores only the element itself, not its descendants.
func IgnoreSelf() h.Attribute {
	return h.Attr("data-ignore__self", "")
}

// IgnoreMorph prevents the element from being morphed during backend patches.
func IgnoreMorph() h.Attribute {
	return h.Attr("data-ignore-morph", "")
}

// PreserveAttr preserves specified attribute values during DOM morphing.
func PreserveAttr(attrs ...string) h.Attribute {
	return h.Attr("data-preserve-attr", strings.Join(attrs, " "))
}

// Effect runs an expression reactively whenever dependencies change.
// Multiple values are joined with "; ".
func Effect(values ...Value) h.Attribute {
	stmts := make([]string, len(values))
	for i, v := range values {
		stmts[i] = js.ToJS(v.expr)
	}
	return h.Attr("data-effect", strings.Join(stmts, "; "))
}

// Init runs an expression when the element loads into the DOM.
func Init(values ...Value) h.Attribute {
	stmts := make([]string, len(values))
	for i, v := range values {
		stmts[i] = js.ToJS(v.expr)
	}
	return h.Attr("data-init", strings.Join(stmts, "; "))
}

// Peek wraps a Value in @peek(() => expr) for debugging.
func Peek(action Value) Value {
	return V(ActionPeek(action.expr))
}

// Computed creates a read-only signal computed from an expression.
//
//	ds.Computed("total", ds.Raw("$price * $quantity"))
//	=> data-computed:total="$price * $quantity"
func Computed(name string, expression Value) *NamedBuilder {
	b := &NamedBuilder{attrBase: newAttr("data-computed:" + name)}
	b.addStmt(js.ToJS(expression.expr))
	return b
}

// Ref creates a signal referencing a DOM element.
//
//	ds.Ref("myElement")  =>  data-ref:myElement
func Ref(signalName string) *NamedBuilder {
	return &NamedBuilder{attrBase: newAttr("data-ref:" + signalName)}
}

// Style sets an inline CSS style property reactively.
//
//	ds.Style("background-color", ds.Raw("$isRed ? 'red' : 'blue'"))
//	=> data-style:background-color="$isRed ? 'red' : 'blue'"
func Style(property string, expression Value) h.Attribute {
	return h.Attr("data-style:"+property, js.ToJS(expression.expr))
}

// Styles sets multiple inline CSS styles reactively using object syntax.
func Styles(styles map[string]string) h.Attribute {
	return h.Attr("data-style", js.ToJS(js.JSON(styles)))
}

// Attrs sets multiple HTML attributes reactively using object syntax.
func Attrs(attrs map[string]string) h.Attribute {
	return h.Attr("data-attr", js.ToJS(js.JSON(attrs)))
}

// Classes sets multiple CSS classes conditionally using object syntax.
func Classes(classes map[string]string) h.Attribute {
	return h.Attr("data-class", js.ToJS(js.JSON(classes)))
}

// JsonSignalsDebug displays reactive JSON-stringified signals for debugging.
func JsonSignalsDebug(options *FilterOptions) *SignalsBuilder {
	b := &SignalsBuilder{attrBase: newAttr("data-json-signals")}
	if options != nil {
		var sb strings.Builder
		options.appendJS(&sb)
		b.addStmt(sb.String())
	}
	return b
}

// FilterOptions specifies regex patterns for filtering signals.
type FilterOptions struct {
	IncludeReg *string
	ExcludeReg *string
}

// appendJS writes the FilterOptions as a JavaScript object with regex literals.
// Output format: {include: /pattern/, exclude: /pattern/}
func (o *FilterOptions) appendJS(sb *strings.Builder) {
	sb.WriteString("{")
	needComma := false
	if o.IncludeReg != nil {
		sb.WriteString("include: /")
		sb.WriteString(*o.IncludeReg)
		sb.WriteString("/")
		needComma = true
	}
	if o.ExcludeReg != nil {
		if needComma {
			sb.WriteString(", ")
		}
		sb.WriteString("exclude: /")
		sb.WriteString(*o.ExcludeReg)
		sb.WriteString("/")
	}
	sb.WriteString("}")
}

// SetAll creates a @setAll(value, filter) Datastar action.
func SetAll(value Value, options *FilterOptions) Value {
	return V(ActionSetAll(value.expr, options))
}

// ToggleAll creates a @toggleAll(filter) Datastar action.
func ToggleAll(options *FilterOptions) Value {
	return V(ActionToggleAll(options))
}
