package hx

import (
	"strconv"
	"strings"
	"time"

	"github.com/jeffh/htmlgen/h"
)

type triggerSpec struct {
	event     string
	filter    string
	modifiers []string
}

func (s *triggerSpec) String() string {
	var sb strings.Builder
	sb.WriteString(s.event)
	if s.filter != "" {
		sb.WriteString("[")
		sb.WriteString(s.filter)
		sb.WriteString("]")
	}
	for _, mod := range s.modifiers {
		sb.WriteString(" ")
		sb.WriteString(mod)
	}
	return sb.String()
}

// TriggerBuilder builds an hx-trigger attribute value.
//
// TriggerBuilder implements h.AttrBuilder, so it can be passed directly to
// element methods like B.Div without calling a terminator method.
type TriggerBuilder struct {
	triggers []triggerSpec
}

// Attribute returns the hx-trigger attribute. Joins all trigger specs with ", ".
func (t *TriggerBuilder) Attribute() h.Attribute {
	parts := make([]string, len(t.triggers))
	for i, trigger := range t.triggers {
		parts[i] = trigger.String()
	}
	return h.Attr("hx-trigger", strings.Join(parts, ", "))
}

// last returns a pointer to the most recently added trigger spec.
// Used by modifier methods to apply modifiers to the last spec in the chain.
func (t *TriggerBuilder) last() *triggerSpec {
	return &t.triggers[len(t.triggers)-1]
}

// Trigger creates a new trigger with the specified event.
// Use chained modifier methods to configure further.
//
// Example:
//
//	hx.Trigger("click")
//	hx.Trigger("keyup").Changed().Delay(500*time.Millisecond)
//	hx.Trigger("click").From("#other-element")
func Trigger(event string) *TriggerBuilder {
	return &TriggerBuilder{triggers: []triggerSpec{{event: event}}}
}

// And appends another trigger to the builder. Subsequent modifier methods
// apply to this newly added trigger.
//
// Example:
//
//	hx.Trigger("load").And("click").Delay(1*time.Second)
func (t *TriggerBuilder) And(event string) *TriggerBuilder {
	t.triggers = append(t.triggers, triggerSpec{event: event})
	return t
}

// Once makes the most recently added trigger fire only once.
func (t *TriggerBuilder) Once() *TriggerBuilder {
	s := t.last()
	s.modifiers = append(s.modifiers, "once")
	return t
}

// Changed makes the most recently added trigger fire only when the element's value has changed.
func (t *TriggerBuilder) Changed() *TriggerBuilder {
	s := t.last()
	s.modifiers = append(s.modifiers, "changed")
	return t
}

// Delay adds a delay before the most recently added trigger fires.
// If another event occurs during the delay, the timer resets.
func (t *TriggerBuilder) Delay(d time.Duration) *TriggerBuilder {
	s := t.last()
	s.modifiers = append(s.modifiers, "delay:"+formatDuration(d))
	return t
}

// Throttle limits how often the most recently added trigger can fire.
func (t *TriggerBuilder) Throttle(d time.Duration) *TriggerBuilder {
	s := t.last()
	s.modifiers = append(s.modifiers, "throttle:"+formatDuration(d))
	return t
}

// From specifies that the event should be listened for on a different element.
// Supports extended CSS selectors: document, window, closest <sel>, find <sel>, next, previous.
func (t *TriggerBuilder) From(selector string) *TriggerBuilder {
	s := t.last()
	if strings.Contains(selector, " ") && !strings.HasPrefix(selector, "(") {
		s.modifiers = append(s.modifiers, "from:("+selector+")")
	} else {
		s.modifiers = append(s.modifiers, "from:"+selector)
	}
	return t
}

// FromDocument listens for the event on the document.
func (t *TriggerBuilder) FromDocument() *TriggerBuilder {
	s := t.last()
	s.modifiers = append(s.modifiers, "from:document")
	return t
}

