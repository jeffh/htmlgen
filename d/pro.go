// Pro Datastar attributes require a commercial license from https://data-star.dev/
// These attributes provide additional functionality beyond the free tier.
package d

import (
	"strings"

	"github.com/jeffh/htmlgen/h"
)

// Animate enables reactive animations on element attributes.
// Requires Datastar Pro license.
// Produces: data-animate
func Animate(options ...AttrMutator) h.Attribute {
	return exprAttr("data-animate", options...)
}

// CustomValidity adds custom validation messages to form inputs.
// Empty strings indicate valid; non-empty strings are shown as validation errors.
// Requires Datastar Pro license.
// Example: CustomValidity(Raw("$password === $confirmPassword ? ” : 'Passwords must match'"))
// Produces: data-custom-validity="$password === $confirmPassword ? ” : 'Passwords must match'"
func CustomValidity(expression ...AttrMutator) h.Attribute {
	return exprAttr("data-custom-validity", expression...)
}

// OnRAF executes an expression on every requestAnimationFrame event.
// Requires Datastar Pro license.
// Example: OnRAF(Throttle(100*time.Millisecond), Raw("$frameCount++"))
// Produces: data-on-raf__throttle.100ms="$frameCount++"
func OnRAF(options ...AttrMutator) h.Attribute {
	return exprAttr("data-on-raf", options...)
}

// OnResize runs an expression when element dimensions change.
// Requires Datastar Pro license.
// Example: OnResize(Debounce(200*time.Millisecond), Raw("$width = el.offsetWidth"))
// Produces: data-on-resize__debounce.200ms="$width = el.offsetWidth"
func OnResize(options ...AttrMutator) h.Attribute {
	return exprAttr("data-on-resize", options...)
}

// Persist persists signals in local or session storage.
// Use Session() modifier for session storage instead of local storage.
// Requires Datastar Pro license.
// Example: Persist(nil) or Persist(&FilterOptions{IncludeReg: ptr("user")})
// Produces: data-persist or data-persist="{include: /user/}"
func Persist(options *FilterOptions, modifiers ...AttrMutator) h.Attribute {
	if options == nil && len(modifiers) == 0 {
		return h.Attr("data-persist", "")
	}
	if options == nil {
		return exprAttr("data-persist", modifiers...)
	}
	opts := append([]AttrMutator{FilterOptionsValue(options)}, modifiers...)
	return exprAttr("data-persist", opts...)
}

// PersistKey persists signals using a custom storage key.
// Requires Datastar Pro license.
// Example: PersistKey("mykey", Session())
// Produces: data-persist:mykey__session
func PersistKey(key string, modifiers ...AttrMutator) h.Attribute {
	opts := append([]AttrMutator{appendName(key)}, modifiers...)
	return exprAttr("data-persist:", opts...)
}

// QueryString syncs query parameters to/from signal values.
// Use Filter() to filter empty values, History() for browser history integration.
// Requires Datastar Pro license.
// Example: QueryString(nil) or QueryString(&FilterOptions{IncludeReg: ptr("search")})
// Produces: data-query-string or data-query-string="{include: /search/}"
func QueryString(options *FilterOptions, modifiers ...AttrMutator) h.Attribute {
	if options == nil && len(modifiers) == 0 {
		return h.Attr("data-query-string", "")
	}
	if options == nil {
		return exprAttr("data-query-string", modifiers...)
	}
	opts := append([]AttrMutator{FilterOptionsValue(options)}, modifiers...)
	return exprAttr("data-query-string", opts...)
}

// ReplaceURL replaces the browser URL without page reload.
// Accepts relative or absolute URLs as evaluated expressions.
// Requires Datastar Pro license.
// Example: ReplaceURL(Raw("`/page${$page}`"))
// Produces: data-replace-url="`/page${$page}`"
func ReplaceURL(expression ...AttrMutator) h.Attribute {
	return exprAttr("data-replace-url", expression...)
}

// ScrollIntoView scrolls the element into viewport view.
// Use scroll behavior and alignment modifiers.
// Requires Datastar Pro license.
// Example: ScrollIntoView(Smooth(), VCenter())
// Produces: data-scroll-into-view__smooth__vcenter
func ScrollIntoView(options ...AttrMutator) h.Attribute {
	return exprAttr("data-scroll-into-view", options...)
}

