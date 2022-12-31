package inithook_test

import (
	"context"
	"testing"

	"github.com/ccmonky/pkg/inithook"
)

func TestInitHook(t *testing.T) {
	// e.g. in apierrors lib
	apierrorsLibAppName := ""
	err := inithook.RegisterAttrSetter(inithook.AppName, "apierrors", func(ctx context.Context, value string) error {
		apierrorsLibAppName = value
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	// e.g. in render lib
	renderLibAppName := ""
	err = inithook.RegisterAttrSetter(inithook.AppName, "render", func(ctx context.Context, value string) error {
		renderLibAppName = value
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	// e.g. another lib
	err = inithook.RegisterAttrSetter("attr_not_used", "another", func(ctx context.Context, value string) error {
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	// e.g. in app
	err = inithook.ExecuteAttrSetters(context.Background(), inithook.AppName, "xxx")
	if err != nil {
		t.Fatal(err)
	}
	if apierrorsLibAppName != "xxx" {
		t.Errorf("apierrorsLibAppName should == xxx, got %s", apierrorsLibAppName)
	}
	if renderLibAppName != "xxx" {
		t.Errorf("renderLibAppName should == xxx, got %s", renderLibAppName)
	}
	attrs := inithook.AttrsNotSetted()
	if len(attrs) != 1 {
		t.Error("should == 1")
	}
	if attrs[0] != "attr_not_used" {
		t.Error("attr_not_used should not be executed")
	}
}
