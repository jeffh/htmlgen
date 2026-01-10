package hx

import (
	"encoding/json"
	"strings"

	"github.com/jeffh/htmlgen/h"
)

// Get creates an hx-get attribute that issues a GET request to the specified URL.
func Get(url string) h.Attribute {
	return h.Attr("hx-get", url)
}

// Post creates an hx-post attribute that issues a POST request to the specified URL.
func Post(url string) h.Attribute {
	return h.Attr("hx-post", url)
}

// Put creates an hx-put attribute that issues a PUT request to the specified URL.
func Put(url string) h.Attribute {
	return h.Attr("hx-put", url)
}

// Patch creates an hx-patch attribute that issues a PATCH request to the specified URL.
func Patch(url string) h.Attribute {
	return h.Attr("hx-patch", url)
}

// Delete creates an hx-delete attribute that issues a DELETE request to the specified URL.
func Delete(url string) h.Attribute {
	return h.Attr("hx-delete", url)
}

// Target creates an hx-target attribute that specifies the target element for swapping.
//
// Values:
//   - A CSS selector (e.g., "#myDiv", ".myClass")
//   - "this" for the triggering element
//   - "closest <selector>" for the closest ancestor
//   - "find <selector>" for the first matching descendant
//   - "next" or "next <selector>" for the next sibling
//   - "previous" or "previous <selector>" for the previous sibling
func Target(selector string) h.Attribute {
	return h.Attr("hx-target", selector)
}

// Select creates an hx-select attribute that extracts specific content from the response.
func Select(selector string) h.Attribute {
	return h.Attr("hx-select", selector)
}

// SelectOOB creates an hx-select-oob attribute for out-of-band content selection.
func SelectOOB(selector string) h.Attribute {
	return h.Attr("hx-select-oob", selector)
}

// SwapOOB creates an hx-swap-oob attribute for out-of-band swapping.
//
// Values:
//   - "true" for default behavior (swap by id)
//   - A swap strategy (e.g., "outerHTML", "innerHTML")
//   - A swap strategy with selector (e.g., "outerHTML:#target")
func SwapOOB(value string) h.Attribute {
	return h.Attr("hx-swap-oob", value)
}

// Include creates an hx-include attribute that includes additional elements in the request.
//
// Values:
//   - A CSS selector (e.g., "#form")
//   - "this" for the triggering element
//   - "closest <selector>" for the closest ancestor
func Include(selector string) h.Attribute {
	return h.Attr("hx-include", selector)
}

// Vals creates an hx-vals attribute with JSON-encoded values.
func Vals(values map[string]any) h.Attribute {
	data, err := json.Marshal(values)
	if err != nil {
		panic("hx.Vals: " + err.Error())
	}
	return h.Attr("hx-vals", string(data))
}

// ValsJS creates an hx-vals attribute with JavaScript expressions.
// The values are prefixed with "js:" to indicate they are JavaScript expressions.
func ValsJS(values map[string]string) h.Attribute {
	parts := make([]string, 0, len(values))
	for k, v := range values {
		parts = append(parts, `"`+k+`": `+v)
	}
	return h.Attr("hx-vals", "js:{"+strings.Join(parts, ", ")+"}")
}

// Headers creates an hx-headers attribute with JSON-encoded headers.
func Headers(headers map[string]string) h.Attribute {
	data, err := json.Marshal(headers)
	if err != nil {
		panic("hx.Headers: " + err.Error())
	}
	return h.Attr("hx-headers", string(data))
}

// HeadersJS creates an hx-headers attribute with a JavaScript expression.
// The expression should evaluate to a JavaScript object.
//
// Example:
//
//	hx.HeadersJS("getAuthHeaders()")
func HeadersJS(jsExpr string) h.Attribute {
	return h.Attr("hx-headers", "js:"+jsExpr)
}

// ParamsFilter specifies how parameters should be filtered.
type ParamsFilter string

const (
	// ParamsAll includes all parameters (default).
	ParamsAll ParamsFilter = "*"
	// ParamsNone excludes all parameters.
	ParamsNone ParamsFilter = "none"
)

// Params creates an hx-params attribute that filters request parameters.
//
// Values:
//   - "*" for all parameters (default)
//   - "none" to exclude all parameters
//   - "not <param-list>" to exclude specific params (e.g., "not secret,password")
//   - A comma-separated list to include only those params (e.g., "name,email")
func Params(filter string) h.Attribute {
	return h.Attr("hx-params", filter)
}

// EncodingType specifies the encoding type for requests.
type EncodingType string

const (
	// EncodingForm uses application/x-www-form-urlencoded encoding.
	EncodingForm EncodingType = "application/x-www-form-urlencoded"
	// EncodingMultipart uses multipart/form-data encoding.
	EncodingMultipart EncodingType = "multipart/form-data"
)

// Encoding creates an hx-encoding attribute that sets the request encoding.
func Encoding(encoding EncodingType) h.Attribute {
	return h.Attr("hx-encoding", string(encoding))
}

// Ext creates an hx-ext attribute that specifies extensions to use.
// To ignore a parent's extension, prefix with "ignore:" (e.g., "ignore:debug").
//
// Example:
//
//	hx.Ext("json-enc", "debug")           // use extensions
//	hx.Ext("ignore:debug")                // ignore parent's debug extension
func Ext(extensions ...string) h.Attribute {
	return h.Attr("hx-ext", strings.Join(extensions, ","))
}

// Request creates an hx-request attribute for configuring request behavior.
// The config is a JSON object with options like timeout, credentials, noHeaders.
func Request(config map[string]any) h.Attribute {
	data, err := json.Marshal(config)
	if err != nil {
		panic("hx.Request: " + err.Error())
	}
	return h.Attr("hx-request", string(data))
}

// RequestJS creates an hx-request attribute with a JavaScript expression.
func RequestJS(jsExpr string) h.Attribute {
	return h.Attr("hx-request", "js:"+jsExpr)
}
