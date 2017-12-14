package main

import (
	"os"
	"github.com/splunknova/nova-cli/src"
	"strings"
)

func main() {
	var searchKeywords string
	var args []string

	if len(os.Args) > 1 {
		searchKeywords, args = os.Args[1], os.Args[1:]
		if strings.HasPrefix(searchKeywords, "-") {
			searchKeywords, args = "", os.Args
		}
	} else {
		searchKeywords, args = "", os.Args
	}

	app := src.NewCLI(searchKeywords)
	app.Run(args)
}