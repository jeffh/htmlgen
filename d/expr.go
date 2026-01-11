package d

import (
	"strings"

	"github.com/jeffh/htmlgen/js"
)

// Re-export js types for convenience
type (
	Expr     = js.Expr
	Callable = js.Callable
	Stmt     = js.Stmt
	KV       = js.KV
)

// Value wraps a js.Expr and implements AttrMutator.
// This allows expressions to be used directly with Datastar attribute functions.
type Value struct {
	expr js.Expr
}

// Modify implements AttrMutator
func (v Value) Modify(attr *attrBuilder) {
	attr.AppendStatement(js.ToJS(v.expr))
}

// Expr returns the underlying js.Expr
func (v Value) Expr() js.Expr {
	return v.expr
}

// V wraps a js.Expr as a Value that can be used with Datastar attributes.
func V(expr js.Expr) Value {
	return Value{expr: expr}
}

// Raw injects raw JavaScript code. This is the escape hatch for arbitrary JavaScript.
// Use with caution as this bypasses type safety.
func Raw(s string) Value {
	return Value{expr: js.Raw(s)}
}

// JsonValue creates a JavaScript value from a Go value using JSON encoding.
// Panics if the value cannot be marshaled to JSON.
func JsonValue(value any) Value {
	return Value{expr: js.JSON(value)}
}

// Re-export js value constructors
var (
	// Str creates a JavaScript string literal, properly escaped.
	Str = js.String
	// Int creates a JavaScript number literal from an integer.
	Int = js.Int
	// Int64 creates a JavaScript number literal from an int64.
	Int64 = js.Int64
	// Float creates a JavaScript number literal from a float64.
	Float = js.Float
	// Bool creates a JavaScript boolean literal.
	Bool = js.Bool
	// Null creates a JavaScript null literal.
	Null = js.Null
	// Undefined creates a JavaScript undefined literal.
	Undefined = js.Undefined
	// JSON creates a JavaScript value from a Go value using JSON encoding.
	JSON = js.JSON
	// Array creates a JavaScript array literal from expressions.
	Array = js.Array
	// Object creates a JavaScript object literal from key-value pairs.
	Object = js.Object
	// Pair creates a key-value pair for Object().
	Pair = js.Pair
	// Ident creates a JavaScript identifier reference.
	Ident = js.Ident
	// This creates the special "this" identifier.
	This = js.This
	// ToJS converts an expression to its JavaScript string representation.
	ToJS = js.ToJS
	// ToJSStmt converts a statement to its JavaScript string representation.
	ToJSStmt = js.ToJSStmt
)

// Re-export js operators
var (
	Add             = js.Add
	Sub             = js.Sub
	Mul             = js.Mul
	Div             = js.Div
	Mod             = js.Mod
	Eq              = js.Eq
	NotEq           = js.NotEq
	Lt              = js.Lt
	LtEq            = js.LtEq
	Gt              = js.Gt
	GtEq            = js.GtEq
	JSAnd           = js.And
	JSOr            = js.Or
	JSNot           = js.Not
	Ternary         = js.Ternary
	NullishCoalesce = js.NullishCoalesce
	Group           = js.Group
)

// Re-export js property/method access
var (
	Prop         = js.Prop
	Method       = js.Method
	Index        = js.Index
	Call         = js.Call
	New          = js.New
	OptionalProp = js.OptionalProp
	OptionalCall = js.OptionalCall
)

// Re-export js statements
var (
	ExprStmt   = js.ExprStmt
	Let        = js.Let
	Const      = js.Const
	Assign     = js.Assign
	AddAssign  = js.AddAssign
	SubAssign  = js.SubAssign
	Incr       = js.Incr
	Decr       = js.Decr
	PreIncr    = js.PreIncr
	PreDecr    = js.PreDecr
	PostIncr   = js.PostIncr
	PostDecr   = js.PostDecr
	Return     = js.Return
	ReturnVoid = js.ReturnVoid
)

// Re-export js builtins
var (
	Console      = js.Console
	Document     = js.Document
	JSWindow     = js.Window
	Event        = js.Event
	JSConsoleLog = js.ConsoleLog
	ConsoleError = js.ConsoleError
	EventTarget  = js.EventTarget
	EventValue   = js.EventValue
)

// SignalRef creates a Datastar signal reference: $name
// Use this to reference a signal value in expressions.
// Example: SignalRef("count") produces $count
func SignalRef(name string) Value {
	// Remove $ prefix if already present
	name = strings.TrimPrefix(name, "$")
	return Value{expr: js.Raw("$" + name)}
}

// DatastarAction creates a Datastar action call: @action(args...)
// Example: DatastarAction("get", js.String("/api")) produces @get("/api")
func DatastarAction(name string, args ...js.Expr) js.Callable {
	var sb strings.Builder
	sb.WriteString("@")
	sb.WriteString(name)
	sb.WriteString("(")
	for i, arg := range args {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(js.ToJS(arg))
	}
	sb.WriteString(")")
	return js.Raw(sb.String())
}

// DSGet creates @get(path) Datastar action
func DSGet(path js.Expr) js.Callable {
	return DatastarAction("get", path)
}

// DSPost creates @post(path) Datastar action
func DSPost(path js.Expr) js.Callable {
	return DatastarAction("post", path)
}

