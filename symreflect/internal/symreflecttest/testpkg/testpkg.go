package testpkg

import (
	"fmt"
	"reflect"

	"github.com/cretz/interflect/symreflect"
	"github.com/stretchr/testify/require"
)

//go:noinline
func sayHello(name string) string {
	return fmt.Sprintf("Hello, %v!", name)
}

type myStruct struct{}

var topLevelVar = "foo"

func TestFuncs(req *require.Assertions) {
	// Load types first
	types, err := symreflect.LoadTypes(symreflect.LoadTypesOptions{})
	req.NoError(err)

	// Load functions and confirm present
	funcs, err := symreflect.LoadFuncs(symreflect.LoadFuncsOptions{
		ModuleDataType: types.Named["runtime"]["moduledata"],
	})
	req.NoError(err)
	fn := funcs.Named["github.com/cretz/interflect/symreflect/internal/symreflecttest/testpkg"]["sayHello"]
	req.NotNil(fn)

	// Invoke
	expected := sayHello("World")
	reflected := symreflect.ReflectFunc(reflect.TypeOf(sayHello), fn.Entry())
	actual := reflected.Call([]reflect.Value{reflect.ValueOf("World")})[0].Interface()
	req.Equal(expected, actual)
}

func TestSymbols(req *require.Assertions, expectPresent bool) {
	// Tell the compiler the var is used
	req.NotEmpty(topLevelVar)

	// Load symbols
	symbols, err := symreflect.LoadSymbols(symreflect.LoadSymbolsOptions{})
	if !expectPresent {
		req.Equal(symreflect.ErrNoSymbolTable, err)
		return
	}
	req.NoError(err)

	// Check top level var
	varSymbol := symbols.Named["github.com/cretz/interflect/symreflect/internal/symreflecttest/testpkg"]["topLevelVar"]
	req.NotNil(varSymbol)
	reflectVar := symbols.ReflectValue(reflect.TypeOf(""), varSymbol)
	req.Equal("foo", reflectVar.Elem().Interface())

	// Check top level func
	funcSymbol := symbols.Named["github.com/cretz/interflect/symreflect/internal/symreflecttest/testpkg"]["sayHello"]
	expected := sayHello("World")
	reflectFunc := symbols.ReflectValue(reflect.TypeOf(sayHello), funcSymbol)
	actual := reflectFunc.Call([]reflect.Value{reflect.ValueOf("World")})[0].Interface()
	req.Equal(expected, actual)
}

func TestTypes(req *require.Assertions) {
	// Just tell the compiler that we use these structs
	req.NotNil(&myStruct{})
	req.NotNil(&struct{ _symreflect_test string }{})

	// Check the named
	types, err := symreflect.LoadTypes(symreflect.LoadTypesOptions{})
	req.NoError(err)
	req.NotNil(types.Named["github.com/cretz/interflect/symreflect/internal/symreflecttest/testpkg"]["myStruct"])

	// Check the anon
	found := false
	for _, typ := range types.AnonStructs {
		if typ.NumField() == 1 && typ.Field(0).Name == "_symreflect_test" {
			found = true
			break
		}
	}
	req.True(found)
}
