package html

import (
	"maps"
	"slices"
	"sort"
)

type Attributes []string

func (a *Attributes) Get(key string) (string, bool) {
	L := len(*a)
	for i := 0; i < L; i += 2 {
		k := (*a)[i]
		if k == key {
			return (*a)[i+1], true
		}
	}
	return "", false
}

func (a *Attributes) Index(key string) int {
	L := len(*a)
	for i := 0; i < L; i += 2 {
		k := (*a)[i]
		if k == key {
			return i
		}
	}
	return -1
}

func (a *Attributes) Set(key, value string) {
	idx := a.Index(key)
	if idx >= 0 {
		(*a)[idx+1] = value
	} else {
		*a = append(*a, key, value)
	}
}

func (a *Attributes) Delete(key string) {
	idx := a.Index(key)
	if idx >= 0 {
		*a = slices.Delete(*a, idx, idx+2)
	}
}

func Attrs(kv ...string) Attributes {
	if len(kv)%2 == 0 {
		panic("Attrs(...) expects an even number of arguments")
	}
	return Attributes(kv)
}

func AttrsMap(m map[string]string) Attributes {
	result := make([]string, 0, len(m)*2)
	keys := slices.Collect(maps.Keys(m))
	sort.Strings(keys)
	for _, k := range keys {
		result = append(result, k, m[k])
	}
	return Attributes(result)
}

// Document metadata
func Html(a Attributes, c ...Builder) Builder     { return tag("html", a, c) }
func Head(a Attributes, c ...Builder) Builder     { return tag("head", a, c) }
func Title(a Attributes, c ...Builder) Builder    { return tag("title", a, c) }
func Meta(a Attributes, c ...Builder) Builder     { return tag("meta", a, c) }
func Link(a Attributes, c ...Builder) Builder     { return tag("link", a, c) }
func Style(a Attributes, c ...Builder) Builder    { return tag("style", a, c) }
func Script(a Attributes, c ...Builder) Builder   { return tag("script", a, c) }
func Noscript(a Attributes, c ...Builder) Builder { return tag("noscript", a, c) }
func Base(a Attributes, c ...Builder) Builder     { return tag("base", a, c) }

// Sectioning root
func Body(a Attributes, c ...Builder) Builder { return tag("body", a, c) }

// Content sectioning
func Address(a Attributes, c ...Builder) Builder { return tag("address", a, c) }
func Article(a Attributes, c ...Builder) Builder { return tag("article", a, c) }
func Aside(a Attributes, c ...Builder) Builder   { return tag("aside", a, c) }
func Footer(a Attributes, c ...Builder) Builder  { return tag("footer", a, c) }
func Header(a Attributes, c ...Builder) Builder  { return tag("header", a, c) }
func H1(a Attributes, c ...Builder) Builder      { return tag("h1", a, c) }
func H2(a Attributes, c ...Builder) Builder      { return tag("h2", a, c) }
func H3(a Attributes, c ...Builder) Builder      { return tag("h3", a, c) }
func H4(a Attributes, c ...Builder) Builder      { return tag("h4", a, c) }
func H5(a Attributes, c ...Builder) Builder      { return tag("h5", a, c) }
func H6(a Attributes, c ...Builder) Builder      { return tag("h6", a, c) }
func Hgroup(a Attributes, c ...Builder) Builder  { return tag("hgroup", a, c) }
func Main(a Attributes, c ...Builder) Builder    { return tag("main", a, c) }
func Nav(a Attributes, c ...Builder) Builder     { return tag("nav", a, c) }
func Section(a Attributes, c ...Builder) Builder { return tag("section", a, c) }
func Search(a Attributes, c ...Builder) Builder  { return tag("search", a, c) }

