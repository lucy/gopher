package gopher

import (
	"io"
	"net/http"
	"os"
	"path"
	"sort"
)

type fileHandler struct {
	root http.FileSystem
}

// FileServer returns a handler that servers the contents of the file system
// rooted at root.
func FileServer(root http.FileSystem) Handler {
	return &fileHandler{root}
}

func (f *fileHandler) ServeGopher(w *Writer, req *Request) {
	err := serve(w, req, f.root)
	if err != nil {
		w.srv.logf("FileServer.ServeGopher: %s", err)
	}
}

func serve(rw *Writer, req *Request, fs http.FileSystem) error {
	p := path.Clean("/" + string(req.Content))
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
		w := rw.DirWriter()
		e := DirEntry{}
		for _, fi := range fis {
			e.Type = itemType(fi)
			e.Name = fi.Name()
			if fi.IsDir() {
				e.Name += "/"
			}
			e.Path = path.Join(p, e.Name)
			err := w.LocalEntry(&e)
			if err != nil {
				return err
			}
		}
		return w.Close()
	}
	_, err = io.Copy(rw, f)
	return err
}

func itemType(fi os.FileInfo) byte {
	if fi.IsDir() {
		return '1'
	}
	switch path.Ext(fi.Name()) {
	case ".html":
		return 'h'
	case ".txt":
		return '0'
	case ".gif":
		return 'g'
	case ".png", ".jpg", ".jpeg":
		return 'I'
	}
	return '9'
}

type byName []os.FileInfo

func (s byName) Len() int           { return len(s) }
func (s byName) Less(i, j int) bool { return s[i].Name() < s[j].Name() }
func (s byName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
