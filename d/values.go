package d

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Raw injects a raw JavaScript value into the element. This can be unsafe if user input is used.
func Raw(value string) ValueMutator {
	return ValueMutator{
		AttrMutator:       AttrFunc(func(attr *attrBuilder) { attr.AppendStatement(value) }),
		AttrValueAppender: AttrValueFunc(func(sb *strings.Builder) { sb.WriteString(value) }),
	}
}

func Str(value string) ValueMutator { return JsonValue(value) }

// JsonValue injects a JSON value.
func JsonValue(value any) ValueMutator {
	return ValueMutator{
		AttrMutator: AttrFunc(func(attr *attrBuilder) {
			b, err := json.Marshal(value)
			if err != nil {
				panic(fmt.Errorf("%w: value=%#v", err, value))
			}
			attr.AppendStatement(string(b))
		}),
		AttrValueAppender: AttrValueFunc(func(sb *strings.Builder) {
			b, err := json.Marshal(value)
			if err != nil {
				panic(fmt.Errorf("%w: value=%#v", err, value))
			}
			sb.WriteString(string(b))
		}),
	}
}
func appendName(name string) AttrMutator {
	return AttrFunc(func(attr *attrBuilder) {
		attr.name.WriteString(name)
	})
}

func Navigate(path string, values ...any) string {
	var sb strings.Builder
	sb.Grow(len(path) + 10)
	sb.WriteString("window.location.href = ")
	Str(fmt.Sprintf(path, values...)).Append(&sb)
	return sb.String()
}

// OnSuccess injects a JavaScript expression to be executed when an action succeeds.
func OnSuccess(rawJS string) AttrValueAppender {
	return AttrValueFunc(func(sb *strings.Builder) {
		sb.Grow(len(rawJS) + 10)
		sb.WriteString(".then(() => ")
		sb.WriteString(rawJS)
		sb.WriteString(")")
	})
}

// OnFailure injects a JavaScript expression to be executed when an action fails.
func OnFailure(rawJS string) AttrValueAppender {
	return AttrValueFunc(func(sb *strings.Builder) {
		sb.Grow(len(rawJS) + 10)
		sb.WriteString(".catch((error) => ")
		sb.WriteString(rawJS)
		sb.WriteString(")")
	})
}

// ConsoleLog injects a console.log statement.
func ConsoleLog(values ...AttrValueAppender) AttrMutator {
	return AttrFunc(func(attr *attrBuilder) {
		sb := strings.Builder{}
		sb.Grow(len(values) * 10)
		sb.WriteString("console.log(")
		for i, value := range values {
			if i != 0 {
				sb.WriteString(", ")
			}
			value.Append(&sb)
		}
		sb.WriteString(")")
		attr.AppendStatement(sb.String())
	})
}

func And(actions ...AttrValueAppender) ValueMutator {
	return ValueMutator{
		AttrMutator: AttrFunc(func(attr *attrBuilder) {
			var parts []string
			var temp strings.Builder
			for _, action := range actions {
				action.Append(&temp)
				parts = append(parts, temp.String())
				temp.Reset()
			}
			attr.AppendStatement(strings.Join(parts, " && "))
		}),
		AttrValueAppender: AttrValueFunc(func(sb *strings.Builder) {
			var parts []string
			var temp strings.Builder
			sb.Grow(len(actions) * 10)
			for _, action := range actions {
				action.Append(&temp)
				parts = append(parts, temp.String())
				temp.Reset()
			}
			sb.WriteString(strings.Join(parts, " && "))
		}),
	}
}
