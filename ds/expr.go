package ds

import (
	"strings"

	"github.com/jeffh/htmlgen/js"
)

// Re-export common js types for convenience.
type (
	Expr = js.Expr
	Stmt = js.Stmt
	KV   = js.KV
)

// Value wraps a js.Expr as a Datastar action/expression value. It is the
// argument type accepted by event handlers like OnClick and by the various
// signal-related attribute helpers.
type Value struct {
	expr js.Expr
}

// Expr returns the underlying js.Expr.
func (v Value) Expr() js.Expr { return v.expr }

// Not returns the logical negation of the Value. Values are immutable; a new
// Value is returned.
//
//	ds.SignalRef("open").Not()  =>  !$open
func (v Value) Not() Value { return Value{expr: v.expr.Not()} }

// And returns the logical conjunction of two Values.
//
//	ds.SignalRef("a").And(ds.SignalRef("b"))  =>  ($a && $b)
func (v Value) And(other Value) Value { return Value{expr: v.expr.And(other.expr)} }

// Or returns the logical disjunction of two Values.
//
//	ds.SignalRef("a").Or(ds.SignalRef("b"))  =>  ($a || $b)
func (v Value) Or(other Value) Value { return Value{expr: v.expr.Or(other.expr)} }

// Ternary uses the Value as a condition and returns a conditional expression.
//
//	ds.SignalRef("open").Ternary(ds.Str("Hide"), ds.Str("Show"))
//	=>  ($open ? "Hide" : "Show")
func (v Value) Ternary(ifTrue, ifFalse Value) Value {
	return Value{expr: v.expr.Ternary(ifTrue.expr, ifFalse.expr)}
}

// Do bridges typed js statements into a Datastar Value. Statements are joined
// with "; " by js.Stmts.
//
// Only valid where Datastar accepts a statement list — event handlers
// (data-on:*), data-init and data-effect. A statement list is NOT an
// expression, so a Do value must not be nested inside another expression
// (Ternary, And, Show, Text, ...); doing so produces invalid JavaScript.
//
//	ds.Do(js.Let("n", js.Int(1)), js.Ident("n").Incr())
//	=>  let n = 1; n++
func Do(stmts ...js.Stmt) Value {
	return Value{expr: js.Raw(js.ToJSStmt(js.Stmts(stmts...)))}
}

// V wraps a js.Expr as a Value.
func V(expr js.Expr) Value { return Value{expr: expr} }

// Raw injects raw JavaScript code as a Value. This is the escape hatch — use
// sparingly, as it bypasses type safety.
func Raw(s string) Value { return Value{expr: js.Raw(s)} }

// Str creates a JavaScript string literal Value, properly JSON-escaped.
func Str(s string) Value { return Value{expr: js.String(s)} }

// JsonValue encodes a Go value as JSON and wraps the result as a Value.
func JsonValue(value any) Value { return Value{expr: js.JSON(value)} }

// Re-export the js helpers used by in-repo callers. Other js constructors
// stay in the js package — import it directly.
var (
	Int  = js.Int
	ToJS = js.ToJS
)

// Datastar expression scope identifiers.
//
// Datastar evaluates attribute expressions with "evt" bound to the triggering
// event and "el" bound to the element carrying the attribute. This differs
// from legacy inline handlers (onclick="..."), where the event object is named
// "event" (js.Event). Inside data-on:*, data-init and friends, always use these.
var (
	// Evt is the Datastar event object: evt.
	Evt js.Expr = js.Ident("evt")
	// El is the Datastar element object: el.
	El js.Expr = js.Ident("el")
	// EvtTarget is evt.target.
	//
	//	ds.OnClick(ds.V(ds.EvtTarget.Method("blur")))  =>  data-on:click="evt.target.blur()"
	EvtTarget = Evt.Prop("target")
	// EvtValue is evt.target.value.
	//
	//	ds.Sig("q").SetExpr(ds.EvtValue)  =>  $q = evt.target.value
	EvtValue = Evt.Prop("target").Prop("value")
	// EvtKey is evt.key.
	//
	//	ds.Sig("k").SetExpr(ds.EvtKey)  =>  $k = evt.key
	EvtKey = Evt.Prop("key")
)

// SignalRef creates a Datastar signal reference: $name. Use this to reference
// a signal value in expressions. Any leading "$" is stripped.
//
// It is equivalent to Sig(name).Value(); prefer Sig when you reference the same
// signal more than once.
//
//	ds.SignalRef("open")  =>  $open
func SignalRef(name string) Value {
	return Sig(name).Value()
}

// DatastarAction creates a Datastar action call: @action(args...).
//
//	DatastarAction("get", js.String("/api"))  =>  @get("/api")
func DatastarAction(name string, args ...js.Expr) js.Expr {
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

// actionPeek creates @peek(() => expr). Used by Peek.
func actionPeek(expr js.Expr) js.Expr {
	var sb strings.Builder
	sb.WriteString("@peek(() => ")
	sb.WriteString(js.ToJS(expr))
	sb.WriteString(")")
	return js.Raw(sb.String())
}

// actionSetAll creates @setAll(value, filter). Used by SetAll.
func actionSetAll(value js.Expr, filter *FilterOptions) js.Expr {
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

// actionToggleAll creates @toggleAll(filter). Used by ToggleAll.
func actionToggleAll(filter *FilterOptions) js.Expr {
	var sb strings.Builder
	sb.WriteString("@toggleAll(")
	if filter != nil && (filter.IncludeReg != nil || filter.ExcludeReg != nil) {
		filter.appendJS(&sb)
	}
	sb.WriteString(")")
	return js.Raw(sb.String())
}

// PromiseChain represents a chainable action for HTTP requests (then/catch).
type PromiseChain interface {
	appendChain(sb *strings.Builder)
}

type thenChain struct{ expr js.Expr }

func (t thenChain) appendChain(sb *strings.Builder) {
	sb.WriteString(".then(() => ")
	sb.WriteString(js.ToJS(t.expr))
	sb.WriteString(")")
}

type catchChain struct{ expr js.Expr }

func (c catchChain) appendChain(sb *strings.Builder) {
	sb.WriteString(".catch((error) => ")
	sb.WriteString(js.ToJS(c.expr))
	sb.WriteString(")")
}

// withChains appends .then()/.catch() chains to a Datastar action.
func withChains(action js.Expr, chains ...PromiseChain) js.Expr {
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
