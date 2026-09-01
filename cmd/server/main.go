package main

import (
	"flag"
	"log"

	"github.com/jasutiin/go-database/server"
)

func main() {
	address := flag.String("addr", ":8080", "HTTP server address")
	flag.Parse()

	if err := server.New(*address).ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
