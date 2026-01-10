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
func (b binaryOp) callable() {}

// Add returns left + right
func Add(left, right Expr) Callable { return binaryOp{left, "+", right} }

// Sub returns left - right
func Sub(left, right Expr) Callable { return binaryOp{left, "-", right} }

// Mul returns left * right
func Mul(left, right Expr) Callable { return binaryOp{left, "*", right} }

// Div returns left / right
func Div(left, right Expr) Callable { return binaryOp{left, "/", right} }

// Mod returns left % right
func Mod(left, right Expr) Callable { return binaryOp{left, "%", right} }

// Eq returns left === right (strict equality)
func Eq(left, right Expr) Callable { return binaryOp{left, "===", right} }

// NotEq returns left !== right (strict inequality)
func NotEq(left, right Expr) Callable { return binaryOp{left, "!==", right} }

// LooseEq returns left == right (loose equality)
func LooseEq(left, right Expr) Callable { return binaryOp{left, "==", right} }

// LooseNotEq returns left != right (loose inequality)
func LooseNotEq(left, right Expr) Callable { return binaryOp{left, "!=", right} }

// Lt returns left < right
func Lt(left, right Expr) Callable { return binaryOp{left, "<", right} }

// LtEq returns left <= right
func LtEq(left, right Expr) Callable { return binaryOp{left, "<=", right} }

// Gt returns left > right
func Gt(left, right Expr) Callable { return binaryOp{left, ">", right} }

// GtEq returns left >= right
func GtEq(left, right Expr) Callable { return binaryOp{left, ">=", right} }

// And returns left && right
func And(left, right Expr) Callable { return binaryOp{left, "&&", right} }

// Or returns left || right
func Or(left, right Expr) Callable { return binaryOp{left, "||", right} }

// NullishCoalesce returns left ?? right
func NullishCoalesce(left, right Expr) Callable { return binaryOp{left, "??", right} }

// BitwiseAnd returns left & right
func BitwiseAnd(left, right Expr) Callable { return binaryOp{left, "&", right} }

// BitwiseOr returns left | right
func BitwiseOr(left, right Expr) Callable { return binaryOp{left, "|", right} }

// BitwiseXor returns left ^ right
func BitwiseXor(left, right Expr) Callable { return binaryOp{left, "^", right} }

// ShiftLeft returns left << right
func ShiftLeft(left, right Expr) Callable { return binaryOp{left, "<<", right} }

// ShiftRight returns left >> right
func ShiftRight(left, right Expr) Callable { return binaryOp{left, ">>", right} }

// UnsignedShiftRight returns left >>> right
func UnsignedShiftRight(left, right Expr) Callable { return binaryOp{left, ">>>", right} }

// Instanceof returns left instanceof right
func Instanceof(left, right Expr) Callable { return binaryOp{left, "instanceof", right} }

// In returns left in right
func In(left, right Expr) Callable { return binaryOp{left, "in", right} }

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
func (u unaryOp) callable() {}

// Not returns !expr
func Not(expr Expr) Callable { return unaryOp{"!", expr, true} }

// Neg returns -expr (negation)
func Neg(expr Expr) Callable { return unaryOp{"-", expr, true} }

// Pos returns +expr (unary plus)
func Pos(expr Expr) Callable { return unaryOp{"+", expr, true} }

// BitwiseNot returns ~expr
func BitwiseNot(expr Expr) Callable { return unaryOp{"~", expr, true} }

// Typeof returns typeof expr
func Typeof(expr Expr) Callable { return unaryOp{"typeof ", expr, true} }

// Void returns void expr
func Void(expr Expr) Callable { return unaryOp{"void ", expr, true} }

// Delete returns delete expr
func Delete(expr Expr) Callable { return unaryOp{"delete ", expr, true} }

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
func (t ternaryOp) callable() {}

// Ternary returns cond ? ifTrue : ifFalse
func Ternary(cond, ifTrue, ifFalse Expr) Callable {
	return ternaryOp{cond, ifTrue, ifFalse}
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
func (g groupExpr) callable() {}

// Group wraps an expression in parentheses.
func Group(expr Expr) Callable {
	return groupExpr{expr}
}

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
func (c commaExpr) callable() {}

// Comma creates a comma expression that evaluates all expressions
// and returns the value of the last one.
func Comma(exprs ...Expr) Callable {
	return commaExpr{exprs}
}

// Spread operator

type spreadExpr struct {
	expr Expr
}

func (s spreadExpr) js(sb *strings.Builder) {
	sb.WriteString("...")
	s.expr.js(sb)
}
func (s spreadExpr) callable() {}

// Spread creates a spread expression: ...expr
func Spread(expr Expr) Callable {
	return spreadExpr{expr}
}
