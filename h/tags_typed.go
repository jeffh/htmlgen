package h

// Typed siblings of the variadic element methods in tags.go. Each XxxE takes
// concrete attribute and body parameters instead of ...any, so arguments are
// never boxed into interfaces and a capturing body closure stays on the stack.
// Behavior is otherwise identical; nil attrs and a nil body are both valid.

// HtmlE writes the root <html> element, preceded by the HTML5 doctype, using
// lang="en" unless attrs provides one. It is the non-boxing fast path of Html.
func (b *B) HtmlE(attrs Attributes, body Body) {
	if b.err != nil {
		return
	}
	if _, ok := attrs.Get("lang"); !ok {
		// Full slice expression so appending never writes into a caller-owned
		// slice's spare capacity.
		attrs = append(attrs[:len(attrs):len(attrs)], Attribute{Name: "lang", Value: "en"})
	}
	b.Doctype()
	b.elementTyped("<html", "</html>", attrs, body)
}

// HeadE writes a <head> element. It is the non-boxing fast path of Head.
func (b *B) HeadE(attrs Attributes, body Body) { b.elementTyped("<head", "</head>", attrs, body) }

// TitleE writes a <title> element. It is the non-boxing fast path of Title.
func (b *B) TitleE(attrs Attributes, body Body) { b.elementTyped("<title", "</title>", attrs, body) }

// MetaE writes a self-closing <meta> element. It is the non-boxing fast path of Meta.
func (b *B) MetaE(attrs Attributes) { b.voidElementTyped("<meta", attrs) }

// LinkE writes a self-closing <link> element. It is the non-boxing fast path of Link.
func (b *B) LinkE(attrs Attributes) { b.voidElementTyped("<link", attrs) }

// StyleE writes a <style> element. It is the non-boxing fast path of Style.
func (b *B) StyleE(attrs Attributes, body Body) { b.elementTyped("<style", "</style>", attrs, body) }

// ScriptE writes a <script> element. It is the non-boxing fast path of Script.
func (b *B) ScriptE(attrs Attributes, body Body) { b.elementTyped("<script", "</script>", attrs, body) }

// NoscriptE writes a <noscript> element. It is the non-boxing fast path of Noscript.
func (b *B) NoscriptE(attrs Attributes, body Body) {
	b.elementTyped("<noscript", "</noscript>", attrs, body)
}

// BaseE writes a self-closing <base> element. It is the non-boxing fast path of Base.
func (b *B) BaseE(attrs Attributes) { b.voidElementTyped("<base", attrs) }

// BodyE writes a <body> element. It is the non-boxing fast path of Body.
func (b *B) BodyE(attrs Attributes, body Body) { b.elementTyped("<body", "</body>", attrs, body) }

// AddressE writes an <address> element. It is the non-boxing fast path of Address.
func (b *B) AddressE(attrs Attributes, body Body) {
	b.elementTyped("<address", "</address>", attrs, body)
}

// ArticleE writes an <article> element. It is the non-boxing fast path of Article.
func (b *B) ArticleE(attrs Attributes, body Body) {
	b.elementTyped("<article", "</article>", attrs, body)
}

// AsideE writes an <aside> element. It is the non-boxing fast path of Aside.
func (b *B) AsideE(attrs Attributes, body Body) { b.elementTyped("<aside", "</aside>", attrs, body) }

// FooterE writes a <footer> element. It is the non-boxing fast path of Footer.
func (b *B) FooterE(attrs Attributes, body Body) { b.elementTyped("<footer", "</footer>", attrs, body) }

// HeaderE writes a <header> element. It is the non-boxing fast path of Header.
func (b *B) HeaderE(attrs Attributes, body Body) { b.elementTyped("<header", "</header>", attrs, body) }

// H1E writes a <h1> element. It is the non-boxing fast path of H1.
func (b *B) H1E(attrs Attributes, body Body) { b.elementTyped("<h1", "</h1>", attrs, body) }

