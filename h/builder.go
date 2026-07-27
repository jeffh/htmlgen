package h

import "fmt"

// Body is an element body that streams child content through b.
type Body = func(*B)

// parseArgs separates attributes from the optional body in an element's
// argument list. Later attributes and bodies replace earlier ones. Caller
// attribute slices are never mutated: the first Attributes arg is borrowed
// as-is, and a private copy is made before any merge.
func parseArgs(name string, args []any) (Attributes, Body) {
	var attrs Attributes
	var body Body
	owned := false // attrs is a private copy safe to mutate

	ensureOwned := func() {
		if !owned {
			attrs = append(Attributes(nil), attrs...)
			owned = true
		}
	}

	for _, arg := range args {
		if arg == nil {
			continue
		}
		switch value := arg.(type) {
		case Attributes:
			if attrs == nil {
				attrs = value
			} else {
				ensureOwned()
				attrs.Merge(value)
			}
		case Attribute:
			// Skip zero attributes (e.g., from AttrIf when the condition is false).
			if value.Name == "" {
				continue
			}
			if attrs == nil {
				attrs = Attributes{value}
				owned = true
			} else {
				ensureOwned()
				attrs.set(value.Name, value.Value)
			}
		case AttrBuilder:
			attr := value.Attribute()
			if attr.Name == "" {
				continue
			}
			if attrs == nil {
				attrs = Attributes{attr}
				owned = true
			} else {
				ensureOwned()
				attrs.set(attr.Name, attr.Value)
			}
		case func(*B):
			// A typed-nil body is ignored, like an untyped nil arg.
			if value != nil {
				body = value
			}
		default:
			panic(fmt.Sprintf("htmlgen: unsupported argument type %T for <%s>", arg, name))
		}
	}
	return attrs, body
}
