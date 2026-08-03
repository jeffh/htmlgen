package h

// Html writes the root <html> element, preceded by the HTML5 doctype. It uses
// lang="en" unless a lang attribute is provided.
//
// Every document has exactly one <html> element wrapping a <head> and a <body>.
// The lang attribute tells screen readers which language to pronounce and
// browsers which font and hyphenation rules to apply, so set it when the page
// is not English.
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
//
// Holds metadata about the document rather than content the user sees: the
// <title>, <meta> tags, stylesheet <link>s, and <script> tags. Everything here
// is processed before the page renders.
func (b *B) Head(attrs Attributes, body Body) { b.element("<head", "</head>", attrs, body) }

// Title writes a <title> element.
//
// The document's name, shown in the browser tab, bookmarks, and search results.
// Required in <head>, and it is often the first thing a screen reader announces,
// so make it descriptive and unique per page.
func (b *B) Title(attrs Attributes, body Body) { b.element("<title", "</title>", attrs, body) }

// Meta writes a self-closing <meta> element.
//
// Carries document metadata that no other tag expresses: the character encoding
// (charset="utf-8"), the mobile viewport (name="viewport"), the search-result
// description (name="description"), and social preview cards (property="og:*").
func (b *B) Meta(attrs Attributes) { b.voidElement("<meta", attrs) }

// Link writes a self-closing <link> element.
//
// Relates the document to an external resource. Most commonly a stylesheet
// (rel="stylesheet" href="..."), but also favicons (rel="icon"), canonical URLs
// (rel="canonical"), preloads, and RSS feeds.
func (b *B) Link(attrs Attributes) { b.voidElement("<link", attrs) }

// Style writes a <style> element.
//
// Embeds CSS directly in the document, typically in <head>. Useful for
// critical above-the-fold styles or a single self-contained page; prefer a
// <link> stylesheet when the CSS is shared across pages and should be cached.
// The body is written verbatim, so use Raw and never interpolate untrusted
// input.
func (b *B) Style(attrs Attributes, body Body) { b.element("<style", "</style>", attrs, body) }

// Script writes a <script> element.
//
// Embeds JavaScript inline or loads it from a src URL. Add defer or async on
// external scripts so parsing is not blocked, and type="module" for ES modules.
// Also used with type="application/json" or type="application/ld+json" to embed
// data the page reads later. The body is written verbatim, so use Raw and never
// interpolate untrusted input.
func (b *B) Script(attrs Attributes, body Body) { b.element("<script", "</script>", attrs, body) }

// Noscript writes a <noscript> element.
//
// Content shown only when scripting is unavailable or disabled. Typically a
// short explanatory message or a non-JavaScript fallback for an interactive
// widget.
func (b *B) Noscript(attrs Attributes, body Body) { b.element("<noscript", "</noscript>", attrs, body) }

// Base writes a self-closing <base> element.
//
// Sets the base URL that every relative link and resource in the document
// resolves against, and/or a default target for links. At most one per
// document, and it must appear before any relative URL is used.
func (b *B) Base(attrs Attributes) { b.voidElement("<base", attrs) }

// Body writes a <body> element.
//
// Contains everything the user actually sees: text, images, and controls. One
// per document, immediately after <head>.
func (b *B) Body(attrs Attributes, body Body) { b.element("<body", "</body>", attrs, body) }

// Address writes an <address> element.
//
// Contact information for the nearest <article> or for the document as a whole
// — an author's email, a mailing address, a phone number. Not a generic wrapper
// for any postal address; use <p> for a shipping address in body text.
func (b *B) Address(attrs Attributes, body Body) { b.element("<address", "</address>", attrs, body) }

// Article writes an <article> element.
//
// A self-contained composition that would still make sense if syndicated or
// read alone: a blog post, a news story, a forum post, a product card, a
// comment. Articles can nest, as with comments inside a post.
func (b *B) Article(attrs Attributes, body Body) { b.element("<article", "</article>", attrs, body) }

// Aside writes an <aside> element.
//
// Content tangentially related to what surrounds it: a sidebar, a pull quote, a
// glossary definition, a block of related links or ads. Removing it should not
// break the main content.
func (b *B) Aside(attrs Attributes, body Body) { b.element("<aside", "</aside>", attrs, body) }

