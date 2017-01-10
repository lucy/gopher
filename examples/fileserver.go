package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/lucy/gopher"
)

var addr = flag.String("addr", "127.0.0.1:7070", "address to listen to")
var host = flag.String("host", "127.0.0.1", "external host name to present to clients")
var port = flag.String("port", "7070", "external port to present to clients")
var path = flag.String("path", "", "path to serve files from (required)")

func main() {
	flag.Parse()
	if *path == "" {
		log.Fatal("missing option: -path")
		flag.Usage()
		os.Exit(1)
	}
	srv := &gopher.Server{
		Addr:    *addr,
		ExtHost: *host,
		ExtPort: *port,
		Handler: gopher.FileServer(http.Dir(*path)),
	}
	log.Printf("Listening on %s", *addr)
	log.Fatal(srv.ListenAndServe())
}
