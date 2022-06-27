package symreflect

import (
	"errors"
	"reflect"
	"strings"
	"unsafe"
)

// LoadSymbolsOptions are options for LoadSymbols.
type LoadSymbolsOptions struct {
	// Default is os.Executable().
	ExePath string
}

// Symbols contains the set of symbols
type Symbols struct {
	// First key is Symbol.PackageName, second key is Symbol.Name (which may have
	// a receiver prepended).
	Named map[string]map[string]*Symbol
	// Offset to add to Symbol.Addr to get the actual symbol address.
	ProcessAddrOffset int64
}

// Symbol represents a symbol in the symbol table.
type Symbol struct {
	// Full package name.
	PackageName string
	// Name which may have the receiver prepended (surrounded be parentheses for
	// pointer receivers).
	Name string
	// Address of the symbol, which must be combined with
	// Symbols.ProcessAddrOffset to get the real address.
	Addr uintptr
	Type SymbolType
}

// SymbolType is the type of symbol.
type SymbolType uint8

// Known symbol types.
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
	// Must be in form of <pkg>.[<type>.]<name>. We combine receiver type and name
	// into name here. Some restrictions below.

	// Cannot have comma or dollar sign
	if strings.ContainsAny(raw, ",$") {
		return
	}
	// Cannot end with asterisk
	if strings.HasSuffix(raw, "*") {
		return
	}

	// Take everything after the first got after the last slash
	pkgPathEndIndex := strings.LastIndex(raw, "/")
	pkgPathEndIndex += 1 + strings.Index(raw[pkgPathEndIndex+1:], ".")
	if pkgPathEndIndex < 0 {
		return
	}

	// Name cannot start with dot
	if pkgPathEndIndex+1 >= len(raw) || raw[pkgPathEndIndex+1] == '.' {
		return
	}

	return raw[:pkgPathEndIndex], raw[pkgPathEndIndex+1:]
}

// ReflectValue returns either a pointer to the package var, a function that can
// call the package function, or an invalid value if unknown symbol type.
func (s *Symbols) ReflectValue(typ reflect.Type, sym *Symbol) reflect.Value {
	switch sym.Type {
	case SymbolTypeText:
		return ReflectFunc(typ, sym.Addr+uintptr(s.ProcessAddrOffset))
	case SymbolTypeData:
		return reflect.NewAt(typ, unsafe.Pointer(sym.Addr+uintptr(s.ProcessAddrOffset)))
	default:
		return reflect.Value{}
	}
}
