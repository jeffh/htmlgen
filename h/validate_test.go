package h

import (
	"io"
	"testing"
)

func mustPanic(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	fn()
}

func TestValidName(t *testing.T) {
	valid := []string{
		"a", "A", "div", "my-element", "data-on:click__debounce.500ms",
		"xlink:href", "data-bind:query__event.input.change", "h1", "SVG",
	}
	for _, name := range valid {
		if !validName(name) {
			t.Errorf("validName(%q) = false, want true", name)
		}
	}
	invalid := []string{
		"", " ", "1a", "-a", ":a", ".a", "_a", "div>",
		"div><script>alert(1)</script", "a b", "a=b", "a\"b", "a'b",
		"a/b", "a<b", "a\x00b", "héllo", "a\nb",
	}
	for _, name := range invalid {
		if validName(name) {
			t.Errorf("validName(%q) = true, want false", name)
		}
	}
}

func TestInvalidTagNamesPanic(t *testing.T) {
	for _, name := range []string{"", "div><script>alert(1)</script", "a b", "1div"} {
		mustPanic(t, func() {
			_ = Render(io.Discard, func(b *B) { b.El(name, nil, nil) })
		})
		mustPanic(t, func() {
			_ = Render(io.Discard, func(b *B) { b.VoidEl(name, nil) })
		})
	}
}

func TestValidCustomTagNames(t *testing.T) {
	got := RenderString(func(b *B) {
		b.El("my-widget", Attrs("class", "x"), func(b *B) { b.Text("hi") })
		b.VoidEl("spacer-el", nil)
	})
	want := `<my-widget class="x">hi</my-widget><spacer-el/>`
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestInvalidAttrNamesPanic(t *testing.T) {
	const bad = `onload=alert(1) x`
	mustPanic(t, func() { Attr(bad, "v") })
	mustPanic(t, func() { AttrIf(true, bad, "v") })
	mustPanic(t, func() { Attrs(bad, "v") })
	mustPanic(t, func() { AttrsMap(map[string]string{bad: "v"}) })
	mustPanic(t, func() {
		attrs := Attrs("class", "x")
		attrs.Set(bad, "v")
	})
	mustPanic(t, func() {
		attrs := Attrs("class", "x")
		attrs.SetDefault(bad, "v")
	})
}

func TestAttrIfFalseSkipsValidation(t *testing.T) {
	// A false condition returns a zero Attribute without validating the name.
	if got := AttrIf(false, "not a valid name", "v"); got != (Attribute{}) {
		t.Fatalf("AttrIf(false) = %#v", got)
	}
}
