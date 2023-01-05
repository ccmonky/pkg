package utils_test

import (
	"testing"

	"github.com/ccmonky/pkg/utils"
)

func TestUlid(t *testing.T) {
	id, err := utils.Ulid()
	if err != nil {
		t.Fatal(err)
	}
	if len(id) != 26 {
		t.Fatal("should ==")
	}
	t.Log(id)
	//t.Fatal(1)
}
