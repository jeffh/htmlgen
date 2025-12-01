package h

import "io"

// Render writes the HTML representation of the given Builder to w.
// Returns nil if b is nil.
func Render(w io.Writer, b Builder) error {
	if b == nil {
		return nil
	}
	return b.Build(NewWriter(w))
}

// RenderIndent writes the HTML representation of the given Builder to w
// with indentation for readability. The indent parameter specifies the string
// to use for each indentation level (e.g., "  " for two spaces or "\t" for tabs).
// Returns nil if b is nil.
func RenderIndent(w io.Writer, indent string, b Builder) error {
	if b == nil {
		return nil
	}
	writer := NewWriter(w)
	writer.SetIndent(indent)
	return b.Build(writer)
}

// Fragment creates a Builder that renders its children without a wrapping element.
func Fragment(c ...Builder) Builder { return &fragmentBuilder{c} }

// Text creates a Builder that renders HTML-escaped text content.
func Text(value string) Builder { return &textBuilder{value, false} }

// Raw creates a Builder that renders unescaped HTML content.
// Use with caution as this can introduce XSS vulnerabilities.
func Raw(value string) Builder { return &textBuilder{value, true} }

// Html creates the root <html> element with DOCTYPE declaration.
// Sets lang="en" by default if not specified in attributes.
func Html(a Attributes, c ...Builder) Builder { return &htmlTagBuilder{a, c} }

// Head creates a <head> element for document metadata.
func Head(a Attributes, c ...Builder) Builder { return tag("head", a, c) }

// Title creates a <title> element for the document title.
func Title(a Attributes, c ...Builder) Builder { return tag("title", a, c) }

// Meta creates a <meta> element for document metadata.
func Meta(a Attributes, c ...Builder) Builder { return stag("meta", a, c) }

// Link creates a <link> element for external resources.
func Link(a Attributes, c ...Builder) Builder { return stag("link", a, c) }

// Style creates a <style> element for embedded CSS.
func Style(a Attributes, c ...Builder) Builder { return tag("style", a, c) }

// Script creates a <script> element for JavaScript.
func Script(a Attributes, c ...Builder) Builder { return tag("script", a, c) }

// Noscript creates a <noscript> element for fallback content.
func Noscript(a Attributes, c ...Builder) Builder { return tag("noscript", a, c) }

// Base creates a <base> element for document base URL.
func Base(a Attributes, c ...Builder) Builder { return stag("base", a, c) }

// Body creates a <body> element as the document's sectioning root.
func Body(a Attributes, c ...Builder) Builder { return tag("body", a, c) }

// Address creates an <address> element for contact information.
func Address(a Attributes, c ...Builder) Builder { return tag("address", a, c) }

// Article creates an <article> element for self-contained content.
func Article(a Attributes, c ...Builder) Builder { return tag("article", a, c) }

// Aside creates an <aside> element for tangentially related content.
func Aside(a Attributes, c ...Builder) Builder { return tag("aside", a, c) }

// Footer creates a <footer> element for footer content.
func Footer(a Attributes, c ...Builder) Builder { return tag("footer", a, c) }

// Header creates a <header> element for introductory content.
func Header(a Attributes, c ...Builder) Builder { return tag("header", a, c) }

// H1 creates an <h1> heading element.
func H1(a Attributes, c ...Builder) Builder { return tag("h1", a, c) }

// H2 creates an <h2> heading element.
func H2(a Attributes, c ...Builder) Builder { return tag("h2", a, c) }

// H3 creates an <h3> heading element.
func H3(a Attributes, c ...Builder) Builder { return tag("h3", a, c) }

// H4 creates an <h4> heading element.
func H4(a Attributes, c ...Builder) Builder { return tag("h4", a, c) }

// H5 creates an <h5> heading element.
func H5(a Attributes, c ...Builder) Builder { return tag("h5", a, c) }

// H6 creates an <h6> heading element.
func H6(a Attributes, c ...Builder) Builder { return tag("h6", a, c) }

// Hgroup creates an <hgroup> element for heading groups.
func Hgroup(a Attributes, c ...Builder) Builder { return tag("hgroup", a, c) }

// Main creates a <main> element for the dominant content.
func Main(a Attributes, c ...Builder) Builder { return tag("main", a, c) }

