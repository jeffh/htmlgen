package h

// Html writes the root <html> element, preceded by the HTML5 doctype. It uses
// lang="en" unless a lang attribute is provided.
func (b *B) Html(attrs Attributes, body Body) {
	if b.err != nil {
		return
	}
	if _, ok := attrs.Get("lang"); !ok {
		// Full slice expression so appending never writes into a caller-owned
		// slice's spare capacity.
		attrs = append(attrs[:len(attrs):len(attrs)], Attribute{Name: "lang", Value: "en"})
	}
	b.Doctype()
	b.element("<html", "</html>", attrs, body)
}

// Head writes a <head> element.
func (b *B) Head(attrs Attributes, body Body) { b.element("<head", "</head>", attrs, body) }

// Title writes a <title> element.
func (b *B) Title(attrs Attributes, body Body) { b.element("<title", "</title>", attrs, body) }

// Meta writes a self-closing <meta> element.
func (b *B) Meta(attrs Attributes) { b.voidElement("<meta", attrs) }

// Link writes a self-closing <link> element.
func (b *B) Link(attrs Attributes) { b.voidElement("<link", attrs) }

// Style writes a <style> element.
func (b *B) Style(attrs Attributes, body Body) { b.element("<style", "</style>", attrs, body) }

// Script writes a <script> element.
func (b *B) Script(attrs Attributes, body Body) { b.element("<script", "</script>", attrs, body) }

// Noscript writes a <noscript> element.
func (b *B) Noscript(attrs Attributes, body Body) { b.element("<noscript", "</noscript>", attrs, body) }

// Base writes a self-closing <base> element.
func (b *B) Base(attrs Attributes) { b.voidElement("<base", attrs) }

// Body writes a <body> element.
func (b *B) Body(attrs Attributes, body Body) { b.element("<body", "</body>", attrs, body) }

// Address writes an <address> element.
func (b *B) Address(attrs Attributes, body Body) { b.element("<address", "</address>", attrs, body) }

// Article writes an <article> element.
func (b *B) Article(attrs Attributes, body Body) { b.element("<article", "</article>", attrs, body) }

// Aside writes an <aside> element.
func (b *B) Aside(attrs Attributes, body Body) { b.element("<aside", "</aside>", attrs, body) }

// Footer writes a <footer> element.
func (b *B) Footer(attrs Attributes, body Body) { b.element("<footer", "</footer>", attrs, body) }

// Header writes a <header> element.
func (b *B) Header(attrs Attributes, body Body) { b.element("<header", "</header>", attrs, body) }

// H1 writes an <h1> element.
func (b *B) H1(attrs Attributes, body Body) { b.element("<h1", "</h1>", attrs, body) }

// H2 writes an <h2> element.
func (b *B) H2(attrs Attributes, body Body) { b.element("<h2", "</h2>", attrs, body) }

// H3 writes an <h3> element.
func (b *B) H3(attrs Attributes, body Body) { b.element("<h3", "</h3>", attrs, body) }

// H4 writes an <h4> element.
func (b *B) H4(attrs Attributes, body Body) { b.element("<h4", "</h4>", attrs, body) }

// H5 writes an <h5> element.
func (b *B) H5(attrs Attributes, body Body) { b.element("<h5", "</h5>", attrs, body) }

// H6 writes an <h6> element.
func (b *B) H6(attrs Attributes, body Body) { b.element("<h6", "</h6>", attrs, body) }

// Hgroup writes an <hgroup> element.
func (b *B) Hgroup(attrs Attributes, body Body) { b.element("<hgroup", "</hgroup>", attrs, body) }

// Main writes a <main> element.
func (b *B) Main(attrs Attributes, body Body) { b.element("<main", "</main>", attrs, body) }

// Nav writes a <nav> element.
func (b *B) Nav(attrs Attributes, body Body) { b.element("<nav", "</nav>", attrs, body) }

// Section writes a <section> element.
func (b *B) Section(attrs Attributes, body Body) { b.element("<section", "</section>", attrs, body) }

// Search writes a <search> element.
func (b *B) Search(attrs Attributes, body Body) { b.element("<search", "</search>", attrs, body) }

// Blockquote writes a <blockquote> element.
func (b *B) Blockquote(attrs Attributes, body Body) {
	b.element("<blockquote", "</blockquote>", attrs, body)
}