// DSPut creates @put(path) Datastar action
func DSPut(path js.Expr) js.Callable {
	return DatastarAction("put", path)
}

// DSDelete creates @delete(path) Datastar action
func DSDelete(path js.Expr) js.Callable {
	return DatastarAction("delete", path)
}

// DSPatch creates @patch(path) Datastar action
func DSPatch(path js.Expr) js.Callable {
	return DatastarAction("patch", path)
}

// DSPeek creates @peek(() => expr) Datastar action for debugging
func DSPeek(expr js.Expr) js.Callable {
	var sb strings.Builder
	sb.WriteString("@peek(() => ")
	sb.WriteString(js.ToJS(expr))
	sb.WriteString(")")
	return js.Raw(sb.String())
}

// DSSetAll creates @setAll(value, filter) Datastar action
func DSSetAll(value js.Expr, filter *FilterOptions) js.Callable {
	var sb strings.Builder
	sb.WriteString("@setAll(")
	sb.WriteString(js.ToJS(value))
	if filter != nil && (filter.IncludeReg != nil || filter.ExcludeReg != nil) {
		sb.WriteString(", ")
		filter.appendJS(&sb)
	}
	sb.WriteString(")")
	return js.Raw(sb.String())
}

// DSToggleAll creates @toggleAll(filter) Datastar action
func DSToggleAll(filter *FilterOptions) js.Callable {
	var sb strings.Builder
	sb.WriteString("@toggleAll(")
	if filter != nil && (filter.IncludeReg != nil || filter.ExcludeReg != nil) {
		filter.appendJS(&sb)
	}
	sb.WriteString(")")
	return js.Raw(sb.String())
}

// DSClipboard creates @clipboard(text) Datastar Pro action
func DSClipboard(text js.Expr) js.Callable {
	return DatastarAction("clipboard", text)
}

// DSClipboardBase64 creates @clipboard(text, true) Datastar Pro action for Base64-decoded content
func DSClipboardBase64(text js.Expr) js.Callable {
	return DatastarAction("clipboard", text, js.Bool(true))
}

// DSFit creates @fit(v, oldMin, oldMax, newMin, newMax) Datastar Pro action
func DSFit(v, oldMin, oldMax, newMin, newMax js.Expr) js.Callable {
	return DatastarAction("fit", v, oldMin, oldMax, newMin, newMax)
}

// DSFitClamped creates @fit(v, oldMin, oldMax, newMin, newMax, true) with clamping
func DSFitClamped(v, oldMin, oldMax, newMin, newMax js.Expr) js.Callable {
	return DatastarAction("fit", v, oldMin, oldMax, newMin, newMax, js.Bool(true))
}

// DSFitRounded creates @fit(v, oldMin, oldMax, newMin, newMax, false, true) with rounding
func DSFitRounded(v, oldMin, oldMax, newMin, newMax js.Expr) js.Callable {
	return DatastarAction("fit", v, oldMin, oldMax, newMin, newMax, js.Bool(false), js.Bool(true))
}

// DSFitClampedRounded creates @fit(v, oldMin, oldMax, newMin, newMax, true, true) with clamping and rounding
func DSFitClampedRounded(v, oldMin, oldMax, newMin, newMax js.Expr) js.Callable {
	return DatastarAction("fit", v, oldMin, oldMax, newMin, newMax, js.Bool(true), js.Bool(true))
}

// PromiseChain represents a chainable action for HTTP requests (then/catch)
type PromiseChain interface {
	appendChain(sb *strings.Builder)
}

// thenChain represents .then(() => expr)
type thenChain struct {
	expr js.Expr
}

func (t thenChain) appendChain(sb *strings.Builder) {
	sb.WriteString(".then(() => ")
	sb.WriteString(js.ToJS(t.expr))
	sb.WriteString(")")
}

// catchChain represents .catch((error) => expr)
type catchChain struct {
	expr js.Expr
}

func (c catchChain) appendChain(sb *strings.Builder) {
	sb.WriteString(".catch((error) => ")
	sb.WriteString(js.ToJS(c.expr))
	sb.WriteString(")")
}

// ThenChain creates a .then() chain for successful request handling
func ThenChain(expr js.Expr) PromiseChain {
	return thenChain{expr}
}

// CatchChain creates a .catch() chain for error handling
func CatchChain(expr js.Expr) PromiseChain {
	return catchChain{expr}
}

// WithChains adds promise chains to a Datastar action, returning a new Callable
func WithChains(action js.Callable, chains ...PromiseChain) js.Callable {
	if len(chains) == 0 {
		return action
	}
	var sb strings.Builder
	sb.WriteString(js.ToJS(action))
	for _, chain := range chains {
		chain.appendChain(&sb)
	}
	return js.Raw(sb.String())
}

// ExprMutator wraps a js.Expr to satisfy AttrMutator
type ExprMutator struct {
	Expr js.Expr
}

func (e ExprMutator) Modify(attr *attrBuilder) {
	attr.AppendStatement(js.ToJS(e.Expr))
}

// E wraps a js.Expr to use as an AttrMutator
// Example: OnClick(PreventDefault(), E(Signal("count")))
func E(expr js.Expr) ExprMutator {
	return ExprMutator{Expr: expr}
}