// FromWindow listens for the event on the window.
func (t *TriggerBuilder) FromWindow() *TriggerBuilder {
	s := t.last()
	s.modifiers = append(s.modifiers, "from:window")
	return t
}

// FromClosest listens for the event on the closest ancestor matching the selector.
func (t *TriggerBuilder) FromClosest(selector string) *TriggerBuilder {
	s := t.last()
	s.modifiers = append(s.modifiers, "from:closest "+selector)
	return t
}

// FromFind listens for the event on the first descendant matching the selector.
func (t *TriggerBuilder) FromFind(selector string) *TriggerBuilder {
	s := t.last()
	s.modifiers = append(s.modifiers, "from:find "+selector)
	return t
}

// FromNext listens for the event on the next sibling element.
func (t *TriggerBuilder) FromNext(selector ...string) *TriggerBuilder {
	s := t.last()
	if len(selector) > 0 {
		s.modifiers = append(s.modifiers, "from:next "+selector[0])
	} else {
		s.modifiers = append(s.modifiers, "from:next")
	}
	return t
}

// FromPrevious listens for the event on the previous sibling element.
func (t *TriggerBuilder) FromPrevious(selector ...string) *TriggerBuilder {
	s := t.last()
	if len(selector) > 0 {
		s.modifiers = append(s.modifiers, "from:previous "+selector[0])
	} else {
		s.modifiers = append(s.modifiers, "from:previous")
	}
	return t
}

// TriggerTarget filters events to only those whose target matches the selector.
// Note: This is different from hx-target; it filters by event.target.
func (t *TriggerBuilder) TriggerTarget(selector string) *TriggerBuilder {
	s := t.last()
	s.modifiers = append(s.modifiers, "target:"+selector)
	return t
}

// Consume prevents the event from triggering requests on parent elements.
func (t *TriggerBuilder) Consume() *TriggerBuilder {
	s := t.last()
	s.modifiers = append(s.modifiers, "consume")
	return t
}

// QueueMode specifies how events should be queued.
type QueueMode string

const (
	// QueueFirst queues the first event to arrive.
	QueueFirst QueueMode = "first"
	// QueueLast queues the last event to arrive (default).
	QueueLast QueueMode = "last"
	// QueueAll queues all events.
	QueueAll QueueMode = "all"
	// QueueNone disables queuing.
	QueueNone QueueMode = "none"
)

// Queue specifies how events should be queued during a request.
func (t *TriggerBuilder) Queue(mode QueueMode) *TriggerBuilder {
	s := t.last()
	s.modifiers = append(s.modifiers, "queue:"+string(mode))
	return t
}

// Filter adds a JavaScript filter expression to the most recently added trigger.
// The expression should evaluate to true/false.
//
// Example:
//
//	hx.Trigger("click").Filter("ctrlKey")
//	hx.Trigger("keyup").Filter("keyCode==13")
func (t *TriggerBuilder) Filter(jsExpr string) *TriggerBuilder {
	s := t.last()
	s.filter = jsExpr
	return t
}

// TriggerAttr creates a simple hx-trigger attribute with a single event.
// For complex triggers with modifiers, use Trigger() builder.
func TriggerAttr(event string) h.Attribute {
	return h.Attr("hx-trigger", event)
}

// TriggerLoad creates an hx-trigger="load" attribute.
func TriggerLoad() h.Attribute {
	return h.Attr("hx-trigger", "load")
}

// TriggerRevealed creates an hx-trigger="revealed" attribute that fires
// when the element is scrolled into view.
func TriggerRevealed() h.Attribute {
	return h.Attr("hx-trigger", "revealed")
}

// IntersectTriggerBuilder builds an hx-trigger attribute for the "intersect" event.
//
// IntersectTriggerBuilder implements h.AttrBuilder, so it can be passed directly
// to element methods.
type IntersectTriggerBuilder struct {
	spec triggerSpec
}

