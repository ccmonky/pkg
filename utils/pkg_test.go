package utils_test

import (
	"net/http"
	"testing"

	"github.com/ccmonky/pkg/utils"
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
		path := utils.PkgPath(tc.object)
		name := utils.TypeName(tc.object)
		if path != tc.pkgPath {
			t.Fatalf("%v should ==, got %s", tc, path)
		}
		if name != tc.typeName {
			t.Fatalf("%v should ==, got %s", tc, name)
		}
	}

	alias := resourceAlias("")
	path := utils.PkgPath(alias)
	if path != "" {
		t.Fatalf("alias pkg path should ==, got %s", path)
	}
	name := utils.TypeName(alias)
	if name != "string" {
		t.Fatalf("alias type name should ==, got %s", name)
	}
	path = utils.PkgPath(&alias)
	if path != "" {
		t.Fatalf("alias pkg path should ==, got %s", path)
	}
	name = utils.TypeName(&alias)
	if name != "string" {
		t.Fatalf("alias type name should ==, got %s", name)
	}
}

func abc() {}

func TestFuncName(t *testing.T) {
	assert.Equalf(t, "abc", utils.FuncName(abc), "func name of abc")
	assert.Equalf(t, "TypeName", utils.FuncName(utils.TypeName), "func name of utils.TypeName")
	assert.Equalf(t, pkgPath+".abc", utils.PkgFuncName(abc), "func name of utils.TypeName")
	assert.Equalf(t, "github.com/ccmonky/utils.TypeName", utils.PkgFuncName(utils.TypeName), "func name of utils.TypeName")
}

type iface interface {
	xxx()
}

func TestInterfaceName(t *testing.T) {
	assert.Equalf(t, "iface", utils.InterfaceName((*iface)(nil)), "pkg interface name")
}

var pkgPath = "github.com/ccmonky/pkg_test"
