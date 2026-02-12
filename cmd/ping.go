package main

import (
	"fmt"
	"os"

	"github.com/Polshkrev/binutils"
	"github.com/Polshkrev/gopolutils"
	"github.com/Polshkrev/gopolutils/collections/safe"
	"github.com/Polshkrev/otvet"
)

func main() {
	if len(os.Args) < 1 {
		fmt.Fprintln(os.Stderr, gopolutils.NewException("No URLs provided.").Error())
		os.Exit(1)
	}
	var statusCodes collections.View[otvet.StatusCode] = gopolutils.Must(binutils.Ping(os.Args[1:]...))
	var i int
	for i = range statusCodes.Collect() {
		var code otvet.StatusCode = statusCodes.Collect()[i]
		fmt.Println(gopolutils.Must(otvet.StatusText(code)))
	}
}
