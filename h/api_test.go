package h

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
)

type staticAttr struct {
	attr Attribute
}

func (a staticAttr) Attribute() Attribute {
	return a.attr
}

func TestStreamingElementsAndEscaping(t *testing.T) {
	got := RenderString(func(b *B) {
		b.Div(
			Attrs("class", "first", "title", `"<&`),
			Attr("class", "last"),
			staticAttr{Attr("data-value", "a&b")},
			func(b *B) {
				b.Text(`"<tag>" & '` + "\x00")
				b.Br(Attr("hidden", ""))
				b.Raw("<strong>safe</strong>")
			},
		)
	})
	want := `<div class="last" title="&#34;&lt;&amp;" data-value="a&amp;b">&#34;&lt;tag&gt;&#34; &amp; &#39;�<br hidden/><strong>safe</strong></div>`
	if got != want {
		t.Fatalf("RenderString() = %q, want %q", got, want)
	}
}

func TestFormattedContent(t *testing.T) {
	got := RenderString(func(b *B) {
		b.Textf("%s %d", "<value>", 3)
		b.Rawf("<i>%s</i>", "raw")
	})
	want := "&lt;value&gt; 3<i>raw</i>"
	if got != want {
		t.Fatalf("formatted output = %q, want %q", got, want)
	}
}

func TestBodyLastWinsAndNilIgnored(t *testing.T) {
	var body Body = func(b *B) { b.Text("last") }
	got := RenderString(func(b *B) {
		b.Div(
			nil,
			func(b *B) { b.Text("first") },
			body,
		)
	})
	if got != "<div>last</div>" {
		t.Fatalf("output = %q", got)
	}
}

func TestUnknownElementArgumentPanics(t *testing.T) {
	defer func() {
		value := recover()
		if value == nil {
			t.Fatal("expected panic")
		}
		if !strings.Contains(value.(string), "unsupported argument type int for <div>") {
			t.Fatalf("panic = %q", value)
		}
	}()
	Render(io.Discard, func(b *B) {
		b.Div(42)
	})
}

func TestTypedNilBodyIsIgnored(t *testing.T) {
	var nilBody Body
	got := RenderString(func(b *B) {
		b.Div(func(b *B) { b.Text("kept") }, nilBody)
	})
	if got != "<div>kept</div>" {
		t.Fatalf("output = %q, want %q", got, "<div>kept</div>")
	}
}

func TestCallerAttributesNotMutated(t *testing.T) {
	shared := Attrs("class", "a")
	got := RenderString(func(b *B) {
		b.Div(shared, Attr("class", "b"))
	})
	if got != `<div class="b"></div>` {
		t.Fatalf("output = %q", got)
	}
	if value, _ := shared.Get("class"); value != "a" {
		t.Fatalf("caller Attributes mutated: class = %q, want %q", value, "a")
	}

	// Html must not write its lang default into a caller slice's spare capacity.
	shared = make(Attributes, 0, 4)
	shared = append(shared, Attribute{Name: "class", Value: "page"})
	RenderString(func(b *B) { b.Html(shared) })
	if len(shared) != 1 {
		t.Fatalf("caller Attributes grew: %v", shared)
	}
}

func TestZeroAttributeArgSkipped(t *testing.T) {
	got := RenderString(func(b *B) {
		b.Div(Attrs("id", "x"), AttrIf(false, "disabled", ""))
	})
	if got != `<div id="x"></div>` {
		t.Fatalf("output = %q", got)
	}
}

func TestVoidElementsIgnoreBodies(t *testing.T) {
	called := false
	body := func(b *B) {
		called = true
		b.Text("ignored")
	}
	got := RenderString(func(b *B) {
		b.Area(body)
		b.Base(body)
		b.Br(body)
		b.Col(body)
		b.Embed(body)
		b.Hr(body)
		b.Img(body)
		b.Input(body)
		b.Link(body)
		b.Meta(body)
		b.Source(body)
		b.Track(body)
		b.Wbr(body)
	})
	want := "<area/><base/><br/><col/><embed/><hr/><img/><input/><link/><meta/><source/><track/><wbr/>"
	if got != want {
		t.Fatalf("void output = %q, want %q", got, want)
	}
	if called {
		t.Fatal("void element ran its body")
	}
}