// Text content
func Blockquote(a Attributes, c ...Builder) Builder { return tag("blockquote", a, c) }
func Dd(a Attributes, c ...Builder) Builder         { return tag("dd", a, c) }
func Div(a Attributes, c ...Builder) Builder        { return tag("div", a, c) }
func Dl(a Attributes, c ...Builder) Builder         { return tag("dl", a, c) }
func Dt(a Attributes, c ...Builder) Builder         { return tag("dt", a, c) }
func Figcaption(a Attributes, c ...Builder) Builder { return tag("figcaption", a, c) }
func Figure(a Attributes, c ...Builder) Builder     { return tag("figure", a, c) }
func Hr(a Attributes, c ...Builder) Builder         { return tag("hr", a, c) }
func Li(a Attributes, c ...Builder) Builder         { return tag("li", a, c) }
func Menu(a Attributes, c ...Builder) Builder       { return tag("menu", a, c) }
func Ol(a Attributes, c ...Builder) Builder         { return tag("ol", a, c) }
func P(a Attributes, c ...Builder) Builder          { return tag("p", a, c) }
func Pre(a Attributes, c ...Builder) Builder        { return tag("pre", a, c) }
func Ul(a Attributes, c ...Builder) Builder         { return tag("ul", a, c) }

// Inline text semantics
func A(a Attributes, c ...Builder) Builder      { return tag("a", a, c) }
func Abbr(a Attributes, c ...Builder) Builder   { return tag("abbr", a, c) }
func B(a Attributes, c ...Builder) Builder      { return tag("b", a, c) }
func Bdi(a Attributes, c ...Builder) Builder    { return tag("bdi", a, c) }
func Bdo(a Attributes, c ...Builder) Builder    { return tag("bdo", a, c) }
func Br(a Attributes, c ...Builder) Builder     { return tag("br", a, c) }
func Cite(a Attributes, c ...Builder) Builder   { return tag("cite", a, c) }
func Code(a Attributes, c ...Builder) Builder   { return tag("code", a, c) }
func Data(a Attributes, c ...Builder) Builder   { return tag("data", a, c) }
func Dfn(a Attributes, c ...Builder) Builder    { return tag("dfn", a, c) }
func Em(a Attributes, c ...Builder) Builder     { return tag("em", a, c) }
func I(a Attributes, c ...Builder) Builder      { return tag("i", a, c) }
func Kbd(a Attributes, c ...Builder) Builder    { return tag("kbd", a, c) }
func Mark(a Attributes, c ...Builder) Builder   { return tag("mark", a, c) }
func Q(a Attributes, c ...Builder) Builder      { return tag("q", a, c) }
func Rp(a Attributes, c ...Builder) Builder     { return tag("rp", a, c) }
func Rt(a Attributes, c ...Builder) Builder     { return tag("rt", a, c) }
func Ruby(a Attributes, c ...Builder) Builder   { return tag("ruby", a, c) }
func S(a Attributes, c ...Builder) Builder      { return tag("s", a, c) }
func Samp(a Attributes, c ...Builder) Builder   { return tag("samp", a, c) }
func Small(a Attributes, c ...Builder) Builder  { return tag("small", a, c) }
func Span(a Attributes, c ...Builder) Builder   { return tag("span", a, c) }
func Strong(a Attributes, c ...Builder) Builder { return tag("strong", a, c) }
func Sub(a Attributes, c ...Builder) Builder    { return tag("sub", a, c) }
func Sup(a Attributes, c ...Builder) Builder    { return tag("sup", a, c) }
func Time(a Attributes, c ...Builder) Builder   { return tag("time", a, c) }
func U(a Attributes, c ...Builder) Builder      { return tag("u", a, c) }
func Var(a Attributes, c ...Builder) Builder    { return tag("var", a, c) }
func Wbr(a Attributes, c ...Builder) Builder    { return tag("wbr", a, c) }

// Image and multimedia
func Area(a Attributes, c ...Builder) Builder  { return tag("area", a, c) }
func Audio(a Attributes, c ...Builder) Builder { return tag("audio", a, c) }
func Img(a Attributes, c ...Builder) Builder   { return tag("img", a, c) }
func Map(a Attributes, c ...Builder) Builder   { return tag("map", a, c) }
func Track(a Attributes, c ...Builder) Builder { return tag("track", a, c) }
func Video(a Attributes, c ...Builder) Builder { return tag("video", a, c) }

