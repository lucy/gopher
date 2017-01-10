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

// Entry writes a DirEntry to the current connection.
func (dw *DirWriter) Entry(e *DirEntry) error {
	_, err := fmt.Fprintf(dw, "%c%s\t%s\t%s\t%s\n",
		e.Type, e.Name, e.Path, e.Host, e.Port)
	return err
}

// errors for LocalEntry
var errMissingExtHost = errors.New("missing Server.ExtHost")
var errMissingExtPort = errors.New("missing Server.ExtPort")

// LocalEntry works like Entry but uses ExtHost and ExtPort from the server
// instead of the ones supplied in e.
func (dw *DirWriter) LocalEntry(e *DirEntry) error {
	e.Host = dw.w.srv.ExtHost
	if e.Host == "" {
		return errMissingExtHost
	}
	e.Port = dw.w.srv.ExtPort
	if e.Port == "" {
		return errMissingExtPort
	}
	return dw.Entry(e)
}
