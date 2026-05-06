package js

// Pre-defined global identifiers.
var (
	// Console is the console object.
	Console = Ident("console")
	// Document is the document object.
	Document = Ident("document")
	// Window is the window object.
	Window = Ident("window")
	// Event is the event object in handlers (the event parameter).
	Event = Ident("event")
	// EventThis is the 'this' value in event handlers (the element).
	EventThis = This()
	// Location is the window.location object.
	Location = Ident("location")
	// History is the window.history object.
	History = Ident("history")
	// Navigator is the window.navigator object.
	Navigator = Ident("navigator")
	// LocalStorage is the localStorage object.
	LocalStorage = Ident("localStorage")
	// SessionStorage is the sessionStorage object.
	SessionStorage = Ident("sessionStorage")
	// JSON_ is the JSON global object (underscore to avoid conflict with JSON()).
	JSON_ = Ident("JSON")
	// Math is the Math global object.
	Math = Ident("Math")
	// Date is the Date constructor.
	Date = Ident("Date")
	// Promise is the Promise constructor.
	Promise = Ident("Promise")
	// Object_ is the Object constructor.
	Object_ = Ident("Object")
	// Array_ is the Array constructor.
	Array_ = Ident("Array")
)

// Console methods

// ConsoleLog creates console.log(args...).
func ConsoleLog(args ...Expr) Expr { return Console.Method("log", args...) }

// ConsoleError creates console.error(args...).
func ConsoleError(args ...Expr) Expr { return Console.Method("error", args...) }

// ConsoleWarn creates console.warn(args...).
func ConsoleWarn(args ...Expr) Expr { return Console.Method("warn", args...) }

// ConsoleInfo creates console.info(args...).
func ConsoleInfo(args ...Expr) Expr { return Console.Method("info", args...) }

// ConsoleDebug creates console.debug(args...).
func ConsoleDebug(args ...Expr) Expr { return Console.Method("debug", args...) }

// ConsoleTable creates console.table(data).
func ConsoleTable(data Expr) Expr { return Console.Method("table", data) }

// ConsoleClear creates console.clear().
func ConsoleClear() Expr { return Console.Method("clear") }

// Document methods

// GetElementById creates document.getElementById(id).
func GetElementById(id Expr) Expr { return Document.Method("getElementById", id) }

// QuerySelector creates document.querySelector(selector).
func QuerySelector(selector Expr) Expr { return Document.Method("querySelector", selector) }

// QuerySelectorAll creates document.querySelectorAll(selector).
func QuerySelectorAll(selector Expr) Expr { return Document.Method("querySelectorAll", selector) }

// CreateElement creates document.createElement(tag).
func CreateElement(tag Expr) Expr { return Document.Method("createElement", tag) }

// CreateTextNode creates document.createTextNode(text).
func CreateTextNode(text Expr) Expr { return Document.Method("createTextNode", text) }

// GetElementsByClassName creates document.getElementsByClassName(className).
func GetElementsByClassName(className Expr) Expr {
	return Document.Method("getElementsByClassName", className)
}

// GetElementsByTagName creates document.getElementsByTagName(tagName).
func GetElementsByTagName(tagName Expr) Expr {
	return Document.Method("getElementsByTagName", tagName)
}

// Window methods

// Alert creates alert(message).
func Alert(message Expr) Expr { return Ident("alert").Call(message) }

// Confirm creates confirm(message).
func Confirm(message Expr) Expr { return Ident("confirm").Call(message) }

// Prompt creates prompt(message, defaultValue?).
func Prompt(message Expr, defaultValue ...Expr) Expr {
	args := make([]Expr, 1, 1+len(defaultValue))
	args[0] = message
	args = append(args, defaultValue...)
	return Ident("prompt").Call(args...)
}

// SetTimeout creates setTimeout(callback, delay).
func SetTimeout(callback, delay Expr) Expr {
	return Ident("setTimeout").Call(callback, delay)
}

// SetInterval creates setInterval(callback, interval).
func SetInterval(callback, interval Expr) Expr {
	return Ident("setInterval").Call(callback, interval)
}

// ClearTimeout creates clearTimeout(id).
func ClearTimeout(id Expr) Expr { return Ident("clearTimeout").Call(id) }

// ClearInterval creates clearInterval(id).
func ClearInterval(id Expr) Expr { return Ident("clearInterval").Call(id) }

// RequestAnimationFrame creates requestAnimationFrame(callback).
func RequestAnimationFrame(callback Expr) Expr {
	return Ident("requestAnimationFrame").Call(callback)
}

