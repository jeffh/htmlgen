package d

import (
	"strconv"
	"strings"
	"time"
)

// Prevents the default event behavior.
func PreventDefault() AttrMutator {
	return AttrFunc(func(attr *attrBuilder) {
		attr.name.WriteString("__prevent")
	})
}

// Prevents the event from being triggered multiple times.
// Only works for built-in events.
func Once() AttrMutator {
	return AttrFunc(func(attr *attrBuilder) {
		attr.name.WriteString("__once")
	})
}

// Specifies that the event should not be canceled.
// Only works for built-in events.
func Passive() AttrMutator {
	return AttrFunc(func(attr *attrBuilder) {
		attr.name.WriteString("__passive")
	})
}

// Specifies that the event should not be captured.
// Only works for built-in events.
func Capture() AttrMutator {
	return AttrFunc(func(attr *attrBuilder) {
		attr.name.WriteString("__capture")
	})
}

type SignalCasing string

const (
	CamelCase  SignalCasing = "camel"  // myEvent
	KebabCase  SignalCasing = "kebab"  // my-event
	SnakeCase  SignalCasing = "snake"  // my_event
	PascalCase SignalCasing = "pascal" // MyEvent
)

// Specifies the casing of the signal/event name
func Case(casing SignalCasing) AttrMutator {
	return AttrFunc(func(attr *attrBuilder) {
		attr.name.Grow(len(attr.name.String()) + 10)
		attr.name.WriteString("__case.")
		attr.name.WriteString(string(casing))
	})
}

type TimingMutatorFunc func(*strings.Builder)

// NoLeading prevents the first trigger from firing immediately (for debounce)
func NoLeading() TimingMutatorFunc {
	return func(sb *strings.Builder) {
		sb.WriteString(".noleading")
	}
}

// NoTrailing prevents the final trigger from firing after the delay (for debounce)
func NoTrailing() TimingMutatorFunc {
	return func(sb *strings.Builder) {
		sb.WriteString(".notrailing")
	}
}

// Leading causes the first trigger to fire immediately (for throttle)
func Leading() TimingMutatorFunc {
	return func(sb *strings.Builder) {
		sb.WriteString(".leading")
	}
}

// Trailing causes the final trigger to fire after the delay (for throttle)
func Trailing() TimingMutatorFunc {
	return func(sb *strings.Builder) {
		sb.WriteString(".trailing")
	}
}

// Specifies the delay of the event triggered
func Delay(d time.Duration) AttrMutator {
	return AttrFunc(func(attr *attrBuilder) {
		delayStr := d.String()
		attr.name.Grow(len(attr.name.String()) + len(delayStr) + 10)
		attr.name.WriteString("__delay.")
		attr.name.WriteString(delayStr)
	})
}

// Debounces the events triggered
func Debounce(d time.Duration, opts ...TimingMutatorFunc) AttrMutator {
	return AttrFunc(func(attr *attrBuilder) {
		delayStr := d.String()
		attr.name.Grow(len(attr.name.String()) + len(delayStr) + 10)
		attr.name.WriteString("__debounce.")
		attr.name.WriteString(delayStr)
		for _, opt := range opts {
			opt(&attr.name)
		}
	})
}

// Throttles the events triggered
func Throttle(d time.Duration, opts ...TimingMutatorFunc) AttrMutator {
	return AttrFunc(func(attr *attrBuilder) {
		delayStr := d.String()
		attr.name.Grow(len(attr.name.String()) + len(delayStr) + 10)
		attr.name.WriteString("__throttle.")
		attr.name.WriteString(delayStr)
		for _, opt := range opts {
			opt(&attr.name)
		}
	})
}

// Wraps expression in document.startViewTransition() when View Transition API is supported
func ViewTransition() AttrMutator {
	return AttrFunc(func(attr *attrBuilder) {
		attr.name.WriteString("__viewtransition")
	})
}

// Attaches event listener to the window element
func Window() AttrMutator {
	return AttrFunc(func(attr *attrBuilder) {
		attr.name.WriteString("__window")
	})
}

// Triggers when the event is triggered outside of the element
func Outside() AttrMutator {
	return AttrFunc(func(attr *attrBuilder) {
		attr.name.WriteString("__outside")
	})
}

// Stops the event from propagating
func StopPropagation() AttrMutator {
	return AttrFunc(func(attr *attrBuilder) {
		attr.name.WriteString("__stop")
	})
}

// IfMissing only patches signals if the keys don't already exist.
// Used with Signal/Signals attributes.
func IfMissing() AttrMutator {
	return AttrFunc(func(attr *attrBuilder) {
		attr.name.WriteString("__ifmissing")
	})
}

// Self modifier for data-ignore - ignores only the element, not descendants.
func Self() AttrMutator {
	return AttrFunc(func(attr *attrBuilder) {
		attr.name.WriteString("__self")
	})
}

// Terse outputs compact JSON without whitespace for data-json-signals debug.
func Terse() AttrMutator {
	return AttrFunc(func(attr *attrBuilder) {
		attr.name.WriteString("__terse")
	})
}

// Half triggers intersection at 50% element visibility.
func Half() AttrMutator {
	return AttrFunc(func(attr *attrBuilder) {
		attr.name.WriteString("__half")
	})
}

// Full triggers intersection at 100% element visibility.
func Full() AttrMutator {
	return AttrFunc(func(attr *attrBuilder) {
		attr.name.WriteString("__full")
	})
}

// Exit triggers when element exits viewport instead of entering.
// Used with data-on-intersect to detect when an element leaves view.
func Exit() AttrMutator {
	return AttrFunc(func(attr *attrBuilder) {
		attr.name.WriteString("__exit")
	})
}

// Threshold sets a custom visibility threshold (0.0-1.0) for intersection.
// Use instead of Half() or Full() for precise control.
// Example: Threshold(0.25) triggers at 25% visibility.
func Threshold(value float64) AttrMutator {
	return AttrFunc(func(attr *attrBuilder) {
		valueStr := strconv.FormatFloat(value, 'f', -1, 64)
		attr.name.WriteString("__threshold.")
		attr.name.WriteString(valueStr)
	})
}

type DurationMutatorFunc func(*strings.Builder)

// DurationLeading causes the first interval to fire immediately.
func DurationLeading() DurationMutatorFunc {
	return func(sb *strings.Builder) {
		sb.WriteString(".leading")
	}
}

// Duration sets the interval timing for data-on-interval.
// Default is 1 second if not specified.
// Example: Duration(500*time.Millisecond) produces __duration.500ms
// Example: Duration(2*time.Second, DurationLeading()) produces __duration.2s.leading
func Duration(d time.Duration, opts ...DurationMutatorFunc) AttrMutator {
	return AttrFunc(func(attr *attrBuilder) {
		durationStr := d.String()
		attr.name.Grow(len(attr.name.String()) + len(durationStr) + 15)
		attr.name.WriteString("__duration.")
		attr.name.WriteString(durationStr)
		for _, opt := range opts {
			opt(&attr.name)
		}
	})
}
