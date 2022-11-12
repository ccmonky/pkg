package pkg_test

import (
	"testing"

	"github.com/ccmonky/pkg"
)

func TestUlid(t *testing.T) {
	id, err := pkg.Ulid()
	if err != nil {
		t.Fatal(err)
	}
	if len(id) != 26 {
		t.Fatal("should ==")
	}
	t.Log(id)
	//t.Fatal(1)
}
