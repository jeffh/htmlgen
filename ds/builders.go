package ds

import (
	"strconv"
	"strings"
	"time"

	"github.com/jeffh/htmlgen/h"
	"github.com/jeffh/htmlgen/js"
)

// attrBase is the shared payload for fluent Datastar attribute builders.
// It tracks the accumulating attribute name (which modifiers append to) and
// the JavaScript statements that form the attribute value.
type attrBase struct {
	name       strings.Builder
	statements []string
}

func (a *attrBase) isTagArg() {}

// Attribute returns the finished h.Attribute. Statements are joined with
// "; " into the attribute value.
func (a *attrBase) Attribute() h.Attribute {
	return h.Attr(a.name.String(), strings.Join(a.statements, "; "))
}

func (a *attrBase) addStmt(s string) {
	a.statements = append(a.statements, s)
}

func (a *attrBase) addValue(v Value) {
	a.statements = append(a.statements, js.ToJS(v.expr))
}

func newAttr(name string) *attrBase {
	var b attrBase
	b.name.WriteString(name)
	return &b
}

// SignalCasing controls the casing of signal/event names emitted by Datastar.
type SignalCasing string

const (
	CamelCase  SignalCasing = "camel"  // myEvent
	KebabCase  SignalCasing = "kebab"  // my-event
	SnakeCase  SignalCasing = "snake"  // my_event
	PascalCase SignalCasing = "pascal" // MyEvent
)

// TimingOption modifies a Debounce/Throttle/Duration name segment.
type TimingOption func(*strings.Builder)

// Leading causes the first trigger to also fire immediately (for Debounce).
// Debounce only fires on the trailing edge by default.
func Leading() TimingOption {
	return func(sb *strings.Builder) { sb.WriteString(".leading") }
}

// NoTrailing prevents the final trigger from firing after the delay (for Debounce).
func NoTrailing() TimingOption {
	return func(sb *strings.Builder) { sb.WriteString(".notrailing") }
}

// NoLeading prevents the first trigger from firing immediately (for Throttle).
// Throttle fires on the leading edge by default.
func NoLeading() TimingOption {
	return func(sb *strings.Builder) { sb.WriteString(".noleading") }
}

// Trailing causes the final trigger to also fire after the delay (for Throttle).
// Throttle only fires on the leading edge by default.
func Trailing() TimingOption {
	return func(sb *strings.Builder) { sb.WriteString(".trailing") }
}

// DurationLeading causes the first interval to fire immediately (for Duration).
func DurationLeading() TimingOption {
	return func(sb *strings.Builder) { sb.WriteString(".leading") }
}

func writeTiming(sb *strings.Builder, segment string, d time.Duration, opts []TimingOption) {
	sb.WriteString(segment)
	sb.WriteString(d.String())
	for _, opt := range opts {
		opt(sb)
	}
}

// EventBuilder builds a Datastar event-handler attribute (data-on:click,
// data-on:submit, data-on:custom, etc.) with chainable modifiers.
type EventBuilder struct {
	*attrBase
}

func newEventBuilder(name string, actions []Value) *EventBuilder {
	b := &EventBuilder{attrBase: newAttr(name)}
	for _, a := range actions {
		b.addValue(a)
	}
	return b
}

// Then appends additional JavaScript actions to be executed in order with the
// initial actions.
func (b *EventBuilder) Then(actions ...Value) *EventBuilder {
	for _, a := range actions {
		b.addValue(a)
	}
	return b
}

// PreventDefault appends "__prevent" — calls event.preventDefault().
func (b *EventBuilder) PreventDefault() *EventBuilder {
	b.name.WriteString("__prevent")
	return b
}

// StopPropagation appends "__stop" — stops event propagation.
func (b *EventBuilder) StopPropagation() *EventBuilder {
	b.name.WriteString("__stop")
	return b
}

// Once appends "__once" — fire only once.
func (b *EventBuilder) Once() *EventBuilder {
	b.name.WriteString("__once")
	return b
}

// Passive appends "__passive" — declares a passive listener.
func (b *EventBuilder) Passive() *EventBuilder {
	b.name.WriteString("__passive")
	return b
}

// Capture appends "__capture" — declares a capture-phase listener.
func (b *EventBuilder) Capture() *EventBuilder {
	b.name.WriteString("__capture")
	return b
}

// Outside appends "__outside" — fires when the event happens outside the element.
func (b *EventBuilder) Outside() *EventBuilder {
	b.name.WriteString("__outside")
	return b
}

