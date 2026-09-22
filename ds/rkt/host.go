package rkt

import (
	"encoding"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"maps"
	"math"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/jeffh/htmlgen/h"
)

// Version is the Datastar release whose Rocket bundle this package targets.
// Rocket is beta upstream (the v1.0.4 bundle is labeled "Rocket beta.2"), so
// pin the version you test against.
const Version = "1.0.4"

// BundleURL is the CDN URL of the Rocket bundle. It is a superset of
// datastar.js: load it instead of datastar.js, not alongside it.
const BundleURL = "https://cdn.jsdelivr.net/gh/starfederation/datastar@v" + Version + "/bundles/datastar-rocket.js"

// AliasedBundleURL is the CDN URL of the CSP-safe aliased Rocket bundle.
const AliasedBundleURL = "https://cdn.jsdelivr.net/gh/starfederation/datastar@v" + Version + "/bundles/datastar-rocket-aliased.js"

// dateLayout matches JavaScript's Date.prototype.toISOString: UTC, with
// millisecond precision and a "Z" zone.
const dateLayout = "2006-01-02T15:04:05.000Z07:00"

// customElementRe matches a custom element tag name: lowercase, starting with
// a letter and containing at least one hyphen.
var customElementRe = regexp.MustCompile(`^[a-z][a-z0-9._-]*-[a-z0-9._-]*$`)

// ScriptAttrs returns the attributes of the Rocket bundle script tag:
// type="module" src=BundleURL.
func ScriptAttrs() h.Attributes {
	return h.Attrs("type", "module", "src", BundleURL)
}

// Script writes the Rocket bundle script tag. Load it instead of datastar.js.
//
//	<script type="module" src="https://cdn.jsdelivr.net/.../datastar-rocket.js"></script>
func Script(b *h.B) {
	b.Script(ScriptAttrs(), nil)
}

// AliasedScript writes the script tag for the CSP-safe aliased Rocket bundle.
func AliasedScript(b *h.B) {
	b.Script(h.Attrs("type", "module", "src", AliasedBundleURL), nil)
}

// kebab converts a camelCase or PascalCase name to kebab-case, which is how
// Rocket reads prop attributes off a host element. A name that is already
// kebab-case (or a single lowercase word) is returned unchanged.
//
//	kebab("maxItems")  =>  "max-items"
//	kebab("MaxItems")  =>  "max-items"
//	kebab("max-items") =>  "max-items"
func kebab(name string) string {
	var sb strings.Builder
	sb.Grow(len(name) + 4)
	bs := []byte(name)
	for i := 0; i < len(bs); i++ {
		c := bs[i]
		if c < 'A' || c > 'Z' {
			sb.WriteByte(c)
			continue
		}
		if i > 0 && bs[i-1] != '-' {
			prevLower := bs[i-1] >= 'a' && bs[i-1] <= 'z' || bs[i-1] >= '0' && bs[i-1] <= '9'
			nextLower := i+1 < len(bs) && bs[i+1] >= 'a' && bs[i+1] <= 'z'
			if prevLower || nextLower {
				sb.WriteByte('-')
			}
		}
		sb.WriteByte(c - 'A' + 'a')
	}
	return sb.String()
}

// encodeProp renders a Go value using Rocket's prop codecs. See Prop.
func encodeProp(name string, value any) string {
	switch v := value.(type) {
	case nil:
		return ""
	case bool:
		return strconv.FormatBool(v)
	case string:
		return v
	case time.Time:
		return v.UTC().Format(dateLayout)
	case []byte:
		return base64.StdEncoding.EncodeToString(v)
	case int:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case uintptr:
		return strconv.FormatUint(uint64(v), 10)
	case float32:
		return formatFloat(name, float64(v), 32)
	case float64:
		return formatFloat(name, v, 64)
	case fmt.Stringer:
		return v.String()
	case encoding.TextMarshaler:
		text, err := v.MarshalText()
		if err != nil {
			panic(fmt.Sprintf("rkt: Prop(%q): MarshalText failed: %v", name, err))
		}
		return string(text)
	default:
		return marshalJSON(name, value)
	}
}

// formatFloat renders a float as a decimal string, panicking on NaN and Inf,
// which have no decimal form a number prop can carry.
func formatFloat(name string, f float64, bits int) string {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		panic(fmt.Sprintf("rkt: Prop(%q): %v is not representable as a number prop", name, f))
	}
	return strconv.FormatFloat(f, 'f', -1, bits)
}