// Footer writes a <footer> element.
//
// Closing content for its nearest section or article — copyright, author
// details, related links, back-to-top navigation. A page may have several, one
// per <article> or <section>, in addition to the page footer.
func (b *B) Footer(attrs Attributes, body Body) { b.element("<footer", "</footer>", attrs, body) }

// Header writes a <header> element.
//
// Introductory content for its nearest section or article: a site banner with
// logo and nav, or an article's title, byline, and publication date. A page may
// have several.
func (b *B) Header(attrs Attributes, body Body) { b.element("<header", "</header>", attrs, body) }

// H1 writes an <h1> element.
//
// The top-level heading, naming what the page or article is about. Usually one
// per page (or one per <article>). Screen reader users navigate by heading, so
// keep H1 through H6 in order and do not skip levels for visual size — use CSS
// for that.
func (b *B) H1(attrs Attributes, body Body) { b.element("<h1", "</h1>", attrs, body) }

// H2 writes an <h2> element.
//
// A major section heading under the <h1>. The usual level for the top-level
// divisions of a page.
func (b *B) H2(attrs Attributes, body Body) { b.element("<h2", "</h2>", attrs, body) }

// H3 writes an <h3> element.
//
// A subsection heading under an <h2>.
func (b *B) H3(attrs Attributes, body Body) { b.element("<h3", "</h3>", attrs, body) }

// H4 writes an <h4> element.
//
// A fourth-level heading, nested under an <h3>.
func (b *B) H4(attrs Attributes, body Body) { b.element("<h4", "</h4>", attrs, body) }

// H5 writes an <h5> element.
//
// A fifth-level heading, nested under an <h4>. Rarely needed; deep nesting
// often signals the content should be split up.
func (b *B) H5(attrs Attributes, body Body) { b.element("<h5", "</h5>", attrs, body) }

// H6 writes an <h6> element.
//
// The deepest heading level, nested under an <h5>.
func (b *B) H6(attrs Attributes, body Body) { b.element("<h6", "</h6>", attrs, body) }

// Hgroup writes an <hgroup> element.
//
// Groups a single heading with adjacent <p> elements that act as a subtitle,
// tagline, or alternative title, so the subtitle does not become a heading of
// its own in the document outline.
func (b *B) Hgroup(attrs Attributes, body Body) { b.element("<hgroup", "</hgroup>", attrs, body) }

// Main writes a <main> element.
//
// The dominant content of the document — what is unique to this page, excluding
// the site header, nav, sidebar, and footer. One per page, and assistive
// technology offers a "skip to main content" jump to it.
func (b *B) Main(attrs Attributes, body Body) { b.element("<main", "</main>", attrs, body) }

// Nav writes a <nav> element.
//
// A block of major navigation links: the primary menu, a table of contents,
// breadcrumbs, pagination. Reserve it for significant link groups — a few links
// in a footer do not need it.
func (b *B) Nav(attrs Attributes, body Body) { b.element("<nav", "</nav>", attrs, body) }

// Section writes a <section> element.
//
// A thematic grouping of content, normally introduced by a heading: a chapter, a
// tabbed panel, a "Features" block. If the content is self-contained enough to
// stand alone use <article>; if it exists only for styling use <div>.
func (b *B) Section(attrs Attributes, body Body) { b.element("<section", "</section>", attrs, body) }

// Search writes a <search> element.
//
// Wraps the controls of a search or filtering interface — typically a <form>
// with a query <input> and a submit <button>. Communicates the search role to
// assistive technology without needing role="search".
func (b *B) Search(attrs Attributes, body Body) { b.element("<search", "</search>", attrs, body) }

// Blockquote writes a <blockquote> element.
//
// An extended quotation from another source, set off as its own block. Use the
// cite attribute for the source URL, and put a human-readable attribution in a
// following <figcaption> (inside a <figure>) rather than inside the quote.
func (b *B) Blockquote(attrs Attributes, body Body) {
	b.element("<blockquote", "</blockquote>", attrs, body)
}

// Dd writes a <dd> element.
//
// The description, definition, or value for the preceding <dt> term inside a
// <dl>. Several <dd> elements may follow one <dt> when a term has multiple
// meanings.
func (b *B) Dd(attrs Attributes, body Body) { b.element("<dd", "</dd>", attrs, body) }

// Div writes a <div> element.
//
// A generic block container with no meaning of its own — the fallback for
// grouping content purely for layout or styling. Reach for a semantic element
// (<section>, <article>, <nav>, <main>) first, and use <div> when none applies.
func (b *B) Div(attrs Attributes, body Body) { b.element("<div", "</div>", attrs, body) }

