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
			Attrs("class", "first", "title", `"<&`).With(
				Attr("class", "last"),
				staticAttr{Attr("data-value", "a&b")},
			),
			func(b *B) {
				b.Text(`"<tag>" & '` + "\x00")
				b.Br(AttrsOf(Attr("hidden", "")))
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

func TestNilAttrsAndNilBody(t *testing.T) {
	var nilBody Body
	got := RenderString(func(b *B) {
		b.Div(nil, nil)
		b.Span(nil, func(b *B) { b.Text("kept") })
		b.P(nil, nilBody)
	})
	want := "<div></div><span>kept</span><p></p>"
	if got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}

func TestCallerAttributesNotMutated(t *testing.T) {
	shared := Attrs("class", "a")
	got := RenderString(func(b *B) {
		b.Div(shared.With(Attr("class", "b")), nil)
	})
	if got != `<div class="b"></div>` {
		t.Fatalf("output = %q", got)
	}
	if value, _ := shared.Get("class"); value != "a" {
		t.Fatalf("caller Attributes mutated: class = %q, want %q", value, "a")
	}

	// With must not write into a caller slice's spare capacity either.
	shared = make(Attributes, 0, 4)
	shared = append(shared, Attribute{Name: "class", Value: "page"})
	if got := shared.With(Attr("id", "main")); len(got) != 2 {
		t.Fatalf("With result = %v", got)
	}
	if len(shared) != 1 {
		t.Fatalf("caller Attributes grew: %v", shared)
	}

	// Html must not write its lang default into a caller slice's spare capacity.
	RenderString(func(b *B) { b.Html(shared, nil) })
	if len(shared) != 1 {
		t.Fatalf("caller Attributes grew: %v", shared)
	}
}

func TestZeroAttributeSkipped(t *testing.T) {
	got := RenderString(func(b *B) {
		b.Div(Attrs("id", "x").With(AttrIf(false, "disabled", ""), nil), nil)
	})
	if got != `<div id="x"></div>` {
		t.Fatalf("output = %q", got)
	}
}

func TestVoidElements(t *testing.T) {
	got := RenderString(func(b *B) {
		b.Area(nil)
		b.Base(nil)
		b.Br(nil)
		b.Col(nil)
		b.Embed(nil)
		b.Hr(nil)
		b.Img(nil)
		b.Input(nil)
		b.Link(nil)
		b.Meta(nil)
		b.Source(nil)
		b.Track(nil)
		b.Wbr(nil)
	})
	want := "<area/><base/><br/><col/><embed/><hr/><img/><input/><link/><meta/><source/><track/><wbr/>"
	if got != want {
		t.Fatalf("void output = %q, want %q", got, want)
	}
}

func TestHtmlDoctypeAndLanguage(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		got := RenderString(func(b *B) {
			b.Html(nil, func(b *B) {
				b.Body(nil, func(b *B) { b.Text("hello") })
			})
		})
		want := "<!DOCTYPE html>\n<html lang=\"en\"><body>hello</body></html>"
		if got != want {
			t.Fatalf("output = %q, want %q", got, want)
		}
	})

	t.Run("override", func(t *testing.T) {
		got := RenderString(func(b *B) {
			b.Html(Attrs("lang", "fr", "class", "page"), nil)
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
		b.El("user-card", AttrsOf(Attr("name", "Ada")), func(b *B) {
			b.VoidEl("avatar", AttrsOf(Attr("src", "/ada.png")))
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
		b.Html(nil, func(b *B) {
			b.Body(nil, func(b *B) {
				b.H1(nil, func(b *B) { b.Text("Title") })
				b.Img(Attrs("src", "photo.jpg"))
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
				b.Span(nil, func(b *B) { b.Text("x") })
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
		b.Div(nil, func(b *B) {
			b.Text(strings.Repeat("a", flushThreshold))
			sawError = errors.Is(b.Err(), errWrite)
			b.Raw("not-written")
			b.Span(nil, func(b *B) { b.Text("also-not-written") })
		})
		b.Raw("still-not-written")
	})
	if !errors.Is(err, errWrite) {
		t.Fatalf("Render error = %v, want %v", err, errWrite)
	}
	if !sawError {
		t.Fatal("B.Err did not expose the sticky error after the failed flush")
	}
	if strings.Contains(writer.output.String(), "not-written") {
		t.Fatalf("output continued after error: %q", writer.output.String())
	}
}

func TestStickyErrorSurfacesFromFinalFlush(t *testing.T) {
	writer := &failAfterWriter{remain: 3}
	err := Render(writer, func(b *B) {
		b.Div(nil, func(b *B) { b.Text("hello") })
	})
	if !errors.Is(err, errWrite) {
		t.Fatalf("Render error = %v, want %v", err, errWrite)
	}
	if got := writer.output.String(); got != "<di" {
		t.Fatalf("output = %q, want %q", got, "<di")
	}
}

func TestFlushDeliversBufferedOutput(t *testing.T) {
	var output bytes.Buffer
	err := Render(&output, func(b *B) {
		b.Div(nil, func(b *B) {
			b.Text("first")
			if flushErr := b.Flush(); flushErr != nil {
				t.Error(flushErr)
			}
			if got := output.String(); got != "<div>first" {
				t.Errorf("after Flush output = %q, want %q", got, "<div>first")
			}
			b.Text("second")
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if got := output.String(); got != "<div>firstsecond</div>" {
		t.Fatalf("output = %q", got)
	}
}

func TestLargeRawBypassesBuffer(t *testing.T) {
	value := strings.Repeat("x", flushThreshold*2)
	got := RenderString(func(b *B) {
		b.Div(nil, func(b *B) {
			b.Text("a")
			b.Raw(value)
			b.Text("b")
		})
	})
	want := "<div>a" + value + "b</div>"
	if got != want {
		t.Fatalf("large Raw output length = %d, want %d", len(got), len(want))
	}
}

// byteOnlyWriter hides bytes.Buffer's WriteString so writes exercise the
// chunked non-StringWriter path.
type byteOnlyWriter struct {
	buf bytes.Buffer
}

func (w *byteOnlyWriter) Write(p []byte) (int, error) {
	return w.buf.Write(p)
}

func TestLargeRawToNonStringWriter(t *testing.T) {
	value := strings.Repeat("x", flushThreshold*2+7)
	var w byteOnlyWriter
	err := Render(&w, func(b *B) {
		b.Div(nil, func(b *B) {
			b.Text("a")
			b.Raw(value)
			b.Text("b")
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	want := "<div>a" + value + "b</div>"
	if w.buf.String() != want {
		t.Fatalf("output length = %d, want %d", w.buf.Len(), len(want))
	}
}

func TestOutputSpanningManyFlushes(t *testing.T) {
	const rows = 500
	var want strings.Builder
	got := RenderString(func(b *B) {
		for i := 0; i < rows; i++ {
			b.Li(Attrs("class", "row"), func(b *B) { b.Text("a&b") })
			want.WriteString(`<li class="row">a&amp;b</li>`)
		}
	})
	if got != want.String() {
		t.Fatalf("output length = %d, want %d", len(got), want.Len())
	}
}

func TestIndentedOutputSpanningManyFlushes(t *testing.T) {
	var output strings.Builder
	var want strings.Builder
	err := RenderIndent(&output, "  ", func(b *B) {
		for i := 0; i < 500; i++ {
			b.P(nil, func(b *B) { b.Text("x") })
			want.WriteString("<p>\n  x\n</p>\n")
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	if output.String() != want.String() {
		t.Fatalf("indented output length = %d, want %d", output.Len(), want.Len())
	}
}

func TestLargeEscapedValuesDoNotGrowBuffer(t *testing.T) {
	const size = 1 << 20
	text := strings.Repeat(`"`, size)
	attrValue := strings.Repeat("&", size)

	peak := 0
	observe := func(b *B) {
		if cap(b.buf) > peak {
			peak = cap(b.buf)
		}
	}

	var output bytes.Buffer
	err := Render(&output, func(b *B) {
		b.Div(Attrs("data-big", attrValue), func(b *B) {
			observe(b)
			b.Text("before")
			b.Text(text)
			observe(b)
			b.Text("after")
			observe(b)
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	var want strings.Builder
	want.WriteString(`<div data-big="`)
	want.WriteString(strings.Repeat("&amp;", size))
	want.WriteString(`">before`)
	want.WriteString(strings.Repeat("&#34;", size))
	want.WriteString("after</div>")
	if output.String() != want.String() {
		t.Fatalf("output length = %d, want %d", output.Len(), want.Len())
	}

	// Escaping the values through the buffer would have grown it to ~5 MiB;
	// a sub-threshold value may legitimately expand it to a few flush chunks.
	if limit := flushThreshold * 6; peak > limit {
		t.Fatalf("buffer grew to %d bytes, want <= %d", peak, limit)
	}
}

func TestNestedAndCustomElements(t *testing.T) {
	got := RenderString(func(b *B) {
		b.Div(Attrs("class", "card"), func(b *B) {
			b.H1(nil, func(b *B) { b.Text("<title>") })
			b.Img(Attrs("src", "a.png"))
			b.Br(nil)
			b.Span(nil, nil)
			b.El("user-card", Attrs("name", "Ada"), func(b *B) {
				b.VoidEl("avatar", nil)
			})
		})
	})
	want := `<div class="card"><h1>&lt;title&gt;</h1><img src="a.png"/><br/><span></span>` +
		`<user-card name="Ada"><avatar/></user-card></div>`
	if got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}

func TestElementBodyDoesNotAllocate(t *testing.T) {
	label := "captured"
	allocs := testing.AllocsPerRun(100, func() {
		Render(io.Discard, func(b *B) {
			b.Div(nil, func(b *B) { b.Text(label) })
		})
	})
	if allocs != 0 {
		t.Fatalf("Div with a capturing closure allocated %v times, want 0", allocs)
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
		b.openTag("<div", "</div>", nil)
		b.Text("open")
	})
	if string(got) != "<div>open</div>" {
		t.Fatalf("defensive close output = %q", got)
	}
}

func TestBuilderPoolReuse(t *testing.T) {
	first := getBuilder(io.Discard, "")
	first.openTags = append(first.openTags, "</unused>")
	first.buf = append(first.buf, "leftover"...)
	putBuilder(first)
	second := getBuilder(io.Discard, "")
	defer putBuilder(second)
	if first != second {
		t.Skip("sync.Pool may discard entries")
	}
	if len(second.openTags) != 0 || len(second.buf) != 0 || second.Err() != nil {
		t.Fatalf("pooled B was not reset: %#v", second)
	}
}

func TestOversizedBufferIsNotPooled(t *testing.T) {
	b := getBuilder(io.Discard, "")
	b.buf = make([]byte, 0, maxPooledBuffer+1)
	putBuilder(b)
	if cap(b.buf) > maxPooledBuffer {
		t.Fatalf("oversized buffer retained: cap = %d", cap(b.buf))
	}
	if cap(b.buf) != startingBuffer {
		t.Fatalf("buffer not re-primed: cap = %d, want %d", cap(b.buf), startingBuffer)
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

func TestAttrsOfAndWith(t *testing.T) {
	got := AttrsOf(
		Attr("class", "one"),
		staticAttr{Attr("data-x", "1")},
		Attr("class", "two"),
		AttrIf(false, "hidden", ""),
		nil,
	)
	want := Attributes{
		{Name: "class", Value: "two"},
		{Name: "data-x", Value: "1"},
	}
	if len(got) != len(want) {
		t.Fatalf("AttrsOf = %#v, want %#v", got, want)
	}
	for index := range want {
		if got[index] != want[index] {
			t.Fatalf("AttrsOf = %#v, want %#v", got, want)
		}
	}

	if got := AttrsOf(); got != nil {
		t.Fatalf("AttrsOf() = %#v, want nil", got)
	}
	if got := Attributes(nil).With(AttrIf(false, "hidden", "")); got != nil {
		t.Fatalf("With(zero attribute) = %#v, want nil", got)
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
