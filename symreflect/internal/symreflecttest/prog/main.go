package main

import (
	"log"
	"os"

	"github.com/cretz/interflect/symreflect/internal/symreflecttest/testpkg"
	"github.com/stretchr/testify/require"
)

func main() {
	req := require.New(testingT{})
	req.Equal(2, len(os.Args))
	switch os.Args[1] {
	case "funcs":
		testpkg.TestFuncs(req)
	case "symbols":
		testpkg.TestSymbols(req, true)
	case "types":
		testpkg.TestTypes(req)
	default:
		req.Fail("unknown command")
	}
}

type testingT struct{}

func (testingT) Errorf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

func (testingT) FailNow() {
	log.Fatal()
}