// CancelAnimationFrame creates cancelAnimationFrame(id).
func CancelAnimationFrame(id Expr) Expr { return Ident("cancelAnimationFrame").Call(id) }

// Fetch creates fetch(url, options?).
func Fetch(url Expr, options ...Expr) Expr {
	args := make([]Expr, 1, 1+len(options))
	args[0] = url
	args = append(args, options...)
	return Ident("fetch").Call(args...)
}

// Event helpers

// PreventDefault creates event.preventDefault().
func PreventDefault() Expr { return Event.Method("preventDefault") }

// StopPropagation creates event.stopPropagation().
func StopPropagation() Expr { return Event.Method("stopPropagation") }

// StopImmediatePropagation creates event.stopImmediatePropagation().
func StopImmediatePropagation() Expr { return Event.Method("stopImmediatePropagation") }

// EventTarget creates event.target.
func EventTarget() Expr { return Event.Prop("target") }

// EventCurrentTarget creates event.currentTarget.
func EventCurrentTarget() Expr { return Event.Prop("currentTarget") }

// EventValue creates event.target.value.
func EventValue() Expr { return Event.Prop("target").Prop("value") }

// EventChecked creates event.target.checked.
func EventChecked() Expr { return Event.Prop("target").Prop("checked") }

// EventKey creates event.key.
func EventKey() Expr { return Event.Prop("key") }

// EventCode creates event.code.
func EventCode() Expr { return Event.Prop("code") }

// EventKeyCode creates event.keyCode (deprecated but common).
func EventKeyCode() Expr { return Event.Prop("keyCode") }

// EventWhich creates event.which (deprecated but common).
func EventWhich() Expr { return Event.Prop("which") }

// EventShiftKey creates event.shiftKey.
func EventShiftKey() Expr { return Event.Prop("shiftKey") }

// EventCtrlKey creates event.ctrlKey.
func EventCtrlKey() Expr { return Event.Prop("ctrlKey") }

// EventAltKey creates event.altKey.
func EventAltKey() Expr { return Event.Prop("altKey") }

// EventMetaKey creates event.metaKey.
func EventMetaKey() Expr { return Event.Prop("metaKey") }

// Navigation helpers

// Navigate creates location.href = url.
func Navigate(url Expr) Stmt { return Location.Prop("href").Assign(url) }

// Reload creates location.reload().
func Reload() Expr { return Location.Method("reload") }

// HistoryBack creates history.back().
func HistoryBack() Expr { return History.Method("back") }

// HistoryForward creates history.forward().
func HistoryForward() Expr { return History.Method("forward") }

// HistoryGo creates history.go(delta).
func HistoryGo(delta Expr) Expr { return History.Method("go", delta) }

// HistoryPushState creates history.pushState(state, title, url).
func HistoryPushState(state, title, url Expr) Expr {
	return History.Method("pushState", state, title, url)
}

// HistoryReplaceState creates history.replaceState(state, title, url).
func HistoryReplaceState(state, title, url Expr) Expr {
	return History.Method("replaceState", state, title, url)
}

// Storage helpers — methods on a storage Expr.

// GetItem creates storage.getItem(key).
func (e Expr) GetItem(key Expr) Expr { return e.Method("getItem", key) }

// SetItem creates storage.setItem(key, value).
func (e Expr) SetItem(key, value Expr) Expr { return e.Method("setItem", key, value) }

// RemoveItem creates storage.removeItem(key).
func (e Expr) RemoveItem(key Expr) Expr { return e.Method("removeItem", key) }

// ClearStorage creates storage.clear().
func (e Expr) ClearStorage() Expr { return e.Method("clear") }

// Element helpers — methods on an element Expr.

// Focus creates element.focus().
func (e Expr) Focus() Expr { return e.Method("focus") }

// Blur creates element.blur().
func (e Expr) Blur() Expr { return e.Method("blur") }

// Click creates element.click().
func (e Expr) Click() Expr { return e.Method("click") }

// SelectElement creates element.select().
func (e Expr) SelectElement() Expr { return e.Method("select") }

// AppendChild creates parent.appendChild(child).
func (e Expr) AppendChild(child Expr) Expr { return e.Method("appendChild", child) }

// RemoveChild creates parent.removeChild(child).
func (e Expr) RemoveChild(child Expr) Expr { return e.Method("removeChild", child) }

// InsertBefore creates parent.insertBefore(newNode, referenceNode).
func (e Expr) InsertBefore(newNode, referenceNode Expr) Expr {
	return e.Method("insertBefore", newNode, referenceNode)
}

