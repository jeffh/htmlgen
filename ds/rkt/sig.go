package rkt

import (
	"fmt"
	"strings"

	"github.com/jeffh/htmlgen/ds"
	"github.com/jeffh/htmlgen/js"
)

// Sig is a typed Rocket private signal name. It mirrors ds.Sig, but every
// reference is emitted with the "$$" prefix, which Rocket rewrites to the
// component's private signal path ($._rocket.<tag>.<id>.<name>). Use ds.Sig
// for ordinary global signals.
//
// A leading "$$" or "$" is accepted and stripped, so Sig("count"),
// Sig("$count") and Sig("$$count") are equivalent.
//
//	ds.Show(rkt.Sig("open").Value())     =>  data-show="$$open"
//	ds.OnClick(rkt.Sig("open").Toggle()) =>  data-on:click="($$open = !$$open)"
type Sig string

// Name returns the signal name without its leading "$$" or "$".
//
//	rkt.Sig("count").Name()    =>  "count"
//	rkt.Sig("$$count").Name()  =>  "count"
func (s Sig) Name() string {
	name := string(s)
	if strings.HasPrefix(name, "$$") {
		return name[2:]
	}
	return strings.TrimPrefix(name, "$")
}

// Ref returns the private signal reference expression. It panics on an empty
// signal name, which would otherwise silently emit a bare "$$". Names are
// trusted developer constants and are not otherwise validated.
//
//	rkt.Sig("count").Ref()  =>  $$count
func (s Sig) Ref() js.Expr {
	name := s.Name()
	if name == "" {
		panic("rkt: empty signal name")
	}
	return js.Raw("$$" + name)
}

// Value returns the private signal reference wrapped as a ds.Value, for use
// with attribute helpers such as ds.Show, ds.Text and ds.Class.
//
//	ds.Show(rkt.Sig("open").Value())  =>  data-show="$$open"
func (s Sig) Value() ds.Value { return ds.V(s.Ref()) }

// Not returns the negated private signal reference.
//
//	ds.Show(rkt.Sig("open").Not())  =>  data-show="!$$open"
func (s Sig) Not() ds.Value { return ds.V(s.Ref().Not()) }

// Set returns a ds.Value that assigns the private signal. The value is
// interpreted like ds.SetSignal: a js.Expr or ds.Value is emitted as-is,
// anything else is JSON encoded.
//
// Passing a js.Stmt panics: statements are not values, and JSON-encoding one
// would silently emit "{}". Use SetExpr with an expression, or ds.Do to
// sequence statements.
//
//	rkt.Sig("q").Set("hello")             =>  ($$q = "hello")
//	rkt.Sig("count").Set(0)               =>  ($$count = 0)
//	rkt.Sig("n").Set(rkt.Sig("m").Ref())  =>  ($$n = $$m)
func (s Sig) Set(value any) ds.Value {
	switch v := value.(type) {
	case js.Stmt:
		panic(fmt.Sprintf("rkt: Sig(%q).Set called with a js.Stmt (%q); statements are not values — use SetExpr with a js.Expr, or ds.Do to sequence statements", string(s), js.ToJSStmt(v)))
	case js.Expr:
		return s.SetExpr(v)
	case ds.Value:
		return s.SetExpr(v.Expr())
	default:
		return s.SetExpr(js.JSON(value))
	}
}

// SetExpr returns a ds.Value that assigns the private signal to a JavaScript
// expression. The assignment is parenthesized, like ds.SetSignalExpr, so it
// stays a valid expression when composed with And, Or or Ternary.
//
//	rkt.Sig("n").SetExpr(js.Int(1))  =>  ($$n = 1)
func (s Sig) SetExpr(e js.Expr) ds.Value {
	return ds.SetSignalExpr("$$"+s.Name(), e)
}

// Toggle returns a ds.Value that flips the private signal's truthiness.
//
//	rkt.Sig("open").Toggle()  =>  ($$open = !$$open)
func (s Sig) Toggle() ds.Value { return s.SetExpr(s.Ref().Not()) }

// Clear returns a ds.Value that assigns the empty string to the private
// signal. Single quotes are used so the generated attribute needs no
// HTML-escaped quotes.
//
//	rkt.Sig("q").Clear()  =>  ($$q = '')
func (s Sig) Clear() ds.Value { return s.SetExpr(js.Raw("''")) }

// Eq returns a strict-equality comparison against the private signal.
//
//	rkt.Sig("tab").Eq(js.String("plans"))  =>  ($$tab === "plans")
func (s Sig) Eq(e js.Expr) ds.Value { return ds.V(s.Ref().Eq(e)) }

// NotEq returns a strict-inequality comparison against the private signal.
//
//	rkt.Sig("tab").NotEq(js.String("plans"))  =>  ($$tab !== "plans")
func (s Sig) NotEq(e js.Expr) ds.Value { return ds.V(s.Ref().NotEq(e)) }

// Sub returns a derived private signal named "<name>_<suffix>", for widgets
// that keep a family of related signals. A leading "_" on suffix is ignored,
// so Sub("open") and Sub("_open") produce the same name.
//
//	rkt.Sig("plan").Sub("open")   =>  rkt.Sig("plan_open")
//	rkt.Sig("plan").Sub("_open")  =>  rkt.Sig("plan_open")
func (s Sig) Sub(suffix string) Sig {
	return Sig(s.Name() + "_" + strings.TrimPrefix(suffix, "_"))
}

// SignalRef creates a Rocket private signal reference: $$name. Any leading
// "$$" or "$" is stripped.
//
// It is equivalent to Sig(name).Value(); prefer Sig when you reference the
// same signal more than once.
//
//	rkt.SignalRef("items")  =>  $$items
func SignalRef(name string) ds.Value { return Sig(name).Value() }
