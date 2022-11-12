package pkg

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// PkgPath return package path of value v
func PkgPath(v interface{}) string {
	typ, _ := Indirect(reflect.TypeOf(v))
	return typ.PkgPath()
}

// PkgPath return type name of value v
func TypeName(v interface{}) string {
	typ, _ := Indirect(reflect.TypeOf(v))
	return typ.Name()
}

// PkgPath return package and type name of value v
func PkgTypeName(v interface{}) string {
	typ, _ := Indirect(reflect.TypeOf(v))
	return fmt.Sprintf("%s.%s", typ.PkgPath(), typ.Name())
}

// PkgFuncName return package and func name of func fn
func PkgFuncName(fn interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
}

// FuncName return func name of func fn
func FuncName(fn interface{}) string {
	parts := strings.Split(PkgFuncName(fn), ".")
	return parts[len(parts)-1]
}

// InterfaceName return interface name of an interface value
func InterfaceName(i interface{}) string {
	return reflect.TypeOf(i).Elem().Name()
}

// Indirect returns the item at the end of Indirection.
func Indirect(t reflect.Type) (reflect.Type, bool) {
	var isPtr bool
	for ; t.Kind() == reflect.Ptr; t = t.Elem() {
		isPtr = true
	}
	return t, isPtr
}
