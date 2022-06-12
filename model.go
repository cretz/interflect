package interflect

import (
	"reflect"
	"sync"

	"golang.org/x/tools/go/ssa"
)

type execPackage struct {
	ssa  *ssa.Package
	init sync.Once
	// Top-level functions (including methods) and vars. Vars are always pointers,
	// even if they are pointers to pointers.
	values map[string]reflect.Value
}

func newExecPackage(exec *Execution, ssaPkg *ssa.Package) *execPackage {
	p := &execPackage{ssa: ssaPkg, values: map[string]reflect.Value{}}
	for name, m := range p.ssa.Members {
		switch m := m.(type) {
		case *ssa.Function:
			// Get the existing type
			fnType := exec.conf.PackageReflector.ReflectValue(p.ssa.Pkg.Path(), name).Type()
			// Create an exec function and make a func out of it
			fn := &execFunction{ssa: m, exec: exec}
			p.values[name] = reflect.MakeFunc(fnType, fn.call)
		case *ssa.Global:
			// Create a new pointer to the value type for the global. We have to build
			// this type instead of reference the existing type because we want the
			// declared type not the current value's type.
			p.values[name] = reflect.New(exec.reflectType(m.Type()))
		}
	}
	return p
}

type execFunction struct {
	ssa  *ssa.Function
	exec *Execution
}

func (e *execFunction) call(args []reflect.Value) []reflect.Value {
	panic("TODO")
}
