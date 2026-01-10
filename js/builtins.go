package js

// Pre-defined global identifiers
var (
	// Console is the console object
	Console = Ident("console")
	// Document is the document object
	Document = Ident("document")
	// Window is the window object
	Window = Ident("window")
	// Event is the event object in handlers (the event parameter)
	Event = Ident("event")
	// EventThis is the 'this' value in event handlers (the element)
	EventThis = This()
	// Location is the window.location object
	Location = Ident("location")
	// History is the window.history object
	History = Ident("history")
	// Navigator is the window.navigator object
	Navigator = Ident("navigator")
	// LocalStorage is the localStorage object
	LocalStorage = Ident("localStorage")
	// SessionStorage is the sessionStorage object
	SessionStorage = Ident("sessionStorage")
	// JSON_ is the JSON global object (underscore to avoid conflict with JSON() function)
	JSON_ = Ident("JSON")
	// Math is the Math global object
	Math = Ident("Math")
	// Date is the Date constructor
	Date = Ident("Date")
	// Promise is the Promise constructor
	Promise = Ident("Promise")
	// Object_ is the Object constructor
	Object_ = Ident("Object")
	// Array_ is the Array constructor
	Array_ = Ident("Array")
)

// Console methods

// ConsoleLog creates console.log(args...)
func ConsoleLog(args ...Expr) Callable {
	return Method(Console, "log", args...)
}

// ConsoleError creates console.error(args...)
func ConsoleError(args ...Expr) Callable {
	return Method(Console, "error", args...)
}

// ConsoleWarn creates console.warn(args...)
func ConsoleWarn(args ...Expr) Callable {
	return Method(Console, "warn", args...)
}

// ConsoleInfo creates console.info(args...)
func ConsoleInfo(args ...Expr) Callable {
	return Method(Console, "info", args...)
}

// ConsoleDebug creates console.debug(args...)
func ConsoleDebug(args ...Expr) Callable {
	return Method(Console, "debug", args...)
}

// ConsoleTable creates console.table(data)
func ConsoleTable(data Expr) Callable {
	return Method(Console, "table", data)
}

// ConsoleClear creates console.clear()
func ConsoleClear() Callable {
	return Method(Console, "clear")
}

// Document methods

// GetElementById creates document.getElementById(id)
func GetElementById(id Expr) Callable {
	return Method(Document, "getElementById", id)
}

// QuerySelector creates document.querySelector(selector)
func QuerySelector(selector Expr) Callable {
	return Method(Document, "querySelector", selector)
}

// QuerySelectorAll creates document.querySelectorAll(selector)
func QuerySelectorAll(selector Expr) Callable {
	return Method(Document, "querySelectorAll", selector)
}

// CreateElement creates document.createElement(tag)
func CreateElement(tag Expr) Callable {
	return Method(Document, "createElement", tag)
}

// CreateTextNode creates document.createTextNode(text)
func CreateTextNode(text Expr) Callable {
	return Method(Document, "createTextNode", text)
}

// GetElementsByClassName creates document.getElementsByClassName(className)
func GetElementsByClassName(className Expr) Callable {
	return Method(Document, "getElementsByClassName", className)
}

// GetElementsByTagName creates document.getElementsByTagName(tagName)
func GetElementsByTagName(tagName Expr) Callable {
	return Method(Document, "getElementsByTagName", tagName)
}

// Window methods

// Alert creates alert(message)
func Alert(message Expr) Callable {
	return Call(Ident("alert"), message)
}

// Confirm creates confirm(message)
func Confirm(message Expr) Callable {
	return Call(Ident("confirm"), message)
}

// Prompt creates prompt(message, defaultValue)
func Prompt(message Expr, defaultValue ...Expr) Callable {
	args := make([]Expr, 1, 1+len(defaultValue))
	args[0] = message
	args = append(args, defaultValue...)
	return Call(Ident("prompt"), args...)
}

// SetTimeout creates setTimeout(callback, delay)
func SetTimeout(callback, delay Expr) Callable {
	return Call(Ident("setTimeout"), callback, delay)
}

// SetInterval creates setInterval(callback, interval)
func SetInterval(callback, interval Expr) Callable {
	return Call(Ident("setInterval"), callback, interval)
}

// ClearTimeout creates clearTimeout(id)
func ClearTimeout(id Expr) Callable {
	return Call(Ident("clearTimeout"), id)
}

// ClearInterval creates clearInterval(id)
func ClearInterval(id Expr) Callable {
	return Call(Ident("clearInterval"), id)
}

// RequestAnimationFrame creates requestAnimationFrame(callback)
func RequestAnimationFrame(callback Expr) Callable {
	return Call(Ident("requestAnimationFrame"), callback)
}

// CancelAnimationFrame creates cancelAnimationFrame(id)
func CancelAnimationFrame(id Expr) Callable {
	return Call(Ident("cancelAnimationFrame"), id)
}

