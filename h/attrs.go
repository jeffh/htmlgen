package h

import (
	"maps"
	"slices"
	"sort"
)

// Attribute represents a single HTML attribute as a name-value pair.
type Attribute struct {
	Name  string
	Value string
}

func (a Attribute) isTagArg() {}

// Attr creates a new Attribute with the given name and value.
// Panics if name is empty.
func Attr(name, value string) Attribute {
	if name == "" {
		panic("attribute name cannot be empty")
	}
	return Attribute{Name: name, Value: value}
}

// Attributes is a slice of Attribute values representing HTML element attributes.
// It provides methods for getting, setting, and deleting attributes by name.
type Attributes []Attribute

func (a Attributes) isTagArg() {}

// Attrs creates an Attributes slice from alternating key-value string pairs.
// Panics if an odd number of arguments is provided or if any key is empty.
//
// Example: Attrs("href", "/home", "class", "nav-link")
func Attrs(kv ...string) Attributes {
	if len(kv)%2 != 0 {
		panic("Attrs(...) expects an even number of arguments")
	}
	if len(kv) == 0 {
		return Attributes(nil)
	}
	results := make(Attributes, 0, len(kv)/2)
	for i := 0; i < len(kv); i += 2 {
		if kv[i] == "" {
			panic("attribute name cannot be empty")
		}
		results = append(results, Attribute{Name: kv[i], Value: kv[i+1]})
	}
	return results
}

// AttrsMap creates an Attributes slice from a map of key-value pairs.
// Keys are sorted alphabetically for deterministic output.
func AttrsMap(m map[string]string) Attributes {
	result := make(Attributes, 0, len(m))
	keys := slices.Collect(maps.Keys(m))
	sort.Strings(keys)
	for _, k := range keys {
		result = append(result, Attribute{k, m[k]})
	}
	return Attributes(result)
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
func (a *Attributes) Set(key, value string) {
	idx := a.Index(key)
	if idx >= 0 {
		(*a)[idx].Value = value
	} else {
		*a = append(*a, Attribute{Name: key, Value: value})
	}
}

// SetDefault sets the value for the given attribute name only if
// the attribute does not already exist.
func (a *Attributes) SetDefault(key, value string) {
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
func (a *Attributes) Merge(b Attributes) {
	if len(b) == 0 {
		return
	}
	if len(b) <= 4 {
		// For small b, use linear search
		for _, attr := range b {
			a.Set(attr.Name, attr.Value)
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
