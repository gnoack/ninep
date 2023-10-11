package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gnoack/ninep"
)

var (
	addr    = flag.String("addr", "localhost:8080", "Address to serve HTTP on")
	service = flag.String("service", "sources", "9p service to connect to")
)

func main() {
	fsys, err := ninep.DialFS(*service, ninep.DialFSOpts{})
	if err != nil {
		log.Fatalf("ninep.DialFS(%q): %v", *service, err)
	}

	http.Handle("/", http.FileServer(http.FS(fsys)))
	err = http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatalf("http.ListenAndServe(%q, nil): %v", *addr, err)
	}
}
