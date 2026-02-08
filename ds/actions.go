package ds

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/jeffh/htmlgen/js"
)

type Stringer interface {
	String() string
}

// Get performs a GET request.
// Returns a Value that can be used with event handlers.
func Get(path string, chains ...PromiseChain) Value {
	return requestValue("get", js.String(path), chains...)
}

// Put performs a PUT request.
// Returns a Value that can be used with event handlers.
func Put(path string, chains ...PromiseChain) Value {
	return requestValue("put", js.String(path), chains...)
}

// Post performs a POST request.
// Returns a Value that can be used with event handlers.
func Post(path string, chains ...PromiseChain) Value {
	return requestValue("post", js.String(path), chains...)
}

// Delete performs a DELETE request.
// Returns a Value that can be used with event handlers.
func Delete(path string, chains ...PromiseChain) Value {
	return requestValue("delete", js.String(path), chains...)
}

// Patch performs a PATCH request.
// Returns a Value that can be used with event handlers.
func Patch(path string, chains ...PromiseChain) Value {
	return requestValue("patch", js.String(path), chains...)
}

// GetDynamic performs a GET request with a dynamic path.
// Returns a Value that can be used with event handlers.
func GetDynamic(path Value, chains ...PromiseChain) Value {
	return requestValue("get", path.expr, chains...)
}

// PutDynamic performs a PUT request with a dynamic path.
// Returns a Value that can be used with event handlers.
func PutDynamic(path Value, chains ...PromiseChain) Value {
	return requestValue("put", path.expr, chains...)
}

// PostDynamic performs a POST request with a dynamic path.
// Returns a Value that can be used with event handlers.
func PostDynamic(path Value, chains ...PromiseChain) Value {
	return requestValue("post", path.expr, chains...)
}

// DeleteDynamic performs a DELETE request with a dynamic path.
// Returns a Value that can be used with event handlers.
func DeleteDynamic(path Value, chains ...PromiseChain) Value {
	return requestValue("delete", path.expr, chains...)
}

// PatchDynamic performs a PATCH request with a dynamic path.
// Returns a Value that can be used with event handlers.
func PatchDynamic(path Value, chains ...PromiseChain) Value {
	return requestValue("patch", path.expr, chains...)
}

func requestValue(method string, path js.Expr, chains ...PromiseChain) Value {
	action := DatastarAction(method, path)
	if len(chains) > 0 {
		return Value{expr: WithChains(action, chains...)}
	}
	return Value{expr: action}
}

// GetWithOptions performs a GET request with options.
func GetWithOptions(path string, opts RequestOptionsBuilder, chains ...PromiseChain) Value {
	return requestValueWithOptions("get", js.String(path), opts, chains...)
}

// PostWithOptions performs a POST request with options.
func PostWithOptions(path string, opts RequestOptionsBuilder, chains ...PromiseChain) Value {
	return requestValueWithOptions("post", js.String(path), opts, chains...)
}

// PutWithOptions performs a PUT request with options.
func PutWithOptions(path string, opts RequestOptionsBuilder, chains ...PromiseChain) Value {
	return requestValueWithOptions("put", js.String(path), opts, chains...)
}

// DeleteWithOptions performs a DELETE request with options.
func DeleteWithOptions(path string, opts RequestOptionsBuilder, chains ...PromiseChain) Value {
	return requestValueWithOptions("delete", js.String(path), opts, chains...)
}

// PatchWithOptions performs a PATCH request with options.
func PatchWithOptions(path string, opts RequestOptionsBuilder, chains ...PromiseChain) Value {
	return requestValueWithOptions("patch", js.String(path), opts, chains...)
}

func requestValueWithOptions(method string, path js.Expr, opts RequestOptionsBuilder, chains ...PromiseChain) Value {
	var sb strings.Builder
	sb.WriteString("@")
	sb.WriteString(method)
	sb.WriteString("(")
	sb.WriteString(js.ToJS(path))
	if len(opts.options) > 0 {
		sb.WriteString(", {")
		for i, opt := range opts.options {
			if i > 0 {
				sb.WriteString(", ")
			}
			opt.appendOption(&sb)
		}
		sb.WriteString("}")
	}
	sb.WriteString(")")
	action := js.Raw(sb.String())
	if len(chains) > 0 {
		return Value{expr: WithChains(action, chains...)}
	}
	return Value{expr: action}
}

// RequestOptionsBuilder collects request options.
type RequestOptionsBuilder struct {
	options []RequestOption
}

// RequestOptions creates a builder for HTTP request options.
// Use with GetWithOptions, PostWithOptions, etc.
// Example: GetWithOptions("/api", RequestOptions().ContentType("form").OpenWhenHidden(true))
func RequestOptions() RequestOptionsBuilder {
	return RequestOptionsBuilder{}
}

// ContentType sets the request content type.
// Values: "json" (default) or "form"
func (b RequestOptionsBuilder) ContentType(ct string) RequestOptionsBuilder {
	b.options = append(b.options, requestOptionFunc(func(sb *strings.Builder) {
		sb.WriteString("contentType: ")
		sb.WriteString(strconv.Quote(ct))
	}))
	return b
}

