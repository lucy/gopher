package gopher

import (
	"log"
	"net/http"

	"github.com/lucy/gopher"
)

func ExampleFileServer() {
	srv := &gopher.Server{Addr: "127.0.0.1:7070",
		ExtHost: "127.0.0.1", ExtPort: "7070",
		Handler: gopher.FileServer(http.Dir("/usr/share/doc"))}
	log.Fatal(srv.ListenAndServe())
}
