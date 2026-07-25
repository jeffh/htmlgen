package ds

import (
	"fmt"

	"github.com/jeffh/htmlgen/js"
)

// Navigate generates JavaScript that navigates to a URL by assigning to
// window.location.href.
//
//	Navigate("/users/%d", 42)  =>  window.location.href = "/users/42"
func Navigate(path string, values ...any) string {
	return js.ToJSStmt(js.Window.Prop("location").Prop("href").Assign(
		js.String(fmt.Sprintf(path, values...)),
	))
}

// ConsoleLog returns a Value that calls console.log on the given expressions.
func ConsoleLog(values ...js.Expr) Value {
	return Value{expr: js.ConsoleLog(values...)}
}

// And combines multiple expressions with && and returns the result.
//
//	And(Raw("$a"), Raw("$b"))  =>  ($a && $b)
func And(actions ...js.Expr) js.Expr {
	if len(actions) == 0 {
		return js.Bool(true)
	}
	result := actions[0]
	for i := 1; i < len(actions); i++ {
		result = result.And(actions[i])
	}
	return result
}

// AndValue combines multiple expressions with && and wraps the result as a Value.
func AndValue(actions ...js.Expr) Value {
	return Value{expr: And(actions...)}
}

// Confirm guards actions behind a native confirm() dialog: the then values only
// run when the user accepts. The message is emitted as a properly escaped
// JavaScript string literal. Multiple then values are folded left with && using
// the same semantics as And.
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
	return Value{expr: And(exprs...)}
}
