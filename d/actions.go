package d

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type Stringer interface {
	String() string
}

// Performs a GET request.
func Get(path string, values ...AttrValueAppender) ValueMutator {
	return request("get", path, values...)
}

// Performs a PUT request.
func Put(path string, values ...AttrValueAppender) ValueMutator {
	return request("put", path, values...)
}

// Performs a POST request.
func Post(path string, values ...AttrValueAppender) ValueMutator {
	return request("post", path, values...)
}

// Performs a DELETE request.
func Delete(path string, values ...AttrValueAppender) ValueMutator {
	return request("delete", path, values...)
}

// Performs a PATCH request.
func Patch(path string, values ...AttrValueAppender) ValueMutator {
	return request("patch", path, values...)
}

// Performs a GET request.
func GetDynamic(path AttrValueAppender, values ...AttrValueAppender) ValueMutator {
	return requestDynamicPath("get", path, values...)
}

// Performs a PUT request.
func PutDynamic(path AttrValueAppender, values ...AttrValueAppender) ValueMutator {
	return requestDynamicPath("put", path, values...)
}

// Performs a POST request.
func PostDynamic(path AttrValueAppender, values ...AttrValueAppender) ValueMutator {
	return requestDynamicPath("post", path, values...)
}

// Performs a DELETE request.
func DeleteDynamic(path AttrValueAppender, values ...AttrValueAppender) ValueMutator {
	return requestDynamicPath("delete", path, values...)
}

// Performs a PATCH request with a dynamic path.
func PatchDynamic(path AttrValueAppender, values ...AttrValueAppender) ValueMutator {
	return requestDynamicPath("patch", path, values...)
}

func request(method string, path string, options ...AttrValueAppender) ValueMutator {
	return ValueMutator{
		AttrMutator: AttrFunc(func(attr *attrBuilder) {
			sb := strings.Builder{}
			sb.Grow(5 + len(method) + len(path) + 30) // 30 for small strings
			sb.WriteString("@")
			sb.WriteString(method)
			sb.WriteString("(")
			sb.WriteString(strconv.Quote(path))
			sb.WriteString(")")
			for _, opt := range options {
				opt.Append(&sb)
			}
			attr.AppendStatement(sb.String())
		}),
		AttrValueAppender: AttrValueFunc(func(sb *strings.Builder) {
			sb.Grow(5 + len(method) + len(path) + 30) // 30 for small strings
			sb.WriteString("@")
			sb.WriteString(method)
			sb.WriteString("(")
			sb.WriteString(strconv.Quote(path))
			sb.WriteString(")")
			for _, opt := range options {
				opt.Append(sb)
			}
		}),
	}
}

func requestDynamicPath(method string, path AttrValueAppender, options ...AttrValueAppender) ValueMutator {
	return ValueMutator{
		AttrMutator: AttrFunc(func(attr *attrBuilder) {
			sb := strings.Builder{}
			sb.Grow(5 + len(method) + 30) // 30 for small strings
			sb.WriteString("@")
			sb.WriteString(method)
			sb.WriteString("(")
			path.Append(&sb)
			for _, opt := range options {
				opt.Append(&sb)
			}
			sb.WriteString(")")
			attr.AppendStatement(sb.String())
		}),
		AttrValueAppender: AttrValueFunc(func(sb *strings.Builder) {
			sb.Grow(5 + len(method) + 30) // 30 for small strings
			sb.WriteString("@")
			sb.WriteString(method)
			sb.WriteString("(")
			path.Append(sb)
			for _, opt := range options {
				opt.Append(sb)
			}
			sb.WriteString(")")
		}),
	}
}

// RequestOptions builds an options object for HTTP request actions.
// Use with Get, Post, Put, Delete, Patch to configure request behavior.
// Example: Get("/api", RequestOptions(ContentType("form"), OpenWhenHidden(true)))
func RequestOptions(opts ...RequestOption) AttrValueAppender {
	return AttrValueFunc(func(sb *strings.Builder) {
		if len(opts) == 0 {
			return
		}
		sb.WriteString(", {")
		for i, opt := range opts {
			if i > 0 {
				sb.WriteString(", ")
			}
			opt.appendOption(sb)
		}
		sb.WriteString("}")
	})
}

