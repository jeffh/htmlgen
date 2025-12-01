package d

import (
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
			sb.WriteString(")")
			for _, opt := range options {
				opt.Append(&sb)
			}
			attr.AppendStatement(sb.String())
		}),
		AttrValueAppender: AttrValueFunc(func(sb *strings.Builder) {
			sb.Grow(5 + len(method) + 30) // 30 for small strings
			sb.WriteString("@")
			sb.WriteString(method)
			sb.WriteString("(")
			path.Append(sb)
			sb.WriteString(")")
			for _, opt := range options {
				opt.Append(sb)
			}
		}),
	}
}