// Dl writes a <dl> element.
//
// A description list of term/description pairs: a glossary, metadata table,
// key/value summary, or FAQ. Contains <dt> terms each followed by one or more
// <dd> descriptions, optionally wrapped in <div>s.
func (b *B) Dl(attrs Attributes, body Body) { b.element("<dl", "</dl>", attrs, body) }

// Dt writes a <dt> element.
//
// The term, name, or key in a <dl>, described by the <dd> elements that follow
// it.
func (b *B) Dt(attrs Attributes, body Body) { b.element("<dt", "</dt>", attrs, body) }

// Figcaption writes a <figcaption> element.
//
// The caption or legend for its parent <figure>. Must be the figure's first or
// last child. Unlike an image's alt text, which replaces the image, a caption
// is read alongside it.
func (b *B) Figcaption(attrs Attributes, body Body) {
	b.element("<figcaption", "</figcaption>", attrs, body)
}

// Figure writes a <figure> element.
//
// Self-contained referenced content — an image, diagram, code listing, or table
// — that the surrounding text points to and that could move elsewhere on the
// page without changing its meaning. Pair it with a <figcaption>.
func (b *B) Figure(attrs Attributes, body Body) { b.element("<figure", "</figure>", attrs, body) }

// Hr writes a self-closing <hr> element.
//
// A thematic break between paragraph-level content, such as a scene change or a
// shift in topic. It is semantic, not decorative — for a plain dividing line,
// use a CSS border instead.
func (b *B) Hr(attrs Attributes) { b.voidElement("<hr", attrs) }

// Li writes an <li> element.
//
// One item in a <ul>, <ol>, or <menu>. In an ordered list, the value attribute
// overrides its number.
func (b *B) Li(attrs Attributes, body Body) { b.element("<li", "</li>", attrs, body) }

// Menu writes a <menu> element.
//
// A semantic alternative to <ul> for a list of interactive commands, such as a
// toolbar's buttons. Browsers render it identically to <ul>.
func (b *B) Menu(attrs Attributes, body Body) { b.element("<menu", "</menu>", attrs, body) }

// Ol writes an <ol> element.
//
// An ordered list whose sequence matters: numbered steps, rankings, a recipe's
// instructions. Use start to begin at another number, reversed to count down,
// and type to switch to letters or roman numerals.
func (b *B) Ol(attrs Attributes, body Body) { b.element("<ol", "</ol>", attrs, body) }

// P writes a <p> element.
//
// A paragraph — the default block for running prose. Browsers add vertical
// margins between paragraphs, so do not use empty <p> tags for spacing.
func (b *B) P(attrs Attributes, body Body) { b.element("<p", "</p>", attrs, body) }

// Pre writes a <pre> element.
//
// Preformatted text whose whitespace and line breaks are preserved and rendered
// in a monospace font: code blocks, ASCII art, terminal output. Wrap a <code>
// element inside it for source code.
func (b *B) Pre(attrs Attributes, body Body) { b.element("<pre", "</pre>", attrs, body) }

// Ul writes a <ul> element.
//
// An unordered list, for items whose order carries no meaning: feature bullets,
// navigation links, tags. Its direct children must be <li> elements.
func (b *B) Ul(attrs Attributes, body Body) { b.element("<ul", "</ul>", attrs, body) }

// A writes an <a> element.
//
// A hyperlink to another page, a fragment on this page (href="#id"), or another
// scheme such as mailto: or tel:. Use target="_blank" with
// rel="noopener noreferrer", and write link text that describes the destination
// rather than "click here". For an action that changes state rather than
// navigating, use <button>.
func (b *B) A(attrs Attributes, body Body) { b.element("<a", "</a>", attrs, body) }

// Abbr writes an <abbr> element.
//
// An abbreviation or acronym, with the expansion in its title attribute so it
// surfaces on hover and to assistive technology.
func (b *B) Abbr(attrs Attributes, body Body) { b.element("<abbr", "</abbr>", attrs, body) }

// B writes a <b> element.
//
// Draws attention to text without implying extra importance — a product name, a
// keyword in an abstract, the lead sentence of a review. When the text really
// is important or urgent, use <strong>; for purely visual weight, use CSS.
func (b *B) B(attrs Attributes, body Body) { b.element("<b", "</b>", attrs, body) }

