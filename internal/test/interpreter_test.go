package test

import (
	"reflect"
	"testing"

	"github.com/cretz/interflect"
	"github.com/cretz/interflect/internal/test/interpreted"
	"github.com/stretchr/testify/require"
)

func TestSimple(t *testing.T) {
	// Make interpreter
	inter, err := interflect.New(interflect.InterpreterConfig{})
	require.NoError(t, err)
	prog, err := inter.PrepareProgram((interflect.ProgramConfig{
		Packages: []interflect.PackageFilter{
			interflect.InterpretPackageByName("github.com/cretz/interflect/internal/test/interpreted"),
		},
	}))
	require.NoError(t, err)

	// Get function from interpreted package
	fn := prog.ReflectFunc(interpreted.SayHello)
	require.True(t, fn.IsValid())

	// Run
	res := prog.ScheduleCall(fn, []reflect.Value{reflect.ValueOf("world")})
	prog.Run()
	require.Equal(t, "Hello, world!", (<-res)[0].String())
}