// Embedded content
func Embed(a Attributes, c ...Builder) Builder   { return tag("embed", a, c) }
func Iframe(a Attributes, c ...Builder) Builder  { return tag("iframe", a, c) }
func Object(a Attributes, c ...Builder) Builder  { return tag("object", a, c) }
func Picture(a Attributes, c ...Builder) Builder { return tag("picture", a, c) }
func Portal(a Attributes, c ...Builder) Builder  { return tag("portal", a, c) }
func Source(a Attributes, c ...Builder) Builder  { return tag("source", a, c) }

// SVG and MathML
func Svg(a Attributes, c ...Builder) Builder  { return tag("svg", a, c) }
func Math(a Attributes, c ...Builder) Builder { return tag("math", a, c) }

// Scripting
func Canvas(a Attributes, c ...Builder) Builder   { return tag("canvas", a, c) }
func Template(a Attributes, c ...Builder) Builder { return tag("template", a, c) }
func Slot(a Attributes, c ...Builder) Builder     { return tag("slot", a, c) }

// Demarcating edits
func Del(a Attributes, c ...Builder) Builder { return tag("del", a, c) }
func Ins(a Attributes, c ...Builder) Builder { return tag("ins", a, c) }

// Table content
func Caption(a Attributes, c ...Builder) Builder  { return tag("caption", a, c) }
func Col(a Attributes, c ...Builder) Builder      { return tag("col", a, c) }
func Colgroup(a Attributes, c ...Builder) Builder { return tag("colgroup", a, c) }
func Table(a Attributes, c ...Builder) Builder    { return tag("table", a, c) }
func Tbody(a Attributes, c ...Builder) Builder    { return tag("tbody", a, c) }
func Td(a Attributes, c ...Builder) Builder       { return tag("td", a, c) }
func Tfoot(a Attributes, c ...Builder) Builder    { return tag("tfoot", a, c) }
func Th(a Attributes, c ...Builder) Builder       { return tag("th", a, c) }
func Thead(a Attributes, c ...Builder) Builder    { return tag("thead", a, c) }
func Tr(a Attributes, c ...Builder) Builder       { return tag("tr", a, c) }

// Forms
func Button(a Attributes, c ...Builder) Builder   { return tag("button", a, c) }
func Datalist(a Attributes, c ...Builder) Builder { return tag("datalist", a, c) }
func Fieldset(a Attributes, c ...Builder) Builder { return tag("fieldset", a, c) }
func Form(a Attributes, c ...Builder) Builder     { return tag("form", a, c) }
func Input(a Attributes, c ...Builder) Builder    { return tag("input", a, c) }
func Label(a Attributes, c ...Builder) Builder    { return tag("label", a, c) }
func Legend(a Attributes, c ...Builder) Builder   { return tag("legend", a, c) }
func Meter(a Attributes, c ...Builder) Builder    { return tag("meter", a, c) }
func Optgroup(a Attributes, c ...Builder) Builder { return tag("optgroup", a, c) }
func Option(a Attributes, c ...Builder) Builder   { return tag("option", a, c) }
func Output(a Attributes, c ...Builder) Builder   { return tag("output", a, c) }
func Progress(a Attributes, c ...Builder) Builder { return tag("progress", a, c) }
func Select(a Attributes, c ...Builder) Builder   { return tag("select", a, c) }
func Textarea(a Attributes, c ...Builder) Builder { return tag("textarea", a, c) }

// Interactive elements
func Details(a Attributes, c ...Builder) Builder { return tag("details", a, c) }
func Dialog(a Attributes, c ...Builder) Builder  { return tag("dialog", a, c) }
func Summary(a Attributes, c ...Builder) Builder { return tag("summary", a, c) }

// Web Components
func CustomElement(name string, a Attributes, c ...Builder) Builder { return tag(name, a, c) }