// H2E writes a <h2> element. It is the non-boxing fast path of H2.
func (b *B) H2E(attrs Attributes, body Body) { b.elementTyped("<h2", "</h2>", attrs, body) }

// H3E writes a <h3> element. It is the non-boxing fast path of H3.
func (b *B) H3E(attrs Attributes, body Body) { b.elementTyped("<h3", "</h3>", attrs, body) }

// H4E writes a <h4> element. It is the non-boxing fast path of H4.
func (b *B) H4E(attrs Attributes, body Body) { b.elementTyped("<h4", "</h4>", attrs, body) }

// H5E writes a <h5> element. It is the non-boxing fast path of H5.
func (b *B) H5E(attrs Attributes, body Body) { b.elementTyped("<h5", "</h5>", attrs, body) }

// H6E writes a <h6> element. It is the non-boxing fast path of H6.
func (b *B) H6E(attrs Attributes, body Body) { b.elementTyped("<h6", "</h6>", attrs, body) }

// HgroupE writes a <hgroup> element. It is the non-boxing fast path of Hgroup.
func (b *B) HgroupE(attrs Attributes, body Body) { b.elementTyped("<hgroup", "</hgroup>", attrs, body) }

// MainE writes a <main> element. It is the non-boxing fast path of Main.
func (b *B) MainE(attrs Attributes, body Body) { b.elementTyped("<main", "</main>", attrs, body) }

// NavE writes a <nav> element. It is the non-boxing fast path of Nav.
func (b *B) NavE(attrs Attributes, body Body) { b.elementTyped("<nav", "</nav>", attrs, body) }

// SectionE writes a <section> element. It is the non-boxing fast path of Section.
func (b *B) SectionE(attrs Attributes, body Body) {
	b.elementTyped("<section", "</section>", attrs, body)
}

// SearchE writes a <search> element. It is the non-boxing fast path of Search.
func (b *B) SearchE(attrs Attributes, body Body) { b.elementTyped("<search", "</search>", attrs, body) }

// BlockquoteE writes a <blockquote> element. It is the non-boxing fast path of Blockquote.
func (b *B) BlockquoteE(attrs Attributes, body Body) {
	b.elementTyped("<blockquote", "</blockquote>", attrs, body)
}

// DdE writes a <dd> element. It is the non-boxing fast path of Dd.
func (b *B) DdE(attrs Attributes, body Body) { b.elementTyped("<dd", "</dd>", attrs, body) }

// DivE writes a <div> element. It is the non-boxing fast path of Div.
func (b *B) DivE(attrs Attributes, body Body) { b.elementTyped("<div", "</div>", attrs, body) }

// DlE writes a <dl> element. It is the non-boxing fast path of Dl.
func (b *B) DlE(attrs Attributes, body Body) { b.elementTyped("<dl", "</dl>", attrs, body) }

// DtE writes a <dt> element. It is the non-boxing fast path of Dt.
func (b *B) DtE(attrs Attributes, body Body) { b.elementTyped("<dt", "</dt>", attrs, body) }

// FigcaptionE writes a <figcaption> element. It is the non-boxing fast path of Figcaption.
func (b *B) FigcaptionE(attrs Attributes, body Body) {
	b.elementTyped("<figcaption", "</figcaption>", attrs, body)
}

// FigureE writes a <figure> element. It is the non-boxing fast path of Figure.
func (b *B) FigureE(attrs Attributes, body Body) { b.elementTyped("<figure", "</figure>", attrs, body) }

// HrE writes a self-closing <hr> element. It is the non-boxing fast path of Hr.
func (b *B) HrE(attrs Attributes) { b.voidElementTyped("<hr", attrs) }

// LiE writes a <li> element. It is the non-boxing fast path of Li.
func (b *B) LiE(attrs Attributes, body Body) { b.elementTyped("<li", "</li>", attrs, body) }

// MenuE writes a <menu> element. It is the non-boxing fast path of Menu.
func (b *B) MenuE(attrs Attributes, body Body) { b.elementTyped("<menu", "</menu>", attrs, body) }

