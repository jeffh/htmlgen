package hx

import (
	"github.com/jeffh/htmlgen/h"
)

// On creates an hx-on attribute for handling DOM events with inline scripts.
// For standard DOM events (click, change, etc.), use the event name directly.
//
// Example:
//
//	hx.On("click", "alert('clicked')")
//	hx.On("htmx:beforeRequest", "console.log('request starting')")
func On(event string, script string) h.Attribute {
	return h.Attr("hx-on:"+event, script)
}

// OnHTMX creates an hx-on attribute for handling HTMX-specific events.
// The "htmx:" prefix is added automatically.
//
// Example:
//
//	hx.OnHTMX("beforeRequest", "console.log('request starting')")
//	hx.OnHTMX("afterSwap", "initializeComponents()")
func OnHTMX(event string, script string) h.Attribute {
	return h.Attr("hx-on::"+event, script)
}

// Standard DOM event handlers

// OnClick creates an hx-on:click attribute.
func OnClick(script string) h.Attribute {
	return h.Attr("hx-on:click", script)
}

// OnChange creates an hx-on:change attribute.
func OnChange(script string) h.Attribute {
	return h.Attr("hx-on:change", script)
}

// OnSubmit creates an hx-on:submit attribute.
func OnSubmit(script string) h.Attribute {
	return h.Attr("hx-on:submit", script)
}

// OnInput creates an hx-on:input attribute.
func OnInput(script string) h.Attribute {
	return h.Attr("hx-on:input", script)
}

// OnKeydown creates an hx-on:keydown attribute.
func OnKeydown(script string) h.Attribute {
	return h.Attr("hx-on:keydown", script)
}

// OnKeyup creates an hx-on:keyup attribute.
func OnKeyup(script string) h.Attribute {
	return h.Attr("hx-on:keyup", script)
}

// OnFocus creates an hx-on:focus attribute.
func OnFocus(script string) h.Attribute {
	return h.Attr("hx-on:focus", script)
}

// OnBlur creates an hx-on:blur attribute.
func OnBlur(script string) h.Attribute {
	return h.Attr("hx-on:blur", script)
}

// OnMouseover creates an hx-on:mouseover attribute.
func OnMouseover(script string) h.Attribute {
	return h.Attr("hx-on:mouseover", script)
}

// OnMouseout creates an hx-on:mouseout attribute.
func OnMouseout(script string) h.Attribute {
	return h.Attr("hx-on:mouseout", script)
}

// HTMX-specific event handlers
// These handle events in the htmx: namespace

// Request lifecycle events

// OnBeforeRequest creates an hx-on::beforeRequest attribute.
// Fires before an AJAX request is issued.
func OnBeforeRequest(script string) h.Attribute {
	return h.Attr("hx-on::before-request", script)
}

// OnBeforeSend creates an hx-on::beforeSend attribute.
// Fires just before the request is sent.
func OnBeforeSend(script string) h.Attribute {
	return h.Attr("hx-on::before-send", script)
}

// OnAfterRequest creates an hx-on::afterRequest attribute.
// Fires after the request completes (success or failure).
func OnAfterRequest(script string) h.Attribute {
	return h.Attr("hx-on::after-request", script)
}

// OnAfterOnLoad creates an hx-on::afterOnLoad attribute.
// Fires after the response has been processed.
func OnAfterOnLoad(script string) h.Attribute {
	return h.Attr("hx-on::after-on-load", script)
}

// Swap lifecycle events

// OnBeforeSwap creates an hx-on::beforeSwap attribute.
// Fires before the swap is performed.
func OnBeforeSwap(script string) h.Attribute {
	return h.Attr("hx-on::before-swap", script)
}

// OnAfterSwap creates an hx-on::afterSwap attribute.
// Fires after the swap completes.
func OnAfterSwap(script string) h.Attribute {
	return h.Attr("hx-on::after-swap", script)
}

// OnAfterSettle creates an hx-on::afterSettle attribute.
// Fires after the DOM has settled.
func OnAfterSettle(script string) h.Attribute {
	return h.Attr("hx-on::after-settle", script)
}

// Error events

// OnResponseError creates an hx-on::responseError attribute.
// Fires when a response error occurs (non-2xx status).
func OnResponseError(script string) h.Attribute {
	return h.Attr("hx-on::response-error", script)
}

