package h

// Fragment creates a Builder that renders its children without a wrapping element.
func Fragment(children ...Builder) Builder { return &fragmentBuilder{children} }

// Text creates a Builder that renders HTML-escaped text content.
func Text(value string) Builder { return &textBuilder{value, false} }

// Raw creates a Builder that renders unescaped HTML content.
// Use with caution as this can introduce XSS vulnerabilities.
func Raw(value string) Builder { return &textBuilder{value, true} }

// Html creates the root <html> element with DOCTYPE declaration.
// Sets lang="en" by default if not specified in attributes.
func Html(args ...TagArg) Builder {
	attrs, children := parseTagArgs(args)
	return &htmlTagBuilder{attrs, children}
}

// Head creates a <head> element for document metadata.
func Head(args ...TagArg) Builder { return tag("head", args...) }

// Title creates a <title> element for the document title.
func Title(args ...TagArg) Builder { return tag("title", args...) }

// Meta creates a <meta> element for document metadata.
func Meta(args ...TagArg) Builder { return stag("meta", args...) }

// Link creates a <link> element for external resources.
func Link(args ...TagArg) Builder { return stag("link", args...) }

// Style creates a <style> element for embedded CSS.
func Style(args ...TagArg) Builder { return tag("style", args...) }

// Script creates a <script> element for JavaScript.
func Script(args ...TagArg) Builder { return tag("script", args...) }

// Noscript creates a <noscript> element for fallback content.
func Noscript(args ...TagArg) Builder { return tag("noscript", args...) }

// Base creates a <base> element for document base URL.
func Base(args ...TagArg) Builder { return stag("base", args...) }

// Body creates a <body> element as the document's sectioning root.
func Body(args ...TagArg) Builder { return tag("body", args...) }

// Address creates an <address> element for contact information.
func Address(args ...TagArg) Builder { return tag("address", args...) }

// Article creates an <article> element for self-contained content.
func Article(args ...TagArg) Builder { return tag("article", args...) }

// Aside creates an <aside> element for tangentially related content.
func Aside(args ...TagArg) Builder { return tag("aside", args...) }

// Footer creates a <footer> element for footer content.
func Footer(args ...TagArg) Builder { return tag("footer", args...) }

// Header creates a <header> element for introductory content.
func Header(args ...TagArg) Builder { return tag("header", args...) }

// H1 creates an <h1> heading element.
func H1(args ...TagArg) Builder { return tag("h1", args...) }

// H2 creates an <h2> heading element.
func H2(args ...TagArg) Builder { return tag("h2", args...) }

// H3 creates an <h3> heading element.
func H3(args ...TagArg) Builder { return tag("h3", args...) }

// H4 creates an <h4> heading element.
func H4(args ...TagArg) Builder { return tag("h4", args...) }

// H5 creates an <h5> heading element.
func H5(args ...TagArg) Builder { return tag("h5", args...) }

// H6 creates an <h6> heading element.
func H6(args ...TagArg) Builder { return tag("h6", args...) }

// Hgroup creates an <hgroup> element for heading groups.
func Hgroup(args ...TagArg) Builder { return tag("hgroup", args...) }

// Main creates a <main> element for the dominant content.
func Main(args ...TagArg) Builder { return tag("main", args...) }

// Nav creates a <nav> element for navigation links.
func Nav(args ...TagArg) Builder { return tag("nav", args...) }

// Section creates a <section> element for thematic content grouping.
func Section(args ...TagArg) Builder { return tag("section", args...) }

// Search creates a <search> element for search functionality.
func Search(args ...TagArg) Builder { return tag("search", args...) }

// Blockquote creates a <blockquote> element for extended quotations.
func Blockquote(args ...TagArg) Builder { return tag("blockquote", args...) }

// Dd creates a <dd> element for description list definitions.
func Dd(args ...TagArg) Builder { return tag("dd", args...) }

// Div creates a <div> element as a generic container.
func Div(args ...TagArg) Builder { return tag("div", args...) }

// Dl creates a <dl> element for description lists.
func Dl(args ...TagArg) Builder { return tag("dl", args...) }

// Dt creates a <dt> element for description list terms.
func Dt(args ...TagArg) Builder { return tag("dt", args...) }

// Figcaption creates a <figcaption> element for figure captions.
func Figcaption(args ...TagArg) Builder { return tag("figcaption", args...) }

// Figure creates a <figure> element for self-contained content with optional caption.
func Figure(args ...TagArg) Builder { return tag("figure", args...) }

// Hr creates an <hr> element for thematic breaks.
func Hr(args ...TagArg) Builder { return stag("hr", args...) }

// Li creates an <li> element for list items.
func Li(args ...TagArg) Builder { return tag("li", args...) }

// Menu creates a <menu> element for toolbar menus.
func Menu(args ...TagArg) Builder { return tag("menu", args...) }