// Nav creates a <nav> element for navigation links.
func Nav(a Attributes, c ...Builder) Builder { return tag("nav", a, c) }

// Section creates a <section> element for thematic content grouping.
func Section(a Attributes, c ...Builder) Builder { return tag("section", a, c) }

// Search creates a <search> element for search functionality.
func Search(a Attributes, c ...Builder) Builder { return tag("search", a, c) }

// Blockquote creates a <blockquote> element for extended quotations.
func Blockquote(a Attributes, c ...Builder) Builder { return tag("blockquote", a, c) }

// Dd creates a <dd> element for description list definitions.
func Dd(a Attributes, c ...Builder) Builder { return tag("dd", a, c) }

// Div creates a <div> element as a generic container.
func Div(a Attributes, c ...Builder) Builder { return tag("div", a, c) }

// Dl creates a <dl> element for description lists.
func Dl(a Attributes, c ...Builder) Builder { return tag("dl", a, c) }

// Dt creates a <dt> element for description list terms.
func Dt(a Attributes, c ...Builder) Builder { return tag("dt", a, c) }

// Figcaption creates a <figcaption> element for figure captions.
func Figcaption(a Attributes, c ...Builder) Builder { return tag("figcaption", a, c) }

// Figure creates a <figure> element for self-contained content with optional caption.
func Figure(a Attributes, c ...Builder) Builder { return tag("figure", a, c) }

// Hr creates an <hr> element for thematic breaks.
func Hr(a Attributes, c ...Builder) Builder { return stag("hr", a, c) }

// Li creates an <li> element for list items.
func Li(a Attributes, c ...Builder) Builder { return tag("li", a, c) }

// Menu creates a <menu> element for toolbar menus.
func Menu(a Attributes, c ...Builder) Builder { return tag("menu", a, c) }

// Ol creates an <ol> element for ordered lists.
func Ol(a Attributes, c ...Builder) Builder { return tag("ol", a, c) }

// P creates a <p> element for paragraphs.
func P(a Attributes, c ...Builder) Builder { return tag("p", a, c) }

// Pre creates a <pre> element for preformatted text.
func Pre(a Attributes, c ...Builder) Builder { return tag("pre", a, c) }

// Ul creates a <ul> element for unordered lists.
func Ul(a Attributes, c ...Builder) Builder { return tag("ul", a, c) }

// A creates an <a> element for hyperlinks.
func A(a Attributes, c ...Builder) Builder { return tag("a", a, c) }

// Abbr creates an <abbr> element for abbreviations.
func Abbr(a Attributes, c ...Builder) Builder { return tag("abbr", a, c) }

// B creates a <b> element for bold text.
func B(a Attributes, c ...Builder) Builder { return tag("b", a, c) }

// Bdi creates a <bdi> element for bidirectional text isolation.
func Bdi(a Attributes, c ...Builder) Builder { return tag("bdi", a, c) }

// Bdo creates a <bdo> element for bidirectional text override.
func Bdo(a Attributes, c ...Builder) Builder { return tag("bdo", a, c) }

// Br creates a <br> element for line breaks.
func Br(a Attributes, c ...Builder) Builder { return stag("br", a, c) }

// Cite creates a <cite> element for citations.
func Cite(a Attributes, c ...Builder) Builder { return tag("cite", a, c) }

// Code creates a <code> element for code fragments.
func Code(a Attributes, c ...Builder) Builder { return tag("code", a, c) }

// Data creates a <data> element for machine-readable content.
func Data(a Attributes, c ...Builder) Builder { return tag("data", a, c) }

// Dfn creates a <dfn> element for definitions.
func Dfn(a Attributes, c ...Builder) Builder { return tag("dfn", a, c) }

// Em creates an <em> element for emphasized text.
func Em(a Attributes, c ...Builder) Builder { return tag("em", a, c) }

// I creates an <i> element for idiomatic text.
func I(a Attributes, c ...Builder) Builder { return tag("i", a, c) }

// Kbd creates a <kbd> element for keyboard input.
func Kbd(a Attributes, c ...Builder) Builder { return tag("kbd", a, c) }

// Mark creates a <mark> element for highlighted text.
func Mark(a Attributes, c ...Builder) Builder { return tag("mark", a, c) }

