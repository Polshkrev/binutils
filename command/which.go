package main

import (
	"fmt"
	"os"

	"github.com/Polshkrev/binutils"
	"github.com/Polshkrev/gopolutils"
	"github.com/Polshkrev/gopolutils/collections"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, gopolutils.NewException("No arguments provided.").Error())
		os.Exit(1)
	}
	var paths collections.View[string]
	var except *gopolutils.Exception
	paths, except = binutils.Which(os.Args[1:]...)
	if except != nil {
		fmt.Fprintln(os.Stderr, except.Error())
		os.Exit(1)
	}
	var path string
	for _, path = range paths.Collect() {
		fmt.Println(path)
	}
}