// Bdi writes a <bdi> element.
//
// Isolates a span whose text direction is unknown — a user-supplied name, for
// instance — so right-to-left content cannot reorder the text around it.
func (b *B) Bdi(attrs Attributes, body Body) { b.element("<bdi", "</bdi>", attrs, body) }

// Bdo writes a <bdo> element.
//
// Overrides the bidirectional algorithm to force a text direction, via a
// required dir="ltr" or dir="rtl" attribute. Rare outside of typographic demos
// and mixed-script edge cases.
func (b *B) Bdo(attrs Attributes, body Body) { b.element("<bdo", "</bdo>", attrs, body) }

// Br writes a self-closing <br> element.
//
// A line break inside a block of text where the break is part of the content:
// poetry, song lyrics, a postal address. Do not use it to space out paragraphs
// — use separate <p> elements or CSS margins.
func (b *B) Br(attrs Attributes) { b.voidElement("<br", attrs) }

// Cite writes a <cite> element.
//
// The title of a referenced creative work — a book, film, paper, song, or blog
// post. It marks the work's title, not the person who made it.
func (b *B) Cite(attrs Attributes, body Body) { b.element("<cite", "</cite>", attrs, body) }

// Code writes a <code> element.
//
// A fragment of computer code shown inline: a function name, a filename, a CLI
// flag. Nest it inside <pre> for a multi-line block.
func (b *B) Code(attrs Attributes, body Body) { b.element("<code", "</code>", attrs, body) }

// Data writes a <data> element.
//
// Pairs human-readable content with a machine-readable equivalent in its value
// attribute — a product SKU behind a product name, or a sort key behind a
// label. For dates and times, use <time> instead.
func (b *B) Data(attrs Attributes, body Body) { b.element("<data", "</data>", attrs, body) }

// Dfn writes a <dfn> element.
//
// Marks the term being defined at the point where the surrounding text defines
// it. Usually appears once per term, on first use.
func (b *B) Dfn(attrs Attributes, body Body) { b.element("<dfn", "</dfn>", attrs, body) }

// Em writes an <em> element.
//
// Stress emphasis — the words a speaker would lean on, which can change a
// sentence's meaning. Rendered italic by default. For a title or a technical
// term set in italics for convention rather than emphasis, use <i>.
func (b *B) Em(attrs Attributes, body Body) { b.element("<em", "</em>", attrs, body) }

// I writes an <i> element.
//
// Text set apart from the prose by convention rather than by emphasis: a
// taxonomic name, a phrase in another language, a transliteration, a thought.
// Also the common hook for icon fonts. When you mean emphasis, use <em>.
func (b *B) I(attrs Attributes, body Body) { b.element("<i", "</i>", attrs, body) }

// Kbd writes a <kbd> element.
//
// User input entered from a keyboard, voice, or other device — a key name or a
// typed command. Nest one <kbd> per key to render a shortcut like Ctrl+C.
func (b *B) Kbd(attrs Attributes, body Body) { b.element("<kbd", "</kbd>", attrs, body) }

// Mark writes a <mark> element.
//
// Text highlighted for relevance in the current context: search-term matches in
// results, or the line a reader was pointed to. Rendered with a yellow
// background by default.
func (b *B) Mark(attrs Attributes, body Body) { b.element("<mark", "</mark>", attrs, body) }

// Q writes a <q> element.
//
// A short inline quotation. Browsers add the quotation marks themselves, so do
// not type them. Use <blockquote> for anything longer than a sentence or two.
func (b *B) Q(attrs Attributes, body Body) { b.element("<q", "</q>", attrs, body) }

// Rp writes an <rp> element.
//
// Fallback parentheses around ruby text, shown only by browsers that cannot
// render <ruby> annotations.
func (b *B) Rp(attrs Attributes, body Body) { b.element("<rp", "</rp>", attrs, body) }

// Rt writes an <rt> element.
//
// The annotation itself inside a <ruby> — the pronunciation guide rendered above
// or beside the base characters.
func (b *B) Rt(attrs Attributes, body Body) { b.element("<rt", "</rt>", attrs, body) }

// Ruby writes a <ruby> element.
//
// Wraps base text with small pronunciation annotations printed alongside it,
// used mainly for East Asian typography (furigana over kanji, pinyin over
// hanzi). Contains <rt> annotations and optional <rp> fallbacks.
func (b *B) Ruby(attrs Attributes, body Body) { b.element("<ruby", "</ruby>", attrs, body) }

