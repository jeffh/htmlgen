package h

import "iter"

// If returns ifTrue if cond is true, otherwise returns ifElse.
// This enables conditional rendering in builder expressions:
//
//	h.Div(
//	    h.If(user.IsAdmin,
//	        h.Span(h.Text("Admin")),
//	        h.Span(h.Text("User")),
//	    ),
//	)
func If(cond bool, ifTrue, ifElse Builder) Builder {
	if cond {
		return ifTrue
	}
	return ifElse
}

// When returns ifTrue if cond is true, otherwise returns nil.
// Nil builders are safely skipped during rendering.
// This is a convenience wrapper around If for cases without an else branch:
//
//	h.Div(
//	    h.When(showWarning, h.Span(h.Text("Warning!"))),
//	)
func When(cond bool, ifTrue Builder) Builder {
	if cond {
		return ifTrue
	}
	return nil
}

// First returns the first non-nil Builder from the provided arguments.
// Returns nil if all arguments are nil or no arguments are provided.
// Useful for fallback content patterns:
//
//	h.Div(
//	    h.First(customHeader, defaultHeader, h.Text("Untitled")),
//	)
func First(b ...Builder) Builder {
	for _, builder := range b {
		if builder != nil {
			return builder
		}
	}
	return nil
}

// forEachBuilder is a lazy builder that maps over an iterator during Build.
type forEachBuilder[X any] struct {
	seq iter.Seq[X]
	fn  func(X) Builder
}

func (m *forEachBuilder[X]) isTagArg() {}

func (m *forEachBuilder[X]) Build(w *Writer) error {
	for x := range m.seq {
		if b := m.fn(x); b != nil {
			if err := b.Build(w); err != nil {
				return err
			}
		}
	}
	return nil
}

// ForEach creates a Builder that lazily maps over an iterator sequence,
// rendering each resulting Builder in order. Nil builders returned
// by the mapping function are skipped.
//
// The iterator is consumed during rendering, not when ForEach is called.
//
//	h.Ul(
//	    h.ForEach(slices.Values(items), func(item Item) h.Builder {
//	        return h.Li(h.Text(item.Name))
//	    }),
//	)
func ForEach[X any](s iter.Seq[X], fn func(X) Builder) Builder {
	return &forEachBuilder[X]{seq: s, fn: fn}
}

// forEach2Builder is a lazy builder that maps over a two-value iterator during Build.
type forEach2Builder[X, Y any] struct {
	seq iter.Seq2[X, Y]
	fn  func(X, Y) Builder
}

func (m *forEach2Builder[X, Y]) isTagArg() {}

func (m *forEach2Builder[X, Y]) Build(w *Writer) error {
	for x, y := range m.seq {
		if b := m.fn(x, y); b != nil {
			if err := b.Build(w); err != nil {
				return err
			}
		}
	}
	return nil
}

// ForEach2 creates a Builder that lazily maps over a two-value iterator sequence,
// rendering each resulting Builder in order. Nil builders returned
// by the mapping function are skipped.
//
// The iterator is consumed during rendering, not when ForEach2 is called.
// Common use cases include maps.All() for key-value pairs or slices.All()
// for index-value pairs:
//
//	h.Dl(
//	    h.ForEach2(maps.All(definitions), func(term, def string) h.Builder {
//	        return h.Fragment(h.Dt(h.Text(term)), h.Dd(h.Text(def)))
//	    }),
//	)
func ForEach2[X, Y any](s iter.Seq2[X, Y], fn func(X, Y) Builder) Builder {
	return &forEach2Builder[X, Y]{seq: s, fn: fn}
}
