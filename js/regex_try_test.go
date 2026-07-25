package js

import "testing"

func TestRegex(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		flags    string
		expected string
	}{
		{"anchored with flag", "^user_", "i", `/^user_/i`},
		{"no flags", `\d+`, "", `/\d+/`},
		{"multiple flags", ".", "gis", `/./gis`},
		{"escaped slash", `a\/b`, "g", `/a\/b/g`},
		{"empty pattern", "", "", `//`},
		{"verbatim, no escaping", `"quoted"`, "", `/"quoted"/`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := exprString(Regex(tt.pattern, tt.flags))
			if got != tt.expected {
				t.Errorf("Regex(%q, %q) = %q, want %q", tt.pattern, tt.flags, got, tt.expected)
			}
		})
	}
}

func TestRegexInExpression(t *testing.T) {
	got := exprString(Ident("s").Method("match", Regex("x", "g")))
	want := `s.match(/x/g)`
	if got != want {
		t.Errorf("match = %q, want %q", got, want)
	}
	got = exprString(Regex("a", "").Method("test", Ident("s")))
	want = `/a/.test(s)`
	if got != want {
		t.Errorf("test = %q, want %q", got, want)
	}
}

func TestTryCatch(t *testing.T) {
	tests := []struct {
		name      string
		body      []Stmt
		errName   string
		catchBody []Stmt
		expected  string
	}{
		{
			"bound error",
			[]Stmt{ExprStmt(ConsoleLog(Int(1)))},
			"e",
			[]Stmt{ExprStmt(ConsoleError(Ident("e")))},
			`try { console.log(1) } catch (e) { console.error(e) }`,
		},
		{
			"optional binding omitted",
			[]Stmt{Ident("x").Incr()},
			"",
			nil,
			`try { x++ } catch {}`,
		},
		{
			"optional binding with catch body",
			[]Stmt{Ident("x").Incr()},
			"",
			[]Stmt{Ident("y").Incr()},
			`try { x++ } catch { y++ }`,
		},
		{
			"multiple statements",
			[]Stmt{Let("n", Int(0)), Ident("n").Incr()},
			"err",
			[]Stmt{ExprStmt(ConsoleError(Ident("err")))},
			`try { let n = 0; n++ } catch (err) { console.error(err) }`,
		},
		{
			"empty body",
			nil,
			"e",
			nil,
			`try {} catch (e) {}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stmtString(TryCatch(tt.body, tt.errName, tt.catchBody))
			if got != tt.expected {
				t.Errorf("TryCatch = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestTry(t *testing.T) {
	tests := []struct {
		name     string
		body     []Stmt
		expected string
	}{
		{"single statement", []Stmt{ExprStmt(Ident("el").Method("focus"))}, `try { el.focus() } catch {}`},
		{"multiple statements", []Stmt{Ident("a").Incr(), Ident("b").Incr()}, `try { a++; b++ } catch {}`},
		{"no statements", nil, `try {} catch {}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stmtString(Try(tt.body...))
			if got != tt.expected {
				t.Errorf("Try = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestTryInHandler(t *testing.T) {
	got := Handler(Let("n", Int(0)), Try(Ident("n").Incr()))
	want := `let n = 0; try { n++ } catch {}`
	if got != want {
		t.Errorf("Handler = %q, want %q", got, want)
	}
}

func TestTryNested(t *testing.T) {
	got := stmtString(Try(Try(Ident("x").Incr())))
	want := `try { try { x++ } catch {} } catch {}`
	if got != want {
		t.Errorf("nested Try = %q, want %q", got, want)
	}
}
