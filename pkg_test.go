package pkg_test

import (
	"net/http"
	"testing"

	"github.com/ccmonky/pkg"
	"github.com/stretchr/testify/assert"
)

type resourceType struct{}

type resourceFunc func()

type resourceRetype string

type resourceAlias = string

func TestPkgPath(t *testing.T) {
	fn := resourceFunc(func() {})
	ret := resourceRetype("")
	pret := &ret
	ppret := &pret
	builtin := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
	var cases = []struct {
		object   interface{}
		pkgPath  string
		typeName string
	}{
		{
			resourceType{},
			pkgPath,
			"resourceType",
		},
		{
			fn,
			pkgPath,
			"resourceFunc",
		},
		{
			ret,
			pkgPath,
			"resourceRetype",
		},
		{
			&resourceType{},
			pkgPath,
			"resourceType",
		},
		{
			&fn,
			pkgPath,
			"resourceFunc",
		},
		{
			&ret,
			pkgPath,
			"resourceRetype",
		},
		{
			builtin,
			"net/http",
			"HandlerFunc",
		},
		{
			&builtin,
			"net/http",
			"HandlerFunc",
		},
		{
			ppret,
			pkgPath,
			"resourceRetype",
		},
	}
	for _, tc := range cases {
		path := pkg.PkgPath(tc.object)
		name := pkg.TypeName(tc.object)
		if path != tc.pkgPath {
			t.Fatalf("%v should ==, got %s", tc, path)
		}
		if name != tc.typeName {
			t.Fatalf("%v should ==, got %s", tc, name)
		}
	}

	alias := resourceAlias("")
	path := pkg.PkgPath(alias)
	if path != "" {
		t.Fatalf("alias pkg path should ==, got %s", path)
	}
	name := pkg.TypeName(alias)
	if name != "string" {
		t.Fatalf("alias type name should ==, got %s", name)
	}
	path = pkg.PkgPath(&alias)
	if path != "" {
		t.Fatalf("alias pkg path should ==, got %s", path)
	}
	name = pkg.TypeName(&alias)
	if name != "string" {
		t.Fatalf("alias type name should ==, got %s", name)
	}
}

func abc() {}

func TestFuncName(t *testing.T) {
	assert.Equalf(t, "abc", pkg.FuncName(abc), "func name of abc")
	assert.Equalf(t, "TypeName", pkg.FuncName(pkg.TypeName), "func name of pkg.TypeName")
	assert.Equalf(t, pkgPath+".abc", pkg.PkgFuncName(abc), "func name of pkg.TypeName")
	assert.Equalf(t, "github.com/ccmonky/pkg.TypeName", pkg.PkgFuncName(pkg.TypeName), "func name of pkg.TypeName")
}

type iface interface {
	xxx()
}

func TestInterfaceName(t *testing.T) {
	assert.Equalf(t, "iface", pkg.InterfaceName((*iface)(nil)), "pkg interface name")
}

var pkgPath = "github.com/ccmonky/pkg_test"
