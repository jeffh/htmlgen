// Package d provides helpers for building Datastar (https://data-star.dev/) reactive attributes.
//
// This package includes:
//   - Signal management: Signal, Signals, Computed, Bind, BindKey
//   - Event handlers: On, OnClick, OnSubmit, OnInput, OnChange, OnLoad, OnIntersect, OnInterval, OnSignalPatch
//   - Reactive display: Show, Text, Class, Classes, Style, Styles, Attribute, Attrs
//   - DOM control: Ref, Indicator, Ignore, IgnoreSelf, IgnoreMorph, PreserveAttr, Effect, Init
//   - HTTP actions: Get, Post, Put, Patch, Delete (and Dynamic variants)
//   - HTTP options: RequestOptions, ContentType, FilterSignals, Headers, OpenWhenHidden, retry config
//   - Core actions: Peek, SetAll, ToggleAll
//   - Modifiers: Debounce, Throttle, Delay, Duration, Once, PreventDefault, ViewTransition, etc.
//
// Pro attributes (require commercial license) are available in this package but documented
// as requiring a Datastar Pro license: Animate, CustomValidity, OnRAF, OnResize, Persist,
// QueryString, ReplaceURL, ScrollIntoView, ViewTransitionName.
//
// Pro actions (require commercial license): Clipboard, ClipboardBase64, Fit, FitClamped,
// FitRounded, FitClampedRounded.
package d
