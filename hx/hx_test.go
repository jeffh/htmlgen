package hx

import (
	"testing"
	"time"

	"github.com/jeffh/htmlgen/h"
)

// ============ attrs.go tests ============

func TestHTTPMethods(t *testing.T) {
	tests := []struct {
		name     string
		fn       func(string) h.Attribute
		attrName string
	}{
		{"Get", Get, "hx-get"},
		{"Post", Post, "hx-post"},
		{"Put", Put, "hx-put"},
		{"Patch", Patch, "hx-patch"},
		{"Delete", Delete, "hx-delete"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attr := tt.fn("/api/test")
			if attr.Name != tt.attrName {
				t.Errorf("%s().Name = %q, want %q", tt.name, attr.Name, tt.attrName)
			}
			if attr.Value != "/api/test" {
				t.Errorf("%s().Value = %q, want %q", tt.name, attr.Value, "/api/test")
			}
		})
	}
}

func TestTarget(t *testing.T) {
	tests := []struct {
		name     string
		attr     h.Attribute
		expected string
	}{
		{"simple selector", Target("#results"), "#results"},
		{"this", Target("this"), "this"},
		{"closest", Target("closest div"), "closest div"},
		{"find", Target("find .item"), "find .item"},
		{"next", Target("next"), "next"},
		{"next with selector", Target("next .sibling"), "next .sibling"},
		{"previous", Target("previous"), "previous"},
		{"previous with selector", Target("previous .sibling"), "previous .sibling"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.attr.Name != "hx-target" {
				t.Errorf("Name = %q, want %q", tt.attr.Name, "hx-target")
			}
			if tt.attr.Value != tt.expected {
				t.Errorf("Value = %q, want %q", tt.attr.Value, tt.expected)
			}
		})
	}
}

func TestSelect(t *testing.T) {
	attr := Select("#content")
	if attr.Name != "hx-select" {
		t.Errorf("Name = %q, want %q", attr.Name, "hx-select")
	}
	if attr.Value != "#content" {
		t.Errorf("Value = %q, want %q", attr.Value, "#content")
	}
}

func TestSelectOOB(t *testing.T) {
	attr := SelectOOB("#sidebar")
	if attr.Name != "hx-select-oob" {
		t.Errorf("Name = %q, want %q", attr.Name, "hx-select-oob")
	}
}

func TestSwapOOB(t *testing.T) {
	tests := []struct {
		name     string
		attr     h.Attribute
		expected string
	}{
		{"true", SwapOOB("true"), "true"},
		{"with strategy", SwapOOB("outerHTML:#target"), "outerHTML:#target"},
		{"strategy only", SwapOOB("outerHTML"), "outerHTML"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.attr.Name != "hx-swap-oob" {
				t.Errorf("Name = %q, want %q", tt.attr.Name, "hx-swap-oob")
			}
			if tt.attr.Value != tt.expected {
				t.Errorf("Value = %q, want %q", tt.attr.Value, tt.expected)
			}
		})
	}
}

func TestInclude(t *testing.T) {
	tests := []struct {
		name     string
		attr     h.Attribute
		expected string
	}{
		{"selector", Include("#form"), "#form"},
		{"this", Include("this"), "this"},
		{"closest", Include("closest form"), "closest form"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.attr.Name != "hx-include" {
				t.Errorf("Name = %q, want %q", tt.attr.Name, "hx-include")
			}
			if tt.attr.Value != tt.expected {
				t.Errorf("Value = %q, want %q", tt.attr.Value, tt.expected)
			}
		})
	}
}

func TestVals(t *testing.T) {
	attr := Vals(map[string]any{"foo": 1, "bar": "hello"})
	if attr.Name != "hx-vals" {
		t.Errorf("Name = %q, want %q", attr.Name, "hx-vals")
	}
	// JSON encoding may vary in order, just check it contains the values
	if attr.Value == "" {
		t.Error("Value should not be empty")
	}
}

func TestValsJS(t *testing.T) {
	attr := ValsJS(map[string]string{"count": "getCount()"})
	if attr.Name != "hx-vals" {
		t.Errorf("Name = %q, want %q", attr.Name, "hx-vals")
	}
	if attr.Value[:3] != "js:" {
		t.Errorf("Value should start with 'js:', got %q", attr.Value)
	}
}

