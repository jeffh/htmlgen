package ds

import (
	"strings"
	"testing"
)

func TestAnimate(t *testing.T) {
	attr := Animate(Raw("opacity: $visible ? 1 : 0"))
	if attr.Name != "data-animate" {
		t.Errorf("Animate().Name = %q, want %q", attr.Name, "data-animate")
	}
}

func TestCustomValidity(t *testing.T) {
	attr := CustomValidity(Raw("$password === $confirm ? '' : 'Passwords must match'"))
	if attr.Name != "data-custom-validity" {
		t.Errorf("CustomValidity().Name = %q, want %q", attr.Name, "data-custom-validity")
	}
}

func TestOnRAF(t *testing.T) {
	attr := OnRAF(Raw("$frameCount++")).Attribute()
	if attr.Name != "data-on-raf" {
		t.Errorf("OnRAF().Name = %q, want %q", attr.Name, "data-on-raf")
	}
}

func TestOnResize(t *testing.T) {
	attr := OnResize(Raw("$width = el.offsetWidth")).Attribute()
	if attr.Name != "data-on-resize" {
		t.Errorf("OnResize().Name = %q, want %q", attr.Name, "data-on-resize")
	}
}

func TestPersist(t *testing.T) {
	// Without options
	attr := Persist(nil).Attribute()
	if attr.Name != "data-persist" {
		t.Errorf("Persist(nil).Name = %q, want %q", attr.Name, "data-persist")
	}
	if attr.Value != "" {
		t.Errorf("Persist(nil).Value = %q, want empty", attr.Value)
	}

	// With filter options
	include := "^user"
	attr = Persist(&FilterOptions{IncludeReg: &include}).Attribute()
	if !strings.Contains(attr.Value, "include: /^user/") {
		t.Errorf("Persist().Value = %q, should contain filter", attr.Value)
	}

	// With modifiers only
	attr = Persist(nil).Session().Attribute()
	if !strings.Contains(attr.Name, "__session") {
		t.Errorf("Persist(nil).Session().Name = %q, should contain __session", attr.Name)
	}
}

func TestPersistKey(t *testing.T) {
	attr := PersistKey("mykey").Attribute()
	if !strings.HasPrefix(attr.Name, "data-persist:") {
		t.Errorf("PersistKey().Name = %q, should start with data-persist:", attr.Name)
	}
	if !strings.Contains(attr.Name, "mykey") {
		t.Errorf("PersistKey().Name = %q, should contain mykey", attr.Name)
	}

	// With session modifier
	attr = PersistKey("mykey").Session().Attribute()
	if !strings.Contains(attr.Name, "__session") {
		t.Errorf("PersistKey(mykey).Session().Name = %q, should contain __session", attr.Name)
	}
}

func TestQueryString(t *testing.T) {
	// Without options
	attr := QueryString(nil).Attribute()
	if attr.Name != "data-query-string" {
		t.Errorf("QueryString(nil).Name = %q, want %q", attr.Name, "data-query-string")
	}
	if attr.Value != "" {
		t.Errorf("QueryString(nil).Value = %q, want empty", attr.Value)
	}

	// With filter options
	include := "^search"
	attr = QueryString(&FilterOptions{IncludeReg: &include}).Attribute()
	if !strings.Contains(attr.Value, "include: /^search/") {
		t.Errorf("QueryString().Value = %q, should contain filter", attr.Value)
	}

	// With modifiers only
	attr = QueryString(nil).Filter().History().Attribute()
	if !strings.Contains(attr.Name, "__filter") {
		t.Errorf("QueryString(nil).Filter().Name = %q, should contain __filter", attr.Name)
	}
	if !strings.Contains(attr.Name, "__history") {
		t.Errorf("QueryString(nil).History().Name = %q, should contain __history", attr.Name)
	}
}

func TestReplaceURL(t *testing.T) {
	attr := ReplaceURL(Raw("`/page${$page}`"))
	if attr.Name != "data-replace-url" {
		t.Errorf("ReplaceURL().Name = %q, want %q", attr.Name, "data-replace-url")
	}
}

func TestScrollIntoView(t *testing.T) {
	attr := ScrollIntoView().Smooth().VCenter().Attribute()
	if !strings.HasPrefix(attr.Name, "data-scroll-into-view") {
		t.Errorf("ScrollIntoView().Name = %q, should start with data-scroll-into-view", attr.Name)
	}
	if !strings.Contains(attr.Name, "__smooth") {
		t.Errorf("ScrollIntoView().Name = %q, should contain __smooth", attr.Name)
	}
	if !strings.Contains(attr.Name, "__vcenter") {
		t.Errorf("ScrollIntoView().Name = %q, should contain __vcenter", attr.Name)
	}
}

func TestViewTransitionName(t *testing.T) {
	attr := ViewTransitionName(Raw("$itemId"))
	if attr.Name != "data-view-transition" {
		t.Errorf("ViewTransitionName().Name = %q, want %q", attr.Name, "data-view-transition")
	}
}

func TestSession(t *testing.T) {
	attr := Persist(nil).Session().Attribute()
	if !strings.Contains(attr.Name, "__session") {
		t.Errorf("Session() should add __session, got %q", attr.Name)
	}
}

func TestFilter(t *testing.T) {
	attr := QueryString(nil).Filter().Attribute()
	if !strings.Contains(attr.Name, "__filter") {
		t.Errorf("Filter() should add __filter, got %q", attr.Name)
	}
}

func TestHistory(t *testing.T) {
	attr := QueryString(nil).History().Attribute()
	if !strings.Contains(attr.Name, "__history") {
		t.Errorf("History() should add __history, got %q", attr.Name)
	}
}

