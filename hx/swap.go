package hx

import (
	"strings"
	"time"

	"github.com/jeffh/htmlgen/h"
)

// SwapStrategy represents an HTMX swap strategy.
type SwapStrategy string

const (
	// InnerHTML replaces the inner html of the target element (default).
	InnerHTML SwapStrategy = "innerHTML"
	// OuterHTML replaces the entire target element with the response.
	OuterHTML SwapStrategy = "outerHTML"
	// TextContent replaces the text content of the target element, without parsing.
	TextContent SwapStrategy = "textContent"
	// BeforeBegin inserts the response before the target element.
	BeforeBegin SwapStrategy = "beforebegin"
	// AfterBegin inserts the response before the first child of the target element.
	AfterBegin SwapStrategy = "afterbegin"
	// BeforeEnd inserts the response after the last child of the target element.
	BeforeEnd SwapStrategy = "beforeend"
	// AfterEnd inserts the response after the target element.
	AfterEnd SwapStrategy = "afterend"
	// SwapDelete deletes the target element regardless of the response.
	SwapDelete SwapStrategy = "delete"
	// None does not append content from response (useful for response headers only).
	None SwapStrategy = "none"
)

// SwapMod is a modifier for swap behavior.
type SwapMod interface {
	applySwap(*swapBuilder)
}

type swapModFunc func(*swapBuilder)

func (f swapModFunc) applySwap(b *swapBuilder) { f(b) }

type swapBuilder struct {
	strategy  SwapStrategy
	modifiers []string
}

// Swap creates an hx-swap attribute with the given strategy and optional modifiers.
//
// Example:
//
//	hx.Swap(hx.InnerHTML)
//	hx.Swap(hx.OuterHTML, hx.Transition())
//	hx.Swap(hx.InnerHTML, hx.SwapDelay(100*time.Millisecond), hx.SettleDelay(200*time.Millisecond))
func Swap(strategy SwapStrategy, mods ...SwapMod) h.Attribute {
	b := &swapBuilder{strategy: strategy}
	for _, mod := range mods {
		mod.applySwap(b)
	}
	return h.Attr("hx-swap", b.String())
}

func (b *swapBuilder) String() string {
	if len(b.modifiers) == 0 {
		return string(b.strategy)
	}
	var sb strings.Builder
	sb.WriteString(string(b.strategy))
	for _, mod := range b.modifiers {
		sb.WriteString(" ")
		sb.WriteString(mod)
	}
	return sb.String()
}

// Transition enables the View Transitions API for this swap.
func Transition() SwapMod {
	return swapModFunc(func(b *swapBuilder) {
		b.modifiers = append(b.modifiers, "transition:true")
	})
}

// SwapDelay adds a delay before the swap is performed.
func SwapDelay(d time.Duration) SwapMod {
	return swapModFunc(func(b *swapBuilder) {
		b.modifiers = append(b.modifiers, "swap:"+formatDuration(d))
	})
}

// SettleDelay adds a delay after the swap before the settle phase.
func SettleDelay(d time.Duration) SwapMod {
	return swapModFunc(func(b *swapBuilder) {
		b.modifiers = append(b.modifiers, "settle:"+formatDuration(d))
	})
}

// IgnoreTitle prevents HTMX from updating the page title from the response.
func IgnoreTitle() SwapMod {
	return swapModFunc(func(b *swapBuilder) {
		b.modifiers = append(b.modifiers, "ignoreTitle:true")
	})
}

// ScrollPosition represents a scroll position.
type ScrollPosition string

const (
	// Top scrolls to the top of the element.
	Top ScrollPosition = "top"
	// Bottom scrolls to the bottom of the element.
	Bottom ScrollPosition = "bottom"
)

// Scroll scrolls the target element to the specified position after swap.
func Scroll(pos ScrollPosition) SwapMod {
	return swapModFunc(func(b *swapBuilder) {
		b.modifiers = append(b.modifiers, "scroll:"+string(pos))
	})
}

// ScrollTarget scrolls a specific element to the specified position after swap.
func ScrollTarget(selector string, pos ScrollPosition) SwapMod {
	return swapModFunc(func(b *swapBuilder) {
		b.modifiers = append(b.modifiers, "scroll:"+selector+":"+string(pos))
	})
}

// Show scrolls the viewport to show the target element at the specified position.
func Show(pos ScrollPosition) SwapMod {
	return swapModFunc(func(b *swapBuilder) {
		b.modifiers = append(b.modifiers, "show:"+string(pos))
	})
}

// ShowTarget scrolls the viewport to show a specific element at the specified position.
func ShowTarget(selector string, pos ScrollPosition) SwapMod {
	return swapModFunc(func(b *swapBuilder) {
		b.modifiers = append(b.modifiers, "show:"+selector+":"+string(pos))
	})
}

// ShowWindow scrolls the window to the specified position.
func ShowWindow(pos ScrollPosition) SwapMod {
	return swapModFunc(func(b *swapBuilder) {
		b.modifiers = append(b.modifiers, "show:window:"+string(pos))
	})
}

// ShowNone disables automatic scrolling to show the swapped content.
func ShowNone() SwapMod {
	return swapModFunc(func(b *swapBuilder) {
		b.modifiers = append(b.modifiers, "show:none")
	})
}

// FocusScroll controls whether HTMX scrolls to bring a focused element into view.
func FocusScroll(enabled bool) SwapMod {
	return swapModFunc(func(b *swapBuilder) {
		if enabled {
			b.modifiers = append(b.modifiers, "focus-scroll:true")
		} else {
			b.modifiers = append(b.modifiers, "focus-scroll:false")
		}
	})
}

// formatDuration formats a duration for HTMX (e.g., "500ms", "1s").
func formatDuration(d time.Duration) string {
	if d%time.Second == 0 {
		return d.String()
	}
	return d.String()
}
