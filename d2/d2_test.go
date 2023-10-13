package d2_test

import (
	"context"
	"github.com/ccmonky/pkg/d2"
	"os"

	"testing"
)

func TestRender(t *testing.T) {
	svg, err := d2.Render(context.Background(), "templates/arch.d2")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile("arch.svg", svg, 0644)
	if err != nil {
		t.Fatal(err)
	}
}