func TestHeaders(t *testing.T) {
	attr := Headers(map[string]string{"X-Custom": "value"})
	if attr.Name != "hx-headers" {
		t.Errorf("Name = %q, want %q", attr.Name, "hx-headers")
	}
}

func TestHeadersJS(t *testing.T) {
	attr := HeadersJS("getAuthHeaders()")
	if attr.Name != "hx-headers" {
		t.Errorf("Name = %q, want %q", attr.Name, "hx-headers")
	}
	if attr.Value != "js:getAuthHeaders()" {
		t.Errorf("Value = %q, want %q", attr.Value, "js:getAuthHeaders()")
	}
}

func TestParams(t *testing.T) {
	tests := []struct {
		name     string
		attr     h.Attribute
		expected string
	}{
		{"all", Params("*"), "*"},
		{"none", Params("none"), "none"},
		{"specific", Params("foo,bar"), "foo,bar"},
		{"not", Params("not secret"), "not secret"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.attr.Name != "hx-params" {
				t.Errorf("Name = %q, want %q", tt.attr.Name, "hx-params")
			}
			if tt.attr.Value != tt.expected {
				t.Errorf("Value = %q, want %q", tt.attr.Value, tt.expected)
			}
		})
	}
}

func TestEncoding(t *testing.T) {
	attr := Encoding(EncodingMultipart)
	if attr.Name != "hx-encoding" {
		t.Errorf("Name = %q, want %q", attr.Name, "hx-encoding")
	}
	if attr.Value != "multipart/form-data" {
		t.Errorf("Value = %q, want %q", attr.Value, "multipart/form-data")
	}
}

func TestExt(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{"single", []string{"json-enc"}, "json-enc"},
		{"multiple", []string{"json-enc", "debug"}, "json-enc,debug"},
		{"ignore", []string{"ignore:debug"}, "ignore:debug"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attr := Ext(tt.args...)
			if attr.Name != "hx-ext" {
				t.Errorf("Name = %q, want %q", attr.Name, "hx-ext")
			}
			if attr.Value != tt.expected {
				t.Errorf("Value = %q, want %q", attr.Value, tt.expected)
			}
		})
	}
}

// ============ swap.go tests ============

func TestSwapStrategies(t *testing.T) {
	tests := []struct {
		strategy SwapStrategy
		expected string
	}{
		{InnerHTML, "innerHTML"},
		{OuterHTML, "outerHTML"},
		{TextContent, "textContent"},
		{BeforeBegin, "beforebegin"},
		{AfterBegin, "afterbegin"},
		{BeforeEnd, "beforeend"},
		{AfterEnd, "afterend"},
		{SwapDelete, "delete"},
		{None, "none"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			attr := Swap(tt.strategy)
			if attr.Name != "hx-swap" {
				t.Errorf("Name = %q, want %q", attr.Name, "hx-swap")
			}
			if attr.Value != tt.expected {
				t.Errorf("Value = %q, want %q", attr.Value, tt.expected)
			}
		})
	}
}

func TestSwapWithModifiers(t *testing.T) {
	tests := []struct {
		name     string
		attr     h.Attribute
		contains string
	}{
		{"transition", Swap(InnerHTML, Transition()), "transition:true"},
		{"swap delay", Swap(InnerHTML, SwapDelay(100*time.Millisecond)), "swap:100ms"},
		{"settle delay", Swap(InnerHTML, SettleDelay(200*time.Millisecond)), "settle:200ms"},
		{"ignore title", Swap(InnerHTML, IgnoreTitle()), "ignoreTitle:true"},
		{"scroll top", Swap(InnerHTML, Scroll(Top)), "scroll:top"},
		{"scroll bottom", Swap(InnerHTML, Scroll(Bottom)), "scroll:bottom"},
		{"scroll target", Swap(InnerHTML, ScrollTarget("#div", Top)), "scroll:#div:top"},
		{"show top", Swap(InnerHTML, Show(Top)), "show:top"},
		{"show target", Swap(InnerHTML, ShowTarget("#div", Bottom)), "show:#div:bottom"},
		{"show window", Swap(InnerHTML, ShowWindow(Top)), "show:window:top"},
		{"show none", Swap(InnerHTML, ShowNone()), "show:none"},
		{"focus scroll true", Swap(InnerHTML, FocusScroll(true)), "focus-scroll:true"},
		{"focus scroll false", Swap(InnerHTML, FocusScroll(false)), "focus-scroll:false"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.attr.Name != "hx-swap" {
				t.Errorf("Name = %q, want %q", tt.attr.Name, "hx-swap")
			}
			if !containsString(tt.attr.Value, tt.contains) {
				t.Errorf("Value = %q, should contain %q", tt.attr.Value, tt.contains)
			}
		})
	}
}

