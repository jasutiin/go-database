package main

import (
	"fmt"

	"github.com/jasutiin/go-database/engine"
)

func main() {
	opts := &engine.Options{DbName: "name"}
	db, err := engine.Startup(opts)

	if err != nil {
		fmt.Println(err.Error())
	}

	db.Get([]byte("hey"))
	db.Put([]byte("hey"), []byte("what's up"))
}
