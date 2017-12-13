package main

import (
	"os"
	"github.com/splunknova/nova-cli/src"
)

func main() {
	var searchKeywords string
	var args []string

	if len(os.Args) > 1 {
		searchKeywords, args = os.Args[1], os.Args[1:]
	} else {
		searchKeywords, args = "", os.Args
	}

	app := src.Foo(searchKeywords)
	app.Run(args)
}