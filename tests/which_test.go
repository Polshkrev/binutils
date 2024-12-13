package tests

import (
	"testing"

	"github.com/Polshkrev/binutils"
	"github.com/Polshkrev/gopolutils"
)

func TestSingleWhichSucess(test *testing.T) {
	var path string = "which"
	var except *gopolutils.Exception
	_, except = binutils.Which(path)
	if except != nil {
		test.Errorf("%s\n", except.Error())
	}
}

func TestSingleWhichFailure(test *testing.T) {
	var path string = "asd"
	var except *gopolutils.Exception
	_, except = binutils.Which(path)
	if except == nil {
		test.Errorf("%s\n", except.Error())
	}
}
