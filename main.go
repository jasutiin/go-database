package main

import (
	"fmt"
	"os"

	"github.com/jasutiin/go-database/storage"
)

func main() {
	opts := &storage.Options{
		DbName:           "name",
		StorageType:      storage.StorageLSM,
		SkipListMaxLevel: 10,
		SkipListMaxSize:  1000,
	}
	store, err := storage.Startup(opts)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	_, err = store.Get([]byte("hey"))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = store.Put([]byte("hey"), []byte("what's up"))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
