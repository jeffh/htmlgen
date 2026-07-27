// Package hx provides helpers for building HTMX (https://htmx.org/) attributes.
//
// This package includes:
//   - HTTP methods: Get, Post, Put, Patch, Delete
//   - Targeting: Target, Select, SelectOOB, SwapOOB
//   - Swap strategies: Swap with modifiers (Transition, SwapDelay, SettleDelay, etc.)
//   - Triggers: Trigger with modifiers (Once, Changed, Delay, Throttle, From, etc.)
//   - Request config: Include, Vals, ValsJS, Headers, Params, Encoding, Ext
//   - Behavior: Boost, PushURL, ReplaceURL, Confirm, Prompt, Indicator, Sync, etc.
//   - Events: On, OnBeforeRequest, OnAfterSwap, and other HTMX event handlers
//
// Basic usage:
//
//	h.Render(w, func(b *h.B) {
//	    b.Button(
//	        h.AttrsOf(
//	            hx.Get("/api/data"),
//	            hx.Target("#results"),
//	            hx.Swap(hx.InnerHTML),
//	        ),
//	        func(b *h.B) { b.Text("Load") },
//	    )
//
//	    // Complex trigger:
//	    b.Input(h.AttrsOf(
//	        hx.Post("/search"),
//	        hx.Trigger("keyup").Changed().Delay(500*time.Millisecond),
//	        hx.Target("#results"),
//	    ))
//
//	    // Swap with modifiers:
//	    b.Div(
//	        h.AttrsOf(
//	            hx.Get("/content"),
//	            hx.Swap(hx.OuterHTML).Transition().Delay(100*time.Millisecond),
//	        ),
//	        nil,
//	    )
//	})
//
// Trigger and Swap return fluent builders that auto-finalize as h.AttrBuilder,
// so h.AttrsOf and h.Attributes.With collect them without an explicit
// terminator method.
package hx
