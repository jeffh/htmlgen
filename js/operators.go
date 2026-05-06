package js

import "strings"

// Binary operators

type binaryOp struct {
	left  Expr
	op    string
	right Expr
}

func (b binaryOp) js(sb *strings.Builder) {
	sb.WriteString("(")
	b.left.js(sb)
	sb.WriteString(" ")
	sb.WriteString(b.op)
	sb.WriteString(" ")
	b.right.js(sb)
	sb.WriteString(")")
}

func (e Expr) binary(op string, other Expr) Expr {
	return Expr{node: binaryOp{left: e, op: op, right: other}}
}

// Add returns (e + other).
func (e Expr) Add(other Expr) Expr { return e.binary("+", other) }

// Sub returns (e - other).
func (e Expr) Sub(other Expr) Expr { return e.binary("-", other) }

// Mul returns (e * other).
func (e Expr) Mul(other Expr) Expr { return e.binary("*", other) }

// Div returns (e / other).
func (e Expr) Div(other Expr) Expr { return e.binary("/", other) }

// Mod returns (e % other).
func (e Expr) Mod(other Expr) Expr { return e.binary("%", other) }

// Eq returns (e === other) — strict equality.
func (e Expr) Eq(other Expr) Expr { return e.binary("===", other) }

// NotEq returns (e !== other) — strict inequality.
func (e Expr) NotEq(other Expr) Expr { return e.binary("!==", other) }

// LooseEq returns (e == other) — loose equality.
func (e Expr) LooseEq(other Expr) Expr { return e.binary("==", other) }

// LooseNotEq returns (e != other) — loose inequality.
func (e Expr) LooseNotEq(other Expr) Expr { return e.binary("!=", other) }

// Lt returns (e < other).
func (e Expr) Lt(other Expr) Expr { return e.binary("<", other) }

// LtEq returns (e <= other).
func (e Expr) LtEq(other Expr) Expr { return e.binary("<=", other) }

// Gt returns (e > other).
func (e Expr) Gt(other Expr) Expr { return e.binary(">", other) }

// GtEq returns (e >= other).
func (e Expr) GtEq(other Expr) Expr { return e.binary(">=", other) }

// And returns (e && other).
func (e Expr) And(other Expr) Expr { return e.binary("&&", other) }

// Or returns (e || other).
func (e Expr) Or(other Expr) Expr { return e.binary("||", other) }

// NullishCoalesce returns (e ?? other).
func (e Expr) NullishCoalesce(other Expr) Expr { return e.binary("??", other) }

// BitwiseAnd returns (e & other).
func (e Expr) BitwiseAnd(other Expr) Expr { return e.binary("&", other) }

// BitwiseOr returns (e | other).
func (e Expr) BitwiseOr(other Expr) Expr { return e.binary("|", other) }

// BitwiseXor returns (e ^ other).
func (e Expr) BitwiseXor(other Expr) Expr { return e.binary("^", other) }

// ShiftLeft returns (e << other).
func (e Expr) ShiftLeft(other Expr) Expr { return e.binary("<<", other) }

// ShiftRight returns (e >> other).
func (e Expr) ShiftRight(other Expr) Expr { return e.binary(">>", other) }

// UnsignedShiftRight returns (e >>> other).
func (e Expr) UnsignedShiftRight(other Expr) Expr { return e.binary(">>>", other) }

// InstanceOf returns (e instanceof other).
func (e Expr) InstanceOf(other Expr) Expr { return e.binary("instanceof", other) }

// In returns (e in obj).
func (e Expr) In(obj Expr) Expr { return e.binary("in", obj) }

// Unary operators

type unaryOp struct {
	op     string
	expr   Expr
	prefix bool
}

func (u unaryOp) js(sb *strings.Builder) {
	if u.prefix {
		sb.WriteString(u.op)
		u.expr.js(sb)
	} else {
		u.expr.js(sb)
		sb.WriteString(u.op)
	}
}

func unary(op string, e Expr) Expr {
	return Expr{node: unaryOp{op: op, expr: e, prefix: true}}
}

// Not returns !e.
func (e Expr) Not() Expr { return unary("!", e) }

// Neg returns -e (negation).
func (e Expr) Neg() Expr { return unary("-", e) }

// Pos returns +e (unary plus).
func (e Expr) Pos() Expr { return unary("+", e) }

// BitwiseNot returns ~e.
func (e Expr) BitwiseNot() Expr { return unary("~", e) }

// Typeof returns (typeof e).
func (e Expr) Typeof() Expr { return unary("typeof ", e) }

// Void returns (void e).
func (e Expr) Void() Expr { return unary("void ", e) }

// Delete returns (delete e).
func (e Expr) Delete() Expr { return unary("delete ", e) }

// Ternary operator

type ternaryOp struct {
	cond    Expr
	ifTrue  Expr
	ifFalse Expr
}

func (t ternaryOp) js(sb *strings.Builder) {
	sb.WriteString("(")
	t.cond.js(sb)
	sb.WriteString(" ? ")
	t.ifTrue.js(sb)
	sb.WriteString(" : ")
	t.ifFalse.js(sb)
	sb.WriteString(")")
}

// Ternary returns (e ? ifTrue : ifFalse) using e as the condition.
func (e Expr) Ternary(ifTrue, ifFalse Expr) Expr {
	return Expr{node: ternaryOp{cond: e, ifTrue: ifTrue, ifFalse: ifFalse}}
}

// Grouping

type groupExpr struct {
	expr Expr
}

func (g groupExpr) js(sb *strings.Builder) {
	sb.WriteString("(")
	g.expr.js(sb)
	sb.WriteString(")")
}

// Group wraps the expression in parentheses.
func (e Expr) Group() Expr { return Expr{node: groupExpr{e}} }

// Comma expression

type commaExpr struct {
	exprs []Expr
}

func (c commaExpr) js(sb *strings.Builder) {
	sb.WriteString("(")
	for i, e := range c.exprs {
		if i > 0 {
			sb.WriteString(", ")
		}
		e.js(sb)
	}
	sb.WriteString(")")
}

// Comma creates a comma expression that evaluates all expressions
// and returns the value of the last one.
func Comma(exprs ...Expr) Expr {
	return Expr{node: commaExpr{exprs}}
}

// Spread operator

type spreadExpr struct {
	expr Expr
}

func (s spreadExpr) js(sb *strings.Builder) {
	sb.WriteString("...")
	s.expr.js(sb)
}

// Spread returns ...e.
func (e Expr) Spread() Expr { return Expr{node: spreadExpr{e}} }