// ReplaceChild creates parent.replaceChild(newChild, oldChild).
func (e Expr) ReplaceChild(newChild, oldChild Expr) Expr {
	return e.Method("replaceChild", newChild, oldChild)
}

// Remove creates element.remove().
func (e Expr) Remove() Expr { return e.Method("remove") }

// CloneNode creates element.cloneNode(deep).
func (e Expr) CloneNode(deep Expr) Expr { return e.Method("cloneNode", deep) }

// ClassList creates element.classList.
func (e Expr) ClassList() Expr { return e.Prop("classList") }

// ClassListAdd creates element.classList.add(classes...).
func (e Expr) ClassListAdd(classes ...Expr) Expr {
	return e.ClassList().Method("add", classes...)
}

// ClassListRemove creates element.classList.remove(classes...).
func (e Expr) ClassListRemove(classes ...Expr) Expr {
	return e.ClassList().Method("remove", classes...)
}

// ClassListToggle creates element.classList.toggle(className, force?).
func (e Expr) ClassListToggle(className Expr, force ...Expr) Expr {
	args := make([]Expr, 1, 1+len(force))
	args[0] = className
	args = append(args, force...)
	return e.ClassList().Method("toggle", args...)
}

// ClassListContains creates element.classList.contains(className).
func (e Expr) ClassListContains(className Expr) Expr {
	return e.ClassList().Method("contains", className)
}

// ClassListReplace creates element.classList.replace(oldClass, newClass).
func (e Expr) ClassListReplace(oldClass, newClass Expr) Expr {
	return e.ClassList().Method("replace", oldClass, newClass)
}

// GetAttribute creates element.getAttribute(name).
func (e Expr) GetAttribute(name Expr) Expr { return e.Method("getAttribute", name) }

// SetAttribute creates element.setAttribute(name, value).
func (e Expr) SetAttribute(name, value Expr) Expr {
	return e.Method("setAttribute", name, value)
}

// RemoveAttribute creates element.removeAttribute(name).
func (e Expr) RemoveAttribute(name Expr) Expr { return e.Method("removeAttribute", name) }

// HasAttribute creates element.hasAttribute(name).
func (e Expr) HasAttribute(name Expr) Expr { return e.Method("hasAttribute", name) }

// Style creates element.style.
func (e Expr) Style() Expr { return e.Prop("style") }

// SetStyle creates element.style.property = value.
func (e Expr) SetStyle(property string, value Expr) Stmt {
	return e.Style().Prop(property).Assign(value)
}

// JSON helpers

// JSONStringify creates JSON.stringify(value, replacer?, space?).
func JSONStringify(value Expr, args ...Expr) Expr {
	allArgs := make([]Expr, 1, 1+len(args))
	allArgs[0] = value
	allArgs = append(allArgs, args...)
	return JSON_.Method("stringify", allArgs...)
}

// JSONParse creates JSON.parse(text, reviver?).
func JSONParse(text Expr, reviver ...Expr) Expr {
	args := make([]Expr, 1, 1+len(reviver))
	args[0] = text
	args = append(args, reviver...)
	return JSON_.Method("parse", args...)
}

// Common patterns

// ParseInt creates parseInt(string, radix).
func ParseInt(str Expr, radix ...Expr) Expr {
	args := make([]Expr, 1, 1+len(radix))
	args[0] = str
	args = append(args, radix...)
	return Ident("parseInt").Call(args...)
}

// ParseFloat creates parseFloat(string).
func ParseFloat(str Expr) Expr { return Ident("parseFloat").Call(str) }

// IsNaN creates isNaN(value).
func IsNaN(value Expr) Expr { return Ident("isNaN").Call(value) }

// IsFinite creates isFinite(value).
func IsFinite(value Expr) Expr { return Ident("isFinite").Call(value) }

// EncodeURI creates encodeURI(uri).
func EncodeURI(uri Expr) Expr { return Ident("encodeURI").Call(uri) }

// DecodeURI creates decodeURI(uri).
func DecodeURI(uri Expr) Expr { return Ident("decodeURI").Call(uri) }

// EncodeURIComponent creates encodeURIComponent(component).
func EncodeURIComponent(component Expr) Expr {
	return Ident("encodeURIComponent").Call(component)
}

// DecodeURIComponent creates decodeURIComponent(component).
func DecodeURIComponent(component Expr) Expr {
	return Ident("decodeURIComponent").Call(component)
}