// OlE writes an <ol> element. It is the non-boxing fast path of Ol.
func (b *B) OlE(attrs Attributes, body Body) { b.elementTyped("<ol", "</ol>", attrs, body) }

// PE writes a <p> element. It is the non-boxing fast path of P.
func (b *B) PE(attrs Attributes, body Body) { b.elementTyped("<p", "</p>", attrs, body) }

// PreE writes a <pre> element. It is the non-boxing fast path of Pre.
func (b *B) PreE(attrs Attributes, body Body) { b.elementTyped("<pre", "</pre>", attrs, body) }

// UlE writes an <ul> element. It is the non-boxing fast path of Ul.
func (b *B) UlE(attrs Attributes, body Body) { b.elementTyped("<ul", "</ul>", attrs, body) }

// AE writes an <a> element. It is the non-boxing fast path of A.
func (b *B) AE(attrs Attributes, body Body) { b.elementTyped("<a", "</a>", attrs, body) }

// AbbrE writes an <abbr> element. It is the non-boxing fast path of Abbr.
func (b *B) AbbrE(attrs Attributes, body Body) { b.elementTyped("<abbr", "</abbr>", attrs, body) }

// BE writes a <b> element. It is the non-boxing fast path of B.
func (b *B) BE(attrs Attributes, body Body) { b.elementTyped("<b", "</b>", attrs, body) }

// BdiE writes a <bdi> element. It is the non-boxing fast path of Bdi.
func (b *B) BdiE(attrs Attributes, body Body) { b.elementTyped("<bdi", "</bdi>", attrs, body) }

// BdoE writes a <bdo> element. It is the non-boxing fast path of Bdo.
func (b *B) BdoE(attrs Attributes, body Body) { b.elementTyped("<bdo", "</bdo>", attrs, body) }

// BrE writes a self-closing <br> element. It is the non-boxing fast path of Br.
func (b *B) BrE(attrs Attributes) { b.voidElementTyped("<br", attrs) }

// CiteE writes a <cite> element. It is the non-boxing fast path of Cite.
func (b *B) CiteE(attrs Attributes, body Body) { b.elementTyped("<cite", "</cite>", attrs, body) }

// CodeE writes a <code> element. It is the non-boxing fast path of Code.
func (b *B) CodeE(attrs Attributes, body Body) { b.elementTyped("<code", "</code>", attrs, body) }

// DataE writes a <data> element. It is the non-boxing fast path of Data.
func (b *B) DataE(attrs Attributes, body Body) { b.elementTyped("<data", "</data>", attrs, body) }

// DfnE writes a <dfn> element. It is the non-boxing fast path of Dfn.
func (b *B) DfnE(attrs Attributes, body Body) { b.elementTyped("<dfn", "</dfn>", attrs, body) }

// EmE writes an <em> element. It is the non-boxing fast path of Em.
func (b *B) EmE(attrs Attributes, body Body) { b.elementTyped("<em", "</em>", attrs, body) }

// IE writes an <i> element. It is the non-boxing fast path of I.
func (b *B) IE(attrs Attributes, body Body) { b.elementTyped("<i", "</i>", attrs, body) }

// KbdE writes a <kbd> element. It is the non-boxing fast path of Kbd.
func (b *B) KbdE(attrs Attributes, body Body) { b.elementTyped("<kbd", "</kbd>", attrs, body) }

// MarkE writes a <mark> element. It is the non-boxing fast path of Mark.
func (b *B) MarkE(attrs Attributes, body Body) { b.elementTyped("<mark", "</mark>", attrs, body) }

// QE writes a <q> element. It is the non-boxing fast path of Q.
func (b *B) QE(attrs Attributes, body Body) { b.elementTyped("<q", "</q>", attrs, body) }

// RpE writes a <rp> element. It is the non-boxing fast path of Rp.
func (b *B) RpE(attrs Attributes, body Body) { b.elementTyped("<rp", "</rp>", attrs, body) }

