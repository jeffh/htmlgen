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
//	h.Button(
//	    hx.Get("/api/data"),
//	    hx.Target("#results"),
//	    hx.Swap(hx.InnerHTML),
//	    h.Text("Load"),
//	)
//
// Complex trigger:
//
//	h.Input(
//	    hx.Post("/search"),
//	    hx.Trigger("keyup", hx.Changed(), hx.Delay(500*time.Millisecond)),
//	    hx.Target("#results"),
//	)
//
// Swap with modifiers:
//
//	h.Div(
//	    hx.Get("/content"),
//	    hx.Swap(hx.OuterHTML, hx.Transition(), hx.SwapDelay(100*time.Millisecond)),
//	)
package hx
