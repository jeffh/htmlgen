package ds

import (
	"fmt"
	"strings"

	"github.com/jeffh/htmlgen/js"
)

// Sig is a typed Datastar signal name. It removes the need to repeat the "$"
// prefix and string-concatenated signal names across a view, and its methods
// produce Values usable anywhere a Datastar expression is accepted (Show,
// Text, OnClick, Init, ...).
//
// A leading "$" is accepted and stripped, so Sig("open") and Sig("$open") are
// equivalent.
//
//	ds.Show(ds.Sig("open").Value())     =>  data-show="$open"
//	ds.OnClick(ds.Sig("open").Toggle()) =>  data-on:click="($open = !$open)"
type Sig string

// Name returns the signal name without the leading "$".
//
//	ds.Sig("open").Name()   =>  "open"
//	ds.Sig("$open").Name()  =>  "open"
func (s Sig) Name() string { return strings.TrimPrefix(string(s), "$") }

// Ref returns the signal reference expression. It panics on an empty signal
// name, which would otherwise silently emit a bare "$". Names are trusted
// developer constants and are not otherwise validated.
//
//	ds.Sig("open").Ref()  =>  $open
func (s Sig) Ref() js.Expr {
	name := s.Name()
	if name == "" {
		panic("ds: empty signal name")
	}
	return js.Raw("$" + name)
}

// Value returns the signal reference wrapped as a Value, for use with
// attribute helpers such as Show, Text and Class.
//
//	ds.Show(ds.Sig("open").Value())  =>  data-show="$open"
func (s Sig) Value() Value { return Value{expr: s.Ref()} }

// Not returns the negated signal reference.
//
//	ds.Show(ds.Sig("open").Not())  =>  data-show="!$open"
func (s Sig) Not() Value { return Value{expr: s.Ref().Not()} }

// Set returns a Value that assigns the signal. The value is interpreted like
// SetSignal: a js.Expr or Value is emitted as-is, anything else is JSON
// encoded.
//
// Passing a js.Stmt panics: statements are not values, and JSON-encoding one
// would silently emit "{}". Use SetExpr with an expression, or Do to sequence
// statements.
//
//	ds.Sig("q").Set("hello")            =>  ($q = "hello")
//	ds.Sig("n").Set(3)                  =>  ($n = 3)
//	ds.Sig("n").Set(ds.Sig("m").Ref())  =>  ($n = $m)
func (s Sig) Set(value any) Value {
	if stmt, ok := value.(js.Stmt); ok {
		panic(fmt.Sprintf("ds: Sig(%q).Set called with a js.Stmt (%q); statements are not values — use SetExpr with a js.Expr, or ds.Do to sequence statements", string(s), js.ToJSStmt(stmt)))
	}
	return SetSignal(s.Name(), value)
}

// SetExpr returns a Value that assigns the signal to a JavaScript expression.
//
//	ds.Sig("n").SetExpr(js.Int(1))  =>  ($n = 1)
func (s Sig) SetExpr(e js.Expr) Value { return SetSignalExpr(s.Name(), e) }

// Toggle returns a Value that flips the signal's truthiness.
//
//	ds.Sig("open").Toggle()  =>  ($open = !$open)
func (s Sig) Toggle() Value { return s.SetExpr(s.Ref().Not()) }

// Clear returns a Value that assigns the empty string to the signal. Single
// quotes are used so the generated attribute needs no HTML-escaped quotes.
//
//	ds.Sig("q").Clear()  =>  ($q = '')
func (s Sig) Clear() Value { return s.SetExpr(js.Raw("''")) }

// Eq returns a strict-equality comparison against the signal.
//
//	ds.Sig("tab").Eq(js.String("plans"))  =>  ($tab === "plans")
func (s Sig) Eq(e js.Expr) Value { return Value{expr: s.Ref().Eq(e)} }

// NotEq returns a strict-inequality comparison against the signal.
//
//	ds.Sig("tab").NotEq(js.String("plans"))  =>  ($tab !== "plans")
func (s Sig) NotEq(e js.Expr) Value { return Value{expr: s.Ref().NotEq(e)} }

// Sub returns a derived signal named "<name>_<suffix>". It is intended for
// widgets that keep a family of related signals (a combobox with query, open,
// active and above sub-signals, for example). A leading "_" on suffix is
// ignored so both Sub("open") and Sub("_open") produce the same name.
//
//	ds.Sig("plan").Sub("open")   =>  ds.Sig("plan_open")
//	ds.Sig("plan").Sub("_open")  =>  ds.Sig("plan_open")
func (s Sig) Sub(suffix string) Sig {
	return Sig(s.Name() + "_" + strings.TrimPrefix(suffix, "_"))
}