func TestHtmlDoctypeAndLanguage(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		got := RenderString(func(b *B) {
			b.Html(func(b *B) {
				b.Body(func(b *B) { b.Text("hello") })
			})
		})
		want := "<!DOCTYPE html>\n<html lang=\"en\"><body>hello</body></html>"
		if got != want {
			t.Fatalf("output = %q, want %q", got, want)
		}
	})

	t.Run("override", func(t *testing.T) {
		got := RenderString(func(b *B) {
			b.Html(Attrs("lang", "fr", "class", "page"))
		})
		want := "<!DOCTYPE html>\n<html lang=\"fr\" class=\"page\"></html>"
		if got != want {
			t.Fatalf("output = %q, want %q", got, want)
		}
	})
}

func TestDoctypeAndCustomElements(t *testing.T) {
	got := RenderString(func(b *B) {
		b.Doctype()
		b.El("user-card", Attr("name", "Ada"), func(b *B) {
			b.VoidEl("avatar", Attr("src", "/ada.png"))
		})
	})
	want := "<!DOCTYPE html>\n<user-card name=\"Ada\"><avatar src=\"/ada.png\"/></user-card>"
	if got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}

func TestRenderIndent(t *testing.T) {
	var output strings.Builder
	err := RenderIndent(&output, "  ", func(b *B) {
		b.Html(func(b *B) {
			b.Body(func(b *B) {
				b.H1(func(b *B) { b.Text("Title") })
				b.Img(Attr("src", "photo.jpg"))
			})
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	want := "<!DOCTYPE html>\n" +
		"<html lang=\"en\">\n" +
		"  <body>\n" +
		"    <h1>\n" +
		"      Title\n" +
		"    </h1>\n" +
		"    <img src=\"photo.jpg\"/>\n" +
		"  </body>\n" +
		"</html>\n"
	if output.String() != want {
		t.Fatalf("output =\n%q\nwant =\n%q", output.String(), want)
	}
}

func TestAttributeLineWrapping(t *testing.T) {
	var output strings.Builder
	err := RenderIndent(&output, "  ", func(b *B) {
		b.SetMaxLineLength(20)
		b.Div(
			Attrs("class", "one", "data-long", "two"),
			func(b *B) {
				b.Span(func(b *B) { b.Text("x") })
			},
		)
	})
	if err != nil {
		t.Fatal(err)
	}
	want := "<div class=\"one\"\n" +
		"  data-long=\"two\">\n" +
		"  <span>\n" +
		"    x\n" +
		"  </span>\n" +
		"</div>\n"
	if output.String() != want {
		t.Fatalf("output =\n%q\nwant =\n%q", output.String(), want)
	}
}

func TestMaxLineLengthInertWithoutIndent(t *testing.T) {
	got := RenderString(func(b *B) {
		b.SetMaxLineLength(10)
		b.Div(Attrs("class", "one", "data-long", "two"), func(b *B) {
			b.Text("x")
		})
	})
	if strings.Contains(got, "\n") {
		t.Fatalf("compact output contains newline: %q", got)
	}
}

var errWrite = errors.New("write failure")

type failAfterWriter struct {
	output bytes.Buffer
	remain int
}

func (w *failAfterWriter) Write(value []byte) (int, error) {
	if w.remain == 0 {
		return 0, errWrite
	}
	if len(value) > w.remain {
		count, _ := w.output.Write(value[:w.remain])
		w.remain = 0
		return count, errWrite
	}
	count, _ := w.output.Write(value)
	w.remain -= count
	return count, nil
}

func TestStickyErrorStopsOutput(t *testing.T) {
	writer := &failAfterWriter{remain: 8}
	sawError := false
	err := Render(writer, func(b *B) {
		b.Div(func(b *B) {
			b.Text("abcdefghijk")
			sawError = errors.Is(b.Err(), errWrite)
			b.Raw("not-written")
			b.Span(func(b *B) { b.Text("also-not-written") })
		})
		b.Raw("still-not-written")
	})
	if !errors.Is(err, errWrite) {
		t.Fatalf("Render error = %v, want %v", err, errWrite)
	}
	if !sawError {
		t.Fatal("B.Err did not expose the sticky error")
	}
	if strings.Contains(writer.output.String(), "not-written") {
		t.Fatalf("output continued after error: %q", writer.output.String())
	}
}

func TestRenderEntryPointsAndDefensiveClose(t *testing.T) {
	if err := Render(io.Discard, nil); err != nil {
		t.Fatal(err)
	}
	if got := RenderString(nil); got != "" {
		t.Fatalf("RenderString(nil) = %q", got)
	}
	if got := RenderBytes(nil); got != nil {
		t.Fatalf("RenderBytes(nil) = %#v, want nil", got)
	}

	got := RenderBytes(func(b *B) {
		b.openTag("div", nil)
		b.Text("open")
	})
	if string(got) != "<div>open</div>" {
		t.Fatalf("defensive close output = %q", got)
	}
}

func TestBuilderPoolReuse(t *testing.T) {
	first := getBuilder(io.Discard, "")
	first.openTags = append(first.openTags, "unused")
	putBuilder(first)
	second := getBuilder(io.Discard, "")
	defer putBuilder(second)
	if first != second {
		t.Skip("sync.Pool may discard entries")
	}
	if len(second.openTags) != 0 || second.Err() != nil {
		t.Fatalf("pooled B was not reset: %#v", second)
	}
}

func TestAttributesAPI(t *testing.T) {
	attrs := Attrs("class", "one", "id", "main")
	if got, ok := attrs.Get("class"); !ok || got != "one" {
		t.Fatalf("Get(class) = %q, %v", got, ok)
	}
	if attrs.Index("id") != 1 || attrs.Index("missing") != -1 {
		t.Fatal("Index returned an unexpected position")
	}
	attrs.Set("class", "two")
	attrs.Set("title", "hello")
	attrs.SetDefault("class", "ignored")
	attrs.SetDefault("role", "main")
	attrs.Delete("id")
	attrs.Merge(Attrs("title", "world", "data-x", "x"))
	want := Attributes{
		{Name: "class", Value: "two"},
		{Name: "title", Value: "world"},
		{Name: "role", Value: "main"},
		{Name: "data-x", Value: "x"},
	}
	if len(attrs) != len(want) {
		t.Fatalf("attributes = %#v, want %#v", attrs, want)
	}
	for index := range want {
		if attrs[index] != want[index] {
			t.Fatalf("attributes = %#v, want %#v", attrs, want)
		}
	}
}

func TestAttributeConstructors(t *testing.T) {
	if got := AttrIf(false, "hidden", ""); got != (Attribute{}) {
		t.Fatalf("AttrIf(false) = %#v", got)
	}
	if got := AttrIf(true, "hidden", ""); got != (Attribute{Name: "hidden"}) {
		t.Fatalf("AttrIf(true) = %#v", got)
	}
	if got := AttrsMap(map[string]string{"z": "1", "a": "2"}); got[0].Name != "a" || got[1].Name != "z" {
		t.Fatalf("AttrsMap order = %#v", got)
	}

	for _, test := range []struct {
		name string
		fn   func()
	}{
		{"empty Attr", func() { Attr("", "value") }},
		{"odd Attrs", func() { Attrs("name") }},
		{"empty Attrs name", func() { Attrs("", "value") }},
	} {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("expected panic")
				}
			}()
			test.fn()
		})
	}
}