// RtE writes a <rt> element. It is the non-boxing fast path of Rt.
func (b *B) RtE(attrs Attributes, body Body) { b.elementTyped("<rt", "</rt>", attrs, body) }

// RubyE writes a <ruby> element. It is the non-boxing fast path of Ruby.
func (b *B) RubyE(attrs Attributes, body Body) { b.elementTyped("<ruby", "</ruby>", attrs, body) }

// SE writes a <s> element. It is the non-boxing fast path of S.
func (b *B) SE(attrs Attributes, body Body) { b.elementTyped("<s", "</s>", attrs, body) }

// SampE writes a <samp> element. It is the non-boxing fast path of Samp.
func (b *B) SampE(attrs Attributes, body Body) { b.elementTyped("<samp", "</samp>", attrs, body) }

// SmallE writes a <small> element. It is the non-boxing fast path of Small.
func (b *B) SmallE(attrs Attributes, body Body) { b.elementTyped("<small", "</small>", attrs, body) }

// SpanE writes a <span> element. It is the non-boxing fast path of Span.
func (b *B) SpanE(attrs Attributes, body Body) { b.elementTyped("<span", "</span>", attrs, body) }

// StrongE writes a <strong> element. It is the non-boxing fast path of Strong.
func (b *B) StrongE(attrs Attributes, body Body) { b.elementTyped("<strong", "</strong>", attrs, body) }

// SubE writes a <sub> element. It is the non-boxing fast path of Sub.
func (b *B) SubE(attrs Attributes, body Body) { b.elementTyped("<sub", "</sub>", attrs, body) }

// SupE writes a <sup> element. It is the non-boxing fast path of Sup.
func (b *B) SupE(attrs Attributes, body Body) { b.elementTyped("<sup", "</sup>", attrs, body) }

// TimeE writes a <time> element. It is the non-boxing fast path of Time.
func (b *B) TimeE(attrs Attributes, body Body) { b.elementTyped("<time", "</time>", attrs, body) }

// UE writes an <u> element. It is the non-boxing fast path of U.
func (b *B) UE(attrs Attributes, body Body) { b.elementTyped("<u", "</u>", attrs, body) }

// VarE writes a <var> element. It is the non-boxing fast path of Var.
func (b *B) VarE(attrs Attributes, body Body) { b.elementTyped("<var", "</var>", attrs, body) }

// WbrE writes a self-closing <wbr> element. It is the non-boxing fast path of Wbr.
func (b *B) WbrE(attrs Attributes) { b.voidElementTyped("<wbr", attrs) }

// AreaE writes an self-closing <area> element. It is the non-boxing fast path of Area.
func (b *B) AreaE(attrs Attributes) { b.voidElementTyped("<area", attrs) }

// AudioE writes an <audio> element. It is the non-boxing fast path of Audio.
func (b *B) AudioE(attrs Attributes, body Body) { b.elementTyped("<audio", "</audio>", attrs, body) }

// ImgE writes an self-closing <img> element. It is the non-boxing fast path of Img.
func (b *B) ImgE(attrs Attributes) { b.voidElementTyped("<img", attrs) }

// MapE writes a <map> element. It is the non-boxing fast path of Map.
func (b *B) MapE(attrs Attributes, body Body) { b.elementTyped("<map", "</map>", attrs, body) }

// TrackE writes a self-closing <track> element. It is the non-boxing fast path of Track.
func (b *B) TrackE(attrs Attributes) { b.voidElementTyped("<track", attrs) }

// VideoE writes a <video> element. It is the non-boxing fast path of Video.
func (b *B) VideoE(attrs Attributes, body Body) { b.elementTyped("<video", "</video>", attrs, body) }

// EmbedE writes an self-closing <embed> element. It is the non-boxing fast path of Embed.
func (b *B) EmbedE(attrs Attributes) { b.voidElementTyped("<embed", attrs) }

