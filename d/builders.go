package d

import (
	"fmt"
	"strings"

	"github.com/jeffh/htmlgen/h"
)

type AttrMutator interface{ Modify(*attrBuilder) }
type AttrValueAppender interface{ Append(*strings.Builder) }

type AttrFunc func(*attrBuilder)
type AttrValueFunc func(*strings.Builder)

func (dof AttrFunc) Modify(attr *attrBuilder)        { dof(attr) }
func (vaf AttrValueFunc) Append(sb *strings.Builder) { vaf(sb) }

type ValueMutator struct {
	AttrMutator
	AttrValueAppender
}

func ValueMutatorFunc(do func(*strings.Builder)) ValueMutator {
	return ValueMutator{
		AttrMutator: AttrFunc(func(attr *attrBuilder) {
			var sb strings.Builder
			do(&sb)
			attr.AppendStatement(sb.String())
		}),
		AttrValueAppender: AttrValueFunc(do),
	}
}

type attrBuilder struct {
	name       strings.Builder
	statements []string
}

func (attr *attrBuilder) AppendStatement(s string) {
	attr.statements = append(attr.statements, s)
}

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

func exprAttr(name string, options ...AttrMutator) h.Attribute {
	attr := buildAttr(name, options...)
	return h.Attr(attr.name.String(), strings.Join(attr.statements, "; "))
}
