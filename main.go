package main

import (
	"fmt"
	"os"

	"github.com/jasutiin/go-database/engine"
)

func main() {
	opts := &engine.Options{
		DbName:           "name",
		StorageType:      engine.StorageLSM,
		SkipListMaxLevel: 10,
		SkipListMaxSize:  1000,
	}
	db, err := engine.Startup(opts)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = db.Get([]byte("hey"))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = db.Put([]byte("hey"), []byte("what's up"))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