// marshalJSON JSON-encodes a value, panicking on an encoding error.
func marshalJSON(name string, value any) string {
	data, err := json.Marshal(value)
	if err != nil {
		panic(fmt.Sprintf("rkt: Prop(%q): JSON encoding failed: %v", name, err))
	}
	return string(data)
}

// Prop sets a Rocket component prop on a host element. The attribute name is
// converted to kebab-case, because that is how Rocket reads props off the
// host element, and is then validated like any other attribute name.
//
// The value is encoded with the codec matching its Go type:
//
//   - bool            => "true" / "false"
//   - string          => as-is
//   - integers        => decimal
//   - floats          => decimal (NaN and Inf panic)
//   - time.Time       => UTC ISO string, like JavaScript's toISOString
//   - []byte          => standard base64
//   - nil             => "" (an empty attribute)
//   - fmt.Stringer    => String()
//   - TextMarshaler   => MarshalText()
//   - anything else   => JSON (a JSON encoding error panics)
//
// Use the typed encoders (BoolProp, NumberProp, JSONProp, DateProp, BinProp)
// when the codec must not follow from the Go type — JSONProp on a string, for
// example, emits a quoted JSON string rather than the bare text.
//
//	rkt.Prop("maxItems", 5)        =>  max-items="5"
//	rkt.Prop("label", "Count")     =>  label="Count"
func Prop(name string, value any) h.Attribute {
	return h.Attr(kebab(name), encodeProp(name, value))
}

// JSONProp sets a prop with Rocket's json codec, always JSON-encoding the
// value. Unlike Prop, a string becomes a quoted JSON string.
//
//	rkt.JSONProp("items", []int{1, 2})  =>  items="[1,2]"
func JSONProp(name string, v any) h.Attribute {
	return h.Attr(kebab(name), marshalJSON(name, v))
}

// DateProp sets a prop with Rocket's date codec: the UTC ISO string
// JavaScript's Date.prototype.toISOString produces.
//
//	rkt.DateProp("startsAt", t)  =>  starts-at="2026-01-02T03:04:05.000Z"
func DateProp(name string, t time.Time) h.Attribute {
	return h.Attr(kebab(name), t.UTC().Format(dateLayout))
}

// BinProp sets a prop with Rocket's bin codec: standard-alphabet base64.
//
//	rkt.BinProp("blob", []byte("hi"))  =>  blob="aGk="
func BinProp(name string, data []byte) h.Attribute {
	return h.Attr(kebab(name), base64.StdEncoding.EncodeToString(data))
}

// BoolProp sets a prop with Rocket's bool codec: "true" or "false". Note that
// this always writes a value — it is not an HTML boolean attribute.
//
//	rkt.BoolProp("isOpen", true)  =>  is-open="true"
func BoolProp(name string, b bool) h.Attribute {
	return h.Attr(kebab(name), strconv.FormatBool(b))
}

// NumberProp sets a prop with Rocket's number codec: a decimal string. NaN
// and Inf panic.
//
//	rkt.NumberProp("step", 0.5)  =>  step="0.5"
func NumberProp(name string, n float64) h.Attribute {
	return h.Attr(kebab(name), formatFloat(name, n, 64))
}

// Props sets several props at once, sorted by key for deterministic output.
// Each value is encoded by Prop.
//
//	rkt.Props(map[string]any{"step": 2, "label": "Count"})
//	=>  label="Count" step="2"
func Props(props map[string]any) h.Attributes {
	attrs := make(h.Attributes, 0, len(props))
	for _, k := range slices.Sorted(maps.Keys(props)) {
		attrs = append(attrs, Prop(k, props[k]))
	}
	if len(attrs) == 0 {
		return nil
	}
	return attrs
}

// Component writes a Rocket component host element: the custom element whose
// tag was registered with rocket('<tag>', …). Props go in attrs (see Prop and
// Props), and body renders the element's light-DOM children.
//
// It panics unless tag is a valid custom element name: lowercase, starting
// with a letter and containing at least one hyphen.
//
//	rkt.Component(b, "counter-widget", h.AttrsOf(rkt.Prop("step", 2)), nil)
//	=>  <counter-widget step="2"></counter-widget>
func Component(b *h.B, tag string, attrs h.Attributes, body h.Body) {
	if !customElementRe.MatchString(tag) {
		panic("rkt: invalid custom element name " + strconv.Quote(tag) + "; must be lowercase and contain a hyphen")
	}
	b.El(tag, attrs, body)
}