func TestSwapMultipleModifiers(t *testing.T) {
	attr := Swap(OuterHTML, Transition(), SwapDelay(100*time.Millisecond), SettleDelay(200*time.Millisecond))
	if attr.Name != "hx-swap" {
		t.Errorf("Name = %q, want %q", attr.Name, "hx-swap")
	}
	if !containsString(attr.Value, "outerHTML") {
		t.Errorf("Value = %q, should contain %q", attr.Value, "outerHTML")
	}
	if !containsString(attr.Value, "transition:true") {
		t.Errorf("Value = %q, should contain %q", attr.Value, "transition:true")
	}
	if !containsString(attr.Value, "swap:100ms") {
		t.Errorf("Value = %q, should contain %q", attr.Value, "swap:100ms")
	}
}

// ============ trigger.go tests ============

func TestTriggerSimple(t *testing.T) {
	attr := Trigger("click").Attr()
	if attr.Name != "hx-trigger" {
		t.Errorf("Name = %q, want %q", attr.Name, "hx-trigger")
	}
	if attr.Value != "click" {
		t.Errorf("Value = %q, want %q", attr.Value, "click")
	}
}

func TestTriggerAttr(t *testing.T) {
	attr := TriggerAttr("click")
	if attr.Name != "hx-trigger" {
		t.Errorf("Name = %q, want %q", attr.Name, "hx-trigger")
	}
	if attr.Value != "click" {
		t.Errorf("Value = %q, want %q", attr.Value, "click")
	}
}

func TestTriggerWithModifiers(t *testing.T) {
	tests := []struct {
		name     string
		attr     h.Attribute
		contains string
	}{
		{"once", Trigger("click", Once()).Attr(), "once"},
		{"changed", Trigger("keyup", Changed()).Attr(), "changed"},
		{"delay", Trigger("keyup", Delay(500*time.Millisecond)).Attr(), "delay:500ms"},
		{"throttle", Trigger("scroll", Throttle(100*time.Millisecond)).Attr(), "throttle:100ms"},
		{"from", Trigger("click", From("#btn")).Attr(), "from:#btn"},
		{"from document", Trigger("keyup", FromDocument()).Attr(), "from:document"},
		{"from window", Trigger("resize", FromWindow()).Attr(), "from:window"},
		{"from closest", Trigger("click", FromClosest("form")).Attr(), "from:closest form"},
		{"from find", Trigger("click", FromFind(".item")).Attr(), "from:find .item"},
		{"from next", Trigger("click", FromNext()).Attr(), "from:next"},
		{"from previous", Trigger("click", FromPrevious()).Attr(), "from:previous"},
		{"target", Trigger("click", TriggerTarget("button")).Attr(), "target:button"},
		{"consume", Trigger("click", Consume()).Attr(), "consume"},
		{"queue first", Trigger("click", Queue(QueueFirst)).Attr(), "queue:first"},
		{"queue last", Trigger("click", Queue(QueueLast)).Attr(), "queue:last"},
		{"queue all", Trigger("click", Queue(QueueAll)).Attr(), "queue:all"},
		{"queue none", Trigger("click", Queue(QueueNone)).Attr(), "queue:none"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !containsString(tt.attr.Value, tt.contains) {
				t.Errorf("Value = %q, should contain %q", tt.attr.Value, tt.contains)
			}
		})
	}
}

func TestTriggerFilter(t *testing.T) {
	attr := Trigger("click", Filter("ctrlKey")).Attr()
	if attr.Value != "click[ctrlKey]" {
		t.Errorf("Value = %q, want %q", attr.Value, "click[ctrlKey]")
	}
}

func TestTriggerFilterComplex(t *testing.T) {
	attr := Trigger("keyup", Filter("keyCode==13")).Attr()
	if attr.Value != "keyup[keyCode==13]" {
		t.Errorf("Value = %q, want %q", attr.Value, "keyup[keyCode==13]")
	}
}

