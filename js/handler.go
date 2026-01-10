package js

import (
	"strings"

	"github.com/jeffh/htmlgen/h"
)

// Handler builds an inline JavaScript handler string from statements.
// Statements are joined with semicolons.
func Handler(stmts ...Stmt) string {
	if len(stmts) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.Grow(64) // Pre-allocate for typical handler size
	for i, stmt := range stmts {
		if i > 0 {
			sb.WriteString("; ")
		}
		stmt.stmt(&sb)
	}
	return sb.String()
}

// ExprHandler builds an inline JavaScript handler from a single expression.
func ExprHandler(expr Expr) string {
	var sb strings.Builder
	sb.Grow(64)
	expr.js(&sb)
	return sb.String()
}

// ToJS converts an expression to its JavaScript string representation.
func ToJS(expr Expr) string {
	var sb strings.Builder
	expr.js(&sb)
	return sb.String()
}

// ToJSStmt converts a statement to its JavaScript string representation.
func ToJSStmt(stmt Stmt) string {
	var sb strings.Builder
	stmt.stmt(&sb)
	return sb.String()
}

// Mouse Events

// OnClick creates an onclick attribute with the given handler.
func OnClick(stmts ...Stmt) h.Attribute {
	return h.Attr("onclick", Handler(stmts...))
}

// OnDblClick creates an ondblclick attribute with the given handler.
func OnDblClick(stmts ...Stmt) h.Attribute {
	return h.Attr("ondblclick", Handler(stmts...))
}

// OnMouseDown creates an onmousedown attribute with the given handler.
func OnMouseDown(stmts ...Stmt) h.Attribute {
	return h.Attr("onmousedown", Handler(stmts...))
}

// OnMouseUp creates an onmouseup attribute with the given handler.
func OnMouseUp(stmts ...Stmt) h.Attribute {
	return h.Attr("onmouseup", Handler(stmts...))
}

// OnMouseOver creates an onmouseover attribute with the given handler.
func OnMouseOver(stmts ...Stmt) h.Attribute {
	return h.Attr("onmouseover", Handler(stmts...))
}

// OnMouseOut creates an onmouseout attribute with the given handler.
func OnMouseOut(stmts ...Stmt) h.Attribute {
	return h.Attr("onmouseout", Handler(stmts...))
}

// OnMouseEnter creates an onmouseenter attribute with the given handler.
func OnMouseEnter(stmts ...Stmt) h.Attribute {
	return h.Attr("onmouseenter", Handler(stmts...))
}

// OnMouseLeave creates an onmouseleave attribute with the given handler.
func OnMouseLeave(stmts ...Stmt) h.Attribute {
	return h.Attr("onmouseleave", Handler(stmts...))
}

// OnMouseMove creates an onmousemove attribute with the given handler.
func OnMouseMove(stmts ...Stmt) h.Attribute {
	return h.Attr("onmousemove", Handler(stmts...))
}

// OnContextMenu creates an oncontextmenu attribute with the given handler.
func OnContextMenu(stmts ...Stmt) h.Attribute {
	return h.Attr("oncontextmenu", Handler(stmts...))
}

// OnWheel creates an onwheel attribute with the given handler.
func OnWheel(stmts ...Stmt) h.Attribute {
	return h.Attr("onwheel", Handler(stmts...))
}

// Keyboard Events

// OnKeyDown creates an onkeydown attribute with the given handler.
func OnKeyDown(stmts ...Stmt) h.Attribute {
	return h.Attr("onkeydown", Handler(stmts...))
}

// OnKeyUp creates an onkeyup attribute with the given handler.
func OnKeyUp(stmts ...Stmt) h.Attribute {
	return h.Attr("onkeyup", Handler(stmts...))
}

// OnKeyPress creates an onkeypress attribute with the given handler.
// Deprecated: Use OnKeyDown or OnKeyUp instead.
func OnKeyPress(stmts ...Stmt) h.Attribute {
	return h.Attr("onkeypress", Handler(stmts...))
}

// Focus Events

// OnFocus creates an onfocus attribute with the given handler.
func OnFocus(stmts ...Stmt) h.Attribute {
	return h.Attr("onfocus", Handler(stmts...))
}

