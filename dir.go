package gopher

import (
	"errors"
	"fmt"
	"io"
)

type DirWriter struct {
	w *Writer
	io.WriteCloser
}

type DirEntry struct {
	Type byte
	Name string
	Path string
	Host string
	Port string
}

func (dw *DirWriter) Entry(e *DirEntry) error {
	_, err := fmt.Fprintf(dw, "%c%s\t%s\t%s\t%s\n",
		e.Type, e.Name, e.Path, e.Host, e.Port)
	return err
}

func (dw *DirWriter) LocalEntry(e *DirEntry) error {
	e.Host = dw.w.srv.ExtHost
	e.Port = dw.w.srv.ExtPort
	if e.Host == "" {
		return errors.New("DirWriter.LocalEntry: missing ExtHost")
	}
	if e.Port == "" {
		return errors.New("DirWriter.LocalEntry: missing ExtPort")
	}
	return dw.Entry(e)
}