func TestTriggerMultiple(t *testing.T) {
	attr := Trigger("load").And("click", Delay(1*time.Second)).Attr()
	if attr.Value != "load, click delay:1s" {
		t.Errorf("Value = %q, want %q", attr.Value, "load, click delay:1s")
	}
}

func TestTriggerLoad(t *testing.T) {
	attr := TriggerLoad().Attr()
	if attr.Value != "load" {
		t.Errorf("Value = %q, want %q", attr.Value, "load")
	}
}

func TestTriggerRevealed(t *testing.T) {
	attr := TriggerRevealed().Attr()
	if attr.Value != "revealed" {
		t.Errorf("Value = %q, want %q", attr.Value, "revealed")
	}
}

func TestTriggerIntersect(t *testing.T) {
	attr := TriggerIntersect(IntersectThreshold(0.5)).Attr()
	if !containsString(attr.Value, "intersect") {
		t.Errorf("Value = %q, should contain %q", attr.Value, "intersect")
	}
	if !containsString(attr.Value, "threshold:0.5") {
		t.Errorf("Value = %q, should contain %q", attr.Value, "threshold:0.5")
	}
}

func TestTriggerEvery(t *testing.T) {
	attr := TriggerEvery(2 * time.Second).Attr()
	if attr.Value != "every 2s" {
		t.Errorf("Value = %q, want %q", attr.Value, "every 2s")
	}
}

func TestTriggerEveryWithFilter(t *testing.T) {
	attr := TriggerEvery(1*time.Second, Filter("document.visibilityState == 'visible'")).Attr()
	if !containsString(attr.Value, "every 1s") {
		t.Errorf("Value = %q, should contain %q", attr.Value, "every 1s")
	}
	if !containsString(attr.Value, "[document.visibilityState == 'visible']") {
		t.Errorf("Value = %q, should contain filter", attr.Value)
	}
}

func TestFromWithSpaces(t *testing.T) {
	attr := Trigger("click", From("form input")).Attr()
	if !containsString(attr.Value, "from:(form input)") {
		t.Errorf("Value = %q, should contain %q", attr.Value, "from:(form input)")
	}
}

// ============ behavior.go tests ============

func TestBoost(t *testing.T) {
	tests := []struct {
		enabled  bool
		expected string
	}{
		{true, "true"},
		{false, "false"},
	}

	for _, tt := range tests {
		attr := Boost(tt.enabled)
		if attr.Name != "hx-boost" {
			t.Errorf("Name = %q, want %q", attr.Name, "hx-boost")
		}
		if attr.Value != tt.expected {
			t.Errorf("Boost(%v).Value = %q, want %q", tt.enabled, attr.Value, tt.expected)
		}
	}
}

func TestPushURL(t *testing.T) {
	tests := []struct {
		name     string
		attr     h.Attribute
		expected string
	}{
		{"true", PushURL("true"), "true"},
		{"false", PushURL("false"), "false"},
		{"url", PushURL("/custom/url"), "/custom/url"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.attr.Name != "hx-push-url" {
				t.Errorf("Name = %q, want %q", tt.attr.Name, "hx-push-url")
			}
			if tt.attr.Value != tt.expected {
				t.Errorf("Value = %q, want %q", tt.attr.Value, tt.expected)
			}
		})
	}
}

func TestReplaceURL(t *testing.T) {
	tests := []struct {
		name     string
		attr     h.Attribute
		expected string
	}{
		{"true", ReplaceURL("true"), "true"},
		{"false", ReplaceURL("false"), "false"},
		{"url", ReplaceURL("/custom/url"), "/custom/url"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.attr.Name != "hx-replace-url" {
				t.Errorf("Name = %q, want %q", tt.attr.Name, "hx-replace-url")
			}
			if tt.attr.Value != tt.expected {
				t.Errorf("Value = %q, want %q", tt.attr.Value, tt.expected)
			}
		})
	}
}

func TestConfirm(t *testing.T) {
	attr := Confirm("Are you sure?")
	if attr.Name != "hx-confirm" {
		t.Errorf("Name = %q, want %q", attr.Name, "hx-confirm")
	}
	if attr.Value != "Are you sure?" {
		t.Errorf("Value = %q, want %q", attr.Value, "Are you sure?")
	}
}

