package js

import "strings"

// Prop accesses a property on a callable expression.
// Example: Prop(Ident("document"), "body") => document.body
func Prop(obj Callable, name string) Callable {
	return propAccess{obj: obj, prop: name}
}

type propAccess struct {
	obj  Callable
	prop string
}

func (p propAccess) js(sb *strings.Builder) {
	p.obj.js(sb)
	sb.WriteString(".")
	sb.WriteString(p.prop)
}
func (p propAccess) callable() {}

// Index accesses an element by index or computed property.
// Example: Index(Ident("arr"), Int(0)) => arr[0]
// Example: Index(Ident("obj"), String("key")) => obj["key"]
func Index(obj Callable, index Expr) Callable {
	return indexAccess{obj: obj, index: index}
}

type indexAccess struct {
	obj   Callable
	index Expr
}

func (i indexAccess) js(sb *strings.Builder) {
	i.obj.js(sb)
	sb.WriteString("[")
	i.index.js(sb)
	sb.WriteString("]")
}
func (i indexAccess) callable() {}

// Call invokes a callable with arguments.
// Example: Call(Ident("alert"), String("hello")) => alert("hello")
func Call(fn Callable, args ...Expr) Callable {
	return funcCall{fn: fn, args: args}
}

type funcCall struct {
	fn   Callable
	args []Expr
}

func (f funcCall) js(sb *strings.Builder) {
	f.fn.js(sb)
	sb.WriteString("(")
	for i, arg := range f.args {
		if i > 0 {
			sb.WriteString(", ")
		}
		arg.js(sb)
	}
	sb.WriteString(")")
}
func (f funcCall) callable() {}

// Method calls a method on an object with arguments.
// Example: Method(Ident("console"), "log", String("hello")) => console.log("hello")
func Method(obj Callable, method string, args ...Expr) Callable {
	return funcCall{
		fn:   propAccess{obj: obj, prop: method},
		args: args,
	}
}

// New creates a new instance with the new keyword.
// Example: New(Ident("Date")) => new Date()
func New(constructor Callable, args ...Expr) Callable {
	return newExpr{constructor: constructor, args: args}
}

type newExpr struct {
	constructor Callable
	args        []Expr
}

func (n newExpr) js(sb *strings.Builder) {
	sb.WriteString("new ")
	n.constructor.js(sb)
	sb.WriteString("(")
	for i, arg := range n.args {
		if i > 0 {
			sb.WriteString(", ")
		}
		arg.js(sb)
	}
	sb.WriteString(")")
}
func (n newExpr) callable() {}

// OptionalProp accesses a property with optional chaining.
// Example: OptionalProp(Ident("obj"), "foo") => obj?.foo
func OptionalProp(obj Callable, name string) Callable {
	return optionalChain{obj, name}
}

type optionalChain struct {
	obj  Callable
	prop string
}

func (o optionalChain) js(sb *strings.Builder) {
	o.obj.js(sb)
	sb.WriteString("?.")
	sb.WriteString(o.prop)
}
func (o optionalChain) callable() {}

// OptionalCall calls a method with optional chaining.
// Example: OptionalCall(Ident("obj"), "method", args...) => obj?.method(args...)
func OptionalCall(obj Callable, method string, args ...Expr) Callable {
	return optionalMethodCall{obj, method, args}
}

type optionalMethodCall struct {
	obj    Callable
	method string
	args   []Expr
}

func (o optionalMethodCall) js(sb *strings.Builder) {
	o.obj.js(sb)
	sb.WriteString("?.")
	sb.WriteString(o.method)
	sb.WriteString("(")
	for i, arg := range o.args {
		if i > 0 {
			sb.WriteString(", ")
		}
		arg.js(sb)
	}
	sb.WriteString(")")
}
func (o optionalMethodCall) callable() {}
