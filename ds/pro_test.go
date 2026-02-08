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
	attr := OnRAF(Raw("$frameCount++"))
	if attr.Name != "data-on-raf" {
		t.Errorf("OnRAF().Name = %q, want %q", attr.Name, "data-on-raf")
	}
}

func TestOnResize(t *testing.T) {
	attr := OnResize(Raw("$width = el.offsetWidth"))
	if attr.Name != "data-on-resize" {
		t.Errorf("OnResize().Name = %q, want %q", attr.Name, "data-on-resize")
	}
}

func TestPersist(t *testing.T) {
	// Without options
	attr := Persist(nil)
	if attr.Name != "data-persist" {
		t.Errorf("Persist(nil).Name = %q, want %q", attr.Name, "data-persist")
	}
	if attr.Value != "" {
		t.Errorf("Persist(nil).Value = %q, want empty", attr.Value)
	}

	// With filter options
	include := "^user"
	attr = Persist(&FilterOptions{IncludeReg: &include})
	if !strings.Contains(attr.Value, "include: /^user/") {
		t.Errorf("Persist().Value = %q, should contain filter", attr.Value)
	}

	// With modifiers only
	attr = Persist(nil, Session())
	if !strings.Contains(attr.Name, "__session") {
		t.Errorf("Persist(nil, Session()).Name = %q, should contain __session", attr.Name)
	}
}

func TestPersistKey(t *testing.T) {
	attr := PersistKey("mykey")
	if !strings.HasPrefix(attr.Name, "data-persist:") {
		t.Errorf("PersistKey().Name = %q, should start with data-persist:", attr.Name)
	}
	if !strings.Contains(attr.Name, "mykey") {
		t.Errorf("PersistKey().Name = %q, should contain mykey", attr.Name)
	}

	// With session modifier
	attr = PersistKey("mykey", Session())
	if !strings.Contains(attr.Name, "__session") {
		t.Errorf("PersistKey(mykey, Session()).Name = %q, should contain __session", attr.Name)
	}
}

func TestQueryString(t *testing.T) {
	// Without options
	attr := QueryString(nil)
	if attr.Name != "data-query-string" {
		t.Errorf("QueryString(nil).Name = %q, want %q", attr.Name, "data-query-string")
	}
	if attr.Value != "" {
		t.Errorf("QueryString(nil).Value = %q, want empty", attr.Value)
	}

	// With filter options
	include := "^search"
	attr = QueryString(&FilterOptions{IncludeReg: &include})
	if !strings.Contains(attr.Value, "include: /^search/") {
		t.Errorf("QueryString().Value = %q, should contain filter", attr.Value)
	}

	// With modifiers only
	attr = QueryString(nil, Filter(), History())
	if !strings.Contains(attr.Name, "__filter") {
		t.Errorf("QueryString(nil, Filter()).Name = %q, should contain __filter", attr.Name)
	}
	if !strings.Contains(attr.Name, "__history") {
		t.Errorf("QueryString(nil, History()).Name = %q, should contain __history", attr.Name)
	}
}

func TestReplaceURL(t *testing.T) {
	attr := ReplaceURL(Raw("`/page${$page}`"))
	if attr.Name != "data-replace-url" {
		t.Errorf("ReplaceURL().Name = %q, want %q", attr.Name, "data-replace-url")
	}
}

func TestScrollIntoView(t *testing.T) {
	attr := ScrollIntoView(Smooth(), VCenter())
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
	attrName, _ := buildTestAttr("data-persist", Session())
	if !strings.Contains(attrName, "__session") {
		t.Errorf("Session() should add __session, got %q", attrName)
	}
}

func TestFilter(t *testing.T) {
	attrName, _ := buildTestAttr("data-query-string", Filter())
	if !strings.Contains(attrName, "__filter") {
		t.Errorf("Filter() should add __filter, got %q", attrName)
	}
}

func TestHistory(t *testing.T) {
	attrName, _ := buildTestAttr("data-query-string", History())
	if !strings.Contains(attrName, "__history") {
		t.Errorf("History() should add __history, got %q", attrName)
	}
}

func TestSmooth(t *testing.T) {
	attrName, _ := buildTestAttr("data-scroll-into-view", Smooth())
	if !strings.Contains(attrName, "__smooth") {
		t.Errorf("Smooth() should add __smooth, got %q", attrName)
	}
}