// Q creates a <q> element for inline quotations.
func Q(a Attributes, c ...Builder) Builder { return tag("q", a, c) }

// Rp creates an <rp> element for ruby fallback parentheses.
func Rp(a Attributes, c ...Builder) Builder { return tag("rp", a, c) }

// Rt creates an <rt> element for ruby text.
func Rt(a Attributes, c ...Builder) Builder { return tag("rt", a, c) }

// Ruby creates a <ruby> element for ruby annotations.
func Ruby(a Attributes, c ...Builder) Builder { return tag("ruby", a, c) }

// S creates an <s> element for strikethrough text.
func S(a Attributes, c ...Builder) Builder { return tag("s", a, c) }

// Samp creates a <samp> element for sample output.
func Samp(a Attributes, c ...Builder) Builder { return tag("samp", a, c) }

// Small creates a <small> element for side comments.
func Small(a Attributes, c ...Builder) Builder { return tag("small", a, c) }

// Span creates a <span> element as a generic inline container.
func Span(a Attributes, c ...Builder) Builder { return tag("span", a, c) }

// Strong creates a <strong> element for strong importance.
func Strong(a Attributes, c ...Builder) Builder { return tag("strong", a, c) }

// Sub creates a <sub> element for subscript text.
func Sub(a Attributes, c ...Builder) Builder { return tag("sub", a, c) }

// Sup creates a <sup> element for superscript text.
func Sup(a Attributes, c ...Builder) Builder { return tag("sup", a, c) }

// Time creates a <time> element for dates and times.
func Time(a Attributes, c ...Builder) Builder { return tag("time", a, c) }

// U creates a <u> element for underlined text.
func U(a Attributes, c ...Builder) Builder { return tag("u", a, c) }

// Var creates a <var> element for variable names.
func Var(a Attributes, c ...Builder) Builder { return tag("var", a, c) }

// Wbr creates a <wbr> element for word break opportunities.
func Wbr(a Attributes, c ...Builder) Builder { return stag("wbr", a, c) }

// Area creates an <area> element for image map areas.
func Area(a Attributes, c ...Builder) Builder { return stag("area", a, c) }

// Audio creates an <audio> element for sound content.
func Audio(a Attributes, c ...Builder) Builder { return tag("audio", a, c) }

// Img creates an <img> element for images.
func Img(a Attributes, c ...Builder) Builder { return stag("img", a, c) }

// Map creates a <map> element for image maps.
func Map(a Attributes, c ...Builder) Builder { return tag("map", a, c) }

// Track creates a <track> element for media text tracks.
func Track(a Attributes, c ...Builder) Builder { return stag("track", a, c) }

// Video creates a <video> element for video content.
func Video(a Attributes, c ...Builder) Builder { return tag("video", a, c) }

// Embed creates an <embed> element for external content.
func Embed(a Attributes, c ...Builder) Builder { return stag("embed", a, c) }

// Iframe creates an <iframe> element for nested browsing contexts.
func Iframe(a Attributes, c ...Builder) Builder { return tag("iframe", a, c) }

// Object creates an <object> element for external resources.
func Object(a Attributes, c ...Builder) Builder { return tag("object", a, c) }

// Picture creates a <picture> element for responsive images.
func Picture(a Attributes, c ...Builder) Builder { return tag("picture", a, c) }

// Portal creates a <portal> element for embedded pages.
func Portal(a Attributes, c ...Builder) Builder { return tag("portal", a, c) }

// Source creates a <source> element for media sources.
func Source(a Attributes, c ...Builder) Builder { return stag("source", a, c) }

// Svg creates an <svg> element for SVG graphics.
func Svg(a Attributes, c ...Builder) Builder { return tag("svg", a, c) }

// Math creates a <math> element for MathML content.
func Math(a Attributes, c ...Builder) Builder { return tag("math", a, c) }

// Canvas creates a <canvas> element for graphics rendering.
func Canvas(a Attributes, c ...Builder) Builder { return tag("canvas", a, c) }

// Template creates a <template> element for client-side content templates.
func Template(a Attributes, c ...Builder) Builder { return tag("template", a, c) }