// Dd writes a <dd> element.
func (b *B) Dd(attrs Attributes, body Body) { b.element("<dd", "</dd>", attrs, body) }

// Div writes a <div> element.
func (b *B) Div(attrs Attributes, body Body) { b.element("<div", "</div>", attrs, body) }

// Dl writes a <dl> element.
func (b *B) Dl(attrs Attributes, body Body) { b.element("<dl", "</dl>", attrs, body) }

// Dt writes a <dt> element.
func (b *B) Dt(attrs Attributes, body Body) { b.element("<dt", "</dt>", attrs, body) }

// Figcaption writes a <figcaption> element.
func (b *B) Figcaption(attrs Attributes, body Body) {
	b.element("<figcaption", "</figcaption>", attrs, body)
}

// Figure writes a <figure> element.
func (b *B) Figure(attrs Attributes, body Body) { b.element("<figure", "</figure>", attrs, body) }

// Hr writes a self-closing <hr> element.
func (b *B) Hr(attrs Attributes) { b.voidElement("<hr", attrs) }

// Li writes an <li> element.
func (b *B) Li(attrs Attributes, body Body) { b.element("<li", "</li>", attrs, body) }

// Menu writes a <menu> element.
func (b *B) Menu(attrs Attributes, body Body) { b.element("<menu", "</menu>", attrs, body) }

// Ol writes an <ol> element.
func (b *B) Ol(attrs Attributes, body Body) { b.element("<ol", "</ol>", attrs, body) }

// P writes a <p> element.
func (b *B) P(attrs Attributes, body Body) { b.element("<p", "</p>", attrs, body) }

// Pre writes a <pre> element.
func (b *B) Pre(attrs Attributes, body Body) { b.element("<pre", "</pre>", attrs, body) }

// Ul writes a <ul> element.
func (b *B) Ul(attrs Attributes, body Body) { b.element("<ul", "</ul>", attrs, body) }

// A writes an <a> element.
func (b *B) A(attrs Attributes, body Body) { b.element("<a", "</a>", attrs, body) }

// Abbr writes an <abbr> element.
func (b *B) Abbr(attrs Attributes, body Body) { b.element("<abbr", "</abbr>", attrs, body) }

// B writes a <b> element.
func (b *B) B(attrs Attributes, body Body) { b.element("<b", "</b>", attrs, body) }

// Bdi writes a <bdi> element.
func (b *B) Bdi(attrs Attributes, body Body) { b.element("<bdi", "</bdi>", attrs, body) }

// Bdo writes a <bdo> element.
func (b *B) Bdo(attrs Attributes, body Body) { b.element("<bdo", "</bdo>", attrs, body) }

// Br writes a self-closing <br> element.
func (b *B) Br(attrs Attributes) { b.voidElement("<br", attrs) }

// Cite writes a <cite> element.
func (b *B) Cite(attrs Attributes, body Body) { b.element("<cite", "</cite>", attrs, body) }

// Code writes a <code> element.
func (b *B) Code(attrs Attributes, body Body) { b.element("<code", "</code>", attrs, body) }

// Data writes a <data> element.
func (b *B) Data(attrs Attributes, body Body) { b.element("<data", "</data>", attrs, body) }

// Dfn writes a <dfn> element.
func (b *B) Dfn(attrs Attributes, body Body) { b.element("<dfn", "</dfn>", attrs, body) }

// Em writes an <em> element.
func (b *B) Em(attrs Attributes, body Body) { b.element("<em", "</em>", attrs, body) }

// I writes an <i> element.
func (b *B) I(attrs Attributes, body Body) { b.element("<i", "</i>", attrs, body) }

// Kbd writes a <kbd> element.
func (b *B) Kbd(attrs Attributes, body Body) { b.element("<kbd", "</kbd>", attrs, body) }

// Mark writes a <mark> element.
func (b *B) Mark(attrs Attributes, body Body) { b.element("<mark", "</mark>", attrs, body) }

// Q writes a <q> element.
func (b *B) Q(attrs Attributes, body Body) { b.element("<q", "</q>", attrs, body) }

// Rp writes an <rp> element.
func (b *B) Rp(attrs Attributes, body Body) { b.element("<rp", "</rp>", attrs, body) }