func TestInstant(t *testing.T) {
	attrName, _ := buildTestAttr("data-scroll-into-view", Instant())
	if !strings.Contains(attrName, "__instant") {
		t.Errorf("Instant() should add __instant, got %q", attrName)
	}
}

func TestAuto(t *testing.T) {
	attrName, _ := buildTestAttr("data-scroll-into-view", Auto())
	if !strings.Contains(attrName, "__auto") {
		t.Errorf("Auto() should add __auto, got %q", attrName)
	}
}

func TestHStart(t *testing.T) {
	attrName, _ := buildTestAttr("data-scroll-into-view", HStart())
	if !strings.Contains(attrName, "__hstart") {
		t.Errorf("HStart() should add __hstart, got %q", attrName)
	}
}

func TestHCenter(t *testing.T) {
	attrName, _ := buildTestAttr("data-scroll-into-view", HCenter())
	if !strings.Contains(attrName, "__hcenter") {
		t.Errorf("HCenter() should add __hcenter, got %q", attrName)
	}
}

func TestHEnd(t *testing.T) {
	attrName, _ := buildTestAttr("data-scroll-into-view", HEnd())
	if !strings.Contains(attrName, "__hend") {
		t.Errorf("HEnd() should add __hend, got %q", attrName)
	}
}

func TestHNearest(t *testing.T) {
	attrName, _ := buildTestAttr("data-scroll-into-view", HNearest())
	if !strings.Contains(attrName, "__hnearest") {
		t.Errorf("HNearest() should add __hnearest, got %q", attrName)
	}
}

func TestVStart(t *testing.T) {
	attrName, _ := buildTestAttr("data-scroll-into-view", VStart())
	if !strings.Contains(attrName, "__vstart") {
		t.Errorf("VStart() should add __vstart, got %q", attrName)
	}
}

func TestVCenter(t *testing.T) {
	attrName, _ := buildTestAttr("data-scroll-into-view", VCenter())
	if !strings.Contains(attrName, "__vcenter") {
		t.Errorf("VCenter() should add __vcenter, got %q", attrName)
	}
}

func TestVEnd(t *testing.T) {
	attrName, _ := buildTestAttr("data-scroll-into-view", VEnd())
	if !strings.Contains(attrName, "__vend") {
		t.Errorf("VEnd() should add __vend, got %q", attrName)
	}
}

func TestVNearest(t *testing.T) {
	attrName, _ := buildTestAttr("data-scroll-into-view", VNearest())
	if !strings.Contains(attrName, "__vnearest") {
		t.Errorf("VNearest() should add __vnearest, got %q", attrName)
	}
}

func TestFocus(t *testing.T) {
	attrName, _ := buildTestAttr("data-scroll-into-view", Focus())
	if !strings.Contains(attrName, "__focus") {
		t.Errorf("Focus() should add __focus, got %q", attrName)
	}
}

func TestClipboard(t *testing.T) {
	v := Clipboard(JsonValue("Hello, world!"))

	got := ToJS(v.expr)
	expected := `@clipboard("Hello, world!")`
	if got != expected {
		t.Errorf("Clipboard() = %q, want %q", got, expected)
	}

	// Test AttrMutator
	_, attrValue := buildTestAttr("test", v)
	if attrValue != expected {
		t.Errorf("Clipboard().Modify() = %q, want %q", attrValue, expected)
	}
}

func TestClipboardBase64(t *testing.T) {
	v := ClipboardBase64(JsonValue("SGVsbG8="))

	got := ToJS(v.expr)
	expected := `@clipboard("SGVsbG8=", true)`
	if got != expected {
		t.Errorf("ClipboardBase64() = %q, want %q", got, expected)
	}

	// Test AttrMutator
	_, attrValue := buildTestAttr("test", v)
	if attrValue != expected {
		t.Errorf("ClipboardBase64().Modify() = %q, want %q", attrValue, expected)
	}
}

func TestFit(t *testing.T) {
	v := Fit(Raw("$slider"), Raw("0"), Raw("100"), Raw("0"), Raw("255"))

	got := ToJS(v.expr)
	expected := "@fit($slider, 0, 100, 0, 255)"
	if got != expected {
		t.Errorf("Fit() = %q, want %q", got, expected)
	}

	// Test AttrMutator
	_, attrValue := buildTestAttr("test", v)
	if attrValue != expected {
		t.Errorf("Fit().Modify() = %q, want %q", attrValue, expected)
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
