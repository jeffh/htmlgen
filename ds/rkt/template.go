package rkt

import (
	"strings"

	"github.com/jeffh/htmlgen/ds"
	"github.com/jeffh/htmlgen/h"
	"github.com/jeffh/htmlgen/js"
)

// attrName accumulates the attribute name of a fluent Rocket builder.
// Modifier methods append to it; Attribute produces the finished attribute,
// whose name is validated by h.Attr.
type attrName struct {
	name  strings.Builder
	value string
}

func newAttrName(name string) *attrName {
	var a attrName
	a.name.WriteString(name)
	return &a
}

// Attribute returns the finished h.Attribute.
func (a *attrName) Attribute() h.Attribute {
	return h.Attr(a.name.String(), a.value)
}

// For renders a template once per item of source: data-for="item, index in
// source". An empty index omits the index variable, and an empty item and
// index emit the source expression alone, which makes Rocket default the
// names to "item" and "i".
//
// It goes on a <template> element. For panics unless the non-empty names are
// JavaScript identifiers; use ForRaw for an expression this shape cannot
// express.
//
//	rkt.For("item", "i", rkt.SignalRef("items"))  =>  data-for="item, i in $$items"
//	rkt.For("item", "", rkt.SignalRef("items"))   =>  data-for="item in $$items"
//	rkt.For("", "", rkt.SignalRef("items"))       =>  data-for="$$items"
func For(item, index string, source ds.Value) h.Attribute {
	src := js.ToJS(source.Expr())
	if item == "" {
		if index != "" {
			panic("rkt: For called with an index but no item name")
		}
		return h.Attr("data-for", src)
	}
	validateIdent("data-for item name", item)
	if index == "" {
		return h.Attr("data-for", item+" in "+src)
	}
	validateIdent("data-for index name", index)
	return h.Attr("data-for", item+", "+index+" in "+src)
}

// ForRaw is the escape hatch for a data-for expression this package cannot
// build: the expression is emitted verbatim. It goes on a <template> element.
//
//	rkt.ForRaw("item, i in $$items.filter(x => x.done)")
//	=>  data-for="item, i in $$items.filter(x => x.done)"
func ForRaw(expression string) h.Attribute {
	return h.Attr("data-for", expression)
}

// If renders a template when cond is truthy. It goes on a <template> element,
// and may be followed by adjacent sibling <template> elements carrying ElseIf
// and Else — the chain is broken by any other element between them.
//
//	rkt.If(rkt.Sig("open").Value())  =>  data-if="$$open"
func If(cond ds.Value) h.Attribute {
	return h.Attr("data-if", js.ToJS(cond.Expr()))
}

// ElseIf renders a template when every preceding condition in the chain was
// falsy and cond is truthy. It goes on a <template> element that is an
// adjacent sibling of the template carrying If.
//
//	rkt.ElseIf(rkt.Sig("busy").Value())  =>  data-else-if="$$busy"
func ElseIf(cond ds.Value) h.Attribute {
	return h.Attr("data-else-if", js.ToJS(cond.Expr()))
}

// Else renders a template when every preceding condition in the chain was
// falsy. It goes on a <template> element that is an adjacent sibling of the
// template carrying If (or the last ElseIf), and takes no value.
//
//	rkt.Else()  =>  data-else
func Else() h.Attribute {
	return h.Attr("data-else", "")
}

// RefBuilder builds a data-ref:<name> attribute, which stores a reference to
// the element in a signal. The Case modifier may be chained.
type RefBuilder struct {
	*attrName
}

// Case appends "__case.<casing>" — controls how the key in the attribute name
// is converted into a signal name. Keys are camelCased by default.
//
//	rkt.Ref("my-input").Case(ds.CamelCase)  =>  data-ref:my-input__case.camel
func (b *RefBuilder) Case(c ds.SignalCasing) *RefBuilder {
	b.name.WriteString("__case.")
	b.name.WriteString(string(c))
	return b
}

// Ref stores a reference to the element in the named signal, using Rocket's
// key syntax: data-ref:<name>. Modifiers (Case) may be chained.
//
//	rkt.Ref("input")  =>  data-ref:input
func Ref(name string) *RefBuilder {
	return &RefBuilder{attrName: newAttrName("data-ref:" + name)}
}

// Root appends "__root" to an already-built attribute's name, which escapes
// the component's private scope so the attribute binds to the global signal
// root instead of the component's private signals.
//
// It accepts any built attribute, but Rocket honors the modifier on
// data-signals, data-signals:<key>, data-bind, data-bind:<key>,
// data-computed:<key>, data-indicator and data-indicator:<key> — typically
// ds.Signal, ds.Signals, ds.Bind, ds.BindKey, ds.Computed, ds.Indicator and
// ds.IndicatorKey. A nil builder or a zero attribute is returned unchanged.
//
//	rkt.Root(ds.Bind("name"))       =>  data-bind__root="name"
//	rkt.Root(ds.Signal("count", 0)) =>  data-signals:count__root="0"
func Root(b h.AttrBuilder) h.Attribute {
	if b == nil {
		return h.Attribute{}
	}
	attr := b.Attribute()
	if attr.Name == "" {
		return h.Attribute{}
	}
	return h.Attr(attr.Name+"__root", attr.Value)
}