// Ol creates an <ol> element for ordered lists.
func Ol(args ...TagArg) Builder { return tag("ol", args...) }

// P creates a <p> element for paragraphs.
func P(args ...TagArg) Builder { return tag("p", args...) }

// Pre creates a <pre> element for preformatted text.
func Pre(args ...TagArg) Builder { return tag("pre", args...) }

// Ul creates a <ul> element for unordered lists.
func Ul(args ...TagArg) Builder { return tag("ul", args...) }

// A creates an <a> element for hyperlinks.
func A(args ...TagArg) Builder { return tag("a", args...) }

// Abbr creates an <abbr> element for abbreviations.
func Abbr(args ...TagArg) Builder { return tag("abbr", args...) }

// B creates a <b> element for bold text.
func B(args ...TagArg) Builder { return tag("b", args...) }

// Bdi creates a <bdi> element for bidirectional text isolation.
func Bdi(args ...TagArg) Builder { return tag("bdi", args...) }

// Bdo creates a <bdo> element for bidirectional text override.
func Bdo(args ...TagArg) Builder { return tag("bdo", args...) }

// Br creates a <br> element for line breaks.
func Br(args ...TagArg) Builder { return stag("br", args...) }

// Cite creates a <cite> element for citations.
func Cite(args ...TagArg) Builder { return tag("cite", args...) }

// Code creates a <code> element for code fragments.
func Code(args ...TagArg) Builder { return tag("code", args...) }

// Data creates a <data> element for machine-readable content.
func Data(args ...TagArg) Builder { return tag("data", args...) }

// Dfn creates a <dfn> element for definitions.
func Dfn(args ...TagArg) Builder { return tag("dfn", args...) }

// Em creates an <em> element for emphasized text.
func Em(args ...TagArg) Builder { return tag("em", args...) }

// I creates an <i> element for idiomatic text.
func I(args ...TagArg) Builder { return tag("i", args...) }

// Kbd creates a <kbd> element for keyboard input.
func Kbd(args ...TagArg) Builder { return tag("kbd", args...) }

// Mark creates a <mark> element for highlighted text.
func Mark(args ...TagArg) Builder { return tag("mark", args...) }

// Q creates a <q> element for inline quotations.
func Q(args ...TagArg) Builder { return tag("q", args...) }

// Rp creates an <rp> element for ruby fallback parentheses.
func Rp(args ...TagArg) Builder { return tag("rp", args...) }

// Rt creates an <rt> element for ruby text.
func Rt(args ...TagArg) Builder { return tag("rt", args...) }

// Ruby creates a <ruby> element for ruby annotations.
func Ruby(args ...TagArg) Builder { return tag("ruby", args...) }

// S creates an <s> element for strikethrough text.
func S(args ...TagArg) Builder { return tag("s", args...) }

// Samp creates a <samp> element for sample output.
func Samp(args ...TagArg) Builder { return tag("samp", args...) }

// Small creates a <small> element for side comments.
func Small(args ...TagArg) Builder { return tag("small", args...) }

// Span creates a <span> element as a generic inline container.
func Span(args ...TagArg) Builder { return tag("span", args...) }

// Strong creates a <strong> element for strong importance.
func Strong(args ...TagArg) Builder { return tag("strong", args...) }

// Sub creates a <sub> element for subscript text.
func Sub(args ...TagArg) Builder { return tag("sub", args...) }

// Sup creates a <sup> element for superscript text.
func Sup(args ...TagArg) Builder { return tag("sup", args...) }

// Time creates a <time> element for dates and times.
func Time(args ...TagArg) Builder { return tag("time", args...) }

// U creates a <u> element for underlined text.
func U(args ...TagArg) Builder { return tag("u", args...) }

// Var creates a <var> element for variable names.
func Var(args ...TagArg) Builder { return tag("var", args...) }

// Wbr creates a <wbr> element for word break opportunities.
func Wbr(args ...TagArg) Builder { return stag("wbr", args...) }

// Area creates an <area> element for image map areas.
func Area(args ...TagArg) Builder { return stag("area", args...) }

// Audio creates an <audio> element for sound content.
func Audio(args ...TagArg) Builder { return tag("audio", args...) }

// Img creates an <img> element for images.
func Img(args ...TagArg) Builder { return stag("img", args...) }

// Map creates a <map> element for image maps.
func Map(args ...TagArg) Builder { return tag("map", args...) }

// Track creates a <track> element for media text tracks.
func Track(args ...TagArg) Builder { return stag("track", args...) }

// Video creates a <video> element for video content.
func Video(args ...TagArg) Builder { return tag("video", args...) }

// Embed creates an <embed> element for external content.
func Embed(args ...TagArg) Builder { return stag("embed", args...) }