// Slot creates a <slot> element for web component content distribution.
func Slot(a Attributes, c ...Builder) Builder { return tag("slot", a, c) }

// Del creates a <del> element for deleted text.
func Del(a Attributes, c ...Builder) Builder { return tag("del", a, c) }

// Ins creates an <ins> element for inserted text.
func Ins(a Attributes, c ...Builder) Builder { return tag("ins", a, c) }

// Caption creates a <caption> element for table captions.
func Caption(a Attributes, c ...Builder) Builder { return tag("caption", a, c) }

// Col creates a <col> element for table column properties.
func Col(a Attributes, c ...Builder) Builder { return stag("col", a, c) }

// Colgroup creates a <colgroup> element for table column groups.
func Colgroup(a Attributes, c ...Builder) Builder { return tag("colgroup", a, c) }

// Table creates a <table> element for tabular data.
func Table(a Attributes, c ...Builder) Builder { return tag("table", a, c) }

// Tbody creates a <tbody> element for table body content.
func Tbody(a Attributes, c ...Builder) Builder { return tag("tbody", a, c) }

// Td creates a <td> element for table data cells.
func Td(a Attributes, c ...Builder) Builder { return tag("td", a, c) }

// Tfoot creates a <tfoot> element for table footer content.
func Tfoot(a Attributes, c ...Builder) Builder { return tag("tfoot", a, c) }

// Th creates a <th> element for table header cells.
func Th(a Attributes, c ...Builder) Builder { return tag("th", a, c) }

// Thead creates a <thead> element for table header content.
func Thead(a Attributes, c ...Builder) Builder { return tag("thead", a, c) }

// Tr creates a <tr> element for table rows.
func Tr(a Attributes, c ...Builder) Builder { return tag("tr", a, c) }

// Button creates a <button> element for clickable buttons.
func Button(a Attributes, c ...Builder) Builder { return tag("button", a, c) }

// Datalist creates a <datalist> element for input suggestions.
func Datalist(a Attributes, c ...Builder) Builder { return tag("datalist", a, c) }

// Fieldset creates a <fieldset> element for form field groups.
func Fieldset(a Attributes, c ...Builder) Builder { return tag("fieldset", a, c) }

// Form creates a <form> element for user input forms.
func Form(a Attributes, c ...Builder) Builder { return tag("form", a, c) }

// Input creates an <input> element for form inputs.
func Input(a Attributes, c ...Builder) Builder { return stag("input", a, c) }

// Label creates a <label> element for form control labels.
func Label(a Attributes, c ...Builder) Builder { return tag("label", a, c) }

// Legend creates a <legend> element for fieldset captions.
func Legend(a Attributes, c ...Builder) Builder { return tag("legend", a, c) }

// Meter creates a <meter> element for scalar measurements.
func Meter(a Attributes, c ...Builder) Builder { return tag("meter", a, c) }

// Optgroup creates an <optgroup> element for option groups.
func Optgroup(a Attributes, c ...Builder) Builder { return tag("optgroup", a, c) }

// Option creates an <option> element for select options.
func Option(a Attributes, c ...Builder) Builder { return tag("option", a, c) }

// Output creates an <output> element for calculation results.
func Output(a Attributes, c ...Builder) Builder { return tag("output", a, c) }

// Progress creates a <progress> element for progress indicators.
func Progress(a Attributes, c ...Builder) Builder { return tag("progress", a, c) }

// Select creates a <select> element for dropdown lists.
func Select(a Attributes, c ...Builder) Builder { return tag("select", a, c) }

// Textarea creates a <textarea> element for multi-line text input.
func Textarea(a Attributes, c ...Builder) Builder { return tag("textarea", a, c) }

// Details creates a <details> element for disclosure widgets.
func Details(a Attributes, c ...Builder) Builder { return tag("details", a, c) }

// Dialog creates a <dialog> element for modal dialogs.
func Dialog(a Attributes, c ...Builder) Builder { return tag("dialog", a, c) }

// Summary creates a <summary> element for details disclosure.
func Summary(a Attributes, c ...Builder) Builder { return tag("summary", a, c) }

// CustomElement creates a custom HTML element with the given tag name.
func CustomElement(name string, a Attributes, c ...Builder) Builder { return tag(name, a, c) }
