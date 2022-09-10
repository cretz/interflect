package interflect

import (
	"fmt"
	"reflect"
)

type Execution struct {
	conf     ExecutionConfig
	coros    []*coroutine
	packages map[string]*execPackage
}

type ExecutionConfig struct {
	// Required
	Program *Program
	// Default is reflect.Value.MapRange
	MapRanger MapRanger
	// Default is SystemScheduler
	Scheduler Scheduler
	// Value to use for vars or funcs (including methods) accessed. Only applies
	// if called by interpreted code. Key is qualified name (including receiver
	// for methods).
	Intercept map[string]reflect.Value
	// TODO(cretz): Explain default
	PackageReflector PackageReflector
}

type MapRanger func(reflect.Value) MapIter

type MapIter interface {
	Key() reflect.Value
	Value() reflect.Value
	Next() bool
}

type PackageReflector interface {
	ReflectType(pkgName, topLevelName string) reflect.Type
	// This is a pointer for vars and a pointer to a func for functions
	ReflectValue(pkgName, topLevelName string) reflect.Value
}

func NewExecution(conf ExecutionConfig) (*Execution, error) {
	if conf.Program == nil {
		return nil, fmt.Errorf("missing program")
	}
	exec := &Execution{conf: conf}
	if exec.conf.MapRanger == nil {
		exec.conf.MapRanger = func(v reflect.Value) MapIter { return v.MapRange() }
	}
	if exec.conf.Scheduler == nil {
		exec.conf.Scheduler = SystemScheduler
	}
	// Prepare packages
	ssaPackages := exec.conf.Program.ssa.AllPackages()
	exec.packages = make(map[string]*execPackage, len(ssaPackages))
	for _, ssaPackage := range ssaPackages {
		exec.packages[ssaPackage.Pkg.Path()] = newExecPackage(exec, ssaPackage)
	}
	return exec, nil
}

func (*Execution) ReflectFunc(fn interface{}) reflect.Value {
	// TODO(cretz): Extract function name and just use ReflectValue
	panic("TODO")
}

func (*Execution) ReflectValue(pkgName, topLevelName string) reflect.Value {
	panic("TODO")
}

func (e *Execution) Run() {
	for e.RunOnce() {
	}
}

func (e *Execution) RunOnce() (coroutinesRemain bool) {
	// Run until all yielded
	allYielded := false
	for !allYielded {
		allYielded = true
		// Filter done coros in place
		n := 0
		for _, coro := range e.coros {
			wasYielded, alive := coro.runOnce()
			if alive {
				e.coros[n] = coro
				n++
			}
			if !wasYielded {
				allYielded = false
			}
		}
		e.coros = e.coros[:n]
	}
	return len(e.coros) > 0
}

func (*Execution) Go(fn reflect.Value, args []reflect.Value) []reflect.Value {
	panic("TODO")
}

func (*Execution) execPackage(path string) *execPackage {
	// if
	panic("TODO")
}

type coroutine struct {
	sched Coroutine
	frame *frame
}

type frame struct {
}

func (c *coroutine) runOnce() (wasYielded, alive bool) {
	panic("TODO")
}

// func (c *coroutine) call(caller *frame, callPos token.Pos, fn *ssa.Function)
