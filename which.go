package binutils

import (
	"fmt"
	"os/exec"

	"github.com/Polshkrev/gopolutils"
	"github.com/Polshkrev/gopolutils/collections"
)

// Get the system path value of a given string executable.
// Returns the path located in the system path.
// If the given executable's path can not be obtained, an Exception is returned with an empty string.
func which(programme string, path chan<- string, except chan<- *gopolutils.Exception) {
	var result string
	var err error
	result, err = exec.LookPath(programme)
	if err != nil {
		except <- gopolutils.NewException(fmt.Sprintf("executable file '%s' not found in system path.", programme))
		return
	}
	path <- result
	except <- nil
}

// Obtain the system path of a variadic list of executable names.
// Returns a collection of systems paths from the given list of executable names.
// If one of the executable paths can not be obtained, an Exception is returned with a nil data pointer.
func Which(programmes ...string) (collections.View[string], *gopolutils.Exception) {
	var programme string
	var pathChannel chan string = make(chan string)
	var exceptChannel chan *gopolutils.Exception = make(chan *gopolutils.Exception)
	var results collections.Collection[string] = collections.NewArray[string]()
	for _, programme = range programmes {
		defer close(pathChannel)
		defer close(exceptChannel)
		go which(programme, pathChannel, exceptChannel)
		var result string = <-pathChannel
		var except *gopolutils.Exception = <-exceptChannel
		if except != nil {
			return nil, except
		}
		results.Append(result)
	}
	// close(pathChannel)
	// close(exceptChannel)
	return results, nil
}
