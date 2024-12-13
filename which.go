package binutils

import (
	"fmt"
	"os/exec"

	"github.com/Polshkrev/gopolutils"
	"github.com/Polshkrev/gopolutils/collections"
)

func Which(programmes ...string) (collections.Collection[string], *gopolutils.Exception) {
	var programme string
	var results collections.Collection[string] = collections.NewArray[string]()
	for _, programme = range programmes {
		var result string
		var err error
		result, err = exec.LookPath(programme)
		if err != nil {
			return nil, gopolutils.NewException(fmt.Sprintf("executable file '%s' not found in system path.", programme))
		}
		results.Append(result)
	}
	return results, nil
}
