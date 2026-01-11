package hx

import (
	"strings"

	"github.com/jeffh/htmlgen/h"
)

// Boost creates an hx-boost attribute for progressive enhancement.
// When true, links and forms will use AJAX instead of full page loads.
func Boost(enabled bool) h.Attribute {
	if enabled {
		return h.Attr("hx-boost", "true")
	}
	return h.Attr("hx-boost", "false")
}

// PushURL creates an hx-push-url attribute that pushes a URL to the browser history.
//
// Values:
//   - "true" to push the request URL
//   - "false" to disable
//   - A specific URL to push
func PushURL(url string) h.Attribute {
	return h.Attr("hx-push-url", url)
}

// ReplaceURL creates an hx-replace-url attribute that replaces the current URL.
//
// Values:
//   - "true" to replace with the request URL
//   - "false" to disable
//   - A specific URL to use
func ReplaceURL(url string) h.Attribute {
	return h.Attr("hx-replace-url", url)
}

// Confirm creates an hx-confirm attribute that shows a confirmation dialog.
func Confirm(message string) h.Attribute {
	return h.Attr("hx-confirm", message)
}

// Prompt creates an hx-prompt attribute that shows an input prompt.
// The user's input is included as HX-Prompt header in the request.
func Prompt(message string) h.Attribute {
	return h.Attr("hx-prompt", message)
}

// Indicator creates an hx-indicator attribute that specifies an element
// to apply the htmx-request class to during the request.
//
// Values:
//   - A CSS selector (e.g., "#spinner")
//   - "closest <selector>" for the closest ancestor (e.g., "closest .spinner")
func Indicator(selector string) h.Attribute {
	return h.Attr("hx-indicator", selector)
}

// DisabledElt creates an hx-disabled-elt attribute that specifies elements
// to disable during the request.
//
// Values:
//   - A CSS selector (e.g., "button")
//   - "this" for the triggering element
//   - "closest <selector>" for the closest ancestor
func DisabledElt(selector string) h.Attribute {
	return h.Attr("hx-disabled-elt", selector)
}

// SyncStrategy specifies how requests should be synchronized.
type SyncStrategy string

const (
	// SyncDrop drops new requests while a request is in flight.
	SyncDrop SyncStrategy = "drop"
	// SyncAbort aborts the current request when a new one arrives.
	SyncAbort SyncStrategy = "abort"
	// SyncReplace aborts and replaces the current request.
	SyncReplace SyncStrategy = "replace"
	// SyncQueue queues requests to be processed after the current one.
	SyncQueue SyncStrategy = "queue"
	// SyncQueueFirst queues first request only.
	SyncQueueFirst SyncStrategy = "queue first"
	// SyncQueueLast queues last request only (default queue behavior).
	SyncQueueLast SyncStrategy = "queue last"
	// SyncQueueAll queues all requests.
	SyncQueueAll SyncStrategy = "queue all"
)

// Sync creates an hx-sync attribute that controls request synchronization.
//
// The selector can be:
//   - A CSS selector (e.g., "#form")
//   - "this" for the triggering element
//   - "closest <selector>" for the closest ancestor
//
// If selector is empty, only the strategy is used.
func Sync(selector string, strategy SyncStrategy) h.Attribute {
	if selector == "" {
		return h.Attr("hx-sync", string(strategy))
	}
	return h.Attr("hx-sync", selector+":"+string(strategy))
}

// Validate creates an hx-validate attribute that forces form validation.
func Validate(enabled bool) h.Attribute {
	if enabled {
		return h.Attr("hx-validate", "true")
	}
	return h.Attr("hx-validate", "false")
}

// Preserve creates an hx-preserve attribute that preserves the element during swaps.
func Preserve() h.Attribute {
	return h.Attr("hx-preserve", "true")
}

// Disable creates an hx-disable attribute that disables HTMX processing.
func Disable() h.Attribute {
	return h.Attr("hx-disable", "")
}

// Disinherit creates an hx-disinherit attribute that prevents attribute inheritance.
//
// Pass no arguments or "*" to disinherit all attributes, or list specific attributes.
//
// Example:
//
//	hx.Disinherit()                    // disinherit all
//	hx.Disinherit("hx-target")         // disinherit specific
//	hx.Disinherit("hx-target", "hx-swap")
func Disinherit(attrs ...string) h.Attribute {
	if len(attrs) == 0 {
		return h.Attr("hx-disinherit", "*")
	}
	return h.Attr("hx-disinherit", strings.Join(attrs, " "))
}

// Inherit creates an hx-inherit attribute that enables attribute inheritance
// when it's been disabled globally via hx-disinherit.
//
// Pass no arguments or "*" to inherit all attributes, or list specific attributes.
func Inherit(attrs ...string) h.Attribute {
	if len(attrs) == 0 {
		return h.Attr("hx-inherit", "*")
	}
	return h.Attr("hx-inherit", strings.Join(attrs, " "))
}

// HistoryElt creates an hx-history-elt attribute that marks the element
// to snapshot and restore during history navigation.
func HistoryElt() h.Attribute {
	return h.Attr("hx-history-elt", "")
}

// History creates an hx-history attribute that controls history caching.
// Set to false to prevent sensitive data from being stored in history cache.
func History(enabled bool) h.Attribute {
	if enabled {
		return h.Attr("hx-history", "true")
	}
	return h.Attr("hx-history", "false")
}
