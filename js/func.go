package js

import "strings"

// writeArrowParams writes arrow function parameters in the format:
//   - Single param: x
//   - Zero or multiple params: (a, b)
func writeArrowParams(sb *strings.Builder, params []string) {
	if len(params) == 1 {
		sb.WriteString(params[0])
	} else {
		sb.WriteString("(")
		for i, p := range params {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(p)
		}
		sb.WriteString(")")
	}
}

// writeParenParams writes parenthesized parameters: (a, b).
func writeParenParams(sb *strings.Builder, params []string) {
	sb.WriteString("(")
	for i, p := range params {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(p)
	}
	sb.WriteString(")")
}

// writeStmtList writes statements separated by "; ".
func writeStmtList(sb *strings.Builder, stmts []Stmt) {
	for i, s := range stmts {
		if i > 0 {
			sb.WriteString("; ")
		}
		s.stmt(sb)
	}
}

// arrowFuncExpr represents an arrow function with a single expression body.
type arrowFuncExpr struct {
	params []string
	body   Expr
}

func (a arrowFuncExpr) js(sb *strings.Builder) {
	writeArrowParams(sb, a.params)
	sb.WriteString(" => ")
	a.body.js(sb)
}

// ArrowFunc creates an arrow function expression with a single expression body.
//
//	ArrowFunc([]string{"x", "y"}, Ident("x").Add(Ident("y")))
//	=> (x, y) => (x + y)
func ArrowFunc(params []string, body Expr) Expr {
	return Expr{node: arrowFuncExpr{params: params, body: body}}
}

// arrowFuncStmtsExpr represents an arrow function with a statement body.
type arrowFuncStmtsExpr struct {
	params []string
	body   []Stmt
}

func (a arrowFuncStmtsExpr) js(sb *strings.Builder) {
	writeArrowParams(sb, a.params)
	sb.WriteString(" => { ")
	writeStmtList(sb, a.body)
	sb.WriteString(" }")
}

// ArrowFuncStmts creates an arrow function with a statement body.
//
//	ArrowFuncStmts([]string{"e"}, ConsoleLog(Ident("e")).Stmt())
//	=> (e) => { console.log(e) }
func ArrowFuncStmts(params []string, stmts ...Stmt) Expr {
	return Expr{node: arrowFuncStmtsExpr{params: params, body: stmts}}
}

// funcExpr represents an anonymous function expression.
type funcExpr struct {
	params []string
	body   []Stmt
}

func (f funcExpr) js(sb *strings.Builder) {
	sb.WriteString("function")
	writeParenParams(sb, f.params)
	sb.WriteString(" { ")
	writeStmtList(sb, f.body)
	sb.WriteString(" }")
}

// Func creates an anonymous function expression.
//
//	Func([]string{"x", "y"}, Return(Ident("x").Add(Ident("y"))))
//	=> function(x, y) { return (x + y) }
func Func(params []string, stmts ...Stmt) Expr {
	return Expr{node: funcExpr{params: params, body: stmts}}
}

// iifeExpr represents an immediately invoked function expression.
type iifeExpr struct {
	body []Stmt
}

func (i iifeExpr) js(sb *strings.Builder) {
	sb.WriteString("(function() { ")
	writeStmtList(sb, i.body)
	sb.WriteString(" })()")
}

// IIFE creates an immediately invoked function expression.
//
//	IIFE(ConsoleLog(String("hello")).Stmt())
//	=> (function() { console.log("hello") })()
func IIFE(stmts ...Stmt) Expr {
	return Expr{node: iifeExpr{body: stmts}}
}

// templateLiteral represents a template literal with interpolation.
type templateLiteral struct {
	parts []any // alternating strings and Expr
}

func (t templateLiteral) js(sb *strings.Builder) {
	sb.WriteString("`")
	for _, part := range t.parts {
		switch v := part.(type) {
		case string:
			for _, r := range v {
				switch r {
				case '`':
					sb.WriteString("\\`")
				case '\\':
					sb.WriteString("\\\\")
				case '$':
					sb.WriteString("\\$")
				default:
					sb.WriteRune(r)
				}
			}
		case Expr:
			sb.WriteString("${")
			v.js(sb)
			sb.WriteString("}")
		}
	}
	sb.WriteString("`")
}

// Template creates a template literal expression.
// Alternates between string parts and expression parts.
//
//	Template("Hello, ", Ident("name"), "!") => `Hello, ${name}!`
func Template(parts ...any) Expr {
	return Expr{node: templateLiteral{parts}}
}

// awaitExpr represents an await expression.
type awaitExpr struct{ expr Expr }

func (a awaitExpr) js(sb *strings.Builder) {
	sb.WriteString("await ")
	a.expr.js(sb)
}

// Await wraps the expression in an await expression: await e.
func (e Expr) Await() Expr { return Expr{node: awaitExpr{e}} }

// asyncArrowFuncExpr represents an async arrow function with a single
// expression body.
type asyncArrowFuncExpr struct {
	params []string
	body   Expr
}

func (a asyncArrowFuncExpr) js(sb *strings.Builder) {
	sb.WriteString("async ")
	writeArrowParams(sb, a.params)
	sb.WriteString(" => ")
	a.body.js(sb)
}

// AsyncArrowFunc creates an async arrow function with a single expression body.
//
//	AsyncArrowFunc([]string{}, Fetch(String("/api")).Await())
//	=> async () => await fetch("/api")
func AsyncArrowFunc(params []string, body Expr) Expr {
	return Expr{node: asyncArrowFuncExpr{params: params, body: body}}
}

// asyncArrowFuncStmtsExpr represents an async arrow function with a statement body.
type asyncArrowFuncStmtsExpr struct {
	params []string
	body   []Stmt
}

func (a asyncArrowFuncStmtsExpr) js(sb *strings.Builder) {
	sb.WriteString("async ")
	writeArrowParams(sb, a.params)
	sb.WriteString(" => { ")
	writeStmtList(sb, a.body)
	sb.WriteString(" }")
}

// AsyncArrowFuncStmts creates an async arrow function with a statement body.
func AsyncArrowFuncStmts(params []string, stmts ...Stmt) Expr {
	return Expr{node: asyncArrowFuncStmtsExpr{params: params, body: stmts}}
}

// Then creates e.then(onFulfilled).
func (e Expr) Then(onFulfilled Expr) Expr {
	return e.Method("then", onFulfilled)
}

// Catch creates e.catch(onRejected).
func (e Expr) Catch(onRejected Expr) Expr {
	return e.Method("catch", onRejected)
}

// Finally creates e.finally(onFinally).
func (e Expr) Finally(onFinally Expr) Expr {
	return e.Method("finally", onFinally)
}

// PromiseResolve creates Promise.resolve(value).
func PromiseResolve(value Expr) Expr {
	return Promise.Method("resolve", value)
}

// PromiseReject creates Promise.reject(reason).
func PromiseReject(reason Expr) Expr {
	return Promise.Method("reject", reason)
}

// PromiseAll creates Promise.all(iterable).
func PromiseAll(iterable Expr) Expr {
	return Promise.Method("all", iterable)
}

// PromiseRace creates Promise.race(iterable).
func PromiseRace(iterable Expr) Expr {
	return Promise.Method("race", iterable)
}
