package d

import (
	"fmt"
	"strings"

	"github.com/jeffh/htmlgen/h"
)

// AttrMutator modifies an attribute being built (name and/or statements).
type AttrMutator interface{ Modify(*attrBuilder) }

// AttrFunc is a function that implements AttrMutator.
type AttrFunc func(*attrBuilder)

func (f AttrFunc) Modify(attr *attrBuilder) { f(attr) }

// attrBuilder tracks the attribute name and statements during building.
type attrBuilder struct {
	name       strings.Builder
	statements []string
}

// AppendStatement adds a JavaScript statement to the attribute value.
func (attr *attrBuilder) AppendStatement(s string) {
	attr.statements = append(attr.statements, s)
}

// buildAttr creates an attrBuilder with the given name and applies all mutators.
func buildAttr(name string, options ...AttrMutator) *attrBuilder {
	attr := &attrBuilder{
		statements: make([]string, 0, len(options)),
	}
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Sprintf("attr %q: %s", attr.name.String(), r))
		}
	}()
	attr.name.WriteString(name)
	attr.statements = make([]string, 0, len(options))
	for _, opt := range options {
		opt.Modify(attr)
	}
	return attr
}

// exprAttr builds an h.Attribute from a base name and mutators.
// Statements are joined with "; " as the attribute value.
func exprAttr(name string, options ...AttrMutator) h.Attribute {
	attr := buildAttr(name, options...)
	return h.Attr(attr.name.String(), strings.Join(attr.statements, "; "))
}
