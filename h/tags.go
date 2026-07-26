package h

// Html writes the root <html> element, preceded by the HTML5 doctype. It uses
// lang="en" unless a lang attribute is provided.
func (b *B) Html(args ...any) {
	if b.err != nil {
		return
	}
	attrs, body := parseArgs("html", args)
	if _, ok := attrs.Get("lang"); !ok {
		// Full slice expression so appending never writes into a caller-owned
		// slice's spare capacity.
		attrs = append(attrs[:len(attrs):len(attrs)], Attribute{Name: "lang", Value: "en"})
	}
	b.Doctype()
	b.openTag("html", attrs)
	if b.err != nil {
		return
	}
	if body != nil {
		body(b)
	}
	b.closeOneTag()
}

// Head writes a <head> element.
func (b *B) Head(args ...any) { b.element("head", args...) }

// Title writes a <title> element.
func (b *B) Title(args ...any) { b.element("title", args...) }

// Meta writes a self-closing <meta> element.
func (b *B) Meta(args ...any) { b.voidElement("meta", args...) }

// Link writes a self-closing <link> element.
func (b *B) Link(args ...any) { b.voidElement("link", args...) }

// Style writes a <style> element.
func (b *B) Style(args ...any) { b.element("style", args...) }

// Script writes a <script> element.
func (b *B) Script(args ...any) { b.element("script", args...) }

// Noscript writes a <noscript> element.
func (b *B) Noscript(args ...any) { b.element("noscript", args...) }

// Base writes a self-closing <base> element.
func (b *B) Base(args ...any) { b.voidElement("base", args...) }

// Body writes a <body> element.
func (b *B) Body(args ...any) { b.element("body", args...) }

// Address writes an <address> element.
func (b *B) Address(args ...any) { b.element("address", args...) }

// Article writes an <article> element.
func (b *B) Article(args ...any) { b.element("article", args...) }

// Aside writes an <aside> element.
func (b *B) Aside(args ...any) { b.element("aside", args...) }

// Footer writes a <footer> element.
func (b *B) Footer(args ...any) { b.element("footer", args...) }

// Header writes a <header> element.
func (b *B) Header(args ...any) { b.element("header", args...) }

// H1 writes an <h1> element.
func (b *B) H1(args ...any) { b.element("h1", args...) }

// H2 writes an <h2> element.
func (b *B) H2(args ...any) { b.element("h2", args...) }

// H3 writes an <h3> element.
func (b *B) H3(args ...any) { b.element("h3", args...) }

// H4 writes an <h4> element.
func (b *B) H4(args ...any) { b.element("h4", args...) }

// H5 writes an <h5> element.
func (b *B) H5(args ...any) { b.element("h5", args...) }

// H6 writes an <h6> element.
func (b *B) H6(args ...any) { b.element("h6", args...) }

// Hgroup writes an <hgroup> element.
func (b *B) Hgroup(args ...any) { b.element("hgroup", args...) }

// Main writes a <main> element.
func (b *B) Main(args ...any) { b.element("main", args...) }

// Nav writes a <nav> element.
func (b *B) Nav(args ...any) { b.element("nav", args...) }

// Section writes a <section> element.
func (b *B) Section(args ...any) { b.element("section", args...) }

// Search writes a <search> element.
func (b *B) Search(args ...any) { b.element("search", args...) }

// Blockquote writes a <blockquote> element.
func (b *B) Blockquote(args ...any) { b.element("blockquote", args...) }

// Dd writes a <dd> element.
func (b *B) Dd(args ...any) { b.element("dd", args...) }

// Div writes a <div> element.
func (b *B) Div(args ...any) { b.element("div", args...) }

// Dl writes a <dl> element.
func (b *B) Dl(args ...any) { b.element("dl", args...) }

// Dt writes a <dt> element.
func (b *B) Dt(args ...any) { b.element("dt", args...) }

// Figcaption writes a <figcaption> element.
func (b *B) Figcaption(args ...any) { b.element("figcaption", args...) }

// Figure writes a <figure> element.
func (b *B) Figure(args ...any) { b.element("figure", args...) }

// Hr writes a self-closing <hr> element.
func (b *B) Hr(args ...any) { b.voidElement("hr", args...) }

