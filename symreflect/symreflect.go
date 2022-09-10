package symreflect

import (
	"reflect"
	"sync"
)

type reflector struct {
	types   *Types
	symbols *Symbols

	cachedReflectValues     map[string]*reflect.Value
	cachedReflectValuesLock sync.RWMutex
}

var reflectorInst *reflector
var reflectorInitLock sync.Mutex

// TODO(cretz): Does not work during tests due to lack of symbol table
func Reflector() (*reflector, error) {
	reflectorInitLock.Lock()
	defer reflectorInitLock.Unlock()
	if reflectorInst == nil {
		types, err := LoadTypes(LoadTypesOptions{})
		if err != nil {
			return nil, err
		}
		symbols, err := LoadSymbols(LoadSymbolsOptions{})
		if err != nil {
			return nil, err
		}
		reflectorInst = &reflector{
			types:               types,
			symbols:             symbols,
			cachedReflectValues: map[string]*reflect.Value{},
		}
	}
	return reflectorInst, nil
}

func (r *reflector) ReflectType(pkgName, topLevelName string) reflect.Type {
	return r.types.Named[pkgName][topLevelName]
}

func (r *reflector) ReflectValue(pkgName, topLevelName string) reflect.Value {
	symbol := r.symbols.Named[pkgName][topLevelName]
	if symbol == nil {
		return reflect.Value{}
	}
	// Get cached under read lock, but if not there create unlocked and update
	// under write lock (even though it means potential overwrites which is safe)
	fullName := pkgName + "//" + topLevelName
	r.cachedReflectValuesLock.RLock()
	val := r.cachedReflectValues[fullName]
	r.cachedReflectValuesLock.RUnlock()
	if val == nil {
		toCache := r.symbols.ReflectValue(r.ReflectType(pkgName, topLevelName), symbol)
		val = &toCache
		r.cachedReflectValuesLock.Lock()
		r.cachedReflectValues[fullName] = val
		r.cachedReflectValuesLock.Unlock()
	}
	return *val
}
