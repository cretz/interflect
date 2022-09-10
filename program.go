package interflect

import (
	"fmt"
	"go/token"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
)

type Program struct {
	ssa *ssa.Program
}

type ProgramConfig struct {
	// Dependencies are loaded as needed, so only entry package is needed. All
	// packages must have the same Fset.
	Packages []*packages.Package
	// Required and entrypoint must be interpreted. Default for any not matched is
	// PackageNotInterpreted
	PackageModeChecks []PackageModeCheck
}

type PackageModeCheck func(*packages.Package) PackageMode

type PackageMode int

const (
	PackageModeUnknown PackageMode = iota
	// Always interpreted, cannot be imported by non-interpreted and cannot import
	// disallowed
	PackageModeInterpreted
	// Never interpreted, cannot import interpreted
	PackageModeNotInterpreted
	// Cannot be imported by interpreted packages
	PackageModeDisallowed
)

func InterpretPackageByName(name string) PackageModeCheck {
	return func(p *packages.Package) PackageMode {
		if p.PkgPath == name {
			return PackageModeInterpreted
		}
		return PackageModeUnknown
	}
}

func NewProgram(config ProgramConfig) (*Program, error) {
	// Get file set
	var fset *token.FileSet
	for _, pkg := range config.Packages {
		if fset != nil && fset != pkg.Fset {
			return nil, fmt.Errorf("all packages must have same file set")
		}
		fset = pkg.Fset
	}
	if fset == nil {
		return nil, fmt.Errorf("no file set packages")
	}

	// Collect modes
	pkgModes := map[*packages.Package]PackageMode{}
	packages.Visit(config.Packages, nil, func(p *packages.Package) {
		if _, exists := pkgModes[p]; !exists {
			// Find first non-unknown
			for _, check := range config.PackageModeChecks {
				if mode := check(p); mode > PackageModeUnknown {
					pkgModes[p] = mode
					return
				}
			}
			pkgModes[p] = PackageModeNotInterpreted
		}
	})

	// Validate and load SSA packages
	// TODO(cretz): Accept builder mode via config? Default it to something different?
	ssaProg := ssa.NewProgram(fset, 0)
	for pkg, pkgMode := range pkgModes {
		// Must be properly typed
		if pkg.Types == nil || pkg.IllTyped {
			return nil, fmt.Errorf("interpreted package %q is not properly typed", pkg.PkgPath)
		}
		// Create SSA package
		ssaProg.CreatePackage(pkg.Types, pkg.Syntax, pkg.TypesInfo, true)
		switch pkgMode {
		case PackageModeInterpreted:
			// Do not allow disallowed packages
			// TODO(cretz): Also disallow import "C" for interpreted packages?
			for _, depPkg := range pkg.Imports {
				if pkgModes[depPkg] == PackageModeDisallowed {
					return nil, fmt.Errorf("interpreted package %q depends on %q which is configured as disallowed",
						pkg.PkgPath, depPkg.PkgPath)
				}
			}
		case PackageModeNotInterpreted:
			// Do not allow interpreted packages to be imported by non-interpreted
			for _, depPkg := range pkg.Imports {
				if pkgModes[depPkg] == PackageModeInterpreted {
					return nil, fmt.Errorf("non-interpreted package %q depends on %q which is configured as interpreted",
						pkg.PkgPath, depPkg.PkgPath)
				}
			}
		}
	}
	ssaProg.Build()
	return &Program{ssa: ssaProg}, nil
}
