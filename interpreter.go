package interflect

import (
	"reflect"

	"golang.org/x/tools/go/packages"
)

func New(InterpreterConfig) (*Interpreter, error) {
	panic("TODO")
}

type Interpreter struct {
}

type InterpreterConfig struct {
}

func (*Interpreter) PrepareProgram(ProgramConfig) (*Program, error) {
	panic("TODO")
}

type Program struct {
}

type ProgramConfig struct {
	// Default for any not matched is PackageNotInterpreted
	Packages []PackageFilter
	// Default is reflect.Value.MapRange
	MapRanger MapRanger
	// Default is not to use a custom scheduler
	Scheduler Scheduler
	// Value to use for vars or funcs (including methods) accessed. Only applies
	// if called by interpreted code. Key is qualified name (including receiver
	// for methods).
	InterceptInterpreted map[string]reflect.Value
}

type MapRanger func(reflect.Value) MapIter

type MapIter interface {
	Key() reflect.Value
	Value() reflect.Value
	Next() bool
}

// TODO(cretz): Surely will change how this looks as we work on it
type Scheduler interface {
	NewCoroutine() Coroutine
}

type Coroutine interface {
	Complete()
	Yielded() bool
	Select(cases []reflect.SelectCase) (chosen int, recv reflect.Value, recvOK bool)
	Send(ch reflect.Value, arg reflect.Value)
	Recv(ch reflect.Value) (x reflect.Value, ok bool)
}

func (p *Program) ReflectValue(pkgName, topLevelName string) reflect.Value {
	panic("TODO")
}

func (p *Program) ReflectFunc(fn interface{}) reflect.Value {
	panic("TODO")
}

func (p *Program) ReflectType(pkgName, topLevelName string) reflect.Type {
	panic("TODO")
}

func (p *Program) Run() {
	for p.RunOnce() {
	}
}

func (*Program) RunOnce() (coroutinesRemain bool) {
	panic("TODO")
}

func (*Program) ScheduleCall(fn reflect.Value, args []reflect.Value) <-chan []reflect.Value {
	panic("TODO")
}

type PackageFilter func(*packages.Package) PackageMode

type PackageMode int

const (
	PackageModeUnknown PackageMode = iota
	PackageInterpreted
	PackageNotInterpreted
	PackageDisallowed
)

func InterpretPackageByName(name string) PackageFilter {
	return func(p *packages.Package) PackageMode {
		if p.PkgPath == name {
			return PackageInterpreted
		}
		return PackageModeUnknown
	}
}
