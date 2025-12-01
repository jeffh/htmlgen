package d

import (
	"encoding/json"
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

// Defines a set of signals with default values.
// The signals will be encoded as a JSON object.
// The signals name will automatically be prefixed with "$".
func JsonSignals(signals map[string]any) h.Attribute {
	return exprAttr("data-json-signals", JsonValue(signals))
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

func Hide(value ...AttrMutator) h.Attribute {
	if len(value) == 0 {
		return h.Attr("style", "display: none")
	}
	return exprAttr("data-hide", value...)
}

func Attribute(name string, value ...AttrMutator) h.Attribute {
	return exprAttr("data-attr-"+name, value...)
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

type FilterOptions struct {
	IncludeReg *string
	ExcludeReg *string
}

func (o *FilterOptions) append(sb *strings.Builder, prefix string) {
	if o.IncludeReg == nil && o.ExcludeReg == nil {
		return
	}
	sb.WriteString(prefix)
	err := json.NewEncoder(sb).Encode(o.toMap())
	if err != nil {
		panic(err)
	}
}
func (o *FilterOptions) toMap() map[string]any {
	return map[string]any{
		"include": o.IncludeReg,
		"exclude": o.ExcludeReg,
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