// RequestOption is an option for HTTP request actions.
type RequestOption interface {
	appendOption(*strings.Builder)
}

type requestOptionFunc func(*strings.Builder)

func (f requestOptionFunc) appendOption(sb *strings.Builder) { f(sb) }

// ContentType sets the request content type.
// Values: "json" (default) or "form"
func ContentType(ct string) RequestOption {
	return requestOptionFunc(func(sb *strings.Builder) {
		sb.WriteString("contentType: ")
		sb.WriteString(strconv.Quote(ct))
	})
}

// FilterSignals filters which signals are sent with the request.
func FilterSignals(filter *FilterOptions) RequestOption {
	return requestOptionFunc(func(sb *strings.Builder) {
		sb.WriteString("filterSignals: ")
		filter.appendJS(sb)
	})
}

// Selector specifies a CSS selector for form elements (when contentType is 'form').
func Selector(sel string) RequestOption {
	return requestOptionFunc(func(sb *strings.Builder) {
		sb.WriteString("selector: ")
		sb.WriteString(strconv.Quote(sel))
	})
}

// Headers sets custom HTTP headers for the request.
func Headers(headers map[string]string) RequestOption {
	return requestOptionFunc(func(sb *strings.Builder) {
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
	})
}

// OpenWhenHidden keeps the connection alive when the tab is hidden.
func OpenWhenHidden(open bool) RequestOption {
	return requestOptionFunc(func(sb *strings.Builder) {
		sb.WriteString("openWhenHidden: ")
		if open {
			sb.WriteString("true")
		} else {
			sb.WriteString("false")
		}
	})
}

// RetryInterval sets the retry interval in milliseconds (default: 1000).
func RetryInterval(ms int) RequestOption {
	return requestOptionFunc(func(sb *strings.Builder) {
		sb.WriteString("retryInterval: ")
		sb.WriteString(strconv.Itoa(ms))
	})
}

// RetryScaler sets the exponential backoff multiplier (default: 2).
func RetryScaler(scaler float64) RequestOption {
	return requestOptionFunc(func(sb *strings.Builder) {
		sb.WriteString("retryScaler: ")
		sb.WriteString(strconv.FormatFloat(scaler, 'f', -1, 64))
	})
}

// RetryMaxWaitMs sets the maximum wait between retries in milliseconds (default: 30000).
func RetryMaxWaitMs(ms int) RequestOption {
	return requestOptionFunc(func(sb *strings.Builder) {
		sb.WriteString("retryMaxWaitMs: ")
		sb.WriteString(strconv.Itoa(ms))
	})
}

// RetryMaxCount sets the maximum number of retry attempts (default: 10).
func RetryMaxCount(count int) RequestOption {
	return requestOptionFunc(func(sb *strings.Builder) {
		sb.WriteString("retryMaxCount: ")
		sb.WriteString(strconv.Itoa(count))
	})
}

// RequestCancellation sets the request cancellation mode.
// Values: "auto" (default), "disabled"
func RequestCancellation(mode string) RequestOption {
	return requestOptionFunc(func(sb *strings.Builder) {
		sb.WriteString("requestCancellation: ")
		sb.WriteString(strconv.Quote(mode))
	})
}

// Retry sets the retry strategy for the request.
// Values: "auto" (default, retry on network errors), "error" (retry on errors and non-2xx),
// "always" (always retry), "never" (never retry)
func Retry(mode string) RequestOption {
	return requestOptionFunc(func(sb *strings.Builder) {
		sb.WriteString("retry: ")
		sb.WriteString(strconv.Quote(mode))
	})
}

// Payload overrides the request body with custom JSON data.
// Use for POST/PUT/PATCH requests when you need a custom payload.
func Payload(data any) RequestOption {
	return requestOptionFunc(func(sb *strings.Builder) {
		sb.WriteString("body: ")
		b, err := json.Marshal(data)
		if err != nil {
			panic(fmt.Errorf("Payload: %w: value=%#v", err, data))
		}
		sb.WriteString(string(b))
	})
}
