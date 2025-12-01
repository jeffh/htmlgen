package html

import (
	"maps"
	"slices"
	"sort"
)

type Builder interface {
	Build(w *Writer) error
}

func attrsMapToSlice(attrs map[string]string) []string {
	result := make([]string, 0, len(attrs)*2)
	keys := slices.Collect(maps.Keys(attrs))
	sort.Strings(keys)
	for _, k := range keys {
		result = append(result, k, attrs[k])
	}
	return result
}

func tag(name string, attrs []string, children []Builder) Builder {
	return &tagBuilder{
		Name:     name,
		Attrs:    attrs,
		Children: children,
	}
}

type tagBuilder struct {
	Name      string
	Attrs     []string
	Children  []Builder
	SelfClose bool
}

func (b *tagBuilder) Build(w *Writer) error {
	if b.SelfClose {
		return w.SelfClosingTag(b.Name, b.Attrs...)
	}

	if err := w.OpenTag(b.Name, b.Attrs...); err != nil {
		return err
	}

	for _, child := range b.Children {
		if child != nil {
			if err := child.Build(w); err != nil {
				return err
			}
		}
	}

	return w.CloseTag(b.Name)
}