// S writes an <s> element.
//
// Content that is no longer accurate or relevant, such as a sold-out item or a
// superseded price. Rendered with a strikethrough. For tracked edits in a
// document, use <del>.
func (b *B) S(attrs Attributes, body Body) { b.element("<s", "</s>", attrs, body) }

// Samp writes a <samp> element.
//
// Sample output from a program or system — an error message or console
// response. The counterpart to <kbd>, which marks what the user typed.
func (b *B) Samp(attrs Attributes, body Body) { b.element("<samp", "</samp>", attrs, body) }

// Small writes a <small> element.
//
// Side comments and fine print: copyright lines, legal disclaimers, attribution
// notices. It carries that meaning rather than merely shrinking text; use CSS
// font-size for presentation.
func (b *B) Small(attrs Attributes, body Body) { b.element("<small", "</small>", attrs, body) }

// Span writes a <span> element.
//
// A generic inline container with no meaning of its own, for styling or
// scripting a run of text. The inline counterpart to <div>, and the fallback
// when no semantic inline element fits.
func (b *B) Span(attrs Attributes, body Body) { b.element("<span", "</span>", attrs, body) }

// Strong writes a <strong> element.
//
// Content of strong importance, seriousness, or urgency: a warning, a caveat, a
// key term the reader must not miss. Rendered bold by default. For bold styling
// with no added importance, use <b> or CSS.
func (b *B) Strong(attrs Attributes, body Body) { b.element("<strong", "</strong>", attrs, body) }

// Sub writes a <sub> element.
//
// Subscript text, for content where the lowered position is meaningful:
// chemical formulas (H<sub>2</sub>O), mathematical indices, footnote markers.
func (b *B) Sub(attrs Attributes, body Body) { b.element("<sub", "</sub>", attrs, body) }

// Sup writes a <sup> element.
//
// Superscript text, for content where the raised position is meaningful:
// exponents, ordinal suffixes (4<sup>th</sup>), footnote references.
func (b *B) Sup(attrs Attributes, body Body) { b.element("<sup", "</sup>", attrs, body) }

// Time writes a <time> element.
//
// A date, time, or duration, with a machine-readable form in the datetime
// attribute (for example datetime="2024-03-01T09:00"). Lets browsers, search
// engines, and calendar tools parse a value written for humans as "last
// Tuesday".
func (b *B) Time(attrs Attributes, body Body) { b.element("<time", "</time>", attrs, body) }

// U writes a <u> element.
//
// A non-textual annotation rendered as an underline: a misspelling flagged by a
// spell checker, or a proper name marked in Chinese text. Avoid it for anything
// else, since underlined text reads as a link.
func (b *B) U(attrs Attributes, body Body) { b.element("<u", "</u>", attrs, body) }

// Var writes a <var> element.
//
// A variable name in a mathematical expression or programming context — the x
// in an equation, or a placeholder in a command's usage line.
func (b *B) Var(attrs Attributes, body Body) { b.element("<var", "</var>", attrs, body) }

// Wbr writes a self-closing <wbr> element.
//
// An optional word-break opportunity: the browser may wrap here if the line
// needs it, but will not add a hyphen. Useful inside long URLs and identifiers
// that would otherwise overflow.
func (b *B) Wbr(attrs Attributes) { b.voidElement("<wbr", attrs) }

// Area writes a self-closing <area> element.
//
// A clickable hotspot within an image <map>, defined by shape and coords. Give
// each one alt text, as you would an image link.
func (b *B) Area(attrs Attributes) { b.voidElement("<area", attrs) }

// Audio writes an <audio> element.
//
// Embeds sound: music, podcasts, sound effects. Point at a file with src, or
// nest <source> elements to offer several formats, plus <track> for captions.
// Add controls so the player is visible; autoplay with sound is blocked by most
// browsers.
func (b *B) Audio(attrs Attributes, body Body) { b.element("<audio", "</audio>", attrs, body) }

// Img writes a self-closing <img> element.
//
// Embeds an image. Always set alt — descriptive text for meaningful images, or
// alt="" for purely decorative ones. Setting width and height reserves space
// and prevents layout shift; loading="lazy" defers offscreen images.
func (b *B) Img(attrs Attributes) { b.voidElement("<img", attrs) }

