package interflect_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/cretz/interflect"
	"github.com/cretz/interflect/internal/testdata/interpretedpackage1"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/packages"
)

var prog *interflect.Program

func init() {
	pkgs, err := packages.Load(&packages.Config{Mode: packages.LoadAllSyntax},
		"github.com/cretz/interflect/internal/testdata/interpretedpackage1")
	if err != nil {
		panic(err)
	} else if len(pkgs) != 1 {
		panic(fmt.Sprintf("expected 1 package, got %v", len(pkgs)))
	}
	prog, err = interflect.NewProgram(interflect.ProgramConfig{
		Packages: pkgs,
		PackageModeChecks: []interflect.PackageModeCheck{
			interflect.InterpretPackageByName(pkgs[0].PkgPath),
		},
	})
	if err != nil {
		panic(err)
	}
}

func TestSimple(t *testing.T) {
	// Create execution
	exec, err := interflect.NewExecution(interflect.ExecutionConfig{Program: prog})
	require.NoError(t, err)

	// Get function from interpreted package
	fn := exec.ReflectFunc(interpretedpackage1.SayHello)
	require.True(t, fn.IsValid())

	// Run it
	resCh := make(chan []reflect.Value, 1)
	go func() { resCh <- exec.Go(fn, []reflect.Value{reflect.ValueOf("world")}) }()
	exec.Run()
	require.Equal(t, "Hello, world!", (<-resCh)[0].String())
}
