package js

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

// writeJSONString writes a JSON-encoded string directly to the builder,
// avoiding the allocation from json.Marshal.
func writeJSONString(sb *strings.Builder, s string) {
	sb.WriteByte('"')
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		switch r {
		case '"':
			sb.WriteString(`\"`)
		case '\\':
			sb.WriteString(`\\`)
		case '\n':
			sb.WriteString(`\n`)
		case '\r':
			sb.WriteString(`\r`)
		case '\t':
			sb.WriteString(`\t`)
		case '\b':
			sb.WriteString(`\b`)
		case '\f':
			sb.WriteString(`\f`)
		case '<':
			// Match json.Marshal HTML-safe escaping
			sb.WriteString("\\u003c")
		case '>':
			sb.WriteString("\\u003e")
		case '&':
			sb.WriteString("\\u0026")
		default:
			if r < 0x20 {
				sb.WriteString(`\u00`)
				sb.WriteByte("0123456789abcdef"[r>>4])
				sb.WriteByte("0123456789abcdef"[r&0xf])
			} else if r == utf8.RuneError && size == 1 {
				sb.WriteString("\\ufffd")
			} else {
				sb.WriteRune(r)
			}
		}
		i += size
	}
	sb.WriteByte('"')
}

// literal represents a JavaScript literal value emitted verbatim.
type literal struct {
	value string
}

func (l literal) js(sb *strings.Builder) { sb.WriteString(l.value) }

// stringLiteralNode represents a JavaScript string literal that escapes on output.
type stringLiteralNode struct {
	value string
}

func (s stringLiteralNode) js(sb *strings.Builder) { writeJSONString(sb, s.value) }

// String creates a JavaScript string literal, properly escaped using JSON encoding.
func String(s string) Expr { return Expr{node: stringLiteralNode{s}} }

// Int creates a JavaScript number literal from an integer.
func Int(n int) Expr { return Expr{node: literal{strconv.Itoa(n)}} }

// Int64 creates a JavaScript number literal from an int64.
func Int64(n int64) Expr { return Expr{node: literal{strconv.FormatInt(n, 10)}} }

// Float creates a JavaScript number literal from a float64.
func Float(f float64) Expr { return Expr{node: literal{strconv.FormatFloat(f, 'f', -1, 64)}} }

// Bool creates a JavaScript boolean literal.
func Bool(b bool) Expr {
	if b {
		return Expr{node: literal{"true"}}
	}
	return Expr{node: literal{"false"}}
}

// Null creates a JavaScript null literal.
func Null() Expr { return Expr{node: literal{"null"}} }

// Undefined creates a JavaScript undefined literal.
func Undefined() Expr { return Expr{node: literal{"undefined"}} }

// JSON creates a JavaScript value from a Go value using JSON encoding.
// Panics if the value cannot be marshaled to JSON.
func JSON(value any) Expr {
	b, err := json.Marshal(value)
	if err != nil {
		panic(fmt.Errorf("js.JSON: %w: value=%#v", err, value))
	}
	return Expr{node: literal{string(b)}}
}

// arrayLiteral represents a JavaScript array literal.
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

// Array creates a JavaScript array literal from expressions.
func Array(elements ...Expr) Expr { return Expr{node: arrayLiteral{elements}} }

// KV represents a key-value pair for object literals.
type KV struct {
	Key   string
	Value Expr
}

// Pair creates a key-value pair for Object().
func Pair(key string, value Expr) KV { return KV{Key: key, Value: value} }

// objectLiteral represents a JavaScript object literal.
type objectLiteral struct {
	pairs []KV
}

func (o objectLiteral) js(sb *strings.Builder) {
	sb.WriteString("{")
	for i, kv := range o.pairs {
		if i > 0 {
			sb.WriteString(", ")
		}
		writeJSONString(sb, kv.Key)
		sb.WriteString(": ")
		kv.Value.js(sb)
	}
	sb.WriteString("}")
}

// Object creates a JavaScript object literal from key-value pairs.
func Object(pairs ...KV) Expr { return Expr{node: objectLiteral{pairs}} }

// regexLiteral represents a JavaScript regular-expression literal.
type regexLiteral struct {
	pattern string
	flags   string
}

func (r regexLiteral) js(sb *strings.Builder) {
	sb.WriteByte('/')
	sb.WriteString(r.pattern)
	sb.WriteByte('/')
	sb.WriteString(r.flags)
}

// Regex creates a JavaScript regular-expression literal: /pattern/flags.
//
// The pattern and flags are emitted VERBATIM — no escaping or validation is
// performed. The caller is responsible for supplying a valid JS regex literal
// body; in particular any literal "/" inside pattern must already be escaped
// as "\\/", otherwise the emitted literal terminates early and produces
// invalid JavaScript.
//
//	Regex("^user_", "i")   =>  /^user_/i
//	Regex("\\d+", "")      =>  /\d+/
//	Regex("a\\/b", "g")    =>  /a\/b/g
func Regex(pattern, flags string) Expr {
	return Expr{node: regexLiteral{pattern: pattern, flags: flags}}
}

// identifier represents a bare identifier reference.
type identifier string

func (i identifier) js(sb *strings.Builder) { sb.WriteString(string(i)) }

// Ident creates a JavaScript identifier reference.
// This should be used for variable names, not for string literals.
func Ident(name string) Expr { return Expr{node: identifier(name)} }

// This creates the special "this" identifier.
func This() Expr { return Expr{node: identifier("this")} }