// Map writes a <map> element.
//
// Defines an image map: a named set of <area> hotspots that an <img> references
// through its usemap attribute, making different regions of one image link to
// different destinations.
func (b *B) Map(attrs Attributes, body Body) { b.element("<map", "</map>", attrs, body) }

// Track writes a self-closing <track> element.
//
// Supplies a timed text track for <audio> or <video>: captions, subtitles,
// descriptions, or chapters, as a WebVTT file. Set kind, srclang, and label,
// and default on the track to enable initially.
func (b *B) Track(attrs Attributes) { b.voidElement("<track", attrs) }

// Video writes a <video> element.
//
// Embeds a video player. Point at a file with src, or nest <source> elements
// for multiple formats, plus <track> for captions. poster sets the placeholder
// frame; autoplay generally requires muted as well.
func (b *B) Video(attrs Attributes, body Body) { b.element("<video", "</video>", attrs, body) }

// Embed writes a self-closing <embed> element.
//
// Embeds external content handled by a plugin or the browser's own viewer, such
// as a PDF. Largely legacy — prefer <iframe>, <video>, <audio>, or <img>.
func (b *B) Embed(attrs Attributes) { b.voidElement("<embed", attrs) }

// Iframe writes an <iframe> element.
//
// Embeds another HTML document inline: a map, a video player, a payment form, a
// third-party widget. Use the sandbox and allow attributes to limit what the
// framed page can do, and title to name it for screen readers.
func (b *B) Iframe(attrs Attributes, body Body) { b.element("<iframe", "</iframe>", attrs, body) }

// Object writes an <object> element.
//
// Embeds an external resource — a PDF, an SVG, another document — with nested
// <param> elements or fallback content inside. Mostly superseded by the
// dedicated media elements and <iframe>.
func (b *B) Object(attrs Attributes, body Body) { b.element("<object", "</object>", attrs, body) }

// Picture writes a <picture> element.
//
// Wraps several <source> candidates and one required <img> fallback so the
// browser can pick an image by viewport width, pixel density, or format
// support. The tool for art direction and modern formats like AVIF and WebP.
func (b *B) Picture(attrs Attributes, body Body) { b.element("<picture", "</picture>", attrs, body) }

// Portal writes a <portal> element.
//
// An experimental element for embedding a preview of another page that can be
// activated into a full navigation. Not supported by current browsers; included
// for completeness.
func (b *B) Portal(attrs Attributes, body Body) { b.element("<portal", "</portal>", attrs, body) }

// Source writes a self-closing <source> element.
//
// One media or image candidate inside <picture>, <audio>, or <video>. The
// browser picks the first it can use, selecting on type, media, or srcset. List
// candidates from most to least preferred.
func (b *B) Source(attrs Attributes) { b.voidElement("<source", attrs) }

// Svg writes an <svg> element.
//
// Embeds inline vector graphics: icons, logos, charts, diagrams. Unlike an SVG
// loaded through <img>, inline SVG can be styled with CSS and scripted. Give
// meaningful graphics a <title> or aria-label, and mark decorative ones
// aria-hidden="true".
func (b *B) Svg(attrs Attributes, body Body) { b.element("<svg", "</svg>", attrs, body) }

// Math writes a <math> element.
//
// The root of MathML markup for rendering mathematical notation — fractions,
// radicals, matrices, integrals — as real text rather than an image.
func (b *B) Math(attrs Attributes, body Body) { b.element("<math", "</math>", attrs, body) }

// Canvas writes a <canvas> element.
//
// A scriptable bitmap surface drawn on with the 2D or WebGL context: games,
// data visualizations, image editing, generated graphics. Its contents are
// invisible to assistive technology, so put fallback content inside and provide
// an accessible alternative.
func (b *B) Canvas(attrs Attributes, body Body) { b.element("<canvas", "</canvas>", attrs, body) }

// Template writes a <template> element.
//
// Holds markup that is parsed but not rendered until JavaScript clones it —
// a row template for a table, or the shadow DOM contents of a custom element.
// Scripts and images inside do not run or load.
func (b *B) Template(attrs Attributes, body Body) { b.element("<template", "</template>", attrs, body) }

// Slot writes a <slot> element.
//
// A placeholder inside a web component's shadow DOM where the light-DOM
// children the caller supplied are projected. Use the name attribute for named
// slots.
func (b *B) Slot(attrs Attributes, body Body) { b.element("<slot", "</slot>", attrs, body) }

