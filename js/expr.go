package js

import "strings"

// Expr represents a JavaScript expression that produces a value.
// Expressions can be composed into larger expressions or used as statements.
type Expr interface {
	// js writes the JavaScript expression to the builder.
	js(sb *strings.Builder)
}

// Stmt represents a JavaScript statement.
// Statements are complete units of execution.
type Stmt interface {
	// stmt writes the JavaScript statement to the builder.
	// The statement should NOT include a trailing semicolon.
	stmt(sb *strings.Builder)
}

// Callable represents a JavaScript value that can have properties accessed
// and methods called on it (identifiers, objects, function results).
type Callable interface {
	Expr
	callable()
}

// exprStmt wraps an expression to be used as a statement.
type exprStmt struct {
	expr Expr
}

func (e exprStmt) stmt(sb *strings.Builder) {
	e.expr.js(sb)
}

// ExprStmt converts an expression to a statement.
func ExprStmt(expr Expr) Stmt {
	return exprStmt{expr}
}

// rawExpr represents a raw JavaScript expression string.
// This is the escape hatch for arbitrary JS.
type rawExpr string

func (r rawExpr) js(sb *strings.Builder) { sb.WriteString(string(r)) }
func (r rawExpr) callable()              {}

// Raw injects raw JavaScript code. This is the ONLY way to inject arbitrary JS.
// Use with caution as this bypasses type safety.
func Raw(code string) Callable {
	return rawExpr(code)
}
