package binutils

import (
	"fmt"
	"os/exec"

	"github.com/Polshkrev/gopolutils"
	"github.com/Polshkrev/gopolutils/collections"
)

func which(programme string) (string, *gopolutils.Exception) {
	var result string
	var err error
	result, err = exec.LookPath(programme)
	if err != nil {
		return "", gopolutils.NewException(fmt.Sprintf("executable file '%s' not found in system path.", programme))
	}
	return result, nil
}

func Which(programmes ...string) (collections.Collection[string], *gopolutils.Exception) {
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
