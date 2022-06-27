package symreflect

import (
	"reflect"
	"unsafe"
)

//go:linkname resolveTypeOff reflect.resolveTypeOff
func resolveTypeOff(rtype unsafe.Pointer, off int32) unsafe.Pointer

//go:linkname typelinks reflect.typelinks
func typelinks() (sections []unsafe.Pointer, offset [][]int32)

type emptyInterface struct {
	typ  unsafe.Pointer
	word unsafe.Pointer
}

// LoadTypesOptions are options for LoadTypes.
type LoadTypesOptions struct {
}

// Types are types loaded from LoadTypes.
type Types struct {
	Named       map[string]map[string]reflect.Type
	AnonStructs []reflect.Type
}

// LoadTypes loads all known types from the internal typelinks in the
// executable.
func LoadTypes(LoadTypesOptions) (*Types, error) {
	// Inspired by https://github.com/modern-go/reflect2

	types := &Types{Named: map[string]map[string]reflect.Type{}}
	var obj interface{} = reflect.TypeOf(0)
	sections, offset := typelinks()
	for i, offs := range offset {
		rodata := sections[i]
		for _, off := range offs {
			(*emptyInterface)(unsafe.Pointer(&obj)).word = resolveTypeOff(unsafe.Pointer(rodata), off)
			typ := obj.(reflect.Type)
			// Only pointers are named
			if typ.Kind() == reflect.Ptr {
				typ = typ.Elem()
				if typ.PkgPath() != "" && typ.Name() != "" {
					pkgTypes := types.Named[typ.PkgPath()]
					if pkgTypes == nil {
						pkgTypes = map[string]reflect.Type{}
						types.Named[typ.PkgPath()] = pkgTypes
					}
					pkgTypes[typ.Name()] = typ
				}
			} else if typ.Kind() == reflect.Struct {
				types.AnonStructs = append(types.AnonStructs, typ)
			}
		}
	}
	return types, nil
}
