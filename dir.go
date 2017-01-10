package gopher

import (
	"errors"
	"fmt"
	"io"
)

// A DirWriter is used to write gopher directories to a connection.
type DirWriter struct {
	w *Writer
	io.WriteCloser
}

// A DirEntry represents a gopher directory entry.
type DirEntry struct {
	Type byte
	Name string
	Path string
	Host string
	Port string
}

// Entry writes a direntry to the current connection.
func (dw *DirWriter) Entry(e *DirEntry) error {
	_, err := fmt.Fprintf(dw, "%c%s\t%s\t%s\t%s\n",
		e.Type, e.Name, e.Path, e.Host, e.Port)
	return err
}

// LocalEntry works like Entry but uses ExtHost and ExtPort from the server
// instead of the ones supplied in e.
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
