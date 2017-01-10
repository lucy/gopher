package gopher

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"
)

type fileServer struct {
	http.FileSystem
}

func NewFileServer(path string) Handler {
	return &fileServer{http.Dir("./gopher")}
}

func (fs *fileServer) ServeGopher(w *Writer, req *Request) {
	err := fs.serve(w, req)
	if err != nil {
		log.Print(err)
	}
}

func (fs *fileServer) serve(w *Writer, req *Request) error {
	p := path.Clean(string(req.Content))
	f, err := fs.Open(p)
	if err != nil {
		return err
	}
	fi, err := f.Stat()
	if err != nil {
		return err
	}
	if fi.IsDir() {
		fis, err := f.Readdir(-1)
		if err != nil {
			return err
		}
		sort.Sort(byName(fis))
		w := w.DirWriter()
		e := DirEntry{}
		for _, fi := range fis {
			e.Type = itemType(fi)
			e.Name = fi.Name()
			e.Path = path.Join(p, fi.Name())
			err := w.LocalEntry(&e)
			if err != nil {
				return err
			}
		}
		return w.Close()
	}
	_, err = io.Copy(w.Conn, f)
	return err
}

// TODO: make this less ugly
func itemType(fi os.FileInfo) byte {
	if fi.IsDir() {
		return '1'
	}
	n := fi.Name()
	switch {
	case strings.HasSuffix(n, ".html"):
		return 'h'
	case strings.HasSuffix(n, ".txt"):
		return '0'
	case strings.HasSuffix(n, ".gif"):
		return 'g'
	case strings.HasSuffix(n, ".png"),
		strings.HasSuffix(n, ".jpg"),
		strings.HasSuffix(n, ".jpeg"):
		return 'I'
	}
	return '9'
}

type byName []os.FileInfo

func (s byName) Len() int           { return len(s) }
func (s byName) Less(i, j int) bool { return s[i].Name() < s[j].Name() }
func (s byName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