// Attribute returns the hx-trigger attribute.
func (b *IntersectTriggerBuilder) Attribute() h.Attribute {
	return h.Attr("hx-trigger", b.spec.String())
}

// TriggerIntersect creates a trigger that fires when the element intersects the viewport.
func TriggerIntersect() *IntersectTriggerBuilder {
	return &IntersectTriggerBuilder{spec: triggerSpec{event: "intersect"}}
}

// Root specifies a root element for intersection observation.
func (b *IntersectTriggerBuilder) Root(selector string) *IntersectTriggerBuilder {
	b.spec.modifiers = append(b.spec.modifiers, "root:"+selector)
	return b
}

// Threshold specifies the visibility threshold for intersection (0.0 to 1.0).
func (b *IntersectTriggerBuilder) Threshold(value float64) *IntersectTriggerBuilder {
	b.spec.modifiers = append(b.spec.modifiers, "threshold:"+formatFloat(value))
	return b
}

// Once makes the trigger fire only once.
func (b *IntersectTriggerBuilder) Once() *IntersectTriggerBuilder {
	b.spec.modifiers = append(b.spec.modifiers, "once")
	return b
}

// IntervalTriggerBuilder builds an hx-trigger attribute for polling triggers
// (e.g., "every 2s").
//
// IntervalTriggerBuilder implements h.AttrBuilder, so it can be passed directly
// to element methods.
type IntervalTriggerBuilder struct {
	spec triggerSpec
}

// Attribute returns the hx-trigger attribute.
func (b *IntervalTriggerBuilder) Attribute() h.Attribute {
	return h.Attr("hx-trigger", b.spec.String())
}

// TriggerEvery creates a polling trigger that fires at the specified interval.
//
// Example:
//
//	hx.TriggerEvery(2*time.Second)
//	hx.TriggerEvery(1*time.Second).Filter("document.visibilityState == 'visible'")
func TriggerEvery(interval time.Duration) *IntervalTriggerBuilder {
	return &IntervalTriggerBuilder{
		spec: triggerSpec{event: "every " + formatDuration(interval)},
	}
}

// Filter adds a JavaScript filter expression to the polling trigger.
func (b *IntervalTriggerBuilder) Filter(jsExpr string) *IntervalTriggerBuilder {
	b.spec.filter = jsExpr
	return b
}

// Once makes the polling trigger fire only once.
func (b *IntervalTriggerBuilder) Once() *IntervalTriggerBuilder {
	b.spec.modifiers = append(b.spec.modifiers, "once")
	return b
}

// Throttle limits how often the polling trigger can fire.
func (b *IntervalTriggerBuilder) Throttle(d time.Duration) *IntervalTriggerBuilder {
	b.spec.modifiers = append(b.spec.modifiers, "throttle:"+formatDuration(d))
	return b
}

// From specifies the element on which the polling event is observed.
func (b *IntervalTriggerBuilder) From(selector string) *IntervalTriggerBuilder {
	if strings.Contains(selector, " ") && !strings.HasPrefix(selector, "(") {
		b.spec.modifiers = append(b.spec.modifiers, "from:("+selector+")")
	} else {
		b.spec.modifiers = append(b.spec.modifiers, "from:"+selector)
	}
	return b
}

// FromDocument listens for the polling event on the document.
func (b *IntervalTriggerBuilder) FromDocument() *IntervalTriggerBuilder {
	b.spec.modifiers = append(b.spec.modifiers, "from:document")
	return b
}

// FromWindow listens for the polling event on the window.
func (b *IntervalTriggerBuilder) FromWindow() *IntervalTriggerBuilder {
	b.spec.modifiers = append(b.spec.modifiers, "from:window")
	return b
}

// formatFloat formats a float for HTMX attributes.
func formatFloat(f float64) string {
	// For integers, use integer format
	if f == float64(int64(f)) {
		return strconv.FormatInt(int64(f), 10)
	}
	// For decimals, use compact representation
	s := strconv.FormatFloat(f, 'f', -1, 64)
	return s
}
