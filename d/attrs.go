package d

import (
	"strings"

	"github.com/jeffh/htmlgen/h"
)

// Sets a signal to an arbitrary JavaScript expression.
// The signalName will automatically be prefixed with "$".
func SetSignalExpr(signalName string, expression AttrValueAppender) AttrMutator {
	return AttrFunc(func(attr *attrBuilder) {
		sb := strings.Builder{}
		sb.Grow(len(signalName) + 10)
		if strings.HasPrefix(signalName, "$") {
			sb.WriteString(signalName)
		} else {
			sb.WriteString("$")
			sb.WriteString(signalName)
		}
		sb.WriteString(" = ")
		expression.Append(&sb)
		attr.AppendStatement(sb.String())
	})
}

// Sets a signal to a value that is encoded as JSON.
// Use SetSignalExpr if you need to set the signal to a more complex expression.
// The signalName will automatically be prefixed with "$".
func SetSignal(signalName string, jsValue any) AttrMutator {
	if expr, ok := jsValue.(AttrValueAppender); ok {
		return SetSignalExpr(signalName, expr)
	}
	return SetSignalExpr(signalName, JsonValue(jsValue))
}

// Sets an action to be executed when the form is submitted.
// The action will be encoded as a JavaScript expression.
func OnSubmit(options ...AttrMutator) h.Attribute { return exprAttr("data-on:submit", options...) }

// Sets an action to be executed when the input is changed.
// The action will be encoded as a JavaScript expression.
func OnInput(options ...AttrMutator) h.Attribute { return exprAttr("data-on:input", options...) }

// Sets an action to be executed when the change is detected.
// The action will be encoded as a JavaScript expression.
func OnChange(options ...AttrMutator) h.Attribute { return exprAttr("data-on:change", options...) }

// Sets an action to be executed when the element is clicked.
// The action will be encoded as a JavaScript expression.
func OnClick(options ...AttrMutator) h.Attribute { return exprAttr("data-on:click", options...) }

// Sets an action to be executed when the element is loaded.
// The action will be encoded as a JavaScript expression.
func OnLoad(options ...AttrMutator) h.Attribute { return exprAttr("data-on:load", options...) }

// Sets an action to be executed when the event is triggered.
// The action will be encoded as a JavaScript expression.
func On(eventName string, options ...AttrMutator) h.Attribute {
	opts := append([]AttrMutator{appendName(eventName)}, options...)
	return exprAttr("data-on:", opts...)
}

// OnIntersect runs an expression when the element intersects the viewport.
// Use Half() for 50% visibility, Full() for 100% visibility.
// Example: OnIntersect(Once(), Raw("$seen = true"))
// Produces: data-on-intersect__once="$seen = true"
func OnIntersect(options ...AttrMutator) h.Attribute {
	return exprAttr("data-on-intersect", options...)
}

// OnInterval executes an expression at regular intervals.
// Default interval is 1 second. Use Duration() to customize.
// Example: OnInterval(Duration(500*time.Millisecond), Raw("$count++"))
// Produces: data-on-interval__duration.500ms="$count++"
func OnInterval(options ...AttrMutator) h.Attribute {
	return exprAttr("data-on-interval", options...)
}

// OnSignalPatch runs an expression whenever signals are updated.
// A "patch" variable is available containing signal patch details.
// Example: OnSignalPatch(Raw("console.log('Signal changed!')"))
// Produces: data-on-signal-patch="console.log('Signal changed!')"
func OnSignalPatch(options ...AttrMutator) h.Attribute {
	return exprAttr("data-on-signal-patch", options...)
}

// OnSignalPatchFilter filters which signals trigger OnSignalPatch handlers.
// Example: OnSignalPatchFilter(&FilterOptions{IncludeReg: ptr("^counter$")})
// Produces: data-on-signal-patch-filter="{include: /^counter$/}"
func OnSignalPatchFilter(options *FilterOptions) h.Attribute {
	if options == nil {
		return h.Attr("data-on-signal-patch-filter", "")
	}
	return exprAttr("data-on-signal-patch-filter", FilterOptionsValue(options))
}

// Sets a signal to an arbitrary JavaScript expression.
// The signal's default value will be encoded as a JavaScript expression.
// The signal name will automatically be prefixed with "$".
func SignalExpr(name string, defaultExpression AttrValueAppender) h.Attribute {
	return exprAttr("data-signals:", appendName(name), AttrFunc(func(attr *attrBuilder) {
		defaultExpression.Append(&attr.name)
	}))
}

