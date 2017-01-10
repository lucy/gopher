package main

import (
	"log"

	"github.com/lucy/gopher"
)

var host = "127.0.0.1"
var port = "7070"

func main() {
	addr := host + ":" + port
	srv := &gopher.Server{
		Addr:    addr,
		ExtHost: host,
		ExtPort: port,
		Handler: gopher.NewFileServer("./gopher"),
	}
	log.Printf("Listening on %s", addr)
	log.Fatal(srv.ListenAndServe())
}