// OnSendError creates an hx-on::sendError attribute.
// Fires when a network error prevents the request.
func OnSendError(script string) h.Attribute {
	return h.Attr("hx-on::send-error", script)
}

// OnTimeout creates an hx-on::timeout attribute.
// Fires when a request times out.
func OnTimeout(script string) h.Attribute {
	return h.Attr("hx-on::timeout", script)
}

// OnSwapError creates an hx-on::swapError attribute.
// Fires when an error occurs during swap.
func OnSwapError(script string) h.Attribute {
	return h.Attr("hx-on::swap-error", script)
}

// Configuration events

// OnConfigRequest creates an hx-on::configRequest attribute.
// Fires before the request is configured, allowing modification.
func OnConfigRequest(script string) h.Attribute {
	return h.Attr("hx-on::config-request", script)
}

// History events

// OnHistoryRestore creates an hx-on::historyRestore attribute.
// Fires when history is restored.
func OnHistoryRestore(script string) h.Attribute {
	return h.Attr("hx-on::history-restore", script)
}

// OnPushedIntoHistory creates an hx-on::pushedIntoHistory attribute.
// Fires when a URL is pushed to history.
func OnPushedIntoHistory(script string) h.Attribute {
	return h.Attr("hx-on::pushed-into-history", script)
}

// OnReplacedInHistory creates an hx-on::replacedInHistory attribute.
// Fires when a URL is replaced in history.
func OnReplacedInHistory(script string) h.Attribute {
	return h.Attr("hx-on::replaced-in-history", script)
}

// DOM processing events

// OnLoad creates an hx-on::load attribute (HTMX load event).
// Fires when new content is loaded into the DOM.
// Note: This is different from the standard DOM load event.
func OnHTMXLoad(script string) h.Attribute {
	return h.Attr("hx-on::load", script)
}

// OnBeforeProcessNode creates an hx-on::beforeProcessNode attribute.
// Fires before HTMX processes a node.
func OnBeforeProcessNode(script string) h.Attribute {
	return h.Attr("hx-on::before-process-node", script)
}

// OnAfterProcessNode creates an hx-on::afterProcessNode attribute.
// Fires after HTMX processes a node.
func OnAfterProcessNode(script string) h.Attribute {
	return h.Attr("hx-on::after-process-node", script)
}

// Out-of-band events

// OnOOBBeforeSwap creates an hx-on::oobBeforeSwap attribute.
// Fires before an out-of-band swap.
func OnOOBBeforeSwap(script string) h.Attribute {
	return h.Attr("hx-on::oob-before-swap", script)
}

// OnOOBAfterSwap creates an hx-on::oobAfterSwap attribute.
// Fires after an out-of-band swap.
func OnOOBAfterSwap(script string) h.Attribute {
	return h.Attr("hx-on::oob-after-swap", script)
}

// OnOOBErrorNoTarget creates an hx-on::oobErrorNoTarget attribute.
// Fires when an OOB swap target is not found.
func OnOOBErrorNoTarget(script string) h.Attribute {
	return h.Attr("hx-on::oob-error-no-target", script)
}

// Validation events

// OnValidationFailed creates an hx-on::validation:failed attribute.
// Fires when form validation fails.
func OnValidationFailed(script string) h.Attribute {
	return h.Attr("hx-on::validation:failed", script)
}

// OnValidationHalted creates an hx-on::validation:halted attribute.
// Fires when validation halts further processing.
func OnValidationHalted(script string) h.Attribute {
	return h.Attr("hx-on::validation:halted", script)
}

// SSE/WebSocket events

// OnSSEOpen creates an hx-on::sseOpen attribute.
// Fires when an SSE connection is opened.
func OnSSEOpen(script string) h.Attribute {
	return h.Attr("hx-on::sse-open", script)
}

// OnSSEError creates an hx-on::sseError attribute.
// Fires when an SSE connection error occurs.
func OnSSEError(script string) h.Attribute {
	return h.Attr("hx-on::sse-error", script)
}

// OnNoSSESourceError creates an hx-on::noSSESourceError attribute.
// Fires when no SSE source is found.
func OnNoSSESourceError(script string) h.Attribute {
	return h.Attr("hx-on::no-sse-source-error", script)
}