// Defines a signal with a default value.
// The signal's default value will be encoded as a JSON value.
// The signal name will automatically be prefixed with "$".
func Signal(name string, defaultJsValue any) h.Attribute {
	return exprAttr("data-signals:", appendName(name), JsonValue(defaultJsValue))
}

// Signals defines multiple signals with default values using object syntax.
// The signals will be encoded as a JSON object.
// Example: Signals(map[string]any{"foo": 1, "bar": "hello"})
// Produces: data-signals="{\"foo\":1,\"bar\":\"hello\"}"
func Signals(signals map[string]any) h.Attribute {
	return exprAttr("data-signals", JsonValue(signals))
}

// Sets a signal to be used as the value of the element. Updates to the element will be reflected in the signal.
// The signal name will automatically be prefixed with "$".
func Bind(signalName string) h.Attribute {
	return exprAttr("data-bind", Raw(signalName))
}

// Sets a class to be used as the value of the element. Updates to the class will be reflected in the element.
func Class(clsName string, value ...AttrMutator) h.Attribute {
	value = append(value, appendName(clsName))
	return exprAttr("data-class", value...)
}

// Sets the text of the element to be the value of the signal. Updates to the signal will be reflected in the element.
// The signal name will automatically be prefixed with "$".
func Text(value ...AttrMutator) h.Attribute {
	return exprAttr("data-text", value...)
}

func Show(value ...AttrMutator) h.Attribute {
	return exprAttr("data-show", value...)
}

// Hide returns a style attribute that hides the element.
// For reactive hiding, use Show() with a negated condition instead.
func Hide() h.Attribute {
	return h.Attr("style", "display: none")
}

// Attribute sets a single HTML attribute value reactively.
// Example: Attribute("title", Raw("$foo"))
// Produces: data-attr:title="$foo"
func Attribute(name string, value ...AttrMutator) h.Attribute {
	return exprAttr("data-attr:"+name, value...)
}

func Indicator(signalName string) h.Attribute {
	signalName = strings.TrimLeft(signalName, "$")
	return h.Attr("data-indicator", signalName)
}

func Ignore() h.Attribute {
	return h.Attr("data-ignore", "")
}

func Effect(values ...AttrMutator) h.Attribute {
	return exprAttr("data-effect", values...)
}

func Peek(action AttrValueAppender) ValueMutator {
	return ValueMutatorFunc(func(sb *strings.Builder) {
		sb.WriteString("@peek(() => ")
		action.Append(sb)
		sb.WriteString(")")
	})
}

// Computed creates a read-only signal computed from an expression.
// The signal auto-updates when dependencies change.
// Example: Computed("total", Raw("$price * $quantity"))
// Produces: data-computed:total="$price * $quantity"
func Computed(name string, expression AttrValueAppender) h.Attribute {
	return exprAttr("data-computed:", appendName(name), AttrFunc(func(attr *attrBuilder) {
		expression.Append(&attr.name)
	}))
}

// ComputedExpr creates a computed signal with modifiers.
// Example: ComputedExpr("total", Case(CamelCase), Raw("$price * $quantity"))
func ComputedExpr(name string, options ...AttrMutator) h.Attribute {
	opts := append([]AttrMutator{appendName(name)}, options...)
	return exprAttr("data-computed:", opts...)
}

// Init runs an expression when the element loads into the DOM.
// Example: Init(Raw("$count = 1"))
// Produces: data-init="$count = 1"
func Init(options ...AttrMutator) h.Attribute {
	return exprAttr("data-init", options...)
}

// Ref creates a signal referencing a DOM element.
// Example: Ref("myElement")
// Produces: data-ref:myElement
func Ref(signalName string, options ...AttrMutator) h.Attribute {
	opts := append([]AttrMutator{appendName(signalName)}, options...)
	return exprAttr("data-ref:", opts...)
}

// Style sets an inline CSS style property reactively.
// Example: Style("background-color", Raw("$isRed ? 'red' : 'blue'"))
// Produces: data-style:background-color="$isRed ? 'red' : 'blue'"
func Style(property string, expression ...AttrMutator) h.Attribute {
	opts := append([]AttrMutator{appendName(property)}, expression...)
	return exprAttr("data-style:", opts...)
}

// Styles sets multiple inline CSS styles reactively using object syntax.
// Example: Styles(map[string]string{"display": "$hidden ? 'none' : 'block'", "color": "$red ? 'red' : 'green'"})
func Styles(styles map[string]string) h.Attribute {
	return exprAttr("data-style", JsonValue(styles))
}

