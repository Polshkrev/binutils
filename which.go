package binutils

import (
	"fmt"
	"os/exec"

	"github.com/Polshkrev/gopolutils"
	"github.com/Polshkrev/gopolutils/collections"
	"github.com/Polshkrev/gopolutils/collections/safe"
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
// If the given executable's path can not be obtained, an [gopolutils.Exception] is returned with an empty string.
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
// If one of the executable paths can not be obtained, an [gopolutils.Exception] is returned with a nil data pointer.
func Which(programmes ...string) (collections.View[string], *gopolutils.Exception) {
	var results safe.Collection[string] = safe.NewArray[string]()
	var i int
	for i = range programmes {
		var programme string = programmes[i] // Idiomatically more performant than direct access.
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
