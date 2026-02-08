package ds

import (
	"fmt"

	"github.com/jeffh/htmlgen/js"
)

// Navigate generates JavaScript to navigate to a URL.
// Example: Navigate("/users/%d", 42) produces: window.location.href = "/users/42"
func Navigate(path string, values ...any) string {
	return js.ToJSStmt(js.Assign(
		js.Prop(js.Prop(js.Window, "location"), "href"),
		js.String(fmt.Sprintf(path, values...)),
	))
}

// ConsoleLog creates an AttrMutator that logs values to the console.
// Example: ConsoleLog(Signal("value"), Str("clicked"))
func ConsoleLog(values ...js.Expr) AttrMutator {
	return AttrFunc(func(attr *attrBuilder) {
		attr.AppendStatement(js.ToJS(js.ConsoleLog(values...)))
	})
}

// appendName creates an AttrMutator that appends to the attribute name.
func appendName(name string) AttrMutator {
	return AttrFunc(func(attr *attrBuilder) {
		attr.name.WriteString(name)
	})
}

// And creates a JavaScript && expression from multiple expressions.
// Example: And(Raw("$a"), Raw("$b")) produces: $a && $b
func And(actions ...js.Expr) js.Callable {
	if len(actions) == 0 {
		return js.Bool(true)
	}
	result := actions[0]
	for i := 1; i < len(actions); i++ {
		result = js.And(result, actions[i])
	}
	return result.(js.Callable)
}

// AndMutator creates an AttrMutator that combines expressions with &&.
// Example: AndMutator(Raw("$a"), Raw("$b")) produces an AttrMutator.
func AndMutator(actions ...js.Expr) AttrMutator {
	return AttrFunc(func(attr *attrBuilder) {
		attr.AppendStatement(js.ToJS(And(actions...)))
	})
}