// Attrs sets multiple HTML attributes reactively using object syntax.
// Example: Attrs(map[string]string{"title": "$foo", "disabled": "$bar"})
// Produces: data-attr="{\"title\":\"$foo\",\"disabled\":\"$bar\"}"
func Attrs(attrs map[string]string) h.Attribute {
	return exprAttr("data-attr", JsonValue(attrs))
}

// Classes sets multiple CSS classes conditionally using object syntax.
// Example: Classes(map[string]string{"hidden": "$foo", "font-bold": "$bar"})
// Produces: data-class="{\"hidden\":\"$foo\",\"font-bold\":\"$bar\"}"
func Classes(classes map[string]string) h.Attribute {
	return exprAttr("data-class", JsonValue(classes))
}

// IgnoreMorph prevents the element from being morphed during backend patches.
// Produces: data-ignore-morph
func IgnoreMorph() h.Attribute {
	return h.Attr("data-ignore-morph", "")
}

// PreserveAttr preserves specified attribute values during DOM morphing.
// Example: PreserveAttr("open", "class")
// Produces: data-preserve-attr="open class"
func PreserveAttr(attrs ...string) h.Attribute {
	return h.Attr("data-preserve-attr", strings.Join(attrs, " "))
}

// JsonSignalsDebug displays reactive JSON stringified signals for debugging.
// Use with FilterOptions to include/exclude specific signals.
// Example: JsonSignalsDebug(nil) or JsonSignalsDebug(&FilterOptions{IncludeReg: ptr("user")})
// Produces: data-json-signals or data-json-signals="{include: /user/}"
func JsonSignalsDebug(options *FilterOptions, modifiers ...AttrMutator) h.Attribute {
	if options == nil {
		return exprAttr("data-json-signals", modifiers...)
	}
	opts := append([]AttrMutator{FilterOptionsValue(options)}, modifiers...)
	return exprAttr("data-json-signals", opts...)
}

// BindKey creates a two-way data binding using key syntax (signal name in key).
// Example: BindKey("foo", Case(CamelCase))
// Produces: data-bind:foo__case.camel
func BindKey(signalName string, options ...AttrMutator) h.Attribute {
	opts := append([]AttrMutator{appendName(signalName)}, options...)
	return exprAttr("data-bind:", opts...)
}

// IndicatorKey creates a fetch indicator signal using key syntax.
// Example: IndicatorKey("fetching", Case(CamelCase))
// Produces: data-indicator:fetching__case.camel
func IndicatorKey(signalName string, options ...AttrMutator) h.Attribute {
	signalName = strings.TrimLeft(signalName, "$")
	opts := append([]AttrMutator{appendName(signalName)}, options...)
	return exprAttr("data-indicator:", opts...)
}

// IgnoreSelf ignores only the element itself, not its descendants.
// Produces: data-ignore__self
func IgnoreSelf() h.Attribute {
	return h.Attr("data-ignore__self", "")
}

type FilterOptions struct {
	IncludeReg *string
	ExcludeReg *string
}

func (o *FilterOptions) append(sb *strings.Builder, prefix string) {
	if o.IncludeReg == nil && o.ExcludeReg == nil {
		return
	}
	sb.WriteString(prefix)
	o.appendJS(sb)
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

// FilterOptionsValue returns a ValueMutator that outputs FilterOptions as a JS object with regex literals.
func FilterOptionsValue(o *FilterOptions) ValueMutator {
	return ValueMutator{
		AttrMutator: AttrFunc(func(attr *attrBuilder) {
			var sb strings.Builder
			o.appendJS(&sb)
			attr.AppendStatement(sb.String())
		}),
		AttrValueAppender: AttrValueFunc(func(sb *strings.Builder) {
			o.appendJS(sb)
		}),
	}
}

func SetAll(value ValueMutator, options FilterOptions) ValueMutator {
	return ValueMutatorFunc(func(sb *strings.Builder) {
		sb.Grow(len(sb.String()) + 10)
		sb.WriteString("@setAll(")
		value.Append(sb)
		options.append(sb, ", ")
		sb.WriteString(")")
	})
}

func ToggleAll(options FilterOptions) ValueMutator {
	return ValueMutatorFunc(func(sb *strings.Builder) {
		sb.Grow(len(sb.String()) + 10)
		sb.WriteString("@toggleAll(")
		options.append(sb, "")
		sb.WriteString(")")
	})
}
