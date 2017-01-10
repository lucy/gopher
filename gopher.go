package gopher

import (
	"bufio"
	"log"
	"net"
	"net/textproto"
	"time"
)

// A Server defines parameters for running a gopher server.
type Server struct {
	Addr         string        // TCP address to listen on, ":7070" if empty
	ExtHost      string        // External address for this server (required for DirWriter.LocalEntry)
	ExtPort      string        // External port for this server (required for DirWriter.LocalEntry)
	Handler      Handler       // handler to invoke
	ReadTimeout  time.Duration // maximum duration before timing out read of the request<Paste>
	WriteTimeout time.Duration // maximum duration before timing out write of the response<Paste>
	MaxReqBytes  int           // maximum length for a request
	ErrorLog     *log.Logger   // optional logger for errors
}

// A Writer is used for writing responses to requests.
type Writer struct {
	Conn net.Conn
	srv  *Server
}

// DirWriter returns a DirWriter for a Writer.
func (w *Writer) DirWriter() *DirWriter {
	return &DirWriter{w, textproto.NewWriter(bufio.NewWriter(w.Conn)).DotWriter()}
}

// A Request represents a request received by a server.
type Request struct {
	RemoteAddr string
	Content    []byte
}

// A Handler responds to a request.
type Handler interface {
	ServeGopher(w *Writer, request *Request)
}

func (srv *Server) logf(format string, args ...interface{}) {
	if srv.ErrorLog != nil {
		srv.ErrorLog.Printf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}

func (srv *Server) serve(c net.Conn) {
	defer c.Close()
	req := &Request{}
	req.RemoteAddr = c.RemoteAddr().String()
	r := bufio.NewReader(c)
	line, err := r.ReadBytes('\n')
	// trim \n, \r\n
	if l := len(line); l >= 1 && line[l-1] == '\n' {
		line = line[:l-1]
	}
	if l := len(line); l >= 1 && line[l-1] == '\r' {
		line = line[:l-1]
	}
	if err != nil {
		srv.logf("Error reading request from %s: %s", req.RemoteAddr, err)
		return
	}
	req.Content = line
	w := &Writer{c, srv}
	srv.Handler.ServeGopher(w, req)
}

// Serve accepts incoming connections on the listener l, creating a new service
// goroutine for each. The service goroutines read a request and then call
// srv.Handler to reply to them.
func (srv *Server) Serve(l net.Listener) error {
	defer l.Close()
	for {
		c, err := l.Accept()
		if err != nil {
			return err
		}
		go srv.serve(c)
	}
}

// ListenAndServe listens on the TCP network address srv.Addr and calls Serve
// to handle requests.
func (srv *Server) ListenAndServe() error {
	addr := srv.Addr
	if addr == "" {
		addr = "127.0.0.1:7070"
	}
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return srv.Serve(l)
}