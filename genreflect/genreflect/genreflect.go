package genreflect

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/tools/go/packages"
)

type GenerateReflectorConfig struct {
	OutPackage string
	// If empty, all included. Otherwise immediately excluded if not matched in
	// here. Applied before exclude/filter.
	Include []*regexp.Regexp
	// If empty, none excluded. Otherwise immediately excluded if matched in here.
	// Applied after include and before filter.
	Exclude []*regexp.Regexp
	// If present and false, immediately excluded. Applied after include/exclude.
	Filter   func(*packages.Package) bool
	Patterns []string
	Env      []string
}

func (g *GenerateReflectorConfig) include(pkg *packages.Package) bool {
	// If there's an include list and this isn't in, excluded
	if len(g.Include) > 0 {
		matched := false
		for _, include := range g.Include {
			if matched = include.MatchString(pkg.PkgPath); matched {
				break
			}
		}
		if !matched {
			return false
		}
	}
	// If matches any exclusion, excluded
	for _, exclude := range g.Exclude {
		if exclude.MatchString(pkg.PkgPath) {
			return false
		}
	}
	// Check filter if present
	if g.Filter != nil {
		return g.Filter(pkg)
	}
	return true
}

func GenerateReflector(config GenerateReflectorConfig) ([]byte, error) {
	// Load all packages
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedName | packages.NeedImports | packages.NeedDeps | packages.NeedSyntax,
		Env:  config.Env,
	}, config.Patterns...)
	if err != nil {
		return nil, err
	}

	// Filter packages in place then sort
	n := 0
	for _, pkg := range pkgs {
		if config.include(pkg) {
			pkgs[n] = pkg
			n++
		}
	}
	pkgs = pkgs[:n]
	if len(pkgs) == 0 {
		return nil, fmt.Errorf("no packages matched")
	}
	sort.Slice(pkgs, func(i, j int) bool { return pkgs[i].PkgPath < pkgs[j].PkgPath })

	// Go over each building large switches
	importAliases := importAliases{"reflect": "reflect"}
	var typeSwitch, valueSwitch string
PkgLoop:
	for _, pkg := range pkgs {
		// Skip internal
		for _, part := range strings.Split(pkg.PkgPath, "/") {
			if part == "internal" {
				continue PkgLoop
			}
		}
		typeCases, valueCases := loadCases(importAliases, pkg)
		if len(typeCases) > 0 {
			typeSwitch += fmt.Sprintf("case %q: switch topLevelName { %v }\n", pkg.PkgPath, strings.Join(typeCases, "\n"))
		}
		if len(valueCases) > 0 {
			valueSwitch += fmt.Sprintf("case %q: switch topLevelName { %v }\n", pkg.PkgPath, strings.Join(valueCases, "\n"))
		}
	}

	// Build code
	file := "// Code generated by genreflect. DO NOT EDIT.\n\n"
	file += "package " + config.OutPackage + "\n\n"
	file += importAliases.code() + "\n\n"
	file += "// Reflector implements interflect.PackageReflector\n"
	file += "func Reflector() reflector { return reflector{} }\n\n"
	file += "type reflector struct{}\n\n"
	// ReflectType
	file += "func (reflector) ReflectType(pkgName, topLevelName string) reflect.Type {\n"
	if typeSwitch != "" {
		file += "switch pkgName { " + typeSwitch + " }\n"
	}
	file += "return nil\n}\n\n"
	// ReflectValue
	file += "func (reflector) ReflectValue(pkgName, topLevelName string) reflect.Value {\n"
	if valueSwitch != "" {
		file += "switch pkgName { " + valueSwitch + " }\n"
	}
	file += "return reflect.Value{}\n}\n"
	return format.Source([]byte(file))
}

func loadCases(i importAliases, pkg *packages.Package) (typeCases, valueCases []string) {
	for _, file := range pkg.Syntax {
		for _, decl := range file.Decls {
			switch decl := decl.(type) {
			case *ast.FuncDecl:
				// Only exported, top-level functions
				if decl.Name.IsExported() && decl.Recv == nil {
					valueCases = append(valueCases,
						fmt.Sprintf("case %q: return reflect.ValueOf(%v.%v)", decl.Name, i.get(pkg), decl.Name))
				}
			case *ast.GenDecl:
				switch decl.Tok {
				// Exported types
				case token.TYPE:
					for _, spec := range decl.Specs {
						name := spec.(*ast.TypeSpec).Name
						if name.IsExported() {
							typeCases = append(typeCases,
								fmt.Sprintf("case %q: return reflect.TypeOf((*%v.%v)(nil)).Elem()", name, i.get(pkg), name))
						}
					}
				// Exported vars
				case token.VAR:
					for _, spec := range decl.Specs {
						for _, name := range spec.(*ast.ValueSpec).Names {
							if name.IsExported() {
								valueCases = append(valueCases,
									fmt.Sprintf("case %q: return reflect.ValueOf(&%v.%v)", name, i.get(pkg), name))
							}
						}
					}
				}
			}
		}
	}
	sort.Strings(typeCases)
	sort.Strings(valueCases)
	return
}

// Keyed by alias, value is path
type importAliases map[string]string

func (i importAliases) get(pkg *packages.Package) string {
	for n := 0; ; n++ {
		// Only append number if >= 1
		alias := pkg.Name
		if n > 0 {
			alias += strconv.Itoa(n)
		}
		// Check if not there or our alias
		if existingPkgPath, ok := i[alias]; !ok {
			// Not there, add and return
			i[alias] = pkg.PkgPath
			return alias
		} else if existingPkgPath == pkg.PkgPath {
			// Our alias, use
			return alias
		}
	}
}

func (i importAliases) code() string {
	ret := "import (\n"
	for alias, pkgPath := range i {
		if path.Base(pkgPath) == alias {
			alias = ""
		}
		ret += fmt.Sprintf("%v %q\n", alias, pkgPath)
	}
	return ret + ")"
}