// Li writes an <li> element.
func (b *B) Li(args ...any) { b.element("li", args...) }

// Menu writes a <menu> element.
func (b *B) Menu(args ...any) { b.element("menu", args...) }

// Ol writes an <ol> element.
func (b *B) Ol(args ...any) { b.element("ol", args...) }

// P writes a <p> element.
func (b *B) P(args ...any) { b.element("p", args...) }

// Pre writes a <pre> element.
func (b *B) Pre(args ...any) { b.element("pre", args...) }

// Ul writes a <ul> element.
func (b *B) Ul(args ...any) { b.element("ul", args...) }

// A writes an <a> element.
func (b *B) A(args ...any) { b.element("a", args...) }

// Abbr writes an <abbr> element.
func (b *B) Abbr(args ...any) { b.element("abbr", args...) }

// B writes a <b> element.
func (b *B) B(args ...any) { b.element("b", args...) }

// Bdi writes a <bdi> element.
func (b *B) Bdi(args ...any) { b.element("bdi", args...) }

// Bdo writes a <bdo> element.
func (b *B) Bdo(args ...any) { b.element("bdo", args...) }

// Br writes a self-closing <br> element.
func (b *B) Br(args ...any) { b.voidElement("br", args...) }

// Cite writes a <cite> element.
func (b *B) Cite(args ...any) { b.element("cite", args...) }

// Code writes a <code> element.
func (b *B) Code(args ...any) { b.element("code", args...) }

// Data writes a <data> element.
func (b *B) Data(args ...any) { b.element("data", args...) }

// Dfn writes a <dfn> element.
func (b *B) Dfn(args ...any) { b.element("dfn", args...) }

// Em writes an <em> element.
func (b *B) Em(args ...any) { b.element("em", args...) }

// I writes an <i> element.
func (b *B) I(args ...any) { b.element("i", args...) }

// Kbd writes a <kbd> element.
func (b *B) Kbd(args ...any) { b.element("kbd", args...) }

// Mark writes a <mark> element.
func (b *B) Mark(args ...any) { b.element("mark", args...) }

// Q writes a <q> element.
func (b *B) Q(args ...any) { b.element("q", args...) }

// Rp writes an <rp> element.
func (b *B) Rp(args ...any) { b.element("rp", args...) }

// Rt writes an <rt> element.
func (b *B) Rt(args ...any) { b.element("rt", args...) }

// Ruby writes a <ruby> element.
func (b *B) Ruby(args ...any) { b.element("ruby", args...) }

// S writes an <s> element.
func (b *B) S(args ...any) { b.element("s", args...) }

// Samp writes a <samp> element.
func (b *B) Samp(args ...any) { b.element("samp", args...) }

// Small writes a <small> element.
func (b *B) Small(args ...any) { b.element("small", args...) }

// Span writes a <span> element.
func (b *B) Span(args ...any) { b.element("span", args...) }

// Strong writes a <strong> element.
func (b *B) Strong(args ...any) { b.element("strong", args...) }

// Sub writes a <sub> element.
func (b *B) Sub(args ...any) { b.element("sub", args...) }

// Sup writes a <sup> element.
func (b *B) Sup(args ...any) { b.element("sup", args...) }

// Time writes a <time> element.
func (b *B) Time(args ...any) { b.element("time", args...) }

// U writes a <u> element.
func (b *B) U(args ...any) { b.element("u", args...) }

// Var writes a <var> element.
func (b *B) Var(args ...any) { b.element("var", args...) }

// Wbr writes a self-closing <wbr> element.
func (b *B) Wbr(args ...any) { b.voidElement("wbr", args...) }

// Area writes a self-closing <area> element.
func (b *B) Area(args ...any) { b.voidElement("area", args...) }

// Audio writes an <audio> element.
func (b *B) Audio(args ...any) { b.element("audio", args...) }

// Img writes a self-closing <img> element.
func (b *B) Img(args ...any) { b.voidElement("img", args...) }

// Map writes a <map> element.
func (b *B) Map(args ...any) { b.element("map", args...) }

// Track writes a self-closing <track> element.
func (b *B) Track(args ...any) { b.voidElement("track", args...) }

// Video writes a <video> element.
func (b *B) Video(args ...any) { b.element("video", args...) }

