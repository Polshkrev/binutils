package tests

import (
	"testing"

	"github.com/Polshkrev/gopolutils"
)

func TestSinglePingSucess(test *testing.T) {
	var path string = "https://update.code.visualstudio.com/1.99.0/win32-x64-archive/stable"
	var except *gopolutils.Exception
	_, except = binutils.Ping(path)
	if except != nil {
		test.Errorf("%s\n", except.Error())
	}
}

func TestSinglePingFailure(test *testing.T) {
	var path string = "asd"
	var except *gopolutils.Exception
	_, except = binutils.Ping(path)
	if except == nil {
		test.Errorf("The single test case has returned a nil error.\n")
	}
}