// Rt writes an <rt> element.
func (b *B) Rt(attrs Attributes, body Body) { b.element("<rt", "</rt>", attrs, body) }

// Ruby writes a <ruby> element.
func (b *B) Ruby(attrs Attributes, body Body) { b.element("<ruby", "</ruby>", attrs, body) }

// S writes an <s> element.
func (b *B) S(attrs Attributes, body Body) { b.element("<s", "</s>", attrs, body) }

// Samp writes a <samp> element.
func (b *B) Samp(attrs Attributes, body Body) { b.element("<samp", "</samp>", attrs, body) }

// Small writes a <small> element.
func (b *B) Small(attrs Attributes, body Body) { b.element("<small", "</small>", attrs, body) }

// Span writes a <span> element.
func (b *B) Span(attrs Attributes, body Body) { b.element("<span", "</span>", attrs, body) }

// Strong writes a <strong> element.
func (b *B) Strong(attrs Attributes, body Body) { b.element("<strong", "</strong>", attrs, body) }

// Sub writes a <sub> element.
func (b *B) Sub(attrs Attributes, body Body) { b.element("<sub", "</sub>", attrs, body) }

// Sup writes a <sup> element.
func (b *B) Sup(attrs Attributes, body Body) { b.element("<sup", "</sup>", attrs, body) }

// Time writes a <time> element.
func (b *B) Time(attrs Attributes, body Body) { b.element("<time", "</time>", attrs, body) }

// U writes a <u> element.
func (b *B) U(attrs Attributes, body Body) { b.element("<u", "</u>", attrs, body) }

// Var writes a <var> element.
func (b *B) Var(attrs Attributes, body Body) { b.element("<var", "</var>", attrs, body) }

// Wbr writes a self-closing <wbr> element.
func (b *B) Wbr(attrs Attributes) { b.voidElement("<wbr", attrs) }

// Area writes a self-closing <area> element.
func (b *B) Area(attrs Attributes) { b.voidElement("<area", attrs) }

// Audio writes an <audio> element.
func (b *B) Audio(attrs Attributes, body Body) { b.element("<audio", "</audio>", attrs, body) }

// Img writes a self-closing <img> element.
func (b *B) Img(attrs Attributes) { b.voidElement("<img", attrs) }

// Map writes a <map> element.
func (b *B) Map(attrs Attributes, body Body) { b.element("<map", "</map>", attrs, body) }

// Track writes a self-closing <track> element.
func (b *B) Track(attrs Attributes) { b.voidElement("<track", attrs) }

// Video writes a <video> element.
func (b *B) Video(attrs Attributes, body Body) { b.element("<video", "</video>", attrs, body) }

// Embed writes a self-closing <embed> element.
func (b *B) Embed(attrs Attributes) { b.voidElement("<embed", attrs) }

// Iframe writes an <iframe> element.
func (b *B) Iframe(attrs Attributes, body Body) { b.element("<iframe", "</iframe>", attrs, body) }

// Object writes an <object> element.
func (b *B) Object(attrs Attributes, body Body) { b.element("<object", "</object>", attrs, body) }

// Picture writes a <picture> element.
func (b *B) Picture(attrs Attributes, body Body) { b.element("<picture", "</picture>", attrs, body) }

// Portal writes a <portal> element.
func (b *B) Portal(attrs Attributes, body Body) { b.element("<portal", "</portal>", attrs, body) }

// Source writes a self-closing <source> element.
func (b *B) Source(attrs Attributes) { b.voidElement("<source", attrs) }

// Svg writes an <svg> element.
func (b *B) Svg(attrs Attributes, body Body) { b.element("<svg", "</svg>", attrs, body) }

// Math writes a <math> element.
func (b *B) Math(attrs Attributes, body Body) { b.element("<math", "</math>", attrs, body) }

// Canvas writes a <canvas> element.
func (b *B) Canvas(attrs Attributes, body Body) { b.element("<canvas", "</canvas>", attrs, body) }

// Template writes a <template> element.
func (b *B) Template(attrs Attributes, body Body) { b.element("<template", "</template>", attrs, body) }

// Slot writes a <slot> element.
func (b *B) Slot(attrs Attributes, body Body) { b.element("<slot", "</slot>", attrs, body) }

