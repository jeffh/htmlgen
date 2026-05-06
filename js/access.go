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

// propAccess represents obj.prop or obj["prop"].
type propAccess struct {
	obj  Expr
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

// Prop accesses a named property on the expression.
//
//	Ident("document").Prop("body") => document.body
//
// If name is not a valid JS identifier, bracket notation is used instead.
func (e Expr) Prop(name string) Expr {
	return Expr{node: propAccess{obj: e, prop: name}}
}

// indexAccess represents obj[index].
type indexAccess struct {
	obj   Expr
	index Expr
}

func (i indexAccess) js(sb *strings.Builder) {
	i.obj.js(sb)
	sb.WriteString("[")
	i.index.js(sb)
	sb.WriteString("]")
}

// Index returns e[index] — computed property access.
func (e Expr) Index(index Expr) Expr {
	return Expr{node: indexAccess{obj: e, index: index}}
}

// funcCall represents fn(args...).
type funcCall struct {
	fn   Expr
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

// Call invokes the expression as a function with the given arguments.
//
//	Ident("alert").Call(String("hello")) => alert("hello")
func (e Expr) Call(args ...Expr) Expr {
	return Expr{node: funcCall{fn: e, args: args}}
}

// Method calls a method on the expression.
//
//	Ident("console").Method("log", String("hello")) => console.log("hello")
func (e Expr) Method(name string, args ...Expr) Expr {
	return Expr{
		node: funcCall{
			fn:   Expr{node: propAccess{obj: e, prop: name}},
			args: args,
		},
	}
}

// newExpr represents `new Constructor(args...)`.
type newExpr struct {
	constructor Expr
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

// New constructs a new instance using the expression as the constructor.
//
//	Ident("Date").New() => new Date()
func (e Expr) New(args ...Expr) Expr {
	return Expr{node: newExpr{constructor: e, args: args}}
}

// optionalChain represents obj?.prop or obj?.["prop"].
type optionalChain struct {
	obj  Expr
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

// OptionalProp accesses a property using optional chaining.
//
//	Ident("obj").OptionalProp("foo") => obj?.foo
func (e Expr) OptionalProp(name string) Expr {
	return Expr{node: optionalChain{obj: e, prop: name}}
}

// optionalMethodCall represents obj?.method(args...).
type optionalMethodCall struct {
	obj    Expr
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

// OptionalCall calls a method using optional chaining.
//
//	Ident("obj").OptionalCall("method", args...) => obj?.method(args...)
func (e Expr) OptionalCall(method string, args ...Expr) Expr {
	return Expr{node: optionalMethodCall{obj: e, method: method, args: args}}
}
