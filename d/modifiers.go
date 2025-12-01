package d

import (
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

// Specifies the casing of the event
func Case(casing SignalCasing) AttrMutator {
	return AttrFunc(func(attr *attrBuilder) {
		attr.name.Grow(len(attr.name.String()) + 10)
		attr.name.WriteString("__casing.")
		attr.name.WriteString(string(casing))
	})
}

type TimingMutatorFunc func(*strings.Builder)

func NoLeading() TimingMutatorFunc {
	return func(sb *strings.Builder) {
		sb.WriteString(".no-leading")
	}
}

func NoTrailing() TimingMutatorFunc {
	return func(sb *strings.Builder) {
		sb.WriteString(".no-trailing")
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
