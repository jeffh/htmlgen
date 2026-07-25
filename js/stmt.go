package js

import "strings"

// Assignment statements

type assignStmt struct {
	target Expr
	value  Expr
}

func (a assignStmt) stmt(sb *strings.Builder) {
	a.target.js(sb)
	sb.WriteString(" = ")
	a.value.js(sb)
}

// Assign creates an assignment statement: target = value.
func (e Expr) Assign(value Expr) Stmt {
	return Stmt{node: assignStmt{target: e, value: value}}
}

// Compound assignment

type compoundAssign struct {
	target Expr
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

func (e Expr) compoundAssign(op string, value Expr) Stmt {
	return Stmt{node: compoundAssign{target: e, op: op, value: value}}
}

// AddAssign creates: e += value.
func (e Expr) AddAssign(value Expr) Stmt { return e.compoundAssign("+", value) }

// SubAssign creates: e -= value.
func (e Expr) SubAssign(value Expr) Stmt { return e.compoundAssign("-", value) }

// MulAssign creates: e *= value.
func (e Expr) MulAssign(value Expr) Stmt { return e.compoundAssign("*", value) }

// DivAssign creates: e /= value.
func (e Expr) DivAssign(value Expr) Stmt { return e.compoundAssign("/", value) }

// ModAssign creates: e %= value.
func (e Expr) ModAssign(value Expr) Stmt { return e.compoundAssign("%", value) }

// AndAssign creates: e &&= value.
func (e Expr) AndAssign(value Expr) Stmt { return e.compoundAssign("&&", value) }

// OrAssign creates: e ||= value.
func (e Expr) OrAssign(value Expr) Stmt { return e.compoundAssign("||", value) }

// NullishAssign creates: e ??= value.
func (e Expr) NullishAssign(value Expr) Stmt { return e.compoundAssign("??", value) }

// Variable declarations

type varDecl struct {
	kind  string // "let", "const", "var"
	name  string
	value Expr
	init  bool
}

func (v varDecl) stmt(sb *strings.Builder) {
	sb.WriteString(v.kind)
	sb.WriteString(" ")
	sb.WriteString(v.name)
	if v.init {
		sb.WriteString(" = ")
		v.value.js(sb)
	}
}

// Let creates a let declaration: let name = value.
func Let(name string, value Expr) Stmt {
	return Stmt{node: varDecl{kind: "let", name: name, value: value, init: true}}
}

// LetDecl creates a let declaration without initialization: let name.
func LetDecl(name string) Stmt {
	return Stmt{node: varDecl{kind: "let", name: name}}
}

// Const creates a const declaration: const name = value.
func Const(name string, value Expr) Stmt {
	return Stmt{node: varDecl{kind: "const", name: name, value: value, init: true}}
}

// Var creates a var declaration: var name = value.
func Var(name string, value Expr) Stmt {
	return Stmt{node: varDecl{kind: "var", name: name, value: value, init: true}}
}

// VarDecl creates a var declaration without initialization: var name.
func VarDecl(name string) Stmt {
	return Stmt{node: varDecl{kind: "var", name: name}}
}

// Increment / decrement

type incrDecr struct {
	target Expr
	op     string
	pre    bool
}

func (i incrDecr) write(sb *strings.Builder) {
	if i.pre {
		sb.WriteString(i.op)
		i.target.js(sb)
	} else {
		i.target.js(sb)
		sb.WriteString(i.op)
	}
}

func (i incrDecr) js(sb *strings.Builder)   { i.write(sb) }
func (i incrDecr) stmt(sb *strings.Builder) { i.write(sb) }

// Incr creates a post-increment statement: e++.
func (e Expr) Incr() Stmt { return Stmt{node: incrDecr{target: e, op: "++"}} }

// Decr creates a post-decrement statement: e--.
func (e Expr) Decr() Stmt { return Stmt{node: incrDecr{target: e, op: "--"}} }

// PreIncrExpr creates a pre-increment expression: ++e.
func (e Expr) PreIncrExpr() Expr { return Expr{node: incrDecr{target: e, op: "++", pre: true}} }

// PreDecrExpr creates a pre-decrement expression: --e.
func (e Expr) PreDecrExpr() Expr { return Expr{node: incrDecr{target: e, op: "--", pre: true}} }

// PostIncrExpr creates a post-increment expression: e++.
func (e Expr) PostIncrExpr() Expr { return Expr{node: incrDecr{target: e, op: "++"}} }

// PostDecrExpr creates a post-decrement expression: e--.
func (e Expr) PostDecrExpr() Expr { return Expr{node: incrDecr{target: e, op: "--"}} }

// Return statement

type returnStmt struct {
	value  Expr
	hasVal bool
}

func (r returnStmt) stmt(sb *strings.Builder) {
	sb.WriteString("return")
	if r.hasVal {
		sb.WriteString(" ")
		r.value.js(sb)
	}
}

// Return creates a return statement: return value.
func Return(value Expr) Stmt {
	return Stmt{node: returnStmt{value: value, hasVal: true}}
}

// ReturnVoid creates a bare return statement: return.
func ReturnVoid() Stmt {
	return Stmt{node: returnStmt{}}
}

// Throw statement

type throwStmt struct{ value Expr }

func (t throwStmt) stmt(sb *strings.Builder) {
	sb.WriteString("throw ")
	t.value.js(sb)
}

// Throw creates a throw statement: throw value.
func Throw(value Expr) Stmt {
	return Stmt{node: throwStmt{value}}
}

// Break / Continue

type breakStmt struct{ label string }

func (b breakStmt) stmt(sb *strings.Builder) {
	sb.WriteString("break")
	if b.label != "" {
		sb.WriteString(" ")
		sb.WriteString(b.label)
	}
}

// Break creates a break statement.
func Break() Stmt { return Stmt{node: breakStmt{}} }

// BreakLabel creates a labeled break: break label.
func BreakLabel(label string) Stmt { return Stmt{node: breakStmt{label}} }

type continueStmt struct{ label string }

func (c continueStmt) stmt(sb *strings.Builder) {
	sb.WriteString("continue")
	if c.label != "" {
		sb.WriteString(" ")
		sb.WriteString(c.label)
	}
}

// Continue creates a continue statement.
func Continue() Stmt { return Stmt{node: continueStmt{}} }

// ContinueLabel creates a labeled continue: continue label.
func ContinueLabel(label string) Stmt { return Stmt{node: continueStmt{label}} }

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

// If creates an if statement: if (cond) { body... }.
func If(cond Expr, body ...Stmt) Stmt {
	return Stmt{node: ifStmt{cond: cond, body: body}}
}

// IfElse creates an if-else statement.
func IfElse(cond Expr, thenBody []Stmt, elseBody []Stmt) Stmt {
	return Stmt{node: ifStmt{cond: cond, body: thenBody, elseBody: elseBody}}
}

// Try / catch statement

type tryStmt struct {
	body      []Stmt
	errName   string
	catchBody []Stmt
}

// writeStmtBlock writes body as a braced block. An empty body writes "{}" so
// that generated handlers stay compact.
func writeStmtBlock(sb *strings.Builder, body []Stmt) {
	if len(body) == 0 {
		sb.WriteString("{}")
		return
	}
	sb.WriteString("{ ")
	writeStmtList(sb, body)
	sb.WriteString(" }")
}

func (t tryStmt) stmt(sb *strings.Builder) {
	sb.WriteString("try ")
	writeStmtBlock(sb, t.body)
	sb.WriteString(" catch ")
	if t.errName != "" {
		sb.WriteString("(")
		sb.WriteString(t.errName)
		sb.WriteString(") ")
	}
	writeStmtBlock(sb, t.catchBody)
}

// TryCatch creates a try/catch statement. When errName is empty the optional
// catch binding is omitted.
//
//	TryCatch([]Stmt{ExprStmt(ConsoleLog(Int(1)))}, "e", []Stmt{ExprStmt(ConsoleError(Ident("e")))})
//	=>  try { console.log(1) } catch (e) { console.error(e) }
//
//	TryCatch([]Stmt{Ident("x").Incr()}, "", nil)
//	=>  try { x++ } catch {}
func TryCatch(body []Stmt, errName string, catchBody []Stmt) Stmt {
	return Stmt{node: tryStmt{body: body, errName: errName, catchBody: catchBody}}
}

// Try creates a try statement with an empty, unbound catch block. It is sugar
// for TryCatch(body, "", nil) and is useful for best-effort work whose failure
// should be ignored.
//
//	Try(ExprStmt(Ident("el").Method("focus")))
//	=>  try { el.focus() } catch {}
func Try(body ...Stmt) Stmt {
	return TryCatch(body, "", nil)
}

// Statement list

type stmtList []Stmt

func (s stmtList) stmt(sb *strings.Builder) {
	for i, st := range s {
		if i > 0 {
			sb.WriteString("; ")
		}
		st.stmt(sb)
	}
}

// Stmts combines multiple statements (semicolon-separated).
func Stmts(stmts ...Stmt) Stmt {
	return Stmt{node: stmtList(stmts)}
}

// Block statement

type blockStmt struct{ body []Stmt }

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

// Block creates a block statement: { body... }.
func Block(body ...Stmt) Stmt {
	return Stmt{node: blockStmt{body}}
}

// Debugger statement

type debuggerStmt struct{}

func (d debuggerStmt) stmt(sb *strings.Builder) { sb.WriteString("debugger") }

// Debugger creates a debugger statement.
func Debugger() Stmt { return Stmt{node: debuggerStmt{}} }