// IframeE writes an <iframe> element. It is the non-boxing fast path of Iframe.
func (b *B) IframeE(attrs Attributes, body Body) { b.elementTyped("<iframe", "</iframe>", attrs, body) }

// ObjectE writes an <object> element. It is the non-boxing fast path of Object.
func (b *B) ObjectE(attrs Attributes, body Body) { b.elementTyped("<object", "</object>", attrs, body) }

// PictureE writes a <picture> element. It is the non-boxing fast path of Picture.
func (b *B) PictureE(attrs Attributes, body Body) {
	b.elementTyped("<picture", "</picture>", attrs, body)
}

// PortalE writes a <portal> element. It is the non-boxing fast path of Portal.
func (b *B) PortalE(attrs Attributes, body Body) { b.elementTyped("<portal", "</portal>", attrs, body) }

// SourceE writes a self-closing <source> element. It is the non-boxing fast path of Source.
func (b *B) SourceE(attrs Attributes) { b.voidElementTyped("<source", attrs) }

// SvgE writes a <svg> element. It is the non-boxing fast path of Svg.
func (b *B) SvgE(attrs Attributes, body Body) { b.elementTyped("<svg", "</svg>", attrs, body) }

// MathE writes a <math> element. It is the non-boxing fast path of Math.
func (b *B) MathE(attrs Attributes, body Body) { b.elementTyped("<math", "</math>", attrs, body) }

// CanvasE writes a <canvas> element. It is the non-boxing fast path of Canvas.
func (b *B) CanvasE(attrs Attributes, body Body) { b.elementTyped("<canvas", "</canvas>", attrs, body) }

// TemplateE writes a <template> element. It is the non-boxing fast path of Template.
func (b *B) TemplateE(attrs Attributes, body Body) {
	b.elementTyped("<template", "</template>", attrs, body)
}

// SlotE writes a <slot> element. It is the non-boxing fast path of Slot.
func (b *B) SlotE(attrs Attributes, body Body) { b.elementTyped("<slot", "</slot>", attrs, body) }

// DelE writes a <del> element. It is the non-boxing fast path of Del.
func (b *B) DelE(attrs Attributes, body Body) { b.elementTyped("<del", "</del>", attrs, body) }

// InsE writes an <ins> element. It is the non-boxing fast path of Ins.
func (b *B) InsE(attrs Attributes, body Body) { b.elementTyped("<ins", "</ins>", attrs, body) }

// CaptionE writes a <caption> element. It is the non-boxing fast path of Caption.
func (b *B) CaptionE(attrs Attributes, body Body) {
	b.elementTyped("<caption", "</caption>", attrs, body)
}

// ColE writes a self-closing <col> element. It is the non-boxing fast path of Col.
func (b *B) ColE(attrs Attributes) { b.voidElementTyped("<col", attrs) }

// ColgroupE writes a <colgroup> element. It is the non-boxing fast path of Colgroup.
func (b *B) ColgroupE(attrs Attributes, body Body) {
	b.elementTyped("<colgroup", "</colgroup>", attrs, body)
}

// TableE writes a <table> element. It is the non-boxing fast path of Table.
func (b *B) TableE(attrs Attributes, body Body) { b.elementTyped("<table", "</table>", attrs, body) }

// TbodyE writes a <tbody> element. It is the non-boxing fast path of Tbody.
func (b *B) TbodyE(attrs Attributes, body Body) { b.elementTyped("<tbody", "</tbody>", attrs, body) }

// TdE writes a <td> element. It is the non-boxing fast path of Td.
func (b *B) TdE(attrs Attributes, body Body) { b.elementTyped("<td", "</td>", attrs, body) }

// TfootE writes a <tfoot> element. It is the non-boxing fast path of Tfoot.
func (b *B) TfootE(attrs Attributes, body Body) { b.elementTyped("<tfoot", "</tfoot>", attrs, body) }

// ThE writes a <th> element. It is the non-boxing fast path of Th.
func (b *B) ThE(attrs Attributes, body Body) { b.elementTyped("<th", "</th>", attrs, body) }

