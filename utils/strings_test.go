package pkg_test

import (
	"strings"
	"testing"

	"github.com/ccmonky/pkg"
)

func TestIsSpaceOr(t *testing.T) {
	ss := strings.FieldsFunc("parking|  def", pkg.IsSpaceOr('|'))
	if ss[1] != "def" {
		t.Fatal("should ==")
	}
}

func TestToLowerFirst(t *testing.T) {
	if pkg.ToLowerFirst("ProjectMock") != "projectMock" {
		t.Fatal("should ==")
	}
}

func TestSplitByLength(t *testing.T) {
	s := "thissia testl hahafw"
	ss := pkg.SplitByLength(s, 2)
	t.Log(ss)
	if ss[len(ss)-1] != "fw" {
		t.Fatal("should ==")
	}
	s2 := "thissia testl hahaf"
	ss2 := pkg.SplitByLength(s2, 2)
	if ss2[len(ss2)-1] != "f" {
		t.Fatal("should ==", ss2)
	}
	ss3 := pkg.SplitByLength(s, 100)
	if len(ss3) != 1 {
		t.Fatal("should ==")
	}
}

func TestRemoveDuplicatesUnordered(t *testing.T) {
	ss := []string{"test", "ami", "cns2", "audi", "test"}
	ss2 := pkg.RemoveDuplicatesUnordered(ss)
	if len(ss2) != 4 {
		t.Fatal("should ==")
	}
}

func TestCompare(t *testing.T) {
	olds := []string{"1", "2", "3", "4"}
	news := []string{"2", "3", "4", "6"}
	adds, deletes := pkg.Compare(olds, news)
	if !pkg.Contains(adds, "6") {
		t.Fatal("should contains")
	}
	if !pkg.Contains(deletes, "1") {
		t.Fatal("should contains")
	}
}

func TestIndex(t *testing.T) {
	ss := []string{"1", "2", "3", "4"}
	if pkg.Index(ss, "3") != 2 {
		t.Fatal("should ==")
	}
	if pkg.Index(ss, "5") != -1 {
		t.Fatal("should ==")
	}
}
