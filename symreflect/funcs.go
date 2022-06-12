package symreflect

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"unsafe"
)

//go:linkname modulesSlice runtime.modulesSlice
var modulesSlice *reflect.SliceHeader

type LoadFuncsOptions struct {
	// Required. Can be a LoadTypes()["runtime"]["moduledata"].
	ModuleDataType reflect.Type
}

type Funcs struct {
	Named map[string]map[string]*runtime.Func
}

func LoadFuncs(options LoadFuncsOptions) (*Funcs, error) {
	// Inspired by https://github.com/alangpierce/go-forceexport

	if modulesSlice == nil {
		return nil, fmt.Errorf("missing module slice")
	} else if options.ModuleDataType == nil {
		return nil, fmt.Errorf("missing module data type")
	}

	// Reflect the linked var
	modSliceVal := reflect.NewAt(reflect.SliceOf(reflect.PtrTo(options.ModuleDataType)), unsafe.Pointer(modulesSlice)).Elem()

	// Iterate to collect functions
	funcs := map[string]map[string]*runtime.Func{}
	for i := 0; i < modSliceVal.Len(); i++ {
		modDataVal := modSliceVal.Index(i).Elem()

		// Get the lookup table
		pclntable := getUnexportedField(modDataVal.FieldByName("pclntable")).([]byte)

		// Get the function table and loop over it
		ftabVal := modDataVal.FieldByName("ftab")
		for j := 0; j < ftabVal.Len(); j++ {
			// Get offset which may be uintptr or uint32 depending on Go version
			funcOffRaw := getUnexportedField(ftabVal.Index(j).FieldByName("funcoff"))
			off, ok := funcOffRaw.(uintptr)
			if !ok {
				off = uintptr(funcOffRaw.(uint32))
			}
			// Sometimes the offset is out of range, ignore it
			if int(off) >= len(pclntable) {
				continue
			}
			// Convert and separate package from func name
			fn := (*runtime.Func)(unsafe.Pointer(&pclntable[off]))
			fnName := fn.Name()
			if fnName == "" {
				// Ignore unknown functions
				continue
			}
			pkgPathEndIndex := strings.LastIndex(fnName, "/")
			pkgPathEndIndex += 1 + strings.Index(fnName[pkgPathEndIndex+1:], ".")
			if pkgPathEndIndex < 0 {
				continue
			}
			pkgPath, fnName := fnName[:pkgPathEndIndex], fnName[pkgPathEndIndex+1:]
			pkgFuncs := funcs[pkgPath]
			if pkgFuncs == nil {
				pkgFuncs = map[string]*runtime.Func{}
				funcs[pkgPath] = pkgFuncs
			}
			pkgFuncs[fnName] = fn
		}
	}
	return &Funcs{funcs}, nil
}

func ReflectFunc(typ reflect.Type, addr uintptr) reflect.Value {
	// Inspired by by https://github.com/kstenerud/go-subvert
	rFunc := reflect.MakeFunc(typ, nil)
	makeAddressable(&rFunc)
	pFunc := (*unsafe.Pointer)(unsafe.Pointer(rFunc.UnsafeAddr()))
	*pFunc = unsafe.Pointer(addr)
	return rFunc
}

func getUnexportedField(field reflect.Value) interface{} {
	return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Interface()
}
