package js

import "strings"

// Assignment statements

type assignStmt struct {
	target Callable
	value  Expr
}

func (a assignStmt) stmt(sb *strings.Builder) {
	a.target.js(sb)
	sb.WriteString(" = ")
	a.value.js(sb)
}

// Assign creates an assignment statement: target = value
func Assign(target Callable, value Expr) Stmt {
	return assignStmt{target, value}
}

// Compound assignment

type compoundAssign struct {
	target Callable
	op     string
	value  Expr
}

func (c compoundAssign) stmt(sb *strings.Builder) {
	c.target.js(sb)
	sb.WriteString(" ")
	sb.WriteString(c.op)
	sb.WriteString("= ")
	c.value.js(sb)
}

// AddAssign creates: target += value
func AddAssign(target Callable, value Expr) Stmt {
	return compoundAssign{target, "+", value}
}

// SubAssign creates: target -= value
func SubAssign(target Callable, value Expr) Stmt {
	return compoundAssign{target, "-", value}
}

// MulAssign creates: target *= value
func MulAssign(target Callable, value Expr) Stmt {
	return compoundAssign{target, "*", value}
}

// DivAssign creates: target /= value
func DivAssign(target Callable, value Expr) Stmt {
	return compoundAssign{target, "/", value}
}

// ModAssign creates: target %= value
func ModAssign(target Callable, value Expr) Stmt {
	return compoundAssign{target, "%", value}
}

// AndAssign creates: target &&= value
func AndAssign(target Callable, value Expr) Stmt {
	return compoundAssign{target, "&&", value}
}

// OrAssign creates: target ||= value
func OrAssign(target Callable, value Expr) Stmt {
	return compoundAssign{target, "||", value}
}

// NullishAssign creates: target ??= value
func NullishAssign(target Callable, value Expr) Stmt {
	return compoundAssign{target, "??", value}
}

// Variable declarations

type varDecl struct {
	kind  string // "let", "const", "var"
	name  string
	value Expr // nil for declaration without initialization
}

func (v varDecl) stmt(sb *strings.Builder) {
	sb.WriteString(v.kind)
	sb.WriteString(" ")
	sb.WriteString(v.name)
	if v.value != nil {
		sb.WriteString(" = ")
		v.value.js(sb)
	}
}

// Let creates a let declaration: let name = value
func Let(name string, value Expr) Stmt {
	return varDecl{"let", name, value}
}

// LetDecl creates a let declaration without initialization: let name
func LetDecl(name string) Stmt {
	return varDecl{"let", name, nil}
}

// Const creates a const declaration: const name = value
func Const(name string, value Expr) Stmt {
	return varDecl{"const", name, value}
}

// Var creates a var declaration: var name = value
func Var(name string, value Expr) Stmt {
	return varDecl{"var", name, value}
}

// VarDecl creates a var declaration without initialization: var name
func VarDecl(name string) Stmt {
	return varDecl{"var", name, nil}
}

// Increment/Decrement

type incrDecr struct {
	target Callable
	op     string
	pre    bool
}

func (i incrDecr) stmt(sb *strings.Builder) {
	if i.pre {
		sb.WriteString(i.op)
		i.target.js(sb)
	} else {
		i.target.js(sb)
		sb.WriteString(i.op)
	}
}

// Expressions versions for use in larger expressions
func (i incrDecr) js(sb *strings.Builder) { i.stmt(sb) }
func (i incrDecr) callable()              {}

// Incr creates: target++ (post-increment statement)
func Incr(target Callable) Stmt { return incrDecr{target, "++", false} }

// Decr creates: target-- (post-decrement statement)
func Decr(target Callable) Stmt { return incrDecr{target, "--", false} }

// PreIncr creates: ++target (usable as expression)
func PreIncr(target Callable) Callable { return incrDecr{target, "++", true} }

// PreDecr creates: --target (usable as expression)
func PreDecr(target Callable) Callable { return incrDecr{target, "--", true} }

// PostIncr creates: target++ (usable as expression)
func PostIncr(target Callable) Callable { return incrDecr{target, "++", false} }