// Del writes a <del> element.
//
// Text removed in a revision of the document, as in a changelog or tracked
// edit. Pair it with <ins> for the replacement, and use cite and datetime to
// record why and when. For content that is merely outdated, use <s>.
func (b *B) Del(attrs Attributes, body Body) { b.element("<del", "</del>", attrs, body) }

// Ins writes an <ins> element.
//
// Text added in a revision of the document, the counterpart to <del>. Underlined
// by default, and accepts cite and datetime.
func (b *B) Ins(attrs Attributes, body Body) { b.element("<ins", "</ins>", attrs, body) }

// Caption writes a <caption> element.
//
// A title or description for its <table>, which it must be the first child of.
// It gives screen reader users the table's purpose before they read the cells.
func (b *B) Caption(attrs Attributes, body Body) { b.element("<caption", "</caption>", attrs, body) }

// Col writes a self-closing <col> element.
//
// Describes one column (or span columns) inside a <colgroup>, so a class or
// style can be applied to a whole column without touching every cell.
func (b *B) Col(attrs Attributes) { b.voidElement("<col", attrs) }

// Colgroup writes a <colgroup> element.
//
// Groups the column definitions of a table, containing <col> elements or using
// its own span attribute. Placed after <caption> and before any rows.
func (b *B) Colgroup(attrs Attributes, body Body) { b.element("<colgroup", "</colgroup>", attrs, body) }

// Table writes a <table> element.
//
// Displays genuinely tabular data — rows and columns whose position carries
// meaning. Not a layout tool; use CSS grid or flexbox for page layout. Structure
// it with <caption>, <thead>, <tbody>, and <tfoot> so it stays navigable.
func (b *B) Table(attrs Attributes, body Body) { b.element("<table", "</table>", attrs, body) }

// Tbody writes a <tbody> element.
//
// Groups the data rows of a table, distinct from its header and footer. A table
// may contain several to mark logical row groups, and browsers can scroll a
// tbody independently of a sticky <thead>.
func (b *B) Tbody(attrs Attributes, body Body) { b.element("<tbody", "</tbody>", attrs, body) }

// Td writes a <td> element.
//
// A data cell in a table row. Use colspan and rowspan to merge cells, and
// headers to point at the <th> cells that describe it in complex tables.
func (b *B) Td(attrs Attributes, body Body) { b.element("<td", "</td>", attrs, body) }

// Tfoot writes a <tfoot> element.
//
// Groups the summary rows of a table — totals, averages, footnotes — so they
// are identified as a footer regardless of where they appear in the source.
func (b *B) Tfoot(attrs Attributes, body Body) { b.element("<tfoot", "</tfoot>", attrs, body) }

// Th writes a <th> element.
//
// A header cell labeling a row or column. Set scope="col" or scope="row" so
// screen readers can announce the right header when reading each data cell.
func (b *B) Th(attrs Attributes, body Body) { b.element("<th", "</th>", attrs, body) }

// Thead writes a <thead> element.
//
// Groups the header rows of a table. Browsers repeat it across printed pages,
// and it is the usual target for a sticky header.
func (b *B) Thead(attrs Attributes, body Body) { b.element("<thead", "</thead>", attrs, body) }

// Tr writes a <tr> element.
//
// A single row of a table, containing <th> and/or <td> cells.
func (b *B) Tr(attrs Attributes, body Body) { b.element("<tr", "</tr>", attrs, body) }

// Button writes a <button> element.
//
// A clickable control that performs an action: submitting a form
// (type="submit", the default inside a form), resetting it, or running
// JavaScript (type="button"). Use it whenever the click does something on the
// page; use <a> when it navigates somewhere.
func (b *B) Button(attrs Attributes, body Body) { b.element("<button", "</button>", attrs, body) }

// Datalist writes a <datalist> element.
//
// A list of suggested <option> values for an <input> that references it by
// list="id". Unlike <select>, the suggestions are advisory — the user may still
// type anything.
func (b *B) Datalist(attrs Attributes, body Body) { b.element("<datalist", "</datalist>", attrs, body) }

// Fieldset writes a <fieldset> element.
//
// Groups related form controls under a shared <legend>, which is what makes a
// set of radio buttons or checkboxes comprehensible to screen readers. Its
// disabled attribute disables every control inside at once.
func (b *B) Fieldset(attrs Attributes, body Body) { b.element("<fieldset", "</fieldset>", attrs, body) }

