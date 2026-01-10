package js

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// literal represents a JavaScript literal value.
type literal struct {
	value string
}

func (l literal) js(sb *strings.Builder) { sb.WriteString(l.value) }
func (l literal) callable()              {}

// String creates a JavaScript string literal, properly escaped using JSON encoding.
func String(s string) Callable {
	b, _ := json.Marshal(s)
	return literal{string(b)}
}

// Int creates a JavaScript number literal from an integer.
func Int(n int) Callable {
	return literal{strconv.Itoa(n)}
}

// Int64 creates a JavaScript number literal from an int64.
func Int64(n int64) Callable {
	return literal{strconv.FormatInt(n, 10)}
}

// Float creates a JavaScript number literal from a float64.
func Float(f float64) Callable {
	return literal{strconv.FormatFloat(f, 'f', -1, 64)}
}

// Bool creates a JavaScript boolean literal.
func Bool(b bool) Callable {
	if b {
		return literal{"true"}
	}
	return literal{"false"}
}

// Null creates a JavaScript null literal.
func Null() Callable {
	return literal{"null"}
}

// Undefined creates a JavaScript undefined literal.
func Undefined() Callable {
	return literal{"undefined"}
}

// JSON creates a JavaScript value from a Go value using JSON encoding.
// Panics if the value cannot be marshaled to JSON.
func JSON(value any) Callable {
	b, err := json.Marshal(value)
	if err != nil {
		panic(fmt.Errorf("js.JSON: %w: value=%#v", err, value))
	}
	return literal{string(b)}
}

// Array creates a JavaScript array literal from expressions.
func Array(elements ...Expr) Callable {
	return arrayLiteral{elements}
}

type arrayLiteral struct {
	elements []Expr
}

func (a arrayLiteral) js(sb *strings.Builder) {
	sb.WriteString("[")
	for i, el := range a.elements {
		if i > 0 {
			sb.WriteString(", ")
		}
		el.js(sb)
	}
	sb.WriteString("]")
}
func (a arrayLiteral) callable() {}

// Object creates a JavaScript object literal from key-value pairs.
func Object(pairs ...KV) Callable {
	return objectLiteral{pairs}
}

// KV represents a key-value pair for object literals.
type KV struct {
	Key   string
	Value Expr
}

// Pair creates a key-value pair for Object().
func Pair(key string, value Expr) KV {
	return KV{Key: key, Value: value}
}

type objectLiteral struct {
	pairs []KV
}

func (o objectLiteral) js(sb *strings.Builder) {
	sb.WriteString("{")
	for i, kv := range o.pairs {
		if i > 0 {
			sb.WriteString(", ")
		}
		// Quote the key using JSON encoding for safety
		b, _ := json.Marshal(kv.Key)
		sb.WriteString(string(b))
		sb.WriteString(": ")
		kv.Value.js(sb)
	}
	sb.WriteString("}")
}
func (o objectLiteral) callable() {}

// Ident creates a JavaScript identifier reference.
// This should be used for variable names, not for string literals.
func Ident(name string) Callable {
	return identifier(name)
}

type identifier string

func (i identifier) js(sb *strings.Builder) { sb.WriteString(string(i)) }
func (i identifier) callable()              {}

// This creates the special "this" identifier.
func This() Callable {
	return identifier("this")
}
