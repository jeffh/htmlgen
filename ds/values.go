package ds

import (
	"github.com/jeffh/htmlgen/js"
)

// ConsoleLog returns a Value that calls console.log on the given expressions.
func ConsoleLog(values ...js.Expr) Value {
	return Value{expr: js.ConsoleLog(values...)}
}

// Confirm guards actions behind a native confirm() dialog: the then values only
// run when the user accepts. The message is emitted as a properly escaped
// JavaScript string literal. Multiple then values are folded left with && using
// the same semantics as js.And.
//
//	ds.Confirm("Delete this plan?", ds.Delete("/plans/1"))
//	=>  (confirm("Delete this plan?") && @delete("/plans/1"))
//
//	ds.Confirm("Sure?")  =>  confirm("Sure?")
func Confirm(message string, then ...Value) Value {
	guard := js.Confirm(js.String(message))
	if len(then) == 0 {
		return Value{expr: guard}
	}
	exprs := make([]js.Expr, 0, 1+len(then))
	exprs = append(exprs, guard)
	for _, t := range then {
		exprs = append(exprs, t.expr)
	}
	return Value{expr: js.And(exprs...)}
}
