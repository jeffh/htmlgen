package js

import "strings"

// isValidIdentifier reports whether name is a valid JavaScript identifier
// suitable for dot-notation property access. Empty strings and names
// containing characters outside [A-Za-z0-9_$] (or starting with a digit)
// are rejected so the caller can fall back to bracket notation.
func isValidIdentifier(name string) bool {
	if name == "" {
		return false
	}
	for i, r := range name {
		switch {
		case r >= 'a' && r <= 'z',
			r >= 'A' && r <= 'Z',
			r == '_', r == '$':
			// always allowed
		case r >= '0' && r <= '9':
			if i == 0 {
				return false
			}
		default:
			return false
		}
	}
	return true
}

// Prop accesses a property on a callable expression.
// Example: Prop(Ident("document"), "body") => document.body
// If name contains characters that are invalid for dot-notation access,
// bracket notation with a quoted string is used instead.
// Example: Prop(Ident("obj"), "foo-bar") => obj["foo-bar"]
func Prop(obj Callable, name string) Callable {
	return propAccess{obj: obj, prop: name}
}

type propAccess struct {
	obj  Callable
	prop string
}

func (p propAccess) js(sb *strings.Builder) {
	p.obj.js(sb)
	if isValidIdentifier(p.prop) {
		sb.WriteString(".")
		sb.WriteString(p.prop)
	} else {
		sb.WriteString("[")
		writeJSONString(sb, p.prop)
		sb.WriteString("]")
	}
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
	if isValidIdentifier(o.prop) {
		sb.WriteString("?.")
		sb.WriteString(o.prop)
	} else {
		sb.WriteString("?.[")
		writeJSONString(sb, o.prop)
		sb.WriteString("]")
	}
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
	if isValidIdentifier(o.method) {
		sb.WriteString("?.")
		sb.WriteString(o.method)
	} else {
		sb.WriteString("?.[")
		writeJSONString(sb, o.method)
		sb.WriteString("]")
	}
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
