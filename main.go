package main

import (
	"fmt"

	"github.com/jasutiin/go-database/engine"
)

func main() {
	fmt.Println("Hello, World!")
	databaseEngine := engine.Engine{}
	databaseEngine.Startup()
}
