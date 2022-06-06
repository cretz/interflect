package interflect

import "reflect"

type PackageReflector interface {
	ReflectType(pkgName, topLevelName string) reflect.Type
	ReflectValue(pkgName, topLevelName string) reflect.Value
}