// ViewTransitionName sets explicit view-transition-name for CSS animations.
// Requires Datastar Pro license.
// Example: ViewTransitionName(Raw("$itemId"))
// Produces: data-view-transition="$itemId"
func ViewTransitionName(expression ...AttrMutator) h.Attribute {
	return exprAttr("data-view-transition", expression...)
}

// Pro modifiers

// Session uses session storage instead of local storage for Persist.
func Session() AttrMutator {
	return AttrFunc(func(attr *attrBuilder) {
		attr.name.WriteString("__session")
	})
}

// Filter filters out empty values for QueryString.
func Filter() AttrMutator {
	return AttrFunc(func(attr *attrBuilder) {
		attr.name.WriteString("__filter")
	})
}

// History enables browser history integration for QueryString.
func History() AttrMutator {
	return AttrFunc(func(attr *attrBuilder) {
		attr.name.WriteString("__history")
	})
}

// Scroll behavior modifiers for ScrollIntoView

// Smooth enables smooth scrolling behavior.
func Smooth() AttrMutator {
	return AttrFunc(func(attr *attrBuilder) {
		attr.name.WriteString("__smooth")
	})
}

// Instant enables instant scrolling behavior.
func Instant() AttrMutator {
	return AttrFunc(func(attr *attrBuilder) {
		attr.name.WriteString("__instant")
	})
}

// Auto uses browser default scrolling behavior.
func Auto() AttrMutator {
	return AttrFunc(func(attr *attrBuilder) {
		attr.name.WriteString("__auto")
	})
}

// Horizontal alignment modifiers for ScrollIntoView

// HStart aligns element to start of horizontal viewport.
func HStart() AttrMutator {
	return AttrFunc(func(attr *attrBuilder) {
		attr.name.WriteString("__hstart")
	})
}

// HCenter aligns element to center of horizontal viewport.
func HCenter() AttrMutator {
	return AttrFunc(func(attr *attrBuilder) {
		attr.name.WriteString("__hcenter")
	})
}

// HEnd aligns element to end of horizontal viewport.
func HEnd() AttrMutator {
	return AttrFunc(func(attr *attrBuilder) {
		attr.name.WriteString("__hend")
	})
}

// HNearest aligns element to nearest edge of horizontal viewport.
func HNearest() AttrMutator {
	return AttrFunc(func(attr *attrBuilder) {
		attr.name.WriteString("__hnearest")
	})
}

// Vertical alignment modifiers for ScrollIntoView

// VStart aligns element to start of vertical viewport.
func VStart() AttrMutator {
	return AttrFunc(func(attr *attrBuilder) {
		attr.name.WriteString("__vstart")
	})
}

// VCenter aligns element to center of vertical viewport.
func VCenter() AttrMutator {
	return AttrFunc(func(attr *attrBuilder) {
		attr.name.WriteString("__vcenter")
	})
}

// VEnd aligns element to end of vertical viewport.
func VEnd() AttrMutator {
	return AttrFunc(func(attr *attrBuilder) {
		attr.name.WriteString("__vend")
	})
}

// VNearest aligns element to nearest edge of vertical viewport.
func VNearest() AttrMutator {
	return AttrFunc(func(attr *attrBuilder) {
		attr.name.WriteString("__vnearest")
	})
}

// Focus focuses the element after scrolling into view.
func Focus() AttrMutator {
	return AttrFunc(func(attr *attrBuilder) {
		attr.name.WriteString("__focus")
	})
}

// Pro Actions

// Clipboard copies text to the clipboard.
// Requires Datastar Pro license.
// Example: OnClick(Clipboard(Str("Hello, world!")))
// Produces: data-on:click="@clipboard('Hello, world!')"
func Clipboard(text AttrValueAppender) ValueMutator {
	return ValueMutator{
		AttrMutator: AttrFunc(func(attr *attrBuilder) {
			var sb strings.Builder
			sb.WriteString("@clipboard(")
			text.Append(&sb)
			sb.WriteString(")")
			attr.AppendStatement(sb.String())
		}),
		AttrValueAppender: AttrValueFunc(func(sb *strings.Builder) {
			sb.WriteString("@clipboard(")
			text.Append(sb)
			sb.WriteString(")")
		}),
	}
}

