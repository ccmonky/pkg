package utils_test

import (
	"strings"
	"testing"

	"github.com/ccmonky/pkg/utils"
)

func TestIsSpaceOr(t *testing.T) {
	ss := strings.FieldsFunc("parking|  def", utils.IsSpaceOr('|'))
	if ss[1] != "def" {
		t.Fatal("should ==")
	}
}

func TestToLowerFirst(t *testing.T) {
	if utils.ToLowerFirst("ProjectMock") != "projectMock" {
		t.Fatal("should ==")
	}
}

func TestSplitByLength(t *testing.T) {
	s := "thissia testl hahafw"
	ss := utils.SplitByLength(s, 2)
	t.Log(ss)
	if ss[len(ss)-1] != "fw" {
		t.Fatal("should ==")
	}
	s2 := "thissia testl hahaf"
	ss2 := utils.SplitByLength(s2, 2)
	if ss2[len(ss2)-1] != "f" {
		t.Fatal("should ==", ss2)
	}
	ss3 := utils.SplitByLength(s, 100)
	if len(ss3) != 1 {
		t.Fatal("should ==")
	}
}

func TestRemoveDuplicatesUnordered(t *testing.T) {
	ss := []string{"test", "ami", "cns2", "audi", "test"}
	ss2 := utils.RemoveDuplicatesUnordered(ss)
	if len(ss2) != 4 {
		t.Fatal("should ==")
	}
}

func TestCompare(t *testing.T) {
	olds := []string{"1", "2", "3", "4"}
	news := []string{"2", "3", "4", "6"}
	adds, deletes := utils.Compare(olds, news)
	if !utils.Contains(adds, "6") {
		t.Fatal("should contains")
	}
	if !utils.Contains(deletes, "1") {
		t.Fatal("should contains")
	}
}

func TestIndex(t *testing.T) {
	ss := []string{"1", "2", "3", "4"}
	if utils.Index(ss, "3") != 2 {
		t.Fatal("should ==")
	}
	if utils.Index(ss, "5") != -1 {
		t.Fatal("should ==")
	}
}