// FilterSignals filters which signals are sent with the request.
func (b RequestOptionsBuilder) FilterSignals(filter *FilterOptions) RequestOptionsBuilder {
	b.options = append(b.options, requestOptionFunc(func(sb *strings.Builder) {
		sb.WriteString("filterSignals: ")
		filter.appendJS(sb)
	}))
	return b
}

// Selector specifies a CSS selector for form elements (when contentType is 'form').
func (b RequestOptionsBuilder) Selector(sel string) RequestOptionsBuilder {
	b.options = append(b.options, requestOptionFunc(func(sb *strings.Builder) {
		sb.WriteString("selector: ")
		sb.WriteString(strconv.Quote(sel))
	}))
	return b
}

// Headers sets custom HTTP headers for the request.
func (b RequestOptionsBuilder) Headers(headers map[string]string) RequestOptionsBuilder {
	b.options = append(b.options, requestOptionFunc(func(sb *strings.Builder) {
		sb.WriteString("headers: {")
		i := 0
		for k, v := range headers {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(strconv.Quote(k))
			sb.WriteString(": ")
			sb.WriteString(strconv.Quote(v))
			i++
		}
		sb.WriteString("}")
	}))
	return b
}

// OpenWhenHidden keeps the connection alive when the tab is hidden.
func (b RequestOptionsBuilder) OpenWhenHidden(open bool) RequestOptionsBuilder {
	b.options = append(b.options, requestOptionFunc(func(sb *strings.Builder) {
		sb.WriteString("openWhenHidden: ")
		if open {
			sb.WriteString("true")
		} else {
			sb.WriteString("false")
		}
	}))
	return b
}

// RetryInterval sets the retry interval in milliseconds (default: 1000).
func (b RequestOptionsBuilder) RetryInterval(ms int) RequestOptionsBuilder {
	b.options = append(b.options, requestOptionFunc(func(sb *strings.Builder) {
		sb.WriteString("retryInterval: ")
		sb.WriteString(strconv.Itoa(ms))
	}))
	return b
}

// RetryScaler sets the exponential backoff multiplier (default: 2).
func (b RequestOptionsBuilder) RetryScaler(scaler float64) RequestOptionsBuilder {
	b.options = append(b.options, requestOptionFunc(func(sb *strings.Builder) {
		sb.WriteString("retryScaler: ")
		sb.WriteString(strconv.FormatFloat(scaler, 'f', -1, 64))
	}))
	return b
}

// RetryMaxWaitMs sets the maximum wait between retries in milliseconds (default: 30000).
func (b RequestOptionsBuilder) RetryMaxWaitMs(ms int) RequestOptionsBuilder {
	b.options = append(b.options, requestOptionFunc(func(sb *strings.Builder) {
		sb.WriteString("retryMaxWaitMs: ")
		sb.WriteString(strconv.Itoa(ms))
	}))
	return b
}

// RetryMaxCount sets the maximum number of retry attempts (default: 10).
func (b RequestOptionsBuilder) RetryMaxCount(count int) RequestOptionsBuilder {
	b.options = append(b.options, requestOptionFunc(func(sb *strings.Builder) {
		sb.WriteString("retryMaxCount: ")
		sb.WriteString(strconv.Itoa(count))
	}))
	return b
}

// RequestCancellation sets the request cancellation mode.
// Values: "auto" (default), "disabled"
func (b RequestOptionsBuilder) RequestCancellation(mode string) RequestOptionsBuilder {
	b.options = append(b.options, requestOptionFunc(func(sb *strings.Builder) {
		sb.WriteString("requestCancellation: ")
		sb.WriteString(strconv.Quote(mode))
	}))
	return b
}

// Retry sets the retry strategy for the request.
// Values: "auto" (default, retry on network errors), "error" (retry on errors and non-2xx),
// "always" (always retry), "never" (never retry)
func (b RequestOptionsBuilder) Retry(mode string) RequestOptionsBuilder {
	b.options = append(b.options, requestOptionFunc(func(sb *strings.Builder) {
		sb.WriteString("retry: ")
		sb.WriteString(strconv.Quote(mode))
	}))
	return b
}

// Payload overrides the request body with custom JSON data.
// Use for POST/PUT/PATCH requests when you need a custom payload.
func (b RequestOptionsBuilder) Payload(data any) RequestOptionsBuilder {
	b.options = append(b.options, requestOptionFunc(func(sb *strings.Builder) {
		sb.WriteString("body: ")
		bytes, err := json.Marshal(data)
		if err != nil {
			panic(fmt.Errorf("Payload: %w: value=%#v", err, data))
		}
		sb.WriteString(string(bytes))
	}))
	return b
}

// RequestOption is an option for HTTP request actions.
type RequestOption interface {
	appendOption(*strings.Builder)
}

type requestOptionFunc func(*strings.Builder)

func (f requestOptionFunc) appendOption(sb *strings.Builder) { f(sb) }

// OnSuccess creates a .then() chain for successful request handling.
// This is an alias for ThenChain for backward compatibility.
func OnSuccess(expr Value) PromiseChain {
	return ThenChain(expr.expr)
}

// OnFailure creates a .catch() chain for error handling.
// This is an alias for CatchChain for backward compatibility.
func OnFailure(expr Value) PromiseChain {
	return CatchChain(expr.expr)
}