func TestSmooth(t *testing.T) {
	attr := ScrollIntoView().Smooth().Attribute()
	if !strings.Contains(attr.Name, "__smooth") {
		t.Errorf("Smooth() should add __smooth, got %q", attr.Name)
	}
}

func TestInstant(t *testing.T) {
	attr := ScrollIntoView().Instant().Attribute()
	if !strings.Contains(attr.Name, "__instant") {
		t.Errorf("Instant() should add __instant, got %q", attr.Name)
	}
}

func TestAuto(t *testing.T) {
	attr := ScrollIntoView().Auto().Attribute()
	if !strings.Contains(attr.Name, "__auto") {
		t.Errorf("Auto() should add __auto, got %q", attr.Name)
	}
}

func TestHStart(t *testing.T) {
	attr := ScrollIntoView().HStart().Attribute()
	if !strings.Contains(attr.Name, "__hstart") {
		t.Errorf("HStart() should add __hstart, got %q", attr.Name)
	}
}

func TestHCenter(t *testing.T) {
	attr := ScrollIntoView().HCenter().Attribute()
	if !strings.Contains(attr.Name, "__hcenter") {
		t.Errorf("HCenter() should add __hcenter, got %q", attr.Name)
	}
}

func TestHEnd(t *testing.T) {
	attr := ScrollIntoView().HEnd().Attribute()
	if !strings.Contains(attr.Name, "__hend") {
		t.Errorf("HEnd() should add __hend, got %q", attr.Name)
	}
}

func TestHNearest(t *testing.T) {
	attr := ScrollIntoView().HNearest().Attribute()
	if !strings.Contains(attr.Name, "__hnearest") {
		t.Errorf("HNearest() should add __hnearest, got %q", attr.Name)
	}
}

func TestVStart(t *testing.T) {
	attr := ScrollIntoView().VStart().Attribute()
	if !strings.Contains(attr.Name, "__vstart") {
		t.Errorf("VStart() should add __vstart, got %q", attr.Name)
	}
}

func TestVCenter(t *testing.T) {
	attr := ScrollIntoView().VCenter().Attribute()
	if !strings.Contains(attr.Name, "__vcenter") {
		t.Errorf("VCenter() should add __vcenter, got %q", attr.Name)
	}
}

func TestVEnd(t *testing.T) {
	attr := ScrollIntoView().VEnd().Attribute()
	if !strings.Contains(attr.Name, "__vend") {
		t.Errorf("VEnd() should add __vend, got %q", attr.Name)
	}
}

func TestVNearest(t *testing.T) {
	attr := ScrollIntoView().VNearest().Attribute()
	if !strings.Contains(attr.Name, "__vnearest") {
		t.Errorf("VNearest() should add __vnearest, got %q", attr.Name)
	}
}

func TestFocus(t *testing.T) {
	attr := ScrollIntoView().Focus().Attribute()
	if !strings.Contains(attr.Name, "__focus") {
		t.Errorf("Focus() should add __focus, got %q", attr.Name)
	}
}

func TestClipboard(t *testing.T) {
	v := Clipboard(JsonValue("Hello, world!"))

	got := ToJS(v.expr)
	expected := `@clipboard("Hello, world!")`
	if got != expected {
		t.Errorf("Clipboard() = %q, want %q", got, expected)
	}

	// Test as event-handler value
	attr := OnClick(v).Attribute()
	if attr.Value != expected {
		t.Errorf("OnClick(Clipboard()).Value = %q, want %q", attr.Value, expected)
	}
}

func TestClipboardBase64(t *testing.T) {
	v := ClipboardBase64(JsonValue("SGVsbG8="))

	got := ToJS(v.expr)
	expected := `@clipboard("SGVsbG8=", true)`
	if got != expected {
		t.Errorf("ClipboardBase64() = %q, want %q", got, expected)
	}

	// Test as event-handler value
	attr := OnClick(v).Attribute()
	if attr.Value != expected {
		t.Errorf("OnClick(ClipboardBase64()).Value = %q, want %q", attr.Value, expected)
	}
}

func TestFit(t *testing.T) {
	v := Fit(Raw("$slider"), Raw("0"), Raw("100"), Raw("0"), Raw("255"))

	got := ToJS(v.expr)
	expected := "@fit($slider, 0, 100, 0, 255)"
	if got != expected {
		t.Errorf("Fit() = %q, want %q", got, expected)
	}

	// Test as event-handler value
	attr := OnClick(v).Attribute()
	if attr.Value != expected {
		t.Errorf("OnClick(Fit()).Value = %q, want %q", attr.Value, expected)
	}
}

func TestFitClamped(t *testing.T) {
	v := FitClamped(Raw("$v"), Raw("0"), Raw("100"), Raw("0"), Raw("255"))

	got := ToJS(v.expr)
	expected := "@fit($v, 0, 100, 0, 255, true)"
	if got != expected {
		t.Errorf("FitClamped() = %q, want %q", got, expected)
	}
}

func TestFitRounded(t *testing.T) {
	v := FitRounded(Raw("$v"), Raw("0"), Raw("100"), Raw("0"), Raw("255"))

	got := ToJS(v.expr)
	expected := "@fit($v, 0, 100, 0, 255, false, true)"
	if got != expected {
		t.Errorf("FitRounded() = %q, want %q", got, expected)
	}
}

func TestFitClampedRounded(t *testing.T) {
	v := FitClampedRounded(Raw("$v"), Raw("0"), Raw("100"), Raw("0"), Raw("255"))

	got := ToJS(v.expr)
	expected := "@fit($v, 0, 100, 0, 255, true, true)"
	if got != expected {
		t.Errorf("FitClampedRounded() = %q, want %q", got, expected)
	}
}