// OnBlur creates an onblur attribute with the given handler.
func OnBlur(stmts ...Stmt) h.Attribute {
	return h.Attr("onblur", Handler(stmts...))
}

// OnFocusIn creates an onfocusin attribute with the given handler.
func OnFocusIn(stmts ...Stmt) h.Attribute {
	return h.Attr("onfocusin", Handler(stmts...))
}

// OnFocusOut creates an onfocusout attribute with the given handler.
func OnFocusOut(stmts ...Stmt) h.Attribute {
	return h.Attr("onfocusout", Handler(stmts...))
}

// Form Events

// OnChange creates an onchange attribute with the given handler.
func OnChange(stmts ...Stmt) h.Attribute {
	return h.Attr("onchange", Handler(stmts...))
}

// OnInput creates an oninput attribute with the given handler.
func OnInput(stmts ...Stmt) h.Attribute {
	return h.Attr("oninput", Handler(stmts...))
}

// OnSubmit creates an onsubmit attribute with the given handler.
func OnSubmit(stmts ...Stmt) h.Attribute {
	return h.Attr("onsubmit", Handler(stmts...))
}

// OnReset creates an onreset attribute with the given handler.
func OnReset(stmts ...Stmt) h.Attribute {
	return h.Attr("onreset", Handler(stmts...))
}

// OnSelect creates an onselect attribute with the given handler.
func OnSelect(stmts ...Stmt) h.Attribute {
	return h.Attr("onselect", Handler(stmts...))
}

// OnInvalid creates an oninvalid attribute with the given handler.
func OnInvalid(stmts ...Stmt) h.Attribute {
	return h.Attr("oninvalid", Handler(stmts...))
}

// Window/Document Events

// OnLoad creates an onload attribute with the given handler.
func OnLoad(stmts ...Stmt) h.Attribute {
	return h.Attr("onload", Handler(stmts...))
}

// OnUnload creates an onunload attribute with the given handler.
func OnUnload(stmts ...Stmt) h.Attribute {
	return h.Attr("onunload", Handler(stmts...))
}

// OnBeforeUnload creates an onbeforeunload attribute with the given handler.
func OnBeforeUnload(stmts ...Stmt) h.Attribute {
	return h.Attr("onbeforeunload", Handler(stmts...))
}

// OnError creates an onerror attribute with the given handler.
func OnError(stmts ...Stmt) h.Attribute {
	return h.Attr("onerror", Handler(stmts...))
}

// OnScroll creates an onscroll attribute with the given handler.
func OnScroll(stmts ...Stmt) h.Attribute {
	return h.Attr("onscroll", Handler(stmts...))
}

// OnResize creates an onresize attribute with the given handler.
func OnResize(stmts ...Stmt) h.Attribute {
	return h.Attr("onresize", Handler(stmts...))
}

// OnHashChange creates an onhashchange attribute with the given handler.
func OnHashChange(stmts ...Stmt) h.Attribute {
	return h.Attr("onhashchange", Handler(stmts...))
}

// OnPopState creates an onpopstate attribute with the given handler.
func OnPopState(stmts ...Stmt) h.Attribute {
	return h.Attr("onpopstate", Handler(stmts...))
}

// OnStorage creates an onstorage attribute with the given handler.
func OnStorage(stmts ...Stmt) h.Attribute {
	return h.Attr("onstorage", Handler(stmts...))
}

// OnOnline creates an ononline attribute with the given handler.
func OnOnline(stmts ...Stmt) h.Attribute {
	return h.Attr("ononline", Handler(stmts...))
}

// OnOffline creates an onoffline attribute with the given handler.
func OnOffline(stmts ...Stmt) h.Attribute {
	return h.Attr("onoffline", Handler(stmts...))
}

// Clipboard Events

// OnCopy creates an oncopy attribute with the given handler.
func OnCopy(stmts ...Stmt) h.Attribute {
	return h.Attr("oncopy", Handler(stmts...))
}

// OnCut creates an oncut attribute with the given handler.
func OnCut(stmts ...Stmt) h.Attribute {
	return h.Attr("oncut", Handler(stmts...))
}

