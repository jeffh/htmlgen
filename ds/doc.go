// Package ds provides helpers for building Datastar (https://data-star.dev/) reactive attributes.
//
// The package uses a fluent builder API: each event/signal/etc. constructor
// returns a builder type with chainable modifier methods. Builders implement
// h.AttrBuilder, so h.AttrsOf and h.Attributes.With collect them into the
// h.Attributes an element method takes, with no explicit terminator method.
//
// Common usage:
//
//	h.Render(w, func(b *h.B) {
//	    b.Button(
//	        h.AttrsOf(ds.OnClick(ds.Raw("$count++")).PreventDefault().Debounce(300*time.Millisecond)),
//	        func(b *h.B) { b.Text("Click me") },
//	    )
//	    b.Input(h.AttrsOf(ds.Bind("name").Case(ds.CamelCase)))
//	    b.Div(
//	        h.Attrs("class", "panel").With(ds.Show(ds.Raw("$visible"))),
//	        func(b *h.B) { b.Text("Hello") },
//	    )
//	})
//
// # Builders
//
//   - EventBuilder — returned by OnClick, OnSubmit, OnInput, OnChange, On, OnRAF, OnResize.
//     Methods: PreventDefault, StopPropagation, Once, Passive, Capture, Outside, Window,
//     Document, ViewTransition, Case, Delay, Debounce, Throttle, Then.
//   - IntersectBuilder — returned by OnIntersect. Methods: Once, Half, Full, Exit,
//     Threshold, Delay, Debounce, Throttle, ViewTransition.
//   - IntervalBuilder — returned by OnInterval. Methods: Duration, ViewTransition.
//   - SignalPatchBuilder — returned by OnSignalPatch. Methods: Delay, Debounce, Throttle.
//   - InitBuilder — returned by Init. Methods: Delay, ViewTransition.
//   - BindBuilder — returned by Bind, BindKey. Methods: Case, Prop, Event.
//   - NamedBuilder — returned by IndicatorKey, Signal, SignalExpr, Computed, Ref, Class,
//     MatchMedia (Pro). Method: Case.
//   - SignalsBuilder — returned by Signals. Methods: Case, IfMissing.
//   - JsonSignalsBuilder — returned by JsonSignalsDebug. Method: Terse.
//   - PersistBuilder (Pro) — returned by Persist, PersistKey. Method: Session.
//   - QueryStringBuilder (Pro) — returned by QueryString. Methods: Filter, History.
//   - ScrollBuilder (Pro) — returned by ScrollIntoView. Methods: Smooth, Instant, Auto,
//     HStart, HCenter, HEnd, HNearest, VStart, VCenter, VEnd, VNearest, Focus.
//
// # Plain attributes (return h.Attribute directly)
//
// Show, Text, Hide, Classes, ClassesExpr, Style, Styles, StylesExpr,
// Attribute, Attrs, AttrsExpr, Effect, Indicator, Ignore, IgnoreSelf,
// IgnoreMorph, PreserveAttr, OnSignalPatchFilter, Animate (Pro),
// CustomValidity (Pro), ReplaceURL (Pro), ViewTransitionName (Pro).
//
// The *Expr variants (ClassesExpr, StylesExpr, AttrsExpr) keep map values as
// JavaScript expressions; Classes, Styles and Attrs JSON-encode them into
// quoted string literals, which are always truthy.
//
// # Signals and composition
//
// Sig is a typed signal name: Sig("open").Toggle(), Sig("q").Clear(),
// Sig("plan").Sub("open") (derived "plan_open"). Value has Not, And, Or and
// Ternary combinators, Do bridges js statements into a Value for statement
// positions, and Confirm guards actions behind a confirm() dialog. Datastar
// expressions see evt and el — use Evt, El, EvtTarget, EvtValue and EvtKey,
// not the event-based js package helpers.
//
// # HTTP actions
//
// Get, Post, Put, Patch, Delete and the *Dynamic / *WithOptions variants
// return Value (a wrapped js.Expr). Pass them as the first argument to an
// event handler:
//
//	ds.OnClick(ds.Get("/api/data", ds.OnSuccess(ds.Raw("$loaded = true"))))
//
// RequestOptions is a fluent builder for *WithOptions calls.
package ds
