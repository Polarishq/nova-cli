package main

import (
	"os"
	"github.com/splunknova/nova-cli/src"
)

func main() {
	searchKeywords, args := os.Args[1], os.Args[1:]
	app := src.Foo(searchKeywords)
	app.Run(args)
}