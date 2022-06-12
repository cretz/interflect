package symreflect

import (
	"debug/pe"
	"fmt"
	"os"
	"syscall"
)

var kernel32 *syscall.LazyDLL
var getModuleHandle *syscall.LazyProc

func init() {
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	getModuleHandle = kernel32.NewProc("GetModuleHandleA")
}

func LoadSymbols(options LoadSymbolsOptions) (*Symbols, error) {
	// Inspired by Go's nm CLI

	syms := &Symbols{Named: map[string]map[string]Symbol{}}

	// Set process offset
	addrRaw, _, _ := getModuleHandle.Call(0)
	syms.ProcessAddrOffset = int64(addrRaw) - 0x400000

	// Open exe
	if options.ExePath == "" {
		var err error
		if options.ExePath, err = os.Executable(); err != nil {
			return nil, fmt.Errorf("cannot find exe path: %w", err)
		}
	}
	r, err := os.Open(options.ExePath)
	if err != nil {
		return nil, fmt.Errorf("failed opening exe: %w", err)
	}
	defer r.Close()

	f, err := pe.NewFile(r)
	if err != nil {
		return nil, fmt.Errorf("failed opening PE: %w", err)
	} else if len(f.Symbols) == 0 {
		return nil, ErrNoSymbolTable
	}

	var imageBase uint64
	switch oh := f.OptionalHeader.(type) {
	case *pe.OptionalHeader32:
		imageBase = uint64(oh.ImageBase)
	case *pe.OptionalHeader64:
		imageBase = oh.ImageBase
	}

	for _, s := range f.Symbols {
		const (
			N_UNDEF = 0  // An undefined (extern) symbol
			N_ABS   = -1 // An absolute symbol (e_value is a constant, not an address)
			N_DEBUG = -2 // A debugging symbol
		)

		// Parse the name, and if unparseable do not include the symbol
		pkgName, name := parseSymbolName(s.Name)
		if pkgName == "" && name == "" {
			continue
		}

		sym := Symbol{PackageName: pkgName, Name: name, Addr: uintptr(s.Value)}
		switch s.SectionNumber {
		case N_UNDEF:
			sym.Type = SymbolTypeUndef
		case N_ABS:
			sym.Type = SymbolTypeConst
		case N_DEBUG:
		default:
			if s.SectionNumber < 0 || len(f.Sections) < int(s.SectionNumber) {
				return nil, fmt.Errorf("invalid section number in symbol table")
			}
			sect := f.Sections[s.SectionNumber-1]
			const (
				text  = 0x20
				data  = 0x40
				bss   = 0x80
				permW = 0x80000000
			)
			ch := sect.Characteristics
			switch {
			case ch&text != 0:
				sym.Type = SymbolTypeText
			case ch&data != 0:
				if ch&permW == 0 {
					sym.Type = SymbolTypeDataReadOnly
				} else {
					sym.Type = SymbolTypeData
				}
			case ch&bss != 0:
				sym.Type = SymbolTypeBSS
			}
			sym.Addr += uintptr(imageBase + uint64(sect.VirtualAddress))
		}
		pkgSyms := syms.Named[pkgName]
		if pkgSyms == nil {
			pkgSyms = map[string]Symbol{}
			syms.Named[pkgName] = pkgSyms
		}
		pkgSyms[name] = sym
	}
	return syms, nil
}
