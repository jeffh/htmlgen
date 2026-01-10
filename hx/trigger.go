package hx

import (
	"strconv"
	"strings"
	"time"

	"github.com/jeffh/htmlgen/h"
)

// TriggerMod is a modifier for trigger behavior.
type TriggerMod interface {
	applyTrigger(*triggerSpec)
}

type triggerModFunc func(*triggerSpec)

func (f triggerModFunc) applyTrigger(s *triggerSpec) { f(s) }

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
type TriggerBuilder struct {
	triggers []triggerSpec
}

// Trigger creates a new trigger with the specified event and optional modifiers.
//
// Example:
//
//	hx.Trigger("click")
//	hx.Trigger("keyup", hx.Changed(), hx.Delay(500*time.Millisecond))
//	hx.Trigger("click", hx.From("#other-element"))
func Trigger(event string, mods ...TriggerMod) *TriggerBuilder {
	spec := triggerSpec{event: event}
	for _, mod := range mods {
		mod.applyTrigger(&spec)
	}
	return &TriggerBuilder{triggers: []triggerSpec{spec}}
}

// And adds another trigger to the builder.
//
// Example:
//
//	hx.Trigger("load").And("click", hx.Delay(1*time.Second))
func (t *TriggerBuilder) And(event string, mods ...TriggerMod) *TriggerBuilder {
	spec := triggerSpec{event: event}
	for _, mod := range mods {
		mod.applyTrigger(&spec)
	}
	t.triggers = append(t.triggers, spec)
	return t
}

// Attr returns the hx-trigger attribute.
func (t *TriggerBuilder) Attr() h.Attribute {
	parts := make([]string, len(t.triggers))
	for i, trigger := range t.triggers {
		parts[i] = trigger.String()
	}
	return h.Attr("hx-trigger", strings.Join(parts, ", "))
}

// TriggerAttr creates a simple hx-trigger attribute with a single event.
// For complex triggers with modifiers, use Trigger() builder.
func TriggerAttr(event string) h.Attribute {
	return h.Attr("hx-trigger", event)
}

// Once makes the trigger fire only once.
func Once() TriggerMod {
	return triggerModFunc(func(s *triggerSpec) {
		s.modifiers = append(s.modifiers, "once")
	})
}

// Changed makes the trigger fire only when the element's value has changed.
func Changed() TriggerMod {
	return triggerModFunc(func(s *triggerSpec) {
		s.modifiers = append(s.modifiers, "changed")
	})
}

// Delay adds a delay before the trigger fires.
// If another event occurs during the delay, the timer resets.
func Delay(d time.Duration) TriggerMod {
	return triggerModFunc(func(s *triggerSpec) {
		s.modifiers = append(s.modifiers, "delay:"+formatDuration(d))
	})
}

// Throttle limits how often the trigger can fire.
func Throttle(d time.Duration) TriggerMod {
	return triggerModFunc(func(s *triggerSpec) {
		s.modifiers = append(s.modifiers, "throttle:"+formatDuration(d))
	})
}

// From specifies that the event should be listened for on a different element.
// Supports extended CSS selectors: document, window, closest <sel>, find <sel>, next, previous.
func From(selector string) TriggerMod {
	return triggerModFunc(func(s *triggerSpec) {
		// Handle selectors with spaces by wrapping in parentheses
		if strings.Contains(selector, " ") && !strings.HasPrefix(selector, "(") {
			s.modifiers = append(s.modifiers, "from:("+selector+")")
		} else {
			s.modifiers = append(s.modifiers, "from:"+selector)
		}
	})
}

// FromDocument listens for the event on the document.
func FromDocument() TriggerMod {
	return triggerModFunc(func(s *triggerSpec) {
		s.modifiers = append(s.modifiers, "from:document")
	})
}

// FromWindow listens for the event on the window.
func FromWindow() TriggerMod {
	return triggerModFunc(func(s *triggerSpec) {
		s.modifiers = append(s.modifiers, "from:window")
	})
}

