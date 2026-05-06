package js

import "strings"

// Expr represents a JavaScript expression that produces a value.
// Expr is a fluent value: operator and access methods chain to build
// larger expressions.
//
// The zero value is invalid; construct an Expr via one of the package
// constructors (Int, String, Ident, Raw, etc.).
type Expr struct {
	node exprNode
}

// Stmt represents a JavaScript statement.
//
// Statements are produced by package functions (If, Let, Return, ...) and
// by methods on Expr that yield assignment-like statements (Assign,
// AddAssign, Incr, ...).
type Stmt struct {
	node stmtNode
}

// exprNode is the unexported interface satisfied by every expression
// AST node.
type exprNode interface {
	js(sb *strings.Builder)
}

// stmtNode is the unexported interface satisfied by every statement
// AST node.
type stmtNode interface {
	stmt(sb *strings.Builder)
}

// js writes the JavaScript expression to sb. nil nodes write nothing.
func (e Expr) js(sb *strings.Builder) {
	if e.node == nil {
		return
	}
	e.node.js(sb)
}

// stmt writes the JavaScript statement to sb. nil nodes write nothing.
func (s Stmt) stmt(sb *strings.Builder) {
	if s.node == nil {
		return
	}
	s.node.stmt(sb)
}

// rawNode represents arbitrary JS code injected via Raw().
type rawNode string

func (r rawNode) js(sb *strings.Builder) { sb.WriteString(string(r)) }

// Raw injects raw JavaScript code. This is the ONLY way to inject arbitrary JS.
// Use with caution as this bypasses type safety.
func Raw(code string) Expr {
	return Expr{node: rawNode(code)}
}

// exprStmtNode wraps an expression to be used as a statement.
type exprStmtNode struct{ expr Expr }

func (e exprStmtNode) stmt(sb *strings.Builder) { e.expr.js(sb) }

// Stmt converts an expression to a statement.
func (e Expr) Stmt() Stmt {
	return Stmt{node: exprStmtNode{e}}
}

// ExprStmt converts an expression to a statement. It is equivalent to
// e.Stmt() and is provided as a free function for symmetry with Return,
// Throw, etc.
func ExprStmt(e Expr) Stmt {
	return e.Stmt()
}