// TheadE writes a <thead> element. It is the non-boxing fast path of Thead.
func (b *B) TheadE(attrs Attributes, body Body) { b.elementTyped("<thead", "</thead>", attrs, body) }

// TrE writes a <tr> element. It is the non-boxing fast path of Tr.
func (b *B) TrE(attrs Attributes, body Body) { b.elementTyped("<tr", "</tr>", attrs, body) }

// ButtonE writes a <button> element. It is the non-boxing fast path of Button.
func (b *B) ButtonE(attrs Attributes, body Body) { b.elementTyped("<button", "</button>", attrs, body) }

// DatalistE writes a <datalist> element. It is the non-boxing fast path of Datalist.
func (b *B) DatalistE(attrs Attributes, body Body) {
	b.elementTyped("<datalist", "</datalist>", attrs, body)
}

// FieldsetE writes a <fieldset> element. It is the non-boxing fast path of Fieldset.
func (b *B) FieldsetE(attrs Attributes, body Body) {
	b.elementTyped("<fieldset", "</fieldset>", attrs, body)
}

// FormE writes a <form> element. It is the non-boxing fast path of Form.
func (b *B) FormE(attrs Attributes, body Body) { b.elementTyped("<form", "</form>", attrs, body) }

// InputE writes an self-closing <input> element. It is the non-boxing fast path of Input.
func (b *B) InputE(attrs Attributes) { b.voidElementTyped("<input", attrs) }

// LabelE writes a <label> element. It is the non-boxing fast path of Label.
func (b *B) LabelE(attrs Attributes, body Body) { b.elementTyped("<label", "</label>", attrs, body) }

// LegendE writes a <legend> element. It is the non-boxing fast path of Legend.
func (b *B) LegendE(attrs Attributes, body Body) { b.elementTyped("<legend", "</legend>", attrs, body) }

// MeterE writes a <meter> element. It is the non-boxing fast path of Meter.
func (b *B) MeterE(attrs Attributes, body Body) { b.elementTyped("<meter", "</meter>", attrs, body) }

// OptgroupE writes an <optgroup> element. It is the non-boxing fast path of Optgroup.
func (b *B) OptgroupE(attrs Attributes, body Body) {
	b.elementTyped("<optgroup", "</optgroup>", attrs, body)
}

// OptionE writes an <option> element. It is the non-boxing fast path of Option.
func (b *B) OptionE(attrs Attributes, body Body) { b.elementTyped("<option", "</option>", attrs, body) }

// OutputE writes an <output> element. It is the non-boxing fast path of Output.
func (b *B) OutputE(attrs Attributes, body Body) { b.elementTyped("<output", "</output>", attrs, body) }

// ProgressE writes a <progress> element. It is the non-boxing fast path of Progress.
func (b *B) ProgressE(attrs Attributes, body Body) {
	b.elementTyped("<progress", "</progress>", attrs, body)
}

// SelectE writes a <select> element. It is the non-boxing fast path of Select.
func (b *B) SelectE(attrs Attributes, body Body) { b.elementTyped("<select", "</select>", attrs, body) }

// TextareaE writes a <textarea> element. It is the non-boxing fast path of Textarea.
func (b *B) TextareaE(attrs Attributes, body Body) {
	b.elementTyped("<textarea", "</textarea>", attrs, body)
}

// DetailsE writes a <details> element. It is the non-boxing fast path of Details.
func (b *B) DetailsE(attrs Attributes, body Body) {
	b.elementTyped("<details", "</details>", attrs, body)
}

// DialogE writes a <dialog> element. It is the non-boxing fast path of Dialog.
func (b *B) DialogE(attrs Attributes, body Body) { b.elementTyped("<dialog", "</dialog>", attrs, body) }

// SummaryE writes a <summary> element. It is the non-boxing fast path of Summary.
func (b *B) SummaryE(attrs Attributes, body Body) {
	b.elementTyped("<summary", "</summary>", attrs, body)
}
