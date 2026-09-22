// Package rkt provides helpers for building Datastar Rocket markup.
//
// Rocket is the web-component layer that ships with Datastar. A component is
// registered in JavaScript with rocket('my-tag', {props, setup, render}); the
// registration wires the component's props, its private signals and its local
// actions into Datastar. This package does not generate that rocket(...)
// JavaScript definition — it generates the HTML side: the host element and
// its props, and the Rocket-specific attributes and expressions used inside a
// component's rendered template.
//
// # Requirements
//
// Rocket lives in its own bundle. Load datastar-rocket.js (a superset of
// datastar.js) instead of datastar.js — loading both is wrong:
//
//	h.Render(w, func(b *h.B) {
//	    b.Head(nil, func(b *h.B) { rkt.Script(b) })
//	})
//
// Rocket is BETA upstream (the v1.0.4 bundle is labeled "Rocket beta.2"), so
// its attribute and expression syntax may change between releases. Pin the
// bundle version you test against.
//
// # $$ versus $
//
// Inside a component template, "$$foo" references the component's PRIVATE
// local signal foo, which Rocket rewrites to $._rocket.<tag>.<id>.foo. A
// plain "$foo" is the ordinary global signal and is left untouched. Sig and
// SignalRef in this package emit the "$$" form; use ds.Sig and ds.SignalRef
// for global signals. Root escapes the private scope for attributes that
// would otherwise be scoped to the component.
//
// # Catalog
//
//   - Private signals: Sig (Name, Ref, Value, Not, Set, SetExpr, Toggle,
//     Clear, Eq, NotEq, Sub) and SignalRef.
//   - Local actions: Action (@name(args)) and Dispatch
//     (@dispatchRocket("name", args)).
//   - Template attributes: For, ForRaw, If, ElseIf, Else, Ref (RefBuilder
//     with Case) and Root.
//   - Host element: Component, Prop, Props and the typed prop encoders
//     BoolProp, NumberProp, JSONProp, DateProp, BinProp.
//   - Bundle: Version, BundleURL, AliasedBundleURL, ScriptAttrs, Script,
//     AliasedScript.
//
// Everything that produces an expression returns a ds.Value, so results plug
// straight into ds.OnClick, ds.Text, ds.Show and the rest of the ds package:
//
//	rkt.Component(b, "counter-widget", h.AttrsOf(rkt.Prop("step", 2)), func(b *h.B) {
//	    b.Button(h.AttrsOf(
//	        ds.OnClick(rkt.Action("increment")),
//	        ds.Text(rkt.Sig("count").Value()),
//	    ), nil)
//	})
package rkt
