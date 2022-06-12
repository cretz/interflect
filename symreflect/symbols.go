package symreflect

import (
	"errors"
	"reflect"
)

type LoadSymbolsOptions struct {
	// Default is os.Executable().
	ExePath string
}

type Symbols struct {
	Named             map[string]map[string]Symbol
	ProcessAddrOffset int64
}

type Symbol struct {
	PackageName string
	Name        string
	Addr        uintptr
	Type        SymbolType
}

type SymbolType uint8

const (
	SymbolTypeUnknown SymbolType = iota
	SymbolTypeText
	SymbolTypeData
	SymbolTypeDataReadOnly
	SymbolTypeBSS
	SymbolTypeConst
	SymbolTypeUndef
)

// ErrNoSymbolTable is returned from LoadSymbols when the symbol table has been
// stripped from the executable. This can occur with linker flag -s, when using
// "go run", or most of the time when using "go test".
var ErrNoSymbolTable = errors.New("no symbol table")

// Both results empty if unknown
func parseSymbolName(raw string) (pkgName, name string) {
	panic("TODO")
}

func (s *Symbol) ReflectValue(typ reflect.Type) reflect.Value {
	// TODO(cretz): Can use reflect.NewAt for vars, but for functions we need to
	// use ReflectFunc
	panic("TODO")
}
