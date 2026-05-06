// Package ds provides helpers for building Datastar (https://data-star.dev/) reactive attributes.
//
// The package uses a fluent builder API: each event/signal/etc. constructor
// returns a builder type with chainable modifier methods. Builders implement
// h.AttrBuilder so they can be passed directly to tag functions without an
// explicit terminator method.
//
// Common usage:
//
//	h.Button(
//	    ds.OnClick(ds.Raw("$count++")).PreventDefault().Debounce(300*time.Millisecond),
//	    h.Text("Click me"),
//	)
//
//	h.Input(
//	    ds.Bind("name").Case(ds.CamelCase),
//	)
//
//	h.Div(
//	    ds.Show(ds.Raw("$visible")),
//	    h.Text("Hello"),
//	)
//
// # Builders
//
//   - EventBuilder — returned by OnClick, OnSubmit, OnInput, OnChange, OnLoad, On, OnRAF, OnResize.
//     Methods: PreventDefault, StopPropagation, Once, Passive, Capture, Outside, Window,
//     Self, ViewTransition, IfMissing, Delay, Debounce, Throttle, Then.
//   - IntersectBuilder — returned by OnIntersect. Methods: Once, Half, Full, Exit, Threshold.
//   - IntervalBuilder — returned by OnInterval. Method: Duration.
//   - SignalPatchBuilder — returned by OnSignalPatch. Method: Case.
//   - NamedBuilder — returned by Bind, BindKey, IndicatorKey, Signal, SignalExpr, Computed, Ref. Method: Case.
//   - SignalsBuilder — returned by Signals and JsonSignalsDebug. Methods: Case, IfMissing, Terse.
//   - PersistBuilder (Pro) — returned by Persist, PersistKey. Method: Session.
//   - QueryStringBuilder (Pro) — returned by QueryString. Methods: Filter, History.
//   - ScrollBuilder (Pro) — returned by ScrollIntoView. Methods: Smooth, Instant, Auto,
//     HStart, HCenter, HEnd, HNearest, VStart, VCenter, VEnd, VNearest, Focus.
//
// # Plain attributes (return h.Attribute directly)
//
// Show, Text, Hide, Class, Classes, Style, Styles, Attribute, Attrs, Effect, Init,
// Indicator, Ignore, IgnoreSelf, IgnoreMorph, PreserveAttr, OnSignalPatchFilter,
// Animate (Pro), CustomValidity (Pro), ReplaceURL (Pro), ViewTransitionName (Pro).
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
