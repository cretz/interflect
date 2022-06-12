package symreflect_test

import (
	"testing"

	"github.com/cretz/interflect/symreflect"
	"github.com/stretchr/testify/require"
)

type myStruct struct{}

func TestLoadTypesInTest(t *testing.T) {
	// Just tells the compiler that we use these struct
	require.NotNil(t, &myStruct{})
	require.NotNil(t, &struct{ _symflect_test string }{})

	// Check the named
	types, err := symreflect.LoadTypes(symreflect.LoadTypesOptions{})
	require.NoError(t, err)
	require.NotNil(t, types.Named["github.com/cretz/interflect/symreflect_test"]["myStruct"])

	// Check the anon
	found := false
	for _, typ := range types.AnonStructs {
		if typ.NumField() == 1 && typ.Field(0).Name == "_symflect_test" {
			found = true
			break
		}
	}
	require.True(t, found)
}

func TestLoadTypesInProg(t *testing.T) {
	// TODO(cretz): Do standard Go build of prog that does what TestLoadTypesInTest does
	t.Skip("TODO")
}