// Del writes a <del> element.
func (b *B) Del(attrs Attributes, body Body) { b.element("<del", "</del>", attrs, body) }

// Ins writes an <ins> element.
func (b *B) Ins(attrs Attributes, body Body) { b.element("<ins", "</ins>", attrs, body) }

// Caption writes a <caption> element.
func (b *B) Caption(attrs Attributes, body Body) { b.element("<caption", "</caption>", attrs, body) }

// Col writes a self-closing <col> element.
func (b *B) Col(attrs Attributes) { b.voidElement("<col", attrs) }

// Colgroup writes a <colgroup> element.
func (b *B) Colgroup(attrs Attributes, body Body) { b.element("<colgroup", "</colgroup>", attrs, body) }

// Table writes a <table> element.
func (b *B) Table(attrs Attributes, body Body) { b.element("<table", "</table>", attrs, body) }

// Tbody writes a <tbody> element.
func (b *B) Tbody(attrs Attributes, body Body) { b.element("<tbody", "</tbody>", attrs, body) }

// Td writes a <td> element.
func (b *B) Td(attrs Attributes, body Body) { b.element("<td", "</td>", attrs, body) }

// Tfoot writes a <tfoot> element.
func (b *B) Tfoot(attrs Attributes, body Body) { b.element("<tfoot", "</tfoot>", attrs, body) }

// Th writes a <th> element.
func (b *B) Th(attrs Attributes, body Body) { b.element("<th", "</th>", attrs, body) }

// Thead writes a <thead> element.
func (b *B) Thead(attrs Attributes, body Body) { b.element("<thead", "</thead>", attrs, body) }

// Tr writes a <tr> element.
func (b *B) Tr(attrs Attributes, body Body) { b.element("<tr", "</tr>", attrs, body) }

// Button writes a <button> element.
func (b *B) Button(attrs Attributes, body Body) { b.element("<button", "</button>", attrs, body) }

// Datalist writes a <datalist> element.
func (b *B) Datalist(attrs Attributes, body Body) { b.element("<datalist", "</datalist>", attrs, body) }

// Fieldset writes a <fieldset> element.
func (b *B) Fieldset(attrs Attributes, body Body) { b.element("<fieldset", "</fieldset>", attrs, body) }

// Form writes a <form> element.
func (b *B) Form(attrs Attributes, body Body) { b.element("<form", "</form>", attrs, body) }

// Input writes a self-closing <input> element.
func (b *B) Input(attrs Attributes) { b.voidElement("<input", attrs) }

// Label writes a <label> element.
func (b *B) Label(attrs Attributes, body Body) { b.element("<label", "</label>", attrs, body) }

// Legend writes a <legend> element.
func (b *B) Legend(attrs Attributes, body Body) { b.element("<legend", "</legend>", attrs, body) }

// Meter writes a <meter> element.
func (b *B) Meter(attrs Attributes, body Body) { b.element("<meter", "</meter>", attrs, body) }

// Optgroup writes an <optgroup> element.
func (b *B) Optgroup(attrs Attributes, body Body) { b.element("<optgroup", "</optgroup>", attrs, body) }

// Option writes an <option> element.
func (b *B) Option(attrs Attributes, body Body) { b.element("<option", "</option>", attrs, body) }

// Output writes an <output> element.
func (b *B) Output(attrs Attributes, body Body) { b.element("<output", "</output>", attrs, body) }

// Progress writes a <progress> element.
func (b *B) Progress(attrs Attributes, body Body) { b.element("<progress", "</progress>", attrs, body) }

// Select writes a <select> element.
func (b *B) Select(attrs Attributes, body Body) { b.element("<select", "</select>", attrs, body) }

// Textarea writes a <textarea> element.
func (b *B) Textarea(attrs Attributes, body Body) { b.element("<textarea", "</textarea>", attrs, body) }

// Details writes a <details> element.
func (b *B) Details(attrs Attributes, body Body) { b.element("<details", "</details>", attrs, body) }

// Dialog writes a <dialog> element.
func (b *B) Dialog(attrs Attributes, body Body) { b.element("<dialog", "</dialog>", attrs, body) }

// Summary writes a <summary> element.
func (b *B) Summary(attrs Attributes, body Body) { b.element("<summary", "</summary>", attrs, body) }
