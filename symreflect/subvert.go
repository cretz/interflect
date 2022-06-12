package symreflect

import (
	"reflect"
	"unsafe"
)

// Most of this file is taken from https://github.com/kstenerud/go-subvert with
// slight mods

func makeAddressable(v *reflect.Value) {
	*getRVFlagPtr(v) |= rvFlagAddr
}

func getRVFlagPtr(v *reflect.Value) *uintptr {
	return (*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(v)) + rvFlagOffset))
}

var (
	rvFlagAddr uintptr
	rvFlagRO   uintptr

	rvFlagOffset uintptr
)

type rvFlagTester struct {
	A   int // reflect/value.go: flagAddr
	a   int // reflect/value.go: flagStickyRO
	int     // reflect/value.go: flagEmbedRO
	// Note: flagRO = flagStickyRO | flagEmbedRO as of go 1.5
}

func init() {
	getFlag := func(v reflect.Value) uintptr {
		return uintptr(reflect.ValueOf(v).FieldByName("flag").Uint())
	}
	getFldFlag := func(v reflect.Value, fieldName string) uintptr {
		return getFlag(v.FieldByName(fieldName))
	}

	if field, ok := reflect.TypeOf(reflect.Value{}).FieldByName("flag"); ok {
		rvFlagOffset = field.Offset
	} else {
		panic("reflect.Value no longer has a flag field")
	}

	v := rvFlagTester{}
	rv := reflect.ValueOf(&v).Elem()
	rvFlagRO = (getFldFlag(rv, "a") | getFldFlag(rv, "int")) ^ getFldFlag(rv, "A")
	if rvFlagRO == 0 {
		panic("reflect.Value.flag no longer has flagEmbedRO or flagStickyRO bit")
	}

	rvFlagAddr = getFlag(reflect.ValueOf(int(1))) ^ getFldFlag(rv, "A")
	if rvFlagAddr == 0 {
		panic("reflect.Value.flag no longer has a flagAddr bit")
	}
}
