package main

import (
	"fmt"
	"os"

	"github.com/Polshkrev/binutils"
	"github.com/Polshkrev/gopolutils"
	"github.com/Polshkrev/gopolutils/collections"
)

func main() {
	if len(os.Args) < 1 {
		fmt.Fprintln(os.Stderr, gopolutils.NewException("No URLs provided.").Error())
		os.Exit(1)
	}
	var statusCodes collections.Collection[uint16] = gopolutils.Must(binutils.Ping(os.Args[1:]...))
	var code uint16
	for _, code = range statusCodes.Collect() {
		fmt.Println(code)
	}
}