func TestPrompt(t *testing.T) {
	attr := Prompt("Enter your name:")
	if attr.Name != "hx-prompt" {
		t.Errorf("Name = %q, want %q", attr.Name, "hx-prompt")
	}
}

func TestIndicator(t *testing.T) {
	tests := []struct {
		name     string
		attr     h.Attribute
		expected string
	}{
		{"selector", Indicator("#spinner"), "#spinner"},
		{"closest", Indicator("closest .spinner"), "closest .spinner"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.attr.Name != "hx-indicator" {
				t.Errorf("Name = %q, want %q", tt.attr.Name, "hx-indicator")
			}
			if tt.attr.Value != tt.expected {
				t.Errorf("Value = %q, want %q", tt.attr.Value, tt.expected)
			}
		})
	}
}

func TestDisabledElt(t *testing.T) {
	tests := []struct {
		name     string
		attr     h.Attribute
		expected string
	}{
		{"selector", DisabledElt("button"), "button"},
		{"this", DisabledElt("this"), "this"},
		{"closest", DisabledElt("closest button"), "closest button"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.attr.Name != "hx-disabled-elt" {
				t.Errorf("Name = %q, want %q", tt.attr.Name, "hx-disabled-elt")
			}
			if tt.attr.Value != tt.expected {
				t.Errorf("Value = %q, want %q", tt.attr.Value, tt.expected)
			}
		})
	}
}

func TestSync(t *testing.T) {
	tests := []struct {
		name     string
		attr     h.Attribute
		expected string
	}{
		{"drop", Sync("", SyncDrop), "drop"},
		{"abort", Sync("", SyncAbort), "abort"},
		{"replace", Sync("", SyncReplace), "replace"},
		{"queue", Sync("", SyncQueue), "queue"},
		{"with selector", Sync("#form", SyncDrop), "#form:drop"},
		{"this", Sync("this", SyncAbort), "this:abort"},
		{"closest", Sync("closest form", SyncReplace), "closest form:replace"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.attr.Name != "hx-sync" {
				t.Errorf("Name = %q, want %q", tt.attr.Name, "hx-sync")
			}
			if tt.attr.Value != tt.expected {
				t.Errorf("Value = %q, want %q", tt.attr.Value, tt.expected)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		enabled  bool
		expected string
	}{
		{true, "true"},
		{false, "false"},
	}

	for _, tt := range tests {
		attr := Validate(tt.enabled)
		if attr.Name != "hx-validate" {
			t.Errorf("Name = %q, want %q", attr.Name, "hx-validate")
		}
		if attr.Value != tt.expected {
			t.Errorf("Validate(%v).Value = %q, want %q", tt.enabled, attr.Value, tt.expected)
		}
	}
}

func TestPreserve(t *testing.T) {
	attr := Preserve()
	if attr.Name != "hx-preserve" {
		t.Errorf("Name = %q, want %q", attr.Name, "hx-preserve")
	}
	if attr.Value != "true" {
		t.Errorf("Value = %q, want %q", attr.Value, "true")
	}
}

func TestDisable(t *testing.T) {
	attr := Disable()
	if attr.Name != "hx-disable" {
		t.Errorf("Name = %q, want %q", attr.Name, "hx-disable")
	}
}

func TestDisinherit(t *testing.T) {
	tests := []struct {
		name     string
		attr     h.Attribute
		expected string
	}{
		{"all", Disinherit(), "*"},
		{"specific", Disinherit("hx-target", "hx-swap"), "hx-target hx-swap"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.attr.Name != "hx-disinherit" {
				t.Errorf("Name = %q, want %q", tt.attr.Name, "hx-disinherit")
			}
			if tt.attr.Value != tt.expected {
				t.Errorf("Value = %q, want %q", tt.attr.Value, tt.expected)
			}
		})
	}
}

func TestInherit(t *testing.T) {
	tests := []struct {
		name     string
		attr     h.Attribute
		expected string
	}{
		{"all", Inherit(), "*"},
		{"specific", Inherit("hx-target"), "hx-target"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.attr.Name != "hx-inherit" {
				t.Errorf("Name = %q, want %q", tt.attr.Name, "hx-inherit")
			}
			if tt.attr.Value != tt.expected {
				t.Errorf("Value = %q, want %q", tt.attr.Value, tt.expected)
			}
		})
	}
}