// Embed writes a self-closing <embed> element.
func (b *B) Embed(args ...any) { b.voidElement("embed", args...) }

// Iframe writes an <iframe> element.
func (b *B) Iframe(args ...any) { b.element("iframe", args...) }

// Object writes an <object> element.
func (b *B) Object(args ...any) { b.element("object", args...) }

// Picture writes a <picture> element.
func (b *B) Picture(args ...any) { b.element("picture", args...) }

// Portal writes a <portal> element.
func (b *B) Portal(args ...any) { b.element("portal", args...) }

// Source writes a self-closing <source> element.
func (b *B) Source(args ...any) { b.voidElement("source", args...) }

// Svg writes an <svg> element.
func (b *B) Svg(args ...any) { b.element("svg", args...) }

// Math writes a <math> element.
func (b *B) Math(args ...any) { b.element("math", args...) }

// Canvas writes a <canvas> element.
func (b *B) Canvas(args ...any) { b.element("canvas", args...) }

// Template writes a <template> element.
func (b *B) Template(args ...any) { b.element("template", args...) }

// Slot writes a <slot> element.
func (b *B) Slot(args ...any) { b.element("slot", args...) }

// Del writes a <del> element.
func (b *B) Del(args ...any) { b.element("del", args...) }

// Ins writes an <ins> element.
func (b *B) Ins(args ...any) { b.element("ins", args...) }

// Caption writes a <caption> element.
func (b *B) Caption(args ...any) { b.element("caption", args...) }

// Col writes a self-closing <col> element.
func (b *B) Col(args ...any) { b.voidElement("col", args...) }

// Colgroup writes a <colgroup> element.
func (b *B) Colgroup(args ...any) { b.element("colgroup", args...) }

// Table writes a <table> element.
func (b *B) Table(args ...any) { b.element("table", args...) }

// Tbody writes a <tbody> element.
func (b *B) Tbody(args ...any) { b.element("tbody", args...) }

// Td writes a <td> element.
func (b *B) Td(args ...any) { b.element("td", args...) }

// Tfoot writes a <tfoot> element.
func (b *B) Tfoot(args ...any) { b.element("tfoot", args...) }

// Th writes a <th> element.
func (b *B) Th(args ...any) { b.element("th", args...) }

// Thead writes a <thead> element.
func (b *B) Thead(args ...any) { b.element("thead", args...) }

// Tr writes a <tr> element.
func (b *B) Tr(args ...any) { b.element("tr", args...) }

// Button writes a <button> element.
func (b *B) Button(args ...any) { b.element("button", args...) }

// Datalist writes a <datalist> element.
func (b *B) Datalist(args ...any) { b.element("datalist", args...) }

// Fieldset writes a <fieldset> element.
func (b *B) Fieldset(args ...any) { b.element("fieldset", args...) }

// Form writes a <form> element.
func (b *B) Form(args ...any) { b.element("form", args...) }

// Input writes a self-closing <input> element.
func (b *B) Input(args ...any) { b.voidElement("input", args...) }

// Label writes a <label> element.
func (b *B) Label(args ...any) { b.element("label", args...) }

// Legend writes a <legend> element.
func (b *B) Legend(args ...any) { b.element("legend", args...) }

// Meter writes a <meter> element.
func (b *B) Meter(args ...any) { b.element("meter", args...) }

// Optgroup writes an <optgroup> element.
func (b *B) Optgroup(args ...any) { b.element("optgroup", args...) }

// Option writes an <option> element.
func (b *B) Option(args ...any) { b.element("option", args...) }

// Output writes an <output> element.
func (b *B) Output(args ...any) { b.element("output", args...) }

// Progress writes a <progress> element.
func (b *B) Progress(args ...any) { b.element("progress", args...) }

// Select writes a <select> element.
func (b *B) Select(args ...any) { b.element("select", args...) }

// Textarea writes a <textarea> element.
func (b *B) Textarea(args ...any) { b.element("textarea", args...) }

// Details writes a <details> element.
func (b *B) Details(args ...any) { b.element("details", args...) }

// Dialog writes a <dialog> element.
func (b *B) Dialog(args ...any) { b.element("dialog", args...) }

// Summary writes a <summary> element.
func (b *B) Summary(args ...any) { b.element("summary", args...) }
