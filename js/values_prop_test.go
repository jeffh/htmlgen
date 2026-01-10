package js

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/jeffh/gocheck/check"
	"github.com/jeffh/gocheck/gen"
)

func TestWriteJSONString_MatchesJSONMarshal(t *testing.T) {
	cfg := check.Config{}
	str := gen.OneOf(gen.StringASCII, gen.StringPrintable, gen.StringAny)

	check.Spec(t, cfg, gen.ForAll(str, func(s string) check.B {
		// Get writeJSONString output
		var sb strings.Builder
		writeJSONString(&sb, s)
		got := sb.String()

		// Get json.Marshal output
		want, err := json.Marshal(s)
		if err != nil {
			return check.Fail(func() string {
				return "json.Marshal failed: " + err.Error()
			})
		}

		if got != string(want) {
			return check.Fail(func() string {
				return "mismatch:\n  input:  " + s + "\n  got:    " + got + "\n  want:   " + string(want)
			})
		}
		return true
	}))
}
