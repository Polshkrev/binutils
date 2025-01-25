package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Polshkrev/gopolutils"
	"github.com/Polshkrev/gopolutils/collections"
)

// Concurrently look for the given executable in the system path.
func sendRequest(programme string, path chan<- string, except chan<- error) {
	var result string
	var err error
	result, err = exec.LookPath(programme)
	path <- result
	except <- err
}

// Get the system path value of a given string executable.
// Returns the path located in the system path.
// If the given executable's path can not be obtained, an Exception is returned with an empty string.
func which(programme string) (string, *gopolutils.Exception) {
	var pathChannel chan string = make(chan string, 1)
	var exceptChannel chan error = make(chan error, 1)
	defer close(pathChannel)
	defer close(exceptChannel)
	go sendRequest(programme, pathChannel, exceptChannel)
	var path string = <-pathChannel
	var err error = <-exceptChannel
	if err != nil {
		return "", gopolutils.NewException(fmt.Sprintf("executable file '%s' not found in system path.", programme))
	}
	return path, nil
}

// Obtain the system path of a variadic list of executable names.
// Returns a collection of systems paths from the given list of executable names.
// If one of the executable paths can not be obtained, an Exception is returned with a nil data pointer.
func Which(programmes ...string) (collections.View[string], *gopolutils.Exception) {
	var programme string
	var results collections.Collection[string] = collections.NewArray[string]()
	for _, programme = range programmes {
		var result string
		var except *gopolutils.Exception
		result, except = which(programme)
		if except != nil {
			return nil, except
		}
		results.Append(result)
	}
	return results, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, gopolutils.NewException("No arguments provided.").Error())
		os.Exit(1)
	}
	var paths collections.View[string]
	var except *gopolutils.Exception
	paths, except = Which(os.Args[1:]...)
	if except != nil {
		fmt.Fprintln(os.Stderr, except.Error())
		os.Exit(1)
	}
	var path string
	for _, path = range paths.Collect() {
		fmt.Println(path)
	}
}
