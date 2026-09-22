package rkt

import (
	"regexp"
	"strconv"

	"github.com/jeffh/htmlgen/ds"
	"github.com/jeffh/htmlgen/js"
)

// identRe matches a JavaScript identifier, the shape Rocket requires for
// local action names and for data-for loop variables.
var identRe = regexp.MustCompile(`^[A-Za-z_$][A-Za-z0-9_$]*$`)

// validateIdent panics unless name is a JavaScript identifier. Names are
// written into generated JavaScript unescaped and must never come from
// untrusted input.
func validateIdent(kind, name string) {
	if !identRe.MatchString(name) {
		panic("rkt: invalid " + kind + " " + strconv.Quote(name))
	}
}

// Action calls a component-local action registered through Rocket's
// setup({action}) callback: @name(args...). Inside a component's template
// Rocket rewrites the call to @dispatchRocket("name", args...), which falls
// back to the global action registry when no owning component defines the
// action.
//
// Use Dispatch outside a component's own template, where the rewrite does not
// happen.
//
// Action panics unless name is a JavaScript identifier.
//
//	ds.OnClick(rkt.Action("increment"))            =>  data-on:click="@increment()"
//	ds.OnClick(rkt.Action("add", js.Int(1)))       =>  data-on:click="@add(1)"
func Action(name string, args ...js.Expr) ds.Value {
	validateIdent("action name", name)
	return ds.V(ds.DatastarAction(name, args...))
}

// Dispatch calls a Rocket action by name through the explicit dispatcher:
// @dispatchRocket("name", args...). This is the form Rocket rewrites @name(…)
// into, written out by hand for markup that lives outside a component's own
// template (where no rewrite happens).
//
// Dispatch panics unless name is a JavaScript identifier.
//
//	rkt.Dispatch("increment", js.Int(1))  =>  @dispatchRocket("increment", 1)
func Dispatch(name string, args ...js.Expr) ds.Value {
	validateIdent("action name", name)
	all := make([]js.Expr, 0, len(args)+1)
	all = append(all, js.String(name))
	all = append(all, args...)
	return ds.V(ds.DatastarAction("dispatchRocket", all...))
}
