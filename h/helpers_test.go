package h

import "testing"

func TestNativeControlFlowAndComposition(t *testing.T) {
	items := []string{"one", "hidden", "three"}
	list := func(b *B, values []string) {
		b.Ul(func(b *B) {
			for _, item := range values {
				if item == "hidden" {
					continue
				}
				b.Li(func(b *B) { b.Text(item) })
			}
		})
	}

	got := RenderString(func(b *B) {
		list(b, items)
	})
	want := "<ul><li>one</li><li>three</li></ul>"
	if got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}