// Fetch creates fetch(url, options)
func Fetch(url Expr, options ...Expr) Callable {
	args := make([]Expr, 1, 1+len(options))
	args[0] = url
	args = append(args, options...)
	return Call(Ident("fetch"), args...)
}

// Event helpers

// PreventDefault creates event.preventDefault()
func PreventDefault() Callable {
	return Method(Event, "preventDefault")
}

// StopPropagation creates event.stopPropagation()
func StopPropagation() Callable {
	return Method(Event, "stopPropagation")
}

// StopImmediatePropagation creates event.stopImmediatePropagation()
func StopImmediatePropagation() Callable {
	return Method(Event, "stopImmediatePropagation")
}

// EventTarget creates event.target
func EventTarget() Callable {
	return Prop(Event, "target")
}

// EventCurrentTarget creates event.currentTarget
func EventCurrentTarget() Callable {
	return Prop(Event, "currentTarget")
}

// EventValue creates event.target.value (common for input handlers)
func EventValue() Callable {
	return Prop(Prop(Event, "target"), "value")
}

// EventChecked creates event.target.checked (common for checkbox handlers)
func EventChecked() Callable {
	return Prop(Prop(Event, "target"), "checked")
}

// EventKey creates event.key (for keyboard events)
func EventKey() Callable {
	return Prop(Event, "key")
}

// EventCode creates event.code (for keyboard events)
func EventCode() Callable {
	return Prop(Event, "code")
}

// EventKeyCode creates event.keyCode (for keyboard events, deprecated but common)
func EventKeyCode() Callable {
	return Prop(Event, "keyCode")
}

// EventWhich creates event.which (for keyboard events, deprecated but common)
func EventWhich() Callable {
	return Prop(Event, "which")
}

// EventShiftKey creates event.shiftKey
func EventShiftKey() Callable {
	return Prop(Event, "shiftKey")
}

// EventCtrlKey creates event.ctrlKey
func EventCtrlKey() Callable {
	return Prop(Event, "ctrlKey")
}

// EventAltKey creates event.altKey
func EventAltKey() Callable {
	return Prop(Event, "altKey")
}

// EventMetaKey creates event.metaKey
func EventMetaKey() Callable {
	return Prop(Event, "metaKey")
}

// Navigation helpers

// Navigate creates location.href = url
func Navigate(url Expr) Stmt {
	return Assign(Prop(Location, "href"), url)
}

// Reload creates location.reload()
func Reload() Callable {
	return Method(Location, "reload")
}

// HistoryBack creates history.back()
func HistoryBack() Callable {
	return Method(History, "back")
}

// HistoryForward creates history.forward()
func HistoryForward() Callable {
	return Method(History, "forward")
}

// HistoryGo creates history.go(delta)
func HistoryGo(delta Expr) Callable {
	return Method(History, "go", delta)
}

// HistoryPushState creates history.pushState(state, title, url)
func HistoryPushState(state, title, url Expr) Callable {
	return Method(History, "pushState", state, title, url)
}

// HistoryReplaceState creates history.replaceState(state, title, url)
func HistoryReplaceState(state, title, url Expr) Callable {
	return Method(History, "replaceState", state, title, url)
}

// Storage helpers

// GetItem creates storage.getItem(key)
func GetItem(storage Callable, key Expr) Callable {
	return Method(storage, "getItem", key)
}

// SetItem creates storage.setItem(key, value)
func SetItem(storage Callable, key, value Expr) Callable {
	return Method(storage, "setItem", key, value)
}

// RemoveItem creates storage.removeItem(key)
func RemoveItem(storage Callable, key Expr) Callable {
	return Method(storage, "removeItem", key)
}

// ClearStorage creates storage.clear()
func ClearStorage(storage Callable) Callable {
	return Method(storage, "clear")
}

// Focus/Blur helpers

// Focus creates element.focus()
func Focus(element Callable) Callable {
	return Method(element, "focus")
}

// Blur creates element.blur()
func Blur(element Callable) Callable {
	return Method(element, "blur")
}

// FocusThis creates this.focus() (for use in event handlers)
func FocusThis() Callable {
	return Method(EventThis, "focus")
}

// BlurThis creates this.blur() (for use in event handlers)
func BlurThis() Callable {
	return Method(EventThis, "blur")
}

// Click creates element.click()
func Click(element Callable) Callable {
	return Method(element, "click")
}

// Select creates element.select()
func Select(element Callable) Callable {
	return Method(element, "select")
}

// DOM manipulation helpers

// AppendChild creates parent.appendChild(child)
func AppendChild(parent, child Callable) Callable {
	return Method(parent, "appendChild", child)
}

// RemoveChild creates parent.removeChild(child)
func RemoveChild(parent, child Callable) Callable {
	return Method(parent, "removeChild", child)
}