// Form writes a <form> element.
//
// Collects user input and submits it to action using method="get" or
// method="post". Use enctype="multipart/form-data" for file uploads. Forms
// cannot be nested.
func (b *B) Form(attrs Attributes, body Body) { b.element("<form", "</form>", attrs, body) }

// Input writes a self-closing <input> element.
//
// The general-purpose form control; its type attribute selects the behavior —
// text, email, password, number, date, checkbox, radio, file, range, color,
// hidden, and more. Choosing the right type gives mobile users the right
// keyboard and enables built-in validation. Pair each one with a <label>.
func (b *B) Input(attrs Attributes) { b.voidElement("<input", attrs) }

// Label writes a <label> element.
//
// Names a form control, either by wrapping it or by pointing at its id with
// for="...". Beyond accessibility, clicking a label focuses or toggles its
// control, which makes checkboxes and radios far easier to hit.
func (b *B) Label(attrs Attributes, body Body) { b.element("<label", "</label>", attrs, body) }

// Legend writes a <legend> element.
//
// The caption for its <fieldset>, and must be the fieldset's first child.
// Screen readers announce it with each control in the group, so it is what tells
// the user what a radio group is asking.
func (b *B) Legend(attrs Attributes, body Body) { b.element("<legend", "</legend>", attrs, body) }

// Meter writes a <meter> element.
//
// A scalar measurement within a known range, where the range itself is
// meaningful: disk usage, a score, a relevance rating. Set value with min, max,
// and optionally low, high, and optimum so browsers can color it. For a task
// that is running, use <progress>.
func (b *B) Meter(attrs Attributes, body Body) { b.element("<meter", "</meter>", attrs, body) }

// Optgroup writes an <optgroup> element.
//
// Groups related <option> elements inside a <select> under a label, as with
// cities grouped by country. Groups cannot nest.
func (b *B) Optgroup(attrs Attributes, body Body) { b.element("<optgroup", "</optgroup>", attrs, body) }

// Option writes an <option> element.
//
// One choice in a <select>, <optgroup>, or <datalist>. The value attribute is
// what gets submitted when it differs from the visible text; selected marks the
// initial choice.
func (b *B) Option(attrs Attributes, body Body) { b.element("<option", "</option>", attrs, body) }

// Output writes an <output> element.
//
// Displays the result of a calculation or user action — a computed total, a
// slider's current value, a validation message. It is a live region, so screen
// readers announce updates without extra ARIA.
func (b *B) Output(attrs Attributes, body Body) { b.element("<output", "</output>", attrs, body) }

// Progress writes a <progress> element.
//
// Shows how far along a task is: a file upload, a multi-step form, a long
// computation. Set value and max for a determinate bar; omit value for an
// indeterminate spinner. For a static measurement, use <meter>.
func (b *B) Progress(attrs Attributes, body Body) { b.element("<progress", "</progress>", attrs, body) }

// Select writes a <select> element.
//
// A drop-down list of <option> choices, optionally grouped by <optgroup>. Add
// multiple for a multi-select list box, and size to show several rows at once.
func (b *B) Select(attrs Attributes, body Body) { b.element("<select", "</select>", attrs, body) }

// Textarea writes a <textarea> element.
//
// A multi-line free-text control: comments, messages, descriptions. Its initial
// value is its body content, not a value attribute; rows and cols size it, and
// maxlength caps the input.
func (b *B) Textarea(attrs Attributes, body Body) { b.element("<textarea", "</textarea>", attrs, body) }

// Details writes a <details> element.
//
// A disclosure widget the user expands and collapses — an FAQ answer, an
// advanced-options panel, a stack trace. The first child <summary> is the
// always-visible label; the open attribute starts it expanded. Works without
// JavaScript.
func (b *B) Details(attrs Attributes, body Body) { b.element("<details", "</details>", attrs, body) }

// Dialog writes a <dialog> element.
//
// A modal or non-modal dialog box: a confirmation prompt, a settings panel, a
// lightbox. showModal() opens it with a backdrop and a focus trap, so it
// handles Escape and focus management natively.
func (b *B) Dialog(attrs Attributes, body Body) { b.element("<dialog", "</dialog>", attrs, body) }

// Summary writes a <summary> element.
//
// The visible heading of a <details> widget, which it must be the first child
// of. Clicking it toggles the disclosure.
func (b *B) Summary(attrs Attributes, body Body) { b.element("<summary", "</summary>", attrs, body) }