// Window appends "__window" — attaches the listener to the window.
func (b *EventBuilder) Window() *EventBuilder {
	b.name.WriteString("__window")
	return b
}

// Document appends "__document" — attaches the listener to the document.
// Useful for events that are only available on document and do not bubble.
func (b *EventBuilder) Document() *EventBuilder {
	b.name.WriteString("__document")
	return b
}

// ViewTransition appends "__viewtransition" — wraps work in document.startViewTransition().
func (b *EventBuilder) ViewTransition() *EventBuilder {
	b.name.WriteString("__viewtransition")
	return b
}

// Case appends "__case.<casing>" — converts the event name casing.
func (b *EventBuilder) Case(c SignalCasing) *EventBuilder {
	b.name.WriteString("__case.")
	b.name.WriteString(string(c))
	return b
}

// Delay appends "__delay.<duration>".
func (b *EventBuilder) Delay(d time.Duration) *EventBuilder {
	writeTiming(&b.name, "__delay.", d, nil)
	return b
}

// Debounce appends "__debounce.<duration>" plus optional Leading/NoTrailing.
func (b *EventBuilder) Debounce(d time.Duration, opts ...TimingOption) *EventBuilder {
	writeTiming(&b.name, "__debounce.", d, opts)
	return b
}

// Throttle appends "__throttle.<duration>" plus optional NoLeading/Trailing.
func (b *EventBuilder) Throttle(d time.Duration, opts ...TimingOption) *EventBuilder {
	writeTiming(&b.name, "__throttle.", d, opts)
	return b
}

// IntersectBuilder builds a data-on-intersect attribute.
type IntersectBuilder struct {
	*attrBase
}

func newIntersectBuilder(actions []Value) *IntersectBuilder {
	b := &IntersectBuilder{attrBase: newAttr("data-on-intersect")}
	for _, a := range actions {
		b.addValue(a)
	}
	return b
}

// Once appends "__once".
func (b *IntersectBuilder) Once() *IntersectBuilder {
	b.name.WriteString("__once")
	return b
}

// Half appends "__half" — fire at 50% visibility.
func (b *IntersectBuilder) Half() *IntersectBuilder {
	b.name.WriteString("__half")
	return b
}

// Full appends "__full" — fire at 100% visibility.
func (b *IntersectBuilder) Full() *IntersectBuilder {
	b.name.WriteString("__full")
	return b
}

// Exit appends "__exit" — fire when leaving the viewport.
func (b *IntersectBuilder) Exit() *IntersectBuilder {
	b.name.WriteString("__exit")
	return b
}

// Threshold appends "__threshold.<percent>" — custom intersection threshold as
// a percentage of visibility (0–100), e.g. Threshold(25) => "__threshold.25".
func (b *IntersectBuilder) Threshold(percent int) *IntersectBuilder {
	b.name.WriteString("__threshold.")
	b.name.WriteString(strconv.Itoa(percent))
	return b
}

// Delay appends "__delay.<duration>".
func (b *IntersectBuilder) Delay(d time.Duration) *IntersectBuilder {
	writeTiming(&b.name, "__delay.", d, nil)
	return b
}

// Debounce appends "__debounce.<duration>" plus optional Leading/NoTrailing.
func (b *IntersectBuilder) Debounce(d time.Duration, opts ...TimingOption) *IntersectBuilder {
	writeTiming(&b.name, "__debounce.", d, opts)
	return b
}

// Throttle appends "__throttle.<duration>" plus optional NoLeading/Trailing.
func (b *IntersectBuilder) Throttle(d time.Duration, opts ...TimingOption) *IntersectBuilder {
	writeTiming(&b.name, "__throttle.", d, opts)
	return b
}

// ViewTransition appends "__viewtransition" — wraps work in document.startViewTransition().
func (b *IntersectBuilder) ViewTransition() *IntersectBuilder {
	b.name.WriteString("__viewtransition")
	return b
}

// IntervalBuilder builds a data-on-interval attribute.
type IntervalBuilder struct {
	*attrBase
}

func newIntervalBuilder(actions []Value) *IntervalBuilder {
	b := &IntervalBuilder{attrBase: newAttr("data-on-interval")}
	for _, a := range actions {
		b.addValue(a)
	}
	return b
}

// Duration appends "__duration.<duration>" plus optional DurationLeading.
func (b *IntervalBuilder) Duration(d time.Duration, opts ...TimingOption) *IntervalBuilder {
	writeTiming(&b.name, "__duration.", d, opts)
	return b
}

