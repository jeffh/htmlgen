package html

import (
	"errors"
	"fmt"
	"html/template"
	"io"
)

var ErrUnknownTagToClose = errors.New("attempted to close tag not opened")

func NewWriter(w io.Writer) *Writer {
	return &Writer{w, make([]string, 0, 16)}
}

type Writer struct {
	w        io.Writer
	openTags []string
}

func (w *Writer) write(values ...string) error {
	for _, v := range values {
		_, err := w.w.Write([]byte(v))
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *Writer) Doctype() error { return w.write("<!DOCTYPE html>\n") }

func (w *Writer) SelfClosingTag(name string, attrvals ...string) error {
	if len(attrvals)%2 != 0 {
		panic("attrvals must be an even number of items")
	}
	if err := w.write("<", name); err != nil {
		return err
	}
	for i := 0; i < len(attrvals); i += 2 {
		val := attrvals[i+1]
		if val == "" {
			if err := w.write(" ", attrvals[i]); err != nil {
				return err
			}
		} else {
			if err := w.write(" ", attrvals[i], "=\"", template.HTMLEscapeString(val), "\""); err != nil {
				return err
			}
		}
	}
	if err := w.write("/>"); err != nil {
		return err
	}
	return nil
}

func (w *Writer) OpenTag(name string, attrvals ...string) error {
	if len(attrvals)%2 != 0 {
		panic("attrvals must be an even number of items")
	}
	if err := w.write("<", name); err != nil {
		return err
	}
	for i := 0; i < len(attrvals); i += 2 {
		val := attrvals[i+1]
		if val == "" {
			if err := w.write(" ", attrvals[i]); err != nil {
				return err
			}
		} else {
			if err := w.write(" ", attrvals[i], "=\"", template.HTMLEscapeString(val), "\""); err != nil {
				return err
			}
		}
	}
	if err := w.write(">"); err != nil {
		return err
	}
	w.openTags = append(w.openTags, name)
	return nil
}

func (w *Writer) Text(txt string) error       { return w.write(template.HTMLEscapeString(txt)) }
func (w *Writer) Raw(unsafeHtml string) error { return w.write(unsafeHtml) }

func (w *Writer) CloseTag(name string) error {
	size := len(w.openTags)
	if size == 0 {
		return fmt.Errorf("%w: %s", ErrUnknownTagToClose, name)
	}
	for i := size - 1; i >= 0; i-- {
		if w.openTags[i] == name {
			for j := size - 1; j >= i; j-- {
				if err := w.write("</", w.openTags[j], ">"); err != nil {
					return err
				}
			}
			w.openTags = w.openTags[:i]
			break
		}
	}
	return nil
}

func (w *Writer) End() error {
	size := len(w.openTags)
	if size == 0 {
		return ErrUnknownTagToClose
	}
	err := w.write("</", w.openTags[size-1], ">")
	if err != nil {
		return err
	}
	w.openTags = w.openTags[:size-1]
	return nil
}

// Close closes all opened tags
func (w *Writer) Close() error {
	for i := len(w.openTags) - 1; i >= 0; i-- {
		if err := w.write("</", w.openTags[i], ">"); err != nil {
			return err
		}
	}
	w.openTags = nil
	return nil
}
