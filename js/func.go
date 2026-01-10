package js

import "strings"

// ArrowFunc creates an arrow function expression with a single expression body.
// Example: ArrowFunc([]string{"x", "y"}, Add(Ident("x"), Ident("y")))
//
//	=> (x, y) => (x + y)
func ArrowFunc(params []string, body Expr) Callable {
	return arrowFuncExpr{params: params, body: body}
}

type arrowFuncExpr struct {
	params []string
	body   Expr
}

func (a arrowFuncExpr) js(sb *strings.Builder) {
	if len(a.params) == 1 {
		sb.WriteString(a.params[0])
	} else {
		sb.WriteString("(")
		for i, p := range a.params {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(p)
		}
		sb.WriteString(")")
	}
	sb.WriteString(" => ")
	a.body.js(sb)
}
func (a arrowFuncExpr) callable() {}

// ArrowFuncStmts creates an arrow function with a statement body.
// Example: ArrowFuncStmts([]string{"e"}, ExprStmt(ConsoleLog(Ident("e"))))
//
//	=> (e) => { console.log(e) }
func ArrowFuncStmts(params []string, stmts ...Stmt) Callable {
	return arrowFuncStmtsExpr{params: params, body: stmts}
}

type arrowFuncStmtsExpr struct {
	params []string
	body   []Stmt
}

func (a arrowFuncStmtsExpr) js(sb *strings.Builder) {
	if len(a.params) == 1 {
		sb.WriteString(a.params[0])
	} else {
		sb.WriteString("(")
		for i, p := range a.params {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(p)
		}
		sb.WriteString(")")
	}
	sb.WriteString(" => { ")
	for i, s := range a.body {
		if i > 0 {
			sb.WriteString("; ")
		}
		s.stmt(sb)
	}
	sb.WriteString(" }")
}
func (a arrowFuncStmtsExpr) callable() {}

// Func creates an anonymous function expression.
// Example: Func([]string{"x", "y"}, Return(Add(Ident("x"), Ident("y"))))
//
//	=> function(x, y) { return (x + y) }
func Func(params []string, stmts ...Stmt) Callable {
	return funcExpr{params: params, body: stmts}
}

type funcExpr struct {
	params []string
	body   []Stmt
}

func (f funcExpr) js(sb *strings.Builder) {
	sb.WriteString("function(")
	for i, p := range f.params {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(p)
	}
	sb.WriteString(") { ")
	for i, s := range f.body {
		if i > 0 {
			sb.WriteString("; ")
		}
		s.stmt(sb)
	}
	sb.WriteString(" }")
}
func (f funcExpr) callable() {}

// IIFE creates an immediately invoked function expression.
// Example: IIFE(ExprStmt(ConsoleLog(String("hello"))))
//
//	=> (function() { console.log("hello") })()
func IIFE(stmts ...Stmt) Callable {
	return iifeExpr{body: stmts}
}

type iifeExpr struct {
	body []Stmt
}

func (i iifeExpr) js(sb *strings.Builder) {
	sb.WriteString("(function() { ")
	for j, s := range i.body {
		if j > 0 {
			sb.WriteString("; ")
		}
		s.stmt(sb)
	}
	sb.WriteString(" })()")
}
func (i iifeExpr) callable() {}

// Template creates a template literal expression.
// Alternates between string parts and expression parts.
// Example: Template("Hello, ", Ident("name"), "!")
//
//	=> `Hello, ${name}!`
func Template(parts ...any) Callable {
	return templateLiteral{parts}
}

type templateLiteral struct {
	parts []any // alternating strings and Expr
}

func (t templateLiteral) js(sb *strings.Builder) {
	sb.WriteString("`")
	for _, part := range t.parts {
		switch v := part.(type) {
		case string:
			// Escape backticks, backslashes, and ${
			for _, r := range v {
				switch r {
				case '`':
					sb.WriteString("\\`")
				case '\\':
					sb.WriteString("\\\\")
				case '$':
					// Escape $ to prevent accidental interpolation
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
func (t templateLiteral) callable() {}

// Await creates an await expression.
// Example: Await(Fetch(String("/api/data")))
//
//	=> await fetch("/api/data")
func Await(expr Expr) Callable {
	return awaitExpr{expr}
}

type awaitExpr struct {
	expr Expr
}

func (a awaitExpr) js(sb *strings.Builder) {
	sb.WriteString("await ")
	a.expr.js(sb)
}
func (a awaitExpr) callable() {}

// AsyncArrowFunc creates an async arrow function with a single expression body.
// Example: AsyncArrowFunc([]string{}, Await(Fetch(String("/api"))))
//
//	=> async () => await fetch("/api")
func AsyncArrowFunc(params []string, body Expr) Callable {
	return asyncArrowFuncExpr{params: params, body: body}
}

type asyncArrowFuncExpr struct {
	params []string
	body   Expr
}

func (a asyncArrowFuncExpr) js(sb *strings.Builder) {
	sb.WriteString("async ")
	if len(a.params) == 1 {
		sb.WriteString(a.params[0])
	} else {
		sb.WriteString("(")
		for i, p := range a.params {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(p)
		}
		sb.WriteString(")")
	}
	sb.WriteString(" => ")
	a.body.js(sb)
}
func (a asyncArrowFuncExpr) callable() {}

// AsyncArrowFuncStmts creates an async arrow function with a statement body.
// Example: AsyncArrowFuncStmts([]string{}, Let("data", Await(Fetch(String("/api")))))
//
//	=> async () => { let data = await fetch("/api") }
func AsyncArrowFuncStmts(params []string, stmts ...Stmt) Callable {
	return asyncArrowFuncStmtsExpr{params: params, body: stmts}
}

type asyncArrowFuncStmtsExpr struct {
	params []string
	body   []Stmt
}

func (a asyncArrowFuncStmtsExpr) js(sb *strings.Builder) {
	sb.WriteString("async ")
	if len(a.params) == 1 {
		sb.WriteString(a.params[0])
	} else {
		sb.WriteString("(")
		for i, p := range a.params {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(p)
		}
		sb.WriteString(")")
	}
	sb.WriteString(" => { ")
	for i, s := range a.body {
		if i > 0 {
			sb.WriteString("; ")
		}
		s.stmt(sb)
	}
	sb.WriteString(" }")
}
func (a asyncArrowFuncStmtsExpr) callable() {}

// PromiseThen creates expr.then(onFulfilled)
func PromiseThen(promise Callable, onFulfilled Expr) Callable {
	return Method(promise, "then", onFulfilled)
}

// PromiseCatch creates expr.catch(onRejected)
func PromiseCatch(promise Callable, onRejected Expr) Callable {
	return Method(promise, "catch", onRejected)
}

// PromiseFinally creates expr.finally(onFinally)
func PromiseFinally(promise Callable, onFinally Expr) Callable {
	return Method(promise, "finally", onFinally)
}

// PromiseResolve creates Promise.resolve(value)
func PromiseResolve(value Expr) Callable {
	return Method(Promise, "resolve", value)
}

// PromiseReject creates Promise.reject(reason)
func PromiseReject(reason Expr) Callable {
	return Method(Promise, "reject", reason)
}

// PromiseAll creates Promise.all(iterable)
func PromiseAll(iterable Expr) Callable {
	return Method(Promise, "all", iterable)
}

// PromiseRace creates Promise.race(iterable)
func PromiseRace(iterable Expr) Callable {
	return Method(Promise, "race", iterable)
}
