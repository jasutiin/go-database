package main

import (
	"fmt"

	"github.com/jasutiin/go-database/engine"
)

func main() {
	db, err := engine.Startup()

	if err != nil {
		fmt.Println(err.Error())
	}

	db.Get([]byte("hey"))
	db.Put([]byte("hey"), []byte("what's up"))
}
