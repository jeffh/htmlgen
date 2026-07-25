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

// ScrollPosition represents a scroll position.
type ScrollPosition string

const (
	// Top scrolls to the top of the element.
	Top ScrollPosition = "top"
	// Bottom scrolls to the bottom of the element.
	Bottom ScrollPosition = "bottom"
)

// SwapBuilder builds an hx-swap attribute value.
//
// SwapBuilder implements h.AttrBuilder, so it can be passed directly to
// element methods like B.Div without an explicit terminator method.
type SwapBuilder struct {
	strategy  SwapStrategy
	modifiers []string
}

// Attribute returns the hx-swap attribute.
func (s *SwapBuilder) Attribute() h.Attribute {
	return h.Attr("hx-swap", s.String())
}

// String returns the rendered hx-swap value (strategy plus modifiers).
func (s *SwapBuilder) String() string {
	if len(s.modifiers) == 0 {
		return string(s.strategy)
	}
	var sb strings.Builder
	sb.WriteString(string(s.strategy))
	for _, mod := range s.modifiers {
		sb.WriteString(" ")
		sb.WriteString(mod)
	}
	return sb.String()
}

// Swap creates an hx-swap builder with the given strategy.
//
// Example:
//
//	hx.Swap(hx.InnerHTML)
//	hx.Swap(hx.OuterHTML).Transition()
//	hx.Swap(hx.InnerHTML).Delay(100*time.Millisecond).SettleDelay(200*time.Millisecond)
func Swap(strategy SwapStrategy) *SwapBuilder {
	return &SwapBuilder{strategy: strategy}
}

// Transition enables the View Transitions API for this swap.
func (s *SwapBuilder) Transition() *SwapBuilder {
	s.modifiers = append(s.modifiers, "transition:true")
	return s
}

// Delay adds a delay before the swap is performed.
func (s *SwapBuilder) Delay(d time.Duration) *SwapBuilder {
	s.modifiers = append(s.modifiers, "swap:"+formatDuration(d))
	return s
}

// SettleDelay adds a delay after the swap before the settle phase.
func (s *SwapBuilder) SettleDelay(d time.Duration) *SwapBuilder {
	s.modifiers = append(s.modifiers, "settle:"+formatDuration(d))
	return s
}

// IgnoreTitle prevents HTMX from updating the page title from the response.
func (s *SwapBuilder) IgnoreTitle() *SwapBuilder {
	s.modifiers = append(s.modifiers, "ignoreTitle:true")
	return s
}

// Scroll scrolls the target element to the specified position after swap.
func (s *SwapBuilder) Scroll(pos ScrollPosition) *SwapBuilder {
	s.modifiers = append(s.modifiers, "scroll:"+string(pos))
	return s
}

// ScrollTarget scrolls a specific element to the specified position after swap.
func (s *SwapBuilder) ScrollTarget(selector string, pos ScrollPosition) *SwapBuilder {
	s.modifiers = append(s.modifiers, "scroll:"+selector+":"+string(pos))
	return s
}

// Show scrolls the viewport to show the target element at the specified position.
func (s *SwapBuilder) Show(pos ScrollPosition) *SwapBuilder {
	s.modifiers = append(s.modifiers, "show:"+string(pos))
	return s
}

// ShowTarget scrolls the viewport to show a specific element at the specified position.
func (s *SwapBuilder) ShowTarget(selector string, pos ScrollPosition) *SwapBuilder {
	s.modifiers = append(s.modifiers, "show:"+selector+":"+string(pos))
	return s
}

// ShowWindow scrolls the window to the specified position.
func (s *SwapBuilder) ShowWindow(pos ScrollPosition) *SwapBuilder {
	s.modifiers = append(s.modifiers, "show:window:"+string(pos))
	return s
}

// ShowNone disables automatic scrolling to show the swapped content.
func (s *SwapBuilder) ShowNone() *SwapBuilder {
	s.modifiers = append(s.modifiers, "show:none")
	return s
}

// FocusScroll controls whether HTMX scrolls to bring a focused element into view.
func (s *SwapBuilder) FocusScroll(enabled bool) *SwapBuilder {
	if enabled {
		s.modifiers = append(s.modifiers, "focus-scroll:true")
	} else {
		s.modifiers = append(s.modifiers, "focus-scroll:false")
	}
	return s
}

// formatDuration formats a duration for HTMX (e.g., "500ms", "1s").
func formatDuration(d time.Duration) string {
	return d.String()
}
