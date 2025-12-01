package html

import (
	"bytes"
	"testing"
)

func TestHtmlWriting(t *testing.T) {
	cases := []struct {
		Desc     string
		Expected string
		Code     func(w *Writer)
	}{
		{
			"Basic HTML Writing",
			`<!DOCTYPE html>
<html lang="en"><head><link rel="stylesheet" href="style.css" rel="preload"/></head><body><script src="script.js" defer></script></body></html>`,
			func(w *Writer) {
				w.Doctype()
				w.BegTag("html", "lang", "en")
				{
					w.BegTag("head")
					{
						Stylesheet(w, "style.css")
					}
					w.End()
					w.BegTag("body")
					{
						Script(w, "script.js")
					}
					w.EndTag("body")
				}
			},
		},
        {
        },
	}

	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)

			w := NewWriter(buf)
			c.Code(w)
			w.Close()

			out := buf.String()

			if out != c.Expected {
				t.Errorf("expected\n%q\n  to equal\n%q", c.Expected, out)
			}
		})
	}

}
