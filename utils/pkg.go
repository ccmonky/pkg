package utils

import (
	"fmt"
	"log"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"

	"golang.org/x/tools/go/packages"
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

// Source returns the `Source` caller's package path, maybe panic
// NOTE:
// 1. should only use in init phase for performance reason!
// 2. will not
func Source() string {
	_, f, _, ok := runtime.Caller(1)
	if !ok {
		log.Panicln("can not get current file path!")
	}
	path := filepath.Dir(filepath.Clean(f))
	source, ok := sourceCache.Load(path)
	if ok {
		return source.(string)
	}
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedModule,
	}, path)
	if err != nil {
		log.Panicln(err)
	}
	if len(pkgs) == 0 {
		log.Panicf("missing package information for path: %s", path)
	}
	sourceCache.Store(path, pkgs[0].PkgPath)
	return pkgs[0].PkgPath
}

var (
	sourceCache sync.Map // map[path]pkgPath
)
