package js

import (
	"strings"
	"testing"
)

func TestTemplatePartTypes(t *testing.T) {
	tests := []struct {
		name      string
		parts     []any
		want      string
		wantPanic string
	}{
		{
			name:  "happy path",
			parts: []any{"Hello, ", Ident("name"), "!"},
			want:  "`Hello, ${name}!`",
		},
		{
			name:      "int",
			parts:     []any{"Hello, ", 42, "!"},
			wantPanic: "js.Template",
		},
		{
			name:      "stmt",
			parts:     []any{"n=", Ident("n").Incr()},
			wantPanic: "js.Template",
		},
		{
			name:      "nil",
			parts:     []any{"n=", nil},
			wantPanic: "js.Template",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if tt.wantPanic == "" {
					if r != nil {
						t.Fatalf("Template(%v) unexpected panic: %#v", tt.parts, r)
					}
					return
				}
				if r == nil {
					t.Fatalf("Template(%v) did not panic, want substring %q", tt.parts, tt.wantPanic)
				}
				msg, ok := r.(string)
				if !ok {
					t.Fatalf("Template(%v) panic value = %#v, want string", tt.parts, r)
				}
				if !strings.Contains(msg, tt.wantPanic) {
					t.Errorf("Template(%v) panic = %q, want substring %q", tt.parts, msg, tt.wantPanic)
				}
			}()

			got := exprString(Template(tt.parts...))
			if tt.wantPanic == "" && got != tt.want {
				t.Errorf("Template(%v) = %q, want %q", tt.parts, got, tt.want)
			}
		})
	}
}