// OnPaste creates an onpaste attribute with the given handler.
func OnPaste(stmts ...Stmt) h.Attribute {
	return h.Attr("onpaste", Handler(stmts...))
}

// Drag Events

// OnDrag creates an ondrag attribute with the given handler.
func OnDrag(stmts ...Stmt) h.Attribute {
	return h.Attr("ondrag", Handler(stmts...))
}

// OnDragStart creates an ondragstart attribute with the given handler.
func OnDragStart(stmts ...Stmt) h.Attribute {
	return h.Attr("ondragstart", Handler(stmts...))
}

// OnDragEnd creates an ondragend attribute with the given handler.
func OnDragEnd(stmts ...Stmt) h.Attribute {
	return h.Attr("ondragend", Handler(stmts...))
}

// OnDragOver creates an ondragover attribute with the given handler.
func OnDragOver(stmts ...Stmt) h.Attribute {
	return h.Attr("ondragover", Handler(stmts...))
}

// OnDragEnter creates an ondragenter attribute with the given handler.
func OnDragEnter(stmts ...Stmt) h.Attribute {
	return h.Attr("ondragenter", Handler(stmts...))
}

// OnDragLeave creates an ondragleave attribute with the given handler.
func OnDragLeave(stmts ...Stmt) h.Attribute {
	return h.Attr("ondragleave", Handler(stmts...))
}

// OnDrop creates an ondrop attribute with the given handler.
func OnDrop(stmts ...Stmt) h.Attribute {
	return h.Attr("ondrop", Handler(stmts...))
}

// Touch Events

// OnTouchStart creates an ontouchstart attribute with the given handler.
func OnTouchStart(stmts ...Stmt) h.Attribute {
	return h.Attr("ontouchstart", Handler(stmts...))
}

// OnTouchMove creates an ontouchmove attribute with the given handler.
func OnTouchMove(stmts ...Stmt) h.Attribute {
	return h.Attr("ontouchmove", Handler(stmts...))
}

// OnTouchEnd creates an ontouchend attribute with the given handler.
func OnTouchEnd(stmts ...Stmt) h.Attribute {
	return h.Attr("ontouchend", Handler(stmts...))
}

// OnTouchCancel creates an ontouchcancel attribute with the given handler.
func OnTouchCancel(stmts ...Stmt) h.Attribute {
	return h.Attr("ontouchcancel", Handler(stmts...))
}

// Pointer Events

// OnPointerDown creates an onpointerdown attribute with the given handler.
func OnPointerDown(stmts ...Stmt) h.Attribute {
	return h.Attr("onpointerdown", Handler(stmts...))
}

// OnPointerUp creates an onpointerup attribute with the given handler.
func OnPointerUp(stmts ...Stmt) h.Attribute {
	return h.Attr("onpointerup", Handler(stmts...))
}

// OnPointerMove creates an onpointermove attribute with the given handler.
func OnPointerMove(stmts ...Stmt) h.Attribute {
	return h.Attr("onpointermove", Handler(stmts...))
}

// OnPointerEnter creates an onpointerenter attribute with the given handler.
func OnPointerEnter(stmts ...Stmt) h.Attribute {
	return h.Attr("onpointerenter", Handler(stmts...))
}

// OnPointerLeave creates an onpointerleave attribute with the given handler.
func OnPointerLeave(stmts ...Stmt) h.Attribute {
	return h.Attr("onpointerleave", Handler(stmts...))
}

// OnPointerCancel creates an onpointercancel attribute with the given handler.
func OnPointerCancel(stmts ...Stmt) h.Attribute {
	return h.Attr("onpointercancel", Handler(stmts...))
}

// OnPointerOver creates an onpointerover attribute with the given handler.
func OnPointerOver(stmts ...Stmt) h.Attribute {
	return h.Attr("onpointerover", Handler(stmts...))
}

// OnPointerOut creates an onpointerout attribute with the given handler.
func OnPointerOut(stmts ...Stmt) h.Attribute {
	return h.Attr("onpointerout", Handler(stmts...))
}

// Media Events

// OnPlay creates an onplay attribute with the given handler.
func OnPlay(stmts ...Stmt) h.Attribute {
	return h.Attr("onplay", Handler(stmts...))
}

