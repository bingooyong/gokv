package main

import (
	"embed"
	"flag"
	"github.com/bingooyong/gokv/server"
	"io/fs"
	"log"
	"net/http"
	"strings"
)

//go:embed html/*
var html embed.FS

var (
	nodes   = flag.String("nodes", "", "comma seperated list of nodes")
	address = flag.String("address", ":4001", "http host:port")

	// global server
	srv *server.Server
)

func main() {
	var members []string
	if len(*nodes) > 0 {
		members = strings.Split(*nodes, ",")
	}
	// create new server
	s, err := server.New(&server.Options{
		Members: members,
	})

	if err != nil {
		log.Fatal(err)
	}

	// set global server
	srv = s
	log.Printf("Local node %s\n", srv.Address())

	// set http handlers
	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/set", setHandler)
	http.HandleFunc("/del", delHandler)

	// extract the embedded html directory
	htmlContent, err := fs.Sub(html, "html")
	if err != nil {
		log.Fatal(err)
	}

	// serve the html directory by default
	http.Handle("/", http.FileServer(http.FS(htmlContent)))

	log.Printf("Listening on %s\n", *address)

	if err := http.ListenAndServe(*address, nil); err != nil {
		log.Fatal(err)
	}
}
