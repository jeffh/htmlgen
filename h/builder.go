package h

import (
	"maps"
	"slices"
	"sort"
)

// Builder is the interface implemented by all HTML node builders.
// Implementations write their HTML representation to the provided Writer.
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

type htmlTagBuilder struct {
	Attrs    Attributes
	Children []Builder
}

func (b *htmlTagBuilder) Build(w *Writer) error {
	const name = "html"
	w.Doctype()
	b.Attrs.SetDefault("lang", "en")
	if err := w.OpenTag(name, b.Attrs); err != nil {
		return err
	}

	for _, child := range b.Children {
		if child != nil {
			if err := child.Build(w); err != nil {
				return err
			}
		}
	}

	return w.CloseTag(name)
}

func tag(name string, attrs Attributes, children []Builder) Builder {
	return &tagBuilder{
		Name:     name,
		Attrs:    attrs,
		Children: children,
	}
}

func stag(name string, attrs Attributes, children []Builder) Builder {
	return &tagBuilder{
		Name:      name,
		Attrs:     attrs,
		Children:  children,
		SelfClose: true,
	}
}

type tagBuilder struct {
	Name      string
	Attrs     Attributes
	Children  []Builder
	SelfClose bool
}

func (b *tagBuilder) Build(w *Writer) error {
	if b.SelfClose && len(b.Children) == 0 {
		return w.SelfClosingTag(b.Name, b.Attrs)
	}

	if err := w.OpenTag(b.Name, b.Attrs); err != nil {
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

type fragmentBuilder struct {
	Children []Builder
}

func (b *fragmentBuilder) Build(w *Writer) error {
	for _, child := range b.Children {
		if child != nil {
			if err := child.Build(w); err != nil {
				return err
			}
		}
	}
	return nil
}

type textBuilder struct {
	Text  string
	IsRaw bool
}

func (b *textBuilder) Build(w *Writer) error {
	if b.IsRaw {
		return w.Raw(b.Text)
	}
	return w.Text(b.Text)
}