// ViewTransition appends "__viewtransition" — wraps work in document.startViewTransition().
func (b *IntervalBuilder) ViewTransition() *IntervalBuilder {
	b.name.WriteString("__viewtransition")
	return b
}

// SignalPatchBuilder builds a data-on-signal-patch attribute.
type SignalPatchBuilder struct {
	*attrBase
}

func newSignalPatchBuilder(actions []Value) *SignalPatchBuilder {
	b := &SignalPatchBuilder{attrBase: newAttr("data-on-signal-patch")}
	for _, a := range actions {
		b.addValue(a)
	}
	return b
}

// Delay appends "__delay.<duration>".
func (b *SignalPatchBuilder) Delay(d time.Duration) *SignalPatchBuilder {
	writeTiming(&b.name, "__delay.", d, nil)
	return b
}

// Debounce appends "__debounce.<duration>" plus optional Leading/NoTrailing.
func (b *SignalPatchBuilder) Debounce(d time.Duration, opts ...TimingOption) *SignalPatchBuilder {
	writeTiming(&b.name, "__debounce.", d, opts)
	return b
}

// Throttle appends "__throttle.<duration>" plus optional NoLeading/Trailing.
func (b *SignalPatchBuilder) Throttle(d time.Duration, opts ...TimingOption) *SignalPatchBuilder {
	writeTiming(&b.name, "__throttle.", d, opts)
	return b
}

// NamedBuilder builds attributes whose name carries a signal name and which
// support the Case modifier (data-indicator, data-indicator:key,
// data-ref:name, data-computed:name, data-signals:name, data-class:name).
type NamedBuilder struct {
	*attrBase
}

// Case appends "__case.<casing>".
func (b *NamedBuilder) Case(c SignalCasing) *NamedBuilder {
	b.name.WriteString("__case.")
	b.name.WriteString(string(c))
	return b
}

// BindBuilder builds data-bind / data-bind:key attributes. It supports the
// Case, Prop and Event modifiers.
type BindBuilder struct {
	*attrBase
}

// Case appends "__case.<casing>".
func (b *BindBuilder) Case(c SignalCasing) *BindBuilder {
	b.name.WriteString("__case.")
	b.name.WriteString(string(c))
	return b
}

// Prop appends "__prop.<name>" — binds to a specific element property instead
// of the default binding. Must not be a read-only property.
//
//	ds.BindKey("is-checked").Prop("checked")  =>  data-bind:is-checked__prop.checked
func (b *BindBuilder) Prop(name string) *BindBuilder {
	b.name.WriteString("__prop.")
	b.name.WriteString(name)
	return b
}

// Event appends "__event.<events...>" — defines which events sync the element
// property back to the signal.
//
//	ds.BindKey("query").Event("input", "change")  =>  data-bind:query__event.input.change
func (b *BindBuilder) Event(events ...string) *BindBuilder {
	b.name.WriteString("__event")
	for _, e := range events {
		b.name.WriteString(".")
		b.name.WriteString(e)
	}
	return b
}

// InitBuilder builds the data-init attribute. It supports the Delay and
// ViewTransition modifiers.
type InitBuilder struct {
	*attrBase
}

// Delay appends "__delay.<duration>" — waits before running the expression.
func (b *InitBuilder) Delay(d time.Duration) *InitBuilder {
	writeTiming(&b.name, "__delay.", d, nil)
	return b
}

// ViewTransition appends "__viewtransition" — wraps work in document.startViewTransition().
func (b *InitBuilder) ViewTransition() *InitBuilder {
	b.name.WriteString("__viewtransition")
	return b
}

// SignalsBuilder builds the data-signals attribute (whole-object form). It
// supports IfMissing and Case.
type SignalsBuilder struct {
	*attrBase
}

// Case appends "__case.<casing>".
func (b *SignalsBuilder) Case(c SignalCasing) *SignalsBuilder {
	b.name.WriteString("__case.")
	b.name.WriteString(string(c))
	return b
}

// IfMissing appends "__ifmissing" — only patch missing keys.
func (b *SignalsBuilder) IfMissing() *SignalsBuilder {
	b.name.WriteString("__ifmissing")
	return b
}

// JsonSignalsBuilder builds the data-json-signals debug attribute. It supports
// the Terse modifier.
type JsonSignalsBuilder struct {
	*attrBase
}

// Terse appends "__terse" — outputs a more compact JSON representation.
func (b *JsonSignalsBuilder) Terse() *JsonSignalsBuilder {
	b.name.WriteString("__terse")
	return b
}
