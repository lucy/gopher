package gopher

import (
	"bufio"
	"log"
	"net"
	"net/textproto"
	"time"
)

type Server struct {
	Addr         string
	ExtHost      string
	ExtPort      string
	Handler      Handler
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	MaxReqBytes  int
	ErrorLog     *log.Logger
}

type Request struct {
	RemoteAddr string
	Content    []byte
}

type Writer struct {
	Conn net.Conn
	srv  *Server
}

func (w *Writer) DirWriter() *DirWriter {
	return &DirWriter{w, textproto.NewWriter(bufio.NewWriter(w.Conn)).DotWriter()}
}

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