// PostDecr creates: target-- (usable as expression)
func PostDecr(target Callable) Callable { return incrDecr{target, "--", false} }

// Return statement

type returnStmt struct {
	value Expr // nil for bare return
}

func (r returnStmt) stmt(sb *strings.Builder) {
	sb.WriteString("return")
	if r.value != nil {
		sb.WriteString(" ")
		r.value.js(sb)
	}
}

// Return creates a return statement: return value
func Return(value Expr) Stmt {
	return returnStmt{value}
}

// ReturnVoid creates a bare return statement: return
func ReturnVoid() Stmt {
	return returnStmt{nil}
}

// Throw statement

type throwStmt struct {
	value Expr
}

func (t throwStmt) stmt(sb *strings.Builder) {
	sb.WriteString("throw ")
	t.value.js(sb)
}

// Throw creates a throw statement: throw value
func Throw(value Expr) Stmt {
	return throwStmt{value}
}

// Break statement

type breakStmt struct {
	label string
}

func (b breakStmt) stmt(sb *strings.Builder) {
	sb.WriteString("break")
	if b.label != "" {
		sb.WriteString(" ")
		sb.WriteString(b.label)
	}
}

// Break creates a break statement
func Break() Stmt {
	return breakStmt{}
}

// BreakLabel creates a break statement with a label: break label
func BreakLabel(label string) Stmt {
	return breakStmt{label}
}

// Continue statement

type continueStmt struct {
	label string
}

func (c continueStmt) stmt(sb *strings.Builder) {
	sb.WriteString("continue")
	if c.label != "" {
		sb.WriteString(" ")
		sb.WriteString(c.label)
	}
}

// Continue creates a continue statement
func Continue() Stmt {
	return continueStmt{}
}

// ContinueLabel creates a continue statement with a label: continue label
func ContinueLabel(label string) Stmt {
	return continueStmt{label}
}

// If statement

type ifStmt struct {
	cond     Expr
	body     []Stmt
	elseBody []Stmt
}

func (i ifStmt) stmt(sb *strings.Builder) {
	sb.WriteString("if (")
	i.cond.js(sb)
	sb.WriteString(") { ")
	for j, s := range i.body {
		if j > 0 {
			sb.WriteString("; ")
		}
		s.stmt(sb)
	}
	sb.WriteString(" }")
	if len(i.elseBody) > 0 {
		sb.WriteString(" else { ")
		for j, s := range i.elseBody {
			if j > 0 {
				sb.WriteString("; ")
			}
			s.stmt(sb)
		}
		sb.WriteString(" }")
	}
}

// If creates an if statement: if (cond) { body... }
func If(cond Expr, body ...Stmt) Stmt {
	return ifStmt{cond: cond, body: body}
}

// IfElse creates an if-else statement: if (cond) { thenBody... } else { elseBody... }
func IfElse(cond Expr, thenBody []Stmt, elseBody []Stmt) Stmt {
	return ifStmt{cond: cond, body: thenBody, elseBody: elseBody}
}

// Statement list

type stmtList []Stmt

func (s stmtList) stmt(sb *strings.Builder) {
	for i, stmt := range s {
		if i > 0 {
			sb.WriteString("; ")
		}
		stmt.stmt(sb)
	}
}

// Stmts combines multiple statements (semicolon-separated).
func Stmts(stmts ...Stmt) Stmt {
	return stmtList(stmts)
}

// Block statement

type blockStmt struct {
	body []Stmt
}

func (b blockStmt) stmt(sb *strings.Builder) {
	sb.WriteString("{ ")
	for i, s := range b.body {
		if i > 0 {
			sb.WriteString("; ")
		}
		s.stmt(sb)
	}
	sb.WriteString(" }")
}

// Block creates a block statement: { body... }
func Block(body ...Stmt) Stmt {
	return blockStmt{body}
}

// Debugger statement

type debuggerStmt struct{}

func (d debuggerStmt) stmt(sb *strings.Builder) {
	sb.WriteString("debugger")
}

// Debugger creates a debugger statement
func Debugger() Stmt {
	return debuggerStmt{}
}