// OnPause creates an onpause attribute with the given handler.
func OnPause(stmts ...Stmt) h.Attribute {
	return h.Attr("onpause", Handler(stmts...))
}

// OnEnded creates an onended attribute with the given handler.
func OnEnded(stmts ...Stmt) h.Attribute {
	return h.Attr("onended", Handler(stmts...))
}

// OnTimeUpdate creates an ontimeupdate attribute with the given handler.
func OnTimeUpdate(stmts ...Stmt) h.Attribute {
	return h.Attr("ontimeupdate", Handler(stmts...))
}

// OnVolumeChange creates an onvolumechange attribute with the given handler.
func OnVolumeChange(stmts ...Stmt) h.Attribute {
	return h.Attr("onvolumechange", Handler(stmts...))
}

// OnSeeking creates an onseeking attribute with the given handler.
func OnSeeking(stmts ...Stmt) h.Attribute {
	return h.Attr("onseeking", Handler(stmts...))
}

// OnSeeked creates an onseeked attribute with the given handler.
func OnSeeked(stmts ...Stmt) h.Attribute {
	return h.Attr("onseeked", Handler(stmts...))
}

// OnLoadedData creates an onloadeddata attribute with the given handler.
func OnLoadedData(stmts ...Stmt) h.Attribute {
	return h.Attr("onloadeddata", Handler(stmts...))
}

// OnLoadedMetadata creates an onloadedmetadata attribute with the given handler.
func OnLoadedMetadata(stmts ...Stmt) h.Attribute {
	return h.Attr("onloadedmetadata", Handler(stmts...))
}

// OnCanPlay creates an oncanplay attribute with the given handler.
func OnCanPlay(stmts ...Stmt) h.Attribute {
	return h.Attr("oncanplay", Handler(stmts...))
}

// OnCanPlayThrough creates an oncanplaythrough attribute with the given handler.
func OnCanPlayThrough(stmts ...Stmt) h.Attribute {
	return h.Attr("oncanplaythrough", Handler(stmts...))
}

// Animation Events

// OnAnimationStart creates an onanimationstart attribute with the given handler.
func OnAnimationStart(stmts ...Stmt) h.Attribute {
	return h.Attr("onanimationstart", Handler(stmts...))
}

// OnAnimationEnd creates an onanimationend attribute with the given handler.
func OnAnimationEnd(stmts ...Stmt) h.Attribute {
	return h.Attr("onanimationend", Handler(stmts...))
}

// OnAnimationIteration creates an onanimationiteration attribute with the given handler.
func OnAnimationIteration(stmts ...Stmt) h.Attribute {
	return h.Attr("onanimationiteration", Handler(stmts...))
}

// Transition Events

// OnTransitionEnd creates an ontransitionend attribute with the given handler.
func OnTransitionEnd(stmts ...Stmt) h.Attribute {
	return h.Attr("ontransitionend", Handler(stmts...))
}

// OnTransitionStart creates an ontransitionstart attribute with the given handler.
func OnTransitionStart(stmts ...Stmt) h.Attribute {
	return h.Attr("ontransitionstart", Handler(stmts...))
}

// OnTransitionRun creates an ontransitionrun attribute with the given handler.
func OnTransitionRun(stmts ...Stmt) h.Attribute {
	return h.Attr("ontransitionrun", Handler(stmts...))
}

// OnTransitionCancel creates an ontransitioncancel attribute with the given handler.
func OnTransitionCancel(stmts ...Stmt) h.Attribute {
	return h.Attr("ontransitioncancel", Handler(stmts...))
}

// Print Events

// OnBeforePrint creates an onbeforeprint attribute with the given handler.
func OnBeforePrint(stmts ...Stmt) h.Attribute {
	return h.Attr("onbeforeprint", Handler(stmts...))
}

// OnAfterPrint creates an onafterprint attribute with the given handler.
func OnAfterPrint(stmts ...Stmt) h.Attribute {
	return h.Attr("onafterprint", Handler(stmts...))
}

// Custom Event

// On creates a custom event handler attribute.
// Example: On("touchstart", stmts...) creates ontouchstart="..."
func On(event string, stmts ...Stmt) h.Attribute {
	return h.Attr("on"+event, Handler(stmts...))
}
