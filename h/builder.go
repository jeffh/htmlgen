package h

// TagArg is a marker interface for types that can be passed to tag functions.
// Valid types are: Attributes, Attribute, and Builder.
type TagArg interface {
	isTagArg()
}

// Builder is the interface implemented by all HTML node builders.
// Implementations write their HTML representation to the provided Writer.
type Builder interface {
	TagArg
	Build(w *Writer) error
}

type htmlTagBuilder struct {
	Attrs    Attributes
	Children []Builder
}

func (b *htmlTagBuilder) isTagArg() {}
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

// parseTagArgs separates attributes from children in a variadic argument list.
// Multiple Attributes/Attribute are merged (later values override earlier ones).
// All Builder arguments become children.
func parseTagArgs(args []TagArg) (Attributes, []Builder) {
	var attrs Attributes
	var children []Builder

	for _, arg := range args {
		if arg == nil {
			continue
		}
		switch v := arg.(type) {
		case Attributes:
			if attrs == nil {
				attrs = v
			} else {
				attrs.Merge(v)
			}
		case Attribute:
			if attrs == nil {
				attrs = Attributes{v}
			} else {
				attrs.Set(v.Name, v.Value)
			}
		case Builder:
			children = append(children, v)
		}
	}
	return attrs, children
}

func tag(name string, args ...TagArg) Builder {
	attrs, children := parseTagArgs(args)
	return &tagBuilder{
		Name:     name,
		Attrs:    attrs,
		Children: children,
	}
}

func stag(name string, args ...TagArg) Builder {
	attrs, children := parseTagArgs(args)
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

func (b *tagBuilder) isTagArg() {}
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

func (b *fragmentBuilder) isTagArg() {}
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

func (b *textBuilder) isTagArg() {}
func (b *textBuilder) Build(w *Writer) error {
	if b.IsRaw {
		return w.Raw(b.Text)
	}
	return w.Text(b.Text)
}
