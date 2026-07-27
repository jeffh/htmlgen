package h

import (
	"maps"
	"slices"
	"sort"
	"strconv"
)

// Attribute represents a single HTML attribute as a name-value pair.
//
// Name is written to output verbatim, without escaping or validation. Values
// constructed through Attr, Attrs, AttrsMap, AttrIf, Set, or SetDefault are
// validated; an Attribute built directly as a struct literal is trusted, and
// its Name must not be derived from untrusted input.
type Attribute struct {
	Name  string
	Value string
}

// validName reports whether name is a valid attribute or element name:
// an ASCII letter followed by ASCII letters, digits, '_', '.', ':', or '-'.
func validName(name string) bool {
	if name == "" {
		return false
	}
	c := name[0]
	if !('a' <= c && c <= 'z' || 'A' <= c && c <= 'Z') {
		return false
	}
	for i := 1; i < len(name); i++ {
		switch c := name[i]; {
		case 'a' <= c && c <= 'z', 'A' <= c && c <= 'Z', '0' <= c && c <= '9',
			c == '_', c == '.', c == ':', c == '-':
		default:
			return false
		}
	}
	return true
}

func validateAttrName(name string) {
	if !validName(name) {
		panic("htmlgen: invalid attribute name " + strconv.Quote(name))
	}
}

func validateTagName(name string) {
	if !validName(name) {
		panic("htmlgen: invalid element name " + strconv.Quote(name))
	}
}

// AttrBuilder produces a single Attribute on demand.
// Fluent attribute builders (in companion packages like ds and hx) implement
// this interface so that they can be passed directly to element methods without
// an explicit terminator method.
type AttrBuilder interface {
	Attribute() Attribute
}

// Attr creates a new Attribute with the given name and value.
// Panics if name is not an ASCII letter followed by ASCII letters, digits,
// '_', '.', ':', or '-'; attribute names must never come from untrusted input.
func Attr(name, value string) Attribute {
	validateAttrName(name)
	return Attribute{Name: name, Value: value}
}

// AttrIf returns an Attribute if cond is true, otherwise returns a zero Attribute
// which will be ignored during rendering. This is useful for conditionally
// including attributes:
//
//	b.Button(
//	    h.AttrIf(isDisabled, "disabled", ""),
//	    h.AttrIf(isPrimary, "class", "btn-primary"),
//	    func(b *h.B) { b.Text("Submit") },
//	)
//
// Panics if cond is true and name is not a valid attribute name (see Attr).
func AttrIf(cond bool, name, value string) Attribute {
	if cond {
		validateAttrName(name)
		return Attribute{Name: name, Value: value}
	}
	return Attribute{}
}

// Attributes is a slice of Attribute values representing HTML element attributes.
// It provides methods for getting, setting, and deleting attributes by name.
type Attributes []Attribute

// Attrs creates an Attributes slice from alternating key-value string pairs.
// Panics if an odd number of arguments is provided or if any key is not a
// valid attribute name (see Attr).
//
// Example: Attrs("href", "/home", "class", "nav-link")
func Attrs(kv ...string) Attributes {
	if len(kv)%2 != 0 {
		panic("Attrs(...) expects an even number of arguments")
	}
	if len(kv) == 0 {
		return nil
	}
	results := make(Attributes, 0, len(kv)/2)
	for i := 0; i < len(kv); i += 2 {
		validateAttrName(kv[i])
		results = append(results, Attribute{Name: kv[i], Value: kv[i+1]})
	}
	return results
}

// AttrsMap creates an Attributes slice from a map of key-value pairs.
// Keys are sorted alphabetically for deterministic output.
// Panics if any key is not a valid attribute name (see Attr).
func AttrsMap(m map[string]string) Attributes {
	result := make(Attributes, 0, len(m))
	keys := slices.Collect(maps.Keys(m))
	sort.Strings(keys)
	for _, k := range keys {
		validateAttrName(k)
		result = append(result, Attribute{k, m[k]})
	}
	return result
}

// Get returns the value for the given attribute name and true if found,
// or an empty string and false if not found.
func (a *Attributes) Get(key string) (string, bool) {
	for _, attr := range *a {
		if attr.Name == key {
			return attr.Value, true
		}
	}
	return "", false
}

// Index returns the index of the attribute with the given name,
// or -1 if not found.
func (a *Attributes) Index(key string) int {
	for i, attr := range *a {
		if attr.Name == key {
			return i
		}
	}
	return -1
}

// Set sets the value for the given attribute name.
// If the attribute already exists, its value is updated; otherwise,
// a new attribute is appended.
// Panics if key is not a valid attribute name (see Attr).
func (a *Attributes) Set(key, value string) {
	validateAttrName(key)
	a.set(key, value)
}

// set is Set without name validation, for internal paths that merge
// already-constructed Attribute values.
func (a *Attributes) set(key, value string) {
	idx := a.Index(key)
	if idx >= 0 {
		(*a)[idx].Value = value
	} else {
		*a = append(*a, Attribute{Name: key, Value: value})
	}
}

// SetDefault sets the value for the given attribute name only if
// the attribute does not already exist.
// Panics if key is not a valid attribute name (see Attr).
func (a *Attributes) SetDefault(key, value string) {
	validateAttrName(key)
	idx := a.Index(key)
	if idx < 0 {
		*a = append(*a, Attribute{Name: key, Value: value})
	}
}

// Delete removes the attribute with the given name if it exists.
func (a *Attributes) Delete(key string) {
	idx := a.Index(key)
	if idx >= 0 {
		*a = slices.Delete(*a, idx, idx+1)
	}
}

// Merge merges attributes from b into a. Values from b take precedence
// over existing values in a for attributes with the same name.
// Names in b are not re-validated; b carries already-constructed Attribute
// values, which are trusted (see Attribute).
func (a *Attributes) Merge(b Attributes) {
	if len(b) == 0 {
		return
	}
	if len(b) <= 4 {
		// For small b, use linear search
		for _, attr := range b {
			a.set(attr.Name, attr.Value)
		}
		return
	}
	// Build index map of existing attribute positions
	index := make(map[string]int, len(*a))
	for i, attr := range *a {
		index[attr.Name] = i
	}
	// Merge attributes from b
	for _, attr := range b {
		if idx, exists := index[attr.Name]; exists {
			(*a)[idx].Value = attr.Value
		} else {
			*a = append(*a, attr)
			index[attr.Name] = len(*a) - 1
		}
	}
}