// ClipboardBase64 copies Base64-decoded text to the clipboard.
// Useful for content with special characters, quotes, or code fragments.
// Requires Datastar Pro license.
// Example: OnClick(ClipboardBase64(Str("SGVsbG8sIHdvcmxkIQ==")))
// Produces: data-on:click="@clipboard('SGVsbG8sIHdvcmxkIQ==', true)"
func ClipboardBase64(text AttrValueAppender) ValueMutator {
	return ValueMutator{
		AttrMutator: AttrFunc(func(attr *attrBuilder) {
			var sb strings.Builder
			sb.WriteString("@clipboard(")
			text.Append(&sb)
			sb.WriteString(", true)")
			attr.AppendStatement(sb.String())
		}),
		AttrValueAppender: AttrValueFunc(func(sb *strings.Builder) {
			sb.WriteString("@clipboard(")
			text.Append(sb)
			sb.WriteString(", true)")
		}),
	}
}

// Fit linearly interpolates a value from one range to another.
// Requires Datastar Pro license.
// Syntax: @fit(v, oldMin, oldMax, newMin, newMax)
// Example: Computed("rgb", Fit(Raw("$slider"), Raw("0"), Raw("100"), Raw("0"), Raw("255")))
func Fit(v, oldMin, oldMax, newMin, newMax AttrValueAppender) ValueMutator {
	return fitAction(v, oldMin, oldMax, newMin, newMax, false, false)
}

// FitClamped linearly interpolates with clamping to keep results within the target range.
// Requires Datastar Pro license.
// Syntax: @fit(v, oldMin, oldMax, newMin, newMax, true)
func FitClamped(v, oldMin, oldMax, newMin, newMax AttrValueAppender) ValueMutator {
	return fitAction(v, oldMin, oldMax, newMin, newMax, true, false)
}

// FitRounded linearly interpolates with rounding to nearest integer.
// Requires Datastar Pro license.
// Syntax: @fit(v, oldMin, oldMax, newMin, newMax, false, true)
func FitRounded(v, oldMin, oldMax, newMin, newMax AttrValueAppender) ValueMutator {
	return fitAction(v, oldMin, oldMax, newMin, newMax, false, true)
}

// FitClampedRounded linearly interpolates with both clamping and rounding.
// Requires Datastar Pro license.
// Syntax: @fit(v, oldMin, oldMax, newMin, newMax, true, true)
func FitClampedRounded(v, oldMin, oldMax, newMin, newMax AttrValueAppender) ValueMutator {
	return fitAction(v, oldMin, oldMax, newMin, newMax, true, true)
}

func fitAction(v, oldMin, oldMax, newMin, newMax AttrValueAppender, clamp, round bool) ValueMutator {
	return ValueMutator{
		AttrMutator: AttrFunc(func(attr *attrBuilder) {
			var sb strings.Builder
			writeFit(&sb, v, oldMin, oldMax, newMin, newMax, clamp, round)
			attr.AppendStatement(sb.String())
		}),
		AttrValueAppender: AttrValueFunc(func(sb *strings.Builder) {
			writeFit(sb, v, oldMin, oldMax, newMin, newMax, clamp, round)
		}),
	}
}

func writeFit(sb *strings.Builder, v, oldMin, oldMax, newMin, newMax AttrValueAppender, clamp, round bool) {
	sb.WriteString("@fit(")
	v.Append(sb)
	sb.WriteString(", ")
	oldMin.Append(sb)
	sb.WriteString(", ")
	oldMax.Append(sb)
	sb.WriteString(", ")
	newMin.Append(sb)
	sb.WriteString(", ")
	newMax.Append(sb)
	if clamp || round {
		sb.WriteString(", ")
		if clamp {
			sb.WriteString("true")
		} else {
			sb.WriteString("false")
		}
	}
	if round {
		sb.WriteString(", true")
	}
	sb.WriteString(")")
}
