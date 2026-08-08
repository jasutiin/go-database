package main

import (
	"fmt"

	"github.com/jasutiin/go-database/engine"
)

func main() {
	fmt.Println("Hello, World!")
	engine := engine.Engine{}
	engine.Startup()
}
