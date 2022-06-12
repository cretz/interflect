package symreflect_test

import (
	"testing"

	"github.com/cretz/interflect/symreflect"
	"github.com/stretchr/testify/require"
)

func TestLoadSymbolsInTest(t *testing.T) {
	// Symbols are elided in "go test" builds
	_, err := symreflect.LoadSymbols(symreflect.LoadSymbolsOptions{})
	require.Equal(t, symreflect.ErrNoSymbolTable, err)
}

func TestLoadSymbolsInProg(t *testing.T) {
	// TODO(cretz): Do standard Go build of prog that does what TestLoadSymbolsInTest does
	t.Skip("TODO")
}
