package interflect

import (
	"fmt"
	"go/types"
	"reflect"
	"unsafe"
)

func (e *Execution) reflectType(t types.Type) reflect.Type {
	switch t := t.(type) {
	case *types.Array:
		return reflect.ArrayOf(int(t.Len()), e.reflectType(t.Elem()))
	case *types.Basic:
		switch t.Kind() {
		case types.Bool:
			return reflect.TypeOf(false)
		case types.Int:
			return reflect.TypeOf(int(0))
		case types.Int8:
			return reflect.TypeOf(int8(0))
		case types.Int16:
			return reflect.TypeOf(int16(0))
		case types.Int32:
			return reflect.TypeOf(int32(0))
		case types.Int64:
			return reflect.TypeOf(int64(0))
		case types.Uint:
			return reflect.TypeOf(uint(0))
		case types.Uint8:
			return reflect.TypeOf(uint8(0))
		case types.Uint16:
			return reflect.TypeOf(uint16(0))
		case types.Uint32:
			return reflect.TypeOf(uint32(0))
		case types.Uint64:
			return reflect.TypeOf(uint64(0))
		case types.Uintptr:
			return reflect.TypeOf(uintptr(0))
		case types.Float32:
			return reflect.TypeOf(float32(0))
		case types.Float64:
			return reflect.TypeOf(float64(0))
		case types.Complex64:
			return reflect.TypeOf(complex64(0))
		case types.Complex128:
			return reflect.TypeOf(complex128(0))
		case types.String:
			return reflect.TypeOf("")
		case types.UnsafePointer:
			return reflect.TypeOf(unsafe.Pointer(nil))
		default:
			panic(fmt.Sprintf("unable to get type for basic kind %v", t.Kind()))
		}
	case *types.Chan:
		switch t.Dir() {
		case types.SendOnly:
			return reflect.ChanOf(reflect.SendDir, e.reflectType(t.Elem()))
		case types.RecvOnly:
			return reflect.ChanOf(reflect.RecvDir, e.reflectType(t.Elem()))
		default:
			return reflect.ChanOf(reflect.BothDir, e.reflectType(t.Elem()))
		}
	case *types.Interface:
		panic(fmt.Sprintf("anonymous interface not supported"))

	}
	// TODO(cretz): Finish
	panic("TODO")
}