// Iframe creates an <iframe> element for nested browsing contexts.
func Iframe(args ...TagArg) Builder { return tag("iframe", args...) }

// Object creates an <object> element for external resources.
func Object(args ...TagArg) Builder { return tag("object", args...) }

// Picture creates a <picture> element for responsive images.
func Picture(args ...TagArg) Builder { return tag("picture", args...) }

// Portal creates a <portal> element for embedded pages.
func Portal(args ...TagArg) Builder { return tag("portal", args...) }

// Source creates a <source> element for media sources.
func Source(args ...TagArg) Builder { return stag("source", args...) }

// Svg creates an <svg> element for SVG graphics.
func Svg(args ...TagArg) Builder { return tag("svg", args...) }

// Math creates a <math> element for MathML content.
func Math(args ...TagArg) Builder { return tag("math", args...) }

// Canvas creates a <canvas> element for graphics rendering.
func Canvas(args ...TagArg) Builder { return tag("canvas", args...) }

// Template creates a <template> element for client-side content templates.
func Template(args ...TagArg) Builder { return tag("template", args...) }

// Slot creates a <slot> element for web component content distribution.
func Slot(args ...TagArg) Builder { return tag("slot", args...) }

// Del creates a <del> element for deleted text.
func Del(args ...TagArg) Builder { return tag("del", args...) }

// Ins creates an <ins> element for inserted text.
func Ins(args ...TagArg) Builder { return tag("ins", args...) }

// Caption creates a <caption> element for table captions.
func Caption(args ...TagArg) Builder { return tag("caption", args...) }

// Col creates a <col> element for table column properties.
func Col(args ...TagArg) Builder { return stag("col", args...) }

// Colgroup creates a <colgroup> element for table column groups.
func Colgroup(args ...TagArg) Builder { return tag("colgroup", args...) }

// Table creates a <table> element for tabular data.
func Table(args ...TagArg) Builder { return tag("table", args...) }

// Tbody creates a <tbody> element for table body content.
func Tbody(args ...TagArg) Builder { return tag("tbody", args...) }

// Td creates a <td> element for table data cells.
func Td(args ...TagArg) Builder { return tag("td", args...) }

// Tfoot creates a <tfoot> element for table footer content.
func Tfoot(args ...TagArg) Builder { return tag("tfoot", args...) }

// Th creates a <th> element for table header cells.
func Th(args ...TagArg) Builder { return tag("th", args...) }

// Thead creates a <thead> element for table header content.
func Thead(args ...TagArg) Builder { return tag("thead", args...) }

// Tr creates a <tr> element for table rows.
func Tr(args ...TagArg) Builder { return tag("tr", args...) }

// Button creates a <button> element for clickable buttons.
func Button(args ...TagArg) Builder { return tag("button", args...) }

// Datalist creates a <datalist> element for input suggestions.
func Datalist(args ...TagArg) Builder { return tag("datalist", args...) }

// Fieldset creates a <fieldset> element for form field groups.
func Fieldset(args ...TagArg) Builder { return tag("fieldset", args...) }

// Form creates a <form> element for user input forms.
func Form(args ...TagArg) Builder { return tag("form", args...) }

// Input creates an <input> element for form inputs.
func Input(args ...TagArg) Builder { return stag("input", args...) }

// Label creates a <label> element for form control labels.
func Label(args ...TagArg) Builder { return tag("label", args...) }

// Legend creates a <legend> element for fieldset captions.
func Legend(args ...TagArg) Builder { return tag("legend", args...) }

// Meter creates a <meter> element for scalar measurements.
func Meter(args ...TagArg) Builder { return tag("meter", args...) }

// Optgroup creates an <optgroup> element for option groups.
func Optgroup(args ...TagArg) Builder { return tag("optgroup", args...) }

// Option creates an <option> element for select options.
func Option(args ...TagArg) Builder { return tag("option", args...) }

// Output creates an <output> element for calculation results.
func Output(args ...TagArg) Builder { return tag("output", args...) }

// Progress creates a <progress> element for progress indicators.
func Progress(args ...TagArg) Builder { return tag("progress", args...) }

// Select creates a <select> element for dropdown lists.
func Select(args ...TagArg) Builder { return tag("select", args...) }

// Textarea creates a <textarea> element for multi-line text input.
func Textarea(args ...TagArg) Builder { return tag("textarea", args...) }

// Details creates a <details> element for disclosure widgets.
func Details(args ...TagArg) Builder { return tag("details", args...) }

// Dialog creates a <dialog> element for modal dialogs.
func Dialog(args ...TagArg) Builder { return tag("dialog", args...) }

// Summary creates a <summary> element for details disclosure.
func Summary(args ...TagArg) Builder { return tag("summary", args...) }

// CustomElement creates a custom HTML element with the given tag name.
func CustomElement(name string, args ...TagArg) Builder { return tag(name, args...) }