// InsertBefore creates parent.insertBefore(newNode, referenceNode)
func InsertBefore(parent, newNode, referenceNode Callable) Callable {
	return Method(parent, "insertBefore", newNode, referenceNode)
}

// ReplaceChild creates parent.replaceChild(newChild, oldChild)
func ReplaceChild(parent, newChild, oldChild Callable) Callable {
	return Method(parent, "replaceChild", newChild, oldChild)
}

// Remove creates element.remove()
func Remove(element Callable) Callable {
	return Method(element, "remove")
}

// CloneNode creates element.cloneNode(deep)
func CloneNode(element Callable, deep Expr) Callable {
	return Method(element, "cloneNode", deep)
}

// ClassList helpers

// ClassList creates element.classList
func ClassList(element Callable) Callable {
	return Prop(element, "classList")
}

// ClassListAdd creates element.classList.add(classes...)
func ClassListAdd(element Callable, classes ...Expr) Callable {
	return Method(ClassList(element), "add", classes...)
}

// ClassListRemove creates element.classList.remove(classes...)
func ClassListRemove(element Callable, classes ...Expr) Callable {
	return Method(ClassList(element), "remove", classes...)
}

// ClassListToggle creates element.classList.toggle(className, force?)
func ClassListToggle(element Callable, className Expr, force ...Expr) Callable {
	args := make([]Expr, 1, 1+len(force))
	args[0] = className
	args = append(args, force...)
	return Method(ClassList(element), "toggle", args...)
}

// ClassListContains creates element.classList.contains(className)
func ClassListContains(element Callable, className Expr) Callable {
	return Method(ClassList(element), "contains", className)
}

// ClassListReplace creates element.classList.replace(oldClass, newClass)
func ClassListReplace(element Callable, oldClass, newClass Expr) Callable {
	return Method(ClassList(element), "replace", oldClass, newClass)
}

// Attribute helpers

// GetAttribute creates element.getAttribute(name)
func GetAttribute(element Callable, name Expr) Callable {
	return Method(element, "getAttribute", name)
}

// SetAttribute creates element.setAttribute(name, value)
func SetAttribute(element, name, value Expr) Callable {
	return Method(element.(Callable), "setAttribute", name, value)
}

// RemoveAttribute creates element.removeAttribute(name)
func RemoveAttribute(element Callable, name Expr) Callable {
	return Method(element, "removeAttribute", name)
}

// HasAttribute creates element.hasAttribute(name)
func HasAttribute(element Callable, name Expr) Callable {
	return Method(element, "hasAttribute", name)
}

// Style helpers

// Style creates element.style
func Style(element Callable) Callable {
	return Prop(element, "style")
}

// SetStyle creates element.style.property = value
func SetStyle(element Callable, property string, value Expr) Stmt {
	return Assign(Prop(Style(element), property), value)
}

// JSON helpers

// JSONStringify creates JSON.stringify(value, replacer?, space?)
func JSONStringify(value Expr, args ...Expr) Callable {
	allArgs := make([]Expr, 1, 1+len(args))
	allArgs[0] = value
	allArgs = append(allArgs, args...)
	return Method(JSON_, "stringify", allArgs...)
}

// JSONParse creates JSON.parse(text, reviver?)
func JSONParse(text Expr, reviver ...Expr) Callable {
	args := make([]Expr, 1, 1+len(reviver))
	args[0] = text
	args = append(args, reviver...)
	return Method(JSON_, "parse", args...)
}

// Common patterns

// ParseInt creates parseInt(string, radix)
func ParseInt(str Expr, radix ...Expr) Callable {
	args := make([]Expr, 1, 1+len(radix))
	args[0] = str
	args = append(args, radix...)
	return Call(Ident("parseInt"), args...)
}

// ParseFloat creates parseFloat(string)
func ParseFloat(str Expr) Callable {
	return Call(Ident("parseFloat"), str)
}

// IsNaN creates isNaN(value)
func IsNaN(value Expr) Callable {
	return Call(Ident("isNaN"), value)
}

// IsFinite creates isFinite(value)
func IsFinite(value Expr) Callable {
	return Call(Ident("isFinite"), value)
}

// Encodeuri creates encodeURI(uri)
func EncodeURI(uri Expr) Callable {
	return Call(Ident("encodeURI"), uri)
}

// DecodeURI creates decodeURI(uri)
func DecodeURI(uri Expr) Callable {
	return Call(Ident("decodeURI"), uri)
}

// EncodeURIComponent creates encodeURIComponent(component)
func EncodeURIComponent(component Expr) Callable {
	return Call(Ident("encodeURIComponent"), component)
}

// DecodeURIComponent creates decodeURIComponent(component)
func DecodeURIComponent(component Expr) Callable {
	return Call(Ident("decodeURIComponent"), component)
}
