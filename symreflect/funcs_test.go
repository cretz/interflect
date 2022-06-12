package symreflect_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/cretz/interflect/symreflect"
	"github.com/stretchr/testify/require"
)

//go:noinline
func sayHello(name string) string {
	return fmt.Sprintf("Hello, %v!", name)
}

func TestLoadFuncsInTest(t *testing.T) {
	types, err := symreflect.LoadTypes(symreflect.LoadTypesOptions{})
	require.NoError(t, err)
	funcs, err := symreflect.LoadFuncs(symreflect.LoadFuncsOptions{
		ModuleDataType: types.Named["runtime"]["moduledata"],
	})
	require.NoError(t, err)
	fn := funcs.Named["github.com/cretz/interflect/symreflect_test"]["sayHello"]
	require.NotNil(t, fn)
	expected := sayHello("World")
	reflected := symreflect.ReflectFunc(reflect.TypeOf(sayHello), fn.Entry())
	actual := reflected.Call([]reflect.Value{reflect.ValueOf("World")})[0].Interface()
	require.Equal(t, expected, actual)
}

func TestLoadFuncsInProg(t *testing.T) {
	// TODO(cretz): Do standard Go build of prog that does what TestLoadFuncsInTest does
	t.Skip("TODO")
}