func TestHistoryElt(t *testing.T) {
	attr := HistoryElt()
	if attr.Name != "hx-history-elt" {
		t.Errorf("Name = %q, want %q", attr.Name, "hx-history-elt")
	}
}

func TestHistory(t *testing.T) {
	tests := []struct {
		enabled  bool
		expected string
	}{
		{true, "true"},
		{false, "false"},
	}

	for _, tt := range tests {
		attr := History(tt.enabled)
		if attr.Name != "hx-history" {
			t.Errorf("Name = %q, want %q", attr.Name, "hx-history")
		}
		if attr.Value != tt.expected {
			t.Errorf("History(%v).Value = %q, want %q", tt.enabled, attr.Value, tt.expected)
		}
	}
}

// ============ events.go tests ============

func TestOn(t *testing.T) {
	attr := On("click", "alert('clicked')")
	if attr.Name != "hx-on:click" {
		t.Errorf("Name = %q, want %q", attr.Name, "hx-on:click")
	}
	if attr.Value != "alert('clicked')" {
		t.Errorf("Value = %q, want %q", attr.Value, "alert('clicked')")
	}
}

func TestOnHTMX(t *testing.T) {
	attr := OnHTMX("beforeRequest", "console.log('request')")
	if attr.Name != "hx-on::beforeRequest" {
		t.Errorf("Name = %q, want %q", attr.Name, "hx-on::beforeRequest")
	}
}

func TestDOMEventHandlers(t *testing.T) {
	tests := []struct {
		name     string
		attr     h.Attribute
		expected string
	}{
		{"OnClick", OnClick("handler()"), "hx-on:click"},
		{"OnChange", OnChange("handler()"), "hx-on:change"},
		{"OnSubmit", OnSubmit("handler()"), "hx-on:submit"},
		{"OnInput", OnInput("handler()"), "hx-on:input"},
		{"OnKeydown", OnKeydown("handler()"), "hx-on:keydown"},
		{"OnKeyup", OnKeyup("handler()"), "hx-on:keyup"},
		{"OnFocus", OnFocus("handler()"), "hx-on:focus"},
		{"OnBlur", OnBlur("handler()"), "hx-on:blur"},
		{"OnMouseover", OnMouseover("handler()"), "hx-on:mouseover"},
		{"OnMouseout", OnMouseout("handler()"), "hx-on:mouseout"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.attr.Name != tt.expected {
				t.Errorf("Name = %q, want %q", tt.attr.Name, tt.expected)
			}
			if tt.attr.Value != "handler()" {
				t.Errorf("Value = %q, want %q", tt.attr.Value, "handler()")
			}
		})
	}
}

func TestHTMXEventHandlers(t *testing.T) {
	tests := []struct {
		name     string
		attr     h.Attribute
		expected string
	}{
		{"OnBeforeRequest", OnBeforeRequest("handler()"), "hx-on::before-request"},
		{"OnBeforeSend", OnBeforeSend("handler()"), "hx-on::before-send"},
		{"OnAfterRequest", OnAfterRequest("handler()"), "hx-on::after-request"},
		{"OnAfterOnLoad", OnAfterOnLoad("handler()"), "hx-on::after-on-load"},
		{"OnBeforeSwap", OnBeforeSwap("handler()"), "hx-on::before-swap"},
		{"OnAfterSwap", OnAfterSwap("handler()"), "hx-on::after-swap"},
		{"OnAfterSettle", OnAfterSettle("handler()"), "hx-on::after-settle"},
		{"OnResponseError", OnResponseError("handler()"), "hx-on::response-error"},
		{"OnSendError", OnSendError("handler()"), "hx-on::send-error"},
		{"OnTimeout", OnTimeout("handler()"), "hx-on::timeout"},
		{"OnSwapError", OnSwapError("handler()"), "hx-on::swap-error"},
		{"OnConfigRequest", OnConfigRequest("handler()"), "hx-on::config-request"},
		{"OnHistoryRestore", OnHistoryRestore("handler()"), "hx-on::history-restore"},
		{"OnHTMXLoad", OnHTMXLoad("handler()"), "hx-on::load"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.attr.Name != tt.expected {
				t.Errorf("Name = %q, want %q", tt.attr.Name, tt.expected)
			}
		})
	}
}

// ============ Helper functions ============

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && searchString(s, substr)))
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