// FromClosest listens for the event on the closest ancestor matching the selector.
func FromClosest(selector string) TriggerMod {
	return triggerModFunc(func(s *triggerSpec) {
		s.modifiers = append(s.modifiers, "from:closest "+selector)
	})
}

// FromFind listens for the event on the first descendant matching the selector.
func FromFind(selector string) TriggerMod {
	return triggerModFunc(func(s *triggerSpec) {
		s.modifiers = append(s.modifiers, "from:find "+selector)
	})
}

// FromNext listens for the event on the next sibling element.
func FromNext(selector ...string) TriggerMod {
	return triggerModFunc(func(s *triggerSpec) {
		if len(selector) > 0 {
			s.modifiers = append(s.modifiers, "from:next "+selector[0])
		} else {
			s.modifiers = append(s.modifiers, "from:next")
		}
	})
}

// FromPrevious listens for the event on the previous sibling element.
func FromPrevious(selector ...string) TriggerMod {
	return triggerModFunc(func(s *triggerSpec) {
		if len(selector) > 0 {
			s.modifiers = append(s.modifiers, "from:previous "+selector[0])
		} else {
			s.modifiers = append(s.modifiers, "from:previous")
		}
	})
}

// TriggerTarget filters events to only those whose target matches the selector.
// Note: This is different from hx-target; it filters by event.target.
func TriggerTarget(selector string) TriggerMod {
	return triggerModFunc(func(s *triggerSpec) {
		s.modifiers = append(s.modifiers, "target:"+selector)
	})
}

// Consume prevents the event from triggering requests on parent elements.
func Consume() TriggerMod {
	return triggerModFunc(func(s *triggerSpec) {
		s.modifiers = append(s.modifiers, "consume")
	})
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
func Queue(mode QueueMode) TriggerMod {
	return triggerModFunc(func(s *triggerSpec) {
		s.modifiers = append(s.modifiers, "queue:"+string(mode))
	})
}

// Filter adds a JavaScript filter expression to the trigger.
// The expression should evaluate to true/false.
//
// Example:
//
//	hx.Trigger("click", hx.Filter("ctrlKey"))
//	hx.Trigger("keyup", hx.Filter("keyCode==13"))
func Filter(jsExpr string) TriggerMod {
	return triggerModFunc(func(s *triggerSpec) {
		s.filter = jsExpr
	})
}

// Special trigger constructors

// TriggerLoad creates a trigger that fires on page load.
func TriggerLoad(mods ...TriggerMod) *TriggerBuilder {
	return Trigger("load", mods...)
}

// TriggerRevealed creates a trigger that fires when the element is scrolled into view.
func TriggerRevealed(mods ...TriggerMod) *TriggerBuilder {
	return Trigger("revealed", mods...)
}

// TriggerIntersect creates a trigger that fires when the element intersects the viewport.
// Use IntersectRoot and IntersectThreshold modifiers for configuration.
func TriggerIntersect(mods ...TriggerMod) *TriggerBuilder {
	return Trigger("intersect", mods...)
}

// IntersectRoot specifies a root element for intersection observation.
func IntersectRoot(selector string) TriggerMod {
	return triggerModFunc(func(s *triggerSpec) {
		s.modifiers = append(s.modifiers, "root:"+selector)
	})
}

// IntersectThreshold specifies the visibility threshold for intersection (0.0 to 1.0).
func IntersectThreshold(threshold float64) TriggerMod {
	return triggerModFunc(func(s *triggerSpec) {
		s.modifiers = append(s.modifiers, "threshold:"+formatFloat(threshold))
	})
}

// TriggerEvery creates a polling trigger that fires at the specified interval.
//
// Example:
//
//	hx.TriggerEvery(2*time.Second)
//	hx.TriggerEvery(1*time.Second, hx.Filter("document.visibilityState == 'visible'"))
func TriggerEvery(interval time.Duration, mods ...TriggerMod) *TriggerBuilder {
	spec := triggerSpec{event: "every " + formatDuration(interval)}
	for _, mod := range mods {
		mod.applyTrigger(&spec)
	}
	return &TriggerBuilder{triggers: []triggerSpec{spec}}
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
